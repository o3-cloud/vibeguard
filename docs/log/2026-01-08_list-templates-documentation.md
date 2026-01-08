---
summary: Updated documentation for --list-templates flag across all major docs
event_type: code
sources:
  - README.md
  - docs/CLI-REFERENCE.md
  - docs/GETTING_STARTED.md
tags:
  - documentation
  - cli
  - templates
  - init-command
  - user-guide
---

# Documentation Update: --list-templates Flag

## Summary

Updated comprehensive documentation for the `vibeguard init --list-templates` flag across three key documentation files. This flag allows users to discover available templates without creating a configuration file.

## Changes Made

### README.md
- Added `--template` flag to the init command flags table
- Added `--list-templates` flag with description and default value
- Updated examples to show template discovery and usage

### docs/CLI-REFERENCE.md
- Created new section for `--list-templates` flag
- Documented the output format with example templates
- Updated behavior section to include template listing logic
- Clarified that flag is useful for discovering templates before using `--template`

### docs/GETTING_STARTED.md
- Added `vibeguard init --list-templates` to the Quick Start section
- Placed it before other initialization examples to suggest discovering available options first

## Key Features Documented

- **Purpose**: Discover available templates without creating config
- **Usage**: `vibeguard init --list-templates`
- **Output**: Human-readable table showing all available templates with descriptions
- **Behavior**: Command exits after displaying templates (no side effects)

## Testing

- Ran `go test ./internal/cli -v -run "Init"` - all 8 init tests pass
- Specifically verified `TestRunInit_ListTemplatesFlag` passes
- Feature confirmed working as documented

## Relation to Task

This work completes task **vibeguard-938**: Update documentation for --list-templates flag (CLI-REFERENCE, GETTING_STARTED, README)

The documentation now consistently covers:
1. How to discover templates (`vibeguard init --list-templates`)
2. How to use templates (`vibeguard init -t template-name`)
3. Available template options for different project types
