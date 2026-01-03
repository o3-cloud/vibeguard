---
summary: Analysis of 110 LIVED mutations from Gremlins mutation testing run, identifying test gaps
event_type: code review
sources:
  - mutations.txt
  - docs/adr/ADR-007-adopt-mutation-testing.md
tags:
  - mutation-testing
  - gremlins
  - test-quality
  - code-coverage
  - test-gaps
  - boundary-conditions
---

# Mutation Testing Analysis: LIVED Mutations

Analyzed the mutation testing results from `mutations.txt`. The Gremlins run completed in 5 minutes 43 seconds with the following summary:

- **Killed:** 537 (tests caught the mutation)
- **Lived:** 110 (tests failed to detect the mutation)
- **Not covered:** 169 (code not exercised by tests)
- **Timed out:** 9
- **Test efficacy:** 83.00%
- **Mutator coverage:** 79.29%

## Categories of LIVED Mutations

### 1. CONDITIONALS_BOUNDARY (Most Common)

These mutations change `<` to `<=` or `>` to `>=`. Tests fail to catch these boundary edge cases:

| File | Line | Context |
|------|------|---------|
| `internal/assert/parser.go` | 37:29 | Loop boundary in `formatError()` - `i < len(p.input)` |
| `internal/cli/assist/composer.go` | 85:8 | Section separator logic - `i < len(sections)-1` |
| `internal/cli/assist/sections.go` | 42, 55, 58, 61, 95, 104 | Multiple boundary checks in section formatting |
| `internal/cli/inspector/detector.go` | 67, 130, 170, 222, 260, 298, 340, 393, 418 | File detection depth/limit boundaries |
| `internal/cli/inspector/metadata.go` | 139-428 (multiple) | Version comparison boundaries |
| `internal/config/config.go` | 331, 343, 344, 352 | Configuration validation boundaries |

### 2. CONDITIONALS_NEGATION

These mutations flip `!=` to `==` or negate boolean conditions:

| File | Line | Context |
|------|------|---------|
| `internal/cli/assist/composer.go` | 85:8 | Same separator logic |
| `internal/cli/check.go` | 77:55 | JSON output error handling |
| `internal/cli/inspector/detector.go` | 382, 393, 418, 419 | Directory traversal conditions |
| `internal/cli/inspector/metadata.go` | 672, 708, 710, 711, 712, 753 | Metadata parsing conditions |
| `internal/cli/inspector/tools.go` | 766, 767 | Tool detection logic |

### 3. INCREMENT_DECREMENT

These mutations change `++` to `--` or vice versa:

| File | Line | Context |
|------|------|---------|
| `internal/cli/inspector/detector.go` | 386:11, 389:9 | Depth counting logic |

### 4. ARITHMETIC_BASE / INVERT_NEGATIVES

These mutations change arithmetic operators (`+` to `-`, etc.):

| File | Line | Context |
|------|------|---------|
| `internal/cli/assist/composer.go` | 85:23 | `len(sections)-1` calculation |
| `internal/config/config.go` | 267:18, 344:7 | Index calculations |
| `internal/cli/inspector/metadata.go` | 710:30 | Version arithmetic |

## High-Priority Files Requiring Test Improvements

Based on LIVED mutation density:

1. **`internal/cli/inspector/metadata.go`** - 30+ LIVED mutations
   - Focus: Version comparison edge cases, boundary conditions

2. **`internal/cli/inspector/detector.go`** - 15+ LIVED mutations
   - Focus: File traversal depth, result limits, path matching

3. **`internal/cli/assist/sections.go`** - 10+ LIVED mutations
   - Focus: Section boundary formatting, empty section handling

4. **`internal/config/config.go`** - 8+ LIVED mutations
   - Focus: Cycle detection, configuration validation bounds

5. **`internal/cli/inspector/tools.go`** - 8+ LIVED mutations
   - Focus: Tool detection heuristics, string matching

## Root Cause Analysis

The LIVED mutations cluster around:

1. **Off-by-one scenarios** - Loop bounds using `<` vs `<=`
2. **Empty/single-element collections** - Edge cases not tested
3. **Error path logic** - Success paths tested but error handling not
4. **Depth/limit parameters** - Boundary values (0, 1, max) not exercised

## Recommended Test Improvements

### Immediate Actions

1. Add boundary value tests for loop conditions:
   - Test with empty input, single element, exactly-at-boundary cases

2. Add explicit edge case tests for `internal/cli/inspector/detector.go`:
   - `maxDepth = 0`, `maxDepth = 1`, `maxResults = 1`

3. Test error paths in `internal/cli/check.go`:
   - Simulate JSON formatting failures

### Medium-Term

4. Increase assertion specificity in existing tests:
   - Don't just check "no error", verify exact output values

5. Add property-based tests for metadata parsing:
   - Version comparison edge cases
   - Malformed input handling

## Metrics Target

Current test efficacy: **83.00%**
Target test efficacy: **90%+**

To reach 90%, approximately 45 of the 110 LIVED mutations need to be killed through improved tests.
