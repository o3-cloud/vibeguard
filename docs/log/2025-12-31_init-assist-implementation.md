---
summary: Implemented the init --assist command for AI agent-assisted setup, completing task vibeguard-9mi.10
event_type: code
sources:
  - internal/cli/init.go
  - internal/cli/inspector/prompt.go
  - internal/cli/init_test.go
  - docs/ai-assisted-setup.md
tags:
  - ai-assisted-setup
  - cli
  - init-command
  - inspector
  - prompt-generation
  - feature-implementation
---

# Implemented init --assist Command for AI Agent-Assisted Setup

## Summary

Completed task `vibeguard-9mi.10: Phase 4: CLI Integration - init --assist Command`. This feature enables AI agents (like Claude Code) to receive a comprehensive setup prompt that guides them through generating a valid `vibeguard.yaml` configuration for any detected project.

## Changes Made

### 1. Created `internal/cli/inspector/prompt.go`
- Exported `GenerateSetupPrompt` function (previously only in test file)
- Generates a Claude Code-friendly setup prompt with:
  - Project analysis (type, confidence, detected tools)
  - Recommended checks based on detected tools
  - Configuration requirements and validation rules
  - Language-specific examples
  - Task instructions for AI agents

### 2. Updated `internal/cli/init.go`
- Added `--assist` flag to enable AI-assisted setup mode
- Added `--output` (`-o`) flag to save prompt to a file
- Implemented `runAssist` function that:
  - Validates the target directory
  - Runs the full inspector pipeline (detector, scanner, extractor, recommender)
  - Generates and outputs the setup prompt

### 3. Updated `internal/cli/check.go`
- Extended `ExitError` struct to include a `Message` field
- Allows descriptive error messages with specific exit codes

### 4. Updated `internal/cli/inspector/prompt_test.go`
- Changed to use exported `GenerateSetupPrompt` function
- Removed duplicate implementation

### 5. Created `internal/cli/init_test.go`
- Added comprehensive tests for `runAssist` function:
  - Success case on vibeguard project itself
  - Error case for non-existent directory (exit code 2)
  - Error case for undetectable project type (exit code 2)
  - Success case with `--output` flag
  - Error case for file instead of directory (exit code 2)

## Exit Codes

Per the task specification:
- `0`: Success
- `2`: Invalid directory or undetectable project type
- `3`: Runtime error (tool scanning, metadata extraction failures)

## Usage Examples

```bash
# Generate setup prompt to stdout
vibeguard init --assist

# Save prompt to a file
vibeguard init --assist --output setup-prompt.md

# Inspect a specific directory
vibeguard init --assist /path/to/project
```

## Testing

All tests pass:
- 5 new CLI tests for `runAssist` function
- All existing inspector tests continue to pass
- Manual testing confirmed proper prompt generation

## Next Steps

The following tasks in the AI-assisted setup feature epic remain:
- `vibeguard-9mi.7`: Design Prompt Structure
- `vibeguard-9mi.8`: Implement Composer
- `vibeguard-9mi.9`: Validation Guide Templates
- `vibeguard-9mi.15`: Performance Optimization
