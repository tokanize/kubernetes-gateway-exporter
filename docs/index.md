---
layout: home

hero:
  name: "Kubernetes Gateway Exporter"
  text: "See exactly what your Gateway API exposes"
  tagline: "Turns Gateway → Listener → HTTPRoute → Service relationships into Prometheus and OpenTelemetry metrics — with SLSA Level 3 supply-chain provenance built in."
  image:
    src: /logo.svg
    alt: Kubernetes Gateway Exporter
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: Architecture
      link: /architecture
    - theme: alt
      text: GitHub
      link: https://github.com/tokanize/kubernetes-gateway-exporter

features:
  - title: 🌐 Gateway API-native
    details: Resolves Gateway → Listener → HTTPRoute → Service relationships from the live cluster and exposes them as Prometheus gauges.
  - title: 📊 OpenTelemetry & OTLP
    details: Pushes metrics natively over OTLP/gRPC to an OpenTelemetry Collector, alongside the standard Prometheus /metrics scrape endpoint.
  - title: 🛡️ SLSA Level 3
    details: Every release ships Cosign-signed artifacts, keyless OIDC attestations, and full SPDX SBOMs.
  - title: ⚡ Low Overhead
    details: Cache-backed Kubernetes informers and least-privilege RBAC keep API-server load and memory footprint minimal.
  - title: 🔗 Cross-namespace aware
    details: Resolves cross-namespace routing through ReferenceGrants, so multi-team topologies map correctly.
  - title: ✅ Hardened CI/CD
    details: Trivy, govulncheck call-graph analysis, CodeQL, and OpenSSF Scorecard run on every change.
---

<div class="glow-container">
  <div class="glow glow-1"></div>
  <div class="glow glow-2"></div>
</div>
