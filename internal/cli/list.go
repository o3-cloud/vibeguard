package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured checks",
	Long: `List all checks defined in the configuration file.

This command shows the check IDs, their commands, and dependencies.`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	fmt.Printf("Checks (%d):\n\n", len(cfg.Checks))

	for _, check := range cfg.Checks {
		fmt.Printf("  %s\n", check.ID)

		if verbose {
			fmt.Printf("    Command:  %s\n", check.Run)
			fmt.Printf("    Severity: %s\n", check.Severity)
			fmt.Printf("    Timeout:  %s\n", check.Timeout.AsDuration())
			if len(check.Requires) > 0 {
				fmt.Printf("    Requires: %s\n", strings.Join(check.Requires, ", "))
			}
			if check.Suggestion != "" {
				fmt.Printf("    Suggestion: %s\n", check.Suggestion)
			}
			fmt.Println()
		}
	}

	return nil
}
