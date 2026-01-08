---
summary: Critical spec review of VibeGuard Init Template System v1.0 identified 5 blockers and 5 major issues requiring revisions before implementation
event_type: meeting
sources:
  - docs/specs/init-template-system-spec.md
  - internal/cli/init.go
  - internal/cli/assist/sections.go
  - internal/cli/inspector/recommendations.go
tags:
  - spec-review
  - vibeguard
  - init-template-system
  - architecture
  - blockers
  - implementation-planning
  - flag-validation
  - code-quality
---

# Spec Review: VibeGuard Init Template System (v1.0)

**Status:** NOT READY FOR APPROVAL

**Recommendation:** Return spec for revisions addressing all identified issues before implementation begins (estimated 2-3 days of work).

## Executive Summary

The Init Template System specification is comprehensive and addresses a real problem—simplifying VibeGuard setup and reducing codebase duplication. However, critical gaps between specification promises and actual implementation, combined with architectural issues and mathematically impossible reduction targets, prevent approval in its current form.

The spec requires substantial technical refinement before Phase 1 can begin.

## Five Critical Blockers

### 1. Missing Flag Conflict Validation

**Problem:** Specification claims (lines 471-476) that conflicting flags like `--assist + --template` should error, but `init.go` has no implementation for this validation.

**Current behavior:** Code silently ignores one flag or the other (lines 97-104 in init.go).

**Required fix:** Add explicit mutually-exclusive flag validation:
```go
flagsSet := 0
if initAssist { flagsSet++ }
if initListTemplates { flagsSet++ }
if initTemplate != "" { flagsSet++ }

if flagsSet > 1 {
    return fmt.Errorf("--assist, --list-templates, and --template are mutually exclusive")
}
```

**Impact:** Phase 1 cannot complete without this implementation.

### 2. Self-Registering init() Pattern Creates Non-Deterministic Behavior

**Problem:** Specification proposes using Go's `init()` functions (lines 123-133) for template registration across packages.

**Issues:**
- Go's `init()` execution order is undefined when there are no import dependencies
- Cannot reset global registry between tests (test pollution)
- Violates Go best practices (import side effects)
- Creates non-deterministic behavior in tests and production

**Required fix:** Replace with explicit registration pattern:
```go
func RegisterAllTemplates() {
    Register(goStandardTemplate())
    Register(goMinimalTemplate())
    // ... all templates registered in deterministic order
}
```

Called once from `init.go` initialization.

**Impact:** Current pattern will cause intermittent test failures and hard-to-debug issues.

### 3. Phase 3 Line Reduction Math Doesn't Add Up

**Problem:** Specification claims total 1,504 line reduction (line 435), including 1,844 lines from assist module (line 430), but actual identified deletions total only ~400 lines.

**Evidence:**
- `assist/templates.go`: 391 lines (spec acknowledges this)
- `examples/`: ~200 lines
- **Total identified: ~600 lines**
- **Claimed from assist alone: 1,844 lines**
- **Gap: ~1,244 unaccounted lines**

**Impact:** Either the reduction targets are unrealistic, or the scope is underspecified. Blocks Phase 3 planning.

### 4. Circular Dependency in Phase 3: Templates vs. Recommendations

**Problem:** Phase 3 proposes eliminating `inspector/recommendations.go` (525 lines) to simplify codebase, but the recommendation engine IS how agents select appropriate templates.

**Contradiction:**
- Template system requires agents to analyze projects and choose correct templates
- This analysis and matching requires the recommendation/detection logic
- Deleting recommendations.go breaks `vibeguard init --assist` functionality
- Without `--assist`, agents cannot intelligently discover which template to use

**Current flow:** `vibeguard init --assist` → recommendations engine detects tools → suggests checks → agent converts to template

**Proposed broken flow:** `vibeguard init --assist` → ??? → recommend template → agent selects template

**Impact:** Phase 3 cannot be executed as specified without losing core `--assist` functionality.

### 5. Phase 3A Claims "Zero Breaking Changes" While Making Breaking Changes

**Problem:** Line 376 claims Phase 3A has "Zero Breaking Changes," but the phase includes:

1. **Deleting examples/** - Breaks all documentation links and tutorials referencing example configs
2. **Changing LanguageExamplesSection()** - Modifies assist output structure, affecting agents trained on current format
3. **Removing starterConfig** - Changes default behavior or requires explicit template selection

**Reality:** These ARE breaking changes for documentation and agent workflows.

**Impact:** Misleads stakeholders about migration effort required. Requires honest assessment of breaking changes and migration path.

## Major Issues (Should Address Before Phase 1)

### Template Naming Rules Undefined for Complex Scenarios

Lines 136-165 define naming convention but have gaps:

- **Ambiguity:** Is `node-typescript-react-vite` the right order, or `node-react-typescript-vite`?
- **No precedence rules** for multiple modifiers
- **Real-world impact:** Agents cannot reliably select correct template without clear naming semantics

### Monorepo Support Excluded

Line 585 marks monorepo strategy as "Out of Scope," but this is a significant limitation:

- Modern projects often use Go backend + React frontend + Python scripts
- No path for initializing multi-language monorepos
- Limits real-world adoption significantly

### Test Coverage Insufficient

Phase 1 success criteria (lines 449-455) missing critical tests:

- No integration test: `--assist` → `--list-templates` → `--template X` → `vibeguard check`
- No format validation for `--list-templates` output
- No test verifying assist prompt includes template discovery instruction
- No backward compatibility test for `--template list`

### Risk Mitigations Too Vague

Line 493: "Mitigation: Update documentation and agent guides before release"

**Problem:** Doesn't specify:
- Which README sections?
- Which blog posts or tutorials?
- How to validate changes?
- Who is responsible?

Line 517: "Mitigation: Update all documentation references before deletion"

**Better approach:** "Run `git grep 'examples/.*\.yaml'` and systematically update all matches in docs/ and README files."

Line 533: "Mitigation: A/B test prompt changes with agents before full rollout"

**Problem:** Vaporware without infrastructure:
- What are baseline metrics?
- How do you run parallel agent tests?
- What % regression is acceptable?
- Who monitors results?

### Missing Frameworks in Phase 2

Phase 2 (lines 169-178) is incomplete. Notable gaps:

- **Vue.js** - 3rd most popular JS framework (missing `node-vue-vite`)
- **Svelte** - Growing ecosystem (`node-svelte-kit`)
- **Ruby** - Has examples in assist/templates.go but no templates
- **Java** - Has examples but no templates

## What Spec Gets Right

- **Clear three-mode architecture** (AI-assisted, direct selection, discovery)
- **Solid template structure** (YAML format with clear sections)
- **Good naming convention foundation** (language-variant pattern)
- **Realistic Phase 1 scope** (small, focused changes)
- **Comprehensive documentation** of long-term vision

## What Needs To Happen Before Approval

1. **Add explicit flag conflict validation** to Phase 1 with code examples
2. **Replace self-registering init()** with explicit `RegisterAllTemplates()` function
3. **Recount Phase 3 line reductions** with actual file sizes and realistic targets
4. **Resolve circular dependency** between templates and recommendations
5. **Remove "zero breaking changes" claim** and document actual breaking changes + migration path
6. **Define template naming rules** for complex scenarios (precedence, ordering)
7. **Add integration tests** to Phase 1 success criteria
8. **Expand frameworks list** for Phase 2 (Vue, Svelte, Ruby, Java)
9. **Make risk mitigations concrete** with specific steps, not vague intentions
10. **Address monorepo question** - defer or include in scope?

## Estimated Revision Effort

- **Technical writing:** 1-2 days (recount, clarify, add missing sections)
- **Architecture refinement:** 1 day (resolve circular dependency, flag validation)
- **Testing & validation:** 0.5 days (add integration test specs)

**Total:** 2-3 days before resubmission for approval.

## Next Steps

Option A: Return spec to authors with detailed feedback on all 10 issues.

Option B: Schedule working session with authors to address blockers collaboratively.

Option C: Implement Phase 1 with known issues (NOT RECOMMENDED - will cause rework).

## Key Insight

The template system concept is sound and valuable. Templates as a single source of truth CAN reduce codebase complexity significantly. But the execution plan needs technical rigor and honest accounting of tradeoffs before team commitment.
