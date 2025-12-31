package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter vibeguard.yaml",
	Long: `Create a starter configuration file in the current directory.

The generated file includes common checks for Go projects as a starting point.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// TODO: Implement init command
	// This will be implemented in subsequent tasks
	fmt.Println("Creating vibeguard.yaml...")
	return nil
}
