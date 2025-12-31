---
summary: Implementation specification for AI agent-assisted setup feature enabling automatic VibeGuard configuration generation
event_type: deep dive
sources:
  - docs/log/2025-12-31_ai-agent-assisted-setup-research.md
  - internal/cli/init.go (Current init command)
  - examples/ (Configuration examples)
  - cmd/vibeguard/main.go (Entry point)
tags:
  - agent-assisted-setup
  - implementation-spec
  - automation
  - config-generation
  - ai-integration
  - prompt-engineering
  - cli-design
  - repository-inspection
---

# AI Agent-Assisted Setup Implementation Specification

## 1. Feature Overview

### Purpose
Enable AI coding agents (Claude Code, Cursor, etc.) to automatically generate valid VibeGuard configurations by providing a structured prompt that guides configuration generation based on project analysis.

### Success Criteria
- Prompt enables agents to generate valid vibeguard.yaml without manual iteration
- Detection accurately identifies project type and existing tools
- Generated configurations pass schema validation
- Works with any LLM-based coding agent without agent-specific integration
- Reduces onboarding time from hours to minutes

### Key Design Principles
1. **Prompt-Based Interface** - Text output, no binary dependencies or special protocols
2. **Context-Aware** - Analyzes actual project structure and tools
3. **Clear Constraints** - Explicit rules for valid configurations
4. **Human-Reviewable** - Generated configs are YAML and human-readable
5. **Graceful Degradation** - Works even with partial information

---

## 2. Technical Architecture

### High-Level Flow

```
User Project
     ↓
vibeguard init --assist
     ↓
Repository Inspector (detect project type, tools, structure)
     ↓
Prompt Composer (create AI-readable setup guide)
     ↓
Output to stdout (or file)
     ↓
[User pipes to Claude Code / Cursor / other agent]
     ↓
AI Agent reads prompt
     ↓
AI Agent generates vibeguard.yaml
     ↓
User reviews and commits config
     ↓
vibeguard runs checks
```

### Core Components

#### 2.1 Repository Inspector
**Responsibility:** Analyze project structure and detect tools

**Functionality:**
- Project type detection (Go, Node.js, Python, Rust, etc.)
- Tool detection (linters, test frameworks, CI/CD, formatters)
- Metadata extraction (package info, entrypoints, common dirs)
- Existing configuration discovery

**Implementation Location:** `internal/cli/inspector/` package
- `detector.go` - Project type detection logic
- `tools.go` - Tool scanning and discovery
- `metadata.go` - Project metadata extraction
- `patterns.go` - Common path and file patterns

**Output:** Structured `ProjectAnalysis` data structure

#### 2.2 Prompt Composer
**Responsibility:** Create AI-readable setup instructions

**Functionality:**
- Render project analysis into natural language
- Include check definition templates and examples
- Document validation requirements
- Provide fallback guidance for ambiguous cases
- Format for easy parsing by AI agents

**Implementation Location:** `internal/cli/assist/` package
- `composer.go` - Main prompt composition logic
- `sections.go` - Individual prompt sections
- `templates.go` - Reusable text blocks and examples
- `validator_guide.go` - Validation instructions

**Output:** Complete prompt string ready for piping to AI agent

#### 2.3 CLI Command
**Responsibility:** Wire together inspector and composer

**Implementation Location:** `internal/cli/assist.go` (new file)
- New command handler: `HandleAssist()`
- Flag parsing: `--output` (stdout/file), `--verbose`, etc.
- Error handling for missing project files
- Output formatting

---

## 3. Data Structures

### ProjectAnalysis
```go
type ProjectAnalysis struct {
    ProjectType    ProjectType           // go, node, python, rust, etc.
    Tools          map[string]ToolInfo   // detected tools and versions
    Structure      ProjectStructure      // directories, files, entry points
    Metadata       ProjectMetadata       // package info, configs
    Confidence     float64               // 0-1: confidence in detection
    Recommendations []CheckRecommendation
    Warnings       []string              // ambiguities or missing info
}

type ProjectType string
const (
    Go     ProjectType = "go"
    Node   ProjectType = "node"
    Python ProjectType = "python"
    Rust   ProjectType = "rust"
    Java   ProjectType = "java"
    Ruby   ProjectType = "ruby"
    Unknown ProjectType = "unknown"
)

type ToolInfo struct {
    Name        string    // tool name
    Detected    bool      // whether it was found
    Version     string    // if detectable
    ConfigFile  string    // path to config if found
    Confidence  float64   // 0-1: how certain we are
}

type ProjectStructure struct {
    Root           string
    SourceDirs     []string
    TestDirs       []string
    EntryPoints    []string
    ConfigFiles    []string
}

type ProjectMetadata struct {
    Name        string
    Version     string
    Description string
    Language    string
    BuildTool   string
}

type CheckRecommendation struct {
    ID          string
    Description string
    Rationale   string
    Command     string
    Grok        string
    Assert      string
    Severity    string
}
```

### PromptTemplate
```go
type PromptTemplate struct {
    Header            string
    ProjectAnalysis   string
    Instructions      string
    Examples          string
    ValidationGuide   string
    QuestionGuidance  string
    Footer            string
}
```

---

## 4. Implementation Phases

### Phase 1: Repository Inspector
**Goals:** Detect project type and tools accurately

**Tasks:**
1. Implement `ProjectType` detection
   - Go: look for `go.mod`, `*.go` files
   - Node: look for `package.json`, `node_modules/`
   - Python: look for `pyproject.toml`, `setup.py`, `requirements.txt`
   - Ruby: look for `Gemfile`, `*.rb` files
   - Rust: look for `Cargo.toml`

2. Implement tool detection
   - Go linters: golangci-lint, gofmt, govet
   - Node tools: eslint, prettier, jest, mocha
   - Python tools: black, pylint, pytest, mypy
   - CI/CD: GitHub Actions, GitLab CI, CircleCI, Jenkins
   - Version control: git hooks existing config

3. Create metadata extraction
   - Parse package.json, go.mod, pyproject.toml
   - Identify entry points and main source directories
   - Extract version and name information

4. Unit tests for each detector
   - Test fixtures with sample projects
   - Verify accuracy on Go, Node, Python projects
   - Test edge cases and partial detection

**Acceptance Criteria:**
- Detects 90%+ of tools in standard projects
- Handles missing/incomplete configurations gracefully
- Fast execution (<500ms on typical projects)

### Phase 2: Check Recommendations
**Goals:** Suggest appropriate checks based on detected tools

**Tasks:**
1. Create recommendation engine
   - Map detected tools → suggested checks
   - Define sensible defaults for each tool
   - Handle tool-specific variations

2. Build check templates for each tool category
   - Go: fmt, vet, test with coverage, lint
   - Node: eslint, prettier, test, npm audit
   - Python: black, pylint, pytest, mypy
   - General: git hooks, dependency audits

3. Implement reasoning/rationale generation
   - Explain why each check is recommended
   - Link to tool documentation
   - Show example configuration snippets

4. Unit tests
   - Verify recommendation logic
   - Test template rendering
   - Validate generated check YAML syntax

**Acceptance Criteria:**
- Recommendations cover primary tools for each language
- Generated check YAML is valid and runnable
- Rationale is clear and useful

### Phase 3: Prompt Composition
**Goals:** Create effective AI-readable setup instructions

**Tasks:**
1. Design prompt structure
   - Clear sections with headers
   - Unambiguous constraint specification
   - Rich, diverse examples
   - Explicit validation rules

2. Implement prompt composer
   - Render ProjectAnalysis into natural language
   - Include project-specific examples
   - Add language-specific guidance
   - Insert tool-specific patterns and assertions

3. Create reusable templates
   - Project overview template
   - Check definition examples (Go, Node, Python, etc.)
   - Validation rules template
   - Clarification question guidance

4. Generate sample prompts for review
   - Go project with common tools
   - Node.js project with ESLint, Jest
   - Python project with pytest, black
   - Multi-language project

5. Test with actual agents
   - Feed generated prompt to Claude Code
   - Validate generated configuration
   - Iterate on clarity and completeness

**Acceptance Criteria:**
- Prompt < 4000 tokens (fits in context efficiently)
- Generated configs pass validation on first attempt
- Agents understand what checks to create
- No need for follow-up clarification

### Phase 4: CLI Integration
**Goals:** Wire everything together in user-facing command

**Tasks:**
1. Create `vibeguard init --assist` command
   - Parse flags: `--output`, `--verbose`, `--format`
   - Call inspector and composer
   - Handle errors gracefully
   - Output to stdout or file

2. Error handling
   - Missing go.mod/package.json → suggest running in wrong directory
   - Undetectable projects → provide fallback prompt
   - Permission errors → clear error messages

3. Verbose output mode
   - Show what was detected
   - Display confidence levels
   - List tools found
   - Show any warnings

4. Documentation
   - Update README with `--assist` option
   - Add examples to CONTRIBUTING
   - Document for user manual
   - Create demo prompt examples

5. Integration testing
   - Test on real projects (vibeguard itself, well-known Go/Node projects)
   - Verify prompt generation on various systems
   - Test file output vs stdout

**Acceptance Criteria:**
- `vibeguard init --assist` works on real projects
- Output is clean and properly formatted
- Help text is clear and discoverable
- Documentation is complete

### Phase 5: Refinement & Testing
**Goals:** Test with actual agents and gather feedback

**Tasks:**
1. Test with Claude Code
   - Feed prompt to Claude Code assistant
   - Verify generated configs work
   - Iterate on ambiguities
   - Document findings

2. Test with other agents (if possible)
   - Cursor, other coding tools
   - Verify prompt generalization

3. Test on diverse projects
   - Simple single-tool projects
   - Complex multi-tool projects
   - Minimal projects (edge cases)
   - Projects with unusual structures

4. Refinement loop
   - Identify confusing sections in prompt
   - Improve examples based on failures
   - Add missing guidance
   - Validate all generated configs

5. Performance optimization
   - Profile inspector on large codebases
   - Optimize file scanning
   - Cache tool detection if needed

**Acceptance Criteria:**
- 95%+ success rate on generated configs
- All configs pass validation
- Agents rarely need follow-up prompts
- Fast execution (<1s typical)

---

## 5. CLI Design

### Command Structure

```bash
# Generate setup prompt and output to stdout
vibeguard init --assist

# Generate and save to file
vibeguard init --assist --output setup-guide.txt

# Verbose mode with detection details
vibeguard init --assist --verbose

# Pipe to AI agent
vibeguard init --assist | pbcopy  # macOS
vibeguard init --assist | xclip   # Linux
```

### Help Text

```
vibeguard init --assist

Usage:
  vibeguard init --assist [flags]

Flags:
  -o, --output string      Save prompt to file instead of stdout
  -v, --verbose            Show detection details and confidence levels
  -h, --help               Show this help message

Description:
  Analyze your project and generate an AI-readable setup guide.
  Use this to help coding agents (Claude Code, Cursor, etc.)
  automatically generate a valid vibeguard configuration.

Examples:
  # Output to terminal (copy to Claude Code)
  vibeguard init --assist

  # Save to file
  vibeguard init --assist --output guide.txt

  # See detection details
  vibeguard init --assist --verbose
```

### Exit Codes
- `0` - Success, prompt generated
- `2` - Invalid directory (no valid project detected)
- `3` - Runtime error (disk read, etc.)

### Predefined Templates

For users without AI agents or who prefer a guided approach, provide predefined configuration templates that can be selected interactively or via flag.

```bash
# List available templates
vibeguard init --template list

# Use a specific template
vibeguard init --template go-standard
vibeguard init --template node-typescript
vibeguard init --template python-poetry

# Interactive template selection
vibeguard init --template
```

**Built-in Templates:**

| Template | Description | Checks Included |
|----------|-------------|-----------------|
| `go-standard` | Standard Go project | fmt, vet, lint, test with coverage |
| `go-minimal` | Minimal Go checks | fmt, vet, test |
| `node-typescript` | TypeScript Node.js project | eslint, prettier, tsc, jest |
| `node-javascript` | JavaScript Node.js project | eslint, prettier, jest |
| `python-poetry` | Python with Poetry | black, pylint, mypy, pytest |
| `python-pip` | Python with pip | black, flake8, pytest |
| `rust-cargo` | Rust project | fmt, clippy, test |
| `generic` | Language-agnostic basics | git hooks, file checks |

**Template Selection Flow:**
1. If `--template` provided without value, show interactive picker
2. If `--template <name>` provided, use that template directly
3. Template is customized with detected project metadata (name, paths)
4. User reviews generated config and can edit before saving

**Implementation Location:** `internal/cli/templates/` package
- `registry.go` - Template registration and lookup
- `templates/*.yaml` - Predefined template files
- `customizer.go` - Apply project-specific values to templates

---

## 6. Prompt Template Structure

### Section 1: Header & Overview
```
# VibeGuard AI Agent Setup Guide

You are being asked to help set up VibeGuard policy enforcement
for a software project. VibeGuard is a declarative policy tool
that runs quality checks and assertions on code.

This guide will help you understand the project structure,
existing tools, and how to generate a valid configuration.
```

### Section 2: Project Analysis
```
## Project Analysis

**Project Type:** Go
**Confidence:** 95%
**Main Tools Detected:**
- go mod (package management)
- golangci-lint (linting)
- Go test (testing)
- Git (version control)

**Project Structure:**
- Source code: cmd/, internal/, pkg/
- Tests: **/*_test.go
- Go version: 1.21
```

### Section 3: Configuration Requirements
```
## Configuration Requirements

A valid vibeguard.yaml must contain:

### 1. Variables (optional)
Global variables for interpolation in commands and assertions.

Example:
variables:
  go_packages: "./cmd/... ./internal/... ./pkg/..."
  test_dir: "./..."

### 2. Checks (required)
Array of check definitions. Each check must have:
- id: Unique identifier (alphanumeric + underscore)
- run: Shell command to execute
- (optional) grok: Pattern to extract data from output
- (optional) assert: Condition that must be true
- (optional) requires: IDs of checks that must pass first
- (optional) severity: "error" or "warning"
```

### Section 4: Language-Specific Examples
```
## Go-Specific Examples

### Format Check
id: go_fmt
run: gofmt -l -w {{.go_packages}}
assert: stdout == ""
severity: error
suggestion: "Run 'gofmt -l -w ./...' to fix formatting"

### Lint Check
id: go_lint
run: golangci-lint run {{.go_packages}}
assert: exit_code == 0
severity: error

### Test Check
id: go_test
run: go test -v -coverprofile=coverage.out {{.go_packages}}
grok: 'coverage: (?P<coverage_pct>\d+\.\d+)%'
assert: coverage_pct >= 70
severity: error
suggestion: "Coverage below 70%. Add tests to improve coverage."
```

### Section 5: Validation Rules
```
## Validation Rules

Your generated configuration must:

1. Be valid YAML syntax
2. Include at least one check
3. Each check must have unique id
4. Each check must have non-empty run command
5. All requires references must point to existing checks
6. No circular dependencies
7. All variables used in {{.var}} must be defined
8. Severity must be "error" or "warning"

**DO NOT:**
- Include comments in generated YAML
- Add extra top-level keys beyond variables, checks
- Use undefined variables
- Create checks for non-existent tools
```

### Section 6: Multi-Turn Conversation Protocol
```
## Conversation Protocol

This section defines how the AI agent should interact with the user
during configuration generation. Follow this protocol for consistent,
helpful interactions.

### Asking Clarifying Questions

If the project is ambiguous, ask clarifying questions BEFORE generating:

1. **Multiple test frameworks detected?**
   "I see both Jest and Mocha. Which test runner should I use?"

2. **Optional tools not detected?**
   "Should I include a dependency vulnerability check?"

3. **Coverage threshold unknown?**
   "What's the minimum code coverage required?"

Format questions as a JSON block for easy parsing:

\`\`\`json
{
  "status": "needs_clarification",
  "questions": [
    {
      "id": "test_framework",
      "question": "Which test framework should I configure?",
      "options": ["jest", "mocha", "vitest"],
      "default": "jest",
      "reason": "Multiple test frameworks detected in package.json"
    }
  ]
}
\`\`\`

### Reporting Validation Failures

If the generated configuration fails validation, report errors clearly:

\`\`\`json
{
  "status": "validation_failed",
  "errors": [
    {
      "field": "checks[2].assert",
      "message": "Invalid assertion syntax: missing operator",
      "suggestion": "Change 'coverage_pct 70' to 'coverage_pct >= 70'"
    }
  ],
  "revised_config": "... corrected YAML here ..."
}
\`\`\`

### Iterative Refinement

When the user requests changes to a generated config:

1. **Acknowledge** the requested change
2. **Show only the diff** or affected section, not the entire config
3. **Explain** why the change was made
4. **Validate** the updated config before presenting

Example response format:

\`\`\`
I'll update the coverage threshold from 70% to 80%.

Changed check `go_test`:
- assert: coverage_pct >= 70
+ assert: coverage_pct >= 80

The updated configuration passes validation.
\`\`\`

### Requesting Additional Context

If the agent needs more information about the project:

\`\`\`json
{
  "status": "needs_context",
  "requests": [
    {
      "type": "file_content",
      "path": "Makefile",
      "reason": "To understand existing build commands"
    },
    {
      "type": "command_output",
      "command": "go list ./...",
      "reason": "To enumerate all Go packages"
    }
  ]
}
\`\`\`

### Success Response

When configuration is complete and validated:

\`\`\`json
{
  "status": "complete",
  "config_valid": true,
  "checks_count": 5,
  "summary": "Generated 5 checks: go_fmt, go_vet, go_lint, go_test, go_build"
}
\`\`\`

Followed by the complete YAML configuration in a code block.
```

### Section 7: Generated Configuration Template
```
## Your Task

Based on the project analysis above, generate a vibeguard.yaml
configuration that:

1. Defines appropriate variables for this project
2. Creates checks for the detected tools
3. Follows the syntax rules described above
4. Includes helpful suggestions for each check

Output the configuration in a YAML code block:

\`\`\`yaml
# vibeguard.yaml
variables:
  # ... your variables ...

checks:
  # ... your checks ...
\`\`\`

If you have clarifying questions, ask them first before generating.
```

---

## 7. Testing Strategy

### Unit Tests

**Repository Inspector Tests:**
- Test project type detection on sample projects (Go, Node, Python, Rust)
- Verify tool detection (present and absent cases)
- Test metadata extraction accuracy
- Edge cases: empty directories, missing files, malformed configs

**Check Recommendation Tests:**
- Verify recommendations for each tool type
- Test template rendering with various ProjectAnalysis inputs
- Validate generated YAML syntax

**Prompt Composer Tests:**
- Test prompt generation with various inputs
- Verify section rendering
- Test token count (stay < 4000)
- Validate prompt clarity with heuristic checks

### Integration Tests

**End-to-End Tests:**
- Run `vibeguard init --assist` on test projects
- Capture prompt output
- Feed prompt to test agent (mock or simple validator)
- Verify generated config passes `vibeguard validate`

**Real Project Tests:**
- Test on vibeguard itself
- Test on popular Go projects (Hugo, Kubernetes client libs)
- Test on Node projects (Next.js, Express examples)
- Test on Python projects (Django, Flask examples)

### Manual Testing

**Agent Testing:**
- Feed prompt to Claude Code, verify generation
- Test with Cursor if available
- Document findings and iterate

**User Acceptance Testing:**
- Have team members try feature
- Gather feedback on prompt clarity
- Identify confusing sections
- Iterate on examples

### Success Guidelines

These are general targets to aim for, not strict enforcement criteria:

- Most generated configs should pass validation without manual fixes
- Execution should feel fast (sub-second on typical projects)
- Agents should generally understand the prompt without extensive follow-up
- Users should find the prompt clear and the generated configs useful

### AI Agent Configuration Review

As a final quality step, the AI agent should review the generated configuration and provide improvement suggestions. This creates a feedback loop that helps users understand their config and catch potential issues.

**Review Prompt (included in prompt template):**

```
## Configuration Review

After generating the configuration, review it and provide:

1. **Validation Status**: Does the config pass schema validation?
2. **Completeness Check**: Are there obvious missing checks for detected tools?
3. **Best Practice Suggestions**: Improvements based on common patterns
4. **Potential Issues**: Warnings about checks that might fail or be too strict

Format your review as:

### Configuration Review

**Status:** Valid / Invalid
**Checks Generated:** N checks covering [tools]

**Suggestions:**
- Consider adding `npm audit` for dependency vulnerability scanning
- The coverage threshold of 70% is reasonable, but 80% is common for mature projects
- You might want to add a `build` check to catch compilation errors early

**Potential Issues:**
- The `go_lint` check may fail if golangci-lint is not installed
- Consider adding `--fix` to the formatter check for auto-correction
```

**Benefits:**
- Catches configuration gaps the user might miss
- Educates users about best practices
- Provides actionable next steps
- Creates natural conversation for refinement

---

## 8. Rollout Strategy

### Phase 1: Internal Testing
- Develop and test with vibeguard project
- Have team review prompt quality
- Validate on 3-5 diverse projects
- Fix issues before public release

### Phase 2: Alpha Release
- Include in next minor version (v0.N.0)
- Add feature flag if needed: `--assist`
- Document in README
- Solicit feedback from early adopters

### Phase 3: Feedback Loop
- Collect data on generated config success
- Monitor for pattern/prompt issues
- Refine based on real-world usage
- Document common patterns and anti-patterns

### Phase 4: General Availability
- Remove feature flag (if used)
- Promote in documentation
- Add cookbook examples
- Consider agent integration guides

---

## 9. Risk Mitigation

### Risk: Prompt Too Complex for Agents
**Mitigation:**
- Keep prompt <4000 tokens
- Use clear, simple language
- Provide many concrete examples
- Test early with Claude Code

### Risk: Detection Inaccuracy
**Mitigation:**
- Implement conservative confidence scoring
- Show warnings for ambiguous cases
- Ask clarifying questions when uncertain
- Provide fallback prompts for unknown projects

### Risk: Generated Configs Fail Validation
**Mitigation:**
- Include validation rules in prompt
- Provide explicit constraints
- Test with actual agents in alpha
- Iterate on examples based on failures

### Risk: Performance Issues on Large Codebases
**Mitigation:**
- Profile inspector early
- Implement smart caching for tool detection
- Set reasonable file scan limits
- Consider staged detection (fast → thorough)

### Risk: Users Ignore Generated Config Quality
**Mitigation:**
- Emphasize "AI-assisted" not "AI-generated"
- Require user review and approval
- Document review checklist
- Provide clear validation instructions

---

## 10. Success Criteria & Metrics

### Launch Criteria
- [ ] Repository inspector accuracy >90% on test projects
- [ ] Prompt generation <1 second on typical projects
- [ ] Generated configs pass validation 95%+ of time
- [ ] Prompt clarity validated by team review
- [ ] Comprehensive error messages for failures
- [ ] Documentation and examples complete

### Post-Launch Metrics
- **Usage:** Number of `init --assist` invocations
- **Success Rate:** % of generated configs that pass validation
- **User Feedback:** Quality scores from adopters
- **Follow-ups:** Average questions needed to get valid config
- **Adoption:** % of new users using --assist vs manual config

### Performance Targets
- Inspector execution: <500ms
- Prompt composition: <200ms
- Total end-to-end: <1s
- Memory usage: <100MB on large projects

---

## 11. Related Architecture Decisions

- **ADR-003:** Go enables fast, frictionless CLI tool development
- **ADR-005:** VibeGuard dogfooding validates own features
- **ADR-006:** Git hooks integrate VibeGuard into developer workflows
- **CLAUDE.md:** Project conventions for agent integration and tooling

---

## 12. Dependencies & Resources

### External Tools
- Grok library (for pattern validation/testing)
- Standard Go libraries (filepath, os, io)
- YAML parser (already in use)

### Documentation to Create
- README section: "AI-Assisted Setup"
- User guide: "Getting Started with --assist"
- Developer guide: "How the Setup Prompt Works"
- Example prompts: 3-5 sample outputs for review

---

## Conclusion

This specification provides a complete roadmap for implementing the AI agent-assisted setup feature. The design prioritizes:
- **Clarity** - Unambiguous prompt that guides agents effectively
- **Accuracy** - Smart project detection and tool recognition
- **Simplicity** - Pure text-based interface, no special integrations
- **Reliability** - Validation rules and error handling throughout

The implementation can proceed phase-by-phase with clear success criteria at each stage, enabling rapid iteration and continuous validation with real agents.
