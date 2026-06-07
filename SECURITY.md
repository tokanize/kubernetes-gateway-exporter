# Security Policy

## Reporting a vulnerability

**Please do not open a public GitHub issue for security bugs.**

The preferred channel is GitHub private vulnerability reporting:
1. Go to the repository's **Security** tab.
2. Click **"Report a vulnerability"**.
3. Fill in the description, affected versions, and steps to reproduce.

If you prefer email, reach the maintainer at **hello@toka.email**. Include as
much detail as you can: affected release, steps to reproduce, and your
assessment of impact.

We will make a best-effort acknowledgement within a few business days. There is
no formal SLA — this is an open-source project maintained by a single
maintainer.

---

## Supported versions

The project is pre-1.0. Only the latest release and the `main` branch receive
security fixes.

| Version | Supported |
|---------|-----------|
| Latest release | Yes |
| `main` branch | Yes |
| Older releases | No — please upgrade |

---

## Release verification

Every release ships a cosign keyless image signature, GitHub Artifact
Attestation (build provenance) for the image and binary archives, SPDX SBOMs,
and a signed `checksums.txt`.

See **[docs/verification.md](./docs/verification.md)** for step-by-step
copy-pasteable verification commands.

---

## Recommended branch protection (informational)

These settings cannot be enforced from repository files, but are recommended
for any fork or derivative that wishes to maintain a similar security posture:

- Require at least one pull request review before merging to `main`.
- Require the following status checks to pass before merge:
  - `CI` (build + unit tests)
  - `CodeQL` (static analysis)
  - `govulncheck` (Go vulnerability database)
  - `Trivy` (container and dependency scanning)
- Optionally enforce linear history and/or signed commits.

These are recommendations, not guarantees about this repository's current
configuration.
