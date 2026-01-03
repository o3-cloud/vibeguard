---
summary: Added file field support to CheckRecommendation structs and prompt generation
event_type: code
sources:
  - internal/cli/assist/composer.go
  - internal/cli/inspector/recommendations.go
  - internal/cli/inspector/prompt.go
  - internal/cli/assist/sections.go
  - internal/cli/assist/validator_guide.go
tags:
  - ai-assisted-setup
  - prompt-generation
  - file-field
  - check-recommendation
  - bug-fix
---

# File Field Support in Prompt Generation

## Context

Task `vibeguard-tcf` identified that the `file` field was documented in `validator_guide.go` and supported in the Check schema (`config/schema.go:20`), but was missing from the prompt generation system. This meant AI agents wouldn't know they could use the `file` field to read check output from files instead of command stdout.

## Investigation

The `file` field was already documented in:
- `validator_guide.go` lines 127-128: "**file** (string): File to read output from instead of command stdout"
- `sections.go` line 147 in `ConfigRequirementsSection()`: "- **file:** Path to read output from instead of command stdout"

However, the `CheckRecommendation` structs that carry recommendation data through the system did not include the `File` field:
- `assist/composer.go` `CheckRecommendation` struct (lines 30-41)
- `inspector/recommendations.go` `CheckRecommendation` struct (lines 6-20)

## Changes Made

1. **Added `File` field to `assist/composer.go` CheckRecommendation struct** (line 35):
   ```go
   File string // File to read output from instead of command stdout
   ```

2. **Added `File` field to `inspector/recommendations.go` CheckRecommendation struct** (line 11):
   ```go
   File string // File to read output from instead of command stdout
   ```

3. **Updated `inspector/prompt.go` conversion function** to pass through the File field (line 67):
   ```go
   File: r.File,
   ```

4. **Updated `assist/sections.go` RecommendationsSection()** to display File field when present (lines 91-93):
   ```go
   if rec.File != "" {
       sb.WriteString(fmt.Sprintf("**File:** `%s`\n", rec.File))
   }
   ```

## Testing

- All existing tests pass
- Build succeeds with no errors
- The file field is now properly propagated through the entire recommendation pipeline

## Future Considerations

Currently, no built-in recommendations use the `file` field. When tool recommendations are added that write output to files (e.g., coverage reports, test result files), those recommendations should populate the `File` field so AI agents can see examples of its usage.
