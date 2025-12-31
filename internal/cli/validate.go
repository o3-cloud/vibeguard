package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/config"
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
	// Load and validate configuration (Load already validates)
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Print validation success
	fmt.Printf("Configuration is valid (%d checks defined)\n", len(cfg.Checks))

	if verbose {
		fmt.Println("\nChecks:")
		for _, check := range cfg.Checks {
			deps := ""
			if len(check.Requires) > 0 {
				deps = fmt.Sprintf(" (requires: %v)", check.Requires)
			}
			fmt.Printf("  - %s: %s%s\n", check.ID, check.Severity, deps)
		}
	}

	return nil
}
