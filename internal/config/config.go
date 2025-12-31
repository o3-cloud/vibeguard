package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Load reads and parses a VibeGuard configuration file.
// If path is empty, it searches for config files in the default locations.
func Load(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = findConfigFile()
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	cfg.applyDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

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

	return nil
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
