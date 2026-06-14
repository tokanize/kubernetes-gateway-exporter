# metrics-exposition Specification

## Purpose
Define the format and transport mechanism for exposing the mapped Kubernetes route data to observability systems (Prometheus & OpenTelemetry).

> **Implementation Details:** See [ADR-004: OpenTelemetry Metrics](../../../docs/adrs/004-otel-metrics.md) for the technical architecture of the OTLP Observable Gauge.

## Requirements

### Requirement: Prometheus Exporter
The system SHALL provide a `/metrics` HTTP endpoint compliant with Prometheus scraping standards.

#### Scenario: Scrape metrics
- GIVEN the application is running
- WHEN a GET request is made to `/metrics`
- THEN the system returns the `exposed_route_info` gauge
- AND the gauge includes labels: `namespace`, `gateway_name`, `route_name`, `hostname`, `listener_name`, `http_path`, `service_name`, `backend_namespace`, `backend_port`, `lb_type`, and `ip_address`.

#### Scenario: Service has multiple backing Pods
- GIVEN an emitted route relationship targets a Service with multiple Pods or EndpointSlice addresses
- WHEN metrics are collected
- THEN the exporter emits relationship series based on Gateway API configuration
- AND it does not emit one series per Pod or EndpointSlice address.

### Requirement: OpenTelemetry Push
The system SHALL support pushing metrics directly to an OpenTelemetry collector using the OTLP gRPC protocol.

#### Scenario: OTEL metrics push is disabled
- GIVEN `ENABLE_OTEL=false` or unset
- WHEN the system starts
- THEN the OTEL provider is not initialized
- AND no metrics are pushed.

#### Scenario: OTEL metrics push is enabled
- GIVEN `ENABLE_OTEL=true`
- WHEN the system starts
- THEN it registers an `Int64ObservableGauge` with the OTEL SDK
- AND the OTEL SDK periodically invokes a callback
- AND the callback retrieves the current route data from memory
- AND the data is pushed to `localhost:4317` (default) via gRPC.
