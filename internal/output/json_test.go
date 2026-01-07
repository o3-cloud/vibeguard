package output

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/executor"
	"github.com/vibeguard/vibeguard/internal/orchestrator"
)

func TestFormatJSON_AllPassing(t *testing.T) {
	var buf bytes.Buffer

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check:     &config.Check{ID: "fmt"},
				Execution: &executor.Result{Duration: 100 * time.Millisecond},
				Passed:    true,
			},
			{
				Check:     &config.Check{ID: "vet"},
				Execution: &executor.Result{Duration: 300 * time.Millisecond},
				Passed:    true,
			},
		},
		Violations: nil,
		ExitCode:   0,
	}

	err := FormatJSON(&buf, result)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	var output JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if len(output.Checks) != 2 {
		t.Errorf("expected 2 checks, got %d", len(output.Checks))
	}
	if output.Checks[0].ID != "fmt" || output.Checks[0].Status != "passed" {
		t.Errorf("unexpected check[0]: %+v", output.Checks[0])
	}
	if output.Checks[1].ID != "vet" || output.Checks[1].Status != "passed" {
		t.Errorf("unexpected check[1]: %+v", output.Checks[1])
	}
	if len(output.Violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(output.Violations))
	}
	if output.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", output.ExitCode)
	}
}

func TestFormatJSON_WithViolations(t *testing.T) {
	var buf bytes.Buffer

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check:     &config.Check{ID: "fmt"},
				Execution: &executor.Result{Duration: 100 * time.Millisecond},
				Passed:    true,
			},
			{
				Check:     &config.Check{ID: "coverage", Severity: config.SeverityError},
				Execution: &executor.Result{Duration: 900 * time.Millisecond},
				Passed:    false,
			},
		},
		Violations: []*orchestrator.Violation{
			{
				CheckID:    "coverage",
				Severity:   config.SeverityError,
				Command:    "go test -cover ./...",
				Suggestion: "Coverage is 72%, need 80%.",
				Extracted:  map[string]string{"coverage": "72"},
			},
		},
		ExitCode: 1, // Default error exit code
	}

	err := FormatJSON(&buf, result)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	var output JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if len(output.Checks) != 2 {
		t.Errorf("expected 2 checks, got %d", len(output.Checks))
	}
	if output.Checks[0].Status != "passed" {
		t.Errorf("expected check[0] to be passed, got %s", output.Checks[0].Status)
	}
	if output.Checks[1].Status != "failed" {
		t.Errorf("expected check[1] to be failed, got %s", output.Checks[1].Status)
	}
	if len(output.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(output.Violations))
	}
	if output.Violations[0].ID != "coverage" {
		t.Errorf("expected violation id 'coverage', got %s", output.Violations[0].ID)
	}
	if output.Violations[0].Severity != "error" {
		t.Errorf("expected severity 'error', got %s", output.Violations[0].Severity)
	}
	if output.Violations[0].Extracted["coverage"] != "72" {
		t.Errorf("expected extracted coverage '72', got %s", output.Violations[0].Extracted["coverage"])
	}
	if output.ExitCode != 1 { // Default error exit code
		t.Errorf("expected exit code %d, got %d", 1, output.ExitCode)
	}
}

func TestFormatJSON_DurationInMilliseconds(t *testing.T) {
	var buf bytes.Buffer

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check:     &config.Check{ID: "test"},
				Execution: &executor.Result{Duration: 1500 * time.Millisecond},
				Passed:    true,
			},
		},
		Violations: nil,
		ExitCode:   0,
	}

	err := FormatJSON(&buf, result)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	var output JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if output.Checks[0].DurationMS != 1500 {
		t.Errorf("expected duration_ms 1500, got %d", output.Checks[0].DurationMS)
	}
}

func TestFormatJSON_CancelledStatus(t *testing.T) {
	var buf bytes.Buffer

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check:     &config.Check{ID: "fmt"},
				Execution: &executor.Result{Duration: 100 * time.Millisecond},
				Passed:    true,
			},
			{
				Check:     &config.Check{ID: "vet"},
				Execution: &executor.Result{Duration: 250 * time.Millisecond, Cancelled: true},
				Passed:    false,
			},
		},
		Violations: nil,
		ExitCode:   1, // Default error exit code
	}

	err := FormatJSON(&buf, result)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	var output JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if len(output.Checks) != 2 {
		t.Errorf("expected 2 checks, got %d", len(output.Checks))
	}
	if output.Checks[0].ID != "fmt" || output.Checks[0].Status != "passed" {
		t.Errorf("unexpected check[0]: %+v", output.Checks[0])
	}
	if output.Checks[1].ID != "vet" || output.Checks[1].Status != "cancelled" {
		t.Errorf("expected check[1] status to be 'cancelled', got: %+v", output.Checks[1])
	}
	if output.ExitCode != 1 { // Default error exit code
		t.Errorf("expected exit code %d, got %d", 1, output.ExitCode)
	}
}

func TestFormatJSON_WithFixField(t *testing.T) {
	var buf bytes.Buffer

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check:     &config.Check{ID: "coverage", Severity: config.SeverityError},
				Execution: &executor.Result{Duration: 900 * time.Millisecond},
				Passed:    false,
			},
		},
		Violations: []*orchestrator.Violation{
			{
				CheckID:    "coverage",
				Severity:   config.SeverityError,
				Command:    "go test -cover ./...",
				Suggestion: "Coverage is 72%, need 80%.",
				Fix:        "Add tests to improve coverage",
				Extracted:  map[string]string{"coverage": "72"},
			},
		},
		ExitCode: 1, // Default error exit code
	}

	err := FormatJSON(&buf, result)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	var output JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if len(output.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(output.Violations))
	}

	v := output.Violations[0]
	if v.Suggestion != "Coverage is 72%, need 80%." {
		t.Errorf("expected suggestion 'Coverage is 72%%, need 80%%.', got %q", v.Suggestion)
	}
	if v.Fix != "Add tests to improve coverage" {
		t.Errorf("expected fix 'Add tests to improve coverage', got %q", v.Fix)
	}
}

func TestFormatJSON_WithTags(t *testing.T) {
	var buf bytes.Buffer

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check: &config.Check{
					ID:   "fmt",
					Tags: []string{"format", "fast", "pre-commit"},
				},
				Execution: &executor.Result{Duration: 100 * time.Millisecond},
				Passed:    true,
			},
			{
				Check: &config.Check{
					ID:   "security",
					Tags: []string{"security", "slow"},
				},
				Execution: &executor.Result{Duration: 500 * time.Millisecond},
				Passed:    true,
			},
			{
				Check: &config.Check{
					ID: "lint",
				},
				Execution: &executor.Result{Duration: 200 * time.Millisecond},
				Passed:    true,
			},
		},
		Violations: nil,
		ExitCode:   0,
	}

	err := FormatJSON(&buf, result)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	var output JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if len(output.Checks) != 3 {
		t.Errorf("expected 3 checks, got %d", len(output.Checks))
	}

	// Check fmt with tags
	if output.Checks[0].ID != "fmt" {
		t.Errorf("expected check[0] id 'fmt', got %q", output.Checks[0].ID)
	}
	if len(output.Checks[0].Tags) != 3 {
		t.Errorf("expected check[0] to have 3 tags, got %d", len(output.Checks[0].Tags))
	}
	if output.Checks[0].Tags[0] != "format" || output.Checks[0].Tags[1] != "fast" || output.Checks[0].Tags[2] != "pre-commit" {
		t.Errorf("unexpected tags for check[0]: %v", output.Checks[0].Tags)
	}

	// Check security with tags
	if output.Checks[1].ID != "security" {
		t.Errorf("expected check[1] id 'security', got %q", output.Checks[1].ID)
	}
	if len(output.Checks[1].Tags) != 2 {
		t.Errorf("expected check[1] to have 2 tags, got %d", len(output.Checks[1].Tags))
	}

	// Check lint without tags (should be empty, not nil)
	if output.Checks[2].ID != "lint" {
		t.Errorf("expected check[2] id 'lint', got %q", output.Checks[2].ID)
	}
	if output.Checks[2].Tags == nil && len(output.Checks[2].Tags) != 0 {
		t.Errorf("expected check[2] to have empty tags, got %v", output.Checks[2].Tags)
	}
}
