---
summary: Implementation plan for adding prompts section to vibeguard.yaml
event_type: deep dive
sources:
  - docs/log/2026-01-05_prompts-feature-research.md
  - internal/config/schema.go
  - internal/cli/root.go
tags:
  - implementation-plan
  - prompts
  - cli
  - feature
---

# Implementation Plan: Prompts Feature

## Overview

Add a `prompts:` section to vibeguard.yaml that defines named prompts invokable via `vibeguard prompt <name>`, with tight integration to checks via `suggest_prompt` field.

## Phase 1: Core Schema & CLI (MVP)

**Goal:** Basic prompt definition and output

### 1.1 Schema Extension

**File:** `internal/config/schema.go`

```go
type Config struct {
    Version string
    Vars    map[string]string
    Checks  []Check
    Prompts map[string]Prompt  // New field
}

type Prompt struct {
    Description string   `yaml:"description"`
    Template    string   `yaml:"template"`
}
```

**Decision needed:** Use `map[string]Prompt` (keyed by name) or `[]Prompt` with ID field like checks?

- Map is cleaner for lookup: `config.Prompts["fix-coverage"]`
- Array matches existing checks pattern and preserves order

**Recommendation:** Map—prompts don't need ordering or dependencies like checks.

### 1.2 Config Validation

**File:** `internal/config/validate.go`

- Prompt names must be valid identifiers (alphanumeric + hyphen + underscore)
- Template must not be empty
- Description is optional but recommended

### 1.3 CLI Command

**File:** `internal/cli/prompt.go` (new)

```go
var promptCmd = &cobra.Command{
    Use:   "prompt [name]",
    Short: "Output a named prompt for use with claude -p",
}
```

Subcommands/flags:
- `vibeguard prompt --list` - List all prompts with descriptions
- `vibeguard prompt <name>` - Output interpolated prompt to stdout
- `vibeguard prompt <name> --json` - Output as JSON with metadata

### 1.4 Interpolation

Reuse existing `InterpolateWithVars()` from `internal/config/interpolate.go`.

Prompts get access to:
- All `vars:` from config
- Command-line `--set key=value` overrides

**Estimated scope:** ~200-300 lines of new code

---

## Phase 2: Check Integration

**Goal:** Checks can suggest prompts, values flow through

### 2.1 Add `suggest_prompt` Field to Check

**File:** `internal/config/schema.go`

```go
type Check struct {
    // ... existing fields
    SuggestPrompt string `yaml:"suggest_prompt"` // New field
}
```

### 2.2 Validation

- If `suggest_prompt` is set, verify the prompt exists in config
- Warn if suggestion text manually references prompts (redundant)

### 2.3 Output Formatting

**File:** `internal/output/formatter.go`

When a check fails and has `suggest_prompt`:
```
✗ coverage: Coverage is 72%
  Fix: vibeguard prompt fix-coverage | claude -p
```

### 2.4 Context Passing (Key Decision)

How do check's extracted values reach the prompt?

**Option A: Context File (Recommended)**
```
.vibeguard/context.json
{"check_id": "coverage", "extracted": {"coverage": "72"}, "vars": {...}}
```

- `vibeguard check` writes context on failure
- `vibeguard prompt` reads and merges into template vars
- Clean separation, debuggable (can inspect the file)

**Option B: Environment Variables**
```bash
VIBEGUARD_EXTRACTED_coverage=72 vibeguard prompt fix-coverage
```

- More Unix-y
- Harder to debug, clutters environment

**Option C: Explicit Piping**
```bash
vibeguard check coverage --json | vibeguard prompt fix-coverage --from-check
```

- Most explicit
- More typing for users

**Recommendation:** Option A (context file) with Option C available for scripting.

---

## Phase 3: Advanced Features

### 3.1 Prompt Arguments

```yaml
prompts:
  explain-check:
    args:
      - name: check_id
        required: true
    template: |
      Explain the {{.check_id}} check failure.
```

CLI: `vibeguard prompt explain-check --set check_id=coverage`

### 3.2 Stdin Capture

```yaml
prompts:
  review-diff:
    template: |
      Review this diff:
      ```
      {{.stdin}}
      ```
```

CLI: `git diff | vibeguard prompt review-diff`

### 3.3 Auto-Fix Mode

```bash
vibeguard check coverage --fix-with-prompt
```

1. Run check
2. On failure, find `suggest_prompt`
3. Execute `vibeguard prompt <name> | claude -p`
4. Optionally re-run check

**Decision needed:** Is auto-invoking claude too magical? Should it just print the command?

### 3.4 Prompt Composition (Optional)

```yaml
prompts:
  base-context:
    template: "Project: {{.project_name}}, Language: Go"

  fix-coverage:
    includes: [base-context]
    template: |
      {{.base_context}}
      Fix coverage...
```

**Decision needed:** Is this worth the complexity? Could just use vars for shared text.

---

## Implementation Order

| Phase | Scope | Depends On |
|-------|-------|------------|
| 1.1 | Schema extension | - |
| 1.2 | Validation | 1.1 |
| 1.3 | CLI command | 1.1, 1.2 |
| 1.4 | Interpolation | 1.3 |
| 2.1 | suggest_prompt field | 1.* |
| 2.2 | Validation | 2.1 |
| 2.3 | Output formatting | 2.1 |
| 2.4 | Context passing | 2.3 |
| 3.* | Advanced features | 2.* |

## Open Decisions Summary

1. **Map vs Array for prompts?** → Recommend map
2. **Context passing mechanism?** → Recommend context file + explicit piping
3. **Auto-fix behavior?** → Print command vs execute?
4. **Prompt composition?** → Worth the complexity?
5. **Proceed with implementation?** → Start with Phase 1 MVP?

## Files to Create/Modify

**New files:**
- `internal/cli/prompt.go` - CLI command
- `internal/prompt/prompt.go` - Prompt processing logic (optional, could be in cli)

**Modified files:**
- `internal/config/schema.go` - Add Prompt type
- `internal/config/validate.go` - Prompt validation
- `internal/config/interpolate.go` - Maybe extend for stdin
- `internal/output/formatter.go` - suggest_prompt output
- `internal/orchestrator/orchestrator.go` - Write context file on failure

## Success Criteria

Phase 1 complete when:
```bash
$ cat vibeguard.yaml
prompts:
  hello:
    template: "Hello {{.name}}"

$ vibeguard prompt hello --set name=World
Hello World

$ vibeguard prompt --list
hello    (no description)
```

Phase 2 complete when:
```bash
$ vibeguard check coverage
✗ coverage: Coverage is 72%
  Fix: vibeguard prompt fix-coverage | claude -p

$ vibeguard prompt fix-coverage
# Template populated with coverage=72 from last check
```
