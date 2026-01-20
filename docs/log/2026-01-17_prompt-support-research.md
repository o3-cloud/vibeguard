---
summary: Explored vibeguard codebase to understand architecture and identify integration points for adding prompt support
event_type: research
sources:
  - docs/adr/ADR-003-adopt-golang.md
  - docs/adr/ADR-005-adopt-vibeguard.md
  - internal/config/schema.go
  - internal/orchestrator/orchestrator.go
tags:
  - prompts
  - architecture
  - feature-planning
  - configuration
  - integration-points
---

# VibeGuard Prompt Support Research

## Overview

Analyzed the vibeguard codebase to understand its current architecture and identify suitable integration points for adding prompt support.

## What is VibeGuard

VibeGuard is a **lightweight, composable policy enforcement system** designed for seamless integration with CI/CD pipelines, agent loops, and Cloud Code workflows. Written in Go as a single-binary deployment tool, it:

- Orchestrates external tools via shell commands
- Evaluates assertions against tool output using grok patterns
- Enforces policies declaratively through YAML configuration
- Provides intelligent, actionable violations when checks fail
- Supports parallel execution with dependency management
- Integrates with LLMs for nuanced policy evaluation
- Follows "Silence is Success" principle (no output when all checks pass)

## Current Architecture

The system follows a layered architecture:

```
CLI Layer (Cobra)
    ↓
Configuration Layer (YAML parsing, validation, interpolation)
    ↓
Orchestration Layer (DAG execution with topological sort)
    ↓
Execution Layer (Shell commands with timeout handling)
    ↓
Analysis Layer (Grok pattern extraction + assertion evaluation)
    ↓
Output Layer (Text and JSON formatting)
```

### Current Policy Structure (Check-based)

Checks are defined in YAML with:
- **id**: Unique identifier
- **run**: Shell command to execute
- **grok**: Pattern extraction rules
- **assert**: Expression evaluation on extracted data
- **severity**: error or warning
- **suggestion**: Actionable help text
- **fix**: Suggested fix command
- **requires**: Dependency list
- **tags**: Categorization
- **timeout**: Execution timeout

Example:
```yaml
checks:
  - id: format-check
    run: gofmt -l .
    assert: "lines == 0"
    severity: error
    suggestion: "Run 'gofmt -w .' to fix formatting"
```

## Key Implementation Details

### Configuration Layer
- **Location**: `internal/config/schema.go` and `internal/config/config.go`
- Supports variable interpolation via `{{.variable_name}}`
- Validates config structure and detects dependency cycles
- Uses three-state DFS for cycle detection

### Orchestration Layer
- **Location**: `internal/orchestrator/orchestrator.go` and `internal/orchestrator/graph.go`
- Implements Kahn's algorithm for topological sort
- Creates execution levels for parallel processing
- Default: 4 parallel checks (configurable)
- Supports fail-fast behavior

### Analysis Layer
- **Location**: `internal/assert/` (lexer, parser, eval)
- Custom expression parser with AST
- Supports: arithmetic, comparison, logical operators
- Integrates grok pattern extraction (elastic/go-grok library)

## Implementation Plan: `vibeguard prompt` Subcommand

Based on detailed research into existing CLI patterns, the following design is recommended:

### Feature Specification

**Command**: `vibeguard prompt [prompt-id]`

**Behavior**:
- Without arguments: Lists all available prompts with their IDs
- With prompt ID: Outputs the raw prompt content (for piping to other tools)
- With `-v` flag: Shows full details including descriptions and tags

### YAML Configuration Structure

```yaml
version: "1"

prompts:
  - id: code-review
    description: "System prompt for code review assistance"
    content: |
      You are an expert Go code reviewer. When reviewing code:
      1. Check for idiomatic Go patterns
      2. Look for potential bugs and edge cases
      3. Suggest performance improvements
      4. Verify error handling

  - id: test-generator
    description: "Prompt for generating unit tests"
    content: |
      Generate comprehensive unit tests for the provided Go code.
      Use table-driven tests where appropriate.
      Cover edge cases and error conditions.
    tags: [testing, generation]
```

### Code Changes Required

#### 1. Config Schema Changes
**File**: `internal/config/schema.go`

Add new `Prompt` struct:
```go
type Prompt struct {
    ID          string   `yaml:"id"`
    Description string   `yaml:"description,omitempty"`
    Content     string   `yaml:"content"`
    Tags        []string `yaml:"tags,omitempty"`
}
```

Add `Prompts` field to `Config`:
```go
type Config struct {
    Version  string            `yaml:"version"`
    Vars     map[string]string `yaml:"vars"`
    Prompts  []Prompt          `yaml:"prompts,omitempty"`  // NEW
    Checks   []Check           `yaml:"checks"`
    yamlRoot interface{}       `yaml:"-"`
}
```

#### 2. Config Validation
**File**: `internal/config/config.go`

Add validation in `Validate()` method:
- Verify all prompt IDs are unique
- Validate prompt IDs follow naming convention (alphanumeric, hyphens, underscores)
- Verify each prompt has content defined

#### 3. CLI Command Implementation
**File**: `internal/cli/prompt.go` (NEW)

Command implementation following the pattern from `tags.go` and `check.go`:
```go
var promptCmd = &cobra.Command{
    Use:   "prompt [prompt-id]",
    Short: "Display prompt content",
    Long: `Display prompt configurations from vibeguard.yaml.

Without a prompt ID, lists all available prompts.
With a prompt ID, outputs the prompt content.

Examples:
  vibeguard prompt              List all prompts
  vibeguard prompt code-review  Output the code-review prompt content
  vibeguard prompt -v           List prompts with descriptions`,
    Args: cobra.MaximumNArgs(1),
    RunE: runPrompt,
}

func init() {
    rootCmd.AddCommand(promptCmd)
}
```

Implementation features:
- Load config using `config.Load(configFile)`
- Support optional positional argument for prompt ID
- Use `cobra.MaximumNArgs(1)` for optional argument
- Output raw prompt content when ID provided (for piping)
- List all prompts when no ID provided
- Support `-v`/`--verbose` flag for detailed output
- Support `--json` flag for JSON output

#### 4. Tests
**File**: `internal/cli/prompt_test.go` (NEW)

Test cases (following `tags_test.go` pattern):
- List all prompts
- Read specific prompt content
- Handle prompt not found error
- Handle empty prompts list
- Verbose output mode
- JSON output format

### Implementation Steps

1. **Add Prompt struct to config schema** (`internal/config/schema.go`)
   - Define `Prompt` struct with id, description, content, tags
   - Add `Prompts []Prompt` field to `Config`

2. **Enhance config validation** (`internal/config/config.go`)
   - Validate unique prompt IDs
   - Validate prompt ID naming convention
   - Ensure content exists

3. **Create prompt command** (`internal/cli/prompt.go`)
   - Define `promptCmd` with optional `[prompt-id]` argument
   - Register with root command
   - Implement listing (no argument case)
   - Implement reading specific prompt (with ID argument)
   - Support verbose and JSON output modes

4. **Add comprehensive tests** (`internal/cli/prompt_test.go`)
   - Follow existing test patterns from `tags_test.go`
   - Test list all prompts
   - Test read specific prompt
   - Test verbose/JSON output
   - Test error cases

5. **Add built-in prompts** to the system
   - Create `vibeguard prompt init` as a built-in prompt for initialization assistance
   - This replaces the `vibeguard init --assist` workflow
   - Move assist logic into prompt content structure
   - Prompts can be defined in code as defaults or loaded from config

### Built-in Prompts

The system should support built-in prompts that are available by default:

**`init` prompt**: Initialization and setup assistance
```yaml
prompts:
  - id: init
    description: "Guidance for initializing vibeguard configuration"
    content: |
      You are an expert in helping users set up vibeguard policy enforcement.

      Guide them through:
      1. Detecting their project type (Go, Node.js, Python, Rust, etc.)
      2. Recommending appropriate checks based on their project
      3. Explaining what each check does
      4. Creating a initial vibeguard.yaml configuration
      5. Testing the configuration

      Ask clarifying questions to understand their needs and preferences.
```

This allows users to access initialization guidance via:
```bash
vibeguard prompt init
```

### Design Rationale

- **Simple command structure** (`vibeguard prompt [prompt-id]`) follows established patterns from `tags` and `check` commands
- **Optional argument pattern** uses `cobra.MaximumNArgs(1)` (consistent with existing `check` command)
- **Raw output for reading** allows piping to other tools/LLMs
- **Global flags inherited** (`--config`, `--verbose`, `--json`) for consistency
- **Config schema extension** is backward compatible (optional `prompts:` field)
- **Flat command** is simpler than nested subcommands, easier to discover and use
- **Built-in prompts** replace `--assist` flags, making the system more discoverable and prompt-driven
- **Prompt-centric UI** aligns with LLM-driven workflows and agent integration

## Implementation Completed

### Files Created
- `internal/cli/prompt.go`: Main CLI command implementation (86 lines)
- `internal/cli/prompt_test.go`: Comprehensive test suite (265 lines)

### Files Modified
- `internal/config/schema.go`: Added Prompt struct and Prompts field to Config
- `internal/config/config.go`: Added validatePrompts() and FindPromptNodeLine() methods
- `vibeguard.yaml`: Added 4 built-in prompts (init, code-review, security-audit, test-generator)

### Key Implementation Details

#### Command Structure
```bash
vibeguard prompt              # List all prompts
vibeguard prompt init         # Output init prompt content
vibeguard prompt -v           # List with descriptions
vibeguard prompt --json       # JSON format
vibeguard prompt init | llm   # Pipe to external tools
```

#### Prompt Schema
```yaml
prompts:
  - id: prompt-id             # Required: unique identifier
    description: "..."        # Optional: human-readable description
    content: |                # Required: full prompt text
      Multi-line prompt
      content here...
    tags: [tag1, tag2]        # Optional: categorization
```

#### Validation Rules
- Prompt IDs must be unique
- IDs follow naming convention: start with letter/underscore, alphanumeric/hyphens/underscores
- Each prompt must have content defined
- Tags must be lowercase alphanumeric with hyphens
- Optional fields (description, tags) are omitted from YAML if empty

#### Built-in Prompts Added
1. **init**: Initialization and setup guidance (replaces `vibeguard init --assist`)
2. **code-review**: Go code review assistance
3. **security-audit**: Security vulnerability analysis
4. **test-generator**: Unit test generation

### Test Coverage

7 test cases implemented:
- TestRunPrompt_ListAll: Basic listing functionality
- TestRunPrompt_ListAllVerbose: Verbose output with descriptions
- TestRunPrompt_ReadSpecific: Read specific prompt content
- TestRunPrompt_NotFound: Error handling for missing prompts
- TestRunPrompt_NoPrompts: Error when no prompts defined
- TestRunPrompt_JSONOutput: JSON format validation
- TestRunPrompt_ConfigNotFound: Config loading errors

### Design Patterns Used

1. **Optional Positional Argument**: `cobra.MaximumNArgs(1)` pattern from existing `check` command
2. **Global Flag Inheritance**: Respects `--config`, `--verbose`, `--json` global flags
3. **Simple Output for Piping**: Raw text output when reading specific prompt (no formatting)
4. **Config Validation Pattern**: Follows existing check validation approach with DFS and line number tracking

### Backward Compatibility

- Prompts field is optional in YAML (`prompts,omitempty`)
- Existing configs without prompts continue to work
- `vibeguard init --assist` remains functional (can be deprecated in future)
- No breaking changes to existing CLI or config schema

### Usage Examples

**For Users**:
```bash
# Discover available prompts
vibeguard prompt

# Get specific guidance
vibeguard prompt init

# Use with LLM tools
vibeguard prompt code-review | llm "Review this:" < myfile.go
```

**For CI/CD Integration**:
```bash
# Export prompt for use in shell scripts
export REVIEW_PROMPT=$(vibeguard prompt code-review)

# Use in Makefile
review:
	@vibeguard prompt code-review | llm
```

**For Agent Integration**:
```python
import subprocess
import json

# Get all available prompts
result = subprocess.run(['vibeguard', 'prompt', '--json'],
                       capture_output=True, text=True)
prompts = json.loads(result.stdout)

# Read specific prompt
result = subprocess.run(['vibeguard', 'prompt', 'init'],
                       capture_output=True, text=True)
init_prompt = result.stdout
```

### Future Extensibility

The design supports:
1. **Additional Built-in Prompts**:
   - `performance-review`: Performance optimization
   - `dependency-audit`: Dependency security
   - `architecture-review`: System design
   - `migration-guide`: Upgrade instructions

2. **Custom User Prompts**: Users can add their own prompts to vibeguard.yaml

3. **Prompt Composition**: Future feature to combine prompts or generate checks from prompts

4. **LLM Integration Layer**: Potential integration with LLM APIs for dynamic prompt enhancement

### Related ADRs

- ADR-003: Go as primary implementation language
- ADR-005: Adopt VibeGuard for policy enforcement
- ADR-006: Git pre-commit hook integration
- Future ADR-009: Adopt Prompts for LLM-Driven Workflows
