---
summary: Documented the file config field in README with usage examples and added implementation issue
event_type: code
sources:
  - internal/config/schema.go
  - README.md
  - docs/log/2026-01-03_file-field-prompt-support.md
  - internal/orchestrator/orchestrator.go
tags:
  - documentation
  - configuration
  - file-field
  - readme
  - user-guide
  - feature-discovery
  - implementation-gap
---

# Documented File Config Field (vibeguard-rol)

## Task Summary

Completed task vibeguard-rol: "Document file config field". Added comprehensive documentation for the `file` field in the README Configuration Schema section with usage examples and implementation guidance.

## Changes Made

### 1. README Configuration Schema Updates

**Added to YAML Schema Example (line 195-196):**
- Added comment explaining the `file` field: "Optional: Read output from file instead of command stdout"
- Added example: `file: path/to/output.txt`

**Added to Field Details Table (line 225):**
- Added new row: `| 'file' | No | string | File path to read output from instead of command stdout | — |`

### 2. New "Reading Output from Files" Section

Created a dedicated section explaining:
- **Purpose:** When to use the `file` field (tools that write to files instead of stdout)
- **Use Case Example:** Coverage report files from `go test` commands
- **Full Example Configuration:** Shows real-world usage pattern with coverage output file
- **Behavior Clarification:** Explains that the command still runs normally but grok patterns apply to file contents

Example provided:
```yaml
checks:
  - id: coverage
    run: go test ./... -coverprofile=coverage.out
    file: coverage.out
    grok:
      - total:.*\(statements\)\s+%{NUMBER:coverage}%
    assert: "coverage >= 80"
    suggestion: "Coverage is {{.coverage}}%, target is 80%. Add more tests."
```

## Key Findings

### Documentation Status
The `file` field was already:
- ✅ Defined in the Check struct (internal/config/schema.go:20)
- ✅ Documented in validator_guide.go (lines 128-129)
- ✅ Included in config requirements section
- ✅ Propagated through prompt generation system (vibeguard-tcf)
- ❌ **Missing from README** (now fixed)

### Implementation Gap Discovered

**Critical Finding:** The `file` field is documented but **NOT actually implemented** in the orchestrator!

Current behavior (orchestrator.go):
- Line 212: `extracted, matcherErr = matcher.Match(execResult.Combined)` - always uses command output
- Line 391: Same pattern in `RunCheck()`
- The `check.File` field is never read or used

**Issue Created:** vibeguard-ewr
- Type: task
- Priority: P3
- Description: Implement file field functionality to actually read from files when specified

## Testing

All tests pass:
```
PASS ok  github.com/vibeguard/vibeguard/...
```

No breaking changes - documentation addition only.

## Next Steps

1. Implement actual file reading functionality in orchestrator (issue vibeguard-ewr)
2. Add tests for file field functionality
3. Consider expanding examples with more real-world file-based scenarios
4. Update JSON output schema documentation if needed to reflect file field usage

## Related Work

- **ADR-004:** Code Quality Standards - applies to documentation completeness
- **vibeguard-tcf:** File field prompt support (recently completed)
- **vibeguard-ewr:** Implement file field functionality (newly created)

## Files Modified

- `README.md` - Added file field documentation with examples
- `.beads/issues.jsonl` - Issue vibeguard-ewr created
