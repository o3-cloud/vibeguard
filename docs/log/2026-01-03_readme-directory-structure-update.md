---
summary: Successfully updated README.md to match the actual codebase structure and list all Architecture Decision Records
event_type: code
sources:
  - docs/adr/ADR-001-adopt-beads.md
  - docs/adr/ADR-002-adopt-conventional-commits.md
  - docs/adr/ADR-003-adopt-golang.md
  - docs/adr/ADR-004-code-quality-standards.md
  - docs/adr/ADR-005-adopt-vibeguard.md
tags:
  - documentation
  - readme
  - structure
  - project-structure
  - architecture-decisions
  - adr
---

# README Directory Structure Update - Task vibeguard-ali

## Summary

Completed the P1 bug fix task to update README.md directory structure to match the actual codebase layout.

## Changes Made

### Project Structure Section
Updated the Project Structure diagram to accurately reflect the current directory layout:

**Added directories:**
- `bin/` - Compiled binary output
- `spikes/` - Research and prototyping work with subdirectories (config, executor, opa, orchestrator)
- `internal/version/` - Version information and constants package

**Added internal packages:**
- `executor/` - Check execution engine
- `grok/` - Grok pattern extraction and matching
- `orchestrator/` - Check orchestration and dependency management

**Updated docs/ structure:**
- Added `sample-prompts/` subdirectory
- Added `SECURITY.md` file reference
- Added `CHANGELOG.md` file reference

### Architecture Decisions Section
Updated to list all seven existing ADRs instead of just the first three:
- ADR-001: Adopt Beads for AI Agent Task Management
- ADR-002: Adopt Conventional Commits
- ADR-003: Adopt Go as the Primary Implementation Language
- ADR-004: Establish Code Quality Standards and Tooling
- ADR-005: Adopt VibeGuard for Policy Enforcement in CI/CD
- ADR-006: Integrate VibeGuard as Git Pre-Commit Hook for Policy Enforcement
- ADR-007: Adopt Gremlins for Mutation Testing

## Verification

All seven ADRs were verified to exist in `docs/adr/`:
- ✓ ADR-001-adopt-beads.md
- ✓ ADR-002-adopt-conventional-commits.md
- ✓ ADR-003-adopt-golang.md
- ✓ ADR-004-code-quality-standards.md
- ✓ ADR-005-adopt-vibeguard.md
- ✓ ADR-006-integrate-vibeguard-as-claude-code-hook.md
- ✓ ADR-007-adopt-mutation-testing.md

## Commit

Committed with message: `docs(readme): update project structure and architecture decisions sections`

Commit hash: a0b27d3

## Impact

The README.md now accurately reflects the current project structure, making it easier for new contributors and maintainers to understand the codebase organization. The complete list of ADRs provides better visibility into architectural decisions made during the project.

## Related Decisions

- ADR-002: Adopt Conventional Commits (used for commit message formatting)
- ADR-004: Establish Code Quality Standards and Tooling (documentation as part of code quality)
