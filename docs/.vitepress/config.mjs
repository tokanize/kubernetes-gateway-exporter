import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "K8s Gateway Exporter",
  description: "Gateway API exposure-inventory exporter for Kubernetes. Maps Gateway→Listener→HTTPRoute→Service into Prometheus/OTel metrics.",
  base: "/kubernetes-gateway-exporter/", // Important for GitHub Pages deployment
  ignoreDeadLinks: true,
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/kubernetes-gateway-exporter/logo.svg' }]
  ],
  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Architecture', link: '/architecture' },
      { text: 'GitHub', link: 'https://github.com/tokanize/kubernetes-gateway-exporter' }
    ],
    sidebar: [
      {
        text: 'Overview',
        items: [
          { text: 'Architecture', link: '/architecture' },
          { text: 'Security Scanning', link: '/security-scanning' },
          { text: 'Verification', link: '/verification' },
          { text: 'Testing on Kind', link: '/testing-on-kind' }
        ]
      },
      {
        text: 'Architecture Decision Records',
        collapsed: false,
        items: [
          { text: '001: K8s Informers', link: '/adrs/001-k8s-informers' },
          { text: '002: IP Resolution', link: '/adrs/002-ip-resolution' },
          { text: '003: Testing Strategy', link: '/adrs/003-testing-strategy' },
          { text: '004: OTel Metrics', link: '/adrs/004-otel-metrics' },
          { text: '005: Security', link: '/adrs/005-security' },
          { text: '006: Metric Labels', link: '/adrs/006-metric-labels-expansion' }
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
