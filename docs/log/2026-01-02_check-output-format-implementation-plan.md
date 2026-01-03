---
summary: Implementation plan for vibeguard check output formatting improvements - add fix field, WARN/FAIL headers, and Advisory line
event_type: deep dive
sources:
  - docs/log/2026-01-02_check-output-messaging-investigation.md
  - internal/output/formatter.go
  - internal/orchestrator/orchestrator.go
  - internal/config/config.go
  - internal/output/json.go
tags:
  - implementation-plan
  - cli
  - output-formatting
  - vibeguard-check
  - ai-agent
  - schema-change
---

# Check Output Format Implementation Plan

Based on the investigation in `2026-01-02_check-output-messaging-investigation.md`, this plan outlines the implementation steps to improve vibeguard check output for AI agents.

## Goal

Make vibeguard check output more actionable for AI agents by:
1. Separating "what failed" (`suggestion`) from "how to fix" (`fix`)
2. Using `WARN` vs `FAIL` headers based on severity
3. Always showing `Advisory:` line to make blocking status explicit

## Implementation Phases

### Phase 1: Schema Changes

#### Step 1.1: Add `fix` field to Check struct
- File: `internal/config/schema.go`
- Add `Fix string` field to `Check` struct
- Add YAML tag: `yaml:"fix,omitempty"`

#### Step 1.2: Add `Fix` field to Violation struct
- File: `internal/orchestrator/orchestrator.go`
- Add `Fix string` field to `Violation` struct

#### Step 1.3: Update ALL violation creation sites
- File: `internal/orchestrator/orchestrator.go`
- **Location 1**: `Run()` function (lines 254-273) - main check execution
- **Location 2**: `Run()` function (lines 163-169) - dependency skip violations
- **Location 3**: `RunCheck()` function (lines 430-443) - single check execution
- All three sites must populate the new `Fix` field from check config

#### Step 1.4: Update interpolation for fix field
- File: `internal/config/interpolate.go`
- Ensure `InterpolateWithExtracted()` is called for both `suggestion` and `fix` fields
- Both fields should support `{{.extracted}}` template values

### Phase 2: JSON Output Updates

#### Step 2.1: Update JSONViolation struct
- File: `internal/output/json.go`
- Add `Fix string` field to `JSONViolation` struct with tag `json:"fix,omitempty"`

#### Step 2.2: Update FormatJSON
- Ensure `Fix` field is populated in JSON output

### Phase 3: Text Formatter Updates

#### Step 3.1: Update `formatViolation()` for quiet mode
- File: `internal/output/formatter.go` (lines 75-88)
- Use `WARN` header when severity is "warning", `FAIL` otherwise
- Output logic:
  - If suggestion: show it (no prefix)
  - If fix: show `Fix: {fix}`
  - If neither suggestion nor fix: show `Fix: {run command}` as fallback
  - Always show `Advisory:` line based on severity

#### Step 3.2: Update `formatVerbose()` for verbose mode
- File: `internal/output/formatter.go` (lines 48-73)
- **Architecture note**: Verbose mode reads from `Check` struct, not `Violation` struct
- Options:
  - a) Refactor to use `Violation` struct (more work, better consistency)
  - b) Read `Fix` from `Check` directly (simpler, may miss interpolated values)
- Recommended: Option (a) - refactor to use Violation for consistency
- Update header: use `✗ %-15s WARN` vs `✗ %-15s FAIL` based on severity
- Show `Advisory:` line for failed checks

#### Step 3.3: Handle timeout format
- Timeout violations use hardcoded suggestion: "Check timed out..."
- For timeouts: show the timeout message as suggestion, use check's `fix` if present
- If no `fix` on check, show `Fix: {run command}` as fallback

### Phase 4: Testing

#### Step 4.1: Update existing unit tests
- Update output format expectations for new WARN/FAIL headers
- Update JSON output tests in `json_test.go`
- Update formatter tests for Advisory line

#### Step 4.2: Add new unit tests (table-driven)
| suggestion | fix | severity | expected output |
|------------|-----|----------|-----------------|
| set | empty | error | suggestion + fallback fix + Advisory: blocks |
| empty | set | error | Fix: {fix} + Advisory: blocks |
| set | set | error | suggestion + Fix: {fix} + Advisory: blocks |
| set | set | warning | suggestion + Fix: {fix} + Advisory: does not block |
| empty | empty | error | Fix: {run} + Advisory: blocks |

#### Step 4.3: Integration tests
- Test end-to-end with actual check execution
- Verify JSON and text output consistency

### Phase 5: Config Updates

#### Step 5.1: Update vibeguard.yaml
- File: `vibeguard.yaml`
- Add `fix` field to checks where appropriate
- Use `{{.extracted}}` values in `suggestion` for assertion checks
- Example for test-coverage:
  ```yaml
  suggestion: "Coverage is {{.coverage}}%, need 80%"
  fix: "Add tests to improve coverage"
  ```

## Output Format Reference

### Error (blocking)
```text
FAIL  lint (error)

  Fix: golangci-lint run ./...
  Advisory: blocks commit
```

### Warning (non-blocking)
```text
WARN  test (warning)

  Fix: go test ./...
  Advisory: does not block commit
```

### Assertion failure
```text
FAIL  test-coverage (error)

  Coverage is 72%, need 80%
  Fix: Add tests to improve coverage
  Advisory: blocks commit
```

### Timeout
```text
FAIL  test (timeout after 10s)

  Check timed out. Consider increasing timeout or check for hanging tests.
  Fix: go test ./...
  Advisory: blocks commit
```

### Empty suggestion and fix (fallback)
```text
FAIL  lint (error)

  Fix: golangci-lint run ./...
  Advisory: blocks commit
```

When both `suggestion` and `fix` are empty, fall back to showing the `run` command as the fix.

## Acceptance Criteria

- [ ] `fix` field added to `config.Check` struct
- [ ] `Fix` field added to `Violation` struct
- [ ] All three violation creation sites populate `Fix`
- [ ] `InterpolateWithExtracted()` called for both `suggestion` and `fix`
- [ ] `JSONViolation` struct updated with `Fix` field
- [ ] `formatViolation()` uses `WARN` vs `FAIL` based on severity
- [ ] `formatViolation()` outputs suggestion + fix appropriately
- [ ] `formatViolation()` falls back to run command when suggestion and fix are empty
- [ ] `formatViolation()` always shows `Advisory:` line
- [ ] `formatVerbose()` refactored to use Violation struct
- [ ] `formatVerbose()` uses WARN/FAIL and shows Advisory line
- [ ] Timeout violations handled correctly
- [ ] Unit tests updated with table-driven approach
- [ ] Integration tests added
- [ ] `vibeguard.yaml` uses new schema

## Next Steps

Start with Phase 1 (schema changes) as it establishes the foundation for all other changes.
