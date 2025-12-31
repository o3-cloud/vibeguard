---
summary: Implemented fail-fast mode with context cancellation for vibeguard orchestrator (vibeguard-v3m.4)
event_type: code
sources:
  - internal/orchestrator/orchestrator.go
  - internal/executor/executor.go
  - internal/output/formatter.go
  - internal/output/json.go
  - internal/orchestrator/orchestrator_test.go
tags:
  - orchestrator
  - fail-fast
  - context-cancellation
  - parallel-execution
  - vibeguard-v3m
  - phase-3
---

# Fail-Fast Mode Implementation (vibeguard-v3m.4)

Implemented complete fail-fast mode with context cancellation to stop in-flight checks when an error-severity failure occurs. This completes task vibeguard-v3m.4 from Phase 3: Orchestration.

## Background

The fail-fast feature was partially implemented in vibeguard-v3m.3 (parallel execution), but had limitations:
1. In-flight checks continued running until completion instead of being cancelled
2. No indication in results that fail-fast was triggered
3. No user-facing output to show that execution stopped early

## Implementation Details

### Core Changes

1. **Context Cancellation for In-Flight Checks** (`orchestrator.go`)
   - Added `failFastCtx, cancelFailFast := context.WithCancel(ctx)` to create a cancellable context
   - When fail-fast triggers, `cancelFailFast()` is called to cancel all in-flight checks
   - The errgroup's context is derived from `failFastCtx`, so cancellation propagates to all goroutines
   - Handled `context.Canceled` error gracefully when fail-fast is triggered

2. **FailFastTriggered Field in RunResult** (`orchestrator.go`)
   - Added `FailFastTriggered bool` field to `RunResult` struct
   - Allows callers to know if execution was stopped early due to fail-fast
   - Set to `true` when an error-severity check fails with fail-fast enabled

3. **Cancelled Field in Executor Result** (`executor.go`)
   - Added `Cancelled bool` field to `Result` struct
   - Distinguishes between timeout (DeadlineExceeded) and cancellation (Canceled)
   - Exit code set to -1 for cancelled checks (vs 3 for timeout)
   - Updated `String()` method to show "cancelled" status

4. **Output Formatter Updates** (`formatter.go`, `json.go`)
   - Human-readable output shows "Execution stopped early due to --fail-fast" when triggered
   - Verbose mode shows cancelled checks with "⊘" symbol
   - JSON output includes `fail_fast_triggered` field and "cancelled" status for checks

### Key Design Decisions

1. **Error Severity Only** - Fail-fast only triggers on error-severity violations, not warnings
2. **Level Boundary Behavior** - Checks in the same level as the failing check may or may not complete (race condition), but subsequent levels are skipped
3. **Graceful Cancellation** - Cancelled checks return cleanly with a special status, no errors propagated

### Test Coverage

Added 8 new tests for fail-fast functionality:

**Orchestrator Tests:**
- `TestRun_FailFast_SetsFailFastTriggeredFlag` - Verifies flag is set when fail-fast triggers
- `TestRun_NoFailFast_FailFastTriggeredFalse` - Flag stays false when fail-fast is disabled
- `TestRun_FailFast_WarningSeverityDoesNotTrigger` - Warning severity doesn't trigger fail-fast
- `TestRun_FailFast_CancelsLongRunningChecks` - Verifies long-running checks are cancelled quickly
- `TestRun_FailFast_AllChecksPassDoesNotTrigger` - Flag stays false when all checks pass

**Executor Tests:**
- `TestExecute_ContextCancelled_SetsCancelledFlag` - Cancelled context sets flag
- `TestExecute_NormalCompletion_CancelledFlagFalse` - Normal completion keeps flag false
- `TestResult_String_Cancelled` - String representation shows "cancelled"

## Findings

1. **Existing Infrastructure** - The `--fail-fast` flag and basic logic were already wired up; the main work was adding context cancellation and visibility
2. **errgroup Context Behavior** - The errgroup's `WithContext()` creates a context that cancels when any goroutine returns an error, but we needed our own cancel function for fail-fast since goroutines return `nil` after triggering fail-fast
3. **Distinction Between Timeout and Cancellation** - Important to distinguish these two cases in executor results for proper reporting

## Files Changed

- `internal/orchestrator/orchestrator.go` - Context cancellation and FailFastTriggered field
- `internal/executor/executor.go` - Cancelled field and context.Canceled handling
- `internal/output/formatter.go` - Human-readable fail-fast output
- `internal/output/json.go` - JSON fail-fast output
- `internal/orchestrator/orchestrator_test.go` - 5 new fail-fast tests
- `internal/executor/executor_test.go` - 3 new cancellation tests

## Next Steps

- Close vibeguard-v3m.4 task
- Phase 3: Orchestration is now complete (all 5 tasks done)
- Continue with Phase 2: Grok + Assertions (vibeguard-c9m)
