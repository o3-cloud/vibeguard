---
summary: Implemented recommendation engine for AI-assisted setup feature that maps detected tools to suggested VibeGuard checks
event_type: code
sources:
  - internal/cli/inspector/recommendations.go
  - internal/cli/inspector/recommendations_test.go
  - docs/log/2025-12-31_agent-assisted-setup-implementation-spec.md
  - internal/cli/inspector/tools.go
tags:
  - vibeguard
  - ai-assisted-setup
  - recommendation-engine
  - inspector
  - check-generation
  - phase-2
---

# Recommendation Engine Implementation (vibeguard-9mi.5)

Completed implementation of the recommendation engine for the AI agent-assisted setup feature. This engine maps detected development tools to suggested VibeGuard check configurations.

## Implementation Overview

Created `internal/cli/inspector/recommendations.go` with the following components:

### Core Types

- **CheckRecommendation**: Struct representing a suggested check with all configuration fields:
  - ID, Description, Rationale
  - Command, Grok patterns, Assert expression
  - Severity, Suggestion template
  - Requires (dependencies), Timeout
  - Category, Tool, Priority

- **Recommender**: Main engine that generates recommendations based on project type and detected tools

### Key Methods

1. `NewRecommender(projectType, tools)` - Creates a recommender instance
2. `Recommend()` - Generates all recommendations sorted by priority
3. `RecommendForCategory(category)` - Filters recommendations by category

### Utility Functions

- `DeduplicateRecommendations()` - Removes duplicate check IDs
- `FilterByTools()` - Filters recommendations to specific tools
- `GroupByCategory()` - Groups recommendations by category (lint, format, test, etc.)

## Supported Tools

### Go
- golangci-lint (lint)
- gofmt (format)
- go vet (lint)
- go test (test + coverage)
- goimports (format)
- go build (project-level)

### Node.js
- eslint (lint)
- prettier (format)
- jest (test + coverage)
- mocha (test)
- vitest (test)
- typescript (typecheck)
- npm audit (security)

### Python
- black (format)
- pylint (lint)
- ruff (lint)
- flake8 (lint)
- pytest (test + coverage)
- mypy (typecheck)

### Hook Tools
- pre-commit, husky, lefthook - No check recommendations (meta-tools)

## Design Decisions

### Priority-Based Ordering
Recommendations are sorted by priority (lower = higher):
- Build: 5
- Format: 10-11
- Vet/Static Analysis: 15
- Lint: 20
- Typecheck: 25
- Test: 30
- Coverage: 35
- Build scripts: 40
- Security: 50

This ensures fast, deterministic checks run first.

### Grok + Assert for Coverage
Coverage checks include grok patterns and assertions for threshold enforcement:
```go
Grok:   []string{"coverage: %{NUMBER:coverage}%"},
Assert: "coverage >= 70",
```

### Deduplication Strategy
Multiple tools may recommend checks with the same ID (e.g., both ruff and pylint recommend "lint"). Deduplication keeps the first occurrence, allowing the higher-level code to control which tool takes precedence.

## Test Coverage

Created comprehensive unit tests in `recommendations_test.go`:
- Project type recommendations (Go, Node, Python)
- Tool-specific recommendations
- Priority sorting
- Category filtering
- Deduplication
- Required fields validation
- Hook tools exclusion
- Multiple test frameworks handling

All 22 new tests pass.

## Integration Points

The Recommender integrates with:
- `Detector` (from detector.go) - Gets project type
- `ToolScanner` (from tools.go) - Gets detected tools

Next phase (vibeguard-9mi.6) will create check templates that render these recommendations into YAML format.

## Files Changed

- **Added**: `internal/cli/inspector/recommendations.go` (~500 lines)
- **Added**: `internal/cli/inspector/recommendations_test.go` (~450 lines)
