# health-probes Specification

## Purpose
Define the readiness and liveness endpoints used by Kubernetes to manage the Pod lifecycle.

## Requirements

### Requirement: Liveness Probe
The system SHALL provide a lightweight endpoint to confirm the HTTP server is running.

#### Scenario: Server is alive
- GIVEN the application is running
- WHEN a GET request is made to `/healthz`
- THEN the system returns `200 OK` immediately.

### Requirement: Readiness Probe
The system SHALL provide an endpoint to confirm it is ready to serve metrics. It is only ready when the Kubernetes Informer cache has fully synced.

#### Scenario: Cache is syncing
- GIVEN the application just started
- AND the controller-runtime cache is still building the initial `LIST`
- WHEN a GET request is made to `/readyz`
- THEN the system returns `503 Service Unavailable`.

#### Scenario: Cache is synchronized
- GIVEN the application cache is fully synced with the Kube-APIServer
- WHEN a GET request is made to `/readyz`
- THEN the system returns `200 OK`.
