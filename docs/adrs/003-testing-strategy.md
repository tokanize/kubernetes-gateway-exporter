# ADR-003: Testing Strategy

**Date:** 2026-06-06  
**Status:** Accepted

## Context
A Kubernetes SRE/infrastructure component requires extreme reliability. Testing strategies that heavily rely on mocked Kubernetes clients often fail to catch actual schema mismatches or subtle informer caching issues.

## Decision
We mandate a multi-tiered, highly effective testing strategy:

1. **Integration Testing (`envtest`):**
   - Informer startup, cache synchronization, and CRD integration SHOULD be tested using `sigs.k8s.io/controller-runtime/pkg/envtest` when test assets are available.

2. **Unit Testing (Pure Logic):**
   - Pure mapping functions, metric label generation, and Gateway API relationship logic MUST use idiomatic table-driven tests.
   - Kubernetes object graph tests MAY use controller-runtime's fake client, with the Gateway API types registered in the runtime scheme.

3. **Metrics Contract Testing:**
   - Prometheus collectors MUST be registered with a pedantic registry and gathered in tests.
   - Tests MUST verify the complete metric label set.

## Consequences
- **Positive:** High confidence in production. Tests will catch API version skew issues early.
- **Negative:** Full `envtest` coverage requires additional local control-plane binaries and is slower than fake-client tests.

## Traceability
All test files (`*_test.go`) MUST adhere to these rules and reference this document in their package or suite comments.

Manual end-to-end verification without a Gateway controller is documented in `docs/testing-on-kind.md`.
