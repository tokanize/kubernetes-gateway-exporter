# Architecture

The **Kubernetes Gateway Exporter** publishes an inventory of the routing
relationships your Gateway API resources define. For every accepted route it
answers one question — *what is exposed, where, and how?* — and emits the answer
as the `exposed_route_info` metric for Prometheus and OpenTelemetry.

## How it works

The exporter watches Gateway API resources through cache-backed informers and
resolves the relationship chain:

`Gateway → Listener → HTTPRoute → Service`

<img src="./gateway-architecture.svg"
     alt="Gateway to Listener to HTTPRoute to Service, resolved from the informer cache into the exposed_route_info metric"
     style="max-width: 720px; width: 100%; margin: 1.5rem 0;" />

It maps the **logical** routing graph only. It does not enumerate Pods, Pod IPs,
Endpoints, or EndpointSlices behind a Service — a Service with five backing Pods
still produces one series per valid listener / hostname / path / backend
combination, not five Pod-level series.

For the complete label set and configuration, see the
[Metrics reference](./metrics.md).

## Documentation strategy

This repository uses a hybrid documentation approach:

- **OpenSpec** (`openspec/specs/`) defines the functional requirements as
  GIVEN / WHEN / THEN scenarios.
- **ADRs** (`docs/adrs/`) capture the technical architecture and the reasoning
  behind each decision.

The functional specs link to their corresponding technical ADRs, so every core
component in the Go source traces back to a recorded decision.

## Architecture Decision Records

- [ADR-001: Kubernetes Informers Strategy](./adrs/001-k8s-informers.md) — cache-backed informers for resource mapping.
- [ADR-002: Gateway Address Resolution](./adrs/002-ip-resolution.md) — portable address extraction from `Gateway.status.addresses`.
- [ADR-003: Testing Strategy](./adrs/003-testing-strategy.md) — table-driven mapping tests and Prometheus contract tests.
- [ADR-004: OpenTelemetry Metrics](./adrs/004-otel-metrics.md) — the OTLP push model and Observable Gauge implementation.
- [ADR-005: Security Hardening](./adrs/005-security.md) — HTTP timeouts, distroless containers, and runtime hardening.
- [ADR-006: Exposure Metric Identity](./adrs/006-metric-labels-expansion.md) — the flattened Gateway API relationship and its metric labels.

## API specification

The [OpenAPI 3.0 specification](https://github.com/tokanize/kubernetes-gateway-exporter/blob/main/docs/api/openapi.yaml)
defines the management and health endpoints exposed by the exporter
(`/healthz`, `/readyz`, `/metrics`).

## Further reading

- [Getting Started](./getting-started.md) — install with Helm and confirm the metric is served.
- [Testing on kind](./testing-on-kind.md) — verify the exporter without a Gateway controller.
- [OpenSpec requirements](https://github.com/tokanize/kubernetes-gateway-exporter/tree/main/openspec/specs) — functional specs that drive the implementation.
