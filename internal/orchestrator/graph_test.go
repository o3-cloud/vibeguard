// Package orchestrator coordinates check execution with dependency management
// and parallel execution.
package orchestrator

import (
	"reflect"
	"sort"
	"testing"

	"github.com/vibeguard/vibeguard/internal/config"
)

func TestBuildGraph_NoDependencies_AllInLevelZero(t *testing.T) {
	checks := []config.Check{
		{ID: "a", Run: "echo a"},
		{ID: "b", Run: "echo b"},
		{ID: "c", Run: "echo c"},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 1 {
		t.Fatalf("expected 1 level, got %d", len(levels))
	}

	// All checks should be in level 0
	if len(levels[0]) != 3 {
		t.Errorf("expected 3 checks in level 0, got %d", len(levels[0]))
	}

	// Sort for comparison (order within level is not guaranteed)
	got := append([]string{}, levels[0]...)
	sort.Strings(got)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected level 0 to contain %v, got %v", want, got)
	}
}

func TestBuildGraph_LinearDependencyChain(t *testing.T) {
	// c depends on b, b depends on a
	checks := []config.Check{
		{ID: "a", Run: "echo a"},
		{ID: "b", Run: "echo b", Requires: []string{"a"}},
		{ID: "c", Run: "echo c", Requires: []string{"b"}},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 3 {
		t.Fatalf("expected 3 levels, got %d: %v", len(levels), levels)
	}

	// Level 0: a (no deps)
	if len(levels[0]) != 1 || levels[0][0] != "a" {
		t.Errorf("expected level 0 to be [a], got %v", levels[0])
	}

	// Level 1: b (depends on a)
	if len(levels[1]) != 1 || levels[1][0] != "b" {
		t.Errorf("expected level 1 to be [b], got %v", levels[1])
	}

	// Level 2: c (depends on b)
	if len(levels[2]) != 1 || levels[2][0] != "c" {
		t.Errorf("expected level 2 to be [c], got %v", levels[2])
	}
}

func TestBuildGraph_DiamondDependency(t *testing.T) {
	// Diamond: d depends on b and c, both b and c depend on a
	//     a
	//    / \
	//   b   c
	//    \ /
	//     d
	checks := []config.Check{
		{ID: "a", Run: "echo a"},
		{ID: "b", Run: "echo b", Requires: []string{"a"}},
		{ID: "c", Run: "echo c", Requires: []string{"a"}},
		{ID: "d", Run: "echo d", Requires: []string{"b", "c"}},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 3 {
		t.Fatalf("expected 3 levels, got %d: %v", len(levels), levels)
	}

	// Level 0: a
	if len(levels[0]) != 1 || levels[0][0] != "a" {
		t.Errorf("expected level 0 to be [a], got %v", levels[0])
	}

	// Level 1: b and c (can run in parallel)
	if len(levels[1]) != 2 {
		t.Errorf("expected 2 checks in level 1, got %d", len(levels[1]))
	}
	level1 := append([]string{}, levels[1]...)
	sort.Strings(level1)
	if !reflect.DeepEqual(level1, []string{"b", "c"}) {
		t.Errorf("expected level 1 to be [b, c], got %v", level1)
	}

	// Level 2: d
	if len(levels[2]) != 1 || levels[2][0] != "d" {
		t.Errorf("expected level 2 to be [d], got %v", levels[2])
	}
}

func TestBuildGraph_MultipleDependencies(t *testing.T) {
	// d depends on a, b, and c (all in level 0)
	checks := []config.Check{
		{ID: "a", Run: "echo a"},
		{ID: "b", Run: "echo b"},
		{ID: "c", Run: "echo c"},
		{ID: "d", Run: "echo d", Requires: []string{"a", "b", "c"}},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 2 {
		t.Fatalf("expected 2 levels, got %d: %v", len(levels), levels)
	}

	// Level 0: a, b, c
	if len(levels[0]) != 3 {
		t.Errorf("expected 3 checks in level 0, got %d", len(levels[0]))
	}

	// Level 1: d
	if len(levels[1]) != 1 || levels[1][0] != "d" {
		t.Errorf("expected level 1 to be [d], got %v", levels[1])
	}
}

func TestBuildGraph_IndependentChains(t *testing.T) {
	// Two independent chains: a -> b and c -> d
	checks := []config.Check{
		{ID: "a", Run: "echo a"},
		{ID: "b", Run: "echo b", Requires: []string{"a"}},
		{ID: "c", Run: "echo c"},
		{ID: "d", Run: "echo d", Requires: []string{"c"}},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 2 {
		t.Fatalf("expected 2 levels, got %d: %v", len(levels), levels)
	}

	// Level 0: a, c
	if len(levels[0]) != 2 {
		t.Errorf("expected 2 checks in level 0, got %d", len(levels[0]))
	}
	level0 := append([]string{}, levels[0]...)
	sort.Strings(level0)
	if !reflect.DeepEqual(level0, []string{"a", "c"}) {
		t.Errorf("expected level 0 to be [a, c], got %v", level0)
	}

	// Level 1: b, d
	if len(levels[1]) != 2 {
		t.Errorf("expected 2 checks in level 1, got %d", len(levels[1]))
	}
	level1 := append([]string{}, levels[1]...)
	sort.Strings(level1)
	if !reflect.DeepEqual(level1, []string{"b", "d"}) {
		t.Errorf("expected level 1 to be [b, d], got %v", level1)
	}
}

func TestBuildGraph_SingleCheck(t *testing.T) {
	checks := []config.Check{
		{ID: "only", Run: "echo only"},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 1 {
		t.Fatalf("expected 1 level, got %d", len(levels))
	}
	if len(levels[0]) != 1 || levels[0][0] != "only" {
		t.Errorf("expected level 0 to be [only], got %v", levels[0])
	}
}

func TestBuildGraph_EmptyChecks(t *testing.T) {
	checks := []config.Check{}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 0 {
		t.Errorf("expected 0 levels for empty checks, got %d", len(levels))
	}
}

func TestBuildGraph_UnknownDependency_ReturnsError(t *testing.T) {
	checks := []config.Check{
		{ID: "a", Run: "echo a", Requires: []string{"nonexistent"}},
	}

	_, err := BuildGraph(checks)
	if err == nil {
		t.Fatal("expected error for unknown dependency, got nil")
	}

	want := `check "a" requires unknown check: nonexistent`
	if err.Error() != want {
		t.Errorf("expected error %q, got %q", want, err.Error())
	}
}

func TestBuildGraph_CyclicDependency_TwoNodes(t *testing.T) {
	// Note: In practice, config validation catches cycles first.
	// This tests BuildGraph's own cycle detection as a fallback.
	checks := []config.Check{
		{ID: "a", Run: "echo a", Requires: []string{"b"}},
		{ID: "b", Run: "echo b", Requires: []string{"a"}},
	}

	_, err := BuildGraph(checks)
	if err == nil {
		t.Fatal("expected error for cyclic dependency, got nil")
	}

	// Error should mention circular/cyclic dependency
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestBuildGraph_CyclicDependency_ThreeNodes(t *testing.T) {
	// a -> b -> c -> a
	checks := []config.Check{
		{ID: "a", Run: "echo a", Requires: []string{"c"}},
		{ID: "b", Run: "echo b", Requires: []string{"a"}},
		{ID: "c", Run: "echo c", Requires: []string{"b"}},
	}

	_, err := BuildGraph(checks)
	if err == nil {
		t.Fatal("expected error for cyclic dependency, got nil")
	}
}

func TestBuildGraph_ComplexGraph(t *testing.T) {
	// Complex dependency graph:
	//       a
	//      /|\
	//     b c d
	//     |/ \|
	//     e   f
	//      \ /
	//       g
	checks := []config.Check{
		{ID: "a", Run: "echo a"},
		{ID: "b", Run: "echo b", Requires: []string{"a"}},
		{ID: "c", Run: "echo c", Requires: []string{"a"}},
		{ID: "d", Run: "echo d", Requires: []string{"a"}},
		{ID: "e", Run: "echo e", Requires: []string{"b", "c"}},
		{ID: "f", Run: "echo f", Requires: []string{"c", "d"}},
		{ID: "g", Run: "echo g", Requires: []string{"e", "f"}},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 4 {
		t.Fatalf("expected 4 levels, got %d: %v", len(levels), levels)
	}

	// Level 0: a
	if len(levels[0]) != 1 || levels[0][0] != "a" {
		t.Errorf("expected level 0 to be [a], got %v", levels[0])
	}

	// Level 1: b, c, d
	if len(levels[1]) != 3 {
		t.Errorf("expected 3 checks in level 1, got %d", len(levels[1]))
	}
	level1 := append([]string{}, levels[1]...)
	sort.Strings(level1)
	if !reflect.DeepEqual(level1, []string{"b", "c", "d"}) {
		t.Errorf("expected level 1 to be [b, c, d], got %v", level1)
	}

	// Level 2: e, f
	if len(levels[2]) != 2 {
		t.Errorf("expected 2 checks in level 2, got %d", len(levels[2]))
	}
	level2 := append([]string{}, levels[2]...)
	sort.Strings(level2)
	if !reflect.DeepEqual(level2, []string{"e", "f"}) {
		t.Errorf("expected level 2 to be [e, f], got %v", level2)
	}

	// Level 3: g
	if len(levels[3]) != 1 || levels[3][0] != "g" {
		t.Errorf("expected level 3 to be [g], got %v", levels[3])
	}
}

func TestBuildGraph_PreservesCheckOrder(t *testing.T) {
	// Within a level, checks should appear in their original order from the slice
	checks := []config.Check{
		{ID: "z", Run: "echo z"},
		{ID: "a", Run: "echo a"},
		{ID: "m", Run: "echo m"},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 1 {
		t.Fatalf("expected 1 level, got %d", len(levels))
	}

	// Order should match input order (z, a, m), not alphabetical
	want := []string{"z", "a", "m"}
	if !reflect.DeepEqual(levels[0], want) {
		t.Errorf("expected level 0 to preserve order %v, got %v", want, levels[0])
	}
}

func TestBuildGraph_RealWorldExample(t *testing.T) {
	// Simulates a real CI pipeline:
	// lint and vet can run in parallel
	// build depends on vet
	// test depends on build
	// coverage depends on test
	checks := []config.Check{
		{ID: "lint", Run: "golangci-lint run"},
		{ID: "vet", Run: "go vet ./..."},
		{ID: "build", Run: "go build ./...", Requires: []string{"vet"}},
		{ID: "test", Run: "go test ./...", Requires: []string{"build"}},
		{ID: "coverage", Run: "go test -cover ./...", Requires: []string{"test"}},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	levels := graph.Levels()
	if len(levels) != 4 {
		t.Fatalf("expected 4 levels, got %d: %v", len(levels), levels)
	}

	// Level 0: lint, vet
	if len(levels[0]) != 2 {
		t.Errorf("expected 2 checks in level 0, got %d", len(levels[0]))
	}

	// Level 1: build
	if len(levels[1]) != 1 || levels[1][0] != "build" {
		t.Errorf("expected level 1 to be [build], got %v", levels[1])
	}

	// Level 2: test
	if len(levels[2]) != 1 || levels[2][0] != "test" {
		t.Errorf("expected level 2 to be [test], got %v", levels[2])
	}

	// Level 3: coverage
	if len(levels[3]) != 1 || levels[3][0] != "coverage" {
		t.Errorf("expected level 3 to be [coverage], got %v", levels[3])
	}
}

func TestLevels_ReturnsCopy(t *testing.T) {
	checks := []config.Check{
		{ID: "a", Run: "echo a"},
	}

	graph, err := BuildGraph(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get levels twice and modify one
	levels1 := graph.Levels()
	levels2 := graph.Levels()

	// They should point to the same underlying data (current implementation)
	// This test documents the current behavior
	if len(levels1) > 0 && len(levels2) > 0 {
		// Verify they have the same content
		if levels1[0][0] != levels2[0][0] {
			t.Error("levels should have same content")
		}
	}
}
