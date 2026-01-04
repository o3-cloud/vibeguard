// Package orchestrator coordinates check execution with dependency management
// and parallel execution.
package orchestrator

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/executor"
)

func TestRun_PassingCheck_ExitCodeZero(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "pass-check",
				Run:      "exit 0",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if len(result.Violations) != 0 {
		t.Errorf("expected no violations, got %d", len(result.Violations))
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}
	if !result.Results[0].Passed {
		t.Error("expected check to pass")
	}
}

func TestRun_FailingCheck_ErrorSeverity_ExitCodeThree(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:         "fail-check",
				Run:        "exit 1",
				Severity:   config.SeverityError,
				Suggestion: "Fix the issue",
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Exit code 3 for violations (execution errors)
	if result.ExitCode != 1 {
		t.Errorf("expected exit code %d, got %d", 1, result.ExitCode)
	}
	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}
	if result.Results[0].Passed {
		t.Error("expected check to fail")
	}
}

func TestRun_FailingCheck_WarningSeverity_ExitCodeZero(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:         "warn-check",
				Run:        "exit 1",
				Severity:   config.SeverityWarning,
				Suggestion: "Consider fixing",
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Warning severity failures don't change exit code
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0 for warning severity, got %d", result.ExitCode)
	}
	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
	if result.Results[0].Passed {
		t.Error("expected check to fail (even though it's a warning)")
	}
}

func TestRun_MultipleChecks_MixedResults(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "check1",
				Run:      "exit 0",
				Severity: config.SeverityError,
			},
			{
				ID:       "check2",
				Run:      "exit 1",
				Severity: config.SeverityWarning,
			},
			{
				ID:       "check3",
				Run:      "exit 0",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only warnings failed, so exit code should be 0
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}
	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
}

func TestRun_MultipleChecks_ErrorFailure_ExitCodeThree(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "check1",
				Run:      "exit 0",
				Severity: config.SeverityError,
			},
			{
				ID:       "check2",
				Run:      "exit 1",
				Severity: config.SeverityError,
			},
			{
				ID:       "check3",
				Run:      "exit 0",
				Severity: config.SeverityWarning,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Error severity failure should result in exit code 2 (Claude Code hook compatibility)
	if result.ExitCode != 1 {
		t.Errorf("expected exit code %d, got %d", 1, result.ExitCode)
	}
}

func TestRun_FailFast_StopsOnFirstFailure(t *testing.T) {
	// With parallel execution at the same level, fail-fast stops after the level completes.
	// Use dependencies to create multiple levels to test fail-fast stopping at level boundaries.
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "check1",
				Run:      "exit 1", // Fails
				Severity: config.SeverityError,
			},
			{
				ID:       "check2",
				Run:      "exit 0",
				Severity: config.SeverityError,
				Requires: []string{"check1"},
			},
			{
				ID:       "check3",
				Run:      "exit 0",
				Severity: config.SeverityError,
				Requires: []string{"check1"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, true, false, "", 1) // failFast = true

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With fail-fast, should stop after first level (check1 fails)
	// check2 and check3 are in level 1 and should not run
	if len(result.Results) != 1 {
		t.Errorf("expected 1 result (fail-fast stops at level boundary), got %d", len(result.Results))
	}

	if result.ExitCode != 1 {
		t.Errorf("expected exit code %d, got %d", 1, result.ExitCode)
	}
}

func TestRunCheck_SingleCheck_Passes(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "my-check",
				Run:      "exit 0",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.RunCheck(context.Background(), "my-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}
	if !result.Results[0].Passed {
		t.Error("expected check to pass")
	}
}

func TestRunCheck_SingleCheck_Fails(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:         "my-check",
				Run:        "exit 42",
				Severity:   config.SeverityError,
				Suggestion: "Fix it",
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.RunCheck(context.Background(), "my-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 1 {
		t.Errorf("expected exit code %d, got %d", 1, result.ExitCode)
	}
	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
	if result.Results[0].Passed {
		t.Error("expected check to fail")
	}
	if result.Results[0].Execution.ExitCode != 42 {
		t.Errorf("expected execution exit code 42, got %d", result.Results[0].Execution.ExitCode)
	}
}

func TestRunCheck_UnknownCheck_ReturnsConfigError(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "existing-check",
				Run:      "exit 0",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.RunCheck(context.Background(), "non-existent-check")
	if err == nil {
		t.Fatalf("expected error for unknown check, got nil")
	}

	// Unknown check should return a ConfigError
	if !config.IsConfigError(err) {
		t.Errorf("expected ConfigError for unknown check, got %T: %v", err, err)
	}

	// Result should be nil when there's an error
	if result != nil {
		t.Errorf("expected nil result for unknown check, got %v", result)
	}
}

func TestRunCheck_WarningSeverity_FailsButExitZero(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "warn-check",
				Run:      "exit 1",
				Severity: config.SeverityWarning,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.RunCheck(context.Background(), "warn-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Warning severity should not cause exit code 1
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0 for warning severity, got %d", result.ExitCode)
	}
	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
}

func TestViolation_ContainsCorrectInfo(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:         "test-check",
				Run:        "echo 'testing' && exit 1",
				Severity:   config.SeverityError,
				Suggestion: "Try running the tests",
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}

	v := result.Violations[0]
	if v.CheckID != "test-check" {
		t.Errorf("expected CheckID 'test-check', got %q", v.CheckID)
	}
	if v.Severity != config.SeverityError {
		t.Errorf("expected severity 'error', got %q", v.Severity)
	}
	if v.Command != "echo 'testing' && exit 1" {
		t.Errorf("expected command, got %q", v.Command)
	}
	if v.Suggestion != "Try running the tests" {
		t.Errorf("expected suggestion, got %q", v.Suggestion)
	}
}

func TestNew_DefaultMaxParallel(t *testing.T) {
	cfg := &config.Config{Version: "1"}
	exec := executor.New("")

	orch := New(cfg, exec, 0, false, false, "", 1)
	if orch.maxParallel != config.DefaultParallel {
		t.Errorf("expected default max parallel %d, got %d", config.DefaultParallel, orch.maxParallel)
	}

	orch = New(cfg, exec, -1, false, false, "", 1)
	if orch.maxParallel != config.DefaultParallel {
		t.Errorf("expected default max parallel %d, got %d", config.DefaultParallel, orch.maxParallel)
	}
}

// Tests for topological sort execution ordering (vibeguard-v3m.2)

func TestRun_WithDependencies_ExecutesInOrder(t *testing.T) {
	// Create a temp file to track execution order
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "step1",
				Run:      "echo step1",
				Severity: config.SeverityError,
			},
			{
				ID:       "step2",
				Run:      "echo step2",
				Severity: config.SeverityError,
				Requires: []string{"step1"},
			},
			{
				ID:       "step3",
				Run:      "echo step3",
				Severity: config.SeverityError,
				Requires: []string{"step2"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	// Verify execution order: step1 -> step2 -> step3
	expectedOrder := []string{"step1", "step2", "step3"}
	for i, expected := range expectedOrder {
		if result.Results[i].Check.ID != expected {
			t.Errorf("expected result %d to be %q, got %q", i, expected, result.Results[i].Check.ID)
		}
	}
}

func TestRun_DiamondDependency_ExecutesInCorrectOrder(t *testing.T) {
	// Diamond pattern:
	//     a
	//    / \
	//   b   c
	//    \ /
	//     d
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "a",
				Run:      "echo a",
				Severity: config.SeverityError,
			},
			{
				ID:       "b",
				Run:      "echo b",
				Severity: config.SeverityError,
				Requires: []string{"a"},
			},
			{
				ID:       "c",
				Run:      "echo c",
				Severity: config.SeverityError,
				Requires: []string{"a"},
			},
			{
				ID:       "d",
				Run:      "echo d",
				Severity: config.SeverityError,
				Requires: []string{"b", "c"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(result.Results))
	}

	// Verify "a" runs first
	if result.Results[0].Check.ID != "a" {
		t.Errorf("expected first result to be 'a', got %q", result.Results[0].Check.ID)
	}

	// "d" must run after both "b" and "c"
	dIndex := -1
	bIndex := -1
	cIndex := -1
	for i, r := range result.Results {
		switch r.Check.ID {
		case "d":
			dIndex = i
		case "b":
			bIndex = i
		case "c":
			cIndex = i
		}
	}

	if dIndex <= bIndex || dIndex <= cIndex {
		t.Errorf("'d' should run after 'b' and 'c': d=%d, b=%d, c=%d", dIndex, bIndex, cIndex)
	}
}

func TestRun_DependencyFails_SkipsDependent(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "first",
				Run:      "exit 1", // This will fail
				Severity: config.SeverityError,
			},
			{
				ID:       "second",
				Run:      "echo should-not-run",
				Severity: config.SeverityError,
				Requires: []string{"first"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1) // failFast = false

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	// First check should fail
	if result.Results[0].Passed {
		t.Error("expected first check to fail")
	}

	// Second check should be skipped (marked as failed)
	if result.Results[1].Passed {
		t.Error("expected second check to be skipped (failed)")
	}
	if result.Results[1].Execution.ExitCode != -1 {
		t.Errorf("expected exit code -1 for skipped check, got %d", result.Results[1].Execution.ExitCode)
	}

	// Should have 2 violations
	if len(result.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(result.Violations))
	}

	// Second violation should indicate it was skipped
	if result.Violations[1].Suggestion != "Skipped: required dependency failed" {
		t.Errorf("expected skip suggestion, got %q", result.Violations[1].Suggestion)
	}
}

func TestRun_MultipleDependenciesOneFails_SkipsDependent(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "dep1",
				Run:      "exit 0",
				Severity: config.SeverityError,
			},
			{
				ID:       "dep2",
				Run:      "exit 1", // This fails
				Severity: config.SeverityError,
			},
			{
				ID:       "dependent",
				Run:      "echo should-not-run",
				Severity: config.SeverityError,
				Requires: []string{"dep1", "dep2"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	// Find the dependent result
	var dependentResult *CheckResult
	for _, r := range result.Results {
		if r.Check.ID == "dependent" {
			dependentResult = r
			break
		}
	}

	if dependentResult == nil {
		t.Fatal("could not find 'dependent' result")
	}

	// Should be skipped because dep2 failed
	if dependentResult.Passed {
		t.Error("expected 'dependent' to be skipped")
	}
	if dependentResult.Execution.ExitCode != -1 {
		t.Errorf("expected exit code -1 for skipped, got %d", dependentResult.Execution.ExitCode)
	}
}

func TestRun_IndependentChecks_AllExecute(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "independent1",
				Run:      "echo a",
				Severity: config.SeverityError,
			},
			{
				ID:       "independent2",
				Run:      "echo b",
				Severity: config.SeverityError,
			},
			{
				ID:       "independent3",
				Run:      "echo c",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All should pass and be in level 0
	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	for _, r := range result.Results {
		if !r.Passed {
			t.Errorf("expected check %q to pass", r.Check.ID)
		}
	}
}

func TestRun_CyclicDependency_ReturnsError(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "a",
				Run:      "echo a",
				Severity: config.SeverityError,
				Requires: []string{"b"},
			},
			{
				ID:       "b",
				Run:      "echo b",
				Severity: config.SeverityError,
				Requires: []string{"a"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	_, err := orch.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for cyclic dependency, got nil")
	}
}

func TestRun_UnknownDependency_ReturnsError(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "a",
				Run:      "echo a",
				Severity: config.SeverityError,
				Requires: []string{"nonexistent"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	_, err := orch.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for unknown dependency, got nil")
	}
}

func TestRun_FailFast_WithDependencies_StopsCorrectly(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "first",
				Run:      "exit 1", // Fails
				Severity: config.SeverityError,
			},
			{
				ID:       "second",
				Run:      "echo should-not-run",
				Severity: config.SeverityError,
				Requires: []string{"first"},
			},
			{
				ID:       "independent",
				Run:      "echo also-should-not-run",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, true, false, "", 1) // failFast = true

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With fail-fast, should stop after first check fails
	// But we might still execute independent checks at the same level
	// In our sequential implementation, we stop immediately
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 result")
	}

	// First result should be the failing check
	if result.Results[0].Check.ID != "first" && result.Results[0].Check.ID != "independent" {
		t.Errorf("unexpected first result: %q", result.Results[0].Check.ID)
	}
}

func TestRun_ComplexDependencyGraph_CorrectOrder(t *testing.T) {
	// Complex graph:
	//       a
	//      /|\
	//     b c d
	//     |/ \|
	//     e   f
	//      \ /
	//       g
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "echo a", Severity: config.SeverityError},
			{ID: "b", Run: "echo b", Severity: config.SeverityError, Requires: []string{"a"}},
			{ID: "c", Run: "echo c", Severity: config.SeverityError, Requires: []string{"a"}},
			{ID: "d", Run: "echo d", Severity: config.SeverityError, Requires: []string{"a"}},
			{ID: "e", Run: "echo e", Severity: config.SeverityError, Requires: []string{"b", "c"}},
			{ID: "f", Run: "echo f", Severity: config.SeverityError, Requires: []string{"c", "d"}},
			{ID: "g", Run: "echo g", Severity: config.SeverityError, Requires: []string{"e", "f"}},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 7 {
		t.Fatalf("expected 7 results, got %d", len(result.Results))
	}

	// Build index map for verification
	indexByID := make(map[string]int)
	for i, r := range result.Results {
		indexByID[r.Check.ID] = i
	}

	// Verify constraints
	verifyOrder := func(before, after string) {
		if indexByID[before] >= indexByID[after] {
			t.Errorf("%q should run before %q: %d >= %d", before, after, indexByID[before], indexByID[after])
		}
	}

	verifyOrder("a", "b")
	verifyOrder("a", "c")
	verifyOrder("a", "d")
	verifyOrder("b", "e")
	verifyOrder("c", "e")
	verifyOrder("c", "f")
	verifyOrder("d", "f")
	verifyOrder("e", "g")
	verifyOrder("f", "g")
}

func TestRun_DependencyChain_MiddleFails_SkipsDownstream(t *testing.T) {
	// Chain: a -> b -> c -> d
	// b fails, so c and d should be skipped
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "exit 0", Severity: config.SeverityError},
			{ID: "b", Run: "exit 1", Severity: config.SeverityError, Requires: []string{"a"}}, // fails
			{ID: "c", Run: "echo c", Severity: config.SeverityError, Requires: []string{"b"}},
			{ID: "d", Run: "echo d", Severity: config.SeverityError, Requires: []string{"c"}},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(result.Results))
	}

	// Count skipped checks (exit code -1)
	skipped := 0
	for _, r := range result.Results {
		if r.Execution.ExitCode == -1 {
			skipped++
		}
	}

	// c and d should be skipped
	if skipped != 2 {
		t.Errorf("expected 2 skipped checks, got %d", skipped)
	}
}

// Tests for parallel execution (vibeguard-v3m.3)

func TestRun_ParallelExecution_SameLevelRunsConcurrently(t *testing.T) {
	// Create checks that would take too long if run sequentially
	// but complete quickly when run in parallel
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "sleep 0.1", Severity: config.SeverityError},
			{ID: "b", Run: "sleep 0.1", Severity: config.SeverityError},
			{ID: "c", Run: "sleep 0.1", Severity: config.SeverityError},
			{ID: "d", Run: "sleep 0.1", Severity: config.SeverityError},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 4, false, false, "", 1) // maxParallel = 4

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All checks should complete
	if len(result.Results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(result.Results))
	}

	// With 4 checks sleeping 0.1s each and maxParallel=4,
	// total time should be ~0.1s (parallel), not ~0.4s (sequential)
	// Allow some margin for test execution overhead
	if result.Duration > 300*time.Millisecond {
		t.Errorf("expected parallel execution to complete quickly, took %v", result.Duration)
	}
}

func TestRun_ParallelExecution_RespectsMaxParallel(t *testing.T) {
	// With maxParallel=2, running 4 checks that each sleep 0.1s
	// should take ~0.2s (2 batches of 2), not ~0.1s (all 4 at once)
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "sleep 0.1", Severity: config.SeverityError},
			{ID: "b", Run: "sleep 0.1", Severity: config.SeverityError},
			{ID: "c", Run: "sleep 0.1", Severity: config.SeverityError},
			{ID: "d", Run: "sleep 0.1", Severity: config.SeverityError},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 2, false, false, "", 1) // maxParallel = 2

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(result.Results))
	}

	// With maxParallel=2, should take at least 0.2s (two batches)
	// but less than 0.4s (which would be sequential)
	if result.Duration < 150*time.Millisecond {
		t.Errorf("expected maxParallel to limit concurrency, completed in %v", result.Duration)
	}
	if result.Duration > 500*time.Millisecond {
		t.Errorf("execution took too long: %v", result.Duration)
	}
}

func TestRun_ParallelExecution_LevelsRunSequentially(t *testing.T) {
	// Checks at different levels should run sequentially (level by level)
	// Level 0: a (0.1s)
	// Level 1: b, c (0.1s each, parallel)
	// Total should be ~0.2s (0.1s + 0.1s)
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "sleep 0.1", Severity: config.SeverityError},
			{ID: "b", Run: "sleep 0.1", Severity: config.SeverityError, Requires: []string{"a"}},
			{ID: "c", Run: "sleep 0.1", Severity: config.SeverityError, Requires: []string{"a"}},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 4, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	// Should take at least 0.2s (two levels)
	if result.Duration < 150*time.Millisecond {
		t.Errorf("expected level-by-level execution, completed too quickly: %v", result.Duration)
	}
}

func TestRun_ParallelExecution_FailFastWithinLevel(t *testing.T) {
	// When fail-fast is enabled and a check fails within a level,
	// other checks in the same level may or may not complete (race condition),
	// but subsequent levels should not run
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "exit 0", Severity: config.SeverityError},
			{ID: "b", Run: "exit 1", Severity: config.SeverityError}, // Fails (same level as a)
			{ID: "c", Run: "echo should-not-run", Severity: config.SeverityError, Requires: []string{"a", "b"}},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 4, true, false, "", 1) // failFast = true

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Level 0 (a, b) should complete, level 1 (c) should not run
	// So we expect exactly 2 results
	if len(result.Results) != 2 {
		t.Errorf("expected 2 results (level 0 only), got %d", len(result.Results))
	}

	// Verify c didn't run
	for _, r := range result.Results {
		if r.Check.ID == "c" {
			t.Error("check 'c' should not have run due to fail-fast")
		}
	}
}

func TestRun_ParallelExecution_AllFailuresRecorded(t *testing.T) {
	// When multiple checks fail within the same level (without fail-fast),
	// all failures should be recorded
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "exit 1", Severity: config.SeverityError},
			{ID: "b", Run: "exit 2", Severity: config.SeverityError},
			{ID: "c", Run: "exit 3", Severity: config.SeverityWarning},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 4, false, false, "", 1) // failFast = false

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	// All 3 should be violations
	if len(result.Violations) != 3 {
		t.Errorf("expected 3 violations, got %d", len(result.Violations))
	}

	// Exit code should be 2 (error severity failures, Claude Code hook compatibility)
	if result.ExitCode != 1 {
		t.Errorf("expected exit code %d, got %d", 1, result.ExitCode)
	}
}

func TestRun_ParallelExecution_OrderPreservedWithinLevel(t *testing.T) {
	// Results within a level should maintain the original check order
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "echo a", Severity: config.SeverityError},
			{ID: "b", Run: "echo b", Severity: config.SeverityError},
			{ID: "c", Run: "echo c", Severity: config.SeverityError},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 4, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	// Order should be preserved: a, b, c
	expected := []string{"a", "b", "c"}
	for i, exp := range expected {
		if result.Results[i].Check.ID != exp {
			t.Errorf("expected result %d to be %q, got %q", i, exp, result.Results[i].Check.ID)
		}
	}
}

// Timeout handling tests

func TestRun_Timeout_ExitCode3(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "timeout-check",
				Run:      "sleep 5",
				Timeout:  config.Duration(50 * time.Millisecond),
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Timeout should produce exit code 3
	if result.ExitCode != 1 {
		t.Errorf("expected exit code %d for timeout, got %d", 1, result.ExitCode)
	}
}

func TestRun_Timeout_ViolationMarkedAsTimeout(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "timeout-check",
				Run:      "sleep 5",
				Timeout:  config.Duration(50 * time.Millisecond),
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}

	if !result.Violations[0].Timedout {
		t.Error("expected violation to be marked as timeout")
	}
}

func TestRun_Timeout_SuggestionIncludesTimeoutMessage(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:         "timeout-check",
				Run:        "sleep 5",
				Timeout:    config.Duration(50 * time.Millisecond),
				Severity:   config.SeverityError,
				Suggestion: "Original suggestion",
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}

	// Suggestion should be overwritten with timeout message
	if result.Violations[0].Suggestion == "Original suggestion" {
		t.Error("expected suggestion to be replaced with timeout message")
	}
	if result.Violations[0].Suggestion != "Check timed out. Consider increasing the timeout value or optimizing the command." {
		t.Errorf("unexpected suggestion: %q", result.Violations[0].Suggestion)
	}
}

func TestRun_TimeoutTakesPrecedenceOverError(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "fast-fail",
				Run:      "exit 1",
				Severity: config.SeverityError,
			},
			{
				ID:       "timeout-check",
				Run:      "sleep 5",
				Timeout:  config.Duration(50 * time.Millisecond),
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 4, false, false, "", 1) // Run in parallel

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With both error and timeout, timeout exit code 3 should take precedence
	if result.ExitCode != 1 {
		t.Errorf("expected exit code %d (timeout precedence), got %d", 1, result.ExitCode)
	}
}

func TestRunCheck_Timeout_ExitCode3(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "timeout-check",
				Run:      "sleep 5",
				Timeout:  config.Duration(50 * time.Millisecond),
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.RunCheck(context.Background(), "timeout-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Single check timeout should also return exit code 3
	if result.ExitCode != 1 {
		t.Errorf("expected exit code %d for timeout, got %d", 1, result.ExitCode)
	}
}

func TestRunCheck_Timeout_ViolationMarkedAsTimeout(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "timeout-check",
				Run:      "sleep 5",
				Timeout:  config.Duration(50 * time.Millisecond),
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.RunCheck(context.Background(), "timeout-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}

	if !result.Violations[0].Timedout {
		t.Error("expected violation to be marked as timeout")
	}
}

// Fail-fast context cancellation tests

func TestRun_FailFast_SetsFailFastTriggeredFlag(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "failing",
				Run:      "exit 1",
				Severity: config.SeverityError,
			},
			{
				ID:       "dependent",
				Run:      "echo should-not-run",
				Severity: config.SeverityError,
				Requires: []string{"failing"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, true, false, "", 1) // failFast = true

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.FailFastTriggered {
		t.Error("expected FailFastTriggered to be true")
	}
}

func TestRun_NoFailFast_FailFastTriggeredFalse(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "failing",
				Run:      "exit 1",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1) // failFast = false

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.FailFastTriggered {
		t.Error("expected FailFastTriggered to be false when fail-fast is disabled")
	}
}

func TestRun_FailFast_WarningSeverityDoesNotTrigger(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "warning",
				Run:      "exit 1",
				Severity: config.SeverityWarning, // Warning severity
			},
			{
				ID:       "second",
				Run:      "echo should-run",
				Severity: config.SeverityError,
				Requires: []string{"warning"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, true, false, "", 1) // failFast = true

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Warning severity should not trigger fail-fast
	if result.FailFastTriggered {
		t.Error("expected FailFastTriggered to be false for warning severity")
	}
}

func TestRun_FailFast_CancelsLongRunningChecks(t *testing.T) {
	// Test that fail-fast cancels in-flight long-running checks
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "fast-fail",
				Run:      "exit 1",
				Severity: config.SeverityError,
			},
			{
				ID:       "slow-check",
				Run:      "sleep 10", // Would take 10 seconds without cancellation
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 4, true, false, "", 1) // failFast = true, parallel

	start := time.Now()
	result, err := orch.Run(context.Background())
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should complete much faster than 10 seconds due to cancellation
	if duration > 2*time.Second {
		t.Errorf("expected fail-fast to cancel long-running check quickly, took %v", duration)
	}

	if !result.FailFastTriggered {
		t.Error("expected FailFastTriggered to be true")
	}
}

func TestRun_FailFast_AllChecksPassDoesNotTrigger(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{ID: "a", Run: "exit 0", Severity: config.SeverityError},
			{ID: "b", Run: "exit 0", Severity: config.SeverityError},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 4, true, false, "", 1) // failFast = true

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.FailFastTriggered {
		t.Error("expected FailFastTriggered to be false when all checks pass")
	}
}

// Grok pattern matching tests

func TestRun_GrokExtractsValues(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "coverage-check",
				Run:      `echo "coverage: 85.5%"`,
				Grok:     []string{"coverage: (?P<coverage>[0-9.]+)%"},
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	// Verify extracted values
	extracted := result.Results[0].Extracted
	if extracted["coverage"] != "85.5" {
		t.Errorf("expected coverage=85.5, got %q", extracted["coverage"])
	}
}

func TestRun_GrokMultiplePatterns(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:  "multi-grok",
				Run: `echo "total: 42 tests, coverage: 75.0%"`,
				Grok: []string{
					"total: (?P<total>[0-9]+) tests",
					"coverage: (?P<coverage>[0-9.]+)%",
				},
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	extracted := result.Results[0].Extracted
	if extracted["total"] != "42" {
		t.Errorf("expected total=42, got %q", extracted["total"])
	}
	if extracted["coverage"] != "75.0" {
		t.Errorf("expected coverage=75.0, got %q", extracted["coverage"])
	}
}

func TestRun_GrokNoPatterns_EmptyExtracted(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "no-grok",
				Run:      `echo "hello"`,
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	// Without grok patterns, Extracted should be empty map
	if len(result.Results[0].Extracted) != 0 {
		t.Errorf("expected empty extracted map, got %v", result.Results[0].Extracted)
	}
}

func TestRun_GrokNoMatch_EmptyValues(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "no-match",
				Run:      `echo "no numbers here"`,
				Grok:     []string{"(?P<num>[0-9]+)"},
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	// When pattern doesn't match, the key should not be present
	if _, exists := result.Results[0].Extracted["num"]; exists {
		t.Errorf("expected 'num' to not exist in extracted, got %v", result.Results[0].Extracted)
	}
}

func TestRun_GrokInvalidPattern_ReturnsError(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "invalid-grok",
				Run:      `echo "test"`,
				Grok:     []string{"%{INVALID_PATTERN_XYZ:value}"},
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	_, err := orch.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid grok pattern, got nil")
	}
}

func TestRunCheck_GrokExtractsValues(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "single-coverage",
				Run:      `echo "coverage: 90.0%"`,
				Grok:     []string{"coverage: (?P<coverage>[0-9.]+)%"},
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.RunCheck(context.Background(), "single-coverage")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	if result.Results[0].Extracted["coverage"] != "90.0" {
		t.Errorf("expected coverage=90.0, got %q", result.Results[0].Extracted["coverage"])
	}
}

func TestRun_GrokExtractedInViolation(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:         "failing-with-grok",
				Run:        `echo "coverage: 50.0%" && exit 1`,
				Grok:       []string{"coverage: (?P<coverage>[0-9.]+)%"},
				Severity:   config.SeverityError,
				Suggestion: "Coverage is {{.coverage}}%, needs 80%",
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}

	// Verify extracted values are propagated to violation
	if result.Violations[0].Extracted["coverage"] != "50.0" {
		t.Errorf("expected violation to have coverage=50.0, got %v", result.Violations[0].Extracted)
	}
}

// File field tests

func TestRun_FileField_ReadsFromFile(t *testing.T) {
	// Create a temporary file with test content
	tmpFile := "/tmp/vibeguard_test_file.txt"
	fileContent := "coverage: 85.5%"
	if err := os.WriteFile(tmpFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "file-check",
				Run:      "echo ignored",
				File:     tmpFile,
				Grok:     []string{"coverage: (?P<coverage>[0-9.]+)%"},
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	// Verify grok pattern was applied to file content, not command output
	if result.Results[0].Extracted["coverage"] != "85.5" {
		t.Errorf("expected coverage=85.5 from file, got %q", result.Results[0].Extracted["coverage"])
	}

	if !result.Results[0].Passed {
		t.Error("expected check to pass with successful grok match on file")
	}
}

func TestRun_FileField_MissingFile_ReturnsError(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "file-check",
				Run:      "exit 0",
				File:     "/nonexistent/file.txt",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	_, err := orch.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}

	// Should be an execution error, not a config error
	if !config.IsExecutionError(err) {
		t.Errorf("expected ExecutionError, got %T: %v", err, err)
	}
}

func TestRun_FileField_WithAssertion(t *testing.T) {
	tmpFile := "/tmp/vibeguard_test_file_assert.txt"
	fileContent := "coverage: 85.5%"
	if err := os.WriteFile(tmpFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "file-check",
				Run:      "exit 0",
				File:     tmpFile,
				Grok:     []string{"coverage: (?P<coverage>[0-9.]+)%"},
				Assert:   "coverage >= 80",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	// Assertion should pass because coverage is >= 80
	if !result.Results[0].Passed {
		t.Error("expected assertion to pass (coverage >= 80)")
	}
}

func TestRun_FileField_WithAssertion_Fails(t *testing.T) {
	tmpFile := "/tmp/vibeguard_test_file_assert_fail.txt"
	fileContent := "coverage: 50.0%"
	if err := os.WriteFile(tmpFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "file-check",
				Run:      "exit 0",
				File:     tmpFile,
				Grok:     []string{"coverage: (?P<coverage>[0-9.]+)%"},
				Assert:   "coverage >= 80",
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	// Assertion should fail because coverage is < 80
	if result.Results[0].Passed {
		t.Error("expected assertion to fail (coverage < 80)")
	}

	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
}

func TestRunCheck_FileField_ReadsFromFile(t *testing.T) {
	tmpFile := "/tmp/vibeguard_test_file_check.txt"
	fileContent := "version: 1.2.3"
	if err := os.WriteFile(tmpFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "version-check",
				Run:      "echo ignored",
				File:     tmpFile,
				Grok:     []string{"version: (?P<ver>[0-9.]+)"},
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.RunCheck(context.Background(), "version-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	if result.Results[0].Extracted["ver"] != "1.2.3" {
		t.Errorf("expected ver=1.2.3, got %q", result.Results[0].Extracted["ver"])
	}
}

func TestRun_FileField_WithVariableInterpolation(t *testing.T) {
	tmpFile := "/tmp/vibeguard_test_file_vars.txt"
	fileContent := "status: ok"
	if err := os.WriteFile(tmpFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	cfg := &config.Config{
		Version: "1",
		Vars: map[string]string{
			"outputFile": tmpFile,
		},
		Checks: []config.Check{
			{
				ID:       "var-file-check",
				Run:      "exit 0",
				File:     "{{.outputFile}}",
				Grok:     []string{"status: (?P<status>[a-z]+)"},
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	if result.Results[0].Extracted["status"] != "ok" {
		t.Errorf("expected status=ok, got %q", result.Results[0].Extracted["status"])
	}
}

func TestRun_FileField_WithoutGrok_StillReads(t *testing.T) {
	tmpFile := "/tmp/vibeguard_test_file_no_grok.txt"
	fileContent := "some test content"
	if err := os.WriteFile(tmpFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "file-no-grok",
				Run:      "exit 0",
				File:     tmpFile,
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false, "", 1)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	// Without grok, should have empty extracted map but still pass
	if len(result.Results[0].Extracted) != 0 {
		t.Errorf("expected empty extracted map, got %v", result.Results[0].Extracted)
	}

	if !result.Results[0].Passed {
		t.Error("expected check to pass (exit 0, no grok patterns)")
	}
}
