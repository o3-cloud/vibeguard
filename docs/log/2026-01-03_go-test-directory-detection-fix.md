---
summary: Implemented Go test directory detection for co-located tests and fixed root directory skipping bug
event_type: code
sources:
  - internal/cli/inspector/metadata.go
  - internal/cli/inspector/metadata_test.go
  - https://go.dev/doc/tutorial/add-a-test
tags:
  - go
  - testing
  - inspector
  - bug-fix
  - ai-assisted-setup
---

# Go Test Directory Detection Fix

## Problem

The `ProjectStructure.TestDirs` field was always empty for Go projects when running `vibeguard init --assist`. This was tracked in task `vibeguard-51e`.

Go conventionally uses co-located tests (`*_test.go` files in the same directory as source files) rather than dedicated test directories like `tests/` or `test/`. The existing implementation only checked for dedicated test directories.

## Root Cause

Two issues were identified:

1. **Missing co-located test detection**: The `extractGoStructure` function only looked for dedicated directories (`test`, `tests`, `testdata`), not for directories containing `*_test.go` files.

2. **Root directory being skipped**: When using `.` as the project root, `filepath.Walk` starts with path `.` which has `info.Name()` equal to `.`. The condition `strings.HasPrefix(name, ".")` was true, causing the entire root directory to be skipped via `filepath.SkipDir`.

## Solution

### 1. Added `findGoTestDirs()` function

New function that walks the project tree and identifies directories containing `*_test.go` files:

```go
func (m *MetadataExtractor) findGoTestDirs() []string
```

Features:
- Walks from project root to find all `*_test.go` files
- Skips hidden directories (`.git`, etc.), `vendor`, `node_modules`, `testdata`
- Only includes directories that are standard Go source directories (`cmd`, `pkg`, `internal`, `lib`, `api`, `app`) or root
- Returns unique, sorted list of directory paths

### 2. Added `isGoSourceDir()` helper

Validates that a directory path is a standard Go source directory:

```go
func (m *MetadataExtractor) isGoSourceDir(relPath string) bool
```

### 3. Fixed root directory skip bug

Changed the directory skip condition to not skip the root directory itself:

```go
// Before (buggy)
if strings.HasPrefix(name, ".") || name == "vendor" ...

// After (fixed)
if path != m.root && (strings.HasPrefix(name, ".") || name == "vendor" ...)
```

## Testing

Added comprehensive tests:

- `TestMetadataExtractor_ExtractGoStructure_ColocatedTests` - verifies co-located test detection
- `TestMetadataExtractor_ExtractGoStructure_RootLevelTests` - verifies root-level test detection
- `TestMetadataExtractor_ExtractGoStructure_ExcludesVendorAndTestdata` - verifies exclusions
- `TestMetadataExtractor_findGoTestDirs` - unit tests for the new function
- `TestMetadataExtractor_isGoSourceDir` - table-driven tests for source dir validation

## Verification

Running `vibeguard init --assist` now shows test directories:

```
### Project Structure:
- Source directories: internal
- Test directories: internal/assert, internal/cli, internal/cli/assist, ...
- Entry points: cmd/vibeguard/main.go
```

## Related

- Task: `vibeguard-51e` - Improve Go test directory detection
- ADR-005: Dogfooding VibeGuard on itself
