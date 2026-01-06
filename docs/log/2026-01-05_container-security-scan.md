---
summary: Container security scan using Grype identified 45 vulnerabilities in vibeguard:latest Docker image
event_type: research
sources:
  - Grype vulnerability scanner output (2026-01-05)
  - Dockerfile source code analysis
  - Ubuntu 24.04 base image security baseline
tags:
  - security
  - docker
  - grype
  - vulnerability-scan
  - container-security
  - base-image
  - cve-analysis
---

# Container Security Scan Results (vibeguard:latest)

## Executive Summary

Ran Grype vulnerability scanner on the `vibeguard:latest` Docker image built from `ubuntu:24.04` base. Identified **45 total vulnerabilities** with the following breakdown:

- **1 High severity**: CVE-2025-68973 in gpgv
- **17 Medium severity**: Various system packages (curl, libexpat1, openssl, git, tar, libpam, util-linux, bsdutils)
- **25 Low severity**: Miscellaneous standard Ubuntu packages
- **2 Negligible severity**: jq library

## Detailed Findings

### High Severity (1)

- **CVE-2025-68973** in `gpgv` (2.4.4-2ubuntu17.3)
  - Comes from base Ubuntu 24.04 image
  - Investigation needed: Is gpgv required or can it be removed?
  - Created issue: **vibeguard-920**

### Medium Severity (17)

**Top affected packages:**
- `curl` / `libcurl*` (3 vulnerabilities): CVE-2025-0167, CVE-2025-10148, CVE-2025-9086
- `libexpat1` (2 vulnerabilities): CVE-2025-59375, CVE-2025-66382
- `openssl` / `libssl3t64` (2 vulnerabilities): CVE-2024-41996, CVE-2025-27587
- `git` / `git-man` (1 vulnerability): CVE-2024-52005
- `tar` (1 vulnerability): CVE-2025-45582
- `libpam*` (4 vulnerabilities): CVE-2025-8941
- `util-linux` / related (1 vulnerability): CVE-2025-14104
- `gpgv` (1 vulnerability): CVE-2025-68972

**Action:** Created issue **vibeguard-921** to analyze and potentially mitigate these

### Low Severity (25)

Distributed across many standard Ubuntu packages. Most are in foundational libraries with EPSS scores < 1%. Created issue **vibeguard-922** for review.

### Negligible Severity (2)

- `jq` / `libjq1`: CVE-2025-9403 (Negligible impact)

## Root Cause Analysis

The Dockerfile uses `ubuntu:24.04` as the base image with minimal additional packages explicitly installed:
- curl (for downloads)
- bash (shell)
- ca-certificates (TLS)
- git (version control)
- jq (JSON processing)
- libc6 (runtime)

Most vulnerabilities are in transitive dependencies from the base image. These include:
- GPG/cryptography tools (gpgv, libgcrypt20)
- System utilities (coreutils, tar, util-linux)
- Network libraries (curl, openssl, libssh)
- PAM authentication system

## Mitigation Options

1. **High-severity CVE-2025-68973 (gpgv)**
   - Audit if gpgv is actually needed by Beads or Claude Code CLI
   - Consider using distroless or minimal base images if GPG not required

2. **Medium-severity packages**
   - Update base image when security patches available
   - Consider alpine-based or minimal distros (larger one-time migration)
   - Keep current ubuntu:24.04 if patch updates are regular

3. **Risk Assessment Context**
   - This is a development container, not production-facing
   - Vulnerabilities in utilities and libraries have low exploitability in isolated container context
   - Regular base image updates will mitigate over time

## Created Issues

- **vibeguard-920** [P2]: Address high-severity CVE-2025-68973 in gpgv
- **vibeguard-921** [P2]: Analyze and mitigate 17 Medium-severity vulnerabilities
- **vibeguard-922** [P3]: Review 25 Low-severity vulnerabilities

## Validation

✅ `vibeguard check` passed with no policy violations

## Next Steps

1. Investigate whether gpgv can be safely removed from base image
2. Evaluate base image alternatives (alpine, distroless, debian-slim)
3. Schedule regular container scans as part of CI/CD pipeline
4. Document acceptable risk levels for development vs. production containers
