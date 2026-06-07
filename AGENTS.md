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
4. **CodeGraph**: Leverage CodeGraph MCP server if available.
5. **Global Audit**: Update all documentation surfaces (README, ADRs, OpenSpec, OpenAPI, Helm, diagrams) when changing names, APIs, or semantics.
6. **Stale References**: Search and destroy outdated names, labels, and versions. Run validators to verify correctness.
7. **Adversarial Review**: Define Acceptance Criteria in your implementation plan. Before finalizing code, spawn a dedicated autonomous subagent to aggressively review for edge cases, semantic flaws, and security holes.
