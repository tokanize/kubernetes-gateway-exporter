package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tokanize/kubernetes-gateway-exporter/internal/middleware"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Setup creates and returns a fully configured http.Server with all application routes.
func Setup(port string, logger *slog.Logger, mgr manager.Manager) *http.Server {
	mux := http.NewServeMux()

	// Expose Prometheus Metrics
	mux.Handle("/metrics", promhttp.Handler())

	// Liveness Probe
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	// Readiness Probe
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		// Do not report ready until the controller-runtime cache (informers) have fully synced.
		if mgr != nil && !mgr.GetCache().WaitForCacheSync(r.Context()) {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "caches not synced")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ready")
	})

	// Wrap mux with standard logging middleware
	handler := middleware.LoggingMiddleware(logger, mux)

	// Configure hardened HTTP server (Slowloris protection)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB limit on headers
	}

	return srv
}

// RunAsync starts the HTTP server in a goroutine and logs the status.
func RunAsync(srv *http.Server, logger *slog.Logger) {
	go func() {
		logger.Info("Starting HTTP server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", slog.Any("error", err))
		}
	}()
}

// ShutdownGracefully blocks until the server successfully shuts down or contexts expire.
func ShutdownGracefully(srv *http.Server, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown error", slog.Any("error", err))
	} else {
		logger.Info("HTTP server gracefully stopped")
	}
}
