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

## Commit and PR etiquette

- Keep commits small and reviewable; one logical change per commit.
- Ensure `make verify` passes before opening a PR.
- Write a clear commit message explaining *why*, not just what changed.
