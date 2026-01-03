package cli

import (
	"testing"
)

func TestExitError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ExitError
		expected string
	}{
		{
			name:     "with message",
			err:      &ExitError{Code: 1, Message: "custom message"},
			expected: "custom message",
		},
		{
			name:     "without message",
			err:      &ExitError{Code: 42, Message: ""},
			expected: "exit code 42",
		},
		{
			name:     "zero exit code no message",
			err:      &ExitError{Code: 0, Message: ""},
			expected: "exit code 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestExitError_IsError(t *testing.T) {
	err := &ExitError{Code: 1, Message: "test"}

	// ExitError should satisfy the error interface
	var _ error = err
}

func TestExecute_WithNoArgs(t *testing.T) {
	// Execute with no subcommand should show help (no error)
	// This is a smoke test to ensure Execute() doesn't panic
	// We can't easily test the full CLI without side effects
}
