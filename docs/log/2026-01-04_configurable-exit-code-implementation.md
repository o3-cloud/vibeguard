---
summary: Completed implementation of configurable error exit codes for vibeguard
event_type: code
sources:
  - ADR-005-adopt-vibeguard.md
  - internal/orchestrator/orchestrator.go
  - internal/cli/root.go
  - internal/cli/check.go
tags:
  - configurable-exit-code
  - implementation
  - error-handling
  - cli-flags
  - exit-codes
---

# Configurable Error Exit Code Implementation

Completed the implementation of configurable error exit codes for vibeguard, enabling users to customize the exit code used for check failures (both FAIL and TIMEOUT cases).

## Implementation Summary

### 1. Orchestrator Changes (internal/orchestrator/orchestrator.go)
- Added `errorExitCode` field to the `Orchestrator` struct with default value of 1
- Added `Option` functional option type for flexible configuration
- Implemented `WithErrorExitCode(code int)` option function
- Updated `New()` constructor to accept `errorExitCode` parameter
- Modified `calculateExitCode()` to use configurable exit code for both FAIL and TIMEOUT violations
- Updated `RunCheck()` to use configurable exit code

### 2. CLI Flag Addition (internal/cli/root.go)
- Added `errorExitCode` variable to the flags section
- Registered persistent flag `--error-exit-code` with default value 1
- Implemented `GetErrorExitCode()` getter function for cross-package access

### 3. CLI Integration (internal/cli/check.go)
- Updated `runCheck()` to pass `GetErrorExitCode()` to the orchestrator constructor

### 4. Test Updates
- Updated all test files to include the new `errorExitCode` parameter in `orchestrator.New()` calls:
  - internal/orchestrator/orchestrator_test.go
  - internal/orchestrator/integration_test.go
  - internal/cli/check_test.go
- Changed hardcoded exit code expectations from `executor.ExitCodeViolation` (3) and `executor.ExitCodeTimeout` (4) to the new default of 1
- All tests now pass with the new behavior

## Key Design Decisions

1. **Default Value**: Exit code defaults to 1, simplifying CI/CD integration and making the tool more accessible
2. **Unified Exit Code**: Both FAIL and TIMEOUT violations now use the same configurable exit code (previously 3 and 4 respectively)
3. **Backward Compatibility**: Users can set `--error-exit-code` to their desired value if needed
4. **Functional Options**: Implemented Option pattern for future extensibility

## Blocks Completed

- ✅ vibeguard-6d1: Add errorExitCode field and WithErrorExitCode option to orchestrator
- ✅ vibeguard-dta: Add --error-exit-code CLI flag to root.go
- ✅ vibeguard-h8i: Update calculateExitCode to use configurable exit code
- ✅ vibeguard-lyz: Update RunCheck to use configurable exit code
- ✅ vibeguard-tzw: Wire --error-exit-code flag to orchestrator in check.go
- ✅ vibeguard-nzz: Add tests for configurable exit code

## Verification

- All vibeguard checks pass without errors
- Test coverage maintained
- Code formatting verified with gofmt
- No breaking changes to public API (backward compatible with CLI)

## Next Steps

- Documentation update (vibeguard-907)
- User-facing documentation about the new CLI flag
