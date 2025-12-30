package models

import (
	"time"
)

// Policy represents a declarative policy that can be evaluated against resources.
type Policy struct {
	// Name is the unique identifier for this policy
	Name string `yaml:"name" json:"name"`

	// Version is the policy schema version
	Version string `yaml:"version" json:"version"`

	// Description explains the purpose of this policy
	Description string `yaml:"description" json:"description"`

	// Rules contains the set of rules to evaluate
	Rules []Rule `yaml:"rules" json:"rules"`

	// Metadata contains optional policy metadata
	Metadata map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// Rule represents a single rule within a policy.
type Rule struct {
	// ID is the unique identifier for this rule
	ID string `yaml:"id" json:"id"`

	// Description explains what this rule checks
	Description string `yaml:"description" json:"description"`

	// Target specifies what this rule applies to (e.g., "git.commit", "deployment")
	Target string `yaml:"target" json:"target"`

	// Condition is a simple expression to evaluate (supports basic operators)
	Condition string `yaml:"condition" json:"condition"`

	// Severity indicates the importance of this rule (e.g., "error", "warning", "info")
	Severity string `yaml:"severity" json:"severity"`

	// Message is displayed when the rule fails
	Message string `yaml:"message" json:"message"`
}

// EvaluationResult represents the result of evaluating a policy against a resource.
type EvaluationResult struct {
	// PolicyName is the name of the policy evaluated
	PolicyName string `json:"policy_name"`

	// Passed indicates if all rules passed
	Passed bool `json:"passed"`

	// RuleResults contains results for each rule
	RuleResults []RuleResult `json:"rule_results"`

	// Timestamp is when the evaluation occurred
	Timestamp time.Time `json:"timestamp"`

	// Error contains any error that occurred during evaluation
	Error string `json:"error,omitempty"`
}

// RuleResult represents the result of evaluating a single rule.
type RuleResult struct {
	// RuleID is the ID of the rule evaluated
	RuleID string `json:"rule_id"`

	// RuleDescription is the description of the rule
	RuleDescription string `json:"rule_description"`

	// Passed indicates if this rule passed
	Passed bool `json:"passed"`

	// Message is the output message (especially if rule failed)
	Message string `json:"message,omitempty"`

	// Severity of the rule
	Severity string `json:"severity"`
}

// Resource represents the input to be evaluated by a policy.
type Resource struct {
	// Type indicates the resource type (e.g., "commit", "deployment", "pr")
	Type string `json:"type"`

	// Data contains the resource data to evaluate
	Data map[string]interface{} `json:"data"`
}
