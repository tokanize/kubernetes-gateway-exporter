# ADR-002: Gateway Address Resolution

**Date:** 2026-06-06  
**Status:** Accepted

## Context
The `exposed_route_info` metric includes an `ip_address` label. Gateway API already defines the portable source for this value in `Gateway.status.addresses`.

Cloud-provider API lookups would require provider-specific SDKs, credentials, naming assumptions, and additional failure handling. An allocated cloud address also does not prove that a Gateway or Route is ready to serve traffic.

## Decision
The exporter reads the first address whose type is `IPAddress` from `Gateway.status.addresses`.

When no IP address is present, the exporter emits the exposure relationship with `ip_address=""`. Route inventory therefore remains available while the Gateway controller is provisioning infrastructure.

The exporter does not query cloud-provider APIs and does not require cloud IAM permissions.

## Consequences
- The implementation remains portable across Gateway API controllers and cloud providers.
- The metric can describe configured exposure before an address is allocated.
- Address visibility follows the status update latency of the Gateway controller.
- Hostname-type Gateway addresses are not copied into the `ip_address` label.

## Traceability
Address extraction is implemented in `internal/kubernetes/mapper/mapper.go` and specified in `openspec/specs/ip-resolution/spec.md`.
