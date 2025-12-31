---
summary: Added test coverage for cancelled status in JSON output
event_type: code
sources:
  - internal/output/json_test.go
  - internal/output/json.go
  - internal/executor/executor.go
tags:
  - testing
  - json-output
  - cancelled-status
  - vibeguard-fp7
  - coverage
---

# Completed vibeguard-fp7: Added Test Coverage for Cancelled Status in JSON Output

## Summary

Added `TestFormatJSON_CancelledStatus` test to `internal/output/json_test.go` to verify that cancelled checks produce correct JSON output with status='cancelled'.

## Changes

- Added new test function `TestFormatJSON_CancelledStatus` in json_test.go (lines 155-197)
- Test verifies:
  - Cancelled checks produce JSON with status='cancelled'
  - Non-cancelled checks still produce correct status
  - Exit code is properly set
  - Duration is captured correctly

## Implementation Details

The test creates a realistic scenario with:
1. One passing check (fmt) - Duration: 100ms
2. One cancelled check (vet) - Duration: 250ms with `Cancelled: true` in executor.Result
3. Verifies the JSON output correctly reflects the cancelled status with exit code `ExitCodeTimeout`

### How Cancelled Status Works

The cancelled status is set when a check is cancelled (e.g., by fail-fast mechanism). The JSON output handler in `json.go` (lines 44-49) prioritizes cancellation checking:

```
if r.Execution.Cancelled {
    return "cancelled"
} else if r.Passed {
    return "passed"
} else {
    return "failed"
}
```

## Test Results

✅ New test passes: PASS
✅ All existing output tests pass: 11 tests
✅ Full test suite passes: All packages

### Test Output

```
=== RUN   TestFormatJSON_CancelledStatus
--- PASS: TestFormatJSON_CancelledStatus (0.00s)
PASS
ok  	github.com/vibeguard/vibeguard/internal/output	0.167s
```

## Coverage Impact

This test fills a gap in the JSON output test suite by covering the cancelled status path, which is used when checks are cancelled due to fail-fast mode or context cancellation.

## Task Status

✅ **Completed** - vibeguard-fp7
- Task: Add test coverage for cancelled status in JSON output
- Status: DONE
- All tests passing, no issues found
