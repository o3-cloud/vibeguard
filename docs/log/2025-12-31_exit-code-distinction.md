---
summary: Completed vibeguard-8p1: Distinguish config errors (exit 2) from execution errors (exit 3)
event_type: code
sources:
  - vibeguard-8p1 Beads task
  - docs/log/2025-12-31_error-message-analysis.md
  - internal/executor/executor.go
  - internal/orchestrator/orchestrator.go
tags:
  - exit-codes
  - bug-fix
  - vibeguard-8p1
  - error-handling
  - ci-cd
---

# Exit Code Distinction Fix

## Summary

Fixed the exit code mapping in VibeGuard to properly distinguish between config-time errors and execution errors, aligning with the specification defined in vibeguard-8p1.

## Changes Made

### 1. Swapped Exit Code Constants (internal/executor/executor.go)
- **Before:**
  - Exit code 2 = `ExitCodeViolation` (check failed)
  - Exit code 3 = `ExitCodeConfigError` (config error)
  - Exit code 4 = `ExitCodeTimeout` (timeout)

- **After:**
  - Exit code 2 = `ExitCodeConfigError` (config-time errors)
  - Exit code 3 = `ExitCodeViolation` (execution errors)
  - Exit code 4 = `ExitCodeTimeout` (timeout)

This change aligns with the project's semantic versioning of exit codes where config-time errors (which block the entire pipeline) should have a lower exit code than execution errors.

### 2. Updated Test Names and Comments (internal/orchestrator/orchestrator_test.go)
- Renamed `TestRun_FailingCheck_ErrorSeverity_ExitCodeTwo` → `TestRun_FailingCheck_ErrorSeverity_ExitCodeThree`
- Renamed `TestRun_MultipleChecks_ErrorFailure_ExitCodeTwo` → `TestRun_MultipleChecks_ErrorFailure_ExitCodeThree`
- Renamed `TestRunCheck_UnknownCheck_ExitCodeThree` → `TestRunCheck_UnknownCheck_ExitCodeTwo`
- Updated comments to reflect the correct exit code semantics

### 3. Test Coverage
- All existing tests pass with the updated constants
- Tests automatically use the correct exit codes because they reference the constants rather than hardcoded values
- No test logic needed to change; the constants handled the mapping

## Exit Code Semantics

The final exit code mapping now correctly represents:

- **Exit 0** - All checks passed successfully
- **Exit 2** - Configuration error (e.g., unknown check ID, invalid config)
- **Exit 3** - One or more error-severity violations detected during execution
- **Exit 4** - Check execution error (timeout, command not found)

## Rationale

This change improves CI/CD integration by making exit codes semantically meaningful:
- Config errors (exit 2) occur before any checks run and indicate a problem with the configuration itself
- Execution errors (exit 3) occur when checks run but fail to meet their criteria
- This distinction allows CI/CD systems to handle these failure modes differently

## Compatibility

The change is compatible with Claude Code hook expectations where exit codes ≥ 2 are blocking, allowing the tool to function as a pre-commit hook or CI/CD check with proper error categorization.

## Test Results

```
✓ All tests pass (70+ tests across the codebase)
✓ Internal/executor tests: PASS
✓ Internal/orchestrator tests: PASS
✓ All other tests: PASS
✓ Build successful
```

## Related Tasks

- Closes: vibeguard-8p1
- Documented in: docs/log/2025-12-31_error-message-analysis.md
