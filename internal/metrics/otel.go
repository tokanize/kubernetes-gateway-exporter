package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tokanize/kubernetes-gateway-exporter/internal/kubernetes/mapper"
	"github.com/tokanize/kubernetes-gateway-exporter/pkg/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// SetupOTEL initializes the OpenTelemetry push exporter and registers the async callback.
// See docs/adrs/004-otel-metrics.md
func SetupOTEL(ctx context.Context, m *mapper.Mapper, logger *slog.Logger) (func(context.Context) error, error) {
	// Initialize OTLP gRPC Exporter
	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("kubernetes-gateway-exporter"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTEL resource: %w", err)
	}

	// Initialize MeterProvider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(15*time.Second))),
		sdkmetric.WithResource(res),
	)

	// Set global MeterProvider
	otel.SetMeterProvider(provider)

	// Create a Meter
	meter := provider.Meter("kubernetes-gateway-exporter")

	// Register the Observable Gauge
	_, err = meter.Int64ObservableGauge(
		"exposed_route_info",
		metric.WithDescription("Information about logical route-to-Service relationships exposed via Kubernetes Gateway API."),
		metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
			start := time.Now()

			type result struct {
				routes []models.ExposedRoute
				err    error
			}
			resCh := make(chan result, 1)

			go func() {
				routes, err := m.GetExposedRoutes(ctx)
				resCh <- result{routes: routes, err: err}
			}()

			var routes []models.ExposedRoute
			select {
			case <-ctx.Done():
				logger.Error("OTEL Callback Error: timeout waiting for exposed routes cache", slog.Any("error", ctx.Err()))
				return ctx.Err()
			case res := <-resCh:
				if res.err != nil {
					logger.Error("OTEL Callback Error: failed to get routes", slog.Any("error", res.err))
					return res.err
				}
				routes = res.routes
			}

			for _, route := range routes {
				o.Observe(1, metric.WithAttributes(
					attribute.String("namespace", route.Namespace),
					attribute.String("gateway_name", route.GatewayName),
					attribute.String("route_name", route.RouteName),
					attribute.String("hostname", route.Hostname),
					attribute.String("listener_name", route.ListenerName),
					attribute.String("http_path", route.HTTPPath),
					attribute.String("service_name", route.ServiceName),
					attribute.String("backend_namespace", route.BackendNamespace),
					attribute.String("backend_port", route.BackendPort),
					attribute.String("lb_type", route.LBType),
					attribute.String("ip_address", route.IPAddress),
				))
			}
			logger.Debug("OTEL metrics collected", slog.Int("routes_count", len(routes)), slog.Duration("duration", time.Since(start)))
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create observable gauge: %w", err)
	}

	return provider.Shutdown, nil
}
