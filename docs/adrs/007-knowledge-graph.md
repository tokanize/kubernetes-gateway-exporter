# ADR-007: graphify Knowledge Graph

**Date:** 2026-06-12  
**Status:** Accepted

## Context
This repository is maintained by multiple LLM assistants (Claude, Codex, Gemini) whose
identity files already instruct them to prefer structural code navigation over raw
`grep`/`find` (see the **Code Intelligence** principle in `AGENTS.md`). That principle
previously pointed only at the optional **CodeGraph MCP** server (`codegraph_explore`),
which is convenient but **local-only**: its index cache (`.codegraph/`) is git-ignored, so
nothing about the codebase's structure is shared or reviewable in the repo itself.

We wanted a code map that is (a) checked in and diffable, (b) usable by every assistant the
repo already supports, and (c) buildable with no external API dependency for the day-to-day
path.

## Decision
We adopt [graphify](https://github.com/safishamsi/graphify) and commit its knowledge graph.

1. **Tooling:** Installed via `uv tool install graphifyy` (the PyPI package is `graphifyy`
   with a double `y`; the CLI command is `graphify`). The build path used in this repo is
   `graphify update .`, which performs AST-only extraction — no LLM backend or API key
   required.
2. **Agent integration:** Each assistant is registered with `graphify <platform> install`
   (`claude` → `CLAUDE.md`, `codex` → `AGENTS.md`, `gemini` → `GEMINI.md`). This appends a
   `## graphify` ruleset to the agent file and installs PreToolUse/BeforeTool hooks.
3. **Local vs. versioned state:** Mirroring the existing AI-tooling convention in
   `.gitignore` (track shared skills/commands, ignore per-agent local state), the graphify
   **CLI and hooks are local state** — each contributor runs `graphify <platform> install`
   themselves. `CONTRIBUTING.md` documents the setup.
4. **What we commit:** Only the three shareable outputs documented by graphify —
   `graphify-out/graph.json`, `graph.html`, and `GRAPH_REPORT.md`. The regenerable AST
   cache (`graphify-out/cache/`, pinned to the graphify version) and internal scratch
   (`manifest.json`, `.graphify_root`, `.graphify_labels.json`) are git-ignored for the same
   reason `.codegraph/` is: they are derived, machine-specific, and churn on every rebuild.

## Consequences
- **Positive:** Agents (and humans) get a queryable, reviewable map of the codebase
  (`graphify query/path/explain`) that ships with the repo. The AST build is free and
  offline. CodeGraph MCP remains usable alongside it as an on-demand structural lookup.
- **Negative:** The committed `graph.html`/`graph.json` add ~660 KB to the tree and must be
  refreshed with `graphify update .` after structural code changes, or they drift. Community
  labels are placeholders (`Community N`) unless a contributor regenerates them with an LLM
  backend (`graphify label .` with an API key set) — acceptable, since traversal does not
  depend on labels.

## Traceability
The integration touches documentation surfaces only: `AGENTS.md`, `CLAUDE.md`, `CODEX.md`,
`GEMINI.md` (graphify rulesets), `README.md` (Code Intelligence section), `CONTRIBUTING.md`
(local setup), and `.gitignore` (commit-vs-ignore split). No application code is affected.
