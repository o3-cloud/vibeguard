package main

import (
	"fmt"
	"os"

	"github.com/vibeguard/vibeguard/internal/cli"
	"github.com/vibeguard/vibeguard/internal/config"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if config.IsConfigError(err) {
			os.Exit(2) // Configuration error
		}
		os.Exit(1)
	}
}
