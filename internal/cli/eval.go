package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vibeguard/vibeguard/internal/policy"
	"github.com/vibeguard/vibeguard/pkg/models"
)

var evalCmd = &cobra.Command{
	Use:   "eval [policy-file] [resource-type]",
	Short: "Evaluate a policy against a resource",
	Long: `Evaluate a policy file against a resource.

Example:
  vibeguard eval policy.yaml commit`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		policyFile := args[0]
		resourceType := args[1]

		// Load the policy
		loader := policy.NewYAMLLoader()
		p, err := loader.Load(policyFile)
		if err != nil {
			return fmt.Errorf("failed to load policy: %w", err)
		}

		// Create a simple resource to evaluate
		resource := &models.Resource{
			Type: resourceType,
			Data: map[string]interface{}{},
		}

		// Evaluate the policy
		runner := policy.NewSimpleRunner()
		result, err := runner.Evaluate(context.Background(), p, resource)
		if err != nil {
			return fmt.Errorf("evaluation failed: %w", err)
		}

		// Output the result as JSON
		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(evalCmd)
}
