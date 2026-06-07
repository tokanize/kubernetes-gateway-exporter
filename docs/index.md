---
layout: home

hero:
  name: "Kubernetes Gateway Exporter"
  text: "Exposure Inventory & Telemetry"
  tagline: "Maps Gateway API routes to Prometheus and OpenTelemetry metrics with native SLSA Level 3 security."
  image:
    src: /logo.svg
    alt: Kubernetes Gateway Exporter
  actions:
    - theme: brand
      text: View Architecture
      link: /architecture
    - theme: alt
      text: View on GitHub
      link: https://github.com/tokanize/kubernetes-gateway-exporter

features:
  - title: 🌐 Gateway API Native
    details: Automatically maps Gateway → Listener → HTTPRoute → Service relationships and exposes them as Prometheus gauges.
  - title: 📊 OTLP / OpenTelemetry
    details: Supports native pushing via gRPC to OpenTelemetry Collectors alongside standard Prometheus /metrics scraping.
  - title: 🛡️ SLSA Level 3
    details: Cryptographically signed artifacts (Cosign), keyless OIDC attestations, and full SBOMs for every release.
  - title: 🚀 Zero Noise
    details: Cache-backed Kubernetes Informers with minimal RBAC privileges ensure high performance and low API server load.
  - title: 🔍 Cross-Namespace Ready
    details: Seamlessly resolves complex networking topologies including ReferenceGrants for cross-namespace routing.
  - title: ✅ CI/CD Hardened
    details: Integrated Trivy scanning, Govulncheck call-graph analysis, and OSSF Scorecard validations.
---

<div class="glow-container">
  <div class="glow glow-1"></div>
  <div class="glow glow-2"></div>
</div>
