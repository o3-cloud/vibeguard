package output

import (
	"encoding/json"
	"io"

	"github.com/vibeguard/vibeguard/internal/orchestrator"
)

// JSONOutput represents the JSON output format.
type JSONOutput struct {
	Checks            []JSONCheck     `json:"checks"`
	Violations        []JSONViolation `json:"violations"`
	ExitCode          int             `json:"exit_code"`
	FailFastTriggered bool            `json:"fail_fast_triggered,omitempty"`
}

// JSONCheck represents a check result in JSON format.
type JSONCheck struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	DurationMS int64  `json:"duration_ms"`
}

// JSONViolation represents a violation in JSON format.
type JSONViolation struct {
	ID         string            `json:"id"`
	Severity   string            `json:"severity"`
	Command    string            `json:"command"`
	Suggestion string            `json:"suggestion,omitempty"`
	Fix        string            `json:"fix,omitempty"`
	Extracted  map[string]string `json:"extracted,omitempty"`
	LogFile    string            `json:"log_file,omitempty"`
}

// FormatJSON outputs the result in JSON format.
func FormatJSON(out io.Writer, result *orchestrator.RunResult) error {
	output := JSONOutput{
		Checks:            make([]JSONCheck, 0, len(result.Results)),
		Violations:        make([]JSONViolation, 0, len(result.Violations)),
		ExitCode:          result.ExitCode,
		FailFastTriggered: result.FailFastTriggered,
	}

	for _, r := range result.Results {
		status := "passed"
		if r.Execution.Cancelled {
			status = "cancelled"
		} else if !r.Passed {
			status = "failed"
		}
		output.Checks = append(output.Checks, JSONCheck{
			ID:         r.Check.ID,
			Status:     status,
			DurationMS: r.Execution.Duration.Milliseconds(),
		})
	}

	for _, v := range result.Violations {
		output.Violations = append(output.Violations, JSONViolation{
			ID:         v.CheckID,
			Severity:   string(v.Severity),
			Command:    v.Command,
			Suggestion: v.Suggestion,
			Fix:        v.Fix,
			Extracted:  v.Extracted,
			LogFile:    v.LogFile,
		})
	}

	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
