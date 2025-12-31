---
summary: Integrated validator_guide.go with prompt.go to eliminate duplicate validation rules and provide comprehensive guidance
event_type: code
sources:
  - internal/cli/inspector/prompt.go
  - internal/cli/assist/validator_guide.go
  - internal/cli/inspector/prompt_test.go
  - internal/cli/init_test.go
tags:
  - ai-assisted-setup
  - prompt-generation
  - validation-rules
  - refactoring
  - vibeguard-1i7
---

# Validator Guide Integration

Completed task vibeguard-1i7: Integrate validator_guide.go with prompt.go generation.

## Problem

The `prompt.go` file in `internal/cli/inspector/` generated its own validation rules section with 10 basic rules, while `validator_guide.go` in `internal/cli/assist/` contained comprehensive validation documentation. This created:

1. **Duplicated information** - Same rules documented in two places
2. **Incomplete coverage** - prompt.go had 10 rules, guide had much more detail
3. **Maintenance burden** - Fixing rules required updating both places
4. **Missing documentation** - The `file` field was documented in validator_guide.go but not in prompt.go

## Solution

Integrated the comprehensive validation guide from `validator_guide.go` into the prompt template:

1. Added import for `github.com/vibeguard/vibeguard/internal/cli/assist` package
2. Replaced inline "Configuration Requirements" and "Validation Rules" sections with `{{.ValidationGuide}}`
3. Used `assist.NewValidationGuide().GetFullGuide()` to inject comprehensive rules
4. Kept Go-Specific Examples section for practical context

## Changes Made

### internal/cli/inspector/prompt.go
- Added import for assist package
- Added `ValidationGuide string` field to template data struct
- Replaced hardcoded validation rules with `{{.ValidationGuide}}` template variable
- Removed duplicate "Configuration Requirements" and "Validation Rules" sections

### internal/cli/inspector/prompt_test.go
- Updated required sections check to verify new section headers:
  - `## YAML Syntax Requirements`
  - `## Check Structure Requirements`
  - `## Dependency Validation Rules`
  - `## Variable Interpolation Rules`
  - `## Explicit DO NOT List`

### internal/cli/init_test.go
- Updated `TestRunAssist_OutputToFile` to check for new section headers

## Result

Generated prompts now include comprehensive validation documentation:
- **Before:** ~8,500 characters with 10 basic rules
- **After:** ~12,638 characters (~3,159 tokens) with full validation guide

The prompt now includes:
- YAML syntax requirements with common errors to avoid
- Complete check structure documentation including the `file` field
- Detailed dependency validation rules with examples
- Variable interpolation rules with grok-extracted values
- Explicit "DO NOT" list for AI agents

## Verification

- All tests pass: `go test ./...`
- Build succeeds: `go build ./...`
- Vibeguard checks pass (except golangci-lint warning due to tool not being installed locally)
