package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/cli/templates"
	"github.com/vibeguard/vibeguard/internal/config"
)

var (
	initForce    bool
	initTemplate string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter vibeguard.yaml",
	Long: `Create a starter configuration file in the current directory.

Use --template to select a predefined template:
  vibeguard init --template list          List available templates
  vibeguard init --template go-standard   Use the Go standard template

Available templates: ` + strings.Join(templates.Names(), ", ") + `

Without --template, creates a default Go project configuration.
Use --force to overwrite an existing configuration file.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing config file")
	initCmd.Flags().StringVarP(&initTemplate, "template", "t", "", "Use a predefined template (use 'list' to see available templates)")
	rootCmd.AddCommand(initCmd)
}

// starterConfig is the default starter configuration for Go projects
const starterConfig = `version: "1"

vars:
  go_packages: "./..."

checks:
  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    suggestion: "Run 'go vet ./...' and fix reported issues"
    timeout: 60s

  - id: fmt
    run: test -z "$(gofmt -l .)"
    severity: error
    suggestion: "Run 'gofmt -w .' to format code"
    timeout: 30s

  - id: test
    run: go test -race -cover {{.go_packages}}
    severity: error
    suggestion: "Run 'go test ./...' and fix failing tests"
    timeout: 300s
    requires:
      - vet
      - fmt

  - id: build
    run: go build {{.go_packages}}
    severity: error
    suggestion: "Run 'go build ./...' and fix compilation errors"
    timeout: 60s
    requires:
      - vet
`

func runInit(cmd *cobra.Command, args []string) error {
	// Handle --template list
	if initTemplate == "list" {
		return listTemplates()
	}

	// Determine which content to use (validate template early)
	var content string
	var templateName string

	if initTemplate != "" {
		// Use specified template
		tmpl, err := templates.Get(initTemplate)
		if err != nil {
			return fmt.Errorf("unknown template %q (use --template list to see available templates)", initTemplate)
		}
		content = tmpl.Content
		templateName = tmpl.Name
	} else {
		// Use default starter config
		content = starterConfig
		templateName = "default (Go)"
	}

	configPath := "vibeguard.yaml"

	// Check if any config file already exists
	if !initForce {
		for _, name := range config.ConfigFileNames {
			if _, err := os.Stat(name); err == nil {
				return fmt.Errorf("configuration file %q already exists (use --force to overwrite)", name)
			}
		}
	}

	// Write the config
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	absPath, _ := filepath.Abs(configPath)
	fmt.Printf("Created %s (template: %s)\n", absPath, templateName)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review and customize the checks in vibeguard.yaml")
	fmt.Println("  2. Run 'vibeguard check' to execute all checks")
	fmt.Println("  3. Run 'vibeguard validate' to verify your configuration")

	return nil
}

func listTemplates() error {
	tmplList := templates.List()
	fmt.Println("Available templates:")
	fmt.Println()
	for _, t := range tmplList {
		fmt.Printf("  %-20s %s\n", t.Name, t.Description)
	}
	fmt.Println()
	fmt.Println("Usage: vibeguard init --template <name>")
	return nil
}
