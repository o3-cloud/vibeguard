package policy

import (
	"context"
	"testing"

	"github.com/vibeguard/vibeguard/pkg/models"
)

func TestSimpleRunner_Evaluate(t *testing.T) {
	tests := []struct {
		name      string
		policy    *models.Policy
		resource  *models.Resource
		wantPassed bool
		wantErr   bool
	}{
		{
			name: "matching_target_passes",
			policy: &models.Policy{
				Name: "test-policy",
				Rules: []models.Rule{
					{
						ID:       "rule1",
						Target:   "commit",
						Severity: "error",
					},
				},
			},
			resource: &models.Resource{
				Type: "commit",
				Data: map[string]interface{}{},
			},
			wantPassed: true,
			wantErr:    false,
		},
		{
			name: "non_matching_target_fails",
			policy: &models.Policy{
				Name: "test-policy",
				Rules: []models.Rule{
					{
						ID:       "rule1",
						Target:   "deployment",
						Severity: "error",
					},
				},
			},
			resource: &models.Resource{
				Type: "commit",
				Data: map[string]interface{}{},
			},
			wantPassed: false,
			wantErr:    false,
		},
		{
			name: "nil_policy_returns_error",
			policy: nil,
			resource: &models.Resource{
				Type: "commit",
				Data: map[string]interface{}{},
			},
			wantPassed: false,
			wantErr:    true,
		},
		{
			name: "nil_resource_returns_error",
			policy: &models.Policy{
				Name: "test-policy",
				Rules: []models.Rule{
					{
						ID:     "rule1",
						Target: "commit",
					},
				},
			},
			resource:   nil,
			wantPassed: false,
			wantErr:    true,
		},
		{
			name: "multiple_rules_all_pass",
			policy: &models.Policy{
				Name: "multi-rule-policy",
				Rules: []models.Rule{
					{
						ID:     "rule1",
						Target: "commit",
					},
					{
						ID:     "rule2",
						Target: "commit",
					},
				},
			},
			resource: &models.Resource{
				Type: "commit",
				Data: map[string]interface{}{},
			},
			wantPassed: true,
			wantErr:    false,
		},
	}

	sr := NewSimpleRunner()
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sr.Evaluate(ctx, tt.policy, tt.resource)

			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != nil && result.Passed != tt.wantPassed {
				t.Errorf("Evaluate() passed = %v, want %v", result.Passed, tt.wantPassed)
			}
		})
	}
}

func TestSimpleRunner_EvaluateRule(t *testing.T) {
	tests := []struct {
		name       string
		rule       models.Rule
		resource   *models.Resource
		wantPassed bool
		wantMsg    string
	}{
		{
			name: "matching_target",
			rule: models.Rule{
				ID:     "rule1",
				Target: "commit",
			},
			resource: &models.Resource{
				Type: "commit",
				Data: map[string]interface{}{},
			},
			wantPassed: true,
		},
		{
			name: "non_matching_target",
			rule: models.Rule{
				ID:     "rule1",
				Target: "deployment",
			},
			resource: &models.Resource{
				Type: "commit",
				Data: map[string]interface{}{},
			},
			wantPassed: false,
		},
	}

	sr := NewSimpleRunner()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sr.evaluateRule(tt.rule, tt.resource)

			if result.Passed != tt.wantPassed {
				t.Errorf("evaluateRule() passed = %v, want %v", result.Passed, tt.wantPassed)
			}
		})
	}
}
