package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/executor"
	"github.com/vibeguard/vibeguard/internal/orchestrator"
	"github.com/vibeguard/vibeguard/internal/output"
)

// ExitError represents an operation that completed but needs a specific exit code.
type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("exit code %d", e.Code)
}

var checkCmd = &cobra.Command{
	Use:   "check [id]",
	Short: "Run checks",
	Long: `Run all configured checks or a specific check by ID.

Examples:
  vibeguard check           Run all checks
  vibeguard check fmt       Run only the 'fmt' check
  vibeguard check -v        Run all checks with verbose output`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	// Create executor and orchestrator
	exec := executor.New("")
	orch := orchestrator.New(cfg, exec, parallel, failFast, verbose)

	// Create formatter - use stderr for Claude Code hook visibility
	formatter := output.New(os.Stderr, verbose)

	// Run checks
	ctx := context.Background()
	var result *orchestrator.RunResult

	if len(args) > 0 {
		// Run single check
		result, err = orch.RunCheck(ctx, args[0])
	} else {
		// Run all checks
		result, err = orch.Run(ctx)
	}
	if err != nil {
		return err
	}

	// Format and output results - use stderr for Claude Code hook visibility
	if jsonOutput {
		if err := output.FormatJSON(os.Stderr, result); err != nil {
			return err
		}
	} else {
		formatter.FormatResult(result)
	}

	// Exit with appropriate code if needed
	// We return an error with the appropriate exit code wrapping
	if result.ExitCode != 0 {
		return &ExitError{Code: result.ExitCode}
	}

	return nil
}
