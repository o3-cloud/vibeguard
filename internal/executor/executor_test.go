// Package executor provides command execution capabilities for checks.
package executor

import (
	"context"
	"testing"
	"time"
)

func TestExecute_ExitCodeZero_Success(t *testing.T) {
	exec := New("")

	result, err := exec.Execute(context.Background(), "test-pass", "exit 0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if !result.Success {
		t.Error("expected Success to be true for exit code 0")
	}
	if result.CheckID != "test-pass" {
		t.Errorf("expected CheckID 'test-pass', got %q", result.CheckID)
	}
}

func TestExecute_ExitCodeNonZero_Failure(t *testing.T) {
	exec := New("")

	result, err := exec.Execute(context.Background(), "test-fail", "exit 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
	if result.Success {
		t.Error("expected Success to be false for exit code 1")
	}
}

func TestExecute_VariousExitCodes(t *testing.T) {
	tests := []struct {
		name     string
		exitCode int
		success  bool
	}{
		{"exit 0 is success", 0, true},
		{"exit 1 is failure", 1, false},
		{"exit 2 is failure", 2, false},
		{"exit 42 is failure", 42, false},
		{"exit 127 is failure", 127, false},
		{"exit 255 is failure", 255, false},
	}

	exec := New("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use shell directly to handle larger exit codes
			result, err := exec.Execute(context.Background(), "test", "exit "+itoa(tt.exitCode))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.ExitCode != tt.exitCode {
				t.Errorf("expected exit code %d, got %d", tt.exitCode, result.ExitCode)
			}
			if result.Success != tt.success {
				t.Errorf("expected Success to be %v for exit code %d", tt.success, tt.exitCode)
			}
		})
	}
}

// itoa converts an int to a string without importing strconv
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func TestExecute_CapturesStdout(t *testing.T) {
	exec := New("")

	result, err := exec.Execute(context.Background(), "test-stdout", "echo 'hello world'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "hello world\n"
	if result.Stdout != expected {
		t.Errorf("expected stdout %q, got %q", expected, result.Stdout)
	}
}

func TestExecute_CapturesStderr(t *testing.T) {
	exec := New("")

	result, err := exec.Execute(context.Background(), "test-stderr", "echo 'error message' >&2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "error message\n"
	if result.Stderr != expected {
		t.Errorf("expected stderr %q, got %q", expected, result.Stderr)
	}
}

func TestExecute_CombinedOutput(t *testing.T) {
	exec := New("")

	result, err := exec.Execute(context.Background(), "test-combined", "echo 'out'; echo 'err' >&2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Combined is stdout + stderr (order may vary due to buffering)
	if result.Combined != "out\nerr\n" && result.Combined != "err\nout\n" {
		// Combined appends stdout first, then stderr
		expected := "out\nerr\n"
		if result.Combined != expected {
			t.Errorf("expected combined %q, got %q", expected, result.Combined)
		}
	}
}

func TestExecute_TracksDuration(t *testing.T) {
	exec := New("")

	result, err := exec.Execute(context.Background(), "test-duration", "sleep 0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Duration should be at least 100ms
	if result.Duration < 100*time.Millisecond {
		t.Errorf("expected duration >= 100ms, got %v", result.Duration)
	}
}

func TestExecute_ContextCancellation(t *testing.T) {
	exec := New("")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result, err := exec.Execute(ctx, "test-timeout", "sleep 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Command should have been killed, resulting in non-zero exit
	if result.Success {
		t.Error("expected Success to be false for timed-out command")
	}
}

func TestExecute_Timeout_SetsExitCode3(t *testing.T) {
	exec := New("")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result, err := exec.Execute(ctx, "test-timeout", "sleep 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Timeout should set exit code 3
	if result.ExitCode != ExitCodeTimeout {
		t.Errorf("expected exit code %d for timeout, got %d", ExitCodeTimeout, result.ExitCode)
	}
}

func TestExecute_Timeout_SetsTimedoutFlag(t *testing.T) {
	exec := New("")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result, err := exec.Execute(ctx, "test-timeout", "sleep 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Timedout flag should be true
	if !result.Timedout {
		t.Error("expected Timedout to be true for timed-out command")
	}
}

func TestExecute_NoTimeout_TimedoutFlagFalse(t *testing.T) {
	exec := New("")

	result, err := exec.Execute(context.Background(), "test-fast", "echo hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Fast command should not have Timedout flag
	if result.Timedout {
		t.Error("expected Timedout to be false for completed command")
	}
}

func TestExecute_NonZeroExit_NotTimeout(t *testing.T) {
	exec := New("")

	result, err := exec.Execute(context.Background(), "test-fail", "exit 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Non-zero exit without timeout should keep original exit code
	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
	if result.Timedout {
		t.Error("expected Timedout to be false for non-timeout failure")
	}
}

func TestResult_String_Timeout(t *testing.T) {
	result := &Result{
		CheckID:  "test-check",
		ExitCode: ExitCodeTimeout,
		Duration: 30 * time.Second,
		Success:  false,
		Timedout: true,
	}

	expected := "test-check: timeout (exit=3, duration=30s)"
	got := result.String()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestExecute_NonZeroExitNotAnError(t *testing.T) {
	exec := New("")

	// This is the key test: non-zero exit codes should NOT return an error
	// They should only set Success=false and ExitCode
	result, err := exec.Execute(context.Background(), "test-no-error", "exit 42")
	if err != nil {
		t.Fatalf("non-zero exit should not return error, got: %v", err)
	}

	if result.Error != nil {
		t.Errorf("result.Error should be nil for exit code failure, got: %v", result.Error)
	}
	if result.ExitCode != 42 {
		t.Errorf("expected exit code 42, got %d", result.ExitCode)
	}
	if result.Success {
		t.Error("expected Success to be false")
	}
}

func TestResult_String(t *testing.T) {
	tests := []struct {
		name     string
		result   *Result
		expected string
	}{
		{
			name: "passed result",
			result: &Result{
				CheckID:  "test-check",
				ExitCode: 0,
				Duration: 100 * time.Millisecond,
				Success:  true,
			},
			expected: "test-check: passed (exit=0, duration=100ms)",
		},
		{
			name: "failed result",
			result: &Result{
				CheckID:  "test-check",
				ExitCode: 1,
				Duration: 200 * time.Millisecond,
				Success:  false,
			},
			expected: "test-check: failed (exit=1, duration=200ms)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.String()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestExecute_RealCommandPass(t *testing.T) {
	exec := New("")

	// Test with a real command that should pass
	result, err := exec.Execute(context.Background(), "true-cmd", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	if !result.Success {
		t.Error("expected Success to be true for 'true' command")
	}
}

func TestExecute_RealCommandFail(t *testing.T) {
	exec := New("")

	// Test with a real command that should fail
	result, err := exec.Execute(context.Background(), "false-cmd", "false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
	if result.Success {
		t.Error("expected Success to be false for 'false' command")
	}
}

func TestExecute_CommandWithOutput_StillReportsExitCode(t *testing.T) {
	exec := New("")

	// A command that produces output but still fails
	result, err := exec.Execute(context.Background(), "output-fail", "echo 'some output' && exit 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ExitCode != 5 {
		t.Errorf("expected exit code 5, got %d", result.ExitCode)
	}
	if result.Success {
		t.Error("expected Success to be false")
	}
	if result.Stdout != "some output\n" {
		t.Errorf("expected stdout 'some output\\n', got %q", result.Stdout)
	}
}

func TestNew_DefaultsToCurrentDir(t *testing.T) {
	exec := New("")

	// Should not panic and workDir should be set
	if exec.workDir == "" {
		t.Error("workDir should default to current directory, not empty")
	}
}
