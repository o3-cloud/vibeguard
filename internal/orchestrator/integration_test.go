// Package orchestrator provides integration tests for real tool execution
package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/executor"
)

// TestIntegration_RealToolExecution_GoFmt tests gofmt execution
// This test requires gofmt to be installed
func TestIntegration_RealToolExecution_GoFmt(t *testing.T) {
	if _, ok := os.LookupEnv("CI"); !ok {
		t.Skip("Skipping integration test: not in CI environment")
	}

	// Create a test Go file that's properly formatted
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(`package main

func main() {
	println("hello")
}
`), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "fmt",
				Run:      "test -z \"$(gofmt -l " + tmpDir + ")\"",
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
		t.Errorf("expected exit code 0 for properly formatted code, got %d", result.ExitCode)
	}
	if len(result.Violations) != 0 {
		t.Errorf("expected no violations, got %d", len(result.Violations))
	}
}

// TestIntegration_RealToolExecution_EchoCommand tests echo command
// This is a simple integration test that should work on all platforms
func TestIntegration_RealToolExecution_EchoCommand(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "echo-test",
				Run:      "echo 'hello world'",
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
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}
	if !result.Results[0].Passed {
		t.Error("expected check to pass")
	}
	if result.Results[0].Execution.Stdout != "hello world\n" {
		t.Errorf("expected stdout 'hello world\\n', got %q", result.Results[0].Execution.Stdout)
	}
}

// TestIntegration_GrokExtraction_RealToolOutput tests grok extraction with actual command output
func TestIntegration_GrokExtraction_RealToolOutput(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "extract-date",
				Run:      "echo 'Date: 2025-12-31 Time: 12:30:45'",
				Grok:     []string{"Date: (?P<date>\\d{4}-\\d{2}-\\d{2})", "Time: (?P<time>\\d{2}:\\d{2}:\\d{2})"},
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

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	extracted := result.Results[0].Extracted
	if extracted["date"] != "2025-12-31" {
		t.Errorf("expected date='2025-12-31', got %q", extracted["date"])
	}
	if extracted["time"] != "12:30:45" {
		t.Errorf("expected time='12:30:45', got %q", extracted["time"])
	}
}

// TestIntegration_GrokExtraction_MultiLineOutput tests grok with multi-line output
func TestIntegration_GrokExtraction_MultiLineOutput(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "extract-multiline",
				Run:      "echo 'Line 1: value1\\nLine 2: value2\\nLine 3: value3'",
				Grok:     []string{"Line 1: (?P<val1>\\w+)", "Line 3: (?P<val3>\\w+)"},
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

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	extracted := result.Results[0].Extracted
	if extracted["val1"] != "value1" {
		t.Errorf("expected val1='value1', got %q", extracted["val1"])
	}
	if extracted["val3"] != "value3" {
		t.Errorf("expected val3='value3', got %q", extracted["val3"])
	}
}

// TestIntegration_GrokExtraction_InViolationSuggestion tests that extracted values are in violations
func TestIntegration_GrokExtraction_InViolationSuggestion(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:         "coverage-check",
				Run:        "echo 'coverage: 50%' && exit 1",
				Grok:       []string{"coverage: (?P<pct>\\d+)"},
				Severity:   config.SeverityError,
				Suggestion: "Coverage check failed",
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

	// Verify that extracted values are available in the violation
	violation := result.Violations[0]
	if violation.Extracted["pct"] != "50" {
		t.Errorf("expected extracted pct=50, got %q", violation.Extracted["pct"])
	}
}

// TestIntegration_DependencyOrdering_RealTools tests execution order with dependencies
func TestIntegration_DependencyOrdering_RealTools(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files to track execution order
	trackerFile := filepath.Join(tmpDir, "execution_order.txt")

	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "first",
				Run:      "echo 'first' >> " + trackerFile,
				Severity: config.SeverityError,
			},
			{
				ID:       "second",
				Run:      "echo 'second' >> " + trackerFile,
				Severity: config.SeverityError,
				Requires: []string{"first"},
			},
			{
				ID:       "third",
				Run:      "echo 'third' >> " + trackerFile,
				Severity: config.SeverityError,
				Requires: []string{"second"},
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
	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	// Verify execution order in results
	expectedOrder := []string{"first", "second", "third"}
	for i, expected := range expectedOrder {
		if result.Results[i].Check.ID != expected {
			t.Errorf("expected result %d to be %q, got %q", i, expected, result.Results[i].Check.ID)
		}
	}
}

// TestIntegration_DependencyOrdering_FailurePreventsDownstream tests that failure skips dependents
func TestIntegration_DependencyOrdering_FailurePreventsDownstream(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "first",
				Run:      "exit 1",
				Severity: config.SeverityError,
			},
			{
				ID:       "second",
				Run:      "exit 0",
				Severity: config.SeverityError,
				Requires: []string{"first"},
			},
		},
	}

	exec := executor.New("")
	orch := New(cfg, exec, 1, false, false)

	result, err := orch.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	// First should fail
	if result.Results[0].Passed {
		t.Error("expected first check to fail")
	}

	// Second should be skipped (exit code -1)
	if result.Results[1].Execution.ExitCode != -1 {
		t.Errorf("expected second check to be skipped (exit code -1), got %d", result.Results[1].Execution.ExitCode)
	}
}

// TestIntegration_TimeoutHandling_CommandExceedsTimeout tests timeout with real command
func TestIntegration_TimeoutHandling_CommandExceedsTimeout(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "timeout-test",
				Run:      "sleep 5",
				Timeout:  config.Duration(100 * time.Millisecond),
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

	// Timeout should result in exit code 4
	if result.ExitCode != executor.ExitCodeTimeout {
		t.Errorf("expected exit code %d for timeout, got %d", executor.ExitCodeTimeout, result.ExitCode)
	}

	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation for timeout, got %d", len(result.Violations))
	}

	if !result.Violations[0].Timedout {
		t.Error("expected violation to be marked as timeout")
	}
}

// TestIntegration_TimeoutHandling_CommandCompletesBefore Timeout tests command completes within timeout
func TestIntegration_TimeoutHandling_CommandCompletesBeforeTimeout(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "quick-test",
				Run:      "echo 'done'",
				Timeout:  config.Duration(5 * time.Second),
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
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}
	if !result.Results[0].Passed {
		t.Error("expected check to pass")
	}
	if result.Results[0].Execution.Timedout {
		t.Error("expected check to not be marked as timed out")
	}
}

// TestIntegration_ComplexWorkflow tests a realistic workflow with multiple checks
func TestIntegration_ComplexWorkflow(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Checks: []config.Check{
			{
				ID:       "setup",
				Run:      "echo 'Setting up environment'",
				Severity: config.SeverityError,
			},
			{
				ID:       "validate",
				Run:      "echo 'Validation passed'",
				Severity: config.SeverityError,
				Requires: []string{"setup"},
			},
			{
				ID:         "test",
				Run:        "echo 'Tests: 42 passed'",
				Grok:       []string{"Tests: (?P<passed>\\d+) passed"},
				Severity:   config.SeverityError,
				Requires:   []string{"validate"},
				Suggestion: "Ran {{.passed}} tests",
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
	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	// Verify all passed
	for _, r := range result.Results {
		if !r.Passed {
			t.Errorf("expected check %q to pass", r.Check.ID)
		}
	}

	// Verify grok extraction worked
	testResult := result.Results[2]
	if testResult.Extracted["passed"] != "42" {
		t.Errorf("expected extracted passed=42, got %q", testResult.Extracted["passed"])
	}
}
