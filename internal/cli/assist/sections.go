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

// ToolingInspectionSection generates instructions for AI agents to inspect existing tooling.
func ToolingInspectionSection(analysis *ProjectAnalysis) PromptSection {
	var sb strings.Builder

	sb.WriteString("## Tooling Inspection Instructions\n\n")
	sb.WriteString("Before generating the configuration, **inspect the existing tool configurations** in this project to understand how they're set up:\n\n")

	// List detected tools with config files that should be inspected
	configFilesToInspect := []string{}
	for _, tool := range analysis.DetectedTools {
		if tool.ConfigFile != "" {
			configFilesToInspect = append(configFilesToInspect, fmt.Sprintf("- **%s**: Read `%s` to understand its configuration", tool.Name, tool.ConfigFile))
		}
	}

	if len(configFilesToInspect) > 0 {
		sb.WriteString("### Configuration Files to Analyze\n\n")
		for _, item := range configFilesToInspect {
			sb.WriteString(item + "\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("### What to Look For\n\n")
	sb.WriteString("When inspecting each configuration file:\n\n")
	sb.WriteString("1. **Enabled rules/checks**: What specific rules or checks are enabled?\n")
	sb.WriteString("2. **Disabled rules**: Are any rules explicitly disabled? (They may conflict with project needs)\n")
	sb.WriteString("3. **Custom settings**: Are there custom paths, thresholds, or exclusions?\n")
	sb.WriteString("4. **Integration points**: Does the tool integrate with other tools (e.g., editor configs, CI)?\n")
	sb.WriteString("5. **Command variations**: Are there specific flags or options being used?\n\n")

	sb.WriteString("Use this understanding to craft VibeGuard checks that:\n")
	sb.WriteString("- Run tools with the same flags/options as configured\n")
	sb.WriteString("- Respect existing exclusions and thresholds\n")
	sb.WriteString("- Maintain consistency with the project's existing standards\n")

	return PromptSection{
		Title:   "Tooling Inspection Instructions",
		Content: sb.String(),
	}
}

// ToolingResearchSection generates instructions for AI agents to research additional tools.
func ToolingResearchSection(projectType string) PromptSection {
	var sb strings.Builder

	sb.WriteString("## Additional Tooling Research\n\n")
	sb.WriteString("Based on the project type, consider recommending additional quality and security tools that aren't currently configured.\n\n")

	// Get language-specific tool suggestions
	suggestions := getToolSuggestions(projectType)
	if len(suggestions) > 0 {
		sb.WriteString("### Suggested Tools to Consider\n\n")
		for _, suggestion := range suggestions {
			sb.WriteString(fmt.Sprintf("#### %s\n", suggestion.Name))
			sb.WriteString(fmt.Sprintf("**Category:** %s\n", suggestion.Category))
			sb.WriteString(fmt.Sprintf("**Purpose:** %s\n", suggestion.Purpose))
			sb.WriteString(fmt.Sprintf("**Value:** %s\n", suggestion.Value))
			if suggestion.Command != "" {
				sb.WriteString(fmt.Sprintf("**Example command:** `%s`\n", suggestion.Command))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("### How to Present Suggestions\n\n")
	sb.WriteString("When suggesting additional tools to the user:\n\n")
	sb.WriteString("1. **Explain the benefit**: Why would this tool help the project?\n")
	sb.WriteString("2. **Assess compatibility**: Will it work well with existing tools?\n")
	sb.WriteString("3. **Provide options**: Let the user decide which suggestions to include\n")
	sb.WriteString("4. **Include installation notes**: If a tool requires installation, mention it in the suggestion field\n\n")

	sb.WriteString("Ask the user: \"Would you like me to include checks for any of these additional tools?\"\n")

	return PromptSection{
		Title:   "Additional Tooling Research",
		Content: sb.String(),
	}
}

// ToolSuggestion represents a suggested tool for a project type.
type ToolSuggestion struct {
	Name     string
	Category string
	Purpose  string
	Value    string
	Command  string
}

// getToolSuggestions returns tool suggestions based on project type.
func getToolSuggestions(projectType string) []ToolSuggestion {
	switch projectType {
	case "go":
		return []ToolSuggestion{
			{
				Name:     "staticcheck",
				Category: "Linter",
				Purpose:  "Advanced static analysis for Go beyond go vet",
				Value:    "Catches bugs, suggests simplifications, and identifies deprecated code",
				Command:  "staticcheck ./...",
			},
			{
				Name:     "gosec",
				Category: "Security",
				Purpose:  "Security-focused static analysis for Go code",
				Value:    "Identifies potential security issues like SQL injection, hardcoded credentials",
				Command:  "gosec ./...",
			},
			{
				Name:     "errcheck",
				Category: "Linter",
				Purpose:  "Checks for unchecked errors in Go code",
				Value:    "Ensures error handling is not accidentally ignored",
				Command:  "errcheck ./...",
			},
			{
				Name:     "ineffassign",
				Category: "Linter",
				Purpose:  "Detects ineffectual assignments in Go code",
				Value:    "Finds variables assigned but never used",
				Command:  "ineffassign ./...",
			},
			{
				Name:     "govulncheck",
				Category: "Security",
				Purpose:  "Checks Go dependencies for known vulnerabilities",
				Value:    "Official Go team tool for CVE scanning in dependencies",
				Command:  "govulncheck ./...",
			},
		}
	case "node":
		return []ToolSuggestion{
			{
				Name:     "tsc --noEmit",
				Category: "Type Checker",
				Purpose:  "TypeScript type checking without emitting files",
				Value:    "Catches type errors before runtime",
				Command:  "npx tsc --noEmit",
			},
			{
				Name:     "npm audit",
				Category: "Security",
				Purpose:  "Scans dependencies for known vulnerabilities",
				Value:    "Identifies and helps fix security issues in dependencies",
				Command:  "npm audit --audit-level=moderate",
			},
			{
				Name:     "depcheck",
				Category: "Dependencies",
				Purpose:  "Finds unused and missing dependencies",
				Value:    "Keeps package.json clean and accurate",
				Command:  "npx depcheck",
			},
			{
				Name:     "madge",
				Category: "Architecture",
				Purpose:  "Detects circular dependencies",
				Value:    "Prevents hard-to-debug circular import issues",
				Command:  "npx madge --circular src/",
			},
		}
	case "python":
		return []ToolSuggestion{
			{
				Name:     "bandit",
				Category: "Security",
				Purpose:  "Security-focused static analysis for Python",
				Value:    "Finds common security issues like hardcoded passwords, SQL injection",
				Command:  "bandit -r .",
			},
			{
				Name:     "safety",
				Category: "Security",
				Purpose:  "Checks Python dependencies for known vulnerabilities",
				Value:    "Identifies CVEs in installed packages",
				Command:  "safety check",
			},
			{
				Name:     "vulture",
				Category: "Code Quality",
				Purpose:  "Finds dead/unused code in Python projects",
				Value:    "Helps remove unused functions, classes, and variables",
				Command:  "vulture .",
			},
			{
				Name:     "radon",
				Category: "Code Quality",
				Purpose:  "Computes code complexity metrics",
				Value:    "Identifies overly complex functions that need refactoring",
				Command:  "radon cc . -a",
			},
		}
	case "rust":
		return []ToolSuggestion{
			{
				Name:     "cargo audit",
				Category: "Security",
				Purpose:  "Audits Cargo.lock for crates with security vulnerabilities",
				Value:    "Identifies known CVEs in Rust dependencies",
				Command:  "cargo audit",
			},
			{
				Name:     "cargo deny",
				Category: "Dependencies",
				Purpose:  "Checks dependencies for licenses, bans, and advisories",
				Value:    "Enforces dependency policies and security requirements",
				Command:  "cargo deny check",
			},
			{
				Name:     "cargo outdated",
				Category: "Dependencies",
				Purpose:  "Shows outdated dependencies",
				Value:    "Helps keep dependencies up to date",
				Command:  "cargo outdated",
			},
		}
	case "ruby":
		return []ToolSuggestion{
			{
				Name:     "brakeman",
				Category: "Security",
				Purpose:  "Static analysis security scanner for Ruby on Rails",
				Value:    "Finds security vulnerabilities in Rails applications",
				Command:  "brakeman",
			},
			{
				Name:     "bundler-audit",
				Category: "Security",
				Purpose:  "Patch-level verification for bundled gems",
				Value:    "Checks for vulnerable gem versions",
				Command:  "bundle audit check --update",
			},
			{
				Name:     "reek",
				Category: "Code Quality",
				Purpose:  "Code smell detector for Ruby",
				Value:    "Identifies code smells that indicate design problems",
				Command:  "reek .",
			},
		}
	case "java":
		return []ToolSuggestion{
			{
				Name:     "SpotBugs",
				Category: "Linter",
				Purpose:  "Static analysis tool for finding bugs in Java",
				Value:    "Finds potential bugs using bug pattern analysis",
				Command:  "mvn spotbugs:check",
			},
			{
				Name:     "OWASP Dependency-Check",
				Category: "Security",
				Purpose:  "Identifies project dependencies with known vulnerabilities",
				Value:    "Scans for CVEs in Java dependencies",
				Command:  "mvn dependency-check:check",
			},
			{
				Name:     "PMD",
				Category: "Linter",
				Purpose:  "Source code analyzer for common programming flaws",
				Value:    "Finds unused variables, empty catch blocks, unnecessary object creation",
				Command:  "mvn pmd:check",
			},
		}
	default:
		return []ToolSuggestion{
			{
				Name:     "gitleaks",
				Category: "Security",
				Purpose:  "Detects hardcoded secrets in git repos",
				Value:    "Prevents accidental commit of API keys, passwords, tokens",
				Command:  "gitleaks detect",
			},
			{
				Name:     "trivy",
				Category: "Security",
				Purpose:  "Comprehensive security scanner",
				Value:    "Scans for vulnerabilities in code, dependencies, and containers",
				Command:  "trivy fs .",
			},
		}
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

After generating the configuration:

1. **Validate the YAML syntax and schema:**
   `+"```bash"+`
   vibeguard validate
   `+"```"+`
   This verifies the configuration file has correct YAML syntax and adheres to the vibeguard schema.

2. **Run the checks to verify they execute properly:**
   `+"```bash"+`
   vibeguard check
   `+"```"+`
   This runs all defined checks and ensures they execute successfully. Fix any failing checks before considering the task complete.

Only consider this task complete when both commands pass without errors.`, projectName),
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
