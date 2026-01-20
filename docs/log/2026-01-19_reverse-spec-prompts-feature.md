---
summary: Reverse specification derived from prompts branch implementation - comprehensive feature definition based on completed Phase 1 and Phase 2 work
event_type: code
sources:
  - internal/cli/prompt.go
  - internal/config/schema.go
  - internal/config/config.go
  - internal/config/events.go
  - internal/orchestrator/orchestrator.go
  - vibeguard.yaml
tags:
  - prompts
  - feature-specification
  - reverse-engineering
  - phase-1
  - phase-2
  - event-handlers
  - cli-design
  - configuration
---

# Reverse Specification: Prompts Feature Implementation

## Overview

Based on comprehensive analysis of all changes in the `prompts` branch (commits bd68969 through 40aebbb), this document defines the complete specification for the Prompt Feature as implemented in VibeGuard. The feature enables storing, retrieving, and distributing guided prompts for AI agents and users, with automatic triggering on check execution outcomes.

## Implementation Summary

**Scope:** Full implementation of Phase 1 (CLI and configuration) and Phase 2 (event handlers and built-in prompts)

**Commits Analyzed:**
- bd68969: Document completion of Phase 1 and Phase 2
- ad2bfc9: P2.4 - Human-Readable Output Formatting
- a4c29aa: P2.3 - P2.8 Event Handler Implementation Complete
- 786d27c: P2.3 - Orchestrator Integration for Event Triggering
- 99ba537: P2.1 & P2.2 - Event Handler Type Definition and Validation
- 83c941b: P1.3 - Configuration Loading & Validation
- 8c8e397: P1.2 - CLI Command Implementation
- 8c8e397: P1.1 - Prompt Data Structure and Schema

**Files Changed:** 36 files modified/created
**Total Insertions:** ~7,940 lines
**New Test Coverage:** 663 lines (prompt_test.go) + 447 lines (events_test.go)

## Phase 1: Core Prompt Feature (Complete)

### Data Model

**Prompt Structure** (internal/config/schema.go):
```go
type Prompt struct {
    ID          string   `yaml:"id"`
    Description string   `yaml:"description,omitempty"`
    Content     string   `yaml:"content"`
    Tags        []string `yaml:"tags,omitempty"`
}
```

**Configuration Integration**:
- Prompts stored in `vibeguard.yaml` under top-level `prompts` section
- Optional field (`omitempty`) - fully backward compatible
- Integrated into `Config` struct alongside existing `checks` and `vars`

### CLI Command Implementation

**Command:** `vibeguard prompt [prompt-id]`

**Behaviors Implemented:**

1. **List Mode** (no arguments):
   ```bash
   vibeguard prompt
   ```
   Outputs alphabetically sorted prompt IDs

2. **Verbose List** (`-v` or `--verbose` flag):
   ```bash
   vibeguard prompt -v
   ```
   Shows ID, description, and tags for each prompt

3. **JSON Output** (`--json` flag):
   ```bash
   vibeguard prompt --json
   ```
   Outputs structured JSON with metadata (excludes content by default)

4. **Specific Prompt** (with ID argument):
   ```bash
   vibeguard prompt init
   ```
   Outputs raw prompt content to stdout (optimized for piping)

**Error Handling:**
- "no prompts defined in configuration" - When prompts section empty
- "prompt not found: {id}" - When requested ID doesn't exist
- "no config file found" - When vibeguard.yaml missing

### Configuration Validation

**Prompt-Specific Rules** (internal/config/config.go):
- All prompt IDs must be unique within config
- IDs must match pattern: start with letter/underscore, contain only alphanumeric/underscore/hyphen
- Tags must be lowercase alphanumeric with hyphens
- Content field is required (non-empty)
- Validation errors include line numbers for debugging

**Integration:**
- Prompts validated as part of standard `Config.Validate()` flow
- Same validation pattern as existing `Check` validation
- Uses `ConfigError` type for consistent error handling

### Built-in Prompts (Phase 1)

Four example prompts defined in vibeguard.yaml:
- **init** - Guidance for initializing VibeGuard configuration
- **code-review** - System prompt for code review assistance
- **security-audit** - Security-focused code analysis
- **test-generator** - Comprehensive unit test generation

### Test Coverage (Phase 1)

**Test File:** internal/cli/prompt_test.go (663 lines)

Tests implemented:
1. `TestRunPrompt_ListAll` - List all prompts functionality
2. `TestRunPrompt_ListAllVerbose` - Verbose output with descriptions/tags
3. `TestRunPrompt_ReadSpecific` - Retrieve and display specific prompt
4. `TestRunPrompt_NotFound` - Error handling for non-existent prompts
5. `TestRunPrompt_NoPrompts` - Error when no prompts defined
6. `TestRunPrompt_JSONOutput` - JSON format output
7. `TestRunPrompt_ConfigNotFound` - Error when config file missing

All test paths covered including error conditions, edge cases, and format variations.

---

## Phase 2: Event Handlers & Built-in Init Prompt (Complete)

### Event Handler Data Model

**Event Handler Structure** (internal/config/events.go):

```go
type EventHandlers struct {
    Success []interface{} `yaml:"success,omitempty"`  // Can be array or string
    Failure []interface{} `yaml:"failure,omitempty"`
    Timeout interface{}   `yaml:"timeout,omitempty"`
}

type Check struct {
    // ... existing fields ...
    On EventHandlers `yaml:"on,omitempty"`
}
```

**Prompt Value Semantics:**
- **Array elements** → Treated as prompt ID references (must exist in prompts section)
- **String values** → Treated as inline content (no ID lookup)
- **Key distinction:** Type determines interpretation, not content

### Event Types & Precedence

**Three Event Types:**

1. **success** - Triggered when check passes (exit code 0, assertions true)
2. **failure** - Triggered when check fails (exit code ≠ 0, assertions false)
3. **timeout** - Triggered when check exceeds timeout duration

**Precedence Rules** (highest to lowest):
- Timeout (if triggered, failure/success are not)
- Failure (if triggered and no timeout, success is not)
- Success (only if no failure or timeout)

### Event Handler Validation

**Validation Rules** (internal/config/events.go):

For prompt ID references (array elements only):
- Must match existing prompt ID in config
- Same ID format rules as check IDs
- Invalid references → ConfigError with line number

For inline content (string values):
- No validation required
- Can be any text
- No ID matching needed

**Edge Cases Handled:**
- Empty event arrays (no prompts triggered)
- Cancelled checks (fail-fast) → no prompts triggered
- Skipped checks (dependency failed) → no prompts triggered
- Mixed arrays with IDs and inline content supported

### Orchestrator Integration

**Changes to** `internal/orchestrator/orchestrator.go`:

- Added event handler evaluation logic
- Prompt triggering on check completion
- Event precedence implementation
- Triggered prompt collection for each check result

**Output Integration:**
- Triggered prompts included in CheckResult
- Both human-readable and JSON formats supported
- Inline display of prompt content with results

### Output Formatting

**Human-Readable Format** (internal/output/formatter.go):

```
FAIL  vet (error)

  exit status 1
  Fix: Review the error message above
  Log: .vibeguard/log/vet.log

  Triggered Prompts (failure):
  [1] init:
      You are an expert in helping users...

  [2] (inline):
      Also remember to run gofmt...

  Advisory: blocks commit
```

**JSON Format** (internal/output/json.go):

```json
{
  "violations": [
    {
      "id": "vet",
      "severity": "error",
      "triggered_prompts": [
        {
          "event": "failure",
          "source": "init",
          "content": "You are an expert..."
        },
        {
          "event": "failure",
          "source": "inline",
          "content": "Also remember..."
        }
      ]
    }
  ]
}
```

### Built-in Init Prompt

**Implementation** (internal/cli/init_prompt.go):

- Init prompt content embedded as constant in binary
- Content mirrors `vibeguard init --assist` instructions
- Fallback logic: User-defined init prompt takes precedence

**Behavior:**
- Available without vibeguard.yaml
- Included in `vibeguard prompt --json` output
- Can be overridden by custom init prompt in config
- Access via `vibeguard prompt init`

### Test Coverage (Phase 2)

**Event Handler Tests** (internal/config/events_test.go - 447 lines):
- Event type validation
- Prompt ID reference validation
- Inline content handling
- Array vs string syntax parsing
- Error cases and edge conditions

**Orchestrator Tests** (internal/orchestrator/orchestrator_test.go):
- Event triggering on success/failure/timeout
- Event precedence enforcement
- Prompt content inclusion in results
- Integration with check execution flow

**Output Tests** (internal/output/formatter_test.go, internal/output/json_test.go):
- Human-readable formatting of triggered prompts
- JSON structure validation
- Multi-prompt display
- Error message formatting

---

## Architecture & Design Decisions

### 1. Configuration-First Approach
- Prompts stored in vibeguard.yaml alongside checks and vars
- Enables version control of prompts with code
- Integrated validation during config load
- Discoverable via standard config mechanisms

### 2. Dual Representation (ID vs Inline)
- Array syntax for prompt ID references
- String syntax for inline content
- Clear, unambiguous distinction at YAML level
- Supports both reusable prompts and custom messages

### 3. Backward Compatibility
- Prompts section entirely optional
- Existing configs load without modification
- Event handlers optional on all checks
- No breaking changes to existing features

### 4. Output Flexibility
- List mode for human discovery
- Verbose mode for additional context
- JSON for programmatic access
- Raw content for piping to other tools

### 5. CLI Integration
- Seamless piping to LLM command-line tools
- Machine-discoverable via `--json` output
- Consistent with Unix philosophy
- Works with agent automation workflows

---

## Configuration Examples

### Phase 1 Example (Basic Prompts)

```yaml
version: "1"

prompts:
  - id: init
    description: "Initialization guidance"
    content: |
      You are an expert in helping users set up VibeGuard...
    tags: [setup, initialization, guidance]

  - id: code-review
    description: "Code review assistance"
    content: |
      You are an expert code reviewer...
    tags: [review, quality, go]

checks:
  - id: vet
    run: go vet ./...
    severity: error
```

### Phase 2 Example (With Event Handlers)

```yaml
version: "1"

prompts:
  - id: fix-guidance
    description: "Guide for fixing issues"
    content: |
      Follow these steps to fix the issue...
    tags: [fix, resolution]

checks:
  - id: vet
    run: go vet ./...
    severity: error
    on:
      failure:
        - fix-guidance                    # Prompt ID reference
        - "Check the error message above" # Inline content
      success: "Great! Code passed vet"
      timeout: "Vet check timed out"

  - id: test
    run: go test ./...
    on:
      failure: [fix-guidance]
      timeout: "Tests taking too long"
```

---

## Feature Completeness Assessment

### Phase 1 Achievements ✅
- ✅ Prompt data structure and schema
- ✅ CLI command with list/verbose/json/specific modes
- ✅ Configuration loading and validation
- ✅ Error handling and reporting
- ✅ Comprehensive test coverage
- ✅ Backward compatibility
- ✅ Built-in example prompts

### Phase 2 Achievements ✅
- ✅ Event handler data model
- ✅ Event validation (ID references and inline content)
- ✅ Orchestrator integration
- ✅ Event precedence rules (timeout > failure > success)
- ✅ Human-readable output formatting
- ✅ JSON output with triggered prompts
- ✅ Built-in init prompt in binary
- ✅ Comprehensive test coverage
- ✅ Edge case handling (cancelled/skipped checks)

---

## Technical Details

### Code Organization

```
internal/cli/
  prompt.go              # Main command handler (181 lines)
  prompt_test.go         # Phase 1 tests (663 lines)
  init_prompt.go         # Built-in init prompt constant

internal/config/
  schema.go              # Prompt type definition
  config.go              # Loading and validation
  events.go              # Event handler types and validation (75 lines)
  events_test.go         # Event handler tests (447 lines)

internal/orchestrator/
  orchestrator.go        # Event triggering logic
  orchestrator_test.go   # Integration tests

internal/output/
  formatter.go           # Human-readable formatting
  formatter_test.go      # Format tests
  json.go                # JSON output with triggered prompts
  json_test.go           # JSON tests

vibeguard.yaml          # Configuration with prompts and event handlers
```

### Key Implementation Details

**Prompt ID Format:**
- Regex: `^[a-zA-Z_][a-zA-Z0-9_-]*$`
- Same pattern as check IDs
- Enforced at validation time

**Tag Format:**
- Regex: `^[a-z][a-z0-9-]*$`
- Lowercase only
- Optional field

**Error Context:**
- All validation errors include file path and line number
- Enables easy location of issues in config
- Uses ConfigError type for consistency

---

## Integration Points

### With Existing Features

**Configuration System:**
- Prompts loaded as part of standard `Config.Load()`
- Validated in `Config.Validate()`
- Same error handling patterns

**Check System:**
- Event handlers optional on all checks
- Works with existing check fields (grok, assert, requires)
- Respects fail-fast behavior
- Compatible with dependency ordering

**Output System:**
- Formatted alongside check results
- Both text and JSON output modes
- No impact on existing output when prompts not used

**CLI Structure:**
- Registered as subcommand of root command
- Follows established Cobra patterns
- Inherits standard flags (--help, --version)

---

## Usage Patterns Enabled

### Agent Discovery
```bash
vibeguard prompt --json | jq '.[].id'
```

### Direct LLM Integration
```bash
vibeguard prompt code-review | llm prompt "Review this code"
```

### Workflow Automation
```bash
# Agent chains commands using both prompts and check results
vibeguard prompt init | tee init_guidance.txt
vibeguard check run | grep -A 5 "Triggered Prompts"
```

### Human Exploration
```bash
vibeguard prompt -v              # Discover available prompts
vibeguard prompt security-audit | less  # Read full content
```

---

## Testing Methodology

### Unit Test Strategy
- Individual component testing (prompt loading, validation, formatting)
- Error path testing (missing files, invalid IDs, malformed YAML)
- Edge case coverage (empty arrays, null values, special characters)

### Integration Testing
- Full pipeline from config load to formatted output
- Event triggering during check execution
- Prompt precedence rules
- Output format validation

### Test Data
- Comprehensive example configs covering all features
- Edge cases (missing fields, special characters, long content)
- Error scenarios (invalid IDs, missing references)

---

## Backward Compatibility Analysis

✅ **Fully Backward Compatible:**
- Prompts section optional (not required in config)
- Event handlers optional on all checks
- Existing configs parse correctly
- No changes to existing check/var functionality
- Graceful degradation when features not used

---

## Derivation Notes

This specification was reverse-engineered from:
1. Complete test suite (showing intended behavior)
2. Implementation code (showing actual behavior)
3. Type definitions (showing data model)
4. Configuration examples (showing usage patterns)
5. Commit messages (showing implementation phases)

All features documented here have corresponding test coverage and are production-ready.

---

## Related Documentation

- **Phase 1 Spec:** docs/specs/prompt-feature-spec.md (sections 1-6, 8, success criteria)
- **Phase 2 Spec:** docs/specs/prompt-feature-spec.md (sections 7, 9, success criteria)
- **Implementation Logs:**
  - docs/log/2026-01-17_prompt-feature-implementation.md
  - docs/log/2026-01-19_p2-4-human-readable-output-formatting.md
  - docs/log/2026-01-19_p2-event-handler-completion.md

---

## Summary

The Prompt Feature represents a complete, production-ready implementation enabling VibeGuard to store, retrieve, and distribute guided prompts for AI agents and users. Phase 1 provides core CLI functionality with configuration-based storage. Phase 2 adds automatic triggering of prompts on check execution outcomes with flexible inline and ID-based references. Both phases include comprehensive testing, validation, error handling, and documentation. The implementation is fully backward compatible and ready for immediate use.

