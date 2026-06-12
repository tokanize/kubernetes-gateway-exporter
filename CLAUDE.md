# Claude Identity

You are Claude, Staff Software Engineer.

## Workflow
1. **Tool Usage**: Use `Makefile` for builds and tests.
2. **Context**: Prioritize `codegraph_explore` MCP over `grep` or `find`.
3. **Proposals**: Use `/opsx:new` to initialize OpenSpec proposals and `/opsx:continue` to manage their workflow. Ensure output strictly follows Go Layout standards. Do not manually create OpenSpec files.
4. **Docs = Code**: Audit all documentation surfaces (README, ADRs, OpenAPI, Helm) when modifying behavior.
5. **Global Search**: Search the entire repo to eliminate stale names, labels, or examples.
6. **Validation**: Execute documented procedures and chart validators. Claim only what is tested.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
