package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

// validCheckID matches alphanumeric characters, underscores, and hyphens.
// IDs must start with a letter or underscore, followed by any combination of
// alphanumeric characters, underscores, or hyphens.
var validCheckID = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`)

// ConfigError represents a configuration error that should result in exit code 2.
type ConfigError struct {
	Message  string
	Cause    error
	LineNum  int    // Line number in the config file (0 if not available)
	FileName string // File name for reference (optional)
}

func (e *ConfigError) Error() string {
	var msg string
	if e.Cause != nil {
		msg = fmt.Sprintf("%s: %v", e.Message, e.Cause)
	} else {
		msg = e.Message
	}
	if e.LineNum > 0 {
		msg = fmt.Sprintf("%s (line %d)", msg, e.LineNum)
	}
	return msg
}

func (e *ConfigError) Unwrap() error {
	return e.Cause
}

// IsConfigError returns true if the given error is a configuration error.
func IsConfigError(err error) bool {
	var configErr *ConfigError
	return errors.As(err, &configErr)
}

// ExecutionError represents an error that occurred during check execution with context.
// This is used for errors in grok patterns and assertions to provide file/line context.
type ExecutionError struct {
	Message   string
	Cause     error
	CheckID   string // ID of the check that failed
	LineNum   int    // Line number in the config file (0 if not available)
	ErrorType string // "grok" or "assert" to distinguish error type
}

func (e *ExecutionError) Error() string {
	var msg string
	if e.Cause != nil {
		msg = fmt.Sprintf("%s in check %q: %v", e.Message, e.CheckID, e.Cause)
	} else {
		msg = fmt.Sprintf("%s in check %q", e.Message, e.CheckID)
	}
	if e.LineNum > 0 {
		msg = fmt.Sprintf("%s (line %d)", msg, e.LineNum)
	}
	return msg
}

func (e *ExecutionError) Unwrap() error {
	return e.Cause
}

// IsExecutionError returns true if the given error is an execution error.
func IsExecutionError(err error) bool {
	var execErr *ExecutionError
	return errors.As(err, &execErr)
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

	data, err := os.ReadFile(path) // #nosec G304 - path is a validated config file from FindConfig
	if err != nil {
		return nil, &ConfigError{Message: "failed to read config file", Cause: err}
	}

	// Parse with nodes to preserve line information
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, &ConfigError{Message: "failed to parse config file", Cause: err}
	}

	var cfg Config
	if err := root.Decode(&cfg); err != nil {
		return nil, &ConfigError{Message: "failed to parse config file", Cause: err}
	}

	// Store the root node for line number lookups during validation
	cfg.yamlRoot = &root

	// Apply defaults
	cfg.applyDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
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
		return &ConfigError{Message: fmt.Sprintf("unsupported config version: %s", c.Version)}
	}

	if len(c.Checks) == 0 {
		return &ConfigError{Message: "no checks defined"}
	}

	checkIDs := make(map[string]bool)
	for i, check := range c.Checks {
		if check.ID == "" {
			return &ConfigError{
				Message: fmt.Sprintf("check at index %d has no id", i),
				LineNum: c.FindCheckNodeLine("", i),
			}
		}
		if !validCheckID.MatchString(check.ID) {
			return &ConfigError{
				Message: fmt.Sprintf("check %q has invalid id format: must start with a letter or underscore, followed by alphanumeric characters, underscores, or hyphens", check.ID),
				LineNum: c.FindCheckNodeLine(check.ID, i),
			}
		}
		if checkIDs[check.ID] {
			return &ConfigError{
				Message: fmt.Sprintf("duplicate check id: %s", check.ID),
				LineNum: c.FindCheckNodeLine(check.ID, i),
			}
		}
		checkIDs[check.ID] = true

		if check.Run == "" {
			return &ConfigError{
				Message: fmt.Sprintf("check %q has no run command", check.ID),
				LineNum: c.FindCheckNodeLine(check.ID, i),
			}
		}

		// Validate severity
		if check.Severity != SeverityError && check.Severity != SeverityWarning {
			return &ConfigError{
				Message: fmt.Sprintf("check %q has invalid severity: %s", check.ID, check.Severity),
				LineNum: c.FindCheckNodeLine(check.ID, i),
			}
		}

		// Validate requires references
		for _, reqID := range check.Requires {
			// Check for self-reference
			if reqID == check.ID {
				return &ConfigError{
					Message: fmt.Sprintf("check %q cannot require itself", check.ID),
					LineNum: c.FindCheckNodeLine(check.ID, i),
				}
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
					return &ConfigError{
						Message: fmt.Sprintf("check %q requires unknown check: %s", check.ID, reqID),
						LineNum: c.FindCheckNodeLine(check.ID, i),
					}
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
	// Also build a map of check ID to index for line number lookup
	graph := make(map[string][]string)
	idToIndex := make(map[string]int)
	for i, check := range c.Checks {
		graph[check.ID] = check.Requires
		idToIndex[check.ID] = i
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
			// Get line number of the first check in the cycle
			lineNum := c.FindCheckNodeLine(cyclePath[0], idToIndex[cyclePath[0]])
			return &ConfigError{
				Message: fmt.Sprintf("cyclic dependency detected: %s", formatCycle(cyclePath)),
				LineNum: lineNum,
			}
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

// FindCheckNodeLine returns the line number of a check in the YAML, or 0 if not found.
func (c *Config) FindCheckNodeLine(checkID string, checkIndex int) int {
	root, ok := c.yamlRoot.(*yaml.Node)
	if !ok || root == nil {
		return 0
	}

	// When unmarshaling into a Node, the root is a DocumentNode (Kind 1)
	// and its first Content element is the actual mapping
	var mapping *yaml.Node
	if root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
		mapping = root.Content[0]
	} else if root.Kind == yaml.MappingNode {
		mapping = root
	} else {
		return 0
	}

	if mapping.Kind != yaml.MappingNode {
		return 0
	}

	for i := 0; i < len(mapping.Content); i += 2 {
		if i+1 >= len(mapping.Content) {
			break
		}
		keyNode := mapping.Content[i]
		valueNode := mapping.Content[i+1]

		if keyNode.Value == "checks" && valueNode.Kind == yaml.SequenceNode {
			// Found the checks sequence, get the check at the given index
			if checkIndex >= 0 && checkIndex < len(valueNode.Content) {
				return valueNode.Content[checkIndex].Line
			}
			break
		}
	}

	return 0
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
