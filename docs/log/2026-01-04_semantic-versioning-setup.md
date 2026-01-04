---
summary: Implemented semantic versioning and changelog automation for release management
event_type: code
sources:
  - https://www.conventionalcommits.org/
  - https://www.npmjs.com/package/commitlint
  - https://www.npmjs.com/package/standard-version
  - https://docs.github.com/en/code-security/dependabot
  - docs/adr/ADR-002-adopt-conventional-commits.md
tags:
  - versioning
  - automation
  - release-management
  - commitlint
  - conventional-changelog
  - dependabot
  - ci-cd
  - conventional-commits
---

# Semantic Versioning and Changelog Automation Setup

## Summary

Completed setup of semantic versioning and changelog automation infrastructure for VibeGuard, addressing issue vibeguard-nh8 (P1 priority). This implementation builds on ADR-002 (Conventional Commits) and enables fully automated release processes.

## Work Completed

### 1. Commitlint Configuration
- Created `.commitlintrc.json` with conventional commit rules
- Configured type-enum to match project's commit types: feat, fix, docs, style, refactor, test, chore, ci, perf
- Set maximum header length to 88 characters (consistent with project style)
- Enforces subject case rules and blank lines between sections

### 2. Conventional Changelog Configuration
- Created `.versionrc.json` for standard-version
- Configured changelog generation with semantic grouping:
  - Features, Bug Fixes, Performance improvements
  - Documentation, Testing, CI/CD changes
  - Code Refactoring (visible in changelog)
  - Code Style and Chores (hidden from changelog)
- Set GitHub URLs for commit/compare/issue links

### 3. Package.json & NPM Setup
- Created `package.json` with minimal configuration
- Included release scripts:
  - `npm run release` - automatic versioning and changelog
  - `npm run release:major/minor/patch` - explicit version bumping
- Dependencies: commitlint, standard-version, conventional-changelog-cli

### 4. Dependabot Configuration
- Created `.github/dependabot.yml` for automated dependency management
- Configured separate update schedules:
  - Go modules: weekly on Monday 03:00 UTC
  - GitHub Actions: weekly on Monday 04:00 UTC
- Pull request limits and conventional commit prefixes for dependency PRs

### 5. GitHub Actions Workflows
- **`commitlint.yml`**: Validates commit messages on PRs
  - Checks all commits against conventional commit rules
  - Fails PR if commits don't follow specification

- **`release.yml`**: Automated release workflow
  - Triggers on main branch pushes (skips if already a release commit)
  - Runs `standard-version` to bump version and update CHANGELOG.md
  - Builds binaries for Linux, macOS (x86_64 + ARM64), Windows
  - Creates GitHub release with binaries and changelog
  - Pushes tags back to repository

## Key Design Decisions

1. **Node.js for Release Tools**: Used npm/Node.js ecosystem for release automation even though project is Go-based, as these tools are language-agnostic and widely adopted.

2. **Automated vs Manual**: Release workflow is fully automated on push to main, eliminating manual release steps and reducing human error.

3. **Multi-platform Binaries**: Release workflow builds for 4 platforms automatically (Linux AMD64, macOS AMD64, macOS ARM64, Windows AMD64).

4. **Dependabot Scoping**: Limited Dependabot to 5 open PRs at a time to avoid overwhelming reviewers.

## Testing & Validation

All vibeguard checks passed:
- ✓ vet (0.5s)
- ✓ fmt (0.0s)
- ✓ lint (1.2s)
- ✓ test (4.2s)
- ✓ test-coverage (4.5s)
- ✓ build (0.3s)
- ✓ mutation (18.6s)

## Integration with Existing Standards

This implementation fully aligns with:
- **ADR-002**: Conventional Commits specification already in VERSIONING.md
- **ADR-004**: Code quality standards and pre-commit hooks
- **VERSIONING.md**: Semantic versioning policy (MAJOR.MINOR.PATCH)
- **CHANGELOG.md**: Keep a Changelog format

## Workflow Integration

The setup enables the following workflows:

1. **Development**: Developers write conventional commits (enforced by pre-commit hooks)
2. **PR Validation**: Commitlint checks all commits in PR via GitHub Actions
3. **Dependency Updates**: Dependabot creates weekly PRs with conventional commit messages
4. **Release**: Push to main automatically triggers version bump, changelog generation, and binary builds

## Next Steps

1. Install dependencies: `npm install`
2. Configure git pre-commit hook if not already done (see ADR-006)
3. Tag initial release when ready
4. Configure GitHub Actions secrets if binary signing is needed
5. Consider adding changelog validation to CI pipeline

## Files Modified/Created

- Created: `.commitlintrc.json`
- Created: `.versionrc.json`
- Created: `package.json`
- Created: `.github/dependabot.yml`
- Created: `.github/workflows/commitlint.yml`
- Created: `.github/workflows/release.yml`

## Related ADRs

- ADR-002: Adopt Conventional Commits
- ADR-004: Code Quality Standards and Tooling
- ADR-006: Integrate VibeGuard as Git Pre-Commit Hook
