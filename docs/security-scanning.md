# Security Scanning

This document describes the automated security controls for `kubernetes-gateway-exporter`. The controls provide layered engineering evidence; none of them independently proves that an artifact is secure.

## Workflow summary

| Workflow | File | What it scans | Triggers | Blocking? |
|---|---|---|---|---|
| **CodeQL** | `codeql.yml` | Go SAST — static data-flow and taint analysis | PR → main, push → main, weekly | Advisory — findings appear in **Security › Code scanning**; does not block merge by default |
| **govulncheck** | `govulncheck.yml` | Known CVEs in *reachable* Go code (call-graph-aware) | PR → main, push → main, weekly | **BLOCKING** — job fails if any reachable vulnerability is found |
| **Trivy fs** | `trivy.yml` | Dependency vulns, secrets, misconfiguration in repo files | PR → main, push → main, weekly | SARIF upload is advisory; a second gate step **blocks on fixable CRITICAL vulns** |
| **Scorecard** | `scorecard.yml` | Repo security posture (branch protection, pinned deps, signed releases, etc.) | push → main, weekly | Advisory — results published to OpenSSF public API and Code scanning |
| **Release validation** | `release.yml` | Formatting, vet, staticcheck, tidy, tests, Helm lint, govulncheck, Trivy repository and image policy | push → `v*` tag | **BLOCKING** — publishing jobs depend on source validation; release aliases and signatures depend on the image gate |

## CodeQL (Go SAST)

CodeQL performs data-flow and taint analysis across the compiled Go source. Findings are uploaded to GitHub's Code scanning dashboard as SARIF. The workflow uses `autobuild` after `actions/setup-go` so the correct toolchain from `go.mod` is available before the build step.

CodeQL results do not block pull requests unless branch protection rules are explicitly configured to require the check.

## govulncheck

`govulncheck` queries the Go vulnerability database (https://vuln.go.dev) and reports only vulnerabilities that are reachable from the compiled call graph. This makes it substantially lower-noise than a raw dependency tree scan, where many reported CVEs affect code paths that are never called.

This job fails when it finds a reachable vulnerability. Whether that failure blocks a merge also depends on the repository ruleset requiring the `Go vulnerability scan` check.

## Trivy filesystem scan

Trivy scans the repository file tree for three classes of issue:

- **vuln** — CVEs in Go modules and OS packages declared in the repo
- **secret** — hardcoded credentials, API keys, and tokens
- **misconfig** — insecure configuration in Dockerfiles, Kubernetes manifests, and Helm charts

Two steps run on every trigger:

1. **SARIF step** (`exit-code: 0`) — always advisory; results are uploaded to Code scanning for triage. Fork PRs skip the upload because they lack `security-events: write` on the base repository, but the step itself still runs.
2. **Gate step** (`exit-code: 1`, `severity: CRITICAL`, `ignore-unfixed: true`) — **blocks the build** only on fixable CRITICAL findings.

**Threshold rationale:** CRITICAL-only and ignore-unfixed are chosen deliberately. Failing on HIGH/MEDIUM findings or on CVEs with no available fix adds noise without providing actionable signal, which leads teams to suppress entire categories of alerts. The gate is intended to catch egregious, patchable issues while keeping the signal-to-noise ratio high enough to act on.

### Image scanning

Trivy scanning of the **built multi-platform container image** is handled by the release workflow. Buildx first publishes the digest under a run-specific quarantine tag. Trivy scans that digest before any semver, `v`-prefixed, SHA, or `latest` release alias is assigned. The workflow then promotes the same digest without rebuilding, signs and attests it, and performs a second digest scan as a registry consistency check.

If the pre-promotion scan fails, the quarantine object can remain addressable in GHCR, but it receives no release alias, signature, provenance attestation, or GitHub Release.

## Scorecard (OpenSSF)

The [OpenSSF Scorecard](https://securityscorecards.dev) evaluates the repository's security hygiene across a set of automated checks (pinned dependencies, branch protection, signed releases, CI presence, etc.) and produces a score from 0–10. Results are published to the OpenSSF public API (powering the Scorecard badge) and uploaded to Code scanning as SARIF.

Scorecard runs only on pushes to the default branch and on the weekly schedule. Running it on pull requests would analyze a transient merge ref rather than the published repository state, making results less meaningful.

## Immutable dependency pins

GitHub Actions are pinned to immutable commit SHAs where practical. Trailing comments such as `# v6` record the upstream release line for human readability. Dependabot is configured to update the SHA pins when reviewed upstream releases become available.

Docker base images use readable tags plus immutable manifest-list digests, for example `golang:1.26-alpine@sha256:...`. The tag communicates intent while the digest prevents an upstream tag mutation from silently changing a build. Dependabot tracks both GitHub Actions and Docker image updates.

The `govulncheck` CLI is also version-pinned. Its vulnerability database is updated independently when the command runs.

## Release signing, provenance, and SBOMs

- **Cosign signatures** bind an immutable image or chart digest, or the final checksum file, to the GitHub Actions OIDC identity of the release workflow.
- **GitHub Artifact Attestations** bind images, charts, binary archives, and SBOM files to the repository, workflow, source ref, and source commit that produced them.
- **SPDX SBOMs** inventory detected packages in the image and source tree. They support incident response and downstream vulnerability analysis; they are not vulnerability scan results.
- **Signed checksums** cover every downloadable binary archive and SBOM in releases produced by the hardened workflow.

Signatures and attestations establish origin and integrity. They do not establish that the signed code is correct, vulnerability-free, reviewed, or policy-compliant.

## Scanner finding triage

The following findings are intentionally not converted into global ignore rules:

- Dockerfile `HEALTHCHECK`: Kubernetes liveness and readiness probes are the runtime health mechanism, and the distroless image has no shell utility suitable for a portable Docker health command.
- Default namespace: Helm templates intentionally omit a fixed namespace so the operator controls it with `--namespace`.
- Trusted registry restriction: the chart defaults to GHCR but keeps `image.repository` configurable for private mirrors and air-gapped environments.

Default chart security does enforce non-root execution, UID/GID 65532, a read-only root filesystem, dropped capabilities, disabled privilege escalation, and `RuntimeDefault` seccomp.

## Repository enforcement

The default branch ruleset requires changes to arrive through a pull request,
all review conversations to be resolved, and these checks to pass against the
latest target branch:

- `Lint & vet`
- `Unit tests`
- `Helm chart lint`
- `Docker build (no push)`
- `Analyze (Go)`
- `Go vulnerability scan`
- `Filesystem scan`

No second approval is required because the repository currently has a single
maintainer. Repository administrators can bypass rules only through the pull
request interface, not by directly pushing to `main`.

An independent tag ruleset prevents updates and deletion of `v*` release tags
and has no bypass actors. A release correction therefore requires a new version
instead of moving an existing tag.

## What these scans do not prove

- **Not a zero-vulnerability guarantee.** Scanners operate on known databases and static heuristics. Unknown vulnerabilities (zero-days), logic bugs, and issues in runtime configuration are outside their scope.
- **Not a compliance or policy audit.** These workflows provide engineering-level signal. They do not constitute a formal security audit, penetration test, or certification against any regulatory framework (SOC 2, PCI-DSS, etc.).
- **Not exhaustive secret detection.** Secret scanning catches many common patterns but is not a substitute for preventing secrets from entering the repository in the first place.
- **Not proof that a failed release published nothing.** A failed image gate can leave an unpromoted quarantine object in GHCR; it is not signed or exposed under a release alias.
