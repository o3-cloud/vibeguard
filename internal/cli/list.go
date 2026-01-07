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

	out := cmd.OutOrStdout()
	_, _ = fmt.Fprintf(out, "Checks (%d):\n\n", len(cfg.Checks))

	for _, check := range cfg.Checks {
		_, _ = fmt.Fprintf(out, "  %s\n", check.ID)

		if verbose {
			if len(check.Tags) > 0 {
				_, _ = fmt.Fprintf(out, "    Tags:     %s\n", strings.Join(check.Tags, ", "))
			}
			_, _ = fmt.Fprintf(out, "    Command:  %s\n", check.Run)
			_, _ = fmt.Fprintf(out, "    Severity: %s\n", check.Severity)
			_, _ = fmt.Fprintf(out, "    Timeout:  %s\n", check.Timeout.AsDuration())
			if len(check.Requires) > 0 {
				_, _ = fmt.Fprintf(out, "    Requires: %s\n", strings.Join(check.Requires, ", "))
			}
			if check.Suggestion != "" {
				_, _ = fmt.Fprintf(out, "    Suggestion: %s\n", check.Suggestion)
			}
			_, _ = fmt.Fprintln(out)
		}
	}

	return nil
}
