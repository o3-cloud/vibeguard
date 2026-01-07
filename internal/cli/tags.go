package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/config"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all unique tags in the configuration",
	Long: `List all unique tags found in the vibeguard configuration.

Tags are used to filter which checks run. This command shows all unique tags
that are used by checks in the configuration file.

Examples:
  vibeguard tags           List all tags
  vibeguard tags -c config.yaml  List tags from specific config`,
	RunE: runTags,
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}

func runTags(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	// Collect unique tags from all checks
	tagSet := make(map[string]bool)
	for _, check := range cfg.Checks {
		for _, tag := range check.Tags {
			tagSet[tag] = true
		}
	}

	// Convert to sorted slice
	var tagList []string
	for tag := range tagSet {
		tagList = append(tagList, tag)
	}
	sort.Strings(tagList)

	// Output tags, one per line
	out := cmd.OutOrStdout()
	for _, tag := range tagList {
		if _, err := fmt.Fprintln(out, tag); err != nil {
			return err
		}
	}

	return nil
}
