package output

import (
	"bytes"
	"testing"
	"time"

	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/executor"
	"github.com/vibeguard/vibeguard/internal/orchestrator"
)

// Use exit code constants for consistency

func TestFormatter_QuietMode_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check:     &config.Check{ID: "fmt"},
				Execution: &executor.Result{Duration: 100 * time.Millisecond},
				Passed:    true,
			},
			{
				Check:     &config.Check{ID: "vet"},
				Execution: &executor.Result{Duration: 200 * time.Millisecond},
				Passed:    true,
			},
		},
		Violations: nil,
		ExitCode:   0,
	}

	f.FormatResult(result)

	// Quiet mode with no violations should produce no output ("silence is success")
	if buf.Len() != 0 {
		t.Errorf("expected no output in quiet mode with no violations, got: %q", buf.String())
	}
}

func TestFormatter_QuietMode_WithViolations(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

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
			},
		},
		ExitCode: executor.ExitCodeViolation,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show the violation
	if !bytes.Contains(buf.Bytes(), []byte("FAIL  coverage (error)")) {
		t.Errorf("expected violation header, got: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("go test -cover ./...")) {
		t.Errorf("expected command in output, got: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Tip: Coverage is 72%, need 80%.")) {
		t.Errorf("expected suggestion in output, got: %q", output)
	}
}

func TestFormatter_VerboseMode_AllPassing(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, true) // verbose mode

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

	f.FormatResult(result)

	output := buf.String()

	// Should show all checks with timing
	if !bytes.Contains(buf.Bytes(), []byte("✓ fmt")) {
		t.Errorf("expected fmt check in output, got: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("passed (0.1s)")) {
		t.Errorf("expected timing in output, got: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("✓ vet")) {
		t.Errorf("expected vet check in output, got: %q", output)
	}
}

func TestFormatter_VerboseMode_WithFailure(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, true) // verbose mode

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check:     &config.Check{ID: "fmt"},
				Execution: &executor.Result{Duration: 100 * time.Millisecond},
				Passed:    true,
			},
			{
				Check: &config.Check{
					ID:         "coverage",
					Suggestion: "Coverage is 72%, need 80%",
				},
				Execution: &executor.Result{Duration: 900 * time.Millisecond},
				Passed:    false,
				Extracted: map[string]string{},
			},
		},
		Violations: []*orchestrator.Violation{
			{
				CheckID:    "coverage",
				Severity:   config.SeverityError,
				Command:    "go test -cover ./...",
				Suggestion: "Coverage is 72%, need 80%",
			},
		},
		ExitCode: executor.ExitCodeViolation,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show passing check
	if !bytes.Contains(buf.Bytes(), []byte("✓ fmt")) {
		t.Errorf("expected fmt check in output, got: %q", output)
	}

	// Should show failing check with FAIL
	if !bytes.Contains(buf.Bytes(), []byte("✗ coverage")) {
		t.Errorf("expected coverage check failure marker, got: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("FAIL")) {
		t.Errorf("expected FAIL in output, got: %q", output)
	}

	// Should show suggestion
	if !bytes.Contains(buf.Bytes(), []byte("Coverage is 72%, need 80%")) {
		t.Errorf("expected suggestion in output, got: %q", output)
	}
}

func TestTruncateCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "short command",
			input:    "go test ./...",
			expected: "go test ./...",
		},
		{
			name:     "long command gets truncated",
			input:    "some very long command that exceeds sixty characters and should be truncated",
			expected: "some very long command that exceeds sixty characters and ...",
		},
		{
			name:     "multiline command collapses to single line",
			input:    "echo hello\necho world",
			expected: "echo hello echo world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateCommand(tt.input)
			if got != tt.expected {
				t.Errorf("truncateCommand(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
