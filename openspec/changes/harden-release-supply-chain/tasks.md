## 1. Release Gates and Image Promotion

- [x] 1.1 Add a shared release validation job and make all publishing jobs depend on it
- [x] 1.2 Build the multi-platform image under a quarantine SHA tag, scan its digest, and promote the unchanged digest to release aliases before signing
- [x] 1.3 Keep the published digest scan and ensure GitHub Release creation remains dependent on image, binary, and chart success

## 2. Artifact Integrity and Provenance

- [x] 2.1 Generate one final checksum manifest covering binary archives and both SBOM files
- [x] 2.2 Sign the final checksum manifest keylessly and publish all verification material
- [x] 2.3 Add build provenance attestations for image and source SBOM files while retaining image, chart, and binary attestations
- [x] 2.4 Strengthen release notes and verification commands with repository and workflow identity constraints

## 3. Dependency and Runtime Hardening

- [x] 3.1 Pin Docker base images to immutable manifest digests with readable version references
- [x] 3.2 Pin govulncheck to a reviewed version instead of resolving `latest` during CI
- [x] 3.3 Add RuntimeDefault seccomp to the Helm workload and verify all baseline security-context controls in rendered manifests

## 4. Repository Policy and Documentation

- [x] 4.1 Create an immutable `v*` release-tag ruleset and verify update and deletion are blocked
- [x] 4.2 Evaluate and apply safe pull-request and required-check rules for the default branch without requiring a second maintainer
- [x] 4.3 Update security-scanning and verification documentation to match SHA pinning, gates, artifact coverage, limitations, and triage decisions
- [x] 4.4 Audit the repository for stale release-security claims and inconsistent pinning guidance

## 5. Validation and Review

- [ ] 5.1 Run formatting, tests, static analysis, Helm lint, OpenSpec validation, workflow lint, and documentation checks
- [ ] 5.2 Perform an independent adversarial review against every release-supply-chain acceptance scenario
