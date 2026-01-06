---
summary: Set up automated release pipeline with goreleaser and GitHub Actions for multi-platform binary distribution
event_type: code
sources:
  - https://goreleaser.com/intro/
  - https://github.com/goreleaser/goreleaser-action
  - https://docs.github.com/en/actions/publishing-packages/publishing-java-packages-with-gradle
tags:
  - release-automation
  - goreleaser
  - github-actions
  - ci-cd
  - multi-platform-builds
  - homebrew
---

# Release Automation Implementation

Successfully implemented automated release pipeline for VibeGuard using goreleaser and GitHub Actions, enabling seamless multi-platform binary distribution.

## Completed Tasks

### 1. Created `.goreleaser.yaml` Configuration
- **Multi-platform build support**: Linux, macOS, and Windows builds for both amd64 and arm64 architectures
- **Archive generation**: Automatic tar.gz archives for Unix systems and zip files for Windows
- **Homebrew formula generation**: Auto-generates both classic brews and homebrew_casks entries
- **Changelog filtering**: Excludes documentation, test, chore, and CI commits from release notes
- **GitHub release creation**: Automatic release creation with checksums and multi-platform binaries
- **Build optimization**: Stripped binaries with `-s -w` ldflags to minimize binary size

### 2. Refactored GitHub Actions Release Workflow
- **Semantic release job**: Runs standard-version to automatically bump versions and update CHANGELOG.md
- **Goreleaser job**: Depends on semantic-release, builds and publishes all platform binaries
- **Notification job**: Provides release status feedback for monitoring
- **Proper job chaining**: Uses GitHub Actions output dependencies to ensure correct execution order

### 3. Tested Goreleaser Configuration
- Successfully built snapshot binaries for all 6 platform/architecture combinations
- Generated valid Homebrew formula for distribution
- Created checksums and archive files
- All builds completed without errors (runtime: ~11 seconds)

### 4. Validation Results
All VibeGuard policy checks passed:
- ✅ vet (0.2s)
- ✅ fmt (0.1s)
- ✅ actionlint (0.1s)
- ✅ lint (2.7s)
- ✅ staticcheck (0.5s)
- ✅ test (7.2s)
- ✅ test-coverage (6.1s)
- ✅ gosec (7.1s)
- ✅ build (0.2s)

## Technical Details

### Environment
- Go version: 1.24.4
- Goreleaser version: 2.13.2 (via Homebrew)
- Project module: github.com/vibeguard/vibeguard

### Build Artifacts Generated
- `vibeguard_linux_amd64`, `vibeguard_linux_arm64`
- `vibeguard_darwin_amd64`, `vibeguard_darwin_arm64` (macOS Intel and Apple Silicon)
- `vibeguard_windows_amd64.exe`, `vibeguard_windows_arm64.exe`

### Homebrew Distribution
Formula includes:
- Multi-architecture support (macOS x86_64, ARM64, Linux x86_64, ARM64)
- License: Apache-2.0
- Auto-detection of platform and architecture

## Configuration Files Modified

1. **`.goreleaser.yaml`** (new)
   - 92 lines of goreleaser configuration
   - Supports version 2 schema

2. **`.github/workflows/release.yml`** (refactored)
   - Split into three jobs: semantic-release, goreleaser, and notify
   - Uses official goreleaser GitHub Action (v6)
   - Proper environment variable passing and token management

## Known Limitations & Future Steps

### Pre-Release Setup Required
- User must create `homebrew-vibeguard` and `homebrew-casks` repositories
- `HOMEBREW_TAP_GITHUB_TOKEN` secret needs to be configured in GitHub Actions settings
- Homebrew formula upload is currently disabled (`skip_upload: false` allows manual testing)

### Deprecation Warnings in Goreleaser
The configuration contains some deprecation notices:
- `snapshot.name_template` → should use `snapshot.version_template`
- `archives.format` → consider newer format specifications
- `brews` → phased out in favor of `homebrew_casks`

These are preserved for backward compatibility but could be cleaned up in future iterations.

## Integration with Existing Systems

### Semantic Versioning
Integrates with existing `.versionrc.json` and standard-version setup for automated version bumping following Conventional Commits (ADR-002).

### Code Quality
Release workflow respects VibeGuard policy enforcement (ADR-005) by requiring all checks to pass before building.

### Project Documentation
Aligns with `RELEASE.md` documentation and complements manual release process defined there.

## Next Steps for Release

1. Verify goreleaser is installed in CI environment (already added to local dev setup)
2. Create Homebrew tap repositories if planning to publish to Homebrew
3. Configure GitHub Actions secrets for Homebrew token if auto-publishing is desired
4. Push a semantic version tag (e.g., `v1.0.0`) to trigger the automated release
5. Monitor the GitHub Actions workflow execution and verify artifacts are created

## References

- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [GoReleaser GitHub Action](https://github.com/goreleaser/goreleaser-action)
- [ADR-002: Conventional Commits](../adr/ADR-002-adopt-conventional-commits.md)
- [ADR-005: Adopt VibeGuard for Policy Enforcement](../adr/ADR-005-adopt-vibeguard.md)
- [RELEASE.md](../../RELEASE.md) - Full release process documentation
