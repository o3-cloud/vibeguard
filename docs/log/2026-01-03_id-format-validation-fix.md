---
summary: Fixed ID format validation mismatch between validator guide documentation and actual config validation
event_type: code
sources:
  - internal/config/config.go
  - internal/cli/assist/validator_guide.go
  - internal/config/config_test.go
tags:
  - bug-fix
  - validation
  - config
  - ai-assisted-setup
  - vibeguard-whg
---

# Fixed ID Format Validation Mismatch (vibeguard-whg)

## Problem

The validator guide (`internal/cli/assist/validator_guide.go`) documented that check IDs "Must be alphanumeric with underscores and hyphens allowed", but the actual config validation in `internal/config/config.go` only checked that the ID was non-empty. This mismatch could cause AI agents to follow stricter rules than required, or miss actual invalid IDs that would work but shouldn't.

## Solution

Added regex-based ID validation to `config.go` that enforces the documented format:

```go
var validCheckID = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`)
```

This validates that IDs:
- Must start with a letter (a-z, A-Z) or underscore
- Can contain letters, numbers, underscores, and hyphens after the first character

Also updated `validator_guide.go` to be more precise about the format:
- Changed from "Must be alphanumeric with underscores and hyphens allowed"
- To "Must start with a letter or underscore" and "Can contain letters, numbers, underscores, and hyphens"

## Changes Made

1. **internal/config/config.go**
   - Added `regexp` import
   - Added `validCheckID` regex pattern
   - Added validation check after the empty ID check in `Validate()`

2. **internal/cli/assist/validator_guide.go**
   - Updated ID format documentation to match actual validation rules

3. **internal/config/config_test.go**
   - Added `TestLoad_InvalidCheckID` with 19 test cases (10 valid, 9 invalid)
   - Added `TestValidCheckIDRegex` for direct regex edge case testing

## Valid ID Examples
- `test`, `TEST`, `myTest`
- `go_test`, `go-test`
- `test123`, `_private`
- `Go_Test-123`

## Invalid ID Examples
- `123test` (starts with number)
- `-test` (starts with hyphen)
- `go test` (contains space)
- `go.test`, `go:test`, `go/test` (contains invalid characters)

## Related

- Bead: vibeguard-whg
- ADR-004: Code Quality Standards
