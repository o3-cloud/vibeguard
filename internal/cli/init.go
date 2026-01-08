package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/cli/inspector"
	"github.com/vibeguard/vibeguard/internal/cli/templates"
	"github.com/vibeguard/vibeguard/internal/config"
)

var (
	initForce         bool
	initTemplate      string
	initAssist        bool
	initOutput        string
	initListTemplates bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter vibeguard.yaml",
	Long: `Create a starter configuration file in the current directory.

Use --template to select a predefined template:
  vibeguard init --template go-standard   Use the Go standard template
  vibeguard init --list-templates         List available templates

Use --assist for AI agent-assisted setup:
  vibeguard init --assist                 Generate a setup prompt for AI agents
  vibeguard init --assist --output p.md   Save the prompt to a file

Available templates: ` + strings.Join(templates.Names(), ", ") + `

Without --template, creates a default Go project configuration.
Use --force to overwrite an existing configuration file.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing config file")
	initCmd.Flags().StringVarP(&initTemplate, "template", "t", "", "Use a predefined template (use --list-templates to see available templates)")
	initCmd.Flags().BoolVar(&initListTemplates, "list-templates", false, "List available templates")
	initCmd.Flags().BoolVar(&initAssist, "assist", false, "Generate an AI agent-assisted setup prompt")
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "", "Output file for --assist mode (default: stdout)")
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
	// Handle --assist mode
	if initAssist {
		return runAssist(cmd, args)
	}

	// Handle --list-templates flag
	if initListTemplates {
		return listTemplates()
	}

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
			return fmt.Errorf("unknown template %q (use --list-templates to see available templates)", initTemplate)
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
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
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

// runAssist generates an AI agent-assisted setup prompt based on project analysis.
func runAssist(cmd *cobra.Command, args []string) error {
	// Get the project root directory
	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	// Verify the directory exists
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return &ExitError{
				Code:    2,
				Message: fmt.Sprintf("directory does not exist: %s", root),
			}
		}
		return &ExitError{
			Code:    3,
			Message: fmt.Sprintf("failed to access directory: %v", err),
		}
	}
	if !info.IsDir() {
		return &ExitError{
			Code:    2,
			Message: fmt.Sprintf("not a directory: %s", root),
		}
	}

	// Project detection is delegated to the AI agent.
	// Generate a minimal setup prompt without project type detection.
	prompt, err := inspector.GenerateSetupPromptWithoutDetection(root)
	if err != nil {
		return &ExitError{
			Code:    3,
			Message: fmt.Sprintf("failed to generate setup prompt: %v", err),
		}
	}

	// Output the prompt
	if initOutput != "" {
		// Write to file
		if err := os.WriteFile(initOutput, []byte(prompt), 0600); err != nil {
			return &ExitError{
				Code:    3,
				Message: fmt.Sprintf("failed to write output file: %v", err),
			}
		}
		absPath, _ := filepath.Abs(initOutput)
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Setup prompt written to: %s\n", absPath)
	} else {
		// Write to stdout
		fmt.Print(prompt)
	}

	return nil
}
