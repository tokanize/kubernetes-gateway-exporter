# ADR-001: Kubernetes Informers Strategy

**Date:** 2026-06-06  
**Status:** Accepted

## Context
The metrics exporter needs an up-to-date representation of `Gateway`, `HTTPRoute`, `Service`, `Namespace`, and `ReferenceGrant` resources without issuing uncached API requests on every scrape.

## Decision
The exporter uses the cache-backed client provided by a controller-runtime manager. Calls made by the mapper are served from informer caches.

To keep scraper requests predictable, the exporter moves Gateway API relationship resolution out of the scrape path. Scrapers read a shallow copy of a precomputed in-memory snapshot. The mapper is registered as an asynchronous `manager.Runnable` in `cmd/exporter/main.go`. It periodically calculates the route mapping in the background (every 15 seconds) and stores the results in a thread-safe in-memory cache protected by a reader-writer lock (`sync.RWMutex`).

Scrapers (Prometheus and OpenTelemetry) read from this in-memory cache. The expensive Kubernetes relationship resolution is moved out of the scrape path. Scrape cost is proportional only to the number of already-resolved metric series, reading a precomputed snapshot.

The mapper reads:

- `gateway.networking.k8s.io/v1`: `Gateway`, `HTTPRoute`, and `ReferenceGrant`.
- `core/v1`: `Service` and `Namespace`.

## Consequences
- **Positive:** Collection avoids direct API-server reads. Scrape path no longer performs graph traversal.
- **Positive:** Garbage collection overhead on scrapes is drastically reduced compared to active resolution.
- **Negative:** Cache memory usage is slightly increased, and changes in Gateway API resources are reflected with a delay of up to the refresh interval (15 seconds).

## Traceability
The manager and cached client are configured in `cmd/exporter/main.go`; relationship reads are implemented in `internal/kubernetes/mapper/mapper.go`.
