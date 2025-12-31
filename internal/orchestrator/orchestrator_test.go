// Package orchestrator coordinates check execution with dependency management
// and parallel execution.
package orchestrator

import (
	"context"
	"testing"

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
	orch := New(cfg, exec, 1, false, false)

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

func TestRun_FailingCheck_ErrorSeverity_ExitCodeOne(t *testing.T) {
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
	orch := New(cfg, exec, 1, false, false)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
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
	orch := New(cfg, exec, 1, false, false)

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
	orch := New(cfg, exec, 1, false, false)

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

func TestRun_MultipleChecks_ErrorFailure_ExitCodeOne(t *testing.T) {
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
	orch := New(cfg, exec, 1, false, false)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Error severity failure should result in exit code 1
	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
}

func TestRun_FailFast_StopsOnFirstFailure(t *testing.T) {
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
				Severity: config.SeverityError,
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, true, false) // failFast = true

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With fail-fast, should stop after check2 fails
	if len(result.Results) != 2 {
		t.Errorf("expected 2 results (fail-fast), got %d", len(result.Results))
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
	orch := New(cfg, exec, 1, false, false)

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
	orch := New(cfg, exec, 1, false, false)

	result, err := orch.RunCheck(context.Background(), "my-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
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

func TestRunCheck_UnknownCheck_ExitCodeTwo(t *testing.T) {
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
	orch := New(cfg, exec, 1, false, false)

	result, err := orch.RunCheck(context.Background(), "non-existent-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Unknown check should return exit code 2
	if result.ExitCode != 2 {
		t.Errorf("expected exit code 2 for unknown check, got %d", result.ExitCode)
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
	orch := New(cfg, exec, 1, false, false)

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
	orch := New(cfg, exec, 1, false, false)

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

	orch := New(cfg, exec, 0, false, false)
	if orch.maxParallel != config.DefaultParallel {
		t.Errorf("expected default max parallel %d, got %d", config.DefaultParallel, orch.maxParallel)
	}

	orch = New(cfg, exec, -1, false, false)
	if orch.maxParallel != config.DefaultParallel {
		t.Errorf("expected default max parallel %d, got %d", config.DefaultParallel, orch.maxParallel)
	}
}
