---
summary: Implemented validation guide templates for AI agent-assisted setup feature
event_type: code
sources:
  - internal/cli/assist/validator_guide.go
  - internal/cli/assist/validator_guide_test.go
  - docs/log/2025-12-31_agent-assisted-setup-implementation-spec.md
  - internal/config/config.go
  - internal/cli/inspector/prompt.go
tags:
  - ai-assisted-setup
  - validation
  - prompt-engineering
  - documentation
  - vibeguard-9mi.9
---

# Validator Guide Templates Implementation

Completed task vibeguard-9mi.9: Phase 3 - Prompt Composition - Validation Guide Templates.

## Overview

Created a new `internal/cli/assist/` package with comprehensive validation rule templates for AI agents generating VibeGuard configurations. The validation guide provides detailed, structured documentation that AI agents can reference when creating vibeguard.yaml files.

## Implementation Details

### Package Structure

Created `internal/cli/assist/validator_guide.go` with:

1. **ValidationGuide struct** - Container for all validation rule sections
2. **NewValidationGuide()** - Factory function
3. **GetFullGuide()** - Returns complete guide as single string

### Validation Rule Sections

Implemented five comprehensive template sections:

1. **YAMLSyntaxRules** - YAML syntax requirements including:
   - Required top-level keys (version, vars, checks)
   - String quoting rules
   - Common syntax errors to avoid
   - Example valid structure

2. **CheckStructureRules** - Check definition requirements:
   - Required fields (id, run)
   - Optional fields (grok, assert, severity, suggestion, requires, timeout, file)
   - Field descriptions and valid values
   - Complete check example

3. **DependencyValidationRules** - Dependency graph rules:
   - No self-references
   - No circular dependencies
   - Examples of invalid and valid dependency patterns
   - Execution order explanation

4. **VariableInterpolationRules** - Variable usage:
   - {{.varname}} syntax documentation
   - Fields where variables can be used
   - Grok-extracted value precedence
   - Common mistakes to avoid

5. **ExplicitDoNotList** - Prohibitions for AI agents:
   - YAML structure constraints
   - Check definition constraints
   - Command safety guidelines
   - General best practices

### Supporting Data Structures

Added helper variables for validation:
- `AssertionOperators` - Supported operators (==, !=, <, <=, >, >=, &&, ||, !)
- `SpecialAssertionVariables` - exit_code, stdout, stderr
- `SupportedSeverities` - error, warning
- `GrokPatternExamples` - Common grok patterns with examples
- `CommonTimeoutValues` - Recommended timeouts by check type

## Test Coverage

Created comprehensive unit tests in `validator_guide_test.go`:
- 14 test functions covering all aspects
- Tests for content presence and structure
- Tests for helper variable completeness
- Readability checks (headers, subheaders, code examples)

All tests pass:
```
=== RUN   TestNewValidationGuide
--- PASS: TestNewValidationGuide (0.00s)
... (14 tests total)
PASS
ok  	github.com/vibeguard/vibeguard/internal/cli/assist	0.166s
```

## Integration Notes

The ValidationGuide can be integrated with the existing prompt.go in the inspector package to provide more detailed validation guidance. The current prompt.go has a basic "Validation Rules" section that can be enhanced or replaced with the more comprehensive templates.

## Next Steps

Related tasks remaining in the AI-assisted setup feature:
- vibeguard-9mi.8: Implement Composer - integrate ValidationGuide into prompt generation
- vibeguard-9mi.7: Design Prompt Structure - use templates to structure prompts

## Files Changed

- Created: `internal/cli/assist/validator_guide.go`
- Created: `internal/cli/assist/validator_guide_test.go`
