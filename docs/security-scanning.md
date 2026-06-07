# Security Scanning

This document describes the automated security scanning layer for `kubernetes-gateway-exporter`. Each workflow runs independently and targets a specific class of risk.

## Workflow summary

| Workflow | File | What it scans | Triggers | Blocking? |
|---|---|---|---|---|
| **CodeQL** | `codeql.yml` | Go SAST — static data-flow and taint analysis | PR → main, push → main, weekly | Advisory — findings appear in **Security › Code scanning**; does not block merge by default |
| **govulncheck** | `govulncheck.yml` | Known CVEs in *reachable* Go code (call-graph-aware) | PR → main, push → main, weekly | **BLOCKING** — job fails if any reachable vulnerability is found |
| **Trivy fs** | `trivy.yml` | Dependency vulns, secrets, misconfiguration in repo files | PR → main, push → main, weekly | SARIF upload is advisory; a second gate step **blocks on fixable CRITICAL vulns** |
| **Scorecard** | `scorecard.yml` | Repo security posture (branch protection, pinned deps, signed releases, etc.) | push → main, weekly | Advisory — results published to OpenSSF public API and Code scanning |

## CodeQL (Go SAST)

CodeQL performs data-flow and taint analysis across the compiled Go source. Findings are uploaded to GitHub's Code scanning dashboard as SARIF. The workflow uses `autobuild` after `actions/setup-go` so the correct toolchain from `go.mod` is available before the build step.

CodeQL results do not block pull requests unless branch protection rules are explicitly configured to require the check.

## govulncheck

`govulncheck` queries the Go vulnerability database (https://vuln.go.dev) and reports only vulnerabilities that are reachable from the compiled call graph. This makes it substantially lower-noise than a raw dependency tree scan, where many reported CVEs affect code paths that are never called.

This job is **blocking by design**. A found vulnerability must either be resolved (upgrade the affected module) or accepted with a documented justification before the PR can merge.

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

Trivy scanning of the **built container image** is handled by the release workflow, not here. This workflow performs filesystem/repo scanning only.

## Scorecard (OpenSSF)

The [OpenSSF Scorecard](https://securityscorecards.dev) evaluates the repository's security hygiene across a set of automated checks (pinned dependencies, branch protection, signed releases, CI presence, etc.) and produces a score from 0–10. Results are published to the OpenSSF public API (powering the Scorecard badge) and uploaded to Code scanning as SARIF.

Scorecard runs only on pushes to the default branch and on the weekly schedule. Running it on pull requests would analyse a transient merge ref rather than the published repository state, making results less meaningful.

## Version tags vs SHA pins

All actions in these workflows are pinned to major version tags (e.g. `actions/checkout@v4`, `github/codeql-action/init@v3`). Major version tags are maintained by action authors and point to the latest compatible release within that major version. Dependabot is responsible for keeping these tags current.

Full SHA pinning (e.g. `actions/checkout@abc1234`) provides a stronger supply-chain guarantee because it prevents silent tag mutation, but it comes with an operational cost: SHA references must be updated manually or via tooling on every upstream patch release, and outdated SHA pins are harder to audit. SHA pinning is a documented future hardening step that should be evaluated when the project has a Dependabot configuration set up for Actions SHA tracking.

## What these scans do not prove

- **Not a zero-vulnerability guarantee.** Scanners operate on known databases and static heuristics. Unknown vulnerabilities (zero-days), logic bugs, and issues in runtime configuration are outside their scope.
- **Not a compliance or policy audit.** These workflows provide engineering-level signal. They do not constitute a formal security audit, penetration test, or certification against any regulatory framework (SOC 2, PCI-DSS, etc.).
- **Not exhaustive secret detection.** Secret scanning catches many common patterns but is not a substitute for preventing secrets from entering the repository in the first place.
