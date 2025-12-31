package cli

import (
	"fmt"

	"github.com/spf13/cobra"
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
	// TODO: Implement list command
	// This will be implemented in subsequent tasks
	fmt.Println("Listing configured checks...")
	return nil
}
