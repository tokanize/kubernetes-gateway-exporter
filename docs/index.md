# Kubernetes Gateway Exporter Documentation

Welcome to the internal documentation for the **Kubernetes Gateway Exporter**. This repository provides an inventory of logical routing relationships exposed through Kubernetes Gateway API resources.

## Core Traceability Map

This project adheres strictly to the "Docs-First" rule. No core application logic is written without architectural justification. Every core component in the Go source code links back to the documents listed here.

> **Note on Documentation Strategy:** This repository uses a hybrid documentation approach. 
> - **[OpenSpec (`openspec/specs/`)](../openspec/specs/)** defines the functional requirements.
> - **[ADRs (`docs/adrs/`)](./adrs/)** define the technical architecture and decisions.
> The functional specs actively link to their corresponding technical ADRs.

### Architecture Decision Records (ADRs)

*   [ADR-001: Kubernetes Informers Strategy](./adrs/001-k8s-informers.md) - Details the choice of using cache-backed informers for resource mapping.
*   [ADR-002: Gateway Address Resolution](./adrs/002-ip-resolution.md) - Explains portable address extraction from `Gateway.status.addresses`.
*   [ADR-003: Testing Strategy](./adrs/003-testing-strategy.md) - Defines table-driven mapping tests and Prometheus contract tests.
*   [ADR-004: OpenTelemetry Strategy](./adrs/004-otel-metrics.md) - Explains the OTLP push model and Observable Gauge implementation for OTEL.
*   [ADR-005: Security Hardening](./adrs/005-security.md) - Documents HTTP timeouts (Slowloris protection), distroless containers, and future recommendations.
*   [ADR-006: Exposure Metric Identity](./adrs/006-metric-labels-expansion.md) - Defines the flattened Gateway API relationship and metric labels.

### API Specifications

*   [OpenAPI 3.0 Specification](./api/openapi.yaml) - Defines the internal management and health endpoints exposed by the exporter (`/healthz`, `/readyz`, `/metrics`).

### Operational Guides

*   [Testing locally on kind](./testing-on-kind.md) - Builds and verifies the exporter without installing a Gateway controller.

### Key Domain Relationships

The core logic of this exporter resolves the following relationship chain using Kubernetes informers:
`Gateway` -> `Listener` -> `HTTPRoute` -> `Service`

The metric represents logical routing relationships. It does not enumerate Pods, Pod IPs, Endpoints, or EndpointSlices behind a Service.

Labels exposed by the resulting Prometheus metric (`exposed_route_info`):
*   `namespace`
*   `gateway_name`
*   `route_name`
*   `route_namespace`
*   `hostname`
*   `listener_name`
*   `http_path`
*   `service_name`
*   `backend_namespace`
*   `backend_port`
*   `lb_type` (`internal` or `external`, heuristically inferred from the GatewayClass name)
*   `ip_address` (read from `Gateway.status.addresses`)
