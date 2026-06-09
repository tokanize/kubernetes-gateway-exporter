package metrics

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tokanize/kubernetes-gateway-exporter/internal/kubernetes/mapper"
)

var (
	exposedRouteInfoDesc = prometheus.NewDesc(
		"exposed_route_info",
		"Information about logical route-to-Service relationships exposed via Kubernetes Gateway API.",
		[]string{"namespace", "gateway_name", "route_name", "hostname", "listener_name", "http_path", "service_name", "backend_namespace", "backend_port", "lb_type", "ip_address"},
		nil,
	)
)

// Exporter is a Prometheus Collector that gathers exposed route metrics.
type Exporter struct {
	mapper *mapper.Mapper
	logger *slog.Logger
}

// NewExporter creates a new Prometheus Collector.
func NewExporter(m *mapper.Mapper, logger *slog.Logger) *Exporter {
	return &Exporter{
		mapper: m,
		logger: logger.With(slog.String("component", "prometheus_collector")),
	}
}

// Describe implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- exposedRouteInfoDesc
}

// Collect implements prometheus.Collector. It is called on every scrape.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	// Scrape context with concrete timeout
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	routes, err := e.mapper.GetExposedRoutes(ctx)
	if err != nil {
		e.logger.Error("Failed to collect exposed routes", slog.Any("error", err))
		return
	}

	for _, route := range routes {
		ch <- prometheus.MustNewConstMetric(
			exposedRouteInfoDesc,
			prometheus.GaugeValue,
			1.0, // Value is always 1 for info metrics
			route.Namespace,
			route.GatewayName,
			route.RouteName,
			route.Hostname,
			route.ListenerName,
			route.HTTPPath,
			route.ServiceName,
			route.BackendNamespace,
			route.BackendPort,
			route.LBType,
			route.IPAddress,
		)
	}
	e.logger.Debug("Prometheus metrics collected", slog.Int("routes_count", len(routes)), slog.Duration("duration", time.Since(start)))
}
