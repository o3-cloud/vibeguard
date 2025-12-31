// Package cli implements the vibeguard command-line interface.
package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Flags
	configFile string
	verbose    bool
	jsonOutput bool
	parallel   int
	failFast   bool
)

// rootCmd is the base command for vibeguard
var rootCmd = &cobra.Command{
	Use:   "vibeguard",
	Short: "A unified code quality and policy enforcement tool",
	Long: `VibeGuard is a unified code quality and policy enforcement tool designed
for CI/CD pipelines and AI agent workflows.

It orchestrates external tools, evaluates assertions against their output,
and emits actionable signals only when violations occur.

Core principles:
  - Silence is Success: No output when all checks pass
  - Tools are Commands: Any CLI tool integrates via shell commands
  - Simple by Default: Exit code pass/fail requires minimal config
  - Actionable Output: Every violation answers "What failed?" and "What should I do?"`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show all check results, not just failures")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().IntVarP(&parallel, "parallel", "p", 4, "Max parallel checks")
	rootCmd.PersistentFlags().BoolVar(&failFast, "fail-fast", false, "Stop on first failure")
}
