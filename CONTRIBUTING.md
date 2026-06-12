# Contributing

## Prerequisites

- Go 1.26+
- Docker
- Helm
- make

## Running checks locally

`make verify` runs the same checks as CI in one command:

```sh
make verify
```

You can also run individual checks:

| Command | What it does |
|---|---|
| `make fmt-check` | Fail if any file is not gofmt-formatted |
| `make lint` | Run `go vet` and staticcheck |
| `make test` | Run all unit tests |
| `make helm-lint` | Lint the Helm chart |
| `make docker-build` | Build the Docker image locally (no push) |

staticcheck is fetched on demand via `go run` — no manual installation needed.

## CI

Pull requests trigger the `CI` workflow, which runs four parallel jobs:

- **Lint & vet** — `fmt-check`, `go vet`, `staticcheck`, `tidy-check`
- **Unit tests** — `go test ./...`
- **Helm chart lint** — `helm lint`
- **Docker build** — validates the Dockerfile builds (no image is pushed)

Security scanning (CodeQL, govulncheck, Trivy) runs via separate workflows.

## Knowledge graph (graphify)

This repo commits a [graphify](https://github.com/safishamsi/graphify) knowledge graph
(`graphify-out/graph.json`, `graph.html`, `GRAPH_REPORT.md`) so assistants can answer
codebase questions from a scoped subgraph instead of grepping the whole tree. See
[ADR-007](./docs/adrs/007-knowledge-graph.md) for the rationale.

The graph artifacts are versioned, but graphify's CLI and per-agent hooks are **local
state** (each contributor installs their own — they are git-ignored, like CodeGraph's
`.codegraph/`). To set it up:

```sh
# 1. Install the CLI (package is `graphifyy` with a double y; the command is `graphify`)
uv tool install graphifyy        # or: pipx install graphifyy

# 2. Register the skill + hooks for your assistant (writes to your agent file + local config)
graphify claude install          # Claude Code   → CLAUDE.md
graphify codex install           # Codex         → AGENTS.md
graphify gemini install          # Gemini CLI    → GEMINI.md
```

After changing code, refresh the committed graph so it stays accurate (AST-only, no API key needed):

```sh
graphify update .
```

Only the three shareable outputs are tracked; the regenerable AST cache
(`graphify-out/cache/`) and internal scratch are git-ignored.

## Commit and PR etiquette

- Keep commits small and reviewable; one logical change per commit.
- Ensure `make verify` passes before opening a PR.
- Write a clear commit message explaining *why*, not just what changed.
