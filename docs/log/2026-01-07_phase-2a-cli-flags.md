---
summary: Completed Phase 2a - Added --tags and --exclude-tags CLI flags to check command
event_type: code
sources:
  - docs/specs/SPEC-tags.md
  - internal/cli/check.go
  - internal/orchestrator/orchestrator.go
tags:
  - tags-feature
  - cli-flags
  - phase-2a
  - tag-filtering
  - implementation
---

# Phase 2a: CLI Flags Implementation Complete

## Summary

Successfully implemented Phase 2a of the tags feature by adding `--tags` and `--exclude-tags` CLI flags to the vibeguard check command. This enables users to filter checks by tags at runtime.

## Implementation Details

### Files Modified
1. **internal/cli/check.go**
   - Added `tags` and `excludeTags` variables at package level
   - Registered StringSlice flags in init() function
   - Updated examples in command help text
   - Modified runCheck() to set TagFilter on orchestrator when flags are specified

2. **internal/orchestrator/orchestrator.go**
   - Verified TagFilter struct and tagFilter field exist
   - Verified SetTagFilter() method exists
   - No changes needed (already implemented by previous phase 1 work)

### CLI Interface

```bash
# Run checks matching ANY of the specified tags (OR logic)
vibeguard check --tags format,lint

# Exclude checks matching ANY of the specified tags
vibeguard check --exclude-tags slow,llm

# Combine inclusion and exclusion
vibeguard check --tags ci --exclude-tags slow
```

### Help Output
Both flags now appear in `vibeguard check --help` with clear descriptions and examples.

## Testing

- All existing tests pass (24 CLI tests)
- New tests for flag parsing already exist and pass:
  - `TestRunCheck_WithTagsFlag`
  - `TestRunCheck_WithExcludeTagsFlag`
- Build succeeds without errors
- Flag help text displays correctly

## Key Decisions

1. **StringSlice flags** - Used spf13/cobra's StringSliceVar for comma-separated tag lists
2. **Filter storage** - Reused existing TagFilter struct and SetTagFilter method from orchestrator
3. **Flag combination** - Both `--tags` and `--exclude-tags` can be used together (inclusion first, then exclusion)

## Status

✅ Phase 2a complete
✅ Ready for Phase 3 (orchestrator filtering implementation)

## Next Steps

Phase 3 will implement the actual filtering logic in orchestrator.Run() to apply these filters to the check selection.
