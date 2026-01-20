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
		ExitCode: 1, // Default error exit code
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show the violation with FAIL header
	if !bytes.Contains(buf.Bytes(), []byte("FAIL  coverage (error)")) {
		t.Errorf("expected violation header, got: %q", output)
	}
	// Should show the suggestion directly (no "Tip:" prefix in new format)
	if !bytes.Contains(buf.Bytes(), []byte("Coverage is 72%, need 80%.")) {
		t.Errorf("expected suggestion in output, got: %q", output)
	}
	// Should show advisory line
	if !bytes.Contains(buf.Bytes(), []byte("Advisory: blocks commit")) {
		t.Errorf("expected advisory in output, got: %q", output)
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
		ExitCode: 1, // Default error exit code
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

	// Should show advisory line
	if !bytes.Contains(buf.Bytes(), []byte("Advisory: blocks commit")) {
		t.Errorf("expected advisory in output, got: %q", output)
	}
}

func TestFormatter_VerboseMode_WithCancelledCheck(t *testing.T) {
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
				Check:     &config.Check{ID: "test"},
				Execution: &executor.Result{Duration: 50 * time.Millisecond, Cancelled: true},
				Passed:    false,
			},
		},
		Violations: nil,
		ExitCode:   0,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show cancelled check
	if !bytes.Contains(buf.Bytes(), []byte("⊘ test")) {
		t.Errorf("expected test check cancelled marker, got: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("cancelled")) {
		t.Errorf("expected cancelled in output, got: %q", output)
	}
}

func TestFormatter_VerboseMode_FailFastTriggered(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, true) // verbose mode

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check:     &config.Check{ID: "lint"},
				Execution: &executor.Result{Duration: 100 * time.Millisecond},
				Passed:    false,
			},
		},
		Violations: []*orchestrator.Violation{
			{
				CheckID:    "lint",
				Severity:   config.SeverityError,
				Command:    "golangci-lint run ./...",
				Suggestion: "Lint errors found",
			},
		},
		FailFastTriggered: true,
		ExitCode:          1, // Default error exit code
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show fail-fast message
	if !bytes.Contains(buf.Bytes(), []byte("Execution stopped early due to --fail-fast")) {
		t.Errorf("expected fail-fast message in output, got: %q", output)
	}
}

func TestFormatter_VerboseMode_WarningViolation(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, true) // verbose mode

	result := &orchestrator.RunResult{
		Results: []*orchestrator.CheckResult{
			{
				Check: &config.Check{
					ID:       "test",
					Severity: config.SeverityWarning,
				},
				Execution: &executor.Result{Duration: 500 * time.Millisecond},
				Passed:    false,
			},
		},
		Violations: []*orchestrator.Violation{
			{
				CheckID:    "test",
				Severity:   config.SeverityWarning,
				Command:    "go test ./...",
				Suggestion: "Some tests are failing",
				Fix:        "go test ./...",
			},
		},
		ExitCode: 0,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show WARN header for warning severity
	if !bytes.Contains(buf.Bytes(), []byte("WARN")) {
		t.Errorf("expected WARN in output, got: %q", output)
	}

	// Should show Fix line
	if !bytes.Contains(buf.Bytes(), []byte("Fix: go test ./...")) {
		t.Errorf("expected Fix in output, got: %q", output)
	}

	// Should show advisory line for warning
	if !bytes.Contains(buf.Bytes(), []byte("Advisory: does not block commit")) {
		t.Errorf("expected warning advisory in output, got: %q", output)
	}
}

func TestFormatter_QuietMode_FailFastTriggered(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

	result := &orchestrator.RunResult{
		Violations: []*orchestrator.Violation{
			{
				CheckID:    "lint",
				Severity:   config.SeverityError,
				Command:    "golangci-lint run ./...",
				Suggestion: "Lint errors found",
			},
		},
		FailFastTriggered: true,
		ExitCode:          1, // Default error exit code
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show fail-fast message
	if !bytes.Contains(buf.Bytes(), []byte("Execution stopped early due to --fail-fast")) {
		t.Errorf("expected fail-fast message in output, got: %q", output)
	}
}

func TestFormatter_QuietMode_WithTriggeredPrompts(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

	result := &orchestrator.RunResult{
		Violations: []*orchestrator.Violation{
			{
				CheckID:    "vet",
				Severity:   config.SeverityError,
				Command:    "go vet ./...",
				Suggestion: "Fix vet errors",
				TriggeredPrompts: []*orchestrator.TriggeredPrompt{
					{
						Event:   "failure",
						Source:  "init",
						Content: "You are an expert in helping users set up VibeGuard.\nGuide them through setup.",
					},
					{
						Event:   "failure",
						Source:  "inline",
						Content: "Also remember to run gofmt before committing",
					},
				},
			},
		},
		ExitCode: 1,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show violation header
	if !bytes.Contains(buf.Bytes(), []byte("FAIL  vet (error)")) {
		t.Errorf("expected violation header, got: %q", output)
	}

	// Should show triggered prompts section
	if !bytes.Contains(buf.Bytes(), []byte("Triggered Prompts (failure):")) {
		t.Errorf("expected triggered prompts section, got: %q", output)
	}

	// Should show first prompt with ID
	if !bytes.Contains(buf.Bytes(), []byte("[1] init:")) {
		t.Errorf("expected first prompt header, got: %q", output)
	}

	// Should show prompt content
	if !bytes.Contains(buf.Bytes(), []byte("You are an expert in helping users set up VibeGuard.")) {
		t.Errorf("expected prompt content, got: %q", output)
	}

	// Should show second prompt with inline label
	if !bytes.Contains(buf.Bytes(), []byte("[2] (inline):")) {
		t.Errorf("expected inline prompt header, got: %q", output)
	}

	// Should show inline content
	if !bytes.Contains(buf.Bytes(), []byte("Also remember to run gofmt before committing")) {
		t.Errorf("expected inline content, got: %q", output)
	}

	// Should show advisory line
	if !bytes.Contains(buf.Bytes(), []byte("Advisory: blocks commit")) {
		t.Errorf("expected advisory line, got: %q", output)
	}
}

func TestFormatter_TriggeredPrompts_SinglePrompt(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

	result := &orchestrator.RunResult{
		Violations: []*orchestrator.Violation{
			{
				CheckID:  "test",
				Severity: config.SeverityError,
				Command:  "go test ./...",
				TriggeredPrompts: []*orchestrator.TriggeredPrompt{
					{
						Event:   "failure",
						Source:  "test-generator",
						Content: "Generate comprehensive unit tests for the failing code.",
					},
				},
			},
		},
		ExitCode: 1,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show the single prompt
	if !bytes.Contains(buf.Bytes(), []byte("[1] test-generator:")) {
		t.Errorf("expected prompt header, got: %q", output)
	}
}

func TestFormatter_TriggeredPrompts_MultilineContent(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

	result := &orchestrator.RunResult{
		Violations: []*orchestrator.Violation{
			{
				CheckID:  "security",
				Severity: config.SeverityError,
				Command:  "gosec ./...",
				TriggeredPrompts: []*orchestrator.TriggeredPrompt{
					{
						Event:   "failure",
						Source:  "security-audit",
						Content: "Line 1: Check for vulnerabilities\nLine 2: Including injection attacks\nLine 3: And access control issues",
					},
				},
			},
		},
		ExitCode: 1,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show multi-line content with proper indentation
	if !bytes.Contains(buf.Bytes(), []byte("      Line 1: Check for vulnerabilities")) {
		t.Errorf("expected first line indented, got: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("      Line 2: Including injection attacks")) {
		t.Errorf("expected second line indented, got: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("      Line 3: And access control issues")) {
		t.Errorf("expected third line indented, got: %q", output)
	}
}

func TestFormatter_TriggeredPrompts_SuccessEvent(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

	result := &orchestrator.RunResult{
		Violations: []*orchestrator.Violation{
			{
				CheckID:  "fmt",
				Severity: config.SeverityError,
				Command:  "go fmt ./...",
				TriggeredPrompts: []*orchestrator.TriggeredPrompt{
					{
						Event:   "success",
						Source:  "code-review",
						Content: "Great! Code formatting passed. Consider code review suggestions.",
					},
				},
			},
		},
		ExitCode: 1,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show success event label
	if !bytes.Contains(buf.Bytes(), []byte("Triggered Prompts (success):")) {
		t.Errorf("expected success event label, got: %q", output)
	}
}

func TestFormatter_TriggeredPrompts_TimeoutEvent(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

	result := &orchestrator.RunResult{
		Violations: []*orchestrator.Violation{
			{
				CheckID:  "integration-test",
				Severity: config.SeverityError,
				Command:  "go test ./tests/integration/...",
				Timedout: true,
				TriggeredPrompts: []*orchestrator.TriggeredPrompt{
					{
						Event:   "timeout",
						Source:  "inline",
						Content: "Integration tests timed out. Check for hanging goroutines.",
					},
				},
			},
		},
		ExitCode: 1,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should show timeout event label
	if !bytes.Contains(buf.Bytes(), []byte("Triggered Prompts (timeout):")) {
		t.Errorf("expected timeout event label, got: %q", output)
	}
}

func TestFormatter_TriggeredPrompts_NoPrompts(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, false) // quiet mode

	result := &orchestrator.RunResult{
		Violations: []*orchestrator.Violation{
			{
				CheckID:          "vet",
				Severity:         config.SeverityError,
				Command:          "go vet ./...",
				TriggeredPrompts: []*orchestrator.TriggeredPrompt{}, // Empty
			},
		},
		ExitCode: 1,
	}

	f.FormatResult(result)

	output := buf.String()

	// Should NOT show triggered prompts section when empty
	if bytes.Contains(buf.Bytes(), []byte("Triggered Prompts")) {
		t.Errorf("expected no triggered prompts section, got: %q", output)
	}

	// Should still show basic violation info
	if !bytes.Contains(buf.Bytes(), []byte("FAIL  vet (error)")) {
		t.Errorf("expected violation header, got: %q", output)
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

// TestFormatViolation_OutputCombinations tests all combinations of suggestion/fix/severity
// as specified in the implementation plan.
func TestFormatViolation_OutputCombinations(t *testing.T) {
	tests := []struct {
		name              string
		violation         *orchestrator.Violation
		expectSuggestion  bool
		expectFix         bool
		expectFixFallback bool // expect "Fix: <run command>" fallback
		expectWarn        bool
		expectFail        bool
		expectedAdvisory  string
	}{
		{
			name: "suggestion set, fix empty, severity error",
			violation: &orchestrator.Violation{
				CheckID:    "test-check",
				Severity:   config.SeverityError,
				Command:    "go test ./...",
				Suggestion: "Add more tests",
				Fix:        "",
			},
			expectSuggestion:  true,
			expectFix:         false,
			expectFixFallback: false,
			expectFail:        true,
			expectedAdvisory:  "Advisory: blocks commit",
		},
		{
			name: "suggestion empty, fix set, severity error",
			violation: &orchestrator.Violation{
				CheckID:    "test-check",
				Severity:   config.SeverityError,
				Command:    "go test ./...",
				Suggestion: "",
				Fix:        "Run: go test ./...",
			},
			expectSuggestion:  false,
			expectFix:         true,
			expectFixFallback: false,
			expectFail:        true,
			expectedAdvisory:  "Advisory: blocks commit",
		},
		{
			name: "suggestion set, fix set, severity error",
			violation: &orchestrator.Violation{
				CheckID:    "test-check",
				Severity:   config.SeverityError,
				Command:    "go test ./...",
				Suggestion: "Coverage is 72%, need 80%",
				Fix:        "Add tests to improve coverage",
			},
			expectSuggestion:  true,
			expectFix:         true,
			expectFixFallback: false,
			expectFail:        true,
			expectedAdvisory:  "Advisory: blocks commit",
		},
		{
			name: "suggestion set, fix set, severity warning",
			violation: &orchestrator.Violation{
				CheckID:    "test-check",
				Severity:   config.SeverityWarning,
				Command:    "go test ./...",
				Suggestion: "Coverage is 72%, need 80%",
				Fix:        "Add tests to improve coverage",
			},
			expectSuggestion:  true,
			expectFix:         true,
			expectFixFallback: false,
			expectWarn:        true,
			expectedAdvisory:  "Advisory: does not block commit",
		},
		{
			name: "suggestion empty, fix empty, severity error (fallback)",
			violation: &orchestrator.Violation{
				CheckID:    "test-check",
				Severity:   config.SeverityError,
				Command:    "go test ./...",
				Suggestion: "",
				Fix:        "",
			},
			expectSuggestion:  false,
			expectFix:         false,
			expectFixFallback: true,
			expectFail:        true,
			expectedAdvisory:  "Advisory: blocks commit",
		},
		{
			name: "timeout violation",
			violation: &orchestrator.Violation{
				CheckID:    "test-check",
				Severity:   config.SeverityError,
				Command:    "go test ./...",
				Suggestion: "Check timed out. Consider increasing the timeout value or optimizing the command.",
				Fix:        "",
				Timedout:   true,
			},
			expectSuggestion:  true,
			expectFix:         false,
			expectFixFallback: false,
			expectFail:        true,
			expectedAdvisory:  "Advisory: blocks commit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := New(&buf, false) // quiet mode

			result := &orchestrator.RunResult{
				Violations: []*orchestrator.Violation{tt.violation},
			}
			f.FormatResult(result)

			output := buf.String()

			// Check WARN/FAIL header
			if tt.expectWarn && !bytes.Contains(buf.Bytes(), []byte("WARN")) {
				t.Errorf("expected WARN header, got: %q", output)
			}
			if tt.expectFail && !bytes.Contains(buf.Bytes(), []byte("FAIL")) {
				t.Errorf("expected FAIL header, got: %q", output)
			}

			// Check suggestion
			if tt.expectSuggestion && tt.violation.Suggestion != "" {
				if !bytes.Contains(buf.Bytes(), []byte(tt.violation.Suggestion)) {
					t.Errorf("expected suggestion %q in output, got: %q", tt.violation.Suggestion, output)
				}
			}

			// Check fix
			if tt.expectFix && tt.violation.Fix != "" {
				expected := "Fix: " + tt.violation.Fix
				if !bytes.Contains(buf.Bytes(), []byte(expected)) {
					t.Errorf("expected fix %q in output, got: %q", expected, output)
				}
			}

			// Check fallback to run command
			if tt.expectFixFallback {
				expected := "Fix: " + tt.violation.Command
				if !bytes.Contains(buf.Bytes(), []byte(expected)) {
					t.Errorf("expected fallback fix %q in output, got: %q", expected, output)
				}
			}

			// Check advisory line
			if !bytes.Contains(buf.Bytes(), []byte(tt.expectedAdvisory)) {
				t.Errorf("expected advisory %q in output, got: %q", tt.expectedAdvisory, output)
			}

			// Check timeout status in header
			if tt.violation.Timedout {
				if !bytes.Contains(buf.Bytes(), []byte("(timeout)")) {
					t.Errorf("expected (timeout) in output, got: %q", output)
				}
			}
		})
	}
}
