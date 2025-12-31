---
summary: Implemented Go template syntax support for suggestions in vibeguard
event_type: code
sources:
  - https://pkg.go.dev/text/template
  - docs/adr/ADR-002-adopt-conventional-commits.md
  - vibeguard-c9m.3
tags:
  - template-rendering
  - suggestions
  - grok-extraction
  - feature-implementation
  - text-template
---

# Templated Suggestions Implementation

## Overview

Completed implementation of vibeguard-c9m.3: Templated suggestions. This feature enables Go template syntax in suggestion fields, making grok-extracted values and config vars available as `{{.varname}}` for dynamic, context-aware error messages.

## Implementation Details

### Core Changes

1. **internal/config/interpolate.go**
   - Upgraded `InterpolateWithExtracted()` from simple string replacement to proper Go template parsing
   - Uses `text/template` package for full template syntax support
   - Merges config vars and extracted values with proper precedence (config vars win on conflicts)
   - Graceful error handling - falls back to original string on parse/execution errors

2. **internal/output/formatter.go**
   - Updated `formatViolation()` to interpolate suggestions with extracted values
   - Now renders templated suggestions in quiet mode (previously only verbose mode had this)
   - Consistent template rendering across all output formats

### Key Features

- **Full Go Template Syntax**: Supports conditionals (`{{if}}...{{end}}`), functions, ranges, and all standard template features
- **Data Availability**: Grok-extracted values and config vars both available in templates
- **Precedence Rules**: Config vars take precedence over extracted values if there's a conflict
- **Error Resilience**: Invalid templates gracefully fall back to original string
- **Dual Output Modes**: Works in both quiet and verbose output modes, plus JSON mode

### Template Examples

```yaml
# Simple value interpolation
suggestion: "Coverage is {{.coverage}}%, needs 80%"

# Complex message with multiple values
suggestion: "Fix {{.error_type}} error in {{.file}}:{{.line}} - {{.message}}"

# With config vars
suggestion: "{{.tool}}: {{.error_msg}}"

# Conditional (Go template syntax)
suggestion: "{{if gt .coverage 80}}Excellent coverage!{{else}}Coverage too low{{end}}"
```

## Testing & Validation

- ✅ All 100+ existing tests pass
- ✅ New template interpolation tests cover edge cases
- ✅ Config var precedence verified
- ✅ Graceful fallback behavior confirmed
- ✅ Build successful with no warnings

## Data Flow

```
Grok Extraction → CheckResult.Extracted
                ↓
Check Fails → Violation created with:
           - Check.Suggestion (template string)
           - Violation.Extracted (grok values)
           ↓
Output Formatting:
  - InterpolateWithExtracted() parses template
  - Merges extracted + config vars
  - Renders final suggestion message
```

## Design Decisions

1. **Template Parsing on Render**: Templates are parsed at output time rather than config load time, allowing flexibility and deferring validation
2. **Graceful Fallback**: Invalid templates return original string rather than failing - supports legacy simple placeholders
3. **Config Var Precedence**: Ensures predictable behavior when same key exists in both sources
4. **Backward Compatible**: Existing simple `{{.varname}}` strings still work (Go templates support this syntax)

## Related ADRs

- ADR-002: Adopt Conventional Commits (for semantic versioning of this feature)
- ADR-004: Code Quality Standards (all tests pass with 70%+ coverage requirement)
- ADR-005: Adopt Vibeguard for Policy Enforcement

## Next Steps

- Consider adding template validation in config validation phase for early error detection
- Document template syntax in user documentation
- Consider caching compiled templates for performance if needed
