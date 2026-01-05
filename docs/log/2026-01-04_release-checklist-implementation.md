---
summary: Implemented RELEASE.md checklist with health metrics and stale issue workflow
event_type: code
sources:
  - RELEASE.md
  - VERSIONING.md
  - .github/workflows/release.yml
  - docs/adr/ADR-004-code-quality-standards.md
tags:
  - release
  - documentation
  - health-metrics
  - stale-issues
  - maintainability
  - beads
  - vibeguard-8e9
---

# RELEASE.md Checklist Implementation

Completed task vibeguard-8e9: "Improve maintainability with RELEASE.md checklist and health metrics"

## What Was Created

Created `/RELEASE.md` with three main sections:

### 1. Pre-Release Validation Checklist

Comprehensive checklist covering:
- Code quality (CI, linting, vibeguard dogfooding, pre-commit hooks)
- Test coverage (unit tests, 70% minimum, integration, mutation testing)
- Documentation (CHANGELOG, breaking changes, ADRs)
- Version verification (version.go, SemVer compliance)
- Dependency checks (go mod tidy, govulncheck, Dependabot alerts)
- Issue triage (P0/P1 blockers, in-progress review, stale issues)

### 2. Health Metrics Tracking

Established targets for ongoing project health:

**Issue Health:**
| Metric | Target |
|--------|--------|
| Open P0/P1 issues | 0 |
| In-progress issues | ≤3 |
| Blocked issues | 0 |
| Stale issues (>30 days) | Monthly review |

**Code Health:**
| Metric | Target |
|--------|--------|
| Test coverage | ≥70% |
| Mutation efficacy | ≥50% |
| Lint violations | 0 |
| Security vulnerabilities | 0 |

### 3. Stale Issue Workflow

Defined stale issue lifecycle:
- **Warning**: 14 days no activity
- **Stale**: 30 days no activity
- **Close candidate**: 60 days no activity

Monthly triage process with clear actions:
1. Identify stale issues via `bd list`
2. Choose action: Prioritize, Update, Close, or Defer
3. Use labels for tracking: `stale-warning`, `stale`, `close-candidate`

## Current Project State

From `bd stats`:
- Total Issues: 129
- Open: 6
- In Progress: 2
- Blocked: 1
- Closed: 121
- Ready to Work: 5
- Avg Lead Time: 12.5 hours

## Verification

All vibeguard checks passed:
- vet, fmt, actionlint, lint, test, test-coverage, build, mutation

## References

- Builds on ADR-004 code quality standards
- Complements VERSIONING.md semantic versioning policy
- Integrates with existing release.yml automation
- Uses beads (ADR-001) for issue tracking commands
