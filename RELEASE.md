# Release Checklist

This document provides a comprehensive pre-release checklist for VibeGuard maintainers. Follow these steps before creating a new release to ensure quality and consistency.

## Pre-Release Validation

### 1. Code Quality

- [ ] All CI checks pass on main branch
- [ ] `golangci-lint run ./...` passes locally
- [ ] `vibeguard check -v` passes (dogfooding)
- [ ] Pre-commit hooks pass (`pre-commit run --all-files`)

### 2. Test Coverage

- [ ] Unit tests pass: `go test ./...`
- [ ] Test coverage meets 70% minimum: `go test -cover ./...`
- [ ] Integration tests pass (if applicable)
- [ ] Mutation testing efficacy ≥50% (weekly CI run)

### 3. Documentation

- [ ] CHANGELOG.md is up to date with all changes
- [ ] Breaking changes documented with migration guides
- [ ] New features documented in relevant docs files
- [ ] ADRs created for significant architectural decisions

### 4. Version Verification

- [ ] Version in `internal/version/version.go` is correct
- [ ] Version follows SemVer (see VERSIONING.md)
- [ ] Pre-release suffix removed for stable releases (e.g., `-dev` → ``)

### 5. Dependency Check

- [ ] `go mod tidy` produces no changes
- [ ] No security vulnerabilities in dependencies: `govulncheck ./...`
- [ ] Dependabot alerts reviewed and addressed

### 6. Issue Triage

- [ ] No P0/P1 issues blocking release
- [ ] In-progress issues reviewed for completion
- [ ] Stale issues triaged (see Health Metrics below)

## Release Process

### Automated Release (Recommended)

The release workflow is automated via GitHub Actions (`.github/workflows/release.yml`):

1. Ensure all pre-release checks pass
2. Merge to main branch
3. GitHub Action automatically:
   - Runs `standard-version` to bump version
   - Updates CHANGELOG.md
   - Creates git tag
   - Builds binaries for all platforms
   - Creates GitHub Release with artifacts

### Manual Release (If Needed)

```bash
# 1. Ensure clean working directory
git status

# 2. Run standard-version
npm run release

# 3. Push with tags
git push --follow-tags origin main

# 4. Build binaries manually (if not using CI)
GOOS=linux GOARCH=amd64 go build -o bin/vibeguard-linux-amd64 ./cmd/vibeguard
GOOS=darwin GOARCH=amd64 go build -o bin/vibeguard-darwin-amd64 ./cmd/vibeguard
GOOS=darwin GOARCH=arm64 go build -o bin/vibeguard-darwin-arm64 ./cmd/vibeguard
GOOS=windows GOARCH=amd64 go build -o bin/vibeguard-windows-amd64.exe ./cmd/vibeguard
```

## Post-Release Verification

- [ ] GitHub Release created with correct version tag
- [ ] All platform binaries attached to release
- [ ] CHANGELOG.md reflects the release
- [ ] Downloaded binaries work correctly on target platforms

## Health Metrics

Track these metrics to maintain project health between releases:

### Issue Health

| Metric | Target | How to Check |
|--------|--------|--------------|
| Open P0/P1 issues | 0 | `bd list --status=open` |
| In-progress issues | ≤3 | `bd list --status=in_progress` |
| Blocked issues | 0 | `bd blocked` |
| Stale issues (>30 days) | Review monthly | See Stale Issue Workflow |

### Code Health

| Metric | Target | How to Check |
|--------|--------|--------------|
| Test coverage | ≥70% | `go test -cover ./...` |
| Mutation efficacy | ≥50% | Weekly CI run |
| Lint violations | 0 | `golangci-lint run ./...` |
| Security vulnerabilities | 0 | `govulncheck ./...` |

### Project Health Commands

```bash
# Quick health check
bd stats                        # Issue database status
bd ready                        # Work available to pick up
bd blocked                      # Issues waiting on dependencies

# Detailed analysis
bd list --status=open           # All open issues
bd list --status=in_progress    # Work in flight
```

## Stale Issue Workflow

Issues become stale when they have no activity for an extended period. Follow this workflow monthly:

### Definition of Stale

- **Warning**: No activity for 14 days
- **Stale**: No activity for 30 days
- **Close candidate**: No activity for 60 days

### Monthly Triage Process

1. **Identify stale issues**
   ```bash
   # List all open issues and review updated dates
   bd list --status=open
   bd list --status=in_progress
   ```

2. **For each stale issue, choose action:**
   - **Prioritize**: If still relevant, add to sprint/milestone
   - **Update**: Add comment explaining delay or new context
   - **Close**: If no longer relevant, close with reason
   - **Defer**: Move to backlog (P4) if low priority

3. **Close stale issues**
   ```bash
   bd close <id> --reason="Stale: No longer relevant after [context]"
   ```

### Labels for Tracking

Use labels to track stale issue status:
- `stale-warning` - 14+ days inactive
- `stale` - 30+ days inactive
- `close-candidate` - 60+ days inactive, awaiting final review

## Emergency Hotfix Process

For critical security or stability fixes:

1. Create hotfix branch from latest release tag
2. Apply minimal fix only
3. Run full test suite
4. Create patch release (e.g., v0.1.1)
5. Cherry-pick to main if applicable

## References

- [VERSIONING.md](VERSIONING.md) - Semantic versioning policy
- [CHANGELOG.md](CHANGELOG.md) - Release history
- [CONTRIBUTING.md](CONTRIBUTING.md) - Development workflow
- [ADR-002](docs/adr/ADR-002-adopt-conventional-commits.md) - Conventional commits
- [ADR-004](docs/adr/ADR-004-code-quality-standards.md) - Code quality standards
