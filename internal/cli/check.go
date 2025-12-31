package cli

import (
	"fmt"

	"github.com/spf13/cobra"
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
	// TODO: Implement check command
	// This will be implemented in subsequent tasks
	if len(args) > 0 {
		fmt.Printf("Running check: %s\n", args[0])
	} else {
		fmt.Println("Running all checks...")
	}
	return nil
}
