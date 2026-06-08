import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "K8s Gateway Exporter",
  description: "Gateway API exposure-inventory exporter for Kubernetes. Maps Gateway→Listener→HTTPRoute→Service into Prometheus/OTel metrics.",
  base: "/kubernetes-gateway-exporter/", // Important for GitHub Pages deployment

  // Dead-link checking stays ON (VitePress default). Links that intentionally
  // leave the docs/ tree (examples/, SECURITY.md, openspec/, raw OpenAPI YAML)
  // use absolute GitHub URLs, which VitePress does not dead-check — so the build
  // stays honest instead of silencing real broken links with ignoreDeadLinks.

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/kubernetes-gateway-exporter/logo.svg' }]
  ],
  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Getting Started', link: '/getting-started' },
      { text: 'Architecture', link: '/architecture' },
      { text: 'Metrics', link: '/metrics' },
      { text: 'GitHub', link: 'https://github.com/tokanize/kubernetes-gateway-exporter' }
    ],
    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Overview', link: '/architecture' },
          { text: 'Install & Run', link: '/getting-started' },
          { text: 'Testing on kind', link: '/testing-on-kind' }
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'Metrics & Labels', link: '/metrics' }
        ]
      },
      {
        text: 'Guides',
        items: [
          { text: 'Security Scanning', link: '/security-scanning' },
          { text: 'Release Verification', link: '/verification' }
        ]
      },
      {
        text: 'Architecture Decision Records',
        collapsed: true,
        items: [
          { text: '001: Kubernetes Informers Strategy', link: '/adrs/001-k8s-informers' },
          { text: '002: Gateway Address Resolution', link: '/adrs/002-ip-resolution' },
          { text: '003: Testing Strategy', link: '/adrs/003-testing-strategy' },
          { text: '004: OpenTelemetry Metrics', link: '/adrs/004-otel-metrics' },
          { text: '005: Security Hardening', link: '/adrs/005-security' },
          { text: '006: Exposure Metric Identity', link: '/adrs/006-metric-labels-expansion' }
        ]
      }
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/tokanize/kubernetes-gateway-exporter' }
    ],
    search: {
      provider: 'local'
    },
    footer: {
      message: 'Released under the Apache 2.0 License.',
      copyright: 'Copyright © 2026 tokanize'
    }
  }
})
