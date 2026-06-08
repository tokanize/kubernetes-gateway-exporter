# ADR-004: OpenTelemetry Metrics

**Date:** 2026-06-06  
**Status:** Accepted

## Context
While Prometheus is a standard for metric scraping, OpenTelemetry (OTEL) is becoming the industry standard for pushing metrics, traces, and logs. We need to support OTEL to allow sending `exposed_route_info` metrics via an OTLP pipeline.

## Decision
We will add support for pushing metrics via an OTLP gRPC exporter.

1. **Activation:** OTEL metrics will be enabled via the `ENABLE_OTEL=true` environment variable. By default, the exporter uses `localhost:4317` for the OTEL Collector endpoint (configurable via standard `OTEL_EXPORTER_OTLP_ENDPOINT` env var).
2. **Mechanism:** Unlike Prometheus (which scrapes an HTTP endpoint), OTEL SDKs support an asynchronous observation model. We will register an `Int64ObservableGauge` named `exposed_route_info`. 
3. **Execution:** The OTEL SDK triggers a callback every 15 seconds. The callback invokes `mapper.GetExposedRoutes()` and observes the same 12 attributes documented for the Prometheus metric.

## Consequences
- **Positive:** We can export metrics directly to observability backends (Datadog, Dynatrace, New Relic) via an OTEL Collector without exposing a scrape endpoint if we choose not to.
- **Negative:** OTEL dependencies add weight to the Go module. `GetExposedRoutes()` is called by both Prometheus scrapes and OTEL collection cycles, with work proportional to the cached route graph.

## Traceability
Code implementing this strategy MUST reference this document. Expected locations: `internal/metrics/otel.go`.
