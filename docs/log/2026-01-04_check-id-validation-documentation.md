---
summary: Documented check ID validation rules in README.md and validator_guide.go
event_type: code
sources:
  - internal/config/config.go:16
  - README.md
  - internal/cli/assist/validator_guide.go
tags:
  - documentation
  - validation
  - config
  - check-id
---

# Check ID Validation Documentation

Completed task vibeguard-r92 to document check ID validation rules that were previously undocumented in user-facing documentation.

## Context

The check ID regex pattern `^[a-zA-Z_][a-zA-Z0-9_-]*$` was defined at `internal/config/config.go:16` with inline comments, but the user-facing documentation didn't explain these constraints to users creating `vibeguard.yaml` configurations.

## Changes Made

### 1. README.md - Field Details Table

Updated the `id` field description in the Configuration Schema section to include the validation rules:

```markdown
| `id` | Yes (per check) | string | Unique check identifier. Must start with a letter or underscore, followed by alphanumeric characters, underscores, or hyphens (regex: `^[a-zA-Z_][a-zA-Z0-9_-]*$`) | — |
```

### 2. validator_guide.go - Explicit DO NOT List

Added explicit constraint for AI agents generating configurations:

```
- DO NOT use invalid check IDs (must match ^[a-zA-Z_][a-zA-Z0-9_-]*$ - start with letter or underscore, then alphanumeric/underscore/hyphen)
```

## Existing Documentation

The `CheckStructureRules` constant in `validator_guide.go` already had good documentation (lines 87-91):
- Must start with a letter or underscore
- Can contain letters, numbers, underscores, and hyphens
- Must be unique across all checks
- Examples: "fmt", "lint", "go-test", "npm_audit", "_private"

## Verification

All vibeguard checks passed:
- vet, fmt, actionlint, lint, test, test-coverage, build, mutation
