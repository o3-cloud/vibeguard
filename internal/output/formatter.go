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
}

// formatVerbose outputs all check results.
func (f *Formatter) formatVerbose(result *orchestrator.RunResult) {
	for _, r := range result.Results {
		if r.Passed {
			fmt.Fprintf(f.out, "✓ %-15s passed (%.1fs)\n",
				r.Check.ID, r.Execution.Duration.Seconds())
		} else {
			fmt.Fprintf(f.out, "✗ %-15s FAIL (%.1fs)\n",
				r.Check.ID, r.Execution.Duration.Seconds())
			// Show suggestion if present
			if r.Check.Suggestion != "" {
				suggestion := config.InterpolateWithExtracted(
					r.Check.Suggestion,
					nil,
					r.Extracted,
				)
				fmt.Fprintf(f.out, "  %s\n", suggestion)
			}
		}
	}
}

// formatViolation outputs a single violation.
func (f *Formatter) formatViolation(v *orchestrator.Violation) {
	fmt.Fprintf(f.out, "FAIL  %s (%s)\n", v.CheckID, v.Severity)
	fmt.Fprintf(f.out, "      > %s\n", truncateCommand(v.Command))
	if v.Suggestion != "" {
		fmt.Fprintf(f.out, "\n      Tip: %s\n", v.Suggestion)
	}
	fmt.Fprintln(f.out)
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
