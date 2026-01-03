---
summary: Improved mutation test efficacy in internal/assert from 88.31% to 98.70% by adding targeted tests
event_type: code
sources:
  - docs/adr/ADR-007-adopt-mutation-testing.md
  - internal/assert/eval_test.go
  - internal/assert/parser.go
tags:
  - mutation-testing
  - test-improvement
  - internal-assert
  - gremlins
  - code-quality
---

# Mutation Testing Improvements for internal/assert

## Summary

Investigated and addressed surviving mutants in `internal/assert` package as part of task `vibeguard-q3i`. Improved test efficacy from **88.31% to 98.70%** by adding 3 new targeted test functions.

## Initial State

9 mutants survived with 88.31% efficacy:
- `eval.go:215:72` - CONDITIONALS_BOUNDARY in `<` comparison
- `eval.go:276:37` - INVERT_NEGATIVES and ARITHMETIC_BASE in formatFloat precision
- `parser.go:37:*` - Multiple mutations in formatError loop conditions

## Analysis

### eval.go:215:72 - Less-Than Boundary
The mutation changed `cmp < 0` to `cmp <= 0` in the less-than comparison operator. Tests didn't verify that equal values return false for `<`.

### eval.go:276:37 - Float Formatting Precision
The mutation changed `-1` (minimum precision) to `1` (one decimal place) in `strconv.FormatFloat`. Tests didn't verify that multi-decimal precision was preserved.

### parser.go:37 - formatError Loop Conditions
Multiple mutations in `for i := 0; i < pos-1 && i < len(p.input); i++`. Tests didn't verify exact pointer positioning.

## Tests Added

### TestEvaluator_LessThanBoundary
Kills boundary mutation by verifying `10 < 10` returns false (equal values).

### TestFormatFloat_Precision
Kills precision mutations by verifying multi-decimal operations:
- `1.5 + 1.25 == 2.75`
- `7 / 4 == 1.75`
- `-3.14159 == -3.14159`

### TestParser_FormatErrorPointerPosition
Kills loop condition mutations by verifying exact `^` pointer positions in error messages at positions 1, 3, 5, and 10.

## Results

- **Before:** 68 killed, 9 lived, 88.31% efficacy
- **After:** 76 killed, 1 lived, 98.70% efficacy
- **Mutants killed:** 8 additional

## Remaining Mutant

One mutant survives: `CONDITIONALS_BOUNDARY at parser.go:37:29` which changes `i < len(p.input)` to `i <= len(p.input)`.

This is an **equivalent mutant** - the defensive bounds check protects against a hypothetical bug where error position exceeds input length. In practice, the lexer always returns positions within bounds, so this code path is never triggered. The mutation doesn't change observable behavior for any reachable input.

## Files Modified

- `internal/assert/eval_test.go` - Added 3 new test functions (~80 lines)

## Related

- ADR-007: Adopt Gremlins for Mutation Testing
- Task: vibeguard-q3i
