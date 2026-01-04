---
summary: Documentation updates for the configurable --error-exit-code feature, completing the epic vibeguard-o7f
event_type: code
sources:
  - docs/JSON-OUTPUT-SCHEMA.md
  - docs/log/2025-12-30_vibeguard-specification.md
  - internal/executor/executor.go
  - docs/log/2026-01-03_configurable-exit-code-implementation-plan.md
tags:
  - documentation
  - exit-codes
  - cli
  - configurable
  - backwards-compatibility
---

# Configurable Exit Code Documentation Update

Completed the final task (vibeguard-907) of epic vibeguard-o7f to document the new `--error-exit-code` CLI flag.

## Changes Made

### 1. docs/JSON-OUTPUT-SCHEMA.md

- Updated top-level field description for `exit_code` to reflect new default (1 instead of 3/4)
- Revised exit code mapping table:
  - Removed separate exit codes 3 (violation) and 4 (timeout)
  - Added exit code 1 as unified configurable error code
- Added new "Configurable Error Exit Code" section documenting the `--error-exit-code` flag
- Updated all JSON examples to show exit_code: 1 instead of 3/4

### 2. docs/log/2025-12-30_vibeguard-specification.md

- Added `--error-exit-code` to the flags table in section 4.2
- Updated exit code table in section 4.4 to reflect new behavior
- Updated the "Claude Code Hook Compatibility" note with guidance for using `--error-exit-code=2`
- Updated execution flow diagram (section 7.2) to show new exit codes
- Updated example output showing exit code 1 instead of 2

### 3. internal/executor/executor.go

- Added deprecation comments to `ExitCodeViolation` and `ExitCodeTimeout` constants
- Clarified that these constants are kept for backwards compatibility in tests

### 4. Test File Updates

Added `//nolint:staticcheck` directives to test files that legitimately use deprecated constants:
- `internal/output/formatter_test.go` (4 instances)
- `internal/output/json_test.go` (6 instances)

## Key Design Decisions

1. **Unified Error Exit Code**: FAIL (error-severity violations) and TIMEOUT now use the same configurable exit code, simplifying CI/CD integration

2. **Default Change**: Default error exit code changed from 3/4 to 1, which is more conventional for CLI tools

3. **Backwards Compatibility**: Legacy exit codes 3 and 4 are deprecated but retained as constants for tests that verify specific behaviors

4. **Test Handling**: Used `//nolint:staticcheck` directives rather than removing deprecation comments, since tests legitimately verify legacy behavior

## Verification

All vibeguard checks pass:
- vet, fmt, actionlint, lint, test, test-coverage, build, mutation

## Related Issues

- Epic: vibeguard-o7f (Configurable error exit code)
- Task: vibeguard-907 (Update documentation for configurable exit code)
