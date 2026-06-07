## Context

The repository already runs CI, CodeQL, govulncheck, Trivy, OpenSSF Scorecard, keyless cosign signing, GitHub Artifact Attestations, and SPDX SBOM generation. The controls are technically functional, but release jobs start concurrently, the image receives release tags and signatures before its release scan passes, SBOM files are not covered by the signed checksum manifest, and release tags are mutable at the repository level.

The project is maintained primarily by one owner, so controls must improve release integrity without requiring an unavailable second approver. Existing GitHub Actions are already pinned to immutable commit SHAs.

## Goals / Non-Goals

**Goals:**

- Make validation a prerequisite for every publishing job.
- Ensure only an image that passed the release scan receives public release aliases, signatures, and provenance.
- Prevent an existing `v*` tag from being moved or deleted.
- Make every downloadable release artifact, including SBOMs, transitively covered by one signed checksum manifest.
- Improve workload and dependency hardening while keeping automated updates practical.
- Make verification commands enforce the expected repository, workflow, source ref, and source digest.

**Non-Goals:**

- Claim a specific SLSA level or reproducible builds.
- Treat signatures or scanner success as proof that software is vulnerability-free.
- Block releases on every HIGH, unfixable, or heuristic scanner finding.
- Add a Docker `HEALTHCHECK`; Kubernetes liveness and readiness probes are the runtime health mechanism.
- Require a second human reviewer for a single-maintainer repository.

## Decisions

### Shared validation job

The release workflow will start with a `validate` job running formatting checks, vet, staticcheck, tidy verification, tests, Helm lint, govulncheck, and the blocking repository Trivy policy. The image, binaries, and chart jobs will declare `needs: validate`.

This duplicates selected CI checks intentionally. A tag can be created from a commit whose earlier checks are stale or bypassed; release validation must be self-contained.

### Quarantine then promote the exact image digest

Buildx will initially push the multi-platform image only under a technical `sha-<commit>` quarantine tag. Trivy will scan the resulting registry digest. After the gate passes, `docker buildx imagetools create` will assign all release aliases to that exact digest. Signing, SBOM generation, and provenance follow promotion.

This does not make the pre-scan image private, because GHCR package visibility applies to the repository. It does ensure failed content is not advertised under release aliases, signed, attested, or included in a GitHub Release. Rebuilding after scanning was rejected because the scanned and published digests could differ.

### Final checksum manifest in the release job

Producing jobs upload binary archives and SBOMs without a partial checksum manifest. The release job downloads all artifacts, generates one sorted SHA-256 manifest covering binary archives and both SBOMs, signs that manifest with cosign keyless signing, and publishes the manifest, signature, and certificate with the artifacts.

Binary archive attestations remain in the binary job. SBOM files receive build provenance attestations in their producing jobs.

### Immutable dependencies

GitHub Actions remain pinned to full commit SHAs with readable version comments. Docker base images will be pinned to manifest-list digests with tag names retained for readability and Dependabot compatibility. `govulncheck` will be installed at a specific reviewed version rather than `@latest`.

### Repository rules

A tag ruleset matching `refs/tags/v*` will block updates and deletion without bypass actors while allowing creation of new release tags. Existing branch rules will be augmented with pull-request and required-status-check enforcement only if the resulting rules permit a single maintainer to merge a green PR without self-approval. Direct pushes by a bypass actor remain an acknowledged governance exception if GitHub cannot express that workflow safely.

### Kubernetes hardening

The chart will set pod-level `seccompProfile.type: RuntimeDefault` and retain container-level non-root, read-only filesystem, and dropped capabilities. Scanner findings that are inapplicable to a namespaced Helm template or distroless Kubernetes workload will be documented or dismissed with rationale rather than hidden by weakening the scan.

### Consumer verification

Documentation will prefer `--repo` and `--signer-workflow` over owner-only attestation verification. Examples will resolve immutable digests first and provide optional `--source-ref` and `--source-digest` enforcement for release audits.

## Risks / Trade-offs

- [A failed scan leaves a quarantine SHA tag in public GHCR] -> The tag is not a release alias, is not signed or attested, and can be removed by scheduled registry retention later.
- [The release workflow duplicates CI work] -> The additional runtime buys an independent release gate and avoids trusting stale branch checks.
- [Digest-pinned base images reduce readability] -> Retain semantic tags and comments next to each digest and use Dependabot for updates.
- [Strict repository rules can lock out a solo maintainer] -> Apply tag immutability immediately; test branch rules through the API before enabling requirements that need another approver.
- [Scanner heuristics generate noisy findings] -> Keep SARIF visibility, use a narrow blocking policy, and document explicit triage decisions.

## Migration Plan

1. Implement and locally validate workflow, chart, Dockerfile, and documentation changes.
2. Apply immutable `v*` tag rules after confirming existing tags remain readable.
3. Apply safe branch status requirements if current check names and solo-maintainer behavior are compatible.
4. Merge through a pull request and create a new release tag; never move `v0.1.2`.
5. Verify the new image, chart, binaries, checksum signature, SBOM checksums, and provenance independently.

Rollback consists of reverting the workflow commit before creating another release tag. The immutable-tag ruleset can be disabled through repository administration if it prevents legitimate new tag creation.

## Open Questions

None. The next release version is intentionally not hard-coded into the implementation; documentation continues to show the latest verified release until a new release completes.
