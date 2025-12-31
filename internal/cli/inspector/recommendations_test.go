package inspector

import (
	"testing"
)

func TestNewRecommender(t *testing.T) {
	tools := []ToolInfo{
		{Name: "golangci-lint", Detected: true},
		{Name: "go test", Detected: true},
	}

	r := NewRecommender(Go, tools)

	if r.projectType != Go {
		t.Errorf("expected project type Go, got %v", r.projectType)
	}
	if len(r.tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(r.tools))
	}
}

func TestRecommender_Recommend_GoProject(t *testing.T) {
	tools := []ToolInfo{
		{Name: "golangci-lint", Detected: true, Confidence: 0.9},
		{Name: "gofmt", Detected: true, Confidence: 1.0},
		{Name: "go vet", Detected: true, Confidence: 1.0},
		{Name: "go test", Detected: true, Confidence: 1.0},
	}

	r := NewRecommender(Go, tools)
	recs := r.Recommend()

	// Should have recommendations for all detected tools plus project-level ones
	if len(recs) == 0 {
		t.Error("expected recommendations, got none")
	}

	// Verify we have expected check IDs
	ids := make(map[string]bool)
	for _, rec := range recs {
		ids[rec.ID] = true
	}

	expectedIDs := []string{"lint", "fmt", "vet", "test", "coverage", "build"}
	for _, id := range expectedIDs {
		if !ids[id] {
			t.Errorf("expected recommendation with ID %q", id)
		}
	}
}

func TestRecommender_Recommend_NodeProject(t *testing.T) {
	tools := []ToolInfo{
		{Name: "eslint", Detected: true, Confidence: 0.9},
		{Name: "prettier", Detected: true, Confidence: 0.9},
		{Name: "jest", Detected: true, Confidence: 0.9},
		{Name: "typescript", Detected: true, Confidence: 0.95},
	}

	r := NewRecommender(Node, tools)
	recs := r.Recommend()

	if len(recs) == 0 {
		t.Error("expected recommendations, got none")
	}

	// Verify we have expected check IDs
	ids := make(map[string]bool)
	for _, rec := range recs {
		ids[rec.ID] = true
	}

	expectedIDs := []string{"lint", "fmt", "test", "typecheck"}
	for _, id := range expectedIDs {
		if !ids[id] {
			t.Errorf("expected recommendation with ID %q", id)
		}
	}
}

func TestRecommender_Recommend_PythonProject(t *testing.T) {
	tools := []ToolInfo{
		{Name: "black", Detected: true, Confidence: 0.9},
		{Name: "pylint", Detected: true, Confidence: 0.9},
		{Name: "pytest", Detected: true, Confidence: 0.9},
		{Name: "mypy", Detected: true, Confidence: 0.9},
	}

	r := NewRecommender(Python, tools)
	recs := r.Recommend()

	if len(recs) == 0 {
		t.Error("expected recommendations, got none")
	}

	// Verify we have expected check IDs
	ids := make(map[string]bool)
	for _, rec := range recs {
		ids[rec.ID] = true
	}

	expectedIDs := []string{"lint", "fmt", "test", "typecheck"}
	for _, id := range expectedIDs {
		if !ids[id] {
			t.Errorf("expected recommendation with ID %q", id)
		}
	}
}

func TestRecommender_Recommend_NoDetectedTools(t *testing.T) {
	tools := []ToolInfo{
		{Name: "eslint", Detected: false},
		{Name: "prettier", Detected: false},
	}

	r := NewRecommender(Node, tools)
	recs := r.Recommend()

	// Should still get project-level recommendations
	if len(recs) == 0 {
		t.Error("expected at least project-level recommendations")
	}

	// Verify no recommendations are for non-detected tools
	for _, rec := range recs {
		if rec.Tool == "eslint" || rec.Tool == "prettier" {
			t.Errorf("got recommendation for non-detected tool: %s", rec.Tool)
		}
	}
}

func TestRecommender_Recommend_SortedByPriority(t *testing.T) {
	tools := []ToolInfo{
		{Name: "go test", Detected: true},       // Priority 30, 35
		{Name: "golangci-lint", Detected: true}, // Priority 20
		{Name: "gofmt", Detected: true},         // Priority 10
	}

	r := NewRecommender(Go, tools)
	recs := r.Recommend()

	// Verify recommendations are sorted by priority
	for i := 1; i < len(recs); i++ {
		if recs[i].Priority < recs[i-1].Priority {
			t.Errorf("recommendations not sorted: %s (priority %d) comes after %s (priority %d)",
				recs[i].ID, recs[i].Priority, recs[i-1].ID, recs[i-1].Priority)
		}
	}
}

func TestRecommender_RecommendForCategory(t *testing.T) {
	tools := []ToolInfo{
		{Name: "golangci-lint", Detected: true},
		{Name: "gofmt", Detected: true},
		{Name: "go test", Detected: true},
	}

	r := NewRecommender(Go, tools)

	// Test lint category
	lintRecs := r.RecommendForCategory("lint")
	for _, rec := range lintRecs {
		if rec.Category != "lint" {
			t.Errorf("expected category lint, got %s", rec.Category)
		}
	}

	// Test format category
	fmtRecs := r.RecommendForCategory("format")
	for _, rec := range fmtRecs {
		if rec.Category != "format" {
			t.Errorf("expected category format, got %s", rec.Category)
		}
	}

	// Test test category
	testRecs := r.RecommendForCategory("test")
	for _, rec := range testRecs {
		if rec.Category != "test" {
			t.Errorf("expected category test, got %s", rec.Category)
		}
	}
}

func TestDeduplicateRecommendations(t *testing.T) {
	recs := []CheckRecommendation{
		{ID: "lint", Tool: "eslint"},
		{ID: "test", Tool: "jest"},
		{ID: "lint", Tool: "pylint"}, // Duplicate ID
		{ID: "fmt", Tool: "prettier"},
	}

	deduped := DeduplicateRecommendations(recs)

	if len(deduped) != 3 {
		t.Errorf("expected 3 unique recommendations, got %d", len(deduped))
	}

	// First lint should be kept (eslint)
	var lintRec *CheckRecommendation
	for i := range deduped {
		if deduped[i].ID == "lint" {
			lintRec = &deduped[i]
			break
		}
	}
	if lintRec == nil || lintRec.Tool != "eslint" {
		t.Error("expected first lint recommendation (eslint) to be kept")
	}
}

func TestFilterByTools(t *testing.T) {
	recs := []CheckRecommendation{
		{ID: "lint", Tool: "eslint"},
		{ID: "fmt", Tool: "prettier"},
		{ID: "test", Tool: "jest"},
		{ID: "typecheck", Tool: "typescript"},
	}

	filtered := FilterByTools(recs, []string{"eslint", "jest"})

	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered recommendations, got %d", len(filtered))
	}

	for _, rec := range filtered {
		if rec.Tool != "eslint" && rec.Tool != "jest" {
			t.Errorf("unexpected tool in filtered results: %s", rec.Tool)
		}
	}
}

func TestGroupByCategory(t *testing.T) {
	recs := []CheckRecommendation{
		{ID: "lint1", Category: "lint"},
		{ID: "lint2", Category: "lint"},
		{ID: "fmt", Category: "format"},
		{ID: "test", Category: "test"},
	}

	groups := GroupByCategory(recs)

	if len(groups["lint"]) != 2 {
		t.Errorf("expected 2 lint recommendations, got %d", len(groups["lint"]))
	}
	if len(groups["format"]) != 1 {
		t.Errorf("expected 1 format recommendation, got %d", len(groups["format"]))
	}
	if len(groups["test"]) != 1 {
		t.Errorf("expected 1 test recommendation, got %d", len(groups["test"]))
	}
}

func TestRecommendation_HasRequiredFields(t *testing.T) {
	tools := []ToolInfo{
		{Name: "golangci-lint", Detected: true},
		{Name: "gofmt", Detected: true},
		{Name: "go vet", Detected: true},
		{Name: "go test", Detected: true},
		{Name: "eslint", Detected: true},
		{Name: "prettier", Detected: true},
		{Name: "jest", Detected: true},
		{Name: "black", Detected: true},
		{Name: "pytest", Detected: true},
	}

	projectTypes := []ProjectType{Go, Node, Python}

	for _, pt := range projectTypes {
		r := NewRecommender(pt, tools)
		recs := r.Recommend()

		for _, rec := range recs {
			if rec.ID == "" {
				t.Errorf("recommendation missing ID for project type %s", pt)
			}
			if rec.Description == "" {
				t.Errorf("recommendation %s missing Description", rec.ID)
			}
			if rec.Rationale == "" {
				t.Errorf("recommendation %s missing Rationale", rec.ID)
			}
			if rec.Command == "" {
				t.Errorf("recommendation %s missing Command", rec.ID)
			}
			if rec.Severity == "" {
				t.Errorf("recommendation %s missing Severity", rec.ID)
			}
			if rec.Severity != "error" && rec.Severity != "warning" {
				t.Errorf("recommendation %s has invalid Severity: %s", rec.ID, rec.Severity)
			}
			if rec.Category == "" {
				t.Errorf("recommendation %s missing Category", rec.ID)
			}
			if rec.Tool == "" {
				t.Errorf("recommendation %s missing Tool", rec.ID)
			}
		}
	}
}

func TestRecommendation_GrokAndAssert(t *testing.T) {
	tools := []ToolInfo{
		{Name: "go test", Detected: true},
	}

	r := NewRecommender(Go, tools)
	recs := r.Recommend()

	// Find coverage recommendation (should have grok and assert)
	var coverageRec *CheckRecommendation
	for i := range recs {
		if recs[i].ID == "coverage" {
			coverageRec = &recs[i]
			break
		}
	}

	if coverageRec == nil {
		t.Fatal("coverage recommendation not found")
	}

	if len(coverageRec.Grok) == 0 {
		t.Error("coverage recommendation should have grok patterns")
	}
	if coverageRec.Assert == "" {
		t.Error("coverage recommendation should have assertion")
	}
	if len(coverageRec.Requires) == 0 {
		t.Error("coverage recommendation should have requires")
	}
	if coverageRec.Suggestion == "" {
		t.Error("coverage recommendation should have suggestion with template")
	}
}

func TestRecommender_RuffVsPylint(t *testing.T) {
	// When both ruff and pylint are detected, both should be recommended
	// (deduplication happens at a higher level based on user preference)
	tools := []ToolInfo{
		{Name: "ruff", Detected: true},
		{Name: "pylint", Detected: true},
	}

	r := NewRecommender(Python, tools)
	recs := r.Recommend()

	// Both should produce lint recommendations
	lintCount := 0
	for _, rec := range recs {
		if rec.Category == "lint" {
			lintCount++
		}
	}

	// We should have 2 lint recommendations (one for each tool)
	// Note: They both have ID "lint" so deduplication would keep only one
	// But the raw recommendations should include both
	if lintCount != 2 {
		t.Errorf("expected 2 lint recommendations (ruff and pylint), got %d", lintCount)
	}
}

func TestRecommender_HookToolsNoRecommendations(t *testing.T) {
	// Hook management tools shouldn't produce check recommendations
	// as they are meta-tools for running other hooks
	tools := []ToolInfo{
		{Name: "pre-commit", Detected: true},
		{Name: "husky", Detected: true},
		{Name: "lefthook", Detected: true},
	}

	r := NewRecommender(Unknown, tools)
	recs := r.Recommend()

	for _, rec := range recs {
		if rec.Tool == "pre-commit" || rec.Tool == "husky" || rec.Tool == "lefthook" {
			t.Errorf("hook tool %s should not produce recommendations", rec.Tool)
		}
	}
}

func TestRecommender_UnknownProjectType(t *testing.T) {
	tools := []ToolInfo{
		{Name: "eslint", Detected: true},
	}

	r := NewRecommender(Unknown, tools)
	recs := r.Recommend()

	// Should still get recommendations for detected tools
	if len(recs) == 0 {
		t.Error("expected recommendations even for unknown project type")
	}

	// Should have eslint recommendation
	found := false
	for _, rec := range recs {
		if rec.Tool == "eslint" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected eslint recommendation")
	}
}

func TestSortRecommendations(t *testing.T) {
	recs := []CheckRecommendation{
		{ID: "test", Priority: 30},
		{ID: "fmt", Priority: 10},
		{ID: "lint", Priority: 20},
		{ID: "build", Priority: 5},
	}

	sortRecommendations(recs)

	expected := []string{"build", "fmt", "lint", "test"}
	for i, rec := range recs {
		if rec.ID != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], rec.ID)
		}
	}
}

func TestRecommender_NpmAuditSecurity(t *testing.T) {
	tools := []ToolInfo{
		{Name: "npm audit", Detected: true},
	}

	r := NewRecommender(Node, tools)
	recs := r.Recommend()

	var securityRec *CheckRecommendation
	for i := range recs {
		if recs[i].ID == "security" {
			securityRec = &recs[i]
			break
		}
	}

	if securityRec == nil {
		t.Fatal("security recommendation not found")
	}

	if securityRec.Category != "security" {
		t.Errorf("expected category security, got %s", securityRec.Category)
	}
	if securityRec.Severity != "warning" {
		t.Errorf("expected severity warning, got %s", securityRec.Severity)
	}
}

func TestRecommender_Goimports(t *testing.T) {
	tools := []ToolInfo{
		{Name: "goimports", Detected: true},
	}

	r := NewRecommender(Go, tools)
	recs := r.Recommend()

	var importsRec *CheckRecommendation
	for i := range recs {
		if recs[i].ID == "imports" {
			importsRec = &recs[i]
			break
		}
	}

	if importsRec == nil {
		t.Fatal("imports recommendation not found")
	}

	if importsRec.Tool != "goimports" {
		t.Errorf("expected tool goimports, got %s", importsRec.Tool)
	}
	if importsRec.Category != "format" {
		t.Errorf("expected category format, got %s", importsRec.Category)
	}
}

func TestRecommender_MultipleTestFrameworks(t *testing.T) {
	// When multiple test frameworks are detected, all should be recommended
	// (user can choose which to keep)
	tools := []ToolInfo{
		{Name: "jest", Detected: true},
		{Name: "mocha", Detected: true},
		{Name: "vitest", Detected: true},
	}

	r := NewRecommender(Node, tools)
	recs := r.Recommend()

	testTools := make(map[string]bool)
	for _, rec := range recs {
		if rec.Category == "test" {
			testTools[rec.Tool] = true
		}
	}

	// All three test tools should have recommendations
	// (they all have ID "test" so would be deduplicated, but raw recommendations include all)
	if !testTools["jest"] {
		t.Error("expected jest recommendation")
	}
	if !testTools["mocha"] {
		t.Error("expected mocha recommendation")
	}
	if !testTools["vitest"] {
		t.Error("expected vitest recommendation")
	}
}

func TestRecommender_Isort(t *testing.T) {
	tools := []ToolInfo{
		{Name: "isort", Detected: true},
	}

	r := NewRecommender(Python, tools)
	recs := r.Recommend()

	var importsRec *CheckRecommendation
	for i := range recs {
		if recs[i].ID == "imports" && recs[i].Tool == "isort" {
			importsRec = &recs[i]
			break
		}
	}

	if importsRec == nil {
		t.Fatal("isort imports recommendation not found")
	}

	if importsRec.Tool != "isort" {
		t.Errorf("expected tool isort, got %s", importsRec.Tool)
	}
	if importsRec.Category != "format" {
		t.Errorf("expected category format, got %s", importsRec.Category)
	}
	if importsRec.Priority != 11 {
		t.Errorf("expected priority 11, got %d", importsRec.Priority)
	}
}

func TestRecommender_PipAuditSecurity(t *testing.T) {
	tools := []ToolInfo{
		{Name: "pip-audit", Detected: true},
	}

	r := NewRecommender(Python, tools)
	recs := r.Recommend()

	var securityRec *CheckRecommendation
	for i := range recs {
		if recs[i].ID == "security" && recs[i].Tool == "pip-audit" {
			securityRec = &recs[i]
			break
		}
	}

	if securityRec == nil {
		t.Fatal("pip-audit security recommendation not found")
	}

	if securityRec.Category != "security" {
		t.Errorf("expected category security, got %s", securityRec.Category)
	}
	if securityRec.Severity != "warning" {
		t.Errorf("expected severity warning, got %s", securityRec.Severity)
	}
	if securityRec.Priority != 50 {
		t.Errorf("expected priority 50, got %d", securityRec.Priority)
	}
}

func TestRecommender_PythonFullToolchain(t *testing.T) {
	// Test a complete Python toolchain with all common tools
	tools := []ToolInfo{
		{Name: "black", Detected: true},
		{Name: "isort", Detected: true},
		{Name: "pylint", Detected: true},
		{Name: "mypy", Detected: true},
		{Name: "pytest", Detected: true},
		{Name: "pip-audit", Detected: true},
	}

	r := NewRecommender(Python, tools)
	recs := r.Recommend()

	// Verify we have all expected categories
	categories := make(map[string]bool)
	for _, rec := range recs {
		categories[rec.Category] = true
	}

	expectedCategories := []string{"format", "lint", "typecheck", "test", "security"}
	for _, cat := range expectedCategories {
		if !categories[cat] {
			t.Errorf("expected category %q in recommendations", cat)
		}
	}

	// Verify tools are present
	toolsFound := make(map[string]bool)
	for _, rec := range recs {
		toolsFound[rec.Tool] = true
	}

	expectedTools := []string{"black", "isort", "pylint", "mypy", "pytest", "pip-audit"}
	for _, tool := range expectedTools {
		if !toolsFound[tool] {
			t.Errorf("expected tool %q in recommendations", tool)
		}
	}
}
