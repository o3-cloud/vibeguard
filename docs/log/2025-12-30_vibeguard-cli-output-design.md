---
summary: Defined CLI output patterns for VibeGuard following the "silence is success" principle, with examples for success, violation, verbose, and JSON output modes.
event_type: deep dive
sources:
  - docs/log/2025-12-30_vibeguard-implementation-patterns.md
  - docs/log/2025-12-30_vibeguard-architecture-spike.md
tags:
  - vibeguard
  - cli
  - ux
  - output-format
  - json
  - silence-is-success
---

# VibeGuard CLI Output Design

## Core Principle

**Silence is success.** When all policies pass, VibeGuard produces no output and exits with code 0. Output only appears when there are violations or when explicitly requested via `--verbose`.

## Output Modes

### 1. Default (Quiet) - Success

```
$ vibeguard check
$ echo $?
0
```

No output. Exit code 0 indicates all policies passed.

### 2. Default (Quiet) - Violation

```
$ vibeguard check
FAIL  coverage-threshold (error)
      Coverage is 72%, minimum required is 80%

      → Add unit tests for uncovered code paths in:
        - internal/executor/executor.go (68%)
        - internal/config/validate.go (45%)

$ echo $?
1
```

Output structure:
- **Status**: `FAIL` with policy ID and severity
- **Message**: Human-readable description of the violation
- **Action prompt**: Actionable next step prefixed with `→`

### 3. Verbose Mode

```
$ vibeguard check --verbose
✓ go-test          passed (1.2s)
✓ golangci-lint    passed (0.8s)
✗ coverage-threshold  FAIL
  Coverage is 72%, minimum required is 80%
```

Shows all tool executions with pass/fail status and duration.

### 4. JSON Output (Machine-Readable)

```
$ vibeguard check --json
{
  "violations": [{
    "policy": "coverage-threshold",
    "severity": "error",
    "message": "Coverage is 72%, minimum required is 80%",
    "suggested_action": "Add unit tests for uncovered code paths"
  }],
  "exit_code": 1
}
```

For tooling integration, CI/CD pipelines, and programmatic consumption.

### 5. JSON Output - Success

```
$ vibeguard check --json
{
  "violations": [],
  "exit_code": 0
}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All policies passed |
| 1 | One or more policy violations |
| 2 | Configuration error |
| 3 | Tool execution error |

## Design Rationale

1. **CI/CD friendly**: Zero output on success means clean logs
2. **Actionable**: Every violation includes a suggested next step
3. **Structured**: JSON mode enables integration with other tools
4. **Progressive disclosure**: Default is minimal, `--verbose` reveals details

## Future Considerations

- `--format=sarif` for IDE integration
- `--fix` to auto-apply suggested fixes where possible
- `--watch` for continuous checking during development
