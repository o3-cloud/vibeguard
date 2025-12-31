---
summary: Implemented topological sort execution ordering in the orchestrator to run checks in dependency order
event_type: code
sources:
  - internal/orchestrator/orchestrator.go
  - internal/orchestrator/orchestrator_test.go
  - internal/orchestrator/graph.go
tags:
  - vibeguard-v3m.2
  - orchestrator
  - topological-sort
  - dependency-management
  - execution-ordering
---

# Topological Sort Execution Ordering (vibeguard-v3m.2)

## Overview

Integrated the existing dependency graph (Kahn's algorithm) into the orchestrator's `Run()` method to execute checks in correct topological order based on their `requires` dependencies.

## Changes Made

### orchestrator.go

1. **Integrated BuildGraph into Run()**: The orchestrator now builds a dependency graph at the start of execution using `BuildGraph()` which implements Kahn's algorithm.

2. **Level-by-level execution**: Checks are now executed level by level according to the topological sort. Within each level, checks can theoretically run in parallel (to be implemented in vibeguard-v3m.3).

3. **Dependency tracking**: Added `passedChecks` map to track which checks have passed, enabling downstream checks to verify their dependencies succeeded.

4. **Skip failed dependencies**: If a check's required dependency failed, the check is skipped with:
   - Exit code: -1
   - Suggestion: "Skipped: required dependency failed"
   - Recorded as a violation

5. **Extracted calculateExitCode helper**: Created a helper method to reduce code duplication between the main Run loop and fail-fast early return.

6. **Fixed defer leak**: Changed timeout cancellation from `defer cancel()` to explicit `cancel()` call after execution to avoid accumulating deferred functions in the loop.

## Test Coverage

Added 12 new test cases for topological sort execution:

- `TestRun_WithDependencies_ExecutesInOrder` - Linear chain a->b->c
- `TestRun_DiamondDependency_ExecutesInCorrectOrder` - Diamond pattern
- `TestRun_DependencyFails_SkipsDependent` - Skip when dep fails
- `TestRun_MultipleDependenciesOneFails_SkipsDependent` - Skip when any dep fails
- `TestRun_IndependentChecks_AllExecute` - No deps, all in level 0
- `TestRun_CyclicDependency_ReturnsError` - Cycle detection
- `TestRun_UnknownDependency_ReturnsError` - Missing dep detection
- `TestRun_FailFast_WithDependencies_StopsCorrectly` - Fail-fast with deps
- `TestRun_ComplexDependencyGraph_CorrectOrder` - 7-node complex graph
- `TestRun_DependencyChain_MiddleFails_SkipsDownstream` - Cascading skip

## Behavior Summary

| Scenario | Behavior |
|----------|----------|
| No dependencies | All checks run in level 0 (original order preserved) |
| Linear chain | Executes in topological order |
| Diamond pattern | Respects all dependency edges |
| Dependency fails | Downstream checks skipped with exit code -1 |
| Multiple deps, one fails | Dependent check skipped |
| Cyclic dependency | Returns error from BuildGraph |
| Unknown dependency | Returns error from BuildGraph |
| Fail-fast mode | Stops execution, returns early |

## Next Steps

- **vibeguard-v3m.3**: Implement parallel execution within levels using errgroup
- **vibeguard-v3m.4**: Implement fail-fast mode refinements
- **vibeguard-v3m.5**: Add timeout handling per-check

## Related

- vibeguard-v3m.1 (completed): Dependency graph construction with Kahn's algorithm
- ADR-003: Go implementation language choice
