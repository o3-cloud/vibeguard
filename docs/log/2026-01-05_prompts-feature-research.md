---
summary: Research on adding a prompts section to vibeguard.yaml for managing Claude Code prompts
event_type: research
sources:
  - internal/config/schema.go
  - internal/cli/assist/composer.go
  - internal/cli/templates/registry.go
  - vibeguard.yaml
tags:
  - feature-research
  - prompts
  - cli
  - configuration
  - claude-code
  - developer-experience
---

# Research: Adding `prompts:` Section to vibeguard.yaml

## Problem Statement

Currently, reusable prompts for Claude Code (`claude -p "..."`) are typically stored in separate PROMPT.txt files or embedded in scripts. This creates:
- Scattered prompt definitions across the project
- No standardized way to manage/version prompts alongside policy checks
- Manual overhead to invoke prompts with the correct context

## Proposed Feature

Add a `prompts:` section to vibeguard.yaml that allows defining named prompts, invokable via `vibeguard prompt <name>`.

### Example Configuration

```yaml
version: "1"

vars:
  project_name: vibeguard

prompts:
  review-pr:
    description: "Review a pull request for code quality and security"
    template: |
      Review this PR for:
      1. Code quality and readability
      2. Security vulnerabilities
      3. Test coverage gaps
      4. Adherence to {{.project_name}} conventions

      Focus on actionable feedback.

  fix-lint:
    description: "Fix linting errors in staged files"
    template: |
      Fix all linting errors in the staged files.
      Run `golangci-lint run` after fixes to verify.

  explain-check:
    description: "Explain a vibeguard check failure"
    args:
      - name: check_id
        required: true
    template: |
      Explain why the vibeguard check "{{.check_id}}" failed.
      Suggest how to fix it based on the check's configuration.
```

### CLI Usage

```bash
# List available prompts
vibeguard prompt --list

# Output a prompt (for piping to claude -p)
vibeguard prompt review-pr | claude -p

# Or direct execution (if we add --exec flag)
vibeguard prompt review-pr --exec

# With arguments
vibeguard prompt explain-check --arg check_id=coverage
```

## Current Architecture Analysis

### Configuration System (`internal/config/`)

The current `Config` struct in `schema.go`:
```go
type Config struct {
    Version string
    Vars    map[string]string
    Checks  []Check
}
```

Adding prompts would extend this to:
```go
type Config struct {
    Version string
    Vars    map[string]string
    Checks  []Check
    Prompts []Prompt  // New field
}

type Prompt struct {
    ID          string
    Description string
    Template    string
    Args        []PromptArg
}

type PromptArg struct {
    Name     string
    Required bool
    Default  string
}
```

### Variable Interpolation

The existing `interpolate.go` already supports `{{.VAR}}` syntax using Go's `text/template`. This can be reused directly for prompt templates, meaning prompts would automatically have access to:
- Global `vars:` from config
- Command-line `--arg` values
- Environment variables (if we add that)

### CLI Structure (`internal/cli/`)

New command would follow the existing pattern in `root.go`:
```go
func init() {
    rootCmd.AddCommand(promptCmd)
}

var promptCmd = &cobra.Command{
    Use:   "prompt [name]",
    Short: "Output or execute a named prompt",
    // ...
}
```

### Template Registry Pattern

The existing `templates/registry.go` provides a pattern we could adapt:
- Registry-based storage
- Name-based lookup
- Validation on registration

## Use Cases

### 1. Standardized Code Review Prompts
```bash
# In CI or pre-merge
vibeguard prompt review-pr | claude -p
```

### 2. Context-Aware Fix Prompts
```bash
# When a check fails, suggest fixes
vibeguard check coverage || vibeguard prompt fix-coverage | claude -p
```

### 3. Project Onboarding
```bash
# New developer orientation
vibeguard prompt explain-architecture | claude -p
```

### 4. Integration with Checks
```yaml
checks:
  - id: lint
    run: golangci-lint run
    suggestion: "Run `vibeguard prompt fix-lint | claude -p` for AI-assisted fixes"
```

### 5. Composable Prompts
```yaml
prompts:
  base-context:
    template: |
      Project: {{.project_name}}
      Language: Go

  review-security:
    includes: [base-context]
    template: |
      {{include "base-context"}}
      Focus on security vulnerabilities...
```

## Check-to-Prompt Integration (Key Pattern)

The most compelling integration is having checks suggest prompts when they fail. The prompt automatically inherits the check's extracted values, creating a context-aware fix workflow.

### Basic Integration

```yaml
prompts:
  fix-coverage:
    description: "Add tests to improve coverage"
    template: |
      The coverage check failed. Current coverage is {{.coverage}}%.
      Target is {{.coverage_target}}%.

      Analyze the uncovered code and add tests to reach the target.

checks:
  - id: coverage
    run: go test -cover ./...
    grok:
      - coverage: coverage:\s+%{NUMBER:coverage}%
    assert: "coverage >= 89"
    suggestion: "Coverage is {{.coverage}}%. Run: vibeguard prompt fix-coverage | claude -p"
```

### Workflow

```bash
$ vibeguard check coverage
✗ coverage: Coverage is 72%. Run: vibeguard prompt fix-coverage | claude -p

$ vibeguard prompt fix-coverage | claude -p
# Claude receives the full context with actual coverage value
```

### Shorthand: `suggest_prompt` Field

Add a dedicated field to reduce boilerplate:

```yaml
checks:
  - id: coverage
    run: go test -cover ./...
    grok:
      - coverage: coverage:\s+%{NUMBER:coverage}%
    assert: "coverage >= 89"
    suggestion: "Coverage is {{.coverage}}%"
    suggest_prompt: fix-coverage  # Auto-appends: "Run: vibeguard prompt fix-coverage | claude -p"
```

### Auto-Fix Mode: `--fix-with-prompt`

Take it further with automatic prompt execution on failure:

```bash
$ vibeguard check coverage --fix-with-prompt
# 1. Runs the check
# 2. On failure, finds the suggest_prompt
# 3. Automatically pipes to claude -p
# 4. Optionally re-runs check to verify fix
```

### Implementation: Passing Extracted Values to Prompts

The check's grok-extracted values need to flow to the prompt. Options:

**Option A: Environment Variables**
```bash
# vibeguard internally sets these before prompt interpolation
VIBEGUARD_coverage=72
VIBEGUARD_check_id=coverage
```

**Option B: Shared Context File**
```bash
# Write to .vibeguard/context.json after check failure
{"check_id": "coverage", "extracted": {"coverage": "72"}, "exit_code": 1}

# Prompt command reads this automatically
vibeguard prompt fix-coverage  # Merges context.json into template vars
```

**Option C: Explicit Passing (Most Flexible)**
```bash
# Check outputs extracted values
vibeguard check coverage --output-vars
# {"coverage": "72", "check_id": "coverage", "passed": false}

# Pipe to prompt
vibeguard check coverage --output-vars | vibeguard prompt fix-coverage --from-stdin
```

### Full Example: Lint Fix Workflow

```yaml
vars:
  project_name: vibeguard

prompts:
  fix-lint:
    description: "Fix linting errors"
    template: |
      The following linting errors were found in {{.project_name}}:

      ```
      {{.lint_output}}
      ```

      Fix all errors. Run `golangci-lint run` to verify.

checks:
  - id: lint
    run: golangci-lint run 2>&1 || true
    grok:
      - lint_output: (?s)(?P<lint_output>.*)  # Capture all output
    assert: "lint_output == ''"
    suggestion: "Linting failed"
    suggest_prompt: fix-lint
```

```bash
$ vibeguard check lint
✗ lint: Linting failed. Run: vibeguard prompt fix-lint | claude -p

$ vibeguard prompt fix-lint | claude -p
# Claude sees the actual lint errors and fixes them
```

### Benefits of This Pattern

1. **Context flows automatically** - Extracted values from failed checks populate prompts
2. **Single source of truth** - Check definitions and fix prompts live together
3. **Progressive automation** - Manual → suggested → auto-fix modes
4. **Discoverable** - `vibeguard prompt --list` shows available fix prompts
5. **Composable** - Same prompt can be suggested by multiple checks

## Alternative Approaches Considered

### A. Keep Prompts External (Status Quo)
- **Pros**: Simple, no changes needed
- **Cons**: Scattered, no versioning with config, no variable interpolation

### B. Separate prompts.yaml File
- **Pros**: Cleaner separation
- **Cons**: Another file to manage, misses integration with checks

### C. Prompts as Special Checks
```yaml
checks:
  - id: review-pr
    type: prompt  # New type
    template: "..."
```
- **Pros**: Reuses existing infrastructure
- **Cons**: Semantic mismatch (prompts aren't "checks")

## Implementation Considerations

### Schema Extension
- Add `Prompt` struct to `schema.go`
- Add validation for prompt IDs (same rules as check IDs)
- Ensure backward compatibility (prompts section optional)

### New CLI Command
- `vibeguard prompt --list` - List all prompts with descriptions
- `vibeguard prompt <name>` - Output interpolated prompt to stdout
- `vibeguard prompt <name> --exec` - Pipe directly to `claude -p`
- `vibeguard prompt <name> --arg key=value` - Pass template arguments

### Template Processing
- Reuse `InterpolateWithVars()` from `interpolate.go`
- Add support for prompt-specific args
- Consider adding `--stdin` to include piped input in template

### Output Modes
- Plain text (default, for piping)
- JSON (`--json`) for programmatic use
- Verbose (`-v`) to show prompt metadata

## Questions to Resolve

1. **Should prompts support `requires:` like checks?** Could chain prompts or require checks to pass first.

2. **Include functionality?** Allow prompts to compose/include other prompts?

3. **Environment variable access?** Should `{{.env.VAR}}` work in prompts?

4. **Stdin integration?** `git diff | vibeguard prompt review-diff` where `{{.stdin}}` captures input?

5. **Caching/history?** Track which prompts were run and their outputs?

## Next Steps

1. Create a design ADR if moving forward
2. Prototype the schema extension
3. Implement basic `vibeguard prompt` command
4. Add variable interpolation
5. Consider `--exec` integration with claude CLI

## Related

- ADR-005: VibeGuard policy enforcement in CI/CD
- `internal/cli/assist/` - Existing prompt generation for `--assist` flag
