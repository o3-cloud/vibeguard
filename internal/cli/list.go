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

This command shows the check IDs, their commands, and dependencies.

Examples:
  vibeguard list           List all checks
  vibeguard list -v        List all checks with verbose output
  vibeguard list --tags security   List only security checks
  vibeguard list --exclude-tags slow   List all checks except slow ones`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringSliceVar(&tags, "tags", nil, "List only checks matching ANY of these tags (comma-separated)")
	listCmd.Flags().StringSliceVar(&excludeTags, "exclude-tags", nil, "Exclude checks matching ANY of these tags (comma-separated)")
}

func runList(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	// Apply tag filtering
	checksToShow := filterChecksForList(cfg.Checks)

	out := cmd.OutOrStdout()
	_, _ = fmt.Fprintf(out, "Checks (%d):\n\n", len(checksToShow))

	for _, check := range checksToShow {
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

// filterChecksForList applies tag-based filtering to checks for list display.
func filterChecksForList(checks []config.Check) []config.Check {
	if len(tags) == 0 && len(excludeTags) == 0 {
		// No filter, return all checks
		return checks
	}

	filtered := []config.Check{}

	for _, check := range checks {
		// Check inclusion: if tags are specified, check must match at least one tag
		if len(tags) > 0 {
			matchesInclude := false
			for _, includeTag := range tags {
				for _, checkTag := range check.Tags {
					if includeTag == checkTag {
						matchesInclude = true
						break
					}
				}
				if matchesInclude {
					break
				}
			}
			if !matchesInclude {
				continue
			}
		}

		// Check exclusion: skip if matches any exclude tag
		if len(excludeTags) > 0 {
			matchesExclude := false
			for _, excludeTag := range excludeTags {
				for _, checkTag := range check.Tags {
					if excludeTag == checkTag {
						matchesExclude = true
						break
					}
				}
				if matchesExclude {
					break
				}
			}
			if matchesExclude {
				continue
			}
		}

		// Check passed both inclusion and exclusion filters
		filtered = append(filtered, check)
	}

	return filtered
}
