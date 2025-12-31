---
summary: Implemented cyclic dependency validation for check requires in vibeguard config
event_type: code
sources:
  - internal/config/config.go
  - internal/config/config_test.go
tags:
  - vibeguard-002
  - validation
  - cyclic-dependency
  - config
  - bug-fix
---

# Cyclic Dependency Validation Implementation

Completed task vibeguard-002: Add cyclic dependency validation for check requires.

## Problem

The existing validation in `config.go` checked for self-references in `requires` (e.g., check A requires A), but did not detect cyclic dependencies involving multiple checks (e.g., A requires B, B requires C, C requires A). This could cause infinite loops during orchestration.

## Solution

Implemented a DFS-based cycle detection algorithm in `validateNoCycles()` method:

1. Build an adjacency list from check IDs to their required check IDs
2. Use three-state tracking: unvisited (0), visiting (1), visited (2)
3. During DFS traversal, if we encounter a node in "visiting" state, a cycle exists
4. Track the current path for informative error messages showing the cycle

## Implementation Details

Added two new functions to `internal/config/config.go`:

- `validateNoCycles()` - Main cycle detection using DFS with state tracking
- `formatCycle()` - Helper to format cycle path for error messages (e.g., "a -> b -> c -> a")

The validation is called at the end of `Validate()` after all other checks pass.

## Test Coverage

Added comprehensive tests in `internal/config/config_test.go`:

- `TestLoad_CyclicDependency_TwoNodes` - A requires B, B requires A
- `TestLoad_CyclicDependency_ThreeNodes` - A requires B, B requires C, C requires A
- `TestLoad_CyclicDependency_PartialCycle` - Cycle in subgraph (B->C->B) not involving entry point
- `TestLoad_NoCycle_ValidDAG` - Valid DAG with multiple dependencies
- `TestLoad_NoCycle_DiamondDependency` - Diamond pattern (valid, no cycle)
- `TestFormatCycle` - Unit test for cycle formatting

## Error Message Example

When a cycle is detected, the error message clearly shows the cycle path:

```
configuration validation failed: cyclic dependency detected: a -> b -> c -> a
```

## Next Steps

- Close vibeguard-002 bead
- This validation will be particularly important when Phase 3 (Orchestration) is implemented
