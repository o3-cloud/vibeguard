---
summary: Completed comprehensive tests for --list-templates flag functionality in init command
event_type: code
sources:
  - internal/cli/init_test.go
  - internal/cli/init.go
tags:
  - testing
  - list-templates
  - init-command
  - flag-functionality
  - test-coverage
  - stdout-capture
---

# Completed vibeguard-937: Update init_test.go tests for --list-templates functionality

## Summary

Enhanced test coverage for the `vibeguard init --list-templates` flag by adding comprehensive unit tests that validate the functionality, output format, and flag handling.

## Changes Made

### 1. Enhanced `TestListTemplates()`
- Modified to capture and verify stdout output
- Validates that output contains expected headers ("Available templates:", "Usage:")
- Previously was just a basic smoke test

### 2. Added `TestRunInit_ListTemplatesFlag()`
- Tests the `--list-templates` flag specifically with `initListTemplates = true`
- Verifies that runInit properly handles the flag and calls listTemplates()
- Complements existing `TestRunInit_ListTemplates()` which tests `--template list`

### 3. Added `TestListTemplates_OutputFormat()`
- Comprehensive output format validation
- Verifies:
  - Output structure has minimum 3 lines
  - First line contains "Available templates:" header
  - At least one template entry exists (checks for "go-" prefix)
  - Usage information is included
  - Proper line-based output formatting

## Technical Details

### Stdout Capture Pattern
Used os.Pipe() to capture stdout during test execution, allowing validation of printed output format and content without affecting test isolation.

### Error Handling
Fixed linter errors by properly handling `w.Close()` return values with blank identifier assignments (`_ = w.Close()`).

## Testing Results

All new tests pass successfully:
- `TestListTemplates` ✓
- `TestRunInit_ListTemplates` ✓
- `TestRunInit_ListTemplatesFlag` ✓
- `TestListTemplates_OutputFormat` ✓

All existing init tests continue to pass with no regressions.

## Pre-existing Issue Found

During vibeguard check execution, identified a pre-existing test failure:
- **Issue**: `TestPromptTokenEstimate` in `internal/cli/assist` failing
- **Reason**: Prompt exceeds 4000 token estimate (~4204 tokens)
- **Created**: vibeguard-968 to track this pre-existing issue
- **Impact**: Does not affect --list-templates functionality

## Compliance

✓ All new tests follow project conventions
✓ Linter passes with no errors
✓ No regressions in existing tests
✓ Proper error handling for closed file descriptors
✓ Follows ADR-004 code quality standards for testing
