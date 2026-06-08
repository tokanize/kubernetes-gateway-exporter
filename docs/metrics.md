# Metrics & Labels

The exporter emits a single metric, `exposed_route_info`, to Prometheus and —
optionally — to OpenTelemetry. It is a gauge whose value is always `1`; the
information lives entirely in the labels.

## `exposed_route_info`

Each series represents one logical exposure relationship: a specific HTTPRoute
attached to a listener and resolved to a backend Service.

| Label | Description |
|-------|-------------|
| `namespace` | HTTPRoute namespace, retained for backward compatibility. |
| `gateway_name` | Referenced parent Gateway name. |
| `route_name` | HTTPRoute name (`metadata.name`). |
| `route_namespace` | Namespace containing the HTTPRoute. |
| `hostname` | Effective hostname intersection between the route and listener; empty when unrestricted. |
| `listener_name` | Accepting Gateway listener. |
| `http_path` | Path from an HTTPRoute match, defaulting to `/`. |
| `service_name` | Backend Service name. |
| `backend_namespace` | Backend namespace, or the HTTPRoute namespace by default. |
| `backend_port` | Backend Service port, or an empty string when omitted. |
| `lb_type` | `internal` or `external`, inferred from the GatewayClass name. |
| `ip_address` | First `IPAddress` from `Gateway.status.addresses`, or an empty string while unavailable. |

### Example series

```text
exposed_route_info{backend_namespace="test-ns",backend_port="8080",gateway_name="my-gateway",hostname="",http_path="/api",ip_address="",lb_type="external",listener_name="http",namespace="test-ns",route_name="my-route",route_namespace="test-ns",service_name="test-svc"} 1
```

### What it does not measure

The metric describes logical routing relationships, not data-plane traffic. It
does not enumerate Pods, Pod IPs, Endpoints, or EndpointSlices behind a Service,
and backend replica count does not multiply the series.

### Cardinality

Approximate cardinality per Gateway/Route pairing:

`accepted listeners × effective hostnames × path matches × valid backendRefs`

See [ADR-006: Exposure Metric Identity](./adrs/006-metric-labels-expansion.md)
for the full identity and attachment rules.

## Configuration

The exporter is configured through environment variables (surfaced as Helm
values where noted).

| Variable | Default | Helm value | Description |
|----------|---------|------------|-------------|
| `ENABLE_OTEL` | `false` | `config.enableOtel` | Enable the OTLP/gRPC push exporter in addition to the Prometheus endpoint. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `localhost:4317` | — | OpenTelemetry Collector endpoint, used when `ENABLE_OTEL=true`. |

When OTEL is enabled, the SDK observes the same labels via an
`Int64ObservableGauge` on a periodic callback. See
[ADR-004: OpenTelemetry Metrics](./adrs/004-otel-metrics.md).

## Endpoints

| Path | Purpose |
|------|---------|
| `/metrics` | Prometheus scrape endpoint. |
| `/healthz` | Liveness probe. |
| `/readyz` | Readiness — returns `200` only after the informer cache has synced. |

The full endpoint schema is in the
[OpenAPI specification](https://github.com/tokanize/kubernetes-gateway-exporter/blob/main/docs/api/openapi.yaml).
