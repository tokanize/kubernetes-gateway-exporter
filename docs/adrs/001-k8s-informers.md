# ADR-001: Kubernetes Informers Strategy

**Date:** 2026-06-06  
**Status:** Accepted

## Context
The metrics exporter needs an up-to-date representation of `Gateway`, `HTTPRoute`, `Service`, `Namespace`, and `ReferenceGrant` resources without issuing uncached API requests on every scrape.

## Decision
The exporter uses the cache-backed client provided by a controller-runtime manager. Calls made by the mapper are served from informer caches after the manager starts and synchronizes them.

The mapper reads:

- `gateway.networking.k8s.io/v1`: `Gateway`, `HTTPRoute`, and `ReferenceGrant`.
- `core/v1`: `Service` and `Namespace`.

The current implementation lists HTTPRoutes from the cache for each collection cycle and performs cached object lookups while resolving relationships. It does not implement custom field indexers.

## Consequences
- **Positive:** Collection avoids direct API-server reads after cache synchronization.
- **Negative:** Collection cost grows with the number of HTTPRoutes and resolved relationships. The cache also increases memory usage.

## Traceability
The manager and cached client are configured in `cmd/exporter/main.go`; relationship reads are implemented in `internal/kubernetes/mapper/mapper.go`.
