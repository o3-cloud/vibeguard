---
summary: Research findings on adding tags to vibeguard.yml for check categorization and filtering
event_type: research
sources:
  - internal/config/schema.go
  - internal/cli/check.go
  - internal/orchestrator/orchestrator.go
  - docs/ai-assisted-setup.md
  - examples/advanced.yaml
tags:
  - vibeguard
  - configuration
  - tags
  - filtering
  - schema-design
  - cli
  - feature-research
---

# Research: Adding Tags to vibeguard.yml

Researched the feasibility and design considerations for adding tags to the vibeguard.yml configuration schema.

## Current State

### Check Schema Structure

The current `Check` struct in `internal/config/schema.go` has no tag or category support:

```go
type Check struct {
    ID         string   `yaml:"id"`
    Run        string   `yaml:"run"`
    Grok       GrokSpec `yaml:"grok"`
    File       string   `yaml:"file"`
    Assert     string   `yaml:"assert"`
    Severity   Severity `yaml:"severity"`
    Suggestion string   `yaml:"suggestion"`
    Fix        string   `yaml:"fix,omitempty"`
    Requires   []string `yaml:"requires"`
    Timeout    Duration `yaml:"timeout"`
}
```

### Current Check Selection

- `vibeguard check` - runs all checks
- `vibeguard check <id>` - runs single check by exact ID match
- No filtering by tag or category

### Evidence of Planned Feature

The `docs/ai-assisted-setup.md` documentation already references categories and priorities that were planned but not implemented:

| Field | Description | Example Values |
|-------|-------------|----------------|
| Category | Check category | `lint`, `format`, `test`, `security` |
| Priority | Execution order | `10` (format), `30` (test) |

### Current Workaround

The `examples/advanced.yaml` uses comments to organize checks into phases:

```yaml
# PHASE 1: Fast Deterministic Checks
  - id: fmt
  - id: vet
# PHASE 2: Testing and Metrics
  - id: test
# PHASE 3: LLM-Powered Checks (Slower)
  - id: llm-architecture-review
```

## Recommended Tag Structure

Based on the documentation patterns:

**Standard Categories:**
- `build` - Build/compilation checks (priority 5)
- `format` - Code formatting (priority 10-11)
- `lint` - Static analysis (priority 15-20)
- `typecheck` - Type checking (priority 25)
- `test` - Unit/integration tests (priority 30)
- `coverage` - Coverage validation (priority 35)
- `security` - Security scanning (priority 50)

**Custom Tags:**
- `fast` / `slow` - Performance categorization
- `ci-only` / `pre-commit` - Execution context
- `go` / `node` / `python` - Language-specific

## Implementation Requirements

### Schema Changes (config/schema.go)

```go
type Check struct {
    // ... existing fields ...
    Tags     []string `yaml:"tags,omitempty"`
    Category string   `yaml:"category,omitempty"`
}
```

### CLI Changes (cli/check.go)

Add `--tags` flag support:
```
vibeguard check --tags format,lint
vibeguard check --tags security
```

### Files Requiring Changes

1. `internal/config/schema.go` - Add Tags field
2. `internal/config/config.go` - Add tag validation
3. `internal/cli/check.go` - Add `--tags` flag
4. `internal/orchestrator/orchestrator.go` - Filter by tags
5. `internal/cli/list.go` - Display tags
6. `internal/output/formatter.go` - Optional tag grouping
7. `internal/output/json.go` - Include tags in JSON

## Implementation Phases

1. **Phase 1 (Core):** Add Tags field, validation, basic filtering
2. **Phase 2 (CLI):** Implement `--tags` flag for check command
3. **Phase 3 (UX):** Update list and output formatters
4. **Phase 4 (Docs):** Implement Category/Priority from ai-assisted-setup.md
5. **Phase 5 (Advanced):** Tag-based grouping, discovery APIs

## Key Findings

1. Tags were conceptually planned (documented) but never implemented
2. Current workaround of comment-based phases shows clear user need
3. Inspector already has category concepts (`ToolCategory` enum)
4. Implementation is straightforward - schema + CLI + orchestrator changes
5. Backwards compatible - Tags field can be optional

## Next Steps

- Create feature issue in Beads for tag implementation
- Decide on tag vs category vs both approaches
- Consider priority field for execution ordering
