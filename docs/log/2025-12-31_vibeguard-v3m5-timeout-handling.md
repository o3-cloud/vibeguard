---
summary: Implemented timeout handling for vibeguard checks with exit code 3
event_type: code
sources:
  - internal/executor/executor.go
  - internal/orchestrator/orchestrator.go
  - internal/executor/executor_test.go
  - internal/orchestrator/orchestrator_test.go
tags:
  - vibeguard
  - timeout
  - executor
  - orchestrator
  - phase-3
---

# Timeout Handling Implementation (vibeguard-v3m.5)

Completed the timeout handling task for vibeguard, implementing per-check timeout with exit code 3 reporting.

## Task Requirements

From beads task `vibeguard-v3m.5`:
- Implement per-check timeout field (default: 30s) - **Already existed**
- Support s/m/h units - **Already existed via time.ParseDuration**
- Kill timed-out processes and report exit code 3 - **Implemented**
- Use context.WithTimeout - **Already existed**

## Changes Made

### Executor (`internal/executor/executor.go`)

1. Added `ExitCodeTimeout = 3` constant
2. Added `Timedout bool` field to `Result` struct
3. Updated `Execute()` to detect `context.DeadlineExceeded` and:
   - Set `Timedout = true`
   - Set `ExitCode = 3`
4. Updated `String()` method to display "timeout" status for timed-out results

### Orchestrator (`internal/orchestrator/orchestrator.go`)

1. Added `Timedout bool` field to `Violation` struct
2. Updated `calculateExitCode()` to return exit code 3 when any violation has timeout (takes precedence over error severity)
3. Updated violation creation in `Run()` and `RunCheck()` to:
   - Set `Timedout` flag from execution result
   - Override suggestion with timeout-specific message

## Exit Code Precedence

| Condition | Exit Code |
|-----------|-----------|
| All checks pass | 0 |
| Error severity violation | 1 |
| Configuration error | 2 |
| Timeout | 3 (takes precedence) |

## Tests Added

### Executor Tests
- `TestExecute_Timeout_SetsExitCode3` - Verifies exit code 3 on timeout
- `TestExecute_Timeout_SetsTimedoutFlag` - Verifies Timedout flag is set
- `TestExecute_NoTimeout_TimedoutFlagFalse` - Verifies flag is false for completed commands
- `TestExecute_NonZeroExit_NotTimeout` - Verifies non-timeout failures keep original exit code
- `TestResult_String_Timeout` - Verifies String() shows "timeout" status

### Orchestrator Tests
- `TestRun_Timeout_ExitCode3` - Exit code 3 for timed out checks
- `TestRun_Timeout_ViolationMarkedAsTimeout` - Violation has Timedout flag
- `TestRun_Timeout_SuggestionIncludesTimeoutMessage` - Custom timeout suggestion
- `TestRun_TimeoutTakesPrecedenceOverError` - Timeout exit code takes precedence
- `TestRunCheck_Timeout_ExitCode3` - Single check timeout returns exit code 3
- `TestRunCheck_Timeout_ViolationMarkedAsTimeout` - Single check violation marked

## Key Findings

The timeout infrastructure was already well-established:
- `Check.Timeout` field existed with `Duration` type
- Default timeout of 30s was already applied
- `context.WithTimeout` was already used in orchestrator
- `exec.CommandContext` properly respected context cancellation

The main gap was **detecting and properly reporting** timeout conditions:
- Context cancellation resulted in signal-based exit codes (-1 or similar)
- No way to distinguish timeout from other failures
- No exit code 3 as specified in requirements

## Implementation Notes

- Timeout detection uses `ctx.Err() == context.DeadlineExceeded` after command execution
- This approach is more reliable than checking exit codes since killed processes can return various codes
- The timeout suggestion overrides user-defined suggestions since the original command didn't complete
