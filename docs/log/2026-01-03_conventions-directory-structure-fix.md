---
summary: Fixed CONVENTIONS.md to reflect actual internal package structure
event_type: code
sources:
  - docs/CONVENTIONS.md
  - internal/assert
  - internal/cli
  - internal/config
  - internal/executor
tags:
  - documentation
  - conventions
  - bug-fix
  - directory-structure
  - package-layout
---

# Fixed CONVENTIONS.md Directory Structure

## Overview

Completed task vibeguard-y2t by updating CONVENTIONS.md (lines 13-28) to accurately reflect the actual internal package structure of VibeGuard.

## Changes Made

### Removed Non-Existent Packages
The documentation referenced packages that do not exist in the codebase:
- `policy` - removed
- `judge` - removed
- `runner` - removed

### Added Actual Packages
Updated the directory structure diagram to include the actual internal packages:
- `assert` - Assertion expression evaluation
- `executor` - Policy execution logic
- `grok` - Grok pattern matching and parsing
- `output` - Output formatting and rendering
- `version` - Version information

### Removed Non-Existent Directories
- `pkg/` directory - Removed from structure and Package Layout section
- `tests/` directory - Removed from structure (unit tests are colocated with code)

## Verification

Performed thorough investigation to confirm changes:

1. **Actual internal directory contents**: assert, cli, config, executor, grok, orchestrator, output, version
2. **Directory verification**: Confirmed no `pkg/` or `tests/` directories exist
3. **Codebase search**: Investigated all references to `pkg/` in the codebase
   - Found only in examples and test data (e.g., validator_guide.go shows example patterns)
   - Not part of VibeGuard's actual structure

## Result

CONVENTIONS.md now accurately reflects the actual codebase structure and will no longer mislead developers or documentation consumers about the project layout.

## Related Issues

- Closes: vibeguard-y2t (Update CONVENTIONS.md directory structure)
