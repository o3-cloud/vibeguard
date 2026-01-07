---
summary: Completed Phase 3 - Implemented tag filtering in orchestrator
event_type: code
sources:
  - docs/specs/SPEC-tags.md
  - internal/orchestrator/orchestrator.go
  - internal/cli/list.go
tags:
  - tags-feature
  - orchestrator-filtering
  - phase-3
  - tag-filtering-logic
  - dependency-handling
---

# Phase 3: Orchestrator Tag Filtering Implementation Complete

## Summary

Successfully implemented Phase 3 of the tags feature by adding tag-based filtering logic to the orchestrator. This enables runtime filtering of checks based on tags while respecting dependencies.

## Implementation Details

### Files Modified
1. **internal/orchestrator/orchestrator.go**
   - Added `filterChecksByTags()` method to apply include/exclude filters
   - Modified `Run()` method to apply filters before dependency graph building
   - Enhanced dependency validation to check for tag-excluded dependencies
   - Added clear warning messages for skipped checks

2. **internal/cli/list.go**
   - Fixed output to use `cmd.OutOrStdout()` instead of `fmt.Printf`
   - Ensures list command output is properly captured in tests
   - Tags now display correctly in verbose mode

### Filtering Logic

```go
// Include filter: run checks matching ANY tag (OR logic)
vibeguard check --tags format,lint    // runs fmt OR lint

// Exclude filter: skip checks matching ANY tag
vibeguard check --exclude-tags slow   // runs all except slow checks

// Combined: apply include first, then exclude
vibeguard check --tags ci --exclude-tags slow  // ci checks except slow
```

### Dependency Handling

When a check depends on an excluded check:
1. The dependent check is skipped (not executed)
2. Clear warning message: "Skipping check 'test': required dependency 'fmt' not in filtered set"
3. Existing dependency failure logic preserved

Example:
```yaml
checks:
  - id: fmt
    tags: [format]
  - id: test
    tags: [test]
    requires: [fmt]
```

Running `vibeguard check --tags test` will skip `test` with warning that `fmt` is not in the filtered set.

## Testing

All tests passing:
- CLI tests: 28 tests pass (including new tag flag tests)
- Orchestrator tests: 60+ tests pass (including parallel execution and dependency tests)
- New test added: `TestRunList_WithTags` validates tags display in list output
- Policy checks: ✅ All vibeguard checks pass

## Key Design Decisions

1. **Strict filtering**: Dependencies are NOT automatically included when filtering by tags
   - Rationale: Explicit is better than implicit
   - Users should tag dependencies appropriately for their use case

2. **OR logic for both include and exclude**:
   - Include: matches ANY specified tag
   - Exclude: skips if matches ANY specified tag
   - Sufficient for most use cases (AND logic not needed)

3. **Silent tag mismatch**: Non-matching tags don't error
   - Allows flexible filtering with typos or unused tags
   - Matches behavior of other tools (golangci-lint, pytest)

4. **Filter applied early**: Filters applied before dependency graph building
   - Prevents graph building errors with excluded checks
   - More efficient than filtering after graph construction

## Status

✅ Phase 3 complete
✅ Ready for Phase 4 (JSON output updates)

## Next Steps

Phase 4 will update JSON output to include tags in results and add --tags flag to list command for filtering display.
