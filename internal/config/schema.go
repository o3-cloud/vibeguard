// Package config provides configuration loading and validation for VibeGuard.
package config

import "time"

// Config represents the complete VibeGuard configuration.
type Config struct {
	Version string            `yaml:"version"`
	Vars    map[string]string `yaml:"vars"`
	Checks  []Check           `yaml:"checks"`
	// yamlRoot stores the parsed YAML node tree for line number lookups (not exported)
	yamlRoot interface{} `yaml:"-"`
}

// Check represents a single check to execute.
type Check struct {
	ID         string   `yaml:"id"`
	Run        string   `yaml:"run"`
	Grok       GrokSpec `yaml:"grok"`
	File       string   `yaml:"file"`
	Assert     string   `yaml:"assert"`
	Severity   Severity `yaml:"severity"`
	Suggestion string   `yaml:"suggestion"`
	Fix        string   `yaml:"fix,omitempty"`
	Requires   []string `yaml:"requires"`
	Tags       []string `yaml:"tags,omitempty"`
	Timeout    Duration `yaml:"timeout"`
}

// Severity represents the severity level of a check failure.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// GrokSpec allows grok to be either a single string or a list of strings.
type GrokSpec []string

// Duration is a wrapper around time.Duration that supports YAML unmarshaling.
type Duration time.Duration

// DefaultTimeout is the default timeout for checks.
const DefaultTimeout = 30 * time.Second

// DefaultParallel is the default number of parallel checks.
const DefaultParallel = 4

// ConfigFileNames is the list of config file names to search for, in order.
var ConfigFileNames = []string{
	"vibeguard.yaml",
	"vibeguard.yml",
	".vibeguard.yaml",
	".vibeguard.yml",
}
