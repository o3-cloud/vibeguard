package templates

func init() {
	Register(Template{
		Name:        "generic",
		Description: "Generic project template with placeholder checks to customize",
		Content: `version: "1"

# Customize these variables for your project
vars:
  source_dir: "src"

checks:
  # Add your formatting check
  # - id: format
  #   run: <your-formatter> --check {{.source_dir}}
  #   severity: error
  #   suggestion: "Run '<your-formatter>' to format your code"
  #   timeout: 60s

  # Add your linting check
  # - id: lint
  #   run: <your-linter> {{.source_dir}}
  #   severity: error
  #   suggestion: "Run '<your-linter> --fix' to fix issues"
  #   timeout: 60s

  # Add your test check
  - id: test
    run: echo "Configure your test command here"
    severity: error
    suggestion: "Configure the test command in vibeguard.yaml"
    timeout: 300s

  # Add your build check
  # - id: build
  #   run: <your-build-command>
  #   severity: error
  #   suggestion: "Fix build errors"
  #   timeout: 120s
`,
	})
}
