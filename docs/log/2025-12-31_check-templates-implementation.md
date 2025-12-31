---
summary: Completed check templates for AI-assisted setup, added isort and pip-audit support for Python
event_type: code
sources:
  - internal/cli/inspector/recommendations.go
  - internal/cli/inspector/tools.go
  - internal/cli/inspector/recommendations_test.go
  - internal/cli/inspector/tools_test.go
  - .beads/issues.jsonl
tags:
  - ai-assisted-setup
  - check-templates
  - python
  - isort
  - pip-audit
  - security
  - recommendations
  - inspector
---

# Check Templates Implementation (vibeguard-9mi.6)

## Overview

Completed task vibeguard-9mi.6: Phase 2: Check Recommendations - Check Templates. This task involved building check templates for each tool category to enable AI-assisted vibeguard setup.

## Analysis Findings

Upon investigation, the existing `recommendations.go` already contained comprehensive check templates for:

### Go Tools (already implemented)
- `golangci-lint` - Linting with auto-fix suggestions
- `gofmt` - Format checking
- `go vet` - Static analysis
- `go test` - Testing with coverage checks
- `goimports` - Import organization

### Node.js Tools (already implemented)
- `eslint` - JavaScript/TypeScript linting
- `prettier` - Code formatting
- `jest`, `mocha`, `vitest` - Testing frameworks
- `typescript` - Type checking
- `npm audit` - Dependency security scanning

### Python Tools (partially implemented, enhanced)
- `black` - Code formatting
- `pylint` - Linting
- `pytest` - Testing with coverage
- `mypy` - Type checking
- `ruff` - Fast Python linting
- `flake8` - Style checking

## Enhancements Made

### 1. Added isort Support
- **File**: `internal/cli/inspector/tools.go`
- **Detection**: `.isort.cfg`, `pyproject.toml [tool.isort]`, `setup.cfg [isort]`, requirements files
- **Category**: Formatter (priority 11, runs early with other formatters)

### 2. Added pip-audit Support
- **File**: `internal/cli/inspector/tools.go`
- **Detection**: Explicit in requirements (high confidence 0.8) or recommended for all Python projects (0.6)
- **Category**: Security (priority 50, matches npm audit pattern)

### 3. New Recommendation Functions
```go
func (r *Recommender) isortRecommendations(tool ToolInfo) []CheckRecommendation
func (r *Recommender) pipAuditRecommendations(tool ToolInfo) []CheckRecommendation
```

## Template Pattern

Each check template follows this structure:
```go
CheckRecommendation{
    ID:          "unique-id",
    Description: "Human-readable description",
    Rationale:   "Why this check matters",
    Command:     "shell command to execute",
    Grok:        []string{}, // Optional output parsing
    Assert:      "",         // Optional assertions
    Severity:    "error|warning",
    Suggestion:  "Actionable fix guidance with {{.template}} support",
    Requires:    []string{}, // Dependencies on other checks
    Timeout:     "",         // Optional timeout
    Category:    "lint|format|test|typecheck|security|build",
    Tool:        "tool-name",
    Priority:    int,        // Lower = higher priority
}
```

## Priority Ordering

- **5**: Build (go build)
- **10**: Format (gofmt, prettier, black)
- **11**: Import organization (goimports, isort)
- **15**: Vet (go vet)
- **20**: Lint (golangci-lint, eslint, pylint, ruff)
- **25**: Type check (tsc, mypy)
- **30**: Test
- **35**: Coverage
- **50**: Security (npm audit, pip-audit)

## Test Coverage

Added comprehensive tests for new tools:
- `TestRecommender_Isort` - Verifies isort recommendations
- `TestRecommender_PipAuditSecurity` - Verifies pip-audit recommendations
- `TestRecommender_PythonFullToolchain` - Integration test for all Python tools
- `TestToolScanner_ScanPythonTools_Isort*` - Detection tests for isort
- `TestToolScanner_ScanPythonTools_PipAudit*` - Detection tests for pip-audit
- `TestToolScanner_ScanPythonTools_FullToolchain` - Full Python project detection

## Related Work

- **vibeguard-9mi.5** (closed): Recommendation Engine - Created the `Recommender` type and core logic
- **vibeguard-9mi.6** (this task): Check Templates - Added tool-specific templates
- **vibeguard-9mi.7** (pending): Prompt Structure Design - Uses these templates

## Remaining Gaps (tracked as separate issues)

The git hooks tools (`pre-commit`, `husky`, `lefthook`) return nil recommendations by design - they are hook managers, not checks themselves. This is intentional behavior documented in the code comments.
