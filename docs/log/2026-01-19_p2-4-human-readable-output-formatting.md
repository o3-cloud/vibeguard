---
summary: P2.4 - Human-Readable Output Formatting for Triggered Prompts
event_type: code
sources:
  - docs/specs/prompt-feature-spec.md
  - internal/output/formatter.go
  - internal/output/json.go
tags:
  - output-formatting
  - triggered-prompts
  - human-readable
  - json-serialization
  - event-handlers
  - phase-2
---

# P2.4 - Human-Readable Output Formatting for Triggered Prompts

## Overview

Completed P2.4 implementation, which adds output formatting capabilities for triggered prompts. This phase builds on the P2.3 orchestrator integration that already evaluates which prompts trigger for each check outcome. P2.4 focuses solely on presenting triggered prompts to users in both human-readable and JSON formats.

## What Was Implemented

### 1. Human-Readable Output Formatting

**File:** `internal/output/formatter.go`

Added `formatTriggeredPrompts()` method to display triggered prompts in violation output:

```
FAIL  vet (error)

  Fix: Review the error message above
  Log: .vibeguard/log/vet.log

  Triggered Prompts (failure):
  [1] init:
      You are an expert in helping users set up VibeGuard.
      Guide them through setup.

  [2] (inline):
      Also remember to run gofmt before committing

  Advisory: blocks commit
```

**Features:**
- Numbered list of triggered prompts
- Source indicator: prompt ID or "(inline)" for inline content
- Event type label (failure, success, timeout)
- Multi-line prompt content with 6-space indentation
- Positioned between Log and Advisory lines
- Graceful handling of empty prompt arrays (no section displayed)

### 2. JSON Output Support

**File:** `internal/output/json.go`

Added JSON serialization for triggered prompts:

**New Types:**
- `JSONTriggeredPrompt`: Represents a single triggered prompt with event, source, and content fields
- Updated `JSONViolation`: Now includes `TriggeredPrompts` field

**JSON Output Example:**
```json
{
  "violations": [
    {
      "id": "vet",
      "severity": "error",
      "command": "go vet ./...",
      "triggered_prompts": [
        {
          "event": "failure",
          "source": "init",
          "content": "You are an expert in helping users set up VibeGuard..."
        },
        {
          "event": "failure",
          "source": "inline",
          "content": "Also remember to run gofmt before committing"
        }
      ]
    }
  ]
}
```

**Implementation:**
- Converts orchestrator TriggeredPrompt objects to JSONTriggeredPrompt format
- Uses omitempty tag to exclude empty arrays from JSON output
- Maintains all event and source information

### 3. Comprehensive Test Coverage

**Formatter Tests:** `internal/output/formatter_test.go`

Added 6 new test cases:
1. `TestFormatter_QuietMode_WithTriggeredPrompts` - Multiple prompts with mixed sources
2. `TestFormatter_TriggeredPrompts_SinglePrompt` - Single prompt display
3. `TestFormatter_TriggeredPrompts_MultilineContent` - Multi-line content indentation
4. `TestFormatter_TriggeredPrompts_SuccessEvent` - Success event display
5. `TestFormatter_TriggeredPrompts_TimeoutEvent` - Timeout event display
6. `TestFormatter_TriggeredPrompts_NoPrompts` - Empty array handling

**JSON Tests:** `internal/output/json_test.go`

Added 4 new test cases:
1. `TestFormatJSON_WithTriggeredPrompts` - Multiple prompts serialization
2. `TestFormatJSON_WithTriggeredPrompts_SinglePrompt` - Single prompt serialization
3. `TestFormatJSON_WithTriggeredPrompts_EmptyArray` - Empty array handling
4. `TestFormatJSON_WithTriggeredPrompts_DifferentEvents` - Event type variety

**Results:** All 43 tests in output package pass ✓

## Key Design Decisions

### 1. Output Position
Triggered prompts are displayed after the Log line and before the Advisory line. This ensures the flow:
- Check result and severity
- Suggestion/Fix guidance
- Log file location
- **Triggered prompts (NEW)**
- Advisory about commit blocking

### 2. Multi-line Handling
Content is split by newlines and indented with 6 spaces (2 for violation indent + 4 for content indent), making prompts clearly distinguished from other output.

### 3. Source Indicators
- Prompt IDs are displayed as-is (e.g., "init", "code-review")
- Inline content is labeled "(inline)" in parentheses
- This distinction helps users understand whether guidance comes from configured prompts or inline messages

### 4. Event Labels
Each triggered prompt section is labeled with its event type:
- "Triggered Prompts (failure):"
- "Triggered Prompts (success):"
- "Triggered Prompts (timeout):"

This matches the event precedence from the orchestrator and helps users understand why prompts were triggered.

### 5. Empty Array Handling
When no prompts are triggered (empty array):
- Formatter: No "Triggered Prompts" section is displayed
- JSON: Field is omitted (due to omitempty tag) for cleaner output
- This prevents clutter when prompts aren't configured

## Code Quality

- All code passes gofmt formatting check
- All code passes golangci-lint with no issues
- All code passes staticcheck with no issues
- Test coverage includes edge cases and all event types

## Alignment with Architecture

**ADR-004: Code Quality Standards**
- Tests written for all code paths
- Code formatted with gofmt
- Linting passes with golangci-lint
- Proper error handling

**ADR-005: VibeGuard Policy Enforcement**
- Feature integrates seamlessly with check execution
- Respects policy definitions in config

**ADR-006: Git Pre-Commit Hook Integration**
- Changes integrate with existing output formatting
- Backward compatible with non-event-handler configurations

## Integration with Previous Work

**Depends on:**
- P2.3: Orchestrator Integration for Event Triggering
  - Provides TriggeredPrompt evaluation logic
  - Populates CheckResult and Violation with triggered prompts
  - Implements event precedence rules

**Enables:**
- P2.5: JSON Output with Triggered Prompts (partial - JSON support complete)
- P2.6: Built-in Init Prompt Implementation
- P2.7+: Event handler tests and integration tests

## Next Steps

P2.4 is complete. The following phases can proceed:
- **P2.5:** JSON Output with Triggered Prompts (partial - JSON serialization complete)
- **P2.6:** Built-in Init Prompt Implementation
- **P2.7:** Comprehensive Event Handler Tests
- **P2.8:** Integration Tests for Event Handlers

All foundation is in place for triggered prompts to display correctly in both human-readable and JSON output formats.

## Files Modified

1. `internal/output/formatter.go` - Added formatTriggeredPrompts() method
2. `internal/output/formatter_test.go` - Added 6 test cases
3. `internal/output/json.go` - Added JSONTriggeredPrompt type and serialization
4. `internal/output/json_test.go` - Added 4 test cases

## Verification

```bash
# Tests
go test ./internal/output/... -v  # All 43 tests pass

# Code quality
vibeguard check  # All checks pass (lint, staticcheck, fmt)
```
