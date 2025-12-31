---
summary: Fixed unknown check error handling in vibeguard-cce
event_type: code
sources: []
tags:
  - bug
  - fix
  - error-handling
  - unknown-check
---

# Fix: Unknown Check Error Handling (vibeguard-cce)

## Summary

Fixed issue where unknown check IDs were silently exiting with code 2 instead of returning a proper error that gets printed to the user.

## Problem

When `RunCheck()` was called with an unknown check ID, it would return:
- `nil` error
- `RunResult` with `ExitCode: 2` (ConfigError)

The `check.go` file would then call `os.Exit(2)` directly, bypassing all error handlers in `main.go`. This meant:
1. No error message was printed to the user
2. The error wasn't catchable by error handling code
3. The behavior was inconsistent with other error handling paths

## Root Cause

The issue was a design inconsistency:
- `RunCheck()` returned `(nil error, RunResult with exit code)` for unknown checks
- Other error paths returned `(error, nil result)`
- The CLI's `runCheck()` function used `os.Exit()` directly instead of returning errors normally
- `main.go` error handlers never saw the error

## Solution

Implemented proper error propagation:

1. **Modified `orchestrator.go:RunCheck()`** - Return `ConfigError` when check not found instead of `(nil, RunResult)`
2. **Added `ExitError` type in `cli/check.go`** - Represents completed checks with non-zero exit codes
3. **Updated `main.go` error handling** - Added `ExitError` type check to handle exit codes properly
4. **Removed direct `os.Exit()` call** - Now returns `ExitError` from `runCheck()` for proper error flow
5. **Updated test expectations** - Changed `TestRunCheck_UnknownCheck_ExitCodeTwo` to `TestRunCheck_UnknownCheck_ReturnsConfigError`

## Files Changed

- `internal/orchestrator/orchestrator.go`: Return `ConfigError` for unknown checks, added `fmt` import
- `internal/cli/check.go`: Added `ExitError` type, proper error handling, removed `os.Exit()`
- `cmd/vibeguard/main.go`: Added `ExitError` handling with `errors.As()`
- `internal/orchestrator/orchestrator_test.go`: Updated test to expect `ConfigError`

## Testing

✅ All unit tests pass
✅ Manual testing confirms unknown check IDs return proper error message with exit code 2
✅ Normal check execution works correctly with exit code 0

## Example Behavior (Before vs After)

**Before (Silent Exit):**
```bash
$ vibeguard check --config ./vibeguard.yaml non-existent-check
$ echo $?
2
# No error message printed
```

**After (Proper Error):**
```bash
$ vibeguard check --config ./vibeguard.yaml non-existent-check
Error: check with ID "non-existent-check" not found
$ echo $?
2
# Error message properly displayed
```

## Key Design Insight

The fix unifies error handling by ensuring all error paths follow the same pattern:
1. Return error from business logic
2. Cobra/main catches the error
3. `main.go` decides the exit code based on error type

This is more consistent with Go conventions and makes error handling testable.
