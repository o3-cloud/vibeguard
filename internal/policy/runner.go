package policy

import (
	"context"
	"fmt"
	"time"

	"github.com/vibeguard/vibeguard/pkg/models"
)

// Runner evaluates policies against resources.
type Runner interface {
	// Evaluate runs a policy against a resource and returns the result
	Evaluate(ctx context.Context, policy *models.Policy, resource *models.Resource) (*models.EvaluationResult, error)
}

// SimpleRunner is a basic implementation of the Runner interface.
// It provides straightforward rule evaluation without advanced features.
type SimpleRunner struct{}

// NewSimpleRunner creates a new SimpleRunner instance.
func NewSimpleRunner() *SimpleRunner {
	return &SimpleRunner{}
}

// Evaluate implements the Runner interface.
// It evaluates each rule in the policy against the provided resource.
func (sr *SimpleRunner) Evaluate(ctx context.Context, policy *models.Policy, resource *models.Resource) (*models.EvaluationResult, error) {
	if policy == nil {
		return &models.EvaluationResult{
			Error: "policy is nil",
		}, fmt.Errorf("policy is nil")
	}

	if resource == nil {
		return &models.EvaluationResult{
			PolicyName: policy.Name,
			Error:      "resource is nil",
		}, fmt.Errorf("resource is nil")
	}

	result := &models.EvaluationResult{
		PolicyName:  policy.Name,
		RuleResults: make([]models.RuleResult, 0, len(policy.Rules)),
		Timestamp:   time.Now(),
	}

	// Evaluate each rule
	allPassed := true
	for _, rule := range policy.Rules {
		ruleResult := sr.evaluateRule(rule, resource)
		result.RuleResults = append(result.RuleResults, ruleResult)

		if !ruleResult.Passed {
			allPassed = false
		}
	}

	result.Passed = allPassed
	return result, nil
}

// evaluateRule evaluates a single rule against a resource.
// This is a simplified implementation that checks basic conditions.
func (sr *SimpleRunner) evaluateRule(rule models.Rule, resource *models.Resource) models.RuleResult {
	result := models.RuleResult{
		RuleID:          rule.ID,
		RuleDescription: rule.Description,
		Severity:        rule.Severity,
	}

	// Basic condition evaluation - this is a simple placeholder
	// A real implementation would parse and evaluate expressions
	// For now, we just check if the resource type matches the target
	if resource.Type != rule.Target {
		result.Passed = false
		result.Message = fmt.Sprintf("resource type %q does not match rule target %q", resource.Type, rule.Target)
		return result
	}

	// If target matches, consider the rule passed for this basic implementation
	result.Passed = true
	result.Message = fmt.Sprintf("rule %q passed", rule.ID)
	return result
}
