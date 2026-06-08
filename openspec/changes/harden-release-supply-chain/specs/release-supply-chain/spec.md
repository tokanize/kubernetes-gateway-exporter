## ADDED Requirements

### Requirement: Validation precedes publication
The release workflow SHALL complete its formatting, static analysis, dependency vulnerability, test, chart lint, and blocking repository scan gates before any image, binary, or chart publishing job starts.

#### Scenario: Validation fails
- **WHEN** any shared release validation gate fails
- **THEN** no release publishing job starts and no release alias, signature, attestation, chart, binary archive, or GitHub Release is produced by that run

### Requirement: Image promotion preserves the scanned digest
The release workflow SHALL scan the built multi-platform image by immutable digest before assigning release aliases, signing it, or attesting it, and SHALL promote the same digest without rebuilding.

#### Scenario: Image scan succeeds
- **WHEN** the quarantine image digest passes the blocking release scan
- **THEN** all release aliases reference that exact digest before it is signed and attested

#### Scenario: Image scan fails
- **WHEN** the quarantine image digest fails the blocking release scan
- **THEN** the workflow does not assign release aliases, sign or attest the digest, or create a GitHub Release

### Requirement: Release tags are immutable
The repository SHALL prevent updates and deletion of tags matching `v*` after creation.

#### Scenario: Existing release tag is moved
- **WHEN** an actor attempts to update an existing `v*` tag to another commit
- **THEN** GitHub rejects the update

#### Scenario: Existing release tag is deleted
- **WHEN** an actor attempts to delete an existing `v*` tag
- **THEN** GitHub rejects the deletion

### Requirement: Downloadable artifacts have verifiable integrity
Every binary archive and SBOM attached to a GitHub Release SHALL be listed in one SHA-256 checksum manifest whose keyless signature can be verified against the repository release workflow identity.

#### Scenario: Consumer verifies release files
- **WHEN** a consumer downloads release archives, SBOMs, the checksum manifest, and its signature material
- **THEN** checksum verification detects any modified covered file and cosign verification establishes the expected workflow identity for the manifest

### Requirement: Published artifacts have repository-bound provenance
The image, chart, binary archives, and SBOMs SHALL have verifiable GitHub Artifact Attestations bound to the repository release workflow and source revision.

#### Scenario: Consumer enforces provenance identity
- **WHEN** a consumer verifies an artifact using the expected repository, signer workflow, source ref, and source digest
- **THEN** verification succeeds only for an attestation matching all enforced values

### Requirement: Build dependencies are immutable
GitHub Actions and container base images used by release and CI builds SHALL be pinned to immutable commit or content digests where supported, with human-readable upstream versions retained in comments or tags.

#### Scenario: Upstream mutable tag changes
- **WHEN** an upstream action tag or container tag is moved without a reviewed repository update
- **THEN** CI and release builds continue using the previously reviewed immutable revision

### Requirement: Kubernetes workload uses baseline runtime hardening
The default Helm-rendered workload SHALL run as non-root with a read-only root filesystem, dropped Linux capabilities, privilege escalation disabled, and the RuntimeDefault seccomp profile.

#### Scenario: Chart is rendered with defaults
- **WHEN** the Helm chart is rendered without security overrides
- **THEN** the resulting Pod and container security contexts contain all required baseline controls

### Requirement: Security documentation matches enforcement
Security documentation SHALL distinguish advisory findings from blocking gates, describe immutable SHA pinning accurately, and state the limits of signatures, provenance, SBOMs, and scanners.

#### Scenario: Maintainer audits documented controls
- **WHEN** a maintainer compares the documentation with active workflows
- **THEN** trigger conditions, blocking thresholds, pinning strategy, artifact coverage, and verification commands match the implemented configuration
