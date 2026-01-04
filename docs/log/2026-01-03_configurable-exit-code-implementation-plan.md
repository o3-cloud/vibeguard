---
summary: Implementation plan for configurable exit codes - unifying FAIL and TIMEOUT under one CLI flag
event_type: code
sources:
  - docs/log/2026-01-03_configurable-exit-codes-research.md
  - docs/log/2025-12-30_vibeguard-specification.md
  - internal/executor/executor.go
  - internal/orchestrator/orchestrator.go
  - internal/cli/root.go
tags:
  - exit-codes
  - cli
  - implementation-plan
  - configuration
---

# Configurable Exit Code Implementation Plan

## Goal

Make the error exit code configurable via a CLI flag, with FAIL (error-severity violations) and TIMEOUT cases using the same exit code. Default to exit code 1.

## Current State

| Scenario | Current Exit Code |
|----------|-------------------|
| Success | 0 |
| Config Error | 2 |
| FAIL (error severity) | 3 |
| TIMEOUT | 4 |

## Target State

| Scenario | Target Exit Code |
|----------|------------------|
| Success | 0 (unchanged) |
| Config Error | 2 (unchanged) |
| FAIL (error severity) | 1 (configurable via `--error-exit-code`) |
| TIMEOUT | 1 (configurable via `--error-exit-code`) |

## Implementation Steps

### Step 1: Add CLI Flag

**File:** `internal/cli/root.go`

Add new flag variable and registration:

```go
var (
    // ... existing vars ...
    errorExitCode int
)

func init() {
    // ... existing flags ...
    rootCmd.PersistentFlags().IntVar(&errorExitCode, "error-exit-code", 1,
        "Exit code for check failures and timeouts (default 1)")
}
```

### Step 2: Export Flag Getter

**File:** `internal/cli/root.go`

Add getter function for orchestrator access:

```go
// GetErrorExitCode returns the configured exit code for failures
func GetErrorExitCode() int {
    return errorExitCode
}
```

### Step 3: Modify Orchestrator Interface

**File:** `internal/orchestrator/orchestrator.go`

Update `Orchestrator` struct to accept configurable exit code:

```go
type Orchestrator struct {
    config        *config.Config
    executor      Executor
    formatter     Formatter
    verbose       bool
    parallel      int
    failFast      bool
    errorExitCode int  // NEW: configurable exit code for failures
}

func New(cfg *config.Config, exec Executor, fmt Formatter, opts ...Option) *Orchestrator {
    o := &Orchestrator{
        config:        cfg,
        executor:      exec,
        formatter:     fmt,
        verbose:       false,
        parallel:      4,
        failFast:      false,
        errorExitCode: 1,  // NEW: default to 1
    }
    for _, opt := range opts {
        opt(o)
    }
    return o
}
```

### Step 4: Add Orchestrator Option

**File:** `internal/orchestrator/orchestrator.go`

Add option function for setting error exit code:

```go
// WithErrorExitCode sets the exit code used for failures and timeouts
func WithErrorExitCode(code int) Option {
    return func(o *Orchestrator) {
        o.errorExitCode = code
    }
}
```

### Step 5: Update calculateExitCode

**File:** `internal/orchestrator/orchestrator.go`

Modify to use configurable exit code:

```go
func (o *Orchestrator) calculateExitCode(violations []*Violation) int {
    for _, v := range violations {
        // Both timeout and error-severity violations use the configured exit code
        if v.Timedout || v.Severity == config.SeverityError {
            return o.errorExitCode
        }
    }
    return executor.ExitCodeSuccess
}
```

### Step 6: Update RunCheck Exit Code Logic

**File:** `internal/orchestrator/orchestrator.go`

Update the `RunCheck` method to use configurable exit code (around lines 488-508):

```go
// In RunCheck(), replace hardcoded exit codes with o.errorExitCode
if violation != nil {
    result.ExitCode = o.errorExitCode
} else {
    result.ExitCode = executor.ExitCodeSuccess
}
```

### Step 7: Wire Up in CLI

**File:** `internal/cli/check.go`

Pass the configured exit code to orchestrator:

```go
func runCheck(cmd *cobra.Command, args []string) error {
    // ... existing code ...

    orch := orchestrator.New(
        cfg,
        exec,
        formatter,
        orchestrator.WithVerbose(verbose),
        orchestrator.WithParallel(parallel),
        orchestrator.WithFailFast(failFast),
        orchestrator.WithErrorExitCode(errorExitCode),  // NEW
    )

    // ... rest of function ...
}
```

### Step 8: Update Executor Constants (Documentation Only)

**File:** `internal/executor/executor.go`

Add comment clarifying that `ExitCodeViolation` and `ExitCodeTimeout` are now deprecated in favor of configurable exit codes:

```go
const (
    ExitCodeSuccess     = 0 // All checks passed
    ExitCodeConfigError = 2 // Configuration error (config-time errors)

    // Deprecated: These are now configurable via --error-exit-code flag
    // Kept for backwards compatibility in tests and internal use
    ExitCodeViolation   = 3 // One or more error-severity violations
    ExitCodeTimeout     = 4 // Check execution error (timeout, command not found)
)
```

### Step 9: Update Tests

**Files to update:**
- `internal/cli/check_test.go`
- `internal/orchestrator/orchestrator_test.go`

Update tests to:
1. Test default exit code is 1
2. Test custom exit code via `WithErrorExitCode` option
3. Test that both FAIL and TIMEOUT use the same configured exit code

Example test:

```go
func TestConfigurableExitCode(t *testing.T) {
    tests := []struct {
        name          string
        errorExitCode int
        hasViolation  bool
        hasTimeout    bool
        wantExitCode  int
    }{
        {"default failure", 1, true, false, 1},
        {"default timeout", 1, false, true, 1},
        {"custom failure", 42, true, false, 42},
        {"custom timeout", 42, false, true, 42},
        {"success unchanged", 1, false, false, 0},
    }
    // ... test implementation ...
}
```

### Step 10: Update Documentation

**Files to update:**
- `docs/JSON-OUTPUT-SCHEMA.md` - Note that exit codes 3/4 are deprecated
- `docs/log/2025-12-30_vibeguard-specification.md` - Update exit code table
- `README.md` (if exists) - Document new flag

## Files Changed Summary

| File | Change Type |
|------|-------------|
| `internal/cli/root.go` | Add flag |
| `internal/cli/check.go` | Pass flag to orchestrator |
| `internal/orchestrator/orchestrator.go` | Add option, update logic |
| `internal/executor/executor.go` | Add deprecation comment |
| `internal/cli/check_test.go` | Update tests |
| `internal/orchestrator/orchestrator_test.go` | Add tests |
| `docs/JSON-OUTPUT-SCHEMA.md` | Update docs |

## Backwards Compatibility

- Default exit code changes from 3/4 to 1
- Users relying on exit codes 3 vs 4 to distinguish FAIL from TIMEOUT can use `--error-exit-code=3` and check output for timeout status
- JSON output still contains `"exit_code"` field reflecting actual exit code used

## Testing Strategy

1. Unit tests for `calculateExitCode` with various configurations
2. Integration tests verifying CLI flag is honored
3. Test that JSON output reflects configured exit code

## Estimated Changes

- ~50 lines of new code
- ~20 lines of test code
- ~10 lines of documentation updates
