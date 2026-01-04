---
summary: Implementation plan for logging check output to .vibeguard/log/ folder
event_type: code
sources:
  - internal/orchestrator/orchestrator.go
  - internal/executor/executor.go
  - internal/output/formatter.go
  - internal/output/json.go
tags:
  - implementation-plan
  - logging
  - ai-agents
  - check-output
  - developer-experience
  - simplification
---

# Check Output Logging Implementation Plan

Implementation plan for writing check command output (stdout + stderr) to `.vibeguard/log/<check-id>.log` files so AI agents can read failure details without re-running commands.

## Objective

Enable AI agents to inspect actual command output from failed (or passed) checks by writing the combined stdout + stderr to persistent log files and referencing those log files in violation output.

## Simplified Requirements

Based on user feedback to keep it simple:
- Write **only** the combined stdout + stderr (already captured in `executor.Result.Combined`)
- No timestamps, headers, or metadata - just the raw output
- Fixed path: `.vibeguard/log/<check-id>.log`
- Always overwrite (latest run only)
- Best-effort (don't fail checks if logging fails)

## Architecture

### Current State

The `executor.Execute()` method (internal/executor/executor.go:55-109) already captures:
- `Stdout` (separate)
- `Stderr` (separate)
- `Combined` (stdout + stderr concatenated at line 95)

This Combined output is already available but not persisted anywhere.

### Proposed Changes

**4 files to modify:**

1. **internal/orchestrator/orchestrator.go** - Add logging call + helper method
2. **internal/output/formatter.go** - Add "Log:" line to text output
3. **internal/output/json.go** - Add "log_file" field to JSON output
4. **.gitignore** - Exclude .vibeguard/ directory

## Detailed Implementation

### 1. Orchestrator: Add Logging (orchestrator.go)

**Location: After line 216** (immediately after check execution):

```go
// Execute the check
execResult, execErr := o.executor.Execute(checkCtx, check.ID, check.Run)
if cancel != nil {
    cancel()
}

// Write output to log file (best-effort, don't fail check if this fails)
if err := o.writeCheckLog(check.ID, execResult.Combined); err != nil {
    // Silently continue - logging is best-effort
}
```

**New helper method** (add to orchestrator.go):

```go
// writeCheckLog writes check output to .vibeguard/log/<check-id>.log
func (o *Orchestrator) writeCheckLog(checkID, output string) error {
    logDir := ".vibeguard/log"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return err
    }

    logPath := filepath.Join(logDir, checkID+".log")
    return os.WriteFile(logPath, []byte(output), 0644)
}
```

**Update Violation struct** (line 38-47):

```go
type Violation struct {
    CheckID    string
    Severity   config.Severity
    Command    string
    Suggestion string
    Fix        string
    Extracted  map[string]string
    Timedout   bool
    LogFile    string  // NEW: path to log file for AI agents
}
```

**Set LogFile when creating violations** (line 303-312):

```go
violation := &Violation{
    CheckID:    check.ID,
    Severity:   check.Severity,
    Command:    check.Run,
    Suggestion: suggestion,
    Fix:        check.Fix,
    Extracted:  result.Extracted,
    Timedout:   execResult.Timedout,
    LogFile:    filepath.Join(".vibeguard/log", check.ID+".log"),  // NEW
}
```

### 2. Text Formatter: Show Log Path (formatter.go)

**Location: After line 154** (after Fix line, before Advisory):

```go
// Show fix if present, otherwise fallback to run command
if v.Fix != "" {
    fix := config.InterpolateWithExtracted(v.Fix, nil, v.Extracted)
    _, _ = fmt.Fprintf(f.out, "  Fix: %s\n", fix)
} else if v.Suggestion == "" {
    _, _ = fmt.Fprintf(f.out, "  Fix: %s\n", v.Command)
}

// NEW: Show log file location
if v.LogFile != "" {
    _, _ = fmt.Fprintf(f.out, "  Log: %s\n", v.LogFile)
}

// Show advisory line
advisory := "blocks commit"
```

### 3. JSON Formatter: Add log_file Field (json.go)

**Update JSONViolation struct** (line 25-33):

```go
type JSONViolation struct {
    ID         string            `json:"id"`
    Severity   string            `json:"severity"`
    Command    string            `json:"command"`
    Suggestion string            `json:"suggestion,omitempty"`
    Fix        string            `json:"fix,omitempty"`
    Extracted  map[string]string `json:"extracted,omitempty"`
    LogFile    string            `json:"log_file,omitempty"`  // NEW
}
```

**Populate LogFile in output** (line 58-66):

```go
for _, v := range result.Violations {
    output.Violations = append(output.Violations, JSONViolation{
        ID:         v.CheckID,
        Severity:   string(v.Severity),
        Command:    v.Command,
        Suggestion: v.Suggestion,
        Fix:        v.Fix,
        Extracted:  v.Extracted,
        LogFile:    v.LogFile,  // NEW
    })
}
```

### 4. Gitignore: Exclude Logs (.gitignore)

```
.vibeguard/
```

## Expected Output

### Text Format

```
FAIL  lint-go (error)

  Fix linting errors in internal/cli/check.go
  Fix: golangci-lint run ./... --fix
  Log: .vibeguard/log/lint-go.log
  Advisory: This check blocks commit
```

### JSON Format

```json
{
  "violations": [
    {
      "id": "lint-go",
      "severity": "error",
      "command": "golangci-lint run ./...",
      "fix": "golangci-lint run ./... --fix",
      "log_file": ".vibeguard/log/lint-go.log"
    }
  ]
}
```

### Log File Content (.vibeguard/log/lint-go.log)

```
internal/cli/check.go:45:2: ineffective assignment to `err` (ineffassign)
        err = nil
        ^
level=warning msg="[runner] Timeout exceeded..."
```

Just the raw combined stdout + stderr - no headers, timestamps, or formatting.

## Implementation Notes

### Design Decisions

1. **Log all checks (pass and fail)**: AI agents may want to inspect passing checks too
2. **Best-effort**: Don't fail checks if logging fails (disk full, permissions, etc.)
3. **No configuration**: Fixed path keeps it simple
4. **Overwrite always**: No rotation/cleanup logic needed
5. **Raw output only**: No metadata, just stdout + stderr

### Edge Cases Handled

- **Missing directory**: `os.MkdirAll` creates it with 0755 permissions
- **Permission errors**: Silently ignored (best-effort)
- **Empty output**: Empty file written (still signals check ran)
- **Concurrent writes**: Each check has unique filename (check.ID)
- **Long check IDs**: Check IDs validated in config, safe for filenames

### What This Enables

When an AI agent sees a violation like:

```
FAIL  lint-go (error)
  Fix linting errors
  Log: .vibeguard/log/lint-go.log
```

It can immediately read `.vibeguard/log/lint-go.log` to see the actual linter output without:
- Re-running the (potentially slow) command
- Risking different results if state changed
- Needing to know how to invoke the command

## Testing Strategy

1. Run `vibeguard check` with a failing check
2. Verify `.vibeguard/log/<check-id>.log` exists and contains output
3. Run with passing check, verify log still written
4. Verify text formatter shows "Log:" line
5. Verify JSON formatter includes "log_file" field
6. Test with check that times out
7. Test with check that has no output (empty log file)

## Files Modified Summary

| File | Change | Lines |
|------|--------|-------|
| internal/orchestrator/orchestrator.go | Add writeCheckLog() method | +10 |
| internal/orchestrator/orchestrator.go | Call writeCheckLog() after execute | +4 |
| internal/orchestrator/orchestrator.go | Add LogFile to Violation struct | +1 |
| internal/orchestrator/orchestrator.go | Set LogFile in violation creation | +1 |
| internal/output/formatter.go | Add Log: line to output | +3 |
| internal/output/json.go | Add log_file to JSONViolation | +1 |
| internal/output/json.go | Set LogFile in JSON output | +1 |
| .gitignore | Exclude .vibeguard/ | +1 |

**Total: ~22 lines of new/modified code**

## Next Steps

1. Implement changes in order listed above
2. Test with vibeguard's own checks (dogfooding)
3. Verify AI agent can read logs during pre-commit hook failure
4. Consider future enhancement: optional JSON log format

## Related

- Research: [docs/log/2026-01-03_check-output-logging-research.md](./2026-01-03_check-output-logging-research.md)
- ADR-006: Git pre-commit hook integration
