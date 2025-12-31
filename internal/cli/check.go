package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/executor"
	"github.com/vibeguard/vibeguard/internal/orchestrator"
	"github.com/vibeguard/vibeguard/internal/output"
)

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

	// Create formatter
	formatter := output.New(os.Stdout, verbose)

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

	// Format and output results
	if jsonOutput {
		output.FormatJSON(os.Stdout, result)
	} else {
		formatter.FormatResult(result)
	}

	// Exit with appropriate code
	if result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}

	return nil
}
