# ADR-005: Security Hardening

**Date:** 2026-06-06  
**Status:** Accepted

## Context
As a Kubernetes infrastructure component with `ClusterRole` read permissions, this exporter requires stringent security measures to prevent privilege escalation, resource exhaustion, and unauthorized access.

## Decisions Implemented

1. **HTTP Server Hardening (Slowloris Protection):**
   - The default Go `http.Server` is vulnerable to Slowloris attacks where a client opens a connection but sends data very slowly, exhausting the server's connection pool.
   - We configured explicit timeouts: `ReadTimeout` (5s), `ReadHeaderTimeout` (2s), `WriteTimeout` (10s), and `IdleTimeout` (120s) to forcefully drop slow or dead connections.

2. **Minimal Distroless Container:**
   - The application is containerized using a multi-stage Docker build.
   - The runtime image is `gcr.io/distroless/static:nonroot`. It contains no shell or package manager, reducing the runtime attack surface.
   - The application runs entirely as `nonroot` (UID 65532).
   - Builder and runtime images are pinned to immutable manifest-list digests;
     readable tags remain alongside the digests for maintenance context.

3. **Kubernetes Runtime Hardening:**
   - The default chart sets `runAsNonRoot`, UID/GID 65532, a read-only root
     filesystem, disabled privilege escalation, dropped capabilities, and the
     `RuntimeDefault` seccomp profile.

4. **Readiness Probe Integrity:**
   - The `/readyz` endpoint ensures the Kubernetes API cache (`controller-runtime` informer) is fully synchronized before returning 200 OK. This prevents routing traffic to an empty instance.

5. **Principle of Least Privilege (RBAC):**
   - The Helm chart grants read-only access to `Gateway`, `HTTPRoute`, `ReferenceGrant`, `Namespace`, and `Service` resources required for relationship resolution.

6. **NetworkPolicy Baseline:**
   - A `NetworkPolicy` limits inbound connections to port 8080. Egress is currently unrestricted because API-server and optional OTEL destinations vary by cluster; operators should narrow it for their environment.

## Future Recommendations (What else needs securing?)

To achieve a "Zero Trust" production state, the following should be implemented in the cluster environment:

1. **mTLS (Mutual TLS):** 
   - Wrap the pod in a Service Mesh (e.g., Istio or Linkerd) to encrypt traffic between Prometheus and this exporter, or configure Prometheus to use TLS client certificates.
