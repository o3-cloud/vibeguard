---
summary: Research on logging check output to .vibeguard folder for AI agent consumption
event_type: research
sources:
  - internal/orchestrator/orchestrator.go
  - internal/executor/executor.go
  - internal/output/formatter.go
  - internal/output/json.go
  - internal/cli/check.go
tags:
  - logging
  - ai-agents
  - check-output
  - developer-experience
  - cli
  - output-capture
---

# Check Output Logging Research

Research into logging vibeguard check output to a `.vibeguard` folder so AI agents can read failure logs instead of rerunning commands.

## Problem Statement

When a vibeguard check fails during a Claude Code hook, the AI agent currently sees only the formatted violation output (suggestion, fix, severity). To understand the actual failure details, the agent must rerun the command manually, which:

1. Wastes time re-executing potentially slow checks
2. May produce different results if state has changed
3. Requires the agent to know how to invoke the command correctly

## Current Architecture

### Output Capture (executor.go:55-109)

Commands are executed with full output capture:

```go
type Result struct {
    CheckID   string
    ExitCode  int
    Stdout    string          // Separate stdout
    Stderr    string          // Separate stderr
    Combined  string          // stdout + stderr
    Duration  time.Duration
    Success   bool
    Timedout  bool
    Cancelled bool
    Error     error
}
```

The `Combined` field contains `stdout + stderr` (line 95).

### Output Analysis (orchestrator.go:75-86)

The `getAnalysisOutput()` method determines what to analyze:
- If `check.File` is specified: reads that file
- Otherwise: uses `execResult.Combined`

### Output Formatting (formatter.go, json.go)

Two output modes exist:
- **Text formatter**: Human-readable, shows violations with suggestion/fix/advisory
- **JSON formatter**: Machine-readable, includes extracted values

Neither currently persists the raw command output.

## Proposed Solution

### Log Directory Structure

```
.vibeguard/
└── logs/
    ├── lint-go.log
    ├── test-unit.log
    └── build-check.log
```

Each log file is named `<check-id>.log` and contains the most recent execution output for that check.

### Log File Content

```
=== Check: lint-go ===
Timestamp: 2026-01-03T14:32:15Z
Command: golangci-lint run ./...
Exit Code: 1
Duration: 4.532s

=== STDOUT ===
internal/cli/check.go:45:2: ineffective assignment to `err` (ineffassign)
        err = nil
        ^

=== STDERR ===
level=warning msg="[runner] Timeout exceeded..."

=== COMBINED ===
[same as above, interleaved]

=== EXTRACTED VALUES ===
file: internal/cli/check.go
line: 45
message: ineffective assignment to `err`

=== ASSERTION ===
Expression: exit_code == 0
Result: false
```

### Violation Output Enhancement

Add a `Log:` line to violation output:

```
FAIL  lint-go (error)
      Suggestion: Fix linting errors in internal/cli/check.go
      Fix: golangci-lint run ./... --fix
      Log: .vibeguard/logs/lint-go.log
      Advisory: This check blocks commit
```

### JSON Output Enhancement

Add `log_file` field to violation objects:

```json
{
  "violations": [
    {
      "id": "lint-go",
      "severity": "error",
      "log_file": ".vibeguard/logs/lint-go.log",
      ...
    }
  ]
}
```

## Implementation Plan

### Phase 1: Core Logging Infrastructure

1. **Create logger package** (`internal/logger/`)
   - `Logger` struct with configurable log directory
   - `WriteCheckLog(checkID string, result executor.Result, extracted map[string]string)` method
   - Auto-create `.vibeguard/logs/` directory

2. **Integration point** (orchestrator.go)
   - After check execution, before result aggregation
   - Log on failure only (default) or always (with flag)

### Phase 2: Output Enhancement

3. **Update Violation struct** (orchestrator.go:38-47)
   - Add `LogFile string` field

4. **Update formatters**
   - Text: Add "Log:" line after "Fix:" line
   - JSON: Add `log_file` field to JSONViolation

### Phase 3: Configuration

5. **CLI flags** (cli/check.go)
   - `--log-output` - Enable logging (default: true for failures)
   - `--log-dir` - Custom log directory (default: `.vibeguard/logs`)
   - `--log-all` - Log passed checks too

6. **Config file support** (config/schema.go)
   - `logging.enabled: bool`
   - `logging.directory: string`
   - `logging.log_passed: bool`

### Phase 4: Cleanup

7. **Log rotation/cleanup**
   - Keep only last N runs per check, or
   - Clean logs older than N days, or
   - Simple: always overwrite (latest only)

## Design Decisions

### Why `.vibeguard/` folder?

- Project-local, doesn't pollute home directory
- Similar to `.git/`, `.vscode/`, etc.
- Should be gitignored (ephemeral data)
- Easy for AI agents to find relative to project root

### Why overwrite (latest only)?

- Simplest implementation
- Prevents unbounded disk growth
- AI agents only care about most recent failure
- Historical logs can be added later if needed

### Why log failures by default?

- Most common use case is debugging failures
- Reduces noise and disk usage
- Can enable for all checks with `--log-all` flag

### Why plain text logs?

- Human-readable for debugging
- Can be parsed by AI agents easily
- Simple to implement
- Consider JSON logs as alternative/option

## Alternatives Considered

### Alternative 1: In-Memory Cache

Store output in memory and expose via CLI subcommand (`vibeguard show-log <check-id>`).

**Pros:**
- No disk I/O
- Simple implementation

**Cons:**
- Lost when process exits
- Doesn't work across hook invocations
- Requires additional CLI complexity

### Alternative 2: SQLite Database

Store logs in `.vibeguard/logs.db`.

**Pros:**
- Query capabilities
- Better storage efficiency
- Can store history easily

**Cons:**
- Added dependency
- More complex
- Overkill for this use case

### Alternative 3: Structured JSON Logs

Use JSON instead of plain text.

**Pros:**
- Machine-readable
- Consistent parsing

**Cons:**
- Less human-readable
- More complex to generate
- Could offer both formats as option

## Next Steps

1. Create implementation plan ADR (optional)
2. Implement Phase 1 (core logging)
3. Test with Claude Code hook workflow
4. Iterate based on AI agent feedback

## Related

- ADR-006: Git pre-commit hook integration
- executor.Result struct captures all needed data
- File field already shows pattern for output redirection
