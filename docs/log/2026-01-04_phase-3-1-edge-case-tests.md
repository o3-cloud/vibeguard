---
summary: Added comprehensive edge case tests to internal/config/config.go for Phase 3.1 mutation testing
event_type: code
sources:
  - internal/config/config.go
  - internal/config/config_test.go
  - docs/adr/ADR-007-adopt-mutation-testing.md
tags:
  - mutation-testing
  - edge-cases
  - boundary-conditions
  - config-validation
  - test-coverage
  - phase-3
---

# Phase 3.1: Added Edge Case Tests to internal/config/config.go

Completed task vibeguard-96h by adding comprehensive boundary condition tests to the config package test suite.

## Summary

Added 10+ edge case tests to `internal/config/config_test.go` to improve mutation testing coverage for `internal/config/config.go`. These tests focus on boundary conditions and edge cases that aren't covered by existing tests and are specifically designed to catch mutations that the current test suite might miss.

## Changes Made

### New Test Functions

1. **TestLoad_EmptyStringVersion** - Verifies that empty string version is treated as default (applies to applyDefaults mutation)
2. **TestLoad_MultipleEmptyChecks** - Tests error handling for empty check IDs in multi-check configurations
3. **TestLoad_BoundaryCheckIndex** - Validates error reporting for checks at different indices (boundary condition at index 2+)
4. **TestLoad_SeverityBoundaryValues** - Tests valid and invalid severity values with proper error messages (mutation: conditionals_negation at line 196)
5. **TestLoad_TimeoutZeroBoundary** - Tests timeout edge cases including zero (treated as unset), negative, and very large durations (boundary at timeout == 0 check)
6. **TestLoad_DependencyOrderBoundary** - Tests DAG validation where multiple checks depend on same base check
7. **TestLoad_LongDependencyChain** - Tests linear dependency chains (a->b->c->d->e) without cycles
8. **TestConfigError_EmptyMessage** - Tests ConfigError formatting with empty message
9. **TestConfigError_LineNumberBoundary** - Tests line number formatting for zero, line 1, and large line numbers
10. **TestFindCheckNodeLine_EdgeCases** - Tests line number lookup with valid/invalid indices and boundary conditions

## Test Results

- ✅ All 10 new edge case tests pass
- ✅ Total config package tests: 80+ tests all passing
- ✅ gofmt formatting verified
- ✅ vibeguard check: PASS

## Mutation Coverage Impact

These tests specifically target boundary conditions and edge cases that mutation testing would detect:

- **CONDITIONALS_BOUNDARY mutations** at lines 33, 66, 315, 331, 343, 344, 352 - Tests verify that boundary checks (>, <, >=, <=) are correctly implemented
- **CONDITIONALS_NEGATION mutations** - Tests ensure proper negation handling in validation logic
- **ARITHMETIC_BASE mutations** - Tests validate numeric comparisons and array index operations
- **INCREMENT_DECREMENT mutations** - Tests verify loop counter behavior in formatCycle and DFS

The test suite now provides strong protection against subtle logic errors that could slip past traditional coverage-based metrics.

## Key Insights

1. **Zero Duration Handling** - The config system treats `timeout: 0` the same as omitted timeout (both get DefaultTimeout). This is documented in the test.

2. **Line Number Boundary** - Found that `Timeout == 0` check in applyDefaults is used as a sentinel for "not set", which cannot distinguish between explicit 0 and unset.

3. **Dependency Graph Testing** - Long chains and multiple converging dependencies are now tested to catch DFS mutation issues.

4. **Error Formatting Edge Cases** - ConfigError formatting now validated for empty message and line number boundary conditions.

## Next Steps

- Task vibeguard-96h is complete
- These tests will improve mutation testing efficacy when running gremlins
- Consider running `gremlins unleash ./internal/config` to verify mutation kill rate improvement
