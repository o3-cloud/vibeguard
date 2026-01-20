---
title: VibeGuard Prompt Feature - Implementation Specification
date: 2026-01-19
status: implemented
version: 1.0
author: Claude Code (Reverse Specification)
---

**Derived from:** `internal/cli/prompt.go`, `internal/config/schema.go`, `vibeguard.yaml`, and comprehensive test suite

## Overview

The Prompt Feature enables VibeGuard to store, retrieve, and distribute guided prompts for AI agents and users. Prompts are stored in the configuration file (`vibeguard.yaml`) and accessed via a dedicated CLI command, allowing seamless integration with LLM tools and agent workflows.

## Core Capabilities

### 1. Prompt Storage in Configuration

Prompts are stored in the `vibeguard.yaml` configuration file under a top-level `prompts` section:

```yaml
version: "1"

prompts:
  - id: init
    description: "Guidance for initializing vibeguard configuration"
    content: |
      You are an expert in helping users set up VibeGuard...
    tags: [setup, initialization, guidance]

  - id: code-review
    description: "System prompt for code review assistance"
    content: |
      You are an expert code reviewer...
    tags: [review, quality, go]
```

### 2. Prompt Data Structure

**Prompt Type Definition** (`internal/config/schema.go`):

```go
type Prompt struct {
    ID          string   `yaml:"id"`
    Description string   `yaml:"description,omitempty"`
    Content     string   `yaml:"content"`
    Tags        []string `yaml:"tags,omitempty"`
}
```

**Fields:**
- `ID` (required) - Unique identifier for the prompt (alphanumeric with hyphens/underscores)
- `Description` (optional) - Human-readable description of the prompt's purpose
- `Content` (required) - The full prompt text (supports multi-line YAML literals)
- `Tags` (optional) - Categorical tags for organizing and filtering prompts

### 3. CLI Command: `vibeguard prompt`

A new top-level command provides access to prompts with multiple usage patterns:

#### Command Signature

```bash
vibeguard prompt [prompt-id]
```

#### Behaviors

**Without Prompt ID** - List all available prompts:
```bash
vibeguard prompt
```

Output:
```
Prompts (4):

  init
  code-review
  security-audit
  test-generator
```

**With `--verbose` / `-v` flag** - List prompts with descriptions and tags:
```bash
vibeguard prompt -v
```

Output:
```
Prompts (4):

  init
    Description: Guidance for initializing vibeguard configuration
    Tags:        setup, initialization, guidance

  code-review
    Description: System prompt for code review assistance
    Tags:        review, quality, go
```

**With `--json` flag** - List prompts in JSON format (machine-readable):
```bash
vibeguard prompt --json
```

Output:
```json
[
  {
    "id": "init",
    "description": "Guidance for initializing vibeguard configuration",
    "tags": ["setup", "initialization", "guidance"]
  },
  {
    "id": "code-review",
    "description": "System prompt for code review assistance",
    "tags": ["review", "quality", "go"]
  }
]
```

**With Specific Prompt ID** - Output raw prompt content (optimized for piping):
```bash
vibeguard prompt init
```

Output:
```
You are an expert in helping users set up VibeGuard.

Guide them through:
1. Detecting their project type (Go, Node.js, Python, Rust, etc.)
2. Recommending appropriate checks based on their project
...
```

#### Integration with LLM Tools

The command is designed for piping to LLM command-line tools:

```bash
vibeguard prompt code-review | llm prompt "Review this Go code:\n\n$(cat main.go)"
```

This allows agents and users to apply stored system prompts directly to LLM services.

### 4. Error Handling

The feature includes comprehensive error handling:

**No Prompts Defined:**
```bash
$ vibeguard prompt
Error: no prompts defined in configuration
```

**Prompt Not Found:**
```bash
$ vibeguard prompt nonexistent
Error: prompt not found: nonexistent
```

**Configuration File Missing:**
```bash
$ vibeguard prompt
Error: no config file found (tried: [vibeguard.yaml vibeguard.yml .vibeguard.yaml .vibeguard.yml])
```

### 5. Configuration System Enhancements

The config system was enhanced to support prompts:

**Enhanced Config Type** (`internal/config/schema.go`):

```go
type Config struct {
    Version string            `yaml:"version"`
    Vars    map[string]string `yaml:"vars"`
    Prompts []Prompt          `yaml:"prompts,omitempty"`
    Checks  []Check           `yaml:"checks"`
    yamlRoot interface{}      `yaml:"-"`
}
```

**Validation Features** (`internal/config/config.go`):
- Validates all prompts have unique IDs
- Validates all prompt IDs follow the same format rules as check IDs (alphanumeric, underscores, hyphens)
- Validates all prompt IDs are formatted correctly (must start with letter or underscore)
- Validates all prompt tags follow lowercase alphanumeric format with hyphens
- Provides line number context for validation errors

### 6. Built-in Prompts

**Phase 1:** The project includes four predefined prompts in `vibeguard.yaml`:
- `init` - Guidance for initializing VibeGuard configuration
- `code-review` - System prompt for code review assistance
- `security-audit` - Security-focused code analysis
- `test-generator` - Comprehensive unit test generation

**Phase 2:** Single built-in `init` prompt embedded in VibeGuard binary
- Provides the same instructions as `vibeguard init --assist` flag
- Always available without requiring `vibeguard.yaml` configuration
- Can be used via `vibeguard prompt init` in any project
- Purpose: Help users understand and initialize VibeGuard for their project

## Implementation Architecture

### File Organization

```
internal/cli/
├── prompt.go          # Main prompt command implementation (106 lines)
├── prompt_test.go     # Comprehensive test suite (310 lines)
├── init.go            # Init command (uses same config system)
└── ...

internal/config/
├── config.go          # Config loading and validation (511 lines)
├── schema.go          # Type definitions including Prompt (66 lines)
└── ...

vibeguard.yaml        # Configuration with prompts section (150 lines)
```

### Execution Flow

```
User/Agent
    ↓
vibeguard prompt [prompt-id]
    ↓
runPrompt() loads config
    ↓
If no prompt ID:
  ├─ listPrompts() with flags (--verbose, --json)
  └─ Output formatted list or JSON
↓
If prompt ID provided:
  ├─ Search prompts for matching ID
  ├─ Found: Output raw content to stdout
  └─ Not found: Return error
```

### Integration Points

**Configuration System:**
- Prompts loaded alongside checks during `config.Load()`
- Validation integrated into `Config.Validate()`
- Line number tracking for error context

**CLI Command Structure:**
- Registered as subcommand of root command
- Follows Cobra command pattern
- Supports standard flags (--help, --version, etc.)

**Error Handling:**
- ConfigError types used for validation issues
- Errors propagate up to CLI with proper formatting
- Line numbers provided for debugging invalid configs

## Test Coverage

### Test Suite (`internal/cli/prompt_test.go`)

**7 Comprehensive Tests:**

1. **`TestRunPrompt_ListAll`**
   - Validates listing all prompts works
   - Confirms all prompt IDs appear in output

2. **`TestRunPrompt_ListAllVerbose`**
   - Validates verbose output includes descriptions
   - Validates tags are displayed correctly

3. **`TestRunPrompt_ReadSpecific`**
   - Validates reading specific prompt by ID
   - Confirms full prompt content is output
   - Tests multi-line prompt content

4. **`TestRunPrompt_NotFound`**
   - Validates error for non-existent prompt
   - Confirms proper error message

5. **`TestRunPrompt_NoPrompts`**
   - Validates error when no prompts are defined
   - Ensures graceful handling of missing prompts section

6. **`TestRunPrompt_JSONOutput`**
   - Validates JSON output format
   - Confirms valid JSON structure
   - Tests with actual JSON flag

7. **`TestRunPrompt_ConfigNotFound`**
   - Validates error when config file is missing
   - Tests proper error propagation

**Coverage:** All code paths tested including error conditions, multi-format output, and edge cases.

## Design Principles

### 1. **Configuration-First Design**
Prompts live in `vibeguard.yaml`, making them:
- Versioned with code (git history)
- Part of project configuration
- Easy for agents to discover and use
- Naturally discoverable via config validation

### 2. **Composable and Modular**
- Each prompt is independent
- Can be used individually or in workflows
- No dependencies between prompts
- Easy to add, remove, or modify prompts

### 3. **CLI-Optimized Output**
- List output is human-readable by default
- Verbose mode adds context without clutter
- JSON output for programmatic access
- Raw content output for piping to other tools

### 4. **Agent-Friendly Integration**
- Prompts are machine-discoverable via `--json`
- Content output is suitable for LLM pipes
- Consistent naming and tagging system
- Clear error messages guide troubleshooting

### 5. **Graceful Degradation**
- Works with or without prompts defined
- Clear errors when prompts are missing
- Doesn't require all optional fields
- Backward compatible with existing configs

## Check Event Handlers (Phase 2 Feature)

### Overview

Prompts can be attached to check execution outcomes, allowing guidance to surface automatically when checks succeed, fail, or timeout. This enables workflows where users receive contextual prompts based on check results.

### Syntax

Event handlers are specified using the `on:` key under a check:

```yaml
checks:
  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    on:
      success: [code-review]                    # Array of prompt IDs
      failure: [init, security-audit]           # Array of prompt IDs
      timeout: "Check timed out. Try again."    # Inline string (not ID)
```

### Event Types

**success** - Triggered when check passes
- Exit code is 0
- Any assertions evaluate to true
- Used for guidance or contextual suggestions on passing checks

**failure** - Triggered when check fails
- Exit code is non-zero, OR
- Assertion evaluates to false
- Applies to both `severity: error` and `severity: warning`
- Provides remediation guidance or next steps

**timeout** - Triggered when check exceeds timeout
- Takes precedence over failure
- Separate event for timeout-specific guidance
- Distinct from failure (execution interrupted, not completed)

### Prompt Value Types

Each event accepts prompt definitions with a clear distinction:

**Array syntax = Prompt ID references:**
```yaml
failure: [init, security-audit]      # Both treated as prompt IDs
failure: [init]                      # Single ID in array
failure:
  - init                             # Single ID in array
  - security-audit                   # Multiple IDs
```

**String syntax = Inline content:**
```yaml
timeout: "This check timed out. Consider increasing the timeout."
failure: init                        # NOT an array -> treated as inline content
failure: "Check logs for details"    # String literal -> inline content
```

**Key Rule:** Only array elements are treated as prompt ID references. A bare string is always inline content, even if it looks like a prompt ID.

### Event Evaluation & Precedence

Event precedence (highest to lowest):
1. **Timeout** - If execution times out, only timeout prompts trigger (not failure)
2. **Failure** - If check fails (and didn't timeout), failure prompts trigger
3. **Success** - If check passes, success prompts trigger

### Output Display

When prompts are triggered, they display inline with check results:

**Human-Readable Format:**
```
FAIL  vet (error)

  exit status 1
  Fix: Review the error message above
  Log: .vibeguard/log/vet.log

  Triggered Prompts (failure):
  [1] init:
      You are an expert in helping users set up VibeGuard...

  [2] (inline):
      Also remember to run gofmt before committing

  Advisory: blocks commit
```

**JSON Format:**
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

### Validation Rules

**Syntax-Based Validation:**
- **Array elements** (e.g., `[init, audit]`) - Validated as prompt ID references. Must exist in `prompts:` section. Missing IDs cause validation error with line number.
- **String values** (e.g., `init` or `"inline text"`) - Always treated as inline content. No ID validation required.

**Prompt ID Requirements (for array elements only):**
- Must match alphanumeric format with underscores/hyphens
- Must be defined in `prompts:` section
- Invalid references cause ConfigError with line number

**Event Names:**
- Only `success`, `failure`, `timeout` are valid
- Other event names are ignored (graceful degradation)

**Edge Cases:**
- Empty event handler arrays are valid (no prompts triggered)
- Cancelled checks (fail-fast) do not trigger prompts
- Skipped checks (dependency failed) do not trigger prompts

### Configuration Examples

**Example 1: Success guidance**
```yaml
checks:
  - id: test
    run: go test ./...
    on:
      success: [code-review]  # Show review suggestions on passing tests
```

**Example 2: Multi-prompt failure handling**
```yaml
checks:
  - id: security-scan
    run: gosec ./...
    severity: error
    on:
      failure:
        - security-audit      # Check for security issues
        - init                # Get setup guidance
```

**Example 3: Timeout guidance with inline content**
```yaml
checks:
  - id: integration-tests
    run: go test -v ./tests/integration/...
    timeout: 30s
    on:
      timeout: "Integration tests timed out. Check for hanging goroutines or network issues."
      failure: [test-generator]  # Generate test cases if tests fail
```

**Example 4: Syntax distinction (important)**
```yaml
checks:
  - id: example1
    run: some command
    on:
      failure: [init]                # ARRAY: "init" is a prompt ID reference

  - id: example2
    run: some command
    on:
      failure: init                  # STRING: "init" is inline content (even if it's a valid ID!)

  - id: example3
    run: some command
    on:
      failure:
        - init                       # ARRAY: "init" is a prompt ID reference
        - "Run fix command here"     # ARRAY: Inline content
```

**Example 5: Mixed inline and ID references**
```yaml
checks:
  - id: lint
    run: golangci-lint run ./...
    on:
      failure:
        - init                                  # ARRAY: Use init prompt ID
        - "Run 'golangci-lint run --fix' to auto-fix issues"  # ARRAY: Inline help
```

### Integration with Existing Features

- Works with all check configurations (grok, assert, requires, tags)
- Respects fail-fast behavior (timeouts don't trigger events on cancelled checks)
- Compatible with dependency ordering (`requires`)
- Supports all prompt features (IDs, descriptions, tags, multi-line content)

## Usage Patterns

### Pattern 1: Agent Discovery
```bash
# Agent discovers available prompts
vibeguard prompt --json | jq '.[].id'

# Output: init, code-review, security-audit, test-generator
```

### Pattern 2: Direct LLM Integration
```bash
# Agent uses prompt with LLM tool
CODE=$(cat myfile.go)
vibeguard prompt code-review | llm prompt "Review this:\n\n$CODE"
```

### Pattern 3: Workflow Orchestration
```bash
# Agent chains commands for complex workflows
PROMPT=$(vibeguard prompt init)
ANALYSIS=$(vibeguard init --assist)
# ... feed to LLM with both prompt and analysis
```

### Pattern 4: Human Exploration
```bash
# User discovers available prompts
vibeguard prompt -v

# User views specific prompt
vibeguard prompt security-audit | less
```

## Backward Compatibility

✅ **Fully Backward Compatible:**
- Prompts section is optional (`omitempty` in schema)
- Existing configs without prompts still load correctly
- No breaking changes to checks or other features
- Graceful error handling when prompts not defined

## Validation Rules

### Prompt ID Validation
- Must start with letter or underscore
- Can contain alphanumeric characters, underscores, hyphens
- Must be unique across all prompts in config
- Same rules as check IDs (enforced via `validCheckID` regex)

### Prompt Tag Validation
- Must start with lowercase letter
- Can contain lowercase alphanumeric and hyphens
- Must follow lowercase convention (enforced via `validTag` regex)

### Prompt Content Validation
- Content field is required (empty content not allowed)
- Supports multi-line content via YAML literals (`|`, `|-`, `>`)
- No length limits

### Configuration Context
- Prompts validated as part of config load process
- Validation failures include line numbers for debugging
- Errors use ConfigError type with file/line context

## Future Extensions

### Phase 2 Potential Enhancements
- Prompt versioning and updates
- Prompt search/filtering by tag
- Prompt composition (combining multiple prompts)
- Template variables in prompts (like checks)
- Prompt dependencies/ordering
- Prompt metadata (author, created date, updated date)

### Phase 3 Potential Features
- Remote prompt registry integration
- Custom prompt commands
- Prompt hot-reload without restart
- Prompt performance metrics
- Interactive prompt editor

## Success Criteria

✅ **Phase 1**
- [ ] Prompts stored in YAML configuration
- [ ] CLI command for listing and retrieving prompts
- [ ] Support for verbose and JSON output formats
- [ ] Raw content output for piping
- [ ] Comprehensive test coverage (7 tests, all passing)
- [ ] Proper error handling and validation
- [ ] Backward compatible with existing configs
- [ ] Works with AI agent integrations
- [ ] Four built-in example prompts included

⏳ **Phase 2:**
- [ ] `on:` syntax for attaching prompts to check events
- [ ] Support for success, failure, timeout events
- [ ] Both prompt ID references and inline content strings
- [ ] Full prompt content displayed inline with results
- [ ] Validation of prompt ID references with line numbers
- [ ] JSON output includes triggered prompts
- [ ] Comprehensive test coverage for event handlers
- [ ] Event precedence rules (timeout > failure > success)
- [ ] Built-in `init` prompt embedded in binary (no config needed)
- [ ] `vibeguard prompt init` available without `vibeguard.yaml`
- [ ] Prompt content matches `vibeguard init --assist` instructions

## Exit Codes

- **0** - Success (prompt displayed, list shown, or JSON generated)
- **1** - Error (config issue, prompt not found, validation failure)

## Configuration Example

**Phase 1** - Complete configuration with all features:

```yaml
version: "1"

vars:
  go_packages: "./..."

prompts:
  - id: init
    description: "Guidance for initializing vibeguard configuration"
    content: |
      You are an expert in helping users set up VibeGuard.
      Guide them through project detection and configuration creation.
    tags: [setup, initialization, guidance]

  - id: code-review
    description: "System prompt for code review assistance"
    content: |
      You are an expert code reviewer with deep knowledge of Go.
      Check for idiomatic patterns, security issues, and test quality.
    tags: [review, quality, go]

  - id: security-audit
    description: "Security-focused code analysis prompt"
    content: |
      You are a security auditor. Find vulnerabilities including:
      - Command injection and path traversal
      - SQL injection and unsafe deserialization
      - Access control and cryptographic weaknesses
    tags: [security, audit, vulnerability]

  - id: test-generator
    description: "Prompt for generating comprehensive unit tests"
    content: |
      Generate comprehensive unit tests for provided Go code.
      Use table-driven tests, cover edge cases and errors.
    tags: [testing, generation, quality]

checks:
  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    timeout: 5s

  - id: test
    run: go test ./...
    severity: error
```

**Phase 2** - With event handlers and custom prompts:

```yaml
version: "1"

vars:
  go_packages: "./..."

prompts:
  - id: security-audit
    description: "Security-focused code analysis"
    content: |
      Check for security vulnerabilities...
    tags: [security, audit]

  - id: test-generator
    description: "Generate comprehensive unit tests"
    content: |
      Write tests covering edge cases...
    tags: [testing, generation]

checks:
  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    on:
      success:
        - "Great! Your code passed vet checks."
      failure: [security-audit]

  - id: test
    run: go test ./...
    severity: error
    on:
      failure: [test-generator]
      timeout: "Tests timed out. Try running smaller test suites."
```

**Built-in `init` prompt** (Phase 2):
- Always available: `vibeguard prompt init`
- Works without `vibeguard.yaml` file
- Provides the same guidance as `vibeguard init --assist`
- No need to define in configuration

## Built-in Init Prompt (Phase 2 Feature)

### Overview

The `init` prompt is embedded in the VibeGuard binary and provides the same guidance as the `vibeguard init --assist` command. This makes it always available without requiring a `vibeguard.yaml` configuration file.

### Usage

```bash
# Works without any config file
vibeguard prompt init

# Display with less or pipe to LLM
vibeguard prompt init | less
vibeguard prompt init | llm prompt "Guide me through setup"

# JSON output (includes built-in prompt)
vibeguard prompt --json
# Output shows "init" prompt even without config file
```

### Implementation Details

**Storage:**
- Prompt content stored as a constant in `internal/cli/init_prompt.go`
- Content mirrors `vibeguard init --assist` instructions

**Fallback Logic:**
- When `vibeguard prompt init` is requested:
  1. First check if "init" exists in `vibeguard.yaml`
  2. If not found, use built-in prompt
  3. User-defined prompts take precedence over built-in

**Behavior:**
- Available even when no `vibeguard.yaml` exists
- Can be overridden by defining custom `init` prompt in config
- Included in `vibeguard prompt --json` output
- Displayed with source indicator "(built-in)" in verbose mode

### Configuration Override

```yaml
prompts:
  - id: init
    description: "Custom initialization guide"
    content: |
      This replaces the built-in init prompt...
```

## References

### Phase 1 Implementation Files
- CLI Command: `internal/cli/prompt.go` (106 lines)
- Test Suite: `internal/cli/prompt_test.go` (310 lines)
- Config Schema: `internal/config/schema.go` (66 lines)
- Config Loading: `internal/config/config.go` (511 lines)
- Main Config: `vibeguard.yaml` (150 lines)

### Phase 2 Implementation Files (Event Handlers & Built-in Init Prompt - Planned)

**Event Handlers:**
- Event Handler Types: `internal/config/events.go` (~100 lines)
- Event Handler Tests: `internal/config/events_test.go` (~150 lines)
- Orchestrator: `internal/orchestrator/orchestrator.go` (modifications to lines 469-501, 698-725)
- Orchestrator Tests: `internal/orchestrator/orchestrator_test.go` (new event handler tests)
- Output Formatter: `internal/output/formatter.go` (new formatTriggeredPrompts method)
- JSON Formatter: `internal/output/json.go` (JSONTriggeredPrompt type)

**Built-in Init Prompt:**
- Prompt Definition: `internal/cli/init_prompt.go` (embedded prompt constant)
- Prompt Command: `internal/cli/prompt.go` (fallback to built-in if not in config)
- Tests: `internal/cli/prompt_test.go` (new test for built-in prompt)

### Related Features
- Config validation system (supports prompts)
- Check system (shares validation patterns)
- Assist system (references prompts)
- Template system (documented separately)

## Summary

The Prompt Feature provides a clean, CLI-integrated system for storing and accessing guided prompts. Prompts are configuration-driven, fully validated, comprehensively tested, and designed for seamless agent integration. The feature is production-ready and ships with four practical example prompts covering initialization, code review, security audit, and test generation use cases.
