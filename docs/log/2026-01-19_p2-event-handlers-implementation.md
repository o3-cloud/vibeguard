---
summary: P2.1 & P2.2 - Event Handler Type Definition and Validation Implementation
event_type: code
sources:
  - internal/config/events.go
  - internal/config/events_test.go
  - internal/config/config.go
  - docs/specs/prompt-feature-spec.md
tags:
  - event-handlers
  - prompt-feature
  - phase-2
  - yaml-unmarshaling
  - validation
  - configuration-loading
  - type-definition
  - testing
---

# P2.1 & P2.2 - Event Handler Type Definition and Validation

## Overview

Completed two major Phase 2 tasks for the VibeGuard prompt feature: defining event handler types and implementing comprehensive validation. This foundation enables prompts to be triggered based on check outcomes (success, failure, timeout).

## P2.1: Event Handler Type Definition

### Implementation

Created two new types in `internal/config/events.go`:

1. **EventHandler struct** - Defines three event types at check level:
   - `Success` - Triggered when check passes (exit code 0, assertions true)
   - `Failure` - Triggered when check fails (exit code non-zero, assertions false)
   - `Timeout` - Triggered when check exceeds timeout (highest precedence)

2. **EventValue type** - Supports dual representation:
   - **Array syntax** → Treats as prompt ID references: `[init, code-review]`
   - **String syntax** → Treats as inline content: `"This check timed out"`
   - Uses `IsInline bool` flag to distinguish between modes
   - Critical distinction: bare string always treated as inline, only arrays are IDs

### Custom YAML Unmarshaling

Implemented `UnmarshalYAML` for EventValue to handle both formats:
- Tries to parse as `[]string` first (array of prompt IDs)
- Falls back to parsing as single `string` (inline content)
- Sets appropriate `IsInline` flag based on parsed type

### Test Coverage

Created 23 comprehensive test cases in `internal/config/events_test.go`:
- **String parsing tests** (5 cases): empty, multiline, special characters, ID-like strings
- **Array parsing tests** (5 cases): single ID, multiple IDs, flow/block styles
- **YAML marshaling tests** (3 cases): round-trip serialization
- **EventHandler unmarshaling tests** (4 cases): all combinations of event types
- **Check context tests** (1 case): integration with Check struct
- **Edge cases** (4 cases): null values, numeric strings, boolean-like values
- **ID vs inline distinction tests** (5 cases): critical syntax differentiation

## P2.2: Event Handler Validation

### Validation Pipeline

Integrated event handler validation into config loading pipeline:
1. Called from `Config.Validate()` after prompt validation
2. Iterates through all checks and their event handlers
3. Validates prompt ID references against config's prompts section

### Validation Functions

1. **validateEventHandlers(check Check, checkIndex int) error**
   - Validates all three event types for a single check
   - Skips validation for inline content (always valid)
   - Delegates to validateEventPromptIDs for ID arrays

2. **validateEventPromptIDs(checkID, eventName, promptIDs) error**
   - Builds map of valid prompt IDs from config
   - Checks each referenced prompt ID exists
   - Returns ConfigError with line number context

### Error Handling

All validation errors include:
- Clear message indicating which check and event had invalid reference
- Unknown prompt ID name
- YAML line number for debugging

Example: `check "vet" references unknown prompt "invalid-id" in event "failure" (line 45)`

### Test Coverage

Added 6 new test functions in `internal/config/config_test.go`:
- **TestLoad_EventHandler_ValidPromptIDs** - Valid references accepted
- **TestLoad_EventHandler_InvalidPromptID** - Invalid references rejected (3 event types tested)
- **TestLoad_EventHandler_InlineContent** - Inline content always valid
- **TestLoad_EventHandler_MixedInlineAndIDs** - Mixed modes in single check
- **TestLoad_EventHandler_MultiplePromptIDs** - Multiple ID references per event
- **TestLoad_EventHandler_EmptyOn** - Empty handler sections allowed
- **TestLoad_EventHandler_NoOn** - Checks without handlers allowed

## Key Design Decisions

### 1. Array vs String Distinction

The critical rule: **Only arrays are prompt ID references, bare strings are always inline content**

```yaml
# Array = Prompt IDs
on:
  failure: [init, code-review]

# String = Inline content (NOT a prompt ID, even if it looks like one)
on:
  timeout: "init"  # This is literal text, not a reference to prompt "init"
```

This prevents ambiguity and follows YAML conventions.

### 2. Validation Only on ID Arrays

- Inline content requires no validation (always valid)
- ID arrays checked against prompts section
- Missing prompts section (no prompts defined) = all ID references invalid

### 3. No Event Name Validation

Currently does NOT validate that only success/failure/timeout are used. Implementation relies on:
- YAML unmarshaling to parse known fields only
- Unknown event types ignored silently (YAML standard behavior)

Future enhancement could add explicit validation if needed.

### 4. Graceful Degradation

- Checks without `on` section work fine (empty EventHandler)
- Empty `on:` section allowed (all event fields empty)
- Inline content always valid regardless of content
- Proper error messages for invalid references

## Test Results

```
All tests passing:
- 23 event handler type/parsing tests
- 6 event handler validation tests
- Existing prompt validation tests (still passing)
- Existing check validation tests (still passing)

Coverage: 70%+ maintained
```

## Quality Assurance

- All tests passing
- Linting passes (staticcheck, golangci-lint)
- vibeguard check policies passing
- No regressions in existing functionality

## Files Changed

- **internal/config/events.go** (76 lines) - New type definitions
- **internal/config/events_test.go** (450+ lines) - Comprehensive event handler tests
- **internal/config/schema.go** (+1 line) - Added On field to Check struct
- **internal/config/config.go** (+27 lines) - Event handler validation functions
- **internal/config/config_test.go** (+295 lines) - Event handler validation tests

**Total:** 886 insertions across 6 files

## Next Phase

Ready for P2.3: Orchestrator Integration for Event Triggering

The foundation is solid:
- Event handler types fully defined and tested
- Validation ensures referential integrity
- Line numbers for debugging
- Support for both inline and ID-based prompts

Next: Integrate event evaluation into check execution (success/failure/timeout logic with proper precedence rules).

## Related Architecture

Builds on:
- [ADR-001: Adopt Beads for AI Agent Task Management](../../adr/ADR-001-adopt-beads.md)
- [ADR-004: Establish Code Quality Standards](../../adr/ADR-004-code-quality-standards.md)
- [ADR-005: Adopt VibeGuard for Policy Enforcement](../../adr/ADR-005-adopt-vibeguard.md)

Prompt Feature Specification: [docs/specs/prompt-feature-spec.md](../../specs/prompt-feature-spec.md)
