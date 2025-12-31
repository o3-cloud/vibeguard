package executor

import (
	"fmt"
	"os"
)

// ReadFile reads file contents to be used instead of command stdout.
// This supports the `file:` field in check configuration.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", path, err)
	}
	return string(data), nil
}
