package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/config"
)

var promptCmd = &cobra.Command{
	Use:   "prompt [prompt-id]",
	Short: "List or retrieve prompts",
	Long: `List or retrieve prompts from the configuration file.

Without a prompt ID, lists all available prompts:
  vibeguard prompt           List all prompts
  vibeguard prompt -v        List prompts with descriptions and tags
  vibeguard prompt --json    List prompts in JSON format

With a prompt ID, outputs the prompt content:
  vibeguard prompt init                    Show the init prompt
  vibeguard prompt init | less             Pipe prompt to less
  vibeguard prompt code-review | llm ...   Pipe to LLM tools

Examples:
  vibeguard prompt                    List all prompts
  vibeguard prompt -v                 List with descriptions
  vibeguard prompt --json             Machine-readable list
  vibeguard prompt init               Get init prompt content
  vibeguard prompt code-review | less Pipe to less`,
	RunE: runPrompt,
}

func init() {
	rootCmd.AddCommand(promptCmd)
}

func runPrompt(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	// If no prompt ID provided, list all prompts
	if len(args) == 0 {
		return listPromptsWithBuiltin(cmd, cfg.Prompts)
	}

	// Find and retrieve specific prompt
	promptID := args[0]
	for _, prompt := range cfg.Prompts {
		if prompt.ID == promptID {
			// Output raw prompt content to stdout
			out := cmd.OutOrStdout()
			_, _ = fmt.Fprint(out, prompt.Content)
			return nil
		}
	}

	// Check for built-in prompts before failing
	if promptID == "init" {
		out := cmd.OutOrStdout()
		_, _ = fmt.Fprint(out, InitPromptContent)
		return nil
	}

	// Prompt not found
	return fmt.Errorf("prompt not found: %s", promptID)
}

// listPromptsWithBuiltin displays prompts including built-in prompts in various formats
func listPromptsWithBuiltin(cmd *cobra.Command, prompts []config.Prompt) error {
	out := cmd.OutOrStdout()

	// Check if we have any prompts at all (including built-in)
	hasConfiguredPrompts := len(prompts) > 0

	// JSON output
	if jsonOutput {
		// If we have configured prompts, output them
		if hasConfiguredPrompts {
			return outputPromptsJSON(out, prompts)
		}
		// Otherwise output built-in prompts
		return outputBuiltinPromptsJSON(out)
	}

	// Human-readable output
	if !hasConfiguredPrompts {
		// No configured prompts, show only built-in
		_, _ = fmt.Fprintf(out, "Prompts (1 built-in):\n\n")
		_, _ = fmt.Fprintf(out, "  init (built-in)\n")
		if verbose {
			_, _ = fmt.Fprintf(out, "    Description: Built-in VibeGuard setup guidance\n\n")
		}
		return nil
	}

	// Show configured prompts
	_, _ = fmt.Fprintf(out, "Prompts (%d):\n\n", len(prompts))

	for _, prompt := range prompts {
		_, _ = fmt.Fprintf(out, "  %s\n", prompt.ID)

		if verbose {
			if prompt.Description != "" {
				_, _ = fmt.Fprintf(out, "    Description: %s\n", prompt.Description)
			}
			if len(prompt.Tags) > 0 {
				_, _ = fmt.Fprintf(out, "    Tags:        %s\n", strings.Join(prompt.Tags, ", "))
			}
			_, _ = fmt.Fprintln(out)
		}
	}

	// Show built-in init prompt
	_, _ = fmt.Fprintf(out, "  init (built-in)\n")
	if verbose {
		_, _ = fmt.Fprintf(out, "    Description: Built-in VibeGuard setup guidance\n\n")
	}

	return nil
}

// outputBuiltinPromptsJSON outputs only built-in prompts in JSON format
func outputBuiltinPromptsJSON(out io.Writer) error {
	jsonPrompts := []map[string]interface{}{
		{
			"id":          "init",
			"description": "Built-in VibeGuard setup guidance",
			"built_in":    true,
		},
	}

	jsonBytes, err := json.MarshalIndent(jsonPrompts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	_, err = fmt.Fprintf(out, "%s\n", string(jsonBytes))
	return err
}

// outputPromptsJSON outputs prompts in JSON format including built-in prompts
func outputPromptsJSON(out io.Writer, prompts []config.Prompt) error {
	// Create JSON-friendly output with only essential fields (omit Content for list)
	jsonPrompts := make([]map[string]interface{}, 0, len(prompts)+1)

	for _, prompt := range prompts {
		item := map[string]interface{}{
			"id": prompt.ID,
		}
		if prompt.Description != "" {
			item["description"] = prompt.Description
		}
		if len(prompt.Tags) > 0 {
			item["tags"] = prompt.Tags
		}
		jsonPrompts = append(jsonPrompts, item)
	}

	// Add built-in init prompt
	jsonPrompts = append(jsonPrompts, map[string]interface{}{
		"id":          "init",
		"description": "Built-in VibeGuard setup guidance",
		"built_in":    true,
	})

	jsonBytes, err := json.MarshalIndent(jsonPrompts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	_, err = fmt.Fprintf(out, "%s\n", string(jsonBytes))
	return err
}
