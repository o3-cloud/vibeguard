package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/vibeguard/vibeguard/internal/cli"
	"github.com/vibeguard/vibeguard/internal/config"
)

func main() {
	if err := cli.Execute(); err != nil {
		// Check if this is an ExitError (check completed but with specific exit code)
		var exitErr *cli.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.Code)
		}

		// Print error message
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		// Check for specific error types
		if config.IsConfigError(err) {
			os.Exit(2) // Configuration error
		}
		os.Exit(1)
	}
}
