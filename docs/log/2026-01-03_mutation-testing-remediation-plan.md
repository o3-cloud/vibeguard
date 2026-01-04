---
summary: Implementation plan to kill 110 LIVED mutations and achieve 90%+ test efficacy
event_type: deep dive
sources:
  - docs/log/2026-01-03_mutation-testing-analysis-lived-mutations.md
  - mutations.txt
  - docs/adr/ADR-007-adopt-mutation-testing.md
  - docs/adr/ADR-004-code-quality-standards.md
tags:
  - mutation-testing
  - test-improvement
  - implementation-plan
  - gremlins
  - code-quality
  - boundary-testing
---

# Mutation Testing Remediation Plan

## Objective

Kill at least 45 of 110 LIVED mutations to raise test efficacy from **83%** to **90%+**.

## Prioritization Strategy

Files are prioritized by:
1. **Mutation density** - More LIVED mutations = higher priority
2. **Code criticality** - Core functionality over edge features
3. **Fix complexity** - Quick wins before complex refactors

---

## Phase 1: Quick Wins (Target: 25 mutations killed)

### 1.1 Parser Boundary Fix
**File:** `internal/assert/parser.go`
**Mutations:** 1 LIVED at line 37:29

```go
// Current: i < len(p.input)
// Mutation: i <= len(p.input) survives
```

**Action:** Add test case for `formatError` with:
- Empty input string
- Position at exact end of input
- Position beyond input length

**Test to add:**
```go
func TestFormatError_BoundaryConditions(t *testing.T) {
    tests := []struct {
        name  string
        input string
        pos   int
    }{
        {"empty input", "", 0},
        {"pos at end", "abc", 3},
        {"pos beyond end", "abc", 5},
        {"single char", "x", 1},
    }
    // ... verify pointer positioning
}
```

---

### 1.2 Composer Section Separator
**File:** `internal/cli/assist/composer.go`
**Mutations:** 4 LIVED at line 85

```go
// Current: if i < len(sections)-1
// Mutations: boundary, negation, arithmetic all survive
```

**Action:** Add tests for `assembleSections` with:
- Empty sections slice
- Single section (no separator needed)
- Two sections (exactly one separator)
- Verify separator placement precisely

**Test to add:**
```go
func TestAssembleSections_EdgeCases(t *testing.T) {
    tests := []struct {
        name     string
        sections []PromptSection
        wantSeps int // expected separator count
    }{
        {"empty", []PromptSection{}, 0},
        {"single", []PromptSection{{Content: "A"}}, 0},
        {"two", []PromptSection{{Content: "A"}, {Content: "B"}}, 1},
        {"three", []PromptSection{{Content: "A"}, {Content: "B"}, {Content: "C"}}, 2},
    }
    for _, tt := range tests {
        result := composer.assembleSections(tt.sections)
        gotSeps := strings.Count(result, "---")
        assert.Equal(t, tt.wantSeps, gotSeps)
    }
}
```

---

### 1.3 Check Command JSON Output
**File:** `internal/cli/check.go`
**Mutations:** 1 LIVED at line 77:55

**Action:** Add test that verifies error handling when JSON formatting fails. Mock or inject a writer that returns an error.

---

### 1.4 Sections Formatting Boundaries
**File:** `internal/cli/assist/sections.go`
**Mutations:** 8 LIVED at lines 42, 55, 58, 61, 91, 95, 104, 227, 234

**Action:** Add boundary tests for each section builder:
- Lines at exactly max length
- Zero items in collections
- Single item collections

---

## Phase 2: Inspector Package (Target: 35 mutations killed)

### 2.1 Detector File Traversal
**File:** `internal/cli/inspector/detector.go`
**Mutations:** 15+ LIVED

**Key mutations to target:**

| Line | Mutation | Test Needed |
|------|----------|-------------|
| 67:16 | BOUNDARY | `maxDepth=0` should return only root files |
| 130:23 | BOUNDARY | Single result limit test |
| 170:23 | BOUNDARY | Pattern match at boundary |
| 382:14 | NEGATION | Non-"." relative path handling |
| 386:11, 389:9 | INCREMENT | Depth counting accuracy |
| 393:15,29 | BOUNDARY/NEGATION | `maxDepth=-1` (unlimited) vs `maxDepth=0` |
| 418:21 | BOUNDARY/NEGATION | Result count at exact limit |

**Test file to create:** `detector_boundary_test.go`

```go
func TestDetector_DepthBoundaries(t *testing.T) {
    // Create temp directory structure:
    // root/
    //   file1.go       (depth 0)
    //   dir1/
    //     file2.go     (depth 1)
    //     dir2/
    //       file3.go   (depth 2)

    tests := []struct {
        maxDepth    int
        wantFiles   int
    }{
        {-1, 3}, // unlimited
        {0, 1},  // root only
        {1, 2},  // root + dir1
        {2, 3},  // all files
    }
}

func TestDetector_ResultLimits(t *testing.T) {
    // Test maxResults = 0, 1, exact count, beyond count
}
```

---

### 2.2 Metadata Version Comparisons
**File:** `internal/cli/inspector/metadata.go`
**Mutations:** 30+ LIVED

This file has the highest mutation density. Most are CONDITIONALS_BOUNDARY in version comparison logic.

**Focus areas:**
- Version string parsing edge cases
- Comparison operators at boundaries (e.g., `>=1.0.0` vs `>1.0.0`)
- Empty/nil version handling

**Test file to create:** `metadata_boundary_test.go`

```go
func TestVersionComparison_Boundaries(t *testing.T) {
    tests := []struct {
        constraint string
        version    string
        want       bool
    }{
        // Boundary: >= vs >
        {">=1.0.0", "1.0.0", true},
        {">1.0.0", "1.0.0", false},
        // Boundary: <= vs <
        {"<=2.0.0", "2.0.0", true},
        {"<2.0.0", "2.0.0", false},
        // Edge cases
        {">=0.0.0", "0.0.0", true},
        {"*", "0.0.1", true},
    }
}
```

---

### 2.3 Tools Detection Logic
**File:** `internal/cli/inspector/tools.go`
**Mutations:** 8+ LIVED at lines 281, 306, 567, 569, 600, 608, 636, 638, 640

**Focus areas:**
- Tool detection heuristics with edge case inputs
- String matching boundaries
- Result limit handling

---

## Phase 3: Config Package (Target: 10 mutations killed)

### 3.1 Cycle Detection
**File:** `internal/config/config.go`
**Mutations:** 8 LIVED at lines 267, 269, 331, 343, 344, 352, 678-680, 692, 697-700, 706

**Focus areas:**
- Cycle detection with single-node cycles
- Empty dependency graphs
- Boundary index calculations

**Test to add:**
```go
func TestCycleDetection_EdgeCases(t *testing.T) {
    tests := []struct {
        name   string
        checks []Check
        wantCycle bool
    }{
        {"self-reference", []Check{{ID: "a", DependsOn: []string{"a"}}}, true},
        {"no deps", []Check{{ID: "a"}}, false},
        {"single chain", []Check{{ID: "a"}, {ID: "b", DependsOn: []string{"a"}}}, false},
    }
}
```

---

## Phase 4: Spikes Package (Target: 8 mutations killed)

### 4.1 Spike Orchestrator/Config
**Files:** `spikes/config/config.go`, `spikes/orchestrator/orchestrator.go`
**Mutations:** 8+ LIVED

Lower priority as these are experimental spike implementations. Address after core packages.

---

## Verification Process

After each phase:

1. **Run targeted mutation tests:**
   ```bash
   gremlins run --tags "phase1" ./internal/assert/...
   ```

2. **Verify mutations killed:**
   ```bash
   grep "LIVED" mutations.txt | wc -l
   ```

3. **Update tracking:**
   - Record killed mutation count
   - Update efficacy percentage
   - Document any unexpected findings

---

## Success Criteria

| Metric | Current | Target |
|--------|---------|--------|
| Test Efficacy | 83.00% | 90%+ |
| LIVED Mutations | 110 | <65 |
| Critical Files Covered | Partial | Full boundary coverage |

---

## Estimated Effort

| Phase | Mutations Targeted | Complexity |
|-------|-------------------|------------|
| Phase 1 | 25 | Low - straightforward boundary tests |
| Phase 2 | 35 | Medium - requires test fixtures |
| Phase 3 | 10 | Medium - cycle detection edge cases |
| Phase 4 | 8 | Low - spike code, less critical |

---

## Related Decisions

- **ADR-007:** Establishes mutation testing with Gremlins
- **ADR-004:** Sets 70% coverage baseline; this plan exceeds that for mutation efficacy
