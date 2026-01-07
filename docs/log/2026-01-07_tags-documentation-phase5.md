---
summary: Completed Phase 5 documentation for vibeguard tags feature
event_type: code
sources:
  - docs/specs/SPEC-tags.md
  - README.md
  - examples/basic.yaml
  - examples/advanced.yaml
  - docs/ai-assisted-setup.md
tags:
  - tags-feature
  - documentation
  - phase-5
  - examples
  - user-guide
---

# Tags Feature Documentation - Phase 5 Complete

## Overview

Completed comprehensive documentation for the vibeguard tags feature as specified in SPEC-tags.md Phase 5. This work ensures users understand tag-based check filtering and provides practical examples for common use cases.

## Changes Made

### 1. README.md Updates

**Global Flags Table:** Added `--tags` and `--exclude-tags` flags documentation

**vibeguard check command:**
- Added tag filtering section with CLI examples
- Documented OR logic for `--tags` flag
- Documented exclusion with `--exclude-tags` flag
- Showed combination of inclusion and exclusion

**New vibeguard tags command:**
- Documents new discovery command
- Shows syntax for listing all available tags

**Check Tags section (comprehensive):**
- Explains YAML schema for tags field
- Lists standard tag conventions (format, lint, test, security, build, fast, slow, pre-commit, ci, llm)
- Shows filtering examples with AND/OR logic table
- Documents dependency handling with tags
- Explains that dependencies are NOT automatically included when filtering
- Includes real example showing skipped checks with warnings

### 2. New examples/basic.yaml

Created minimal but practical example with:
- 5 checks (fmt, vet, lint, test, coverage)
- All checks tagged appropriately
- Clear comments explaining tags
- Filter examples in end comments

Tags used:
- `format`, `lint`, `fast`, `pre-commit`, `test`, `slow`, `ci`

### 3. examples/advanced.yaml Refactored

Converted from comment-based phases to tag-based organization:

**Phase 1 (Fast Deterministic Checks):**
- fmt: `[format, fast, pre-commit]`
- vet: `[lint, fast, pre-commit]`
- lint: `[lint, fast, ci]`

**Phase 2 (Testing and Metrics):**
- test: `[test, slow, ci]`
- coverage: `[test, slow, ci]`
- complexity: `[test, slow, ci]`
- build: `[build, slow, ci]`

**Phase 3 (LLM-Powered Checks):**
- llm-architecture-review: `[llm, slow, ci-only]`
- llm-security-review: `[llm, security, slow, ci-only]`
- llm-pr-quality: `[llm, slow, ci-only]`

Added comprehensive tag-based filtering examples at end of file showing:
- Local development (fast checks)
- Pre-commit hooks
- CI/CD pipelines
- LLM reviews only
- Full validation
- Tag discovery

### 4. docs/ai-assisted-setup.md Updates

Updated all three language examples (Go, Node.js, Python):

**Before:** No tags, just check definitions

**After:** 
- All checks include appropriate tags
- Added "Usage with Tags" section to each example
- Shows practical filtering scenarios (pre-commit, local fast checks, CI pipeline)
- Helps new users understand tag adoption immediately

Example patterns:
- Format checks: `[format, fast, pre-commit]`
- Lint checks: `[lint, fast, ci]`
- Test checks: `[test, slow, ci]`

## Quality Assurance

**vibeguard check verification:**
```
✓ vet             passed
✓ fmt             passed
✓ actionlint      passed
✓ lint            passed
✓ staticcheck     passed
✓ test            passed
✓ test-coverage   passed
✓ gosec           passed
✓ docker          passed
✓ build           passed
```

All 10 checks pass without errors.

## Spec Alignment

✅ **All Phase 5 requirements met:**

1. **README.md updates** - Comprehensive tags section with examples
2. **examples/basic.yaml** - New practical starting point with tags
3. **examples/advanced.yaml** - Converted to tag-based phase organization
4. **Standard tag conventions** - Documented in README with descriptions
5. **Troubleshooting** - Dependency handling explained
6. **ai-assisted-setup.md** - Updated examples with tag recommendations
7. **Backwards compatibility** - All existing configs still work

## Key Design Decisions

**Tags are strictly optional:** Existing configs without tags continue to work unchanged

**Standard conventions, not enforced:** Users can use any tag names following the pattern

**Dependency handling is strict:** Filtered checks don't automatically include their dependencies. If a filtered check depends on an excluded check, it's skipped with a warning

**OR logic for filtering:** `--tags format,lint` matches checks with EITHER tag

**Exclusion also uses OR logic:** `--exclude-tags slow,llm` excludes checks with EITHER tag

## User Impact

Users can now:
- Run only fast checks locally: `vibeguard check --tags fast`
- Run pre-commit hooks: `vibeguard check --tags pre-commit`
- Run CI checks: `vibeguard check --tags ci`
- Skip slow LLM checks: `vibeguard check --exclude-tags llm,slow`
- Discover available tags: `vibeguard tags`

Documentation provides clear examples and best practices for tag adoption.

## Next Steps

Phase 5 is complete. Implementation phases 1-4 (schema, CLI flags, orchestrator filtering, output formatting) should be verified to ensure tag feature end-to-end functionality.
