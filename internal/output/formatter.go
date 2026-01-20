// Package output provides output formatting for check results.
package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/orchestrator"
)

// Formatter handles output formatting.
type Formatter struct {
	out     io.Writer
	verbose bool
}

// New creates a new Formatter.
func New(out io.Writer, verbose bool) *Formatter {
	return &Formatter{
		out:     out,
		verbose: verbose,
	}
}

// FormatResult formats the run result for output.
// In quiet mode (default), only violations are shown.
// In verbose mode, all check results are shown.
func (f *Formatter) FormatResult(result *orchestrator.RunResult) {
	if f.verbose {
		f.formatVerbose(result)
	} else {
		f.formatQuiet(result)
	}
}

// formatQuiet outputs only violations (silence is success).
func (f *Formatter) formatQuiet(result *orchestrator.RunResult) {
	for _, v := range result.Violations {
		f.formatViolation(v)
	}
	if result.FailFastTriggered {
		_, _ = fmt.Fprintf(f.out, "Execution stopped early due to --fail-fast\n")
	}
}

// formatVerbose outputs all check results.
func (f *Formatter) formatVerbose(result *orchestrator.RunResult) {
	// Build a map of violations by check ID for easy lookup
	violationByID := make(map[string]*orchestrator.Violation)
	for _, v := range result.Violations {
		violationByID[v.CheckID] = v
	}

	for _, r := range result.Results {
		if r.Passed {
			_, _ = fmt.Fprintf(f.out, "✓ %-15s passed (%.1fs)\n",
				r.Check.ID, r.Execution.Duration.Seconds())
			if len(r.Check.Tags) > 0 {
				_, _ = fmt.Fprintf(f.out, "  Tags: %s\n", strings.Join(r.Check.Tags, ", "))
			}
			if len(r.TriggeredPrompts) > 0 {
				_, _ = fmt.Fprintln(f.out)
				f.formatTriggeredPrompts(r.TriggeredPrompts)
			}
		} else if r.Execution.Cancelled {
			_, _ = fmt.Fprintf(f.out, "⊘ %-15s cancelled\n", r.Check.ID)
		} else {
			// Get the violation for this check
			v := violationByID[r.Check.ID]
			if v == nil {
				// Fallback if no violation found (shouldn't happen)
				_, _ = fmt.Fprintf(f.out, "✗ %-15s FAIL (%.1fs)\n",
					r.Check.ID, r.Execution.Duration.Seconds())
				continue
			}

			// Use WARN for warning severity, FAIL for everything else
			header := "FAIL"
			if v.Severity == config.SeverityWarning {
				header = "WARN"
			}

			_, _ = fmt.Fprintf(f.out, "✗ %-15s %s (%.1fs)\n",
				r.Check.ID, header, r.Execution.Duration.Seconds())

			if len(r.Check.Tags) > 0 {
				_, _ = fmt.Fprintf(f.out, "  Tags: %s\n", strings.Join(r.Check.Tags, ", "))
			}

			// Show suggestion if present (interpolated with extracted values)
			if v.Suggestion != "" {
				suggestion := config.InterpolateWithExtracted(
					v.Suggestion,
					nil,
					v.Extracted,
				)
				_, _ = fmt.Fprintf(f.out, "  %s\n", suggestion)
			}

			// Show fix if present, otherwise fallback to run command
			if v.Fix != "" {
				fix := config.InterpolateWithExtracted(
					v.Fix,
					nil,
					v.Extracted,
				)
				_, _ = fmt.Fprintf(f.out, "  Fix: %s\n", fix)
			} else if v.Suggestion == "" {
				// Fallback: show command as fix when no suggestion and no fix
				_, _ = fmt.Fprintf(f.out, "  Fix: %s\n", v.Command)
			}

			// Show advisory line
			advisory := "blocks commit"
			if v.Severity == config.SeverityWarning {
				advisory = "does not block commit"
			}
			_, _ = fmt.Fprintf(f.out, "  Advisory: %s\n", advisory)
		}
	}
	if result.FailFastTriggered {
		_, _ = fmt.Fprintf(f.out, "\nExecution stopped early due to --fail-fast\n")
	}
}

// formatViolation outputs a single violation.
func (f *Formatter) formatViolation(v *orchestrator.Violation) {
	// Use WARN for warning severity, FAIL for everything else
	header := "FAIL"
	if v.Severity == config.SeverityWarning {
		header = "WARN"
	}

	// Format the status info (timeout vs severity)
	statusInfo := string(v.Severity)
	if v.Timedout {
		statusInfo = "timeout"
	}

	_, _ = fmt.Fprintf(f.out, "%s  %s (%s)\n\n", header, v.CheckID, statusInfo)

	// Show suggestion if present (interpolated with extracted values)
	if v.Suggestion != "" {
		suggestion := config.InterpolateWithExtracted(
			v.Suggestion,
			nil,
			v.Extracted,
		)
		_, _ = fmt.Fprintf(f.out, "  %s\n", suggestion)
	}

	// Show fix if present, otherwise fallback to run command
	if v.Fix != "" {
		fix := config.InterpolateWithExtracted(
			v.Fix,
			nil,
			v.Extracted,
		)
		_, _ = fmt.Fprintf(f.out, "  Fix: %s\n", fix)
	} else if v.Suggestion == "" {
		// Fallback: show command as fix when no suggestion and no fix
		_, _ = fmt.Fprintf(f.out, "  Fix: %s\n", v.Command)
	}

	// Show log file location if present
	if v.LogFile != "" {
		_, _ = fmt.Fprintf(f.out, "  Log: %s\n", v.LogFile)
	}

	// Show triggered prompts if present
	if len(v.TriggeredPrompts) > 0 {
		_, _ = fmt.Fprintln(f.out)
		f.formatTriggeredPrompts(v.TriggeredPrompts)
	}

	// Show advisory line
	advisory := "blocks commit"
	if v.Severity == config.SeverityWarning {
		advisory = "does not block commit"
	}
	_, _ = fmt.Fprintf(f.out, "  Advisory: %s\n", advisory)

	_, _ = fmt.Fprintln(f.out)
}

// formatTriggeredPrompts outputs triggered prompts in a formatted list.
func (f *Formatter) formatTriggeredPrompts(prompts []*orchestrator.TriggeredPrompt) {
	if len(prompts) == 0 {
		return
	}

	// Get the event type from the first prompt (all should be the same)
	eventType := prompts[0].Event

	_, _ = fmt.Fprintf(f.out, "  Triggered Prompts (%s):\n", eventType)

	for i, p := range prompts {
		// Format source: either prompt ID or "(inline)" for inline content
		source := p.Source
		if source == "inline" {
			source = "(inline)"
		}

		// Print the prompt header with number and source
		_, _ = fmt.Fprintf(f.out, "  [%d] %s:\n", i+1, source)

		// Split content by newlines and indent each line
		lines := strings.Split(p.Content, "\n")
		for _, line := range lines {
			_, _ = fmt.Fprintf(f.out, "      %s\n", line)
		}
	}
}

// truncateCommand shortens long commands for display.
func truncateCommand(cmd string) string {
	// Collapse multiline commands to single line
	cmd = strings.ReplaceAll(cmd, "\n", " ")
	cmd = strings.Join(strings.Fields(cmd), " ")

	if len(cmd) > 60 {
		return cmd[:57] + "..."
	}
	return cmd
}
