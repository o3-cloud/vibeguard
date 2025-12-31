package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/vibeguard/vibeguard/internal/config"
)

var (
	initForce bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter vibeguard.yaml",
	Long: `Create a starter configuration file in the current directory.

The generated file includes common checks for Go projects as a starting point.
Use --force to overwrite an existing configuration file.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing config file")
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
	configPath := "vibeguard.yaml"

	// Check if any config file already exists
	if !initForce {
		for _, name := range config.ConfigFileNames {
			if _, err := os.Stat(name); err == nil {
				return fmt.Errorf("configuration file %q already exists (use --force to overwrite)", name)
			}
		}
	}

	// Check if the file exists and --force is set
	if !initForce {
		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf("configuration file %q already exists (use --force to overwrite)", configPath)
		}
	}

	// Write the starter config
	if err := os.WriteFile(configPath, []byte(starterConfig), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	absPath, _ := filepath.Abs(configPath)
	fmt.Printf("Created %s\n", absPath)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review and customize the checks in vibeguard.yaml")
	fmt.Println("  2. Run 'vibeguard check' to execute all checks")
	fmt.Println("  3. Run 'vibeguard validate' to verify your configuration")

	return nil
}
