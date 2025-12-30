package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vibeguard",
	Short: "VibeGuard - Policy enforcement system for CI/CD and agent workflows",
	Long: `VibeGuard is a lightweight, composable policy enforcement tool designed to integrate
seamlessly with CI/CD pipelines, agent loops, and Cloud Code workflows.

It provides declarative policy definition, judge integration, and flexible enforcement patterns.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Version flag
	rootCmd.Version = "0.1.0"
}
