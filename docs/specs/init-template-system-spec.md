---
title: VibeGuard Init Template System - Master Specification
date: 2026-01-08
status: draft
version: 1.0
author: Claude Code
---

# VibeGuard Init Template System - Master Specification

## Overview

This specification defines a simplified template-driven initialization strategy for VibeGuard. The system provides predefined configuration templates for popular languages and frameworks, with intelligent template discovery delegated to AI agents. **Templates become the single source of truth for project configurations**, enabling both user simplification and significant codebase simplification.

## Goals

1. **Simplify project setup** - Users and agents can quickly bootstrap VibeGuard configs for their projects
2. **Cover popular frameworks** - Provide templates for major language/framework combinations
3. **Delegate discovery to agents** - AI agents analyze projects and select appropriate templates
4. **Maintain tool focus** - VibeGuard remains focused on policy enforcement, not UI/discovery logic
5. **Enable codebase simplification** - Templates as source of truth eliminates duplication across examples, assist system, and recommendations

## Core Concept

The init system operates in three modes:

### Mode 1: AI-Assisted Setup (Agent Workflow)

```bash
# Step 1: Agent analyzes project
vibeguard init --assist > project-analysis.md
# (prompt includes instruction: "To see available templates, run: vibeguard init --list-templates")

# Step 2: Agent lists templates
vibeguard init --list-templates

# Step 3: Agent selects and applies template based on analysis
vibeguard init --template <selected-template>

# Step 4: Agent validates with checks
vibeguard check

# Step 5: Agent fixes failures and repeats until all checks pass
# (if any checks fail)
```

### Mode 2: Direct Template Selection (User/Agent)

```bash
# User or agent explicitly selects a template
vibeguard init --template node-typescript

# Creates vibeguard.yaml with that template
```

### Mode 3: List Available Templates (Discovery)

```bash
# User or agent lists available templates
vibeguard init --list-templates

# Output:
# Available templates:
#   generic                Generic template for custom setup
#   go-minimal            Minimal Go project (vet, fmt)
#   go-standard           Standard Go project (vet, fmt, test, build)
#   node-javascript       JavaScript/Node.js (ESLint, Prettier, tests)
#   node-typescript       TypeScript/Node.js (ESLint, Prettier, typecheck, tests)
#   python-pip            Python with pip (syntax check, tests)
#   python-poetry         Python with Poetry (syntax check, tests)
#   rust-cargo            Rust with Cargo (clippy, fmt, test, build)
```

## Command-Line Interface

### Flags

**Existing flags (unchanged):**
- `--force/-f` - Overwrite existing config file
- `--template/-t <name>` - Apply specific template (also accepts `list` for backward compatibility)
- `--assist` - Generate AI agent-assisted setup prompt
- `--output/-o <path>` - Output file for `--assist` mode

**New flags:**
- `--list-templates` - List available templates (explicit, clearer than `--template list`)

### Usage Examples

```bash
# Create config with default Go template
vibeguard init

# Use specific template
vibeguard init --template node-typescript

# List templates
vibeguard init --list-templates

# AI-assisted setup (generates analysis prompt)
vibeguard init --assist

# Save AI prompt to file
vibeguard init --assist --output setup-prompt.md

# Overwrite existing config
vibeguard init --template python-fastapi --force
```

## Template System Architecture

### Template Registry

Templates are registered via the `internal/cli/templates` package using a self-registering pattern:

```go
type Template struct {
    Name        string  // Unique identifier (e.g., "node-typescript")
    Description string  // Human-readable description
    Content     string  // YAML configuration content
}
```

Each template file calls `Register()` in its `init()` function:

```go
func init() {
    Register(Template{
        Name:        "node-typescript",
        Description: "TypeScript/Node.js with ESLint, Prettier, typecheck, tests",
        Content: `<YAML config...>`,
    })
}
```

### Template Naming Convention

Templates follow a consistent naming pattern:

- **Language with tooling variant:** `<language>-<variant>`
  - Variants: minimal, standard, strict (represent scope/coverage of checks)
  - Examples: `go-minimal`, `go-standard`, `python-pip`, `python-poetry`

- **Language with framework:** `<language>-<framework>`
  - Framework: web framework, build system, or language dialect
  - Examples: `node-express`, `python-django`, `python-fastapi`

- **Language with dialect and framework:** `<language>-<dialect>-<framework>`
  - Dialect: typescript, jsx variants
  - Examples: `node-typescript` (TypeScript dialect for Node), `node-react` (React framework for Node)

- **Language with framework and build tool:** `<language>-<framework>-<tool>`
  - Tool: vite, cra (Create React App), webpack, etc.
  - Examples: `node-react-vite`, `node-react-cra`

**Current templates:**
```
generic
go-minimal
go-standard
node-javascript
node-typescript
python-pip
python-poetry
rust-cargo
```

**Planned additions (Phase 2):**
```
node-react-vite
node-react-cra
node-nextjs
node-express
python-django
python-fastapi
python-flask
go-gin
go-echo
```

## Template Content Structure

All templates are YAML configurations with this structure:

```yaml
version: "1"

vars:
  # Project-specific variables for interpolation
  <var-name>: "<value>"

checks:
  - id: <check-id>
    run: <command>
    severity: error|warning
    suggestion: "<user-friendly message>"
    timeout: <duration>
    requires: [<dependency-ids>]  # Optional
    grok: [<patterns>]            # Optional
    assert: <condition>           # Optional
```

### Template Design Principles

1. **Language-appropriate tools** - Use tools native to the language ecosystem
2. **Reasonable defaults** - Include checks that matter most for that language/framework
3. **Clear suggestions** - Help users understand how to fix failures
4. **Dependency clarity** - Express check dependencies to ensure proper sequencing
5. **Variable flexibility** - Use template variables for customization

### Example: Node.js TypeScript Template

```yaml
version: "1"

vars:
  source_dir: "src"
  min_coverage: "70"

checks:
  - id: format
    run: npx prettier --check .
    severity: error
    suggestion: "Run 'npx prettier --write .' to format your code"
    timeout: 60s

  - id: lint
    run: npx eslint {{.source_dir}} --max-warnings 0
    severity: error
    suggestion: "Run 'npx eslint {{.source_dir}} --fix' to fix linting issues"
    timeout: 60s

  - id: typecheck
    run: npx tsc --noEmit
    severity: error
    suggestion: "Fix TypeScript type errors shown in the output"
    timeout: 120s

  - id: test
    run: npm test -- --passWithNoTests
    severity: error
    suggestion: "Run 'npm test' to diagnose test failures"
    timeout: 300s
    requires:
      - lint
      - typecheck

  - id: build
    run: npm run build
    severity: error
    suggestion: "Run 'npm run build' to diagnose build errors"
    timeout: 120s
    requires:
      - typecheck
```

## Agent Integration

### AI Agent Workflow

When an AI agent initializes VibeGuard on a project:

1. **Analyze** - Run `vibeguard init --assist` to get comprehensive project analysis
2. **Discover** - Run `vibeguard init --list-templates` to see available templates (per assist instructions)
3. **Select** - Agent analyzes project information and selects the best-matching template
4. **Apply** - Run `vibeguard init --template <name>` with selected template
5. **Validate** - Run `vibeguard check` to verify configuration works
6. **Fix** - If checks fail, agent fixes issues and re-runs `vibeguard check`

### Assist Prompt Integration

The `--assist` output includes:

- **Project detection results** - Language, framework, tools found
- **Tool inspection guidance** - How to analyze existing build/test configs
- **Template discovery instruction** - "To see available templates, run: `vibeguard init --list-templates`"
- **Recommended checks** - Checks that should be included based on detected tools
- **Validation rules** - Rules for creating correct YAML configurations
- **Examples** - Language-specific configuration examples
- **Final task** - Instructions to apply template, run `vibeguard check`, and fix failures

The assist prompt provides analysis; agents use this analysis to make intelligent template selection from the discovered list.

## Implementation Phases

### Phase 1: Add --list-templates Flag and Assist Integration

**Timeline:** Immediate (can implement now)
**Scope:** Minimal code change + assist output update
**Risk Level:** Low

#### 1.1 Critical Changes Required

**File: `internal/cli/init.go`**

1. **Error Messages** (HIGH PRIORITY):
   - Line 114: Change error text from `"unknown template %q (use --template list to see available templates)"` to `"unknown template %q (use --list-templates to see available templates)"`
   - Line 46: Update flag help text from `"Use a predefined template (use 'list' to see available templates)"` to `"Use a predefined template (run 'vibeguard init --list-templates' to see available templates)"`

2. **Flag Already Implemented** ✓ (No changes needed):
   - `--list-templates` boolean flag already declared (line 47)
   - Flag handler already in place (lines 97-98)
   - Backward compatibility with `--template list` already maintained (lines 102-104)

**File: `internal/cli/assist/sections.go` (522 lines)**

3. **Add Template Discovery Section** (HIGH PRIORITY):
   - Create new `TemplateDiscoverySection()` function
   - Insert instruction: "To see available templates, run: `vibeguard init --list-templates`"
   - Include example: "Based on your Go project, try: `vibeguard init --template go-standard`"
   - Explain when to use templates vs. custom configs

**File: `internal/cli/assist/composer.go` (170 lines)**

4. **Update Assist Composition** (HIGH PRIORITY):
   - Update `buildSections()` to include new `TemplateDiscoverySection()`
   - Position: After Recommendations section, before ConfigRequirements
   - This guides agents to discover templates before generating custom config

#### 1.2 Documentation Updates (MEDIUM PRIORITY)

**Files affected:**
- `docs/CLI-REFERENCE.md` - Add `--list-templates` documentation
- `docs/GETTING_STARTED.md` - Add template discovery section
- `README.md` - Update AI-assisted setup section with full workflow
- `docs/AGENTS.md` (new or updated) - Document agent workflow for templates
- `examples/pre-commit/README.md` - Add init prerequisites section

#### 1.3 Tests (MEDIUM PRIORITY)

**File: `internal/cli/init_test.go`**

1. Update existing test:
   - `TestRunInit_UnknownTemplate()` should verify error mentions `--list-templates`

2. Add new tests:
   - `TestRunInit_ListTemplatesFlag()` - Verify flag lists templates
   - `TestRunInit_ListTemplatesAndTemplateConflict()` - Test error handling
   - `TestRunInit_ListTemplatesAndAssistConflict()` - Test error handling

**Deliverable:**
- Explicit template listing capability works
- Assist prompt directs agents to discover templates
- Error messages reference `--list-templates`
- All existing tests pass

### Phase 2: Template Expansion

**Timeline:** 1-2 weeks after Phase 1 is stable
**Scope:** Add new template files (no code changes)
**Risk Level:** Very Low

**New templates to add:**

Language-specific variants:
- `node-react-vite` - React with Vite build tool
- `node-react-cra` - React with Create React App
- `node-nextjs` - Next.js (full-stack React)
- `node-express` - Express.js backend
- `python-django` - Django web framework
- `python-fastapi` - FastAPI modern async framework
- `python-flask` - Flask lightweight framework
- `go-gin` - Gin web framework
- `go-echo` - Echo web framework

**Deliverable:** Rich template library covering major frameworks and variations

### Phase 3: Codebase Simplification

**Timeline:** After Phase 1 and 2 are proven stable
**Scope:** Code cleanup and consolidation
**Risk Level:** Medium

This phase enables **1,800+ lines of code reduction** by making templates the source of truth:

#### 3.1 Phase 3A: Quick Wins (Zero Breaking Changes)

1. **Delete `/examples/` directory** (5 files, ~200 lines):
   - Configuration examples duplicate template content
   - Reference actual templates in documentation instead
   - Update `examples/README.md` to explain templates are the canonical source

2. **Delete `assist/templates.go`** (391 lines):
   - Hardcoded YAML language examples duplicate actual templates
   - Update `LanguageExamplesSection()` to reference actual templates
   - Show: "Here's the actual Go template we use:" followed by template content

3. **Consolidate Data Structures** (50 lines):
   - Remove duplicate `CheckRecommendation` types
   - Keep single `Template` struct as canonical

**Estimated Savings:** ~600 lines, 0 breaking changes, **recommend doing immediately after Phase 1 proves stable**

#### 3.2 Phase 3B: Architectural Simplification (Medium Risk)

1. **Remove `starterConfig` constant** (35 lines):
   - Embedded Go config duplicates `go-standard.go` template
   - Update `runInit()` to require explicit template selection
   - Option: Default to `go-standard` automatically for Go projects

2. **Simplify Recommendations System** (425 lines → 50 lines):
   - Current: Generates detailed check recommendations mirroring templates
   - Simplified: Match project type to template name instead
   - Remove ~12 tool-specific recommendation methods

3. **Consolidate Assist Sections** (250+ lines reducible):
   - Remove or merge `RecommendationsSection` (now template matching)
   - Remove `LanguageExamplesSection` (replaced by actual templates)
   - Simplify `ConfigRequirementsSection` (reference templates for examples)

**Estimated Savings:** ~900+ lines

#### 3.3 Phase 3C: Deep Refactoring (Higher Risk)

1. **Complete Rewrite of `inspector/recommendations.go`** (525 lines):
   - Focus on template matching instead of generating recommendations
   - Replace ~12 tool-specific methods with unified template detection

2. **Redesign Assist System Around Templates**:
   - Narrow scope: analyze project → recommend template
   - Remove custom config generation from assist
   - Templates become the expected output path

**Estimated Savings:** ~300+ additional lines

#### Phase 3 Overall Impact

| Module | Before | After Simplification | Reduction |
|--------|--------|-----|-----------|
| init.go | 300 | 265 | -12% |
| assist/ | 5,344 | 3,500 | -35% |
| inspector/ | 525 | 100 | -81% |
| examples/ | 200 | 0 | -100% |
| **TOTAL** | **7,603** | **5,099** | **-33%** |

**Total Code Reduction: 1,504 lines (-20%)**

### Phase 4: Optional Enhancements (Future)

**Possible improvements:**
- Better `--template list` output formatting (grouping by language)
- Template metadata fields for filtering/categorization
- Template versioning and updates
- Community template registry integration
- Interactive template selection for users

## Success Criteria

### Phase 1: Flag and Integration
- [ ] `vibeguard init --list-templates` exits with code 0
- [ ] `vibeguard init --list-templates` outputs all registered templates in `name description` format
- [ ] `vibeguard init --assist` output includes instruction: "To see available templates, run: `vibeguard init --list-templates`"
- [ ] Error message for unknown template reads: "unknown template '<name>' (use --list-templates to see available templates)"
- [ ] `--template list` still works for backward compatibility
- [ ] Conflicting flags produce clear error messages (exit code 1)

### Phase 2: Template Expansion
- [ ] Each new template YAML is valid (runs through `vibeguard validate`)
- [ ] Each template covers its intended language/framework
- [ ] Templates follow naming convention: `<language>-<variant>` or `<language>-<framework>`
- [ ] Each template includes reasonable defaults for its ecosystem

### Overall
- [ ] No existing functionality broken; backward compatible
- [ ] Documentation updated with agent workflow examples
- [ ] Code follows project conventions (ADRs, style guide)

## Error Handling and Flag Combinations

### Conflicting Flags

- **`--assist` + `--template`** - Error: These modes are mutually exclusive
- **`--assist` + `--list-templates`** - Error: Cannot use both together
- **`--template <name>` + `--list-templates`** - Error: Cannot use both together
- **`--template <nonexistent>`** - Error: "unknown template '<name>' (use --list-templates to see available templates)"
- **`--template list` + `--list-templates`** - Error: Cannot use both (--template list is redundant, use --list-templates)

### Exit Codes

- **0** - Success (config created, templates listed, or prompt generated)
- **1** - Flags conflict or invalid combination
- **2** - Template not found
- **3** - File I/O error
- **4** - Invalid configuration (validation failed)

## Risk Assessment and Mitigation

### Phase 1: Low Risk
**Changes:** Error messages, assist integration, tests
**Risks:**
- Assist prompt change might confuse agents familiar with old format (LOW)
- Error message change affects error handling scripts (LOW)

**Mitigation:**
- Update documentation and agent guides before release
- Thorough testing of assist output format
- Version assist prompts if needed

### Phase 2: Very Low Risk
**Changes:** Add new template files only
**Risks:**
- None - additive change only
- New templates don't break existing functionality
- Can test each template independently

**Mitigation:**
- Validate each template with `vibeguard validate`
- Test each template with sample projects
- Run `vibeguard check` on each template

### Phase 3A: Minimal Risk
**Changes:** Delete redundant files and consolidate structures
**Risks:**
- Removing examples breaks documentation links (LOW)
- Removing assist/templates.go affects LanguageExamplesSection rendering (LOW)

**Mitigation:**
- Update all documentation references before deletion
- Implement new LanguageExamplesSection that references templates
- Create replacement examples in docs/guides/
- Thorough testing of assist output

### Phase 3B: Medium Risk
**Changes:** Remove starterConfig, simplify recommendations, consolidate sections
**Risks:**
- Removing default behavior forces users to choose template (MEDIUM)
- Simplifying recommendations might degrade agent assistance (MEDIUM)
- Merging assist sections changes prompt structure (MEDIUM)

**Mitigation:**
- Only proceed after Phase 1-2 proven stable for 2+ weeks
- Implement smart defaults (detect Go → suggest go-standard)
- Extensive testing with real agent workflows
- A/B test prompt changes with agents before full rollout
- Clear migration path in documentation

### Phase 3C: Higher Risk
**Changes:** Complete rewrite of recommendation system and assist redesign
**Risks:**
- Agent workflows may need adjustment (MEDIUM-HIGH)
- Significant code changes increase bug risk (MEDIUM)
- May require adjusting agent integration patterns (MEDIUM)

**Mitigation:**
- Only proceed after Phase 3A-3B proven stable for 1+ month
- Comprehensive test suite for recommendations
- Agent integration testing with multiple agent types
- Gradual rollout to dogfooding first

### Overall Recommendation
- **Phase 1:** Implement immediately - low risk, high value
- **Phase 2:** Implement 1-2 weeks after Phase 1 stable - very low risk
- **Phase 3A:** Implement 2-3 weeks after Phase 2 stable - minimal risk, high value
- **Phase 3B:** Implement 1 month after Phase 3A - medium risk, evaluate carefully
- **Phase 3C:** Consider carefully - only if team committed to refactoring

## Success Metrics

### Phase 1 Success
- Agents receive clear instruction to discover templates
- Error messages guide users to `--list-templates`
- All tests pass with new assertions
- Documentation updated and clear
- No regressions in existing functionality

### Phase 2 Success
- Template library covers 10+ popular frameworks
- Each template validates correctly
- Agents can successfully select appropriate templates
- No existing templates broken or changed

### Phase 3 Success
- Codebase reduced by 1,500+ lines
- Single source of truth (templates)
- Agent workflows unchanged or improved
- Maintenance burden significantly reduced
- All tests pass with improved coverage

## Out of Scope

- Interactive CLI mode for humans (agents handle intelligent selection)
- Template metadata enhancements (future phase)
- Remote template registry (keep templates in codebase)
- Template generation/customization UI
- Template versioning beyond code versions
- Monorepo/multi-language project strategies (future phase)
- Template composition (combining multiple templates)

## Key Design Principles

1. **Templates are canonical** - Single source of truth for configurations
2. **Explicit over implicit** - Users/agents always know what template they're using
3. **Agents are intelligent** - Agents handle discovery and selection, tool provides options
4. **Tool focus** - VibeGuard stays focused on policy enforcement
5. **Backward compatibility** - Existing functionality never broken
6. **Simplification through consolidation** - Templates eliminate duplication

## References

### Implementation
- Current implementation: `internal/cli/init.go`, `internal/cli/templates/`
- Project detection: `internal/cli/inspector/`
- Assist system: `internal/cli/assist/`
- Tests: `internal/cli/init_test.go`

### Architecture Decision Records
- ADR-003: Adopt Go as the Primary Implementation Language
- ADR-004: Establish Code Quality Standards and Tooling
- ADR-005: Adopt VibeGuard for Policy Enforcement in CI/CD
- ADR-006: Integrate VibeGuard as Git Pre-Commit Hook for Policy Enforcement

### Related Documentation
- Research log: `docs/log/2026-01-08_vibeguard-init-template-simplification.md`
- CLI Reference: `docs/CLI-REFERENCE.md`
- Getting Started: `docs/GETTING_STARTED.md`
- Main README: `README.md`

### Appendix: Specification Evolution

This specification consolidates findings from:
1. Initial template strategy research (research-phase)
2. Codebase impact analysis (implementation-phase)
3. Simplification opportunity analysis (optimization-phase)

Version history:
- v1.0 (2026-01-08): Initial master specification consolidating all phases
