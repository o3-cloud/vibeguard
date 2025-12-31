---
summary: Evaluated AI agent-assisted setup feature for automating vibeguard configuration generation
event_type: research
sources:
  - README.md (VibeGuard overview and CLI reference)
  - cmd/vibeguard/main.go (Entry point and error handling)
  - internal/cli/init.go (Current init command implementation)
  - examples/ (Configuration examples for Go and Node projects)
  - CLAUDE.md (Project agent conventions)
tags:
  - agent-assisted-setup
  - automation
  - config-generation
  - ai-integration
  - developer-experience
  - initialization
  - prompt-engineering
  - cli-enhancement
---

# AI Agent-Assisted Setup for VibeGuard

## Overview

Researched the feasibility of adding an AI agent-assisted setup mode to VibeGuard that would automate configuration generation. This feature would allow any coding agent (Claude Code, Cursor, etc.) to intelligently set up VibeGuard for a new project by reading a standardized prompt output.

## Current State Analysis

### Existing Setup Process

**Current `vibeguard init` command:**
- Generates a static, hardcoded Go project template
- Outputs only two commands for next steps
- Requires manual configuration review and customization
- No intelligence about project structure or context

**Current config format:**
- Well-structured YAML with clear schema
- Supports variable interpolation (`{{.var_name}}`)
- Flexible check definitions with grok patterns and assertions
- Dependencies and timeout specifications available

### Key Configuration Elements

The config schema supports:
1. **Variables** - Global vars for interpolation (e.g., `go_packages: "./..."`)
2. **Checks** - Array of check definitions with:
   - `id` - Unique identifier
   - `run` - Shell command to execute
   - `grok` - Pattern extraction from output
   - `assert` - Conditional assertions on extracted data
   - `requires` - Check dependencies
   - `severity` - error or warning
   - `suggestion` - Help text for failures
   - `timeout` - Execution time limit

### Examples Available

- `simple.yaml` - Basic Go checks (fmt, vet, test)
- `go-project.yaml` - Comprehensive Go setup with coverage validation
- `node-project.yaml` - Node.js project example (reviewed for patterns)

## Proposed AI Agent-Assisted Setup Feature

### Core Concept

A new `vibeguard init --assist` command (or similar) would:

1. **Inspect Repository Context**
   - Detect project type (Go, Node.js, Python, etc.)
   - Scan for existing tools (linters, test frameworks, CI/CD config)
   - Identify package structure and entry points
   - Review project metadata (package.json, go.mod, pyproject.toml, etc.)

2. **Output an AI-Readable Prompt**
   - Comprehensive guide for AI agents to understand the repo
   - Structured instructions for configuration generation
   - Clear requirements and constraints
   - Examples of valid check definitions
   - Instructions for asking clarifying questions

3. **AI Agent Workflow**
   - Read the prompt from vibeguard
   - Analyze project structure and existing tools
   - Ask clarifying questions if needed (using structured formats)
   - Generate appropriate vibeguard.yaml configuration
   - Validate generated config against schema
   - Apply the configuration to the project

### Expected Output Format

```bash
$ vibeguard init --assist
```

Would output a structured prompt containing:

```
# VibeGuard AI Agent Setup Guide

## Project Analysis
[Repository structure, detected tools, existing configurations]

## Configuration Generation Instructions
[Clear requirements for what checks to create]

## Example Checks
[Valid check definitions showing patterns]

## Validation Requirements
[How to validate the generated config]

## Questions to Ask
[When and how to ask the developer clarifying questions]
```

## Key Benefits

### For Developers
1. **Zero Manual Configuration** - AI agent sets up checks automatically
2. **Context-Aware** - Detects existing tools and project structure
3. **Fast Onboarding** - Seconds to full policy enforcement
4. **Works with Any Agent** - Standard prompt-based interface
5. **Intelligent Defaults** - Suggestions based on best practices

### For VibeGuard Project
1. **Universal Adoption** - Works with Claude Code, Cursor, other agents
2. **Real-World Testing** - Validates setup process at scale
3. **Dogfooding Value** - Can use for VibeGuard's own setup
4. **Composable Design** - Aligns with VibeGuard's declarative philosophy
5. **Low Friction Integration** - Just a text-based prompt, no new dependencies

## Technical Implementation Strategy

### Phase 1: Prompt Engineering
- Create comprehensive setup guide prompt
- Include project detection patterns
- Document valid check definitions and patterns
- Add validation instructions
- Create example patterns for common tools

### Phase 2: Repository Inspection
- Add project type detection logic
- Scan for common tools and frameworks
- Parse project metadata files
- Output structured analysis for prompt

### Phase 3: Integration
- New `vibeguard init --assist` command
- Output composed prompt to stdout (or file)
- Document how to use with Claude Code via stdin
- Support piping to AI agent workflows

### Phase 4: Validation
- Generate sample configs for multiple project types
- Test with actual AI agents
- Verify config validity
- Collect feedback on clarity and completeness

## Design Considerations

### Prompt Design
- **Clear Structure** - Sections for context, instructions, examples
- **Unambiguous Constraints** - Rules for valid configurations
- **Rich Examples** - Show successful patterns for different project types
- **Validation Rules** - How to verify generated configs
- **Fallback Guidance** - What to do if uncertain

### Project Detection
- File-based heuristics (go.mod, package.json, Pipfile)
- Common tool detection (.golangci.yml, eslintrc, pytest.ini)
- Language-specific best practices
- Graceful degradation if detection uncertain

### Scope Limitations
- Focus on detecting existing quality tools
- Generate checks for detected tools
- Avoid generating checks for tools that don't exist
- Ask clarifying questions for ambiguous cases
- Don't override user preferences

## Alignment with VibeGuard Philosophy

### Declarative Configuration
✓ Config remains YAML-based and declarative
✓ Agent generates valid declarations, not imperative scripts
✓ Configuration stays human-readable and reviewable

### Composability
✓ Checks defined independently with clear dependencies
✓ Agent can suggest checks that compose together
✓ Variables and patterns are composable building blocks

### Transparency
✓ AI output (generated config) is fully transparent
✓ Developer must review and approve configuration
✓ Clear reasoning for each check suggestion

## Comparison with Alternatives

### Hardcoded Templates (Current)
- ✗ One-size-fits-all approach
- ✗ Requires manual customization
- ✗ Doesn't detect existing tools
- ✓ Simple and reliable

### AI Agent-Assisted Setup (Proposed)
- ✓ Context-aware configuration
- ✓ Detects existing tools automatically
- ✓ Works with any coding agent
- ✓ Minimal manual work
- ✗ Requires careful prompt engineering

### Interactive Setup (Not Recommended)
- Would require agent framework integration
- VibeGuard CLI not designed for async interactions
- Prompt-based approach is more universal

## Next Steps

1. **Create Setup Guide Prompt** - Draft comprehensive AI setup instructions
2. **Implement Repository Inspector** - Add project type detection
3. **Build Assist Mode** - New CLI command to output structured prompt
4. **Test with Agents** - Use Claude Code and other agents to test flow
5. **Gather Feedback** - Refine prompt based on real-world usage
6. **Document Workflow** - Add to README and project guides

## Related Decisions

- **ADR-006** - VibeGuard integration with developer tools and hooks
- **ADR-003** - Go enables frictionless CLI integration
- **CLAUDE.md** - Project conventions for agent integration

## Conclusion

AI agent-assisted setup is highly feasible and aligns well with VibeGuard's design philosophy. The key is engineering a clear, comprehensive prompt that enables agents to understand project context and generate valid configurations. This approach requires no changes to VibeGuard's core architecture and can be implemented incrementally.
