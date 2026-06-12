# Graph Report - kubernetes-gateway-exporter  (2026-06-12)

## Corpus Check
- 81 files · ~53,171 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 483 nodes · 489 edges · 76 communities (74 shown, 2 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 24 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6a0f0b5b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]

## God Nodes (most connected - your core abstractions)
1. `Mapper` - 27 edges
2. `Security Scanning` - 11 edges
3. `Testing Locally on kind` - 10 edges
4. `Release Verification Guide` - 10 edges
5. `sameParentReference()` - 9 edges
6. `ADDED Requirements` - 9 edges
7. `T` - 8 edges
8. `Setup()` - 8 edges
9. `Decisions` - 8 edges
10. `main()` - 7 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewMapper()`  [INFERRED]
  cmd/exporter/main.go → internal/kubernetes/mapper/mapper.go
- `main()` --calls--> `NewExporter()`  [INFERRED]
  cmd/exporter/main.go → internal/metrics/exporter.go
- `main()` --calls--> `SetupOTEL()`  [INFERRED]
  cmd/exporter/main.go → internal/metrics/otel.go
- `main()` --calls--> `RunAsync()`  [INFERRED]
  cmd/exporter/main.go → internal/server/server.go
- `main()` --calls--> `Setup()`  [INFERRED]
  cmd/exporter/main.go → internal/server/server.go

## Import Cycles
- None detected.

## Communities (76 total, 2 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.16
Nodes (25): ExposedRoute, Gateway, HTTPBackendRef, HTTPRoute, HTTPRouteRule, Context, Listener, Mapper (+17 more)

### Community 1 - "Community 1"
Cohesion: 0.11
Nodes (17): main(), Handler, Context, Logger, Mapper, Logger, Logger, T (+9 more)

### Community 2 - "Community 2"
Cohesion: 0.18
Nodes (19): Client, Hostname, Logger, T, T, intersectHostnames(), NewMapper(), ptr() (+11 more)

### Community 3 - "Community 3"
Cohesion: 0.10
Nodes (19): Codebase Analysis, Graceful Exit Handling, Guardrails, Phase 10: Archive, Phase 11: Recap & Next Steps, Phase 1: Welcome, Phase 2: Task Selection, Phase 3: Explore Demo (+11 more)

### Community 4 - "Community 4"
Cohesion: 0.10
Nodes (19): Codebase Analysis, Graceful Exit Handling, Guardrails, Phase 10: Archive, Phase 11: Recap & Next Steps, Phase 1: Welcome, Phase 2: Task Selection, Phase 3: Explore Demo (+11 more)

### Community 5 - "Community 5"
Cohesion: 0.10
Nodes (19): Codebase Analysis, Graceful Exit Handling, Guardrails, Phase 10: Archive, Phase 11: Recap & Next Steps, Phase 1: Welcome, Phase 2: Task Selection, Phase 3: Explore Demo (+11 more)

### Community 6 - "Community 6"
Cohesion: 0.10
Nodes (19): Codebase Analysis, Graceful Exit Handling, Guardrails, Phase 10: Archive, Phase 11: Recap & Next Steps, Phase 1: Welcome, Phase 2: Task Selection, Phase 3: Explore Demo (+11 more)

### Community 7 - "Community 7"
Cohesion: 0.10
Nodes (19): ADDED Requirements, Requirement: Build dependencies are immutable, Requirement: Downloadable artifacts have verifiable integrity, Requirement: Image promotion preserves the scanned digest, Requirement: Kubernetes workload uses baseline runtime hardening, Requirement: Published artifacts have repository-bound provenance, Requirement: Release tags are immutable, Requirement: Security documentation matches enforcement (+11 more)

### Community 8 - "Community 8"
Cohesion: 0.13
Nodes (14): dependencies, @fontsource/inter, vitepress, vue, description, name, overrides, esbuild (+6 more)

### Community 9 - "Community 9"
Cohesion: 0.14
Nodes (13): Consumer verification, Context, Decisions, Final checksum manifest in the release job, Goals / Non-Goals, Immutable dependencies, Kubernetes hardening, Migration Plan (+5 more)

### Community 10 - "Community 10"
Cohesion: 0.14
Nodes (13): 1. General & Architecture, 2. APIs & Requirements, 3. Agent Instructions, 4. Security & Releases, 5. Guides, Build and Test, Deploy, Documentation Directory (+5 more)

### Community 11 - "Community 11"
Cohesion: 0.15
Nodes (12): CodeQL (Go SAST), govulncheck, Image scanning, Immutable dependency pins, Release signing, provenance, and SBOMs, Repository enforcement, Scanner finding triage, Scorecard (OpenSSF) (+4 more)

### Community 12 - "Community 12"
Cohesion: 0.18
Nodes (10): Check for context, Ending Discovery, Guardrails, Handling Different Entry Points, OpenSpec Awareness, The Stance, What You Don't Have To Do, What You Might Do (+2 more)

### Community 13 - "Community 13"
Cohesion: 0.18
Nodes (10): Check for context, Ending Discovery, Guardrails, Handling Different Entry Points, OpenSpec Awareness, The Stance, What You Don't Have To Do, What You Might Do (+2 more)

### Community 14 - "Community 14"
Cohesion: 0.18
Nodes (10): 1. Create the cluster and install Gateway API, 2. Build and load the image, 3. Deploy the exporter, 4. Create the test resources, 5. Simulate Gateway controller acceptance, 6. Verify the metric, 7. Verify that acceptance is required, 8. Clean up (+2 more)

### Community 15 - "Community 15"
Cohesion: 0.18
Nodes (10): 1. Pull the container image by digest, 2. Verify the cosign signature of the image, 3. Verify build provenance, 4. Verify checksums and their signature, 5. Find and inspect SBOMs, 6. Verify the published Helm chart, Prerequisites, Release Verification Guide (+2 more)

### Community 16 - "Community 16"
Cohesion: 0.18
Nodes (10): Check for context, Ending Discovery, Guardrails, Handling Different Entry Points, OpenSpec Awareness, The Stance, What You Don't Have To Do, What You Might Do (+2 more)

### Community 17 - "Community 17"
Cohesion: 0.18
Nodes (10): k8s-informers Specification, Purpose, Requirement: Event-driven Caching, Requirement: Resource Mapping, Requirements, Scenario: Backend Service does not exist, Scenario: Cross-namespace backend, Scenario: Initialization (+2 more)

### Community 18 - "Community 18"
Cohesion: 0.20
Nodes (9): ip-resolution Specification, Purpose, Requirement: Provider Independence, Requirement: Retrieve IP from Gateway Status, Requirements, Scenario: Address is missing from Gateway status, Scenario: Gateway address is not allocated yet, Scenario: Gateway exposes only a hostname address (+1 more)

### Community 19 - "Community 19"
Cohesion: 0.20
Nodes (9): metrics-exposition Specification, Purpose, Requirement: OpenTelemetry Push, Requirement: Prometheus Exporter, Requirements, Scenario: OTEL metrics push is disabled, Scenario: OTEL metrics push is enabled, Scenario: Scrape metrics (+1 more)

### Community 20 - "Community 20"
Cohesion: 0.20
Nodes (9): Check for context, Ending Discovery, Guardrails, OpenSpec Awareness, The Stance, What You Don't Have To Do, What You Might Do, When a change exists (+1 more)

### Community 21 - "Community 21"
Cohesion: 0.22
Nodes (8): Common configuration, Confirm it is running, Getting Started, Install with Helm, Next steps, Prerequisites, Push to an OpenTelemetry Collector, Without the Prometheus Operator

### Community 22 - "Community 22"
Cohesion: 0.22
Nodes (8): health-probes Specification, Purpose, Requirement: Liveness Probe, Requirement: Readiness Probe, Requirements, Scenario: Cache is synchronized, Scenario: Cache is syncing, Scenario: Server is alive

### Community 23 - "Community 23"
Cohesion: 0.25
Nodes (7): ADR-006: Exposure Metric Identity, Attachment Rules, Compatibility, Consequences, Context, Decision, Traceability

### Community 24 - "Community 24"
Cohesion: 0.36
Nodes (6): Desc, Logger, Mapper, Metric, Exporter, NewExporter()

### Community 25 - "Community 25"
Cohesion: 0.25
Nodes (7): Cardinality, Configuration, Endpoints, Example series, `exposed_route_info`, Metrics & Labels, What it does not measure

### Community 26 - "Community 26"
Cohesion: 0.29
Nodes (6): API specification, Architecture, Architecture Decision Records, Documentation strategy, Further reading, How it works

### Community 27 - "Community 27"
Cohesion: 0.29
Nodes (6): Capabilities, Impact, Modified Capabilities, New Capabilities, What Changes, Why

### Community 28 - "Community 28"
Cohesion: 0.33
Nodes (5): ADR-001: Kubernetes Informers Strategy, Consequences, Context, Decision, Traceability

### Community 29 - "Community 29"
Cohesion: 0.33
Nodes (5): ADR-002: Gateway Address Resolution, Consequences, Context, Decision, Traceability

### Community 30 - "Community 30"
Cohesion: 0.33
Nodes (5): ADR-003: Testing Strategy, Consequences, Context, Decision, Traceability

### Community 31 - "Community 31"
Cohesion: 0.33
Nodes (5): ADR-004: OpenTelemetry Metrics, Consequences, Context, Decision, Traceability

### Community 32 - "Community 32"
Cohesion: 0.33
Nodes (5): CI, Commit and PR etiquette, Contributing, Prerequisites, Running checks locally

### Community 33 - "Community 33"
Cohesion: 0.33
Nodes (5): 1. Release Gates and Image Promotion, 2. Artifact Integrity and Provenance, 3. Dependency and Runtime Hardening, 4. Repository Policy and Documentation, 5. Validation and Review

### Community 34 - "Community 34"
Cohesion: 0.33
Nodes (5): Recommended branch protection (informational), Release verification, Reporting a vulnerability, Security Policy, Supported versions

### Community 35 - "Community 35"
Cohesion: 0.40
Nodes (4): ADR-005: Security Hardening, Context, Decisions Implemented, Future Recommendations (What else needs securing?)

### Community 36 - "Community 36"
Cohesion: 0.40
Nodes (4): Agent Identities, Autonomous Agent Routing Map, Core Principles, graphify

### Community 37 - "Community 37"
Cohesion: 0.50
Nodes (3): Claude Identity, graphify, Workflow

### Community 38 - "Community 38"
Cohesion: 0.50
Nodes (3): Gemini Identity, graphify, Workflow

## Knowledge Gaps
- **271 isolated node(s):** `name`, `version`, `description`, `docs:dev`, `docs:build` (+266 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewMapper()` connect `Community 2` to `Community 0`, `Community 1`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Why does `Mapper` connect `Community 0` to `Community 2`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Why does `main()` connect `Community 1` to `Community 24`, `Community 2`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **What connects `name`, `version`, `description` to the rest of the system?**
  _271 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.10869565217391304 - nodes in this community are weakly interconnected._
- **Should `Community 3` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `Community 4` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._