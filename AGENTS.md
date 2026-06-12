# Autonomous Agent Routing Map

Welcome to `kubernetes-gateway-exporter`. This repository uses **Spec-Driven Development**. Read your identity file before making changes.

## Agent Identities
- **[Gemini](./GEMINI.md)**
- **[Claude](./CLAUDE.md)**
- **[Codex](./CODEX.md)**

## Core Principles
1. **Docs-First**: See `docs/architecture.md`. No code without architectural justification.
2. **Spec-Driven**: Read `openspec/specs/` before implementing features. Use `/opsx:new` (or `openspec new change`) to initialize new specifications and `/opsx:continue` to drive the workflow. Do NOT manually create or scaffold OpenSpec directories.
3. **Traceability**: Link all code changes to an ADR or OpenSpec scenario.
4. **Code Intelligence**: Prefer structural navigation over raw `grep`/`find`. Query the committed graphify knowledge graph (`graphify query/path/explain`, see the `## graphify` section below) and the CodeGraph MCP server (`codegraph_explore`) when available.
5. **Global Audit**: Update all documentation surfaces (README, ADRs, OpenSpec, OpenAPI, Helm, diagrams) when changing names, APIs, or semantics.
6. **Stale References**: Search and destroy outdated names, labels, and versions. Run validators to verify correctness.
7. **Adversarial Review**: Define Acceptance Criteria in your implementation plan. Before finalizing code, spawn a dedicated autonomous subagent to aggressively review for edge cases, semantic flaws, and security holes.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

When the user types `/graphify`, invoke the `skill` tool with `skill: "graphify"` before doing anything else.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- Dirty graphify-out/ files are expected after hooks or incremental updates; dirty graph files are not a reason to skip graphify. Only skip graphify if the task is about stale or incorrect graph output, or the user explicitly says not to use it.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
