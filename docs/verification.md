# Release Verification Guide

This guide explains how to verify the integrity and provenance of every
kubernetes-gateway-exporter release artifact: the container image, binary
archives, and SBOMs. Verification is optional but recommended when you are
pulling a release into a production environment or an automated pipeline.

For the CI scanning layer (static analysis, container scanning, dependency
auditing) see [docs/security-scanning.md](./security-scanning.md).

---

## Prerequisites

| Tool | Purpose | Install |
|------|---------|---------|
| `docker` with Buildx | Resolve and pull image manifests | [Docker Engine](https://docs.docker.com/engine/install/) or [Docker Desktop](https://docs.docker.com/desktop/) |
| `cosign` | Verify image signature and blob signature | `go install github.com/sigstore/cosign/v2/cmd/cosign@latest` or download a release from [sigstore/cosign](https://github.com/sigstore/cosign/releases) |
| `crane` | Resolve the OCI chart digest | `go install github.com/google/go-containerregistry/cmd/crane@latest` |
| `gh` | Verify GitHub Artifact Attestation (build provenance) | [cli.github.com](https://cli.github.com/) |
| `helm` | Pull and install the verified chart | [helm.sh/docs/intro/install](https://helm.sh/docs/intro/install/) |
| `jq` | Resolve the image digest and inspect SBOM JSON | `brew install jq` / `apt install jq` |
| `sha256sum` | Verify binary checksums | Included on Linux; use `shasum -a 256` on macOS |

You do not need any cosign key material. All signatures use keyless signing
via GitHub OIDC and the public Sigstore transparency log.

---

## 1. Pull the container image by digest

Container image tags are mutable — a tag can be reassigned to a different
image without notice. A digest (the `sha256:…` hash of the image manifest) is
immutable and uniquely identifies a specific build.

Set these once and reuse them across every step below:

```bash
IMAGE=ghcr.io/tokanize/kubernetes-gateway-exporter
VERSION=0.1.2 # Replace with the release you want to verify
```

**Step 1 — resolve the digest for a tag:**

```bash
docker buildx imagetools inspect "${IMAGE}:${VERSION}"
# Capture the manifest digest directly:
DIGEST=$(docker buildx imagetools inspect "${IMAGE}:${VERSION}" \
  --format '{{json .}}' | jq -r '.manifest.digest')
echo "$DIGEST"   # e.g. sha256:abc123…
```

**Step 2 — pull by digest:**

```bash
docker pull "${IMAGE}@${DIGEST}"
```

All subsequent verification steps use `${DIGEST}`, not the tag.

---

## 2. Verify the cosign signature of the image

The release workflow signs the image by digest using cosign keyless signing
(GitHub OIDC). Verification checks that the signature is present in the
Sigstore Rekor transparency log and that the signing identity matches the
`tokanize/kubernetes-gateway-exporter` release workflow.

```bash
cosign verify \
  --certificate-identity-regexp "https://github.com/tokanize/kubernetes-gateway-exporter" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "${IMAGE}@${DIGEST}"
```

A successful run prints the verified payload JSON. An error means either the
signature is absent or the signing identity does not match.

---

## 3. Verify build provenance

GitHub Artifact Attestation records a cryptographic link between the image (or
binary archive) and the exact source commit and workflow run that produced it.

**Image:**

```bash
gh attestation verify "oci://${IMAGE}@${DIGEST}" --owner tokanize
```

**Binary archive (after downloading from the GitHub Release):**

```bash
gh attestation verify \
  "kubernetes-gateway-exporter_v${VERSION}_linux_amd64.tar.gz" \
  --owner tokanize
```

Swap the OS/architecture in the filename as needed. Binary archives are
published for the following combinations:

- `linux_amd64`, `linux_arm64`
- `darwin_amd64`, `darwin_arm64`

---

## 4. Verify checksums and their signature

Each GitHub Release ships `checksums.txt`, `checksums.txt.sig`, and
`checksums.txt.pem`.

**Step 1 — verify the checksum file:**

```bash
sha256sum -c checksums.txt --ignore-missing
```

On macOS use `shasum -a 256 -c checksums.txt --ignore-missing`.

**Step 2 — verify the cosign blob signature of `checksums.txt`:**

```bash
cosign verify-blob \
  --certificate-identity-regexp "https://github.com/tokanize/kubernetes-gateway-exporter" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  --signature checksums.txt.sig \
  --certificate checksums.txt.pem \
  checksums.txt
```

This confirms the checksum file itself was produced by the official release
workflow and has not been tampered with.

---

## 5. Find and inspect SBOMs

Two SPDX SBOMs are attached to every GitHub Release as release assets:

| Asset | Contents |
|-------|----------|
| `sbom-image.spdx.json` | Packages present in the container image |
| `sbom-source.spdx.json` | Go module dependencies from the source tree |

Download from the GitHub Release page, then inspect with `jq`:

```bash
# List all package names in the image SBOM
jq '.packages[].name' sbom-image.spdx.json

# Count packages
jq '.packages | length' sbom-image.spdx.json
```

---

## 6. Verify the published Helm chart

The chart is published as a signed OCI artifact at
`oci://ghcr.io/tokanize/charts/kubernetes-gateway-exporter`. Resolve its digest,
then verify the cosign signature and build provenance exactly as for the image:



```bash
CHART=ghcr.io/tokanize/charts/kubernetes-gateway-exporter

# Resolve the chart digest for a version
CHART_DIGEST=$(crane digest "${CHART}:${VERSION}")

# Verify the signature
cosign verify \
  --certificate-identity-regexp "https://github.com/tokanize/kubernetes-gateway-exporter" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "${CHART}@${CHART_DIGEST}"

# Verify build provenance
gh attestation verify \
  "oci://${CHART}@${CHART_DIGEST}" \
  --owner tokanize
```

Then pull and install the verified chart:

```bash
CHART=ghcr.io/tokanize/charts/kubernetes-gateway-exporter
VERSION=0.1.2 # Replace with the release you want to install

helm install gateway-exporter \
  "oci://${CHART}" --version "${VERSION}"
```

---

## What verification proves

- The image or binary was built by the `tokanize/kubernetes-gateway-exporter`
  GitHub repository and workflow, not by a third party.
- The cosign signature was issued via GitHub OIDC and recorded in the public
  Sigstore Rekor transparency log. No long-lived signing key was involved.
- The build provenance attestation links the artifact to a specific source
  commit and workflow run, making the build reproducible to audit.
- The checksums file was signed by the same workflow, so a matching checksum
  confirms the binary archive has not been altered after release.

## What verification does NOT prove

- That the software is free of security vulnerabilities. Run the tools
  described in [docs/security-scanning.md](./security-scanning.md) for that.
- That your deployment environment is correctly configured or secured.
- That the Gateway API exposure model implemented here is policy-compliant for
  your organisation's specific requirements.

---

*Questions or concerns about a specific release? Open an issue or see
[SECURITY.md](../SECURITY.md) for the vulnerability reporting process.*
