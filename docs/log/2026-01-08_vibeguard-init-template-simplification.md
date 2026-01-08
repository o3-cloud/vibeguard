---
summary: Research on vibeguard init strategy simplification using predefined templates with auto-discovery for existing projects and interactive selection for new projects
event_type: research
sources:
  - https://yeoman.io/
  - https://www.cookiecutter.io/
  - https://create-react-app.dev/docs/custom-templates/
  - https://github.com/vibeguard/vibeguard/internal/cli/init.go
  - https://github.com/vibeguard/vibeguard/internal/cli/templates/registry.go
tags:
  - init-command
  - template-system
  - scaffolding
  - CLI-architecture
  - user-experience
  - AI-agent-integration
  - project-detection
---

# VibeGuard Init Strategy Simplification Research

## Objective

Evaluate how to simplify vibeguard's init command by implementing a predefined template system with intelligent auto-discovery for AI agents and interactive selection for users on new projects.

## Current State Analysis

### Existing Implementation

The vibeguard codebase already has a robust foundation for template-based init:

**File Structure:**
- `internal/cli/init.go` - Main init command logic (262 lines)
- `internal/cli/templates/registry.go` - Template management system (63 lines)
- `internal/cli/templates/` - 8 template implementations:
  - go-standard, go-minimal
  - node-typescript, node-javascript
  - python-poetry, python-pip
  - rust-cargo
  - generic

**Current Template System Architecture:**

The registry pattern is elegant and extensible:
```go
// Each template self-registers via init()
func init() {
    Register(Template{
        Name: "node-typescript",
        Description: "TypeScript/Node.js project with ESLint, Prettier...",
        Content: `<YAML config here>`
    })
}
```

**Current Operating Modes:**

1. **Simple Init** (`vibeguard init`) - Uses hardcoded Go starter config
2. **Template Selection** (`vibeguard init --template go-standard`)
3. **List Templates** (`vibeguard init --template list`)
4. **AI-Assisted Setup** (`vibeguard init --assist`) - Auto-detects project and generates multi-section prompt

### Existing AI-Assisted Capabilities

The `--assist` mode is quite sophisticated:
- Uses `inspector.NewDetector()` to identify project type
- Uses `inspector.NewToolScanner()` to find installed tools
- Uses `inspector.NewMetadataExtractor()` for version info
- Uses `inspector.NewRecommender()` to suggest checks
- Generates comprehensive markdown prompts with:
  - Project analysis section
  - Tool inspection instructions
  - Language-specific recommendations
  - Validation rules and examples
  - Final task instructions

## Industry Best Practices Research

### Yeoman (JavaScript Scaffolding Tool)

**Model:** Generator-based framework
- Uses EJS templating language
- Generators are plugins run via `yo` command
- Generators can programmatically create files and generate tests
- More flexible automation than pure template systems
- Modern user experience with interactive prompts

**Key Insight:** Yeoman is not just templates—it's a framework for generators with programmatic logic.

### Cookiecutter (Python Scaffolding Tool)

**Model:** Configuration-first template system
- Uses Jinja2 templating language
- Over 6,000 pre-built templates available
- Stable, configuration-driven approach
- Easier to debug than programmatic systems
- Focus on simplicity

**Key Insight:** Config-driven systems are more stable and debuggable than programmatic ones.

### Create React App (React/Node.js)

**Model:** Template-based with npm registry integration
- Syntax: `npx create-react-app my-app --template [template-name]`
- Templates follow naming convention: `cra-template-[template-name]`
- Scoped templates supported: `@scope/cra-template-[name]`
- Users discover templates by searching npm registry
- Templates inherit all Create React App features

**Key Insight:** Simple naming conventions and npm registry integration provide discoverability without central documentation.

### Copier (Newer Alternative)

- Python-based alternative to Cookiecutter and Yeoman
- Combines strengths of both approaches
- Can be used as CLI or library
- Project templating with answers stored for updates

## Key Findings & Recommendations

### 1. Template Strategy (SIMPLIFICATION OPPORTUNITY)

**Current State:** Good foundation with 8 templates covering major languages

**Recommendation:** Expand template library with framework-specific variants

**Examples of simplified template naming:**
- **Languages:** `go`, `python`, `node`, `rust`, `java`, `csharp`
- **Node.js Frameworks:** `node-react-vite`, `node-react-cra`, `node-nextjs`, `node-vue`, `node-express`
- **Python Frameworks:** `python-django`, `python-fastapi`, `python-flask`
- **Go Frameworks:** `go-gin`, `go-echo`, `go-chi`
- **Frontend:** `frontend-react`, `frontend-vue`, `frontend-svelte`, `frontend-html`

This aligns with CRA's approach of framework-specific templates with simple naming.

### 2. AI Agent Auto-Discovery (LEVERAGE EXISTING --assist)

**Current Implementation:** Already detects project type and recommends checks

**Enhancement Opportunity:** Automatically select best-matching template

Flow for AI agents:
1. Run `vibeguard init --assist` (with maybe `--auto-template` flag)
2. System detects project structure, language, package manager, build tools
3. Agent receives prompt that includes recommended template name
4. Agent can either:
   - Accept recommendation: `vibeguard init --template <name>`
   - Or use full --assist flow: `vibeguard init --assist` to generate custom config

**Implementation Note:** The detection logic already exists in `inspector/` package. Could extend `NewRecommender()` to recommend template by name, not just checks.

### 3. Interactive Selection for New Projects

**Current State:** Interactive mode not implemented

**Recommendation:** Add `--interactive` flag for users starting new projects

Flow for human users:
```bash
vibeguard init --interactive
# OR just: vibeguard init
# (default could be interactive if no config exists and no --template specified)
```

Interactive flow:
1. Question: "What language/framework are you using?"
2. Show filtered templates for that choice
3. Allow user to preview template or select
4. Apply template to create vibeguard.yaml

**UI Pattern:** Similar to CRA's interactive selection or Yeoman's prompt system.

**Note:** This feature may not be necessary if AI agents handle template selection via --assist.

### 4. Template Discovery

**Current Implementation:** `vibeguard init --template list` shows all templates

**Enhancement Opportunity:** Make discovery more intuitive

Options:
1. **Simple:** Add `vibeguard list templates` command
2. **Better:** Add filtering: `vibeguard list templates --lang=node`
3. **Best:** Add categories in template metadata, show grouped:
   ```text
   Go Templates:
     - go-standard: Standard Go project with vet, fmt, test, build
     - go-minimal: Minimal Go with just vet and fmt

   Node.js Templates:
     - node-typescript: TypeScript/Node with ESLint, Prettier, tests
     - node-javascript: JavaScript/Node with ESLint, Prettier
   ```

### 5. Template Metadata Enhancement

**Current:** Each template has Name, Description, Content

**Recommendation:** Add metadata fields for better categorization and filtering:
```go
type Template struct {
    Name         string    // "node-typescript"
    Description  string    // Human-readable description
    Category     string    // "node" | "python" | "go" | etc
    Keywords     []string  // ["typescript", "eslint", "prettier"]
    RequiredTools []string // ["npm", "node"] for detection
    Content      string    // YAML config
    Version      string    // Template version for updates
}
```

This enables:
- Filtering by language/category
- Better AI agent detection matching
- Version management for template updates

## Implementation Path

### Simplified Approach: Let Agents Handle Discovery

Given that AI agents can intelligently select templates based on project analysis, vibeguard doesn't need built-in discovery logic.

**Workflow:**
1. Agent runs: `vibeguard init --assist`
2. Receives comprehensive project analysis (language, tools, structure, etc.)
3. Agent analyzes output and determines best-matching template
4. Agent runs: `vibeguard init --template <name>` with the selected template
5. Agent runs: `vibeguard check` to execute templated checks
6. If any checks fail, agent fixes the issues and re-runs `vibeguard check`
7. Repeat until all checks pass

This keeps vibeguard simple and focused while delegating intelligence to the agent. The `vibeguard check` step validates that the configuration works for the specific project and any tooling gaps are addressed.

### Phase 1: Add --list-templates Flag
- Add explicit `--list-templates` flag to init command for better discoverability
- Keep backward compatibility with `--template list`
- Minimal code change: register flag, check in runInit, delegate to listTemplates()
- Makes template listing more intuitive for users and agents

### Phase 2: Template Expansion (Quick Win)
- Add 5-10 framework-specific templates to existing registry
- Use existing naming/registration pattern
- No code changes needed, just new template files
- Examples: `node-react-vite`, `node-nextjs`, `python-fastapi`, `go-gin`

### Phase 3: Improve --template list Output (Optional)
- Better formatting of available templates
- Add descriptions to make selection easier for humans
- Could add grouping by language/framework

## Key Design Decisions Already Made (Per ADRs)

The project has established patterns:

**ADR-006:** VibeGuard integrates as git pre-commit hook
- Init command should support both interactive and automated flows
- AI-assisted mode aligns with automation needs

**ADR-004:** Code Quality Standards
- Templates should enforce project-specific quality standards
- Current templates model this well (format, lint, test, build pattern)

## Next Steps

1. **Immediate:** Expand template library with framework variants
   - `node-react-vite`, `node-nextjs`, `node-express`
   - `python-fastapi`, `python-django`
   - `go-gin`, `go-echo`
   - And others based on community need

2. **Future:** Optionally improve `--template list` output formatting for discoverability

## Conclusion

VibeGuard already has the perfect architecture for this simplified approach:

1. **`vibeguard init --assist`** - Provides agents with comprehensive project analysis
2. **`vibeguard init --template <name>`** - Applies any template to create config
3. **Agents decide** - Based on analysis, pick the best matching template

No complex discovery logic needed in vibeguard itself. The tool stays focused on policy enforcement while agents handle the intelligence. The only work required is expanding the template library to cover more frameworks—just new files, no code changes.
