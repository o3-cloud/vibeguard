---
summary: Implemented monthly Docker image rebuild workflow for automatic security patch updates
event_type: code
sources:
  - https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions
  - https://docs.docker.com/build/ci/github-actions/
  - https://github.com/docker/build-push-action
tags:
  - docker
  - security
  - ci-cd
  - github-actions
  - security-patches
  - automation
  - vibeguard-924
---

# Monthly Docker Image Rebuild Implementation

Successfully implemented a GitHub Actions workflow to automatically rebuild the Docker image on a monthly schedule, ensuring base system packages receive security updates.

## Changes Made

### 1. Created `.github/workflows/docker-rebuild.yml`

Implemented a comprehensive GitHub Actions workflow with:

- **Scheduled Trigger**: First day of each month at 00:00 UTC using cron expression
- **Manual Trigger**: `workflow_dispatch` for ad-hoc rebuilds with optional reason input
- **Multi-Platform Build**: Docker Buildx setup for building across architectures
- **Container Registry Push**: GitHub Container Registry (GHCR) authentication and push
- **Semantic Versioning**: Image tags include:
  - Branch references
  - Semantic version tags
  - SHA-based content addressable tags
  - `latest` tag
  - Run-numbered rebuild tags for traceability
- **Build Caching**: Registry-based layer caching for faster subsequent rebuilds
- **Build Summary**: Automated documentation of rebuild trigger, reason, and status

### 2. Fixed actionlint Compliance Issues

Resolved shellcheck violations in the workflow:
- Properly quoted variables in shell script to prevent globbing/word splitting
- Replaced individual echo redirects with grouped output using brace syntax
- Changed shell comparison from `==` to `=` for POSIX compliance
- Grouped all summary echo statements into single redirect to `$GITHUB_STEP_SUMMARY`

## Test Results

All vibeguard compliance checks passed successfully:
- ✓ vet (Go code validation)
- ✓ fmt (Code formatting)
- ✓ actionlint (GitHub Actions workflow validation)
- ✓ lint (golangci-lint)
- ✓ staticcheck (Static code analysis)
- ✓ test (Unit tests)
- ✓ test-coverage (Coverage threshold verification)
- ✓ gosec (Security scanning)
- ✓ docker (Docker build validation)
- ✓ build (Binary compilation)

## Key Features

### Security Focus
- Automated monthly rebuilds ensure base Ubuntu 24.04 image receives latest security patches
- Reduces window of exposure to published CVEs in system packages
- Complements existing in-image mitigation strategies (CVE-2025-68973, CVE-2025-14104)

### Flexible Triggering
- **Scheduled**: Consistent monthly cadence for predictable security updates
- **Manual**: Allows immediate rebuilds for critical security advisories
- **Input Parameters**: Optional reason field for rebuild justification

### Container Registry
- Pushes to ghcr.io for GitHub ecosystem integration
- Accessible to Actions workflows for downstream dependencies
- Supports Docker pull-through caching for faster local builds

### Build Cache Management
- Registry-based caching reduces build times for subsequent rebuilds
- Improves feedback loop for security patch validation
- Supports multi-platform builds without rebuilding from scratch

### Operational Visibility
- Build summaries document each rebuild in GitHub Actions UI
- Trigger type and reason tracked for audit purposes
- Clear success/failure status for monitoring

## Design Decisions

1. **Monthly Schedule**: Balances timely security patching with resource efficiency
2. **GHCR Push**: Leverages GitHub's container registry for tight CI/CD integration
3. **Caching Strategy**: Registry cache improves rebuild performance while maintaining image freshness
4. **Manual Trigger**: Allows security team to initiate immediate rebuilds if critical CVEs are discovered

## Related Issues and Dependencies

- **Resolves**: vibeguard-924 (Implement monthly Docker image rebuild for security patches)
- **Blocks**: No dependent tasks identified
- **Related**: ADR-005 (Adopt Vibeguard for Policy Enforcement), ADR-006 (Git Pre-Commit Hook Integration)

## Next Steps

- Monitor first scheduled rebuild execution (2026-02-01)
- Consider expanding to include vulnerability scanning in the rebuild pipeline
- Track image rebuild metrics for security patch lag analysis
