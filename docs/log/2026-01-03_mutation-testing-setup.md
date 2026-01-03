---
summary: Set up Gremlins mutation testing with minimal config optimized for speed
event_type: code
sources:
  - https://gremlins.dev/latest/
  - docs/adr/ADR-007-adopt-mutation-testing.md
tags:
  - mutation-testing
  - gremlins
  - testing
  - test-quality
  - adr-007
---

# Mutation Testing Setup with Gremlins

Implemented ADR-007 by setting up Gremlins for mutation testing with a minimal configuration optimized for fast execution.

## Key Findings

### Configuration Issues Discovered

1. **YAML config format** - The `.gremlins.yaml` config uses flat flags under `unleash:` key, not nested `mutants:` structure as shown in some documentation examples.

2. **Timeout coefficient critical** - Default timeout coefficient (1) caused all mutations to time out even for fast tests. Increased to 10 for reliable results.

3. **filepath.SplitList bug** - During setup, discovered and fixed a bug in `detector.go` where `filepath.SplitList` was incorrectly used for path depth calculation (it splits PATH-like strings, not file paths).

4. **Benchmark test optimization** - Fixed `TestPerformanceTargets` which was hanging due to nested `testing.Benchmark` calls. Replaced with direct `time.Now()` timing for 175x speedup.

### Mutation Testing Results

Running on `internal/assert` package:
- **Efficacy**: 87.01% (67 killed / 77 runnable)
- **Execution time**: ~13 seconds
- **Surviving mutants**: 10 (indicates areas for test improvement)
- **Timed out**: 8 (some mutations cause infinite loops)

### Surviving Mutants (Weak Test Areas)

Located in:
- `eval.go:215:72` - Boundary condition
- `eval.go:276:37` - Arithmetic/negation
- `lexer.go:216:25` - Boundary condition
- `parser.go:37:*` - Multiple conditions (6 survivors)

These indicate tests that execute the code but don't assert on specific boundary conditions.

## Implementation

### Files Created/Modified

1. **`.gremlins.yaml`** - Minimal config:
   - `timeout-coefficient: 10` (prevents false timeouts)
   - `workers: 4` (parallel execution)
   - Enabled: arithmetic-base, conditionals-boundary, conditionals-negation, increment-decrement, invert-negatives
   - Disabled: invert-assignments, invert-bitwise, invert-bwassign, invert-logical, invert-loopctrl, remove-self-assignments

2. **`vibeguard.yaml`** - Added mutation check:
   - Runs on `internal/assert` package only (for speed)
   - 50% efficacy threshold
   - Warning severity (non-blocking)
   - 60s timeout
   - Requires `test` check to pass first

3. **`docs/adr/ADR-007-adopt-mutation-testing.md`** - Status updated to "Accepted"

## Performance Considerations

- Full codebase mutation testing would take several minutes
- Limited to single package (`internal/assert`) for CI integration
- Can run `gremlins unleash ./...` manually for comprehensive analysis
- PR-diff mode (`--diff`) available for incremental testing

## Next Steps

- Consider adding mutation testing to more critical packages
- Investigate surviving mutants and strengthen assertions
- Set up weekly CI job for full mutation analysis
- Track mutation score trends over time
