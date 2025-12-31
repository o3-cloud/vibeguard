package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigError represents a configuration error that should result in exit code 2.
type ConfigError struct {
	Message string
	Cause   error
}

func (e *ConfigError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *ConfigError) Unwrap() error {
	return e.Cause
}

// IsConfigError returns true if the given error is a configuration error.
func IsConfigError(err error) bool {
	var configErr *ConfigError
	return errors.As(err, &configErr)
}

// Load reads and parses a VibeGuard configuration file.
// If path is empty, it searches for config files in the default locations.
// Returns a ConfigError for any configuration-related errors (exit code 2).
func Load(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = findConfigFile()
		if err != nil {
			return nil, &ConfigError{Message: "no config file found", Cause: err}
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ConfigError{Message: "failed to read config file", Cause: err}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, &ConfigError{Message: "failed to parse config file", Cause: err}
	}

	// Apply defaults
	cfg.applyDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, &ConfigError{Message: "configuration validation failed", Cause: err}
	}

	// Interpolate variables
	cfg.Interpolate()

	return &cfg, nil
}

// findConfigFile searches for a config file in the default locations.
func findConfigFile() (string, error) {
	for _, name := range ConfigFileNames {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
	}
	return "", fmt.Errorf("no config file found (tried: %v)", ConfigFileNames)
}

// applyDefaults sets default values for optional fields.
func (c *Config) applyDefaults() {
	if c.Version == "" {
		c.Version = "1"
	}
	if c.Vars == nil {
		c.Vars = make(map[string]string)
	}

	for i := range c.Checks {
		if c.Checks[i].Severity == "" {
			c.Checks[i].Severity = SeverityError
		}
		if c.Checks[i].Timeout == 0 {
			c.Checks[i].Timeout = Duration(DefaultTimeout)
		}
	}
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error {
	if c.Version != "1" {
		return fmt.Errorf("unsupported config version: %s", c.Version)
	}

	if len(c.Checks) == 0 {
		return fmt.Errorf("no checks defined")
	}

	checkIDs := make(map[string]bool)
	for i, check := range c.Checks {
		if check.ID == "" {
			return fmt.Errorf("check at index %d has no id", i)
		}
		if checkIDs[check.ID] {
			return fmt.Errorf("duplicate check id: %s", check.ID)
		}
		checkIDs[check.ID] = true

		if check.Run == "" {
			return fmt.Errorf("check %q has no run command", check.ID)
		}

		// Validate severity
		if check.Severity != SeverityError && check.Severity != SeverityWarning {
			return fmt.Errorf("check %q has invalid severity: %s", check.ID, check.Severity)
		}

		// Validate requires references
		for _, reqID := range check.Requires {
			// Check for self-reference
			if reqID == check.ID {
				return fmt.Errorf("check %q cannot require itself", check.ID)
			}

			if !checkIDs[reqID] {
				// Check if it exists at all (forward reference)
				found := false
				for _, c := range c.Checks {
					if c.ID == reqID {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("check %q requires unknown check: %s", check.ID, reqID)
				}
			}
		}
	}

	// Validate no cyclic dependencies
	if err := c.validateNoCycles(); err != nil {
		return err
	}

	return nil
}

// validateNoCycles checks for cyclic dependencies in the requires graph.
// It uses DFS with three states: unvisited, visiting (in current path), and visited (fully processed).
func (c *Config) validateNoCycles() error {
	// Build adjacency list: check ID -> list of required check IDs
	graph := make(map[string][]string)
	for _, check := range c.Checks {
		graph[check.ID] = check.Requires
	}

	// Track visited state: 0 = unvisited, 1 = visiting (in current path), 2 = visited
	state := make(map[string]int)

	// Track the current path for error reporting
	var path []string

	// DFS function to detect cycles
	var dfs func(id string) error
	dfs = func(id string) error {
		if state[id] == 2 {
			// Already fully visited, no cycle through this node
			return nil
		}
		if state[id] == 1 {
			// Found a cycle - build cycle description
			cycleStart := -1
			for i, p := range path {
				if p == id {
					cycleStart = i
					break
				}
			}
			cyclePath := append(path[cycleStart:], id)
			return fmt.Errorf("cyclic dependency detected: %s", formatCycle(cyclePath))
		}

		// Mark as visiting
		state[id] = 1
		path = append(path, id)

		// Visit all dependencies
		for _, reqID := range graph[id] {
			if err := dfs(reqID); err != nil {
				return err
			}
		}

		// Mark as fully visited
		state[id] = 2
		path = path[:len(path)-1]
		return nil
	}

	// Run DFS from each unvisited node
	for _, check := range c.Checks {
		if state[check.ID] == 0 {
			if err := dfs(check.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// formatCycle formats a cycle path for display (e.g., "a -> b -> c -> a").
func formatCycle(path []string) string {
	result := path[0]
	for i := 1; i < len(path); i++ {
		result += " -> " + path[i]
	}
	return result
}

// UnmarshalYAML implements custom YAML unmarshaling for GrokSpec.
func (g *GrokSpec) UnmarshalYAML(value *yaml.Node) error {
	// Try single string first
	var single string
	if err := value.Decode(&single); err == nil {
		*g = GrokSpec{single}
		return nil
	}

	// Try list of strings
	var list []string
	if err := value.Decode(&list); err != nil {
		return err
	}
	*g = list
	return nil
}

// UnmarshalYAML implements custom YAML unmarshaling for Duration.
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	*d = Duration(duration)
	return nil
}

// AsDuration returns the Duration as a time.Duration.
func (d Duration) AsDuration() time.Duration {
	return time.Duration(d)
}
