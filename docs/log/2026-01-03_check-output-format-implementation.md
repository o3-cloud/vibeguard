---
summary: Implemented check output format improvements for AI agents with Fix field, WARN/FAIL headers, and Advisory line
event_type: code
sources:
  - docs/log/2026-01-02_check-output-format-implementation-plan.md
  - internal/orchestrator/orchestrator.go
  - internal/output/formatter.go
  - internal/output/json.go
  - internal/config/interpolate.go
tags:
  - ai-agent
  - cli
  - output-formatting
  - vibeguard-check
  - schema-change
  - epic-1zn
---

# Check Output Format Implementation

Completed implementation of the check output format improvements for AI agents (epic vibeguard-1zn). This work makes vibeguard check output more actionable for AI agents by separating "what failed" from "how to fix".

## Changes Made

### Phase 1: Schema Changes

1. **Added Fix field to Violation struct** (vibeguard-1zn.2)
   - File: `internal/orchestrator/orchestrator.go:42`
   - Added `Fix string` field to carry the fix instruction from check config

2. **Updated all violation creation sites** (vibeguard-1zn.3)
   - Three locations in `orchestrator.go` now populate the Fix field:
     - Dependency skip violations (line 169)
     - Main check execution violations (line 266)
     - Single check execution violations (line 443)

3. **Updated interpolation for fix field** (vibeguard-1zn.4)
   - File: `internal/config/interpolate.go:15`
   - Added `Fix` field interpolation with config vars

### Phase 2: JSON Output Updates

4. **Updated JSONViolation struct** (vibeguard-1zn.5)
   - File: `internal/output/json.go:31`
   - Added `Fix string` with `json:"fix,omitempty"` tag

5. **Updated FormatJSON** (vibeguard-1zn.6)
   - File: `internal/output/json.go:64`
   - Populates Fix field in JSON output

### Phase 3: Text Formatter Updates

6. **Updated formatViolation() for quiet mode** (vibeguard-1zn.7)
   - Uses WARN header for warning severity, FAIL for error/timeout
   - Shows suggestion (if present), then Fix line
   - Falls back to run command when no suggestion and no fix
   - Always shows Advisory line ("blocks commit" or "does not block commit")

7. **Refactored formatVerbose()** (vibeguard-1zn.8)
   - Now uses Violation struct for consistency
   - Same WARN/FAIL header logic as quiet mode
   - Shows Fix and Advisory lines for failed checks

8. **Handle timeout violations** (vibeguard-1zn.9)
   - Timeout shows "(timeout)" in status info
   - Timeout message shown as suggestion
   - Fix field from check config still applies

### Phase 4: Testing

9. **Updated existing unit tests** (vibeguard-1zn.10)
   - Updated test expectations for new output format
   - Removed checks for old "Tip:" prefix and command line

10. **Added table-driven tests** (vibeguard-1zn.11)
    - `TestFormatViolation_OutputCombinations` covers all combinations:
      - suggestion set/empty, fix set/empty, severity error/warning
      - timeout violation handling
      - fallback to run command

11. **Added JSON test for Fix field** (vibeguard-1zn.12)
    - `TestFormatJSON_WithFixField` verifies Fix field in JSON output

### Phase 5: Config Updates

12. **Updated vibeguard.yaml** (vibeguard-1zn.13)
    - Added fix field to all checks
    - Updated test-coverage to use `{{.coverage}}` in suggestion

## New Output Format

### Error (blocking)
```
FAIL  lint (error)

  golangci-lint found lint issues
  Fix: golangci-lint run ./...
  Advisory: blocks commit
```

### Warning (non-blocking)
```
WARN  test (warning)

  Some tests are failing
  Fix: go test ./...
  Advisory: does not block commit
```

### Timeout
```
FAIL  test (timeout)

  Check timed out. Consider increasing the timeout value or optimizing the command.
  Advisory: blocks commit
```

## Files Modified

- `internal/orchestrator/orchestrator.go` - Violation struct and creation sites
- `internal/output/formatter.go` - formatViolation() and formatVerbose()
- `internal/output/json.go` - JSONViolation struct and FormatJSON
- `internal/config/interpolate.go` - Fix field interpolation
- `internal/output/formatter_test.go` - Updated and new tests
- `internal/output/json_test.go` - New JSON test
- `vibeguard.yaml` - Added fix fields to all checks

## Next Steps

- Close beads tasks vibeguard-1zn.2 through vibeguard-1zn.13
- Consider closing the parent epic vibeguard-1zn if all tasks complete
