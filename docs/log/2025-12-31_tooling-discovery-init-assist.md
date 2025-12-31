---
summary: Added tooling inspection and research instructions to init --assist command output
event_type: code
sources:
  - internal/cli/assist/sections.go
  - internal/cli/assist/composer.go
  - internal/cli/assist/composer_test.go
  - .beads/vibeguard-43l.jsonl
tags:
  - init-assist
  - ai-assisted-setup
  - tooling-discovery
  - prompt-composition
  - feature
---

# Add Tooling Discovery Instructions to init --assist Output

## Summary

Implemented task `vibeguard-43l` to add tooling discovery instructions to the `init --assist` command output. This feature helps AI agents better understand and configure VibeGuard by:

1. **Inspecting existing tooling** - Agents are now instructed to analyze existing configuration files to understand how tools are configured
2. **Researching additional tooling** - Agents receive language-specific suggestions for quality/security tools they can recommend to users

## Changes Made

### New Sections Added

1. **ToolingInspectionSection** (`sections.go`)
   - Lists detected tools that have configuration files
   - Provides guidance on what to look for in each config file:
     - Enabled rules/checks
     - Disabled rules
     - Custom settings (paths, thresholds, exclusions)
     - Integration points with other tools
     - Command variations and flags

2. **ToolingResearchSection** (`sections.go`)
   - Provides language-specific tool suggestions based on project type
   - Includes suggestions for Go, Node.js, Python, Rust, Ruby, Java, and generic projects
   - Each suggestion includes:
     - Tool name
     - Category (Linter, Security, Code Quality, etc.)
     - Purpose
     - Value proposition
     - Example command
   - Guidance on how to present suggestions to users

### Tool Suggestions by Language

- **Go**: staticcheck, gosec, errcheck, ineffassign, govulncheck
- **Node.js**: tsc --noEmit, npm audit, depcheck, madge
- **Python**: bandit, safety, vulture, radon
- **Rust**: cargo audit, cargo deny, cargo outdated
- **Ruby**: brakeman, bundler-audit, reek
- **Java**: SpotBugs, OWASP Dependency-Check, PMD
- **Generic**: gitleaks, trivy

### Composer Updates

- Updated `buildSections()` to include the two new sections
- Updated `ComposeWithOptions()` to support selective inclusion
- Added `IncludeToolingInspection` and `IncludeToolingResearch` to `ComposerOptions`
- Updated `DefaultComposerOptions()` to include new sections
- Updated `MinimalComposerOptions()` to exclude new sections (keeping prompts small when needed)

### Tests Added

- `TestToolingInspectionSection` - Tests with and without config files
- `TestToolingResearchSection` - Tests all project types
- `TestGetToolSuggestions` - Validates suggestion data structure

## Technical Notes

- Prompt token estimate remains under 4000 tokens (~3811 tokens with new sections)
- All existing tests continue to pass
- The implementation follows the existing section-based composition pattern

## Next Steps

- Related task `vibeguard-esq` will add validation instructions to the output
- Consider adding more tool suggestions as the ecosystem evolves
