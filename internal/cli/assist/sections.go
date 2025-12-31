// Package assist provides AI agent-assisted setup functionality for VibeGuard.
package assist

import (
	"fmt"
	"strings"
)

// PromptSection represents a discrete section of the AI agent setup prompt.
type PromptSection struct {
	Title   string
	Content string
}

// HeaderSection generates the introductory section of the prompt.
func HeaderSection() PromptSection {
	return PromptSection{
		Title: "Header",
		Content: `# VibeGuard AI Agent Setup Guide

You are being asked to help set up VibeGuard policy enforcement for a software project.
VibeGuard is a declarative policy tool that runs quality checks and assertions on code.

This guide will help you understand the project structure, existing tools, and how to
generate a valid configuration.`,
	}
}

// ProjectAnalysisSection generates the project analysis section.
func ProjectAnalysisSection(analysis *ProjectAnalysis) PromptSection {
	var sb strings.Builder

	sb.WriteString("## Project Analysis\n\n")
	sb.WriteString(fmt.Sprintf("**Project Name:** %s\n", analysis.Name))
	sb.WriteString(fmt.Sprintf("**Project Type:** %s\n", analysis.ProjectType))
	sb.WriteString(fmt.Sprintf("**Detection Confidence:** %.0f%%\n", analysis.Confidence*100))
	if analysis.LanguageVersion != "" {
		sb.WriteString(fmt.Sprintf("**Language Version:** %s\n", analysis.LanguageVersion))
	}

	// Tools detected
	if len(analysis.DetectedTools) > 0 {
		sb.WriteString("\n### Main Tools Detected:\n")
		for _, tool := range analysis.DetectedTools {
			if tool.ConfigFile != "" {
				sb.WriteString(fmt.Sprintf("- %s (config: %s)\n", tool.Name, tool.ConfigFile))
			} else {
				sb.WriteString(fmt.Sprintf("- %s\n", tool.Name))
			}
		}
	}

	// Project structure
	sb.WriteString("\n### Project Structure:\n")
	if len(analysis.SourceDirs) > 0 {
		sb.WriteString(fmt.Sprintf("- Source directories: %s\n", strings.Join(analysis.SourceDirs, ", ")))
	}
	if len(analysis.TestDirs) > 0 {
		sb.WriteString(fmt.Sprintf("- Test directories: %s\n", strings.Join(analysis.TestDirs, ", ")))
	}
	if len(analysis.EntryPoints) > 0 {
		sb.WriteString(fmt.Sprintf("- Entry points: %s\n", strings.Join(analysis.EntryPoints, ", ")))
	}

	// Build system
	sb.WriteString("\n### Build System:\n")
	if analysis.BuildOutputDir != "" {
		sb.WriteString(fmt.Sprintf("- Build output: %s\n", analysis.BuildOutputDir))
	} else {
		sb.WriteString("- Build output: default\n")
	}

	return PromptSection{
		Title:   "Project Analysis",
		Content: sb.String(),
	}
}

// RecommendationsSection generates the recommended checks section.
func RecommendationsSection(recommendations []CheckRecommendation) PromptSection {
	var sb strings.Builder

	sb.WriteString("## Recommended Checks\n\n")
	sb.WriteString("Based on the detected tools, here are the recommended checks:\n\n")

	for _, rec := range recommendations {
		sb.WriteString(fmt.Sprintf("### %s (%s)\n", rec.ID, rec.Category))
		sb.WriteString(fmt.Sprintf("**Description:** %s\n", rec.Description))
		sb.WriteString(fmt.Sprintf("**Rationale:** %s\n", rec.Rationale))
		sb.WriteString(fmt.Sprintf("**Command:** `%s`\n", rec.Command))
		sb.WriteString(fmt.Sprintf("**Severity:** %s\n", rec.Severity))
		if len(rec.Grok) > 0 {
			sb.WriteString(fmt.Sprintf("**Grok Patterns:** %s\n", formatGrokPatterns(rec.Grok)))
		}
		if rec.Assert != "" {
			sb.WriteString(fmt.Sprintf("**Assertion:** `%s`\n", rec.Assert))
		}
		if rec.Suggestion != "" {
			sb.WriteString(fmt.Sprintf("**Suggestion on failure:** %s\n", rec.Suggestion))
		}
		if len(rec.Requires) > 0 {
			sb.WriteString(fmt.Sprintf("**Requires:** %s\n", strings.Join(rec.Requires, ", ")))
		}
		sb.WriteString("\n")
	}

	return PromptSection{
		Title:   "Recommended Checks",
		Content: sb.String(),
	}
}

// ConfigRequirementsSection generates the configuration requirements section.
func ConfigRequirementsSection() PromptSection {
	return PromptSection{
		Title: "Configuration Requirements",
		Content: `## Configuration Requirements

A valid vibeguard.yaml must contain:

### 1. Version (required)
` + "```yaml" + `
version: "1"
` + "```" + `

### 2. Variables (optional)
Global variables for interpolation in commands and assertions.

` + "```yaml" + `
vars:
  packages: "./..."
  coverage_threshold: "70"
` + "```" + `

### 3. Checks (required)
Array of check definitions. Each check must have:
- **id:** Unique identifier (alphanumeric + underscore + hyphen)
- **run:** Shell command to execute

Optional fields:
- **grok:** Array of patterns to extract data from output
- **assert:** Condition that must be true
- **requires:** Array of check IDs that must pass first
- **severity:** "error" or "warning" (default: error)
- **suggestion:** Message shown on failure
- **timeout:** Duration string (e.g., "30s", "5m")
- **file:** Path to read output from instead of command stdout`,
	}
}

// LanguageExamplesSection generates language-specific examples.
func LanguageExamplesSection(projectType string) PromptSection {
	examples, ok := languageExamples[projectType]
	if !ok {
		examples = languageExamples["generic"]
	}

	return PromptSection{
		Title:   fmt.Sprintf("%s-Specific Examples", capitalizeFirst(projectType)),
		Content: examples,
	}
}

// ValidationRulesSection generates the validation rules section.
func ValidationRulesSection() PromptSection {
	guide := NewValidationGuide()
	return PromptSection{
		Title:   "Validation Rules",
		Content: guide.GetFullGuide(),
	}
}

// TaskSection generates the final task instruction section.
func TaskSection(projectName string) PromptSection {
	return PromptSection{
		Title: "Your Task",
		Content: fmt.Sprintf(`## Your Task

Based on the project analysis above, generate a vibeguard.yaml configuration that:

1. Includes version: "1"
2. Defines appropriate variables for this project
3. Creates checks for the detected tools
4. Follows the syntax rules described above
5. Includes helpful suggestions for each check
6. Uses appropriate timeouts for each check type

Output the configuration in a YAML code block:

`+"```yaml"+`
# vibeguard.yaml for %s
version: "1"

vars:
  # ... your variables ...

checks:
  # ... your checks ...
`+"```"+`

After generating the configuration, verify it would pass the validation rules listed above.`, projectName),
	}
}

// formatGrokPatterns formats grok patterns for display.
func formatGrokPatterns(patterns []string) string {
	var formatted []string
	for _, p := range patterns {
		formatted = append(formatted, fmt.Sprintf("`%s`", p))
	}
	return strings.Join(formatted, " ")
}

// capitalizeFirst capitalizes the first letter of a string.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
