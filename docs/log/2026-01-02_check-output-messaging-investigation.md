---
summary: Investigated vibeguard check output formatting - proposed adding `fix` field to schema, always showing Advisory line, and using WARN/FAIL headers based on severity
event_type: research
sources:
  - internal/output/formatter.go
  - internal/orchestrator/orchestrator.go
  - internal/config/config.go
  - internal/config/interpolate.go
  - vibeguard.yaml
tags:
  - cli
  - output-formatting
  - user-experience
  - vibeguard-check
  - ai-agent
---

# Vibeguard Check Output Messaging Investigation

Investigated the current output messaging of `vibeguard check` to identify how to make it more direct and explicit about what actions users should take.

## Current Output Format

The current output for a failed check looks like:

```text
FAIL  lint (error)
      > golangci-lint run ./...

      Tip: golangci-lint found lint issues, run 'golangci-lint run ./...' to fix lint issues
```

## Key Files

| File                                     | Purpose                                                                          |
| ---------------------------------------- | -------------------------------------------------------------------------------- |
| `internal/output/formatter.go`           | Main formatting logic for check results - handles both verbose and quiet output |
| `internal/orchestrator/orchestrator.go`  | Creates `Violation` structs with suggestions; contains two creation sites       |
| `internal/config/config.go`              | Check configuration and default severity handling                               |
| `internal/config/interpolate.go`         | Template interpolation for suggestion text with extracted values                |
| `vibeguard.yaml`                         | Check definitions with suggestion text                                          |

## Output Flow

### Violation Creation Sites

There are **two** places where violations are created in the orchestrator:

1. **Main Run() function** (`orchestrator.go:254-273`):
   - Creates violations during parallel check execution
   - Lines 254-267: Creates `Violation` struct with CheckID, Severity, Command, Suggestion, Extracted, Timedout
   - Lines 269-272: Triggers fail-fast cancellation for error severity violations

2. **RunCheck() function** (`orchestrator.go:430-443`):
   - Creates violations for single-check execution mode
   - Same structure as main Run() but for isolated check runs

Both sites use the same timeout suggestion fallback:
```go
suggestion = "Check timed out. Consider increasing the timeout value or optimizing the command."
```

### Output Formatting

The formatter has **two output modes** that must be kept consistent:

1. **Quiet Mode** - `formatViolation()` (`formatter.go:75-88`):
   - Line 77: `FAIL  %s (%s)` - header with check ID and severity
   - Line 78: `> %s` - truncated command (60 chars max)
   - Line 85: `Tip: %s` - suggestion with interpolated values

2. **Verbose Mode** - `formatVerbose()` (`formatter.go:48-73`):
   - Line 57: `✗ %-15s FAIL (%.1fs)` - check ID with duration
   - Line 66: `%s` - suggestion **without** "Tip:" prefix
   - Shows all checks (passed, failed, cancelled), not just violations

Note: Verbose mode outputs the suggestion directly without the "Tip:" prefix, which is inconsistent with quiet mode.

### Template Interpolation

The `InterpolateWithExtracted()` function (`interpolate.go:37-45`) renders Go template strings with:
- Config vars (from `vibeguard.yaml` vars section)
- Extracted values from grok patterns (e.g., `{{.coverage}}`)
- Config vars take precedence over extracted values if there's a conflict

This enables dynamic suggestions like "Coverage is {{.coverage}}%, need 80%".

## Exit Codes

The orchestrator uses specific exit codes (`orchestrator.go:313-337`):

| Exit Code | Meaning   | Use Case                                                                       |
| --------- | --------- | ------------------------------------------------------------------------------ |
| 0         | Success   | All checks passed                                                              |
| 2         | Violation | One or more checks failed (uses 2, not 1, for Claude Code hook compatibility) |
| 4         | Timeout   | One or more checks timed out (takes precedence over error severity)           |

AI agents can use these exit codes to determine action without parsing output.

## Default Severity Behavior

When a check omits the `severity` field, it defaults to `error` (`config.go:141-142`):

```go
if c.Checks[i].Severity == "" {
    c.Checks[i].Severity = SeverityError
}
```

**Implication**: The `test-coverage` check in `vibeguard.yaml` (lines 31-37) has no explicit severity, so it defaults to `error` and will block commits.

## AI Agent Perspective

When an AI agent (like Claude Code) encounters a failed check, it needs:

1. **What to run to fix** - The exact command to remediate the issue
2. **Severity impact** - Does this block (error) or is it advisory (warning)?
3. **What failed** - For assertion failures, what value was found vs expected

### Distinguishing "What Failed" vs "How to Fix"

Current confusion: The document originally conflated the run command with the fix command.

**Different check types have different fix patterns:**

| Check Type | Example         | What Failed                            | How to Fix                                        |
| ---------- | --------------- | -------------------------------------- | ------------------------------------------------- |
| Exit code  | `lint`          | `golangci-lint run ./...` exited non-0 | Run the same command, fix reported issues         |
| Assertion  | `test-coverage` | Coverage was 72%, needed 80%           | Add more tests (not re-running the same command)  |
| Timeout    | `test`          | Command exceeded time limit            | Increase timeout or fix hanging tests             |

For exit-code checks, the fix command is often the same as what failed. For assertion checks, the suggestion should explain what failed, and the fix should explain what to do about it.

### What AI Agents Don't Need

- "Tip:" prefix (agents know it's guidance)
- Explanatory prose around the command

Note: We intentionally add `Advisory:` line even though `FAIL`/`WARN` is in the header. This makes blocking status explicit without requiring agents to know the FAIL=blocks, WARN=doesn't mapping.

## Recommended Format

### Errors (blocking)

```text
FAIL  lint (error)

  Fix: golangci-lint run ./...
  Advisory: blocks commit
```

### Warnings (non-blocking)

```text
WARN  test (warning)

  Fix: go test ./...
  Advisory: does not block commit
```

Always show "Advisory:" line for both errors and warnings to make blocking status explicit.

### Timeouts

```text
FAIL  test (timeout after 10s)

  Fix: go test ./... -v
  Note: Consider increasing timeout or check for hanging tests
  Advisory: blocks commit
```

### Assertion failures

```text
FAIL  test-coverage (error)

  Coverage is 72%, need 80%
  Fix: Add tests to improve coverage
  Advisory: blocks commit
```

Extracted values (like `{{.coverage}}`) are formatted directly in the `suggestion` field. No separate "Got/Need" display - users control the message format.

## Severity Distinction

| Severity | Header | Advisory Line                    | Implication                |
| -------- | ------ | -------------------------------- | -------------------------- |
| error    | `FAIL` | `Advisory: blocks commit`        | Must fix before proceeding |
| warning  | `WARN` | `Advisory: does not block commit`| Can proceed, fix later     |

Key insight: Use `WARN` instead of `FAIL` for warnings, and always include `Advisory:` line to make blocking status explicit.

## Schema Decision: Add `fix` Field

**Problem**: The `suggestion` field has ambiguous semantics:
- For exit-code checks: often redundantly duplicates the `run` command
- For assertion checks: explains what failed, but doesn't say how to fix it

**Decision**: Add a new `fix` field to separate concerns:

| Field | Purpose | Example |
| ----- | ------- | ------- |
| `suggestion` | Explains what failed (can use `{{.extracted}}` values) | "Coverage is {{.coverage}}%, need 80%" |
| `fix` | Action to remediate (command or guidance) | "Add tests to improve coverage" |

### Schema Changes

```yaml
# Assertion check - both fields useful
- id: test-coverage
  run: go test ./... -coverprofile cover.out && go tool cover -func cover.out
  grok:
    - total:.*\(statements\)\s+%{NUMBER:coverage}%
  assert: "coverage >= 80"
  suggestion: "Coverage is {{.coverage}}%, need 80%"  # What failed
  fix: "Add tests to improve coverage"                # How to fix

# Exit-code check - suggestion optional, fix is the command
- id: lint
  run: golangci-lint run ./...
  severity: error
  fix: "golangci-lint run ./..."  # Just the fix command
  # suggestion omitted - failure is obvious from command output
```

### Formatter Output

**Error (blocking):**
```text
FAIL  lint (error)

  Fix: golangci-lint run ./...
  Advisory: blocks commit
```

**Warning (non-blocking):**
```text
WARN  test (warning)

  Fix: go test ./...
  Advisory: does not block commit
```

**Assertion failure with suggestion:**
```text
FAIL  test-coverage (error)

  Coverage is 72%, need 80%
  Fix: Add tests to improve coverage
  Advisory: blocks commit
```

Note: Extracted values like `{{.coverage}}` are interpolated in the `suggestion` field. No auto-generated "Got/Need" - users control the message.

### Backward Compatibility

- `suggestion` remains supported for existing configs
- If only `suggestion` provided (no `fix`), show `suggestion` (no `Fix:` prefix)
- If both provided, show `suggestion` then `Fix: {fix}`
- If only `fix` provided, just show `Fix: {fix}`
- `Advisory:` line is always added based on severity (new behavior)

## Open Questions

1. **Should we add `--output=json` for machine-parseable output?** This may be more useful for AI agents than optimizing text output.

2. **What happens with empty suggestion AND fix?** Should we show a default based on the `run` command?

## Backward Compatibility Considerations

Changing output format may break:
- CI scripts parsing output
- Existing Claude Code hooks expecting specific patterns
- User muscle memory for recognizing output

Consider:
- Adding `--format=v2` flag for new format
- Deprecation period with warnings
- Documentation of format changes

## Acceptance Criteria

The improvement is "done" when:

- [ ] Schema updated: Add `fix` field to `config.Check` struct
- [ ] Schema validation: Allow `suggestion`, `fix`, or both
- [ ] Violation struct updated to carry both `Suggestion` and `Fix`
- [ ] `formatViolation()` uses `WARN` vs `FAIL` based on severity
- [ ] `formatViolation()` outputs `suggestion` then `Fix: {fix}` appropriately
- [ ] `formatViolation()` always shows `Advisory:` line (errors: "blocks commit", warnings: "does not block commit")
- [ ] `formatVerbose()` updated consistently with quiet mode
- [ ] Both violation creation sites pass the new `fix` field
- [ ] Exit codes remain unchanged (0/2/4)
- [ ] Tests updated for new schema and output format
- [ ] `vibeguard.yaml` updated to use new `fix` field with `{{.extracted}}` in suggestions

## Next Steps

### High Priority
1. Add `fix` field to `config.Check` struct in `schema.go`
2. Add `Fix` field to `Violation` struct in `orchestrator.go`
3. Update both violation creation sites to populate `Fix`
4. Update `formatViolation()` to handle `suggestion` + `fix` output
5. Update `formatVerbose()` consistently

### Medium Priority
6. Update `formatViolation()` to use `WARN` vs `FAIL` based on severity
7. Update `vibeguard.yaml` checks to use new schema
8. Add tests for new output format

### Low Priority
9. Consider `--output=json` flag for structured output
10. Document schema changes in user documentation
