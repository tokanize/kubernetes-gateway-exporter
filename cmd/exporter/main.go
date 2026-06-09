package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tokanize/kubernetes-gateway-exporter/internal/kubernetes/mapper"
	"github.com/tokanize/kubernetes-gateway-exporter/internal/metrics"
	"github.com/tokanize/kubernetes-gateway-exporter/internal/server"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(gwv1.Install(scheme))
}

func main() {
	// Initialize JSON Logger for production
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting Kubernetes Gateway Exporter...")

	// 1. Setup Controller Runtime Manager
	config := ctrl.GetConfigOrDie()
	mgr, err := manager.New(config, manager.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
	})
	if err != nil {
		logger.Error("Failed to create manager", slog.Any("error", err))
		os.Exit(1)
	}

	// 2. Setup Mapper
	routeMapper := mapper.NewMapper(mgr.GetClient(), logger)
	if err := mgr.Add(routeMapper); err != nil {
		logger.Error("Failed to add mapper as runnable to manager", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. Setup Metrics Exporters
	// Prometheus
	promExporter := metrics.NewExporter(routeMapper, logger)
	prometheus.MustRegister(promExporter)

	// OpenTelemetry
	if os.Getenv("ENABLE_OTEL") == "true" {
		shutdownOTEL, err := metrics.SetupOTEL(ctx, routeMapper, logger)
		if err != nil {
			logger.Error("Failed to setup OTEL", slog.Any("error", err))
			os.Exit(1)
		}
		defer func() {
			if err := shutdownOTEL(context.Background()); err != nil {
				logger.Error("Failed to shutdown OTEL", slog.Any("error", err))
			}
		}()
		logger.Info("OpenTelemetry push metrics enabled")
	}

	// 4. Setup HTTP Endpoints & Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.Setup(port, logger, mgr, routeMapper)
	server.RunAsync(srv, logger)

	// Start the manager (starts the informers/caches)
	go func() {
		logger.Info("Starting informer caches...")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			logger.Error("Manager exited non-zero", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown setup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	logger.Info("Received shutdown signal", slog.String("signal", sig.String()))

	server.ShutdownGracefully(srv, logger)
}
