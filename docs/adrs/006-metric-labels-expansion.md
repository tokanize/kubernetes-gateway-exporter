# ADR-006: Exposure Metric Identity

**Date:** 2026-06-07  
**Status:** Accepted

## Context
The original `exposed_route_info` labels identified a Service, Gateway, and path, but did not identify the `HTTPRoute`, listener, hostname, backend namespace, or backend port that produced the relationship. Distinct Gateway API configurations could therefore collapse into the same Prometheus series.

## Decision
The metric includes these labels:

- `namespace`: HTTPRoute namespace, retained for resource identity and backward compatibility.
- `gateway_name`: referenced Gateway name.
- `route_name`: HTTPRoute name.
- `hostname`: effective hostname intersection between the Route and listener; empty when unrestricted.
- `listener_name`: accepting Gateway listener.
- `http_path`: path from an HTTPRoute match, defaulting to `/`.
- `service_name`: backend Service name.
- `backend_namespace`: explicit backend namespace or the HTTPRoute namespace by default.
- `backend_port`: backend Service port, or an empty string when omitted.
- `lb_type`: `internal` or `external`, inferred from the GatewayClass name.
- `ip_address`: first `IPAddress` from `Gateway.status.addresses`, or an empty string.

The mapper emits one series per valid listener, effective hostname, path match, backend reference, and Gateway IP value. Identical series are deduplicated.

The metric does not enumerate the Pods or EndpointSlices selected by a Service. Backend replica count does not multiply exposure series.

### Attachment Rules
The mapper requires an `Accepted=True` parent status for the current HTTPRoute generation. It resolves listeners using `sectionName` and `port`, HTTP/HTTPS protocol compatibility, `allowedRoutes.kinds`, namespace selection, and hostname intersection.

Cross-namespace Service backends require a matching `ReferenceGrant`. Each backend is evaluated independently, and the referenced Service and port must exist, so one invalid backend does not hide valid backends in the same rule.

The mapper does not require `Gateway Programmed=True`. This allows configured and accepted exposure relationships to remain visible while infrastructure is provisioning.

### Compatibility
Adding labels changes Prometheus time-series identity and is a breaking metric schema change. During rollout, Prometheus may temporarily contain both old and new series. Dashboards and recording rules must aggregate using explicit labels rather than relying on the former label set.

Approximate cardinality is:

`accepted listeners x effective hostnames x path matches x valid backendRefs`

## Consequences
- Exposure series can be traced back to a specific HTTPRoute and listener.
- Invalid listener attachments and ungranted cross-namespace backends are excluded.
- Cardinality increases with the number of valid routing combinations.
- Address discovery remains provider-independent and may temporarily yield an empty `ip_address`.

## Traceability
Implementation is located in `internal/kubernetes/mapper/mapper.go`, `pkg/models/route.go`, and `internal/metrics/`.
