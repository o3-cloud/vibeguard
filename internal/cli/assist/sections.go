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

	sb.WriteString("## Your Project\n\n")
	sb.WriteString(fmt.Sprintf("**Project Name:** %s\n\n", analysis.Name))
	sb.WriteString("### Initial Analysis Instructions\n\n")
	sb.WriteString("Before proceeding, analyze the project structure to understand:\n\n")
	sb.WriteString("1. **Project Type & Language**: Examine the directory structure, file extensions, and package managers to identify the primary language and framework (Go, Node.js/TypeScript, Python, Rust, Java, Ruby, etc.)\n")
	sb.WriteString("2. **Detected Tools**: Look for configuration files like `go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`, `pom.xml`, `Gemfile` to identify build systems, package managers, and development tools\n")
	sb.WriteString("3. **Project Structure**: Identify source directories, test directories, entry points, and build output locations\n")
	sb.WriteString("4. **Existing Quality Tools**: Check for configuration files for linters, formatters, and test frameworks (eslint, pytest, golangci-lint, cargo, etc.)\n\n")

	return PromptSection{
		Title:   "Your Project",
		Content: sb.String(),
	}
}

// TemplateDiscoverySection generates instructions for discovering and using templates.
func TemplateDiscoverySection(projectType string) PromptSection {
	var sb strings.Builder

	sb.WriteString("## Available Templates\n\n")
	sb.WriteString("Rather than creating a custom configuration from scratch, you can use one of VibeGuard's predefined templates that are optimized for specific languages and frameworks.\n\n")

	sb.WriteString("### Discover Available Templates\n\n")
	sb.WriteString("Run this command to see all available templates:\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("vibeguard init --list-templates\n")
	sb.WriteString("```\n\n")
	sb.WriteString("This will display all templates with their descriptions and supported languages/frameworks.\n\n")

	sb.WriteString("### Selecting a Template\n\n")
	sb.WriteString("Based on your project analysis from the previous step:\n\n")
	sb.WriteString("1. **Match your language/framework**: Find the template that best matches the project type you identified\n")
	sb.WriteString("2. **Review the template**: Use `vibeguard init --list-templates` to see what each template covers\n")
	sb.WriteString("3. **Try the template**: Run `vibeguard init --template <template-name>` to use a predefined template\n\n")

	sb.WriteString("### When to Use Templates\n\n")
	sb.WriteString("Use a predefined template when:\n")
	sb.WriteString("- You want quick setup with reasonable defaults for your language/framework\n")
	sb.WriteString("- Your project matches the template's language/framework\n")
	sb.WriteString("- You're happy with the template's default checks\n\n")

	sb.WriteString("### When to Use Custom Configuration\n\n")
	sb.WriteString("Create a custom configuration when:\n")
	sb.WriteString("- No template matches your exact setup\n")
	sb.WriteString("- You have specific checks not in standard templates\n")
	sb.WriteString("- You want fine-grained control over all settings\n\n")

	sb.WriteString("You can also start with a template and modify it to suit your project's specific needs.")

	return PromptSection{
		Title:   "Available Templates",
		Content: sb.String(),
	}
}

// RecommendationsSection generates the recommended checks section.
// If recommendations is empty, it provides guidance on identifying checks from existing tools.
func RecommendationsSection(recommendations []CheckRecommendation) PromptSection {
	var sb strings.Builder

	if len(recommendations) == 0 {
		// No recommendations provided - guide the agent to identify checks from tools
		sb.WriteString("## Identifying Quality Checks\n\n")
		sb.WriteString("Based on your project analysis, you should identify the appropriate quality checks to include in the VibeGuard configuration.\n\n")
		sb.WriteString("### How to Identify Checks\n\n")
		sb.WriteString("For each quality tool or framework found in your analysis:\n\n")
		sb.WriteString("1. **Determine the command**: What command runs the tool? (e.g., `go test ./...`, `npm test`, `pytest`)\n")
		sb.WriteString("2. **Identify assertions**: What does success look like? (e.g., exit code 0, specific output patterns)\n")
		sb.WriteString("3. **Set severity**: Is this critical (error) or a warning?\n")
		sb.WriteString("4. **Add context**: Write a description explaining why this check is important\n\n")
		sb.WriteString("### Examples of Common Checks\n\n")
		sb.WriteString("- **Tests**: Run your test suite and ensure it passes\n")
		sb.WriteString("- **Linting**: Run code quality/style checkers\n")
		sb.WriteString("- **Type Checking**: Run type checkers (TypeScript, mypy, etc.)\n")
		sb.WriteString("- **Security**: Run security scanners (gosec, bandit, npm audit)\n")
		sb.WriteString("- **Coverage**: Verify test coverage meets threshold\n")
		sb.WriteString("- **Formatting**: Ensure code is properly formatted\n")
		return PromptSection{
			Title:   "Identifying Quality Checks",
			Content: sb.String(),
		}
	}

	// Recommendations provided - show them (backward compatibility)
	sb.WriteString("## Recommended Checks\n\n")
	sb.WriteString("Based on the detected tools, here are the recommended checks:\n\n")

	for _, rec := range recommendations {
		sb.WriteString(fmt.Sprintf("### %s (%s)\n", rec.ID, rec.Category))
		sb.WriteString(fmt.Sprintf("**Description:** %s\n", rec.Description))
		sb.WriteString(fmt.Sprintf("**Rationale:** %s\n", rec.Rationale))
		sb.WriteString(fmt.Sprintf("**Command:** `%s`\n", rec.Command))
		if rec.File != "" {
			sb.WriteString(fmt.Sprintf("**File:** `%s`\n", rec.File))
		}
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

Based on your project analysis above, follow these steps to set up VibeGuard for the %s project:

### Step 1: Choose Your Approach

You have two options:

**Option A: Use a Template (Recommended for most projects)**
- Run: `+"```bash"+`vibeguard init --list-templates`+"```"+`
- Find a template that matches your project type
- Run: `+"```bash"+`vibeguard init --template <template-name>`+"```"+`
- This creates a predefined vibeguard.yaml configuration optimized for your language/framework

**Option B: Create a Custom Configuration**
- If no template fits your needs, create a custom vibeguard.yaml
- Based on your analysis, identify the quality checks and tools to include
- Follow the Configuration Requirements section above
- Ensure each check has a clear purpose and appropriate settings

### Step 2: Validate Your Configuration

Once you have a vibeguard.yaml file, validate it:

`+"```bash"+`
vibeguard validate
`+"```"+`

This verifies the configuration has correct YAML syntax and adheres to the vibeguard schema.

### Step 3: Test Your Checks

Run the checks to ensure they execute properly:

`+"```bash"+`
vibeguard check
`+"```"+`

Review any failures and fix them:
- Adjust commands if tools aren't found
- Update assertions if they don't match your project's output
- Add variable definitions if needed
- Adjust timeouts if checks are timing out

### Step 4: Commit the Configuration

Once all checks pass, save the vibeguard.yaml file to your project repository. This ensures the configuration is version-controlled and used by all team members.

### Success Criteria

Your task is complete when:
- ✅ vibeguard validate passes without errors
- ✅ vibeguard check runs all checks successfully
- ✅ vibeguard.yaml is committed to the repository`, projectName),
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
