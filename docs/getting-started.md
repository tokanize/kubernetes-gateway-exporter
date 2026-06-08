# Getting Started

This guide installs the Kubernetes Gateway Exporter into a cluster that already
has the Gateway API CRDs, and confirms the `exposed_route_info` metric is being
served. To try it end-to-end on a throwaway cluster without a Gateway
controller, follow [Testing on kind](./testing-on-kind.md) instead.

## Prerequisites

- A Kubernetes cluster with the [Gateway API CRDs](https://gateway-api.sigs.k8s.io/guides/) installed
- `kubectl` and `helm` (v3.8+ for OCI registry support)
- Optional: the Prometheus Operator, if you want a `ServiceMonitor` created automatically

## Install with Helm

The chart is published as a signed OCI artifact. Install the latest release:

```bash
VERSION=0.1.2 # Replace with the release you want to install

helm upgrade --install gateway-exporter \
  oci://ghcr.io/tokanize/charts/kubernetes-gateway-exporter --version "${VERSION}" \
  --namespace monitoring --create-namespace
```

To verify the chart signature and build provenance before installing, see
[Release Verification](./verification.md).

### Without the Prometheus Operator

If the `ServiceMonitor` CRD is not installed, disable its creation:

```bash
helm upgrade --install gateway-exporter \
  oci://ghcr.io/tokanize/charts/kubernetes-gateway-exporter --version "${VERSION}" \
  --namespace monitoring --create-namespace \
  --set serviceMonitor.enabled=false
```

## Confirm it is running

```bash
kubectl rollout status \
  deployment/gateway-exporter-kubernetes-gateway-exporter \
  -n monitoring --timeout=90s

kubectl port-forward -n monitoring \
  service/gateway-exporter-kubernetes-gateway-exporter 8080:8080
```

In another terminal:

```bash
curl -fsS http://127.0.0.1:8080/metrics | grep '^exposed_route_info'
```

If no routes are exposed yet, the result is empty — that is expected until an
`Accepted` HTTPRoute resolves to a Service. See the
[Metrics reference](./metrics.md) for the full label set.

## Common configuration

The chart ships least-privilege RBAC, a non-root distroless runtime, and a
`NetworkPolicy` by default. The values you are most likely to change:

| Value | Default | Purpose |
|-------|---------|---------|
| `serviceMonitor.enabled` | `true` | Create a Prometheus Operator `ServiceMonitor`. |
| `serviceMonitor.interval` | `15s` | Scrape interval for the `ServiceMonitor`. |
| `config.enableOtel` | `"false"` | Push metrics over OTLP/gRPC in addition to Prometheus. |
| `image.digest` | `""` | Pin the image by digest (`repository@sha256:…`); takes precedence over the tag. |
| `networkPolicy.enabled` | `true` | Restrict inbound traffic to port 8080. |

### Push to an OpenTelemetry Collector

```bash
helm upgrade --install gateway-exporter \
  oci://ghcr.io/tokanize/charts/kubernetes-gateway-exporter --version "${VERSION}" \
  --namespace monitoring --create-namespace \
  --set config.enableOtel=true
```

By default the SDK targets `localhost:4317`. Override the destination with the
standard `OTEL_EXPORTER_OTLP_ENDPOINT` environment variable. See
[ADR-004: OpenTelemetry Metrics](./adrs/004-otel-metrics.md) and the
[Metrics reference](./metrics.md).

## Next steps

- [Metrics & Labels](./metrics.md) — the full `exposed_route_info` schema.
- [Architecture](./architecture.md) — how routes are resolved.
- [Security Scanning](./security-scanning.md) and [Release Verification](./verification.md).
