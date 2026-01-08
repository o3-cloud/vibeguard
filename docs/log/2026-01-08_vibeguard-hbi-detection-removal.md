---
summary: Task vibeguard-hbi - Removed detection pipeline from runAssist(), delegating to AI agents
event_type: code
sources:
  - docs/specs/init-template-system-spec.md
  - internal/cli/init.go
  - internal/cli/inspector/prompt.go
tags:
  - refactoring
  - init-template-system
  - ai-agents
  - simplification
  - detection-pipeline
  - runAssist
  - code-reduction
---

# Task vibeguard-hbi: Remove Detection Pipeline from runAssist()

## Summary

Successfully completed vibeguard-hbi task by removing the detection pipeline from the `runAssist()` function. Project type detection is now delegated to AI agents, simplifying the initialization flow and reducing code complexity by ~65%.

## Changes Made

### 1. Simplified runAssist() in internal/cli/init.go

Removed ~50 lines of detection code:
- Removed `inspector.NewDetector()` call and project type detection logic
- Removed confidence checking that blocked undetectable projects
- Removed `inspector.NewToolScanner()` tool scanning
- Removed `inspector.NewMetadataExtractor()` metadata and structure extraction
- Removed `inspector.NewRecommender()` and check recommendation generation

The function now calls a new simplified helper instead:
```go
prompt, err := inspector.GenerateSetupPromptWithoutDetection(root)
```

### 2. Created GenerateSetupPromptWithoutDetection() in inspector/prompt.go

New function generates a minimal setup prompt without project detection:
- Creates ProjectAnalysis with "unknown" project type
- Empty tools, source dirs, and other metadata arrays
- Delegates all detection responsibility to the AI agent
- Maintains consistent prompt generation interface

### 3. Updated Test in init_test.go

Modified `TestRunAssist_UndetectableProject`:
- Changed expectation from error to success
- Project detection is now agent responsibility, not a blocker
- Function works on any directory without type detection

## Code Impact

**Before:**
- runAssist(): 78 lines (including detection pipeline)
- 5 inspector package dependencies in init flow
- Hard requirement for detectable projects

**After:**
- runAssist(): 27 lines (65% reduction)
- Direct call to single helper function
- Works on any directory

## Alignment with Specification

This change implements Phase 1 of the Init Template System spec (vibeguard-hbi):
- "Project type detection should no longer be required"
- Agents analyze projects and select appropriate templates
- Template discovery instruction in assist output guides agents

## Verification

✅ Code compiles successfully
✅ All unit tests pass (init, assist, templates)
✅ vibeguard check passes all policy checks
✅ Commit: `0370c28` - refactor: Delegate project detection to AI agent in runAssist()
✅ Beads issue closed and synced

## Next Steps

The removal of detection enables:
- Phase 2: Template expansion (add new templates)
- Phase 3: Codebase simplification (remove duplicate detection logic from assist system)
- Cleaner architecture where agents handle intelligence, tool provides configuration options
