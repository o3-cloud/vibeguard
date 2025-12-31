// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/vibeguard/vibeguard/internal/cli/assist"
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

{{.ValidationGuide}}

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

	// Get comprehensive validation guide from the assist package
	validationGuide := assist.NewValidationGuide()

	data := struct {
		ProjectType     *DetectionResult
		Tools           []ToolInfo
		Metadata        *ProjectMetadata
		Structure       *ProjectStructure
		Recommendations []CheckRecommendation
		ValidationGuide string
	}{
		ProjectType:     projectType,
		Tools:           tools,
		Metadata:        metadata,
		Structure:       structure,
		Recommendations: recommendations,
		ValidationGuide: validationGuide.GetFullGuide(),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return buf.String(), nil
}
