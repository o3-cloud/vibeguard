---
summary: Design decision to replace vibeguard init --assist with built-in vibeguard prompt init for prompt-centric initialization workflow
event_type: code
sources:
  - internal/cli/init.go
  - internal/cli/assist/composer.go
  - docs/adr/ADR-005-adopt-vibeguard.md
tags:
  - prompts
  - initialization
  - user-experience
  - design-decision
  - llm-integration
  - cli-redesign
---

# Prompt-Centric Initialization Design

## Design Decision

Replace the `vibeguard init --assist` flag-based approach with a built-in prompt system:

**Old approach**:
```bash
vibeguard init --assist
```

**New approach**:
```bash
vibeguard prompt init
```

## Rationale

### Problem with Current Approach
- Initialization assistance is hidden behind a flag (`--assist`)
- Mixing initialization logic with flag-based assistance is confusing
- Not discoverable - users must know the flag exists
- Doesn't align with prompt-driven LLM workflows

### Solution: Prompt-Centric UI
- Make prompts first-class citizens in vibeguard
- Built-in prompts provide guidance and assistance
- Users naturally discover help via `vibeguard prompt`
- Extensible foundation for additional prompts

## Built-in Prompts

### 1. `init` Prompt (Replaces `vibeguard init --assist`)
Initialization and setup guidance
```
Guide users through:
1. Detecting project type (Go, Node.js, Python, Rust, etc.)
2. Recommending checks based on project
3. Explaining what each check does
4. Creating initial vibeguard.yaml
5. Testing the configuration
```

### 2. Future Built-in Prompts
- `generate-check` - Help creating new checks
- `security-audit` - Security-focused analysis
- `performance-review` - Performance optimization guidance
- `test-strategy` - Testing approach guidance

## Implementation Approach

### Config Schema
Prompts defined in YAML with:
- `id`: Unique identifier (e.g., "init")
- `description`: Human-readable description
- `content`: Full prompt text for LLM consumption
- `tags`: Optional categorization
- `builtin`: Flag indicating it's a built-in prompt

### CLI Command
```bash
vibeguard prompt [prompt-id]
```

Command behavior:
- Lists all prompts (built-in + config)
- Outputs specific prompt content
- Supports `-v` for verbose (descriptions)
- Supports `--json` for structured output
- Raw output suitable for piping to LLMs

### Built-in Prompt Registry
- Prompts compiled into binary as defaults
- Can be overridden in vibeguard.yaml
- Available even without configuration file
- Ensures initialization guidance always available

## Benefits

1. **Discoverability**: Users find help by exploring `vibeguard prompt`
2. **Consistency**: Prompt-centric workflow aligns with LLM integration
3. **Extensibility**: Foundation for additional guidance prompts
4. **Simplicity**: Clearer mental model (prompts are for guidance/assistance)
5. **Backward Compatible**: Existing `vibeguard init` remains, `--assist` deprecated
6. **Agent-Friendly**: Structured prompts for automation workflows

## Migration Path

1. Implement `vibeguard prompt` command with support for arbitrary prompts
2. Add built-in `init` prompt to registry
3. Deprecate `vibeguard init --assist` in favor of `vibeguard prompt init`
4. Guide users to new approach in documentation
5. Eventually remove `--assist` flag in major version

## Technical Implementation

### Files to Create/Modify
- `internal/config/schema.go`: Add Prompt struct
- `internal/cli/prompt.go`: New prompt command
- `internal/cli/builtin/`: Built-in prompts package
- `vibeguard.yaml`: Examples showing custom prompts

### Built-in Prompts Location
- Prompts compiled into binary via Go `embed` or constants
- Located in `internal/cli/builtin/prompts.go`
- Can be extended without recompilation (user config overrides)

## Next Steps

1. Implement config schema changes
2. Create prompt CLI command
3. Define built-in prompts
4. Update initialization flow to suggest `vibeguard prompt init`
5. Add comprehensive tests
6. Document migration path

## Related ADRs

- ADR-005: Adopt VibeGuard for policy enforcement
- ADR-006: Integrate VibeGuard as Git Pre-Commit Hook
