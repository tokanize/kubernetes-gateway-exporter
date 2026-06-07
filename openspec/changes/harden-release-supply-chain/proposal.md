## Why

The release pipeline produces verifiable signatures, provenance, and SBOMs, but publication can begin before all validation succeeds, release tags can be moved, and documentation understates the existing SHA pinning. These gaps allow partial or superseded artifacts to remain publicly available despite a failed or repeated release.

## What Changes

- Add a shared validation gate that must pass before image, binary, or chart publication starts.
- Scan the built image artifact before registry publication and retain a digest-based scan of the published multi-architecture image.
- Protect release tags from creation bypass, updates, and deletion after publication.
- Include SBOMs in signed release checksums and attest their provenance.
- Pin Docker base images by immutable digest while retaining readable version comments and Dependabot updates.
- Harden the Helm workload security context and triage scanner findings that are not applicable.
- Strengthen verification commands to bind attestations and signatures to this repository, workflow, source ref, and source commit.
- Correct security documentation to state that GitHub Actions are pinned to immutable commit SHAs.

## Capabilities

### New Capabilities
- `release-supply-chain`: Security gates, immutable release identity, artifact integrity, provenance, and consumer verification requirements for published releases.

### Modified Capabilities

None.

## Impact

Affected surfaces include GitHub Actions workflows, repository rulesets, the Dockerfile, Helm chart defaults and templates, release verification documentation, security-scanning documentation, and release artifact composition. No application API or metric semantics change.
