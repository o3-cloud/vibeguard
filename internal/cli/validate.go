package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Long: `Validate the vibeguard configuration file without running any checks.

This command is useful for CI/CD pipelines to catch configuration errors early.`,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	// TODO: Implement validate command
	// This will be implemented in subsequent tasks
	fmt.Println("Validating configuration...")
	return nil
}
