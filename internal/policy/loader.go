package policy

import (
	"fmt"
	"os"

	"github.com/vibeguard/vibeguard/pkg/models"
	"gopkg.in/yaml.v3"
)

// Loader loads policies from various sources.
type Loader interface {
	// Load reads a policy from the specified path
	Load(path string) (*models.Policy, error)
}

// YAMLLoader loads policies from YAML files.
type YAMLLoader struct{}

// NewYAMLLoader creates a new YAMLLoader instance.
func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

// Load reads and parses a YAML policy file.
func (yl *YAMLLoader) Load(path string) (*models.Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %w", err)
	}

	policy := &models.Policy{}
	if err := yaml.Unmarshal(data, policy); err != nil {
		return nil, fmt.Errorf("failed to parse policy YAML: %w", err)
	}

	// Validate policy structure
	if policy.Name == "" {
		return nil, fmt.Errorf("policy name is required")
	}

	if len(policy.Rules) == 0 {
		return nil, fmt.Errorf("policy must contain at least one rule")
	}

	// Validate rules
	for _, rule := range policy.Rules {
		if rule.ID == "" {
			return nil, fmt.Errorf("all rules must have an id")
		}
		if rule.Target == "" {
			return nil, fmt.Errorf("rule %q: target is required", rule.ID)
		}
		if rule.Severity == "" {
			rule.Severity = "error" // Default severity
		}
	}

	return policy, nil
}
