// Package assist provides AI agent-assisted setup functionality for VibeGuard.
// It generates prompts and validation guides to help AI agents create valid configurations.
package assist

// ValidationGuide contains all validation rules for VibeGuard configurations.
// These templates are used to guide AI agents in generating valid vibeguard.yaml files.
type ValidationGuide struct {
	YAMLSyntax         string
	CheckStructure     string
	DependencyRules    string
	InterpolationRules string
	DoNotList          string
}

// NewValidationGuide creates a new ValidationGuide with all rule sections.
func NewValidationGuide() *ValidationGuide {
	return &ValidationGuide{
		YAMLSyntax:         YAMLSyntaxRules,
		CheckStructure:     CheckStructureRules,
		DependencyRules:    DependencyValidationRules,
		InterpolationRules: VariableInterpolationRules,
		DoNotList:          ExplicitDoNotList,
	}
}

// GetFullGuide returns the complete validation guide as a single string.
func (g *ValidationGuide) GetFullGuide() string {
	return g.YAMLSyntax + "\n\n" +
		g.CheckStructure + "\n\n" +
		g.DependencyRules + "\n\n" +
		g.InterpolationRules + "\n\n" +
		g.DoNotList
}

// YAMLSyntaxRules defines YAML syntax requirements for VibeGuard configurations.
const YAMLSyntaxRules = `## YAML Syntax Requirements

The vibeguard.yaml file must be valid YAML. Follow these rules:

### Basic Structure
- The file must start with valid YAML syntax
- Use 2-space indentation consistently
- Strings with special characters must be quoted
- Arrays can use either block style (- item) or flow style ([item1, item2])

### Required Top-Level Keys
The configuration must have these top-level keys:
- **version** (required): Must be "1" (as a string)
- **vars** (optional): Map of variable names to string values
- **checks** (required): Array of check definitions

### String Quoting Rules
Quote strings that contain:
- Colons followed by space (: )
- Special YAML characters: {, }, [, ], ,, &, *, #, ?, |, -, <, >, =, !, %, @, \
- Leading/trailing spaces
- Numbers that should be treated as strings
- Boolean-like values: yes, no, true, false, on, off

### Example Valid Structure
` + "```yaml" + `
version: "1"

vars:
  packages: "./..."
  coverage_threshold: "70"

checks:
  - id: fmt
    run: gofmt -l .
    severity: error
` + "```" + `

### Common Syntax Errors to Avoid
1. Missing quotes around version number: version: 1 (wrong) vs version: "1" (correct)
2. Tabs instead of spaces for indentation
3. Inconsistent indentation levels
4. Missing space after colon in key-value pairs
5. Unquoted strings with special characters`

// CheckStructureRules defines the structure requirements for individual checks.
const CheckStructureRules = `## Check Structure Requirements

Each check in the **checks** array must follow this structure:

### Required Fields
- **id** (string, required): Unique identifier for the check
  - Must start with a letter or underscore
  - Can contain letters, numbers, underscores, and hyphens
  - Must be unique across all checks
  - Examples: "fmt", "lint", "go-test", "npm_audit", "_private"

- **run** (string, required): Shell command to execute
  - Must be non-empty
  - Will be executed via shell (sh -c)
  - Can reference variables using {{.varname}} syntax

### Optional Fields
- **grok** (string or array of strings): Patterns to extract data from command output
  - Uses Grok syntax (similar to Logstash)
  - Common patterns: %{NUMBER:varname}, %{WORD:varname}, %{GREEDYDATA:varname}
  - Extracted values can be used in assertions

- **assert** (string): Condition that must be true for the check to pass
  - Supports comparisons: ==, !=, <, <=, >, >=
  - Supports boolean operators: &&, ||, !
  - References extracted grok values: coverage >= 70
  - References special variables: exit_code == 0, stdout == ""

- **severity** (string): "error" or "warning"
  - Default: "error"
  - "error": Check failure fails the overall run
  - "warning": Check failure is reported but doesn't fail the run

- **suggestion** (string): Message shown when check fails
  - Supports variable interpolation: {{.varname}}
  - Can reference grok-extracted values

- **requires** (array of strings): IDs of checks that must pass first
  - Creates dependency ordering
  - Referenced checks must exist
  - No circular dependencies allowed

- **timeout** (string): Maximum time for check execution
  - Format: Go duration string (e.g., "30s", "5m", "1h")
  - Default: "30s"

- **file** (string): File to read output from instead of command stdout
  - Useful for reading generated reports

### Example Complete Check
` + "```yaml" + `
- id: coverage
  run: go test -cover ./... 2>&1 | tail -1
  grok:
    - "coverage: %{NUMBER:coverage}%"
  assert: "coverage >= 70"
  severity: warning
  suggestion: "Coverage is {{.coverage}}%, target is 70%. Add more tests."
  requires:
    - test
  timeout: 60s
` + "```" + ``

// DependencyValidationRules defines rules for check dependencies.
const DependencyValidationRules = `## Dependency Validation Rules

The **requires** field creates dependencies between checks. These rules must be followed:

### Basic Rules
1. All IDs in **requires** must reference existing checks
2. A check cannot require itself (no self-reference)
3. No circular dependencies allowed

### Circular Dependency Detection
A circular dependency exists when check A requires B, and B (directly or indirectly) requires A.

Examples of INVALID circular dependencies:
` + "```yaml" + `
# Direct cycle: A -> B -> A
checks:
  - id: A
    run: echo A
    requires: [B]
  - id: B
    run: echo B
    requires: [A]

# Indirect cycle: A -> B -> C -> A
checks:
  - id: A
    run: echo A
    requires: [C]
  - id: B
    run: echo B
    requires: [A]
  - id: C
    run: echo C
    requires: [B]
` + "```" + `

### Valid Dependency Patterns
` + "```yaml" + `
# Linear chain: fmt -> vet -> lint -> test
checks:
  - id: fmt
    run: gofmt -l .
  - id: vet
    run: go vet ./...
    requires: [fmt]
  - id: lint
    run: golangci-lint run
    requires: [vet]
  - id: test
    run: go test ./...
    requires: [lint]

# Diamond pattern: A <- B, A <- C, B <- D, C <- D
checks:
  - id: A
    run: echo A
  - id: B
    run: echo B
    requires: [A]
  - id: C
    run: echo C
    requires: [A]
  - id: D
    run: echo D
    requires: [B, C]
` + "```" + `

### Execution Order
- Checks with no dependencies run first (potentially in parallel)
- A check only runs after all its required checks pass
- If a required check fails, dependent checks are skipped`

// VariableInterpolationRules defines how variable interpolation works.
const VariableInterpolationRules = `## Variable Interpolation Rules

Variables defined in **vars** can be referenced throughout the configuration.

### Syntax
- Variables are referenced using Go template syntax: {{.varname}}
- The variable name must match exactly (case-sensitive)
- Variables must be defined in the **vars** section before use

### Where Variables Can Be Used
Variables can be interpolated in these fields:
- **run**: Command to execute
- **assert**: Assertion expression
- **suggestion**: Failure message
- **file**: Output file path
- **grok**: Pattern strings (array elements)

### Example Usage
` + "```yaml" + `
version: "1"

vars:
  packages: "./cmd/... ./internal/... ./pkg/..."
  coverage_min: "70"
  test_timeout: "5m"

checks:
  - id: test
    run: go test {{.packages}} -timeout {{.test_timeout}}
    severity: error
    suggestion: "Tests failed in {{.packages}}"

  - id: coverage
    run: go test {{.packages}} -cover 2>&1 | tail -1
    grok:
      - "coverage: %{NUMBER:coverage}%"
    assert: "coverage >= {{.coverage_min}}"
    suggestion: "Coverage {{.coverage}}% is below {{.coverage_min}}%"
` + "```" + `

### Variable Naming Rules
- Use alphanumeric characters and underscores
- Start with a letter or underscore
- Case-sensitive: {{.Packages}} != {{.packages}}
- Good names: packages, coverage_min, test_timeout
- Avoid: kebab-case ({{.test-timeout}}) or spaces

### Grok-Extracted Values
Values extracted by grok patterns are available in:
- **assert**: Reference extracted values directly (coverage >= 70)
- **suggestion**: Use template syntax ({{.coverage}})

Config vars take precedence over grok-extracted values if names conflict.

### Common Mistakes
1. Using undefined variables: {{.undefined_var}} will remain as literal text
2. Wrong syntax: {.var} or {{ .var }} (must be {{.var}})
3. Referencing grok values in **run** (grok runs after the command)`

// ExplicitDoNotList defines explicit constraints for AI agents.
const ExplicitDoNotList = `## Explicit DO NOT List

When generating a vibeguard.yaml configuration, DO NOT:

### YAML Structure
- DO NOT include YAML comments (# comment) in the generated config
- DO NOT add extra top-level keys beyond version, vars, checks
- DO NOT use YAML anchors and aliases (&anchor, *alias)
- DO NOT use multi-document YAML (--- separator)

### Check Definitions
- DO NOT create checks for tools not detected in the project
- DO NOT use empty strings for required fields (id, run)
- DO NOT use invalid check IDs (must match ^[a-zA-Z_][a-zA-Z0-9_-]*$ - start with letter or underscore, then alphanumeric/underscore/hyphen)
- DO NOT use duplicate check IDs
- DO NOT reference undefined variables
- DO NOT create circular dependencies in requires

### Assertions and Grok
- DO NOT use invalid assertion syntax (missing operators, invalid comparisons)
- DO NOT reference grok-extracted values that don't exist
- DO NOT use unsupported assertion operators
- DO NOT write grok patterns that don't match Grok syntax

### Commands
- DO NOT assume tools are installed without detection evidence
- DO NOT use interactive commands (require user input)
- DO NOT use commands that modify system state destructively
- DO NOT hardcode paths that are project-specific without using variables

### Severity and Timeouts
- DO NOT use severity values other than "error" or "warning"
- DO NOT use invalid timeout formats (use Go duration: "30s", "5m", not "30 seconds")
- DO NOT set unreasonably short timeouts that would cause false failures

### Variable Interpolation
- DO NOT use variables before defining them in vars
- DO NOT mix up variable syntax (use {{.var}}, not $var or ${var})
- DO NOT name variables with special characters or spaces

### General Guidelines
- DO NOT generate configs that would fail vibeguard validate
- DO NOT assume the execution environment has specific tools without verification
- DO NOT create overly complex configurations when simple ones suffice
- DO NOT add checks that provide no value for the project type`

// AssertionOperators lists the supported assertion operators.
var AssertionOperators = []string{
	"==", // Equal
	"!=", // Not equal
	"<",  // Less than
	"<=", // Less than or equal
	">",  // Greater than
	">=", // Greater than or equal
	"&&", // Logical AND
	"||", // Logical OR
	"!",  // Logical NOT
}

// SpecialAssertionVariables lists special variables available in assertions.
var SpecialAssertionVariables = []string{
	"exit_code", // Exit code of the command
	"stdout",    // Standard output of the command
	"stderr",    // Standard error of the command
}

// SupportedSeverities lists the valid severity values.
var SupportedSeverities = []string{
	"error",
	"warning",
}

// GrokPatternExamples provides common grok pattern examples.
var GrokPatternExamples = map[string]string{
	"NUMBER":       "%{NUMBER:varname}",
	"WORD":         "%{WORD:varname}",
	"INT":          "%{INT:varname}",
	"GREEDYDATA":   "%{GREEDYDATA:varname}",
	"IP":           "%{IP:varname}",
	"TIMESTAMP":    "%{TIMESTAMP_ISO8601:varname}",
	"QUOTEDSTRING": "%{QUOTEDSTRING:varname}",
}

// CommonTimeoutValues provides recommended timeout values by check type.
var CommonTimeoutValues = map[string]string{
	"format":    "5s",
	"lint":      "30s",
	"test":      "60s",
	"build":     "60s",
	"coverage":  "60s",
	"security":  "120s",
	"typecheck": "60s",
}
