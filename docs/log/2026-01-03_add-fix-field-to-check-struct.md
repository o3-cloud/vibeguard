---
summary: Added Fix field to Check struct as part of AI agent output improvements
event_type: code
sources:
  - docs/log/2026-01-02_check-output-format-implementation-plan.md
  - internal/config/schema.go
tags:
  - schema
  - check-struct
  - fix-field
  - ai-agent-output
  - vibeguard-1zn
---

# Add Fix Field to Check Struct

Implemented task `vibeguard-1zn.1` - adding a `Fix` field to the Check struct in `internal/config/schema.go`.

## Change Summary

Added the following field to the `Check` struct:

```go
Fix string `yaml:"fix,omitempty"`
```

The field was placed after `Suggestion` to maintain logical grouping - both fields provide guidance to users when a check fails:
- `Suggestion`: explains what failed (the "what")
- `Fix`: provides actionable steps to resolve the issue (the "how")

## Rationale

This change is part of epic `vibeguard-1zn` (Check Output Format Improvements for AI Agents). The goal is to make vibeguard output more actionable for AI agents by:

1. Separating "what failed" (suggestion) from "how to fix" (fix)
2. Using WARN vs FAIL headers based on severity
3. Always showing Advisory line to make blocking status explicit

## Verification

- Build: Passes (`go build ./...`)
- Tests: All 14 packages pass
- Linter: 0 issues (`golangci-lint run ./...`)

## Next Steps

This task unblocks:
- `vibeguard-1zn.2`: Add Fix field to Violation struct
- `vibeguard-1zn.13`: Update vibeguard.yaml with fix field examples
