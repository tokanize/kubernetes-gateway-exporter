# k8s-informers Specification

## Purpose
Define the data caching and mapping strategy to ensure the application does not overload the Kubernetes Control Plane.

> **Implementation Details:** See [ADR-001: Kubernetes Informers](../../../docs/adrs/001-k8s-informers.md) for the architectural decision on `controller-runtime` caching.

## Requirements

### Requirement: Event-driven Caching
The system SHALL NOT continuously poll the Kubernetes API server for resource state. It must use cache-backed Informers.

#### Scenario: Initialization
- GIVEN the application starts
- WHEN the controller-runtime manager connects to the Kube-APIServer
- THEN informer-backed reads initialize caches for the resource types used by the mapper
- AND watches keep those caches synchronized with API-server changes.

### Requirement: Resource Mapping
The system SHALL correctly map the relationship between Gateway API resources and their underlying Services.

#### Scenario: Valid Gateway Route
- GIVEN an `HTTPRoute` references a `Gateway` via `parentRefs`
- AND the `HTTPRoute` routes traffic to a `Service` via `backendRefs`
- AND a matching parent status has `Accepted=True` for the current Route generation
- WHEN metrics are collected
- THEN the system correlates the Gateway, accepting listener, effective hostname, Route path, and Service backend into metric vectors.

#### Scenario: Listener rejects the Route
- GIVEN a listener excludes `HTTPRoute` through protocol, `allowedRoutes.kinds`, namespace rules, or hostname rules
- WHEN metrics are collected
- THEN no exposure relationship is emitted for that listener.

#### Scenario: Cross-namespace backend
- GIVEN an HTTPRoute references a Service in another namespace
- WHEN no matching `ReferenceGrant` exists in the Service namespace
- THEN that backend relationship is not emitted
- AND valid backends in the same Route remain eligible for emission.

#### Scenario: Backend Service does not exist
- GIVEN an HTTPRoute backendRef targets a missing Service or a missing Service port
- WHEN metrics are collected
- THEN that backend relationship is not emitted.
