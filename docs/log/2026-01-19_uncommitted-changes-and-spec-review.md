---
summary: Comprehensive review of uncommitted changes and VibeGuard init template system specification
event_type: deep dive
sources:
  - internal/config/config.go
  - internal/config/schema.go
  - vibeguard.yaml
  - internal/cli/prompt.go
  - internal/cli/prompt_test.go
  - docs/specs/init-template-system-spec.md
tags:
  - specification
  - code-review
  - template-system
  - configuration
  - prompt-feature
  - uncommitted-changes
  - implementation-planning
  - architecture
---

# Review: Uncommitted Changes & Init Template System Specification

## Overview

This review documents the analysis of uncommitted changes in the vibeguard repository and the comprehensive specification for the VibeGuard init template system. The specification outlines a significant architectural simplification through templates-as-source-of-truth approach.

## Uncommitted Changes Summary

### Modified Files

#### 1. `internal/config/config.go`
**Status:** Modified (511 lines)
**Key Additions:**
- Enhanced error handling with `ConfigError` and `ExecutionError` types
- Line number tracking in YAML for precise error reporting
- Comprehensive validation for check dependencies and cycles (DFS-based cycle detection)
- Prompt validation support (`validatePrompts()`)
- YAML node preservation for line number lookup during validation
- Custom YAML unmarshaling for flexible `GrokSpec` (string or list) and `Duration` types

**Quality Assessment:** High - Well-structured with thoughtful error handling and line number context

#### 2. `internal/config/schema.go`
**Status:** Modified (66 lines)
**Key Additions:**
- New `Prompt` type supporting stored prompt templates (ID, description, content, tags)
- Enhanced `Config` struct to include prompts array
- Updated check struct documentation
- Duration and GrokSpec type definitions

**Quality Assessment:** Clean schema design with backward compatibility

#### 3. `vibeguard.yaml`
**Status:** Modified (150 lines)
**Key Additions:**
- New `prompts` section with 4 predefined prompts:
  - `init`: Guidance for initializing vibeguard configuration
  - `code-review`: System prompt for code review assistance
  - `security-audit`: Security-focused code analysis
  - `test-generator`: Comprehensive unit test generation
- Enhanced check definitions with timeout specifications and proper severity levels
- Better variable support with `{{.go_packages}}` interpolation

**Quality Assessment:** Practical and covers essential use cases for agent guidance

### New Files

#### 4. `internal/cli/prompt.go`
**Status:** New (implementation pending)
**Purpose:** Implements prompt feature for CLI
- Expected functionality for retrieving and using stored prompts
- Likely integrates with config prompts section

#### 5. `internal/cli/prompt_test.go`
**Status:** New (tests pending)
**Purpose:** Test coverage for prompt feature

### New Log Files
- `docs/log/2026-01-17_claude-code-hooks-integration.md` - Previous work on hooks
- `docs/log/2026-01-17_prompt-feature-implementation.md` - Feature planning
- `docs/log/2026-01-17_prompt-init-design-decision.md` - Design decisions
- `docs/log/2026-01-17_prompt-support-research.md` - Research findings

## Specification Document Analysis

### `docs/specs/init-template-system-spec.md` (626 lines)

**Title:** VibeGuard Init Template System - Master Specification
**Status:** Draft (v1.0, dated 2026-01-08)
**Scope:** Comprehensive template-driven initialization strategy

### Key Specification Components

#### 1. Core Concept & Three Modes
- **Mode 1:** AI-Assisted Setup (multi-step agent workflow)
- **Mode 2:** Direct Template Selection
- **Mode 3:** List Available Templates

#### 2. Command-Line Interface
**Existing flags (unchanged):**
- `--force/-f` - Overwrite existing config
- `--template/-t <name>` - Apply specific template
- `--assist` - Generate AI agent-assisted setup prompt
- `--output/-o <path>` - Output file for assist mode

**New flags:**
- `--list-templates` - Explicit template listing

#### 3. Template Architecture

**Template Registry Pattern:**
- Self-registering templates via `internal/cli/templates` package
- `Register()` function called in template `init()` functions
- Standardized `Template` struct (Name, Description, Content)

**Naming Convention:**
- `<language>-<variant>` (e.g., `go-minimal`, `go-standard`)
- `<language>-<framework>` (e.g., `node-express`, `python-django`)
- `<language>-<dialect>-<framework>` (e.g., `node-typescript`)

**Current Templates (8):**
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

**Planned Templates (Phase 2):**
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

#### 4. Implementation Phases

**Phase 1: Add --list-templates Flag and Assist Integration** (Immediate)
- Risk Level: Low
- Critical Changes:
  - Update error messages to reference `--list-templates`
  - Update flag help text
  - Add `TemplateDiscoverySection()` to assist prompts
  - Update assist composition in `composer.go`
- Backward compatibility maintained with `--template list`

**Phase 2: Template Expansion** (1-2 weeks after Phase 1)
- Risk Level: Very Low
- Add 9 new templates for popular frameworks
- No code changes required (additive only)

**Phase 3A: Quick Wins - Codebase Simplification** (Zero Breaking Changes)
- Risk Level: Minimal
- Delete `/examples/` directory (~200 lines)
- Delete `assist/templates.go` (~391 lines)
- Consolidate data structures (~50 lines)
- **Estimated Savings:** ~600 lines

**Phase 3B: Architectural Simplification** (Medium Risk)
- Risk Level: Medium
- Remove `starterConfig` constant (~35 lines)
- Simplify recommendations system (425→50 lines)
- Consolidate assist sections (~250+ lines reducible)
- **Estimated Savings:** ~900+ lines

**Phase 3C: Deep Refactoring** (Higher Risk)
- Risk Level: Higher
- Complete rewrite of `inspector/recommendations.go` (~525 lines)
- Redesign assist system around templates
- **Estimated Savings:** ~300+ lines

#### 5. Overall Impact (Phase 3)

| Module | Before | After | Reduction |
|--------|--------|-------|-----------|
| init.go | 300 | 265 | -12% |
| assist/ | 5,344 | 3,500 | -35% |
| inspector/ | 525 | 100 | -81% |
| examples/ | 200 | 0 | -100% |
| **TOTAL** | **7,603** | **5,099** | **-33%** |

**Total Code Reduction: 1,504 lines (-20%)**

#### 6. Success Criteria

**Phase 1:**
- `--list-templates` exits with code 0
- Assist prompt includes template discovery instruction
- Error messages reference `--list-templates`
- Backward compatibility maintained

**Phase 2:**
- Each template is valid YAML
- Templates cover intended language/framework
- Follow naming conventions

**Phase 3:**
- No existing functionality broken
- Code reduced by 1,500+ lines
- Single source of truth (templates)

#### 7. Risk Assessment Summary

- **Phase 1:** LOW risk - error messages and assist integration
- **Phase 2:** VERY LOW risk - additive only
- **Phase 3A:** MINIMAL risk - remove redundant files
- **Phase 3B:** MEDIUM risk - remove defaults and simplify systems
- **Phase 3C:** HIGHER risk - complete rewrites

## Key Findings

### 1. Specification Quality
✅ **Strengths:**
- Comprehensive and well-structured
- Clear phase breakdown with risk assessment
- Backward compatibility preserved
- Detailed success criteria

✅ **Well-thought out:**
- Templates as single source of truth is sound architecture
- Simplification estimates well-documented
- Agent integration workflow clearly defined

### 2. Uncommitted Changes Assessment
✅ **Prompt Feature:**
- Adds valuable guided-setup capability
- Integrates well with existing configuration system
- Extends config schema with backward compatibility

⚠️ **Status:**
- `prompt.go` and `prompt_test.go` are placeholders
- Implementation work needed for full functionality

### 3. Configuration System Enhancement
✅ **Improvements:**
- Enhanced error handling with line numbers
- Cycle detection for dependency validation
- Flexible YAML unmarshaling

### 4. Recommendation
- **Phase 1 should start immediately** - low risk, high value for agent workflows
- **Phase 2 can follow quickly** - additive only
- **Phase 3 requires careful planning** - significant refactoring

## Next Steps

1. **Complete prompt.go implementation** with prompt retrieval and usage
2. **Implement Phase 1 changes** to error messages and assist integration
3. **Add template discovery section** to assist system
4. **Update documentation** to reflect new `--list-templates` workflow
5. **Plan Phase 2 template expansion** with timeline
6. **Evaluate Phase 3 before proceeding** - requires team discussion on refactoring commitment

## Related Decisions

- ADR-003: Adopt Go as the Primary Implementation Language
- ADR-004: Establish Code Quality Standards and Tooling
- ADR-005: Adopt VibeGuard for Policy Enforcement in CI/CD

## Rate Assessment

**Specification Quality:** ⭐⭐⭐⭐⭐ (5/5)
- Comprehensive, well-structured, considers risks and benefits
- Clear phased approach with realistic scope management
- Strong alignment with project goals

**Implementation Readiness:** ⭐⭐⭐⭐ (4/5)
- Specification well-prepared for Phase 1 implementation
- Placeholder files exist for prompt feature
- Phase 1 changes are straightforward and low-risk

**Overall Assessment:** Excellent specification document supporting significant architectural improvement. Ready to proceed with Phase 1 implementation.
