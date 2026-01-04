---
summary: Research on making error exit codes configurable for FAIL and TIMEOUT cases
event_type: research
sources:
  - internal/executor/executor.go
  - internal/orchestrator/orchestrator.go
  - internal/cli/root.go
  - internal/cli/check.go
  - cmd/vibeguard/main.go
tags:
  - exit-codes
  - cli
  - configuration
  - error-handling
  - timeout
  - fail
---

# Configurable Exit Codes Research

Research on making error exit codes configurable, with FAIL and TIMEOUT cases using the same exit code, configurable via CLI flag, defaulting to exit code 1.

## Current Exit Code Architecture

### Exit Code Constants

Defined in `internal/executor/executor.go:13-20`:

```go
const (
    ExitCodeSuccess     = 0 // All checks passed
    ExitCodeConfigError = 2 // Configuration error
    ExitCodeViolation   = 3 // One or more error-severity violations
    ExitCodeTimeout     = 4 // Check execution timeout
)
```

### Exit Code Flow

1. **Orchestrator** (`internal/orchestrator/orchestrator.go:361-382`) - `calculateExitCode()` determines exit code based on violations
2. **CLI** (`internal/cli/check.go:86-87`) - Wraps exit code in `ExitError`
3. **Main** (`cmd/vibeguard/main.go:12-28`) - Calls `os.Exit()` with the code

### Current Logic

```go
func (o *Orchestrator) calculateExitCode(violations []*Violation) int {
    hasTimeout := false
    hasError := false

    for _, v := range violations {
        if v.Timedout {
            hasTimeout = true
        }
        if v.Severity == config.SeverityError {
            hasError = true
        }
    }

    if hasTimeout {
        return executor.ExitCodeTimeout  // Currently 4
    }
    if hasError {
        return executor.ExitCodeViolation  // Currently 3
    }
    return executor.ExitCodeSuccess
}
```

## Key Findings

### 1. FAIL and TIMEOUT Are Separate Exit Codes

- **FAIL** (error severity violation) → Exit code 3
- **TIMEOUT** → Exit code 4
- Timeout takes precedence over error severity

### 2. CLI Flag Pattern

Existing flags in `internal/cli/root.go:53-61`:

```go
rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file")
rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show all check results")
rootCmd.PersistentFlags().IntVarP(&parallel, "parallel", "p", 4, "Max parallel checks")
```

### 3. Exit Error Wrapper

`internal/cli/check.go:16-27` - Clean integration point for configurable exit codes:

```go
type ExitError struct {
    Code    int
    Message string
}
```

## Implementation Approach

### Option A: CLI Flag (Recommended)

Add `--error-exit-code` flag to unify FAIL and TIMEOUT exit codes:

```go
// In root.go
var errorExitCode int

rootCmd.PersistentFlags().IntVar(&errorExitCode, "error-exit-code", 1,
    "Exit code for FAIL and TIMEOUT errors (default 1)")
```

### Changes Required

1. **`internal/cli/root.go`** - Add `errorExitCode` flag variable and registration
2. **`internal/cli/check.go`** - Pass configured exit code to orchestrator or remap exit codes
3. **`internal/orchestrator/orchestrator.go`** - Modify `calculateExitCode()` to accept configurable error exit code
4. **Tests** - Update tests that expect specific exit codes (3, 4) to handle configurable values

### Exit Code Mapping

| Scenario | Current | Proposed (with `--error-exit-code=1`) |
|----------|---------|---------------------------------------|
| Success | 0 | 0 |
| Config Error | 2 | 2 |
| FAIL (error severity) | 3 | 1 (configurable) |
| TIMEOUT | 4 | 1 (configurable) |

### Key Design Decision

Unifying FAIL and TIMEOUT under one configurable exit code simplifies integration with CI/CD systems that only care about pass/fail, while preserving the distinction in output messages.

## Files to Modify

1. `internal/cli/root.go` - Add flag
2. `internal/cli/check.go` - Pass/use configured exit code
3. `internal/orchestrator/orchestrator.go` - Accept configurable exit code parameter
4. `internal/executor/executor.go` - Keep constants but may add helper
5. `docs/JSON-OUTPUT-SCHEMA.md` - Document new behavior
6. Test files - Update exit code assertions

## Next Steps

1. Create implementation plan
2. Add CLI flag `--error-exit-code` with default value 1
3. Modify orchestrator to use configurable exit code for both FAIL and TIMEOUT
4. Update tests
5. Update documentation
