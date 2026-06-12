# Codex Identity

You are Codex, AI Pair Programmer.

## Workflow
1. **Testing**: Generate robust, table-driven unit tests.
2. **Security**: Ensure HTTP servers include Slowloris timeouts (`docs/adrs/005-security.md`).
3. **Verification**: Always run `go fmt ./...` and `make test`.
4. **Docs = Code**: Audit all documentation surfaces (README, ADRs, OpenAPI, Helm) when modifying behavior.
5. **OpenSpec**: Use `/opsx:new` to initialize new specifications and `/opsx:continue` to manage them. Do NOT create OpenSpec directories manually.
6. **Global Search**: Search the entire repo to eliminate stale names, labels, or examples.
7. **Validation**: Run OpenSpec and Helm validators. Claim only what is tested.

## graphify

This project commits a knowledge graph at `graphify-out/`. For codebase questions, run `graphify query "<question>"` (also `graphify path "<A>" "<B>"` and `graphify explain "<concept>"`) before grepping raw files, and run `graphify update .` after code changes. The full ruleset lives in the `## graphify` section of [AGENTS.md](./AGENTS.md).
