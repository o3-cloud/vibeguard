---
summary: Fixed missing file/line context for grok and assertion errors in RunCheck() method
event_type: code
sources:
  - internal/orchestrator/orchestrator.go
  - internal/config/config.go
tags:
  - error-handling
  - grok
  - assertions
  - debugging
  - vibeguard-trb
  - code-quality
---

# Fix: Add File/Line Context to Grok and Assert Errors

## Completed Work

Fixed issue **vibeguard-trb**: "Grok and assert errors lack file/line context"

### Problem
The `RunCheck()` method in the orchestrator was not wrapping grok and assertion errors with file/line context information from the YAML configuration. While the `Run()` method properly wrapped these errors with `ExecutionError` containing line numbers and error types, `RunCheck()` was returning raw errors without this context.

### Solution
Updated `RunCheck()` method to:
1. Track the check index when finding a check by ID (line 343-350)
2. Wrap grok compilation errors with `ExecutionError` containing file/line context (lines 376-386)
3. Wrap grok matching errors with `ExecutionError` containing file/line context (lines 388-399)
4. Wrap assertion evaluation errors with `ExecutionError` containing file/line context (lines 407-416)

### Implementation Details
- Used existing `FindCheckNodeLine()` method to obtain accurate YAML line numbers
- Reused the `ExecutionError` type with `ErrorType` field for distinction
- Maintained consistency with error handling patterns in the `Run()` method
- Added comments to document the error wrapping for clarity

### Testing Results
- All 276 existing tests pass without regression
- No new test failures introduced
- Implementation verified through full test suite execution

### Files Modified
- `internal/orchestrator/orchestrator.go` (RunCheck method, lines 338-419)

### Technical Notes
The infrastructure for file/line context was already well-established through:
- `ExecutionError` type with LineNum and ErrorType fields
- `FindCheckNodeLine()` method for YAML node lookup
- Consistent error formatting through `Error()` method

The fix simply ensures this infrastructure is used consistently across all error paths in RunCheck(), matching the behavior of the Run() method.
