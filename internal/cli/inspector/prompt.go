// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// GenerateSetupPrompt creates a Claude Code-friendly setup prompt based on inspection results.
func GenerateSetupPrompt(
	projectType *DetectionResult,
	tools []ToolInfo,
	metadata *ProjectMetadata,
	structure *ProjectStructure,
	recommendations []CheckRecommendation,
) (string, error) {
	const promptTemplate = `# VibeGuard AI Agent Setup Guide

You are being asked to help set up VibeGuard policy enforcement for a software project.
VibeGuard is a declarative policy tool that runs quality checks and assertions on code.

This guide will help you understand the project structure, existing tools, and how to
generate a valid configuration.

---

## Project Analysis

**Project Name:** {{.Metadata.Name}}
**Project Type:** {{.ProjectType.Type}}
**Detection Confidence:** {{printf "%.0f" (mul .ProjectType.Confidence 100)}}%
**Language Version:** {{.Metadata.Version}}

### Main Tools Detected:
{{range .Tools}}{{if .Detected}}- {{.Name}}{{if .ConfigFile}} (config: {{.ConfigFile}}){{end}}
{{end}}{{end}}

### Project Structure:
- Source directories: {{join .Structure.SourceDirs ", "}}
- Test directories: {{join .Structure.TestDirs ", "}}
- Entry points: {{join .Structure.EntryPoints ", "}}

### Build System:
- Build output: {{if .Structure.BuildOutputDir}}{{.Structure.BuildOutputDir}}{{else}}default{{end}}

---

## Recommended Checks

Based on the detected tools, here are the recommended checks:

{{range .Recommendations}}
### {{.ID}} ({{.Category}})
**Description:** {{.Description}}
**Rationale:** {{.Rationale}}
**Command:** ` + "`{{.Command}}`" + `
**Severity:** {{.Severity}}
{{if .Grok}}**Grok Patterns:** {{range .Grok}}` + "`{{.}}`" + ` {{end}}{{end}}
{{if .Assert}}**Assertion:** ` + "`{{.Assert}}`" + `{{end}}
{{if .Suggestion}}**Suggestion on failure:** {{.Suggestion}}{{end}}
{{if .Requires}}**Requires:** {{join .Requires ", "}}{{end}}

{{end}}

---

## Configuration Requirements

A valid vibeguard.yaml must contain:

### 1. Version (required)
` + "```yaml" + `
version: "1"
` + "```" + `

### 2. Variables (optional)
Global variables for interpolation in commands and assertions.

` + "```yaml" + `
vars:
  go_packages: "./..."
  test_dir: "./..."
` + "```" + `

### 3. Checks (required)
Array of check definitions. Each check must have:
- **id:** Unique identifier (alphanumeric + underscore + hyphen)
- **run:** Shell command to execute

Optional fields:
- **grok:** Array of patterns to extract data from output (uses Grok syntax)
- **assert:** Condition that must be true (e.g., "coverage >= 70")
- **requires:** Array of check IDs that must pass first
- **severity:** "error" or "warning" (default: error)
- **suggestion:** Message shown on failure (supports {{"{{"}}` + "`.variable`" + `{{"}}"}} templating)
- **timeout:** Duration string (e.g., "30s", "5m")

---

## Go-Specific Examples

### Format Check
` + "```yaml" + `
- id: fmt
  run: test -z "$(gofmt -l .)"
  severity: error
  suggestion: "Run 'gofmt -w .' to format code"
  timeout: 5s
` + "```" + `

### Lint Check
` + "```yaml" + `
- id: lint
  run: golangci-lint run {{"{{"}}` + "`.go_packages`" + `{{"}}"}}
  severity: error
  suggestion: "Fix linting issues. Run 'golangci-lint run --fix' for auto-fixes."
  timeout: 30s
` + "```" + `

### Test with Coverage
` + "```yaml" + `
- id: test
  run: go test {{"{{"}}` + "`.go_packages`" + `{{"}}"}}
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 60s

- id: coverage
  run: go test ./... -coverprofile cover.out && go tool cover -func cover.out
  grok:
    - total:.*\(statements\)\s+%{NUMBER:coverage}%
  assert: "coverage >= 70"
  severity: warning
  suggestion: "Coverage is {{"{{"}}` + "`.coverage`" + `{{"}}"}}%, target is 70%. Add tests to improve."
  requires:
    - test
  timeout: 60s
` + "```" + `

### Build Check with Dependency
` + "```yaml" + `
- id: build
  run: go build {{"{{"}}` + "`.go_packages`" + `{{"}}"}}
  severity: error
  suggestion: "Fix compilation errors"
  timeout: 30s
  requires:
    - vet
` + "```" + `

---

## Validation Rules

Your generated configuration must:

1. Be valid YAML syntax
2. Have ` + "`version: \"1\"`" + ` at the top level
3. Include at least one check in the ` + "`checks`" + ` array
4. Each check must have a unique ` + "`id`" + `
5. Each check must have a non-empty ` + "`run`" + ` command
6. All ` + "`requires`" + ` references must point to existing check IDs
7. No circular dependencies in ` + "`requires`" + `
8. All variables used in double-curly-brace syntax (e.g., .var) must be defined in vars
9. ` + "`severity`" + ` must be "error" or "warning"
10. Grok patterns must be valid Grok syntax

**DO NOT:**
- Include YAML comments in the generated config
- Add extra top-level keys beyond version, vars, checks
- Use undefined variables
- Create checks for tools not detected in the project

---

## Your Task

Based on the project analysis above, generate a vibeguard.yaml configuration that:

1. Includes version: "1"
2. Defines appropriate variables for this project
3. Creates checks for the detected tools
4. Follows the syntax rules described above
5. Includes helpful suggestions for each check
6. Uses appropriate timeouts for each check type

Output the configuration in a YAML code block:

` + "```yaml" + `
# vibeguard.yaml for {{.Metadata.Name}}
version: "1"

vars:
  # ... your variables ...

checks:
  # ... your checks ...
` + "```" + `

After generating the configuration, verify it would pass the validation rules listed above.
`

	// Template functions
	funcMap := template.FuncMap{
		"join": strings.Join,
		"mul": func(a float64, b float64) float64 {
			return a * b
		},
	}

	tmpl, err := template.New("prompt").Funcs(funcMap).Parse(promptTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	data := struct {
		ProjectType     *DetectionResult
		Tools           []ToolInfo
		Metadata        *ProjectMetadata
		Structure       *ProjectStructure
		Recommendations []CheckRecommendation
	}{
		ProjectType:     projectType,
		Tools:           tools,
		Metadata:        metadata,
		Structure:       structure,
		Recommendations: recommendations,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return buf.String(), nil
}
