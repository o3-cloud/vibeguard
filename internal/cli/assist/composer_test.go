package assist

import (
	"strings"
	"testing"
)

func TestNewComposer(t *testing.T) {
	analysis := &ProjectAnalysis{
		Name:        "test-project",
		ProjectType: "go",
		Confidence:  0.95,
	}
	recs := []CheckRecommendation{
		{ID: "fmt", Description: "Format check", Category: "format"},
	}

	composer := NewComposer(analysis, recs)

	if composer.analysis != analysis {
		t.Error("NewComposer should store the analysis")
	}
	if len(composer.recommendations) != 1 {
		t.Error("NewComposer should store recommendations")
	}
}

func TestComposerCompose(t *testing.T) {
	analysis := &ProjectAnalysis{
		Name:            "github.com/example/myapp",
		ProjectType:     "go",
		Confidence:      1.0,
		LanguageVersion: "1.22",
		DetectedTools: []ToolInfo{
			{Name: "golangci-lint", ConfigFile: ".golangci.yml", Detected: true},
			{Name: "gofmt", Detected: true},
		},
		SourceDirs:     []string{"cmd", "internal"},
		TestDirs:       []string{"internal"},
		EntryPoints:    []string{"cmd/main.go"},
		BuildOutputDir: "bin",
	}
	recs := []CheckRecommendation{
		{
			ID:          "fmt",
			Description: "Check Go code formatting",
			Rationale:   "Consistent formatting",
			Command:     "gofmt -l .",
			Severity:    "error",
			Category:    "format",
		},
		{
			ID:          "lint",
			Description: "Run linter",
			Rationale:   "Catch bugs",
			Command:     "golangci-lint run",
			Severity:    "error",
			Requires:    []string{"fmt"},
			Category:    "lint",
		},
	}

	composer := NewComposer(analysis, recs)
	prompt := composer.Compose()

	// Check that all expected sections are present
	expectedSections := []string{
		"# VibeGuard AI Agent Setup Guide",
		"## Your Project",
		"**Project Name:** github.com/example/myapp",
		"analyze the project structure",
		"golangci-lint",
		"## Recommended Checks",
		"fmt (format)",
		"lint (lint)",
		"## Configuration Requirements",
		"## Go-Specific Examples",
		"## YAML Syntax Requirements",
		"## Your Task",
	}

	for _, expected := range expectedSections {
		if !strings.Contains(prompt, expected) {
			t.Errorf("Prompt should contain %q", expected)
		}
	}
}

func TestComposerComposeWithOptions(t *testing.T) {
	analysis := &ProjectAnalysis{
		Name:        "test-project",
		ProjectType: "node",
		Confidence:  0.8,
	}
	recs := []CheckRecommendation{}

	composer := NewComposer(analysis, recs)

	t.Run("minimal options", func(t *testing.T) {
		opts := MinimalComposerOptions()
		prompt := composer.ComposeWithOptions(opts)

		// Should include header, analysis, requirements, task
		if !strings.Contains(prompt, "# VibeGuard AI Agent Setup Guide") {
			t.Error("Minimal should include header")
		}
		if !strings.Contains(prompt, "## Your Project") {
			t.Error("Minimal should include analysis")
		}
		if !strings.Contains(prompt, "## Your Task") {
			t.Error("Minimal should include task")
		}

		// Should NOT include validation guide (too detailed)
		if strings.Contains(prompt, "## YAML Syntax Requirements") {
			t.Error("Minimal should not include validation details")
		}
	})

	t.Run("default options", func(t *testing.T) {
		opts := DefaultComposerOptions()
		prompt := composer.ComposeWithOptions(opts)

		// Should include all sections
		if !strings.Contains(prompt, "## YAML Syntax Requirements") {
			t.Error("Default should include validation rules")
		}
	})

	t.Run("selective options", func(t *testing.T) {
		opts := ComposerOptions{
			IncludeHeader:   true,
			IncludeAnalysis: true,
			IncludeTask:     true,
		}
		prompt := composer.ComposeWithOptions(opts)

		if !strings.Contains(prompt, "# VibeGuard AI Agent Setup Guide") {
			t.Error("Should include header")
		}
		if strings.Contains(prompt, "## Recommended Checks") {
			t.Error("Should not include recommendations when disabled")
		}
	})
}

func TestLanguageExamplesSection(t *testing.T) {
	testCases := []struct {
		projectType     string
		expectedContent string
	}{
		{"go", "## Go-Specific Examples"},
		{"node", "## Node.js-Specific Examples"},
		{"python", "## Python-Specific Examples"},
		{"rust", "## Rust-Specific Examples"},
		{"ruby", "## Ruby-Specific Examples"},
		{"java", "## Java-Specific Examples"},
		{"unknown", "## Generic Examples"},
	}

	for _, tc := range testCases {
		t.Run(tc.projectType, func(t *testing.T) {
			section := LanguageExamplesSection(tc.projectType)
			if !strings.Contains(section.Content, tc.expectedContent) {
				t.Errorf("Expected %q section for project type %q", tc.expectedContent, tc.projectType)
			}
		})
	}
}

func TestProjectAnalysisSection(t *testing.T) {
	analysis := &ProjectAnalysis{
		Name:            "my-app",
		ProjectType:     "python",
		Confidence:      0.85,
		LanguageVersion: "3.11",
		DetectedTools: []ToolInfo{
			{Name: "black", ConfigFile: "pyproject.toml", Detected: true},
			{Name: "pytest", Detected: true},
		},
		SourceDirs:     []string{"src"},
		TestDirs:       []string{"tests"},
		EntryPoints:    []string{"src/main.py"},
		BuildOutputDir: "dist",
	}

	section := ProjectAnalysisSection(analysis)

	// New format delegates analysis to the agent, so check for guidance instead
	expectedContents := []string{
		"**Project Name:** my-app",
		"Initial Analysis Instructions",
		"Project Type & Language",
		"Detected Tools",
		"Project Structure",
		"Existing Quality Tools",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(section.Content, expected) {
			t.Errorf("Project analysis section should contain %q", expected)
		}
	}
}

func TestRecommendationsSection(t *testing.T) {
	recs := []CheckRecommendation{
		{
			ID:          "test",
			Description: "Run tests",
			Rationale:   "Ensure code works",
			Command:     "pytest",
			Grok:        []string{"%{NUMBER:coverage}%"},
			Assert:      "coverage >= 70",
			Severity:    "error",
			Suggestion:  "Fix failing tests",
			Requires:    []string{"lint"},
			Category:    "test",
		},
	}

	section := RecommendationsSection(recs)

	expectedContents := []string{
		"### test (test)",
		"**Description:** Run tests",
		"**Rationale:** Ensure code works",
		"**Command:** `pytest`",
		"**Severity:** error",
		"**Grok Patterns:** `%{NUMBER:coverage}%`",
		"**Assertion:** `coverage >= 70`",
		"**Suggestion on failure:** Fix failing tests",
		"**Requires:** lint",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(section.Content, expected) {
			t.Errorf("Recommendations section should contain %q", expected)
		}
	}
}

func TestHeaderSection(t *testing.T) {
	section := HeaderSection()

	if section.Title != "Header" {
		t.Errorf("Expected title 'Header', got %q", section.Title)
	}
	if !strings.Contains(section.Content, "# VibeGuard AI Agent Setup Guide") {
		t.Error("Header section should contain title")
	}
	if !strings.Contains(section.Content, "policy enforcement") {
		t.Error("Header section should mention policy enforcement")
	}
}

func TestConfigRequirementsSection(t *testing.T) {
	section := ConfigRequirementsSection()

	expectedContents := []string{
		"## Configuration Requirements",
		"version: \"1\"",
		"vars:",
		"Checks (required)",
		"**id:**",
		"**run:**",
		"**grok:**",
		"**assert:**",
		"**requires:**",
		"**severity:**",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(section.Content, expected) {
			t.Errorf("Config requirements section should contain %q", expected)
		}
	}
}

func TestTaskSection(t *testing.T) {
	section := TaskSection("my-project")

	if !strings.Contains(section.Content, "## Your Task") {
		t.Error("Task section should have title")
	}
	if !strings.Contains(section.Content, "my-project") {
		t.Error("Task section should mention project name")
	}
	if !strings.Contains(section.Content, "Option A: Use a Template") {
		t.Error("Task section should describe template approach")
	}
	if !strings.Contains(section.Content, "vibeguard init --list-templates") {
		t.Error("Task section should show how to list templates")
	}
}

func TestCapitalizeFirst(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"go", "Go"},
		{"node", "Node"},
		{"", ""},
		{"PYTHON", "PYTHON"},
		{"rust", "Rust"},
	}

	for _, tc := range testCases {
		result := capitalizeFirst(tc.input)
		if result != tc.expected {
			t.Errorf("capitalizeFirst(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestFormatGrokPatterns(t *testing.T) {
	patterns := []string{"%{NUMBER:count}", "%{WORD:name}"}
	result := formatGrokPatterns(patterns)

	if result != "`%{NUMBER:count}` `%{WORD:name}`" {
		t.Errorf("formatGrokPatterns returned unexpected result: %q", result)
	}
}

func TestToolingInspectionSection(t *testing.T) {
	t.Run("with config files", func(t *testing.T) {
		analysis := &ProjectAnalysis{
			Name:        "test-project",
			ProjectType: "go",
			DetectedTools: []ToolInfo{
				{Name: "golangci-lint", ConfigFile: ".golangci.yml", Detected: true},
				{Name: "gofmt", Detected: true}, // no config file
			},
		}

		section := ToolingInspectionSection(analysis)

		expectedContents := []string{
			"## Tooling Inspection Instructions",
			"### Configuration Files to Analyze",
			"- **golangci-lint**: Read `.golangci.yml`",
			"### What to Look For",
			"Enabled rules/checks",
			"Disabled rules",
			"Custom settings",
		}

		for _, expected := range expectedContents {
			if !strings.Contains(section.Content, expected) {
				t.Errorf("Tooling inspection section should contain %q", expected)
			}
		}

		// gofmt should NOT be in config files to analyze (no config file)
		if strings.Contains(section.Content, "- **gofmt**: Read") {
			t.Error("Tools without config files should not be listed in config files to analyze")
		}
	})

	t.Run("without config files", func(t *testing.T) {
		analysis := &ProjectAnalysis{
			Name:        "test-project",
			ProjectType: "go",
			DetectedTools: []ToolInfo{
				{Name: "gofmt", Detected: true},
				{Name: "go test", Detected: true},
			},
		}

		section := ToolingInspectionSection(analysis)

		// Should still have the section but no config files list
		if !strings.Contains(section.Content, "## Tooling Inspection Instructions") {
			t.Error("Should still contain title")
		}
		if strings.Contains(section.Content, "### Configuration Files to Analyze") {
			t.Error("Should not have config files section when no tools have configs")
		}
	})
}

func TestToolingResearchSection(t *testing.T) {
	testCases := []struct {
		projectType     string
		expectedTools   []string
		unexpectedTools []string
	}{
		{
			projectType:   "go",
			expectedTools: []string{"staticcheck", "gosec", "errcheck", "govulncheck"},
		},
		{
			projectType:   "node",
			expectedTools: []string{"tsc --noEmit", "npm audit", "depcheck", "madge"},
		},
		{
			projectType:   "python",
			expectedTools: []string{"bandit", "safety", "vulture", "radon"},
		},
		{
			projectType:   "rust",
			expectedTools: []string{"cargo audit", "cargo deny", "cargo outdated"},
		},
		{
			projectType:   "ruby",
			expectedTools: []string{"brakeman", "bundler-audit", "reek"},
		},
		{
			projectType:   "java",
			expectedTools: []string{"SpotBugs", "OWASP Dependency-Check", "PMD"},
		},
		{
			projectType:   "unknown",
			expectedTools: []string{"gitleaks", "trivy"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.projectType, func(t *testing.T) {
			section := ToolingResearchSection(tc.projectType)

			if !strings.Contains(section.Content, "## Additional Tooling Research") {
				t.Error("Should contain section title")
			}

			for _, tool := range tc.expectedTools {
				if !strings.Contains(section.Content, tool) {
					t.Errorf("Section for %s should contain tool %q", tc.projectType, tool)
				}
			}

			// Check that it has the guidance for presenting suggestions
			if !strings.Contains(section.Content, "### How to Present Suggestions") {
				t.Error("Should contain guidance on presenting suggestions")
			}
			if !strings.Contains(section.Content, "Would you like me to include checks") {
				t.Error("Should contain prompt for user about including additional tools")
			}
		})
	}
}

func TestGetToolSuggestions(t *testing.T) {
	// Test that each project type returns non-empty suggestions
	projectTypes := []string{"go", "node", "python", "rust", "ruby", "java", "unknown"}

	for _, pt := range projectTypes {
		suggestions := getToolSuggestions(pt)
		if len(suggestions) == 0 {
			t.Errorf("getToolSuggestions(%q) returned empty slice", pt)
		}

		// Verify each suggestion has required fields
		for _, s := range suggestions {
			if s.Name == "" {
				t.Errorf("Suggestion for %s has empty Name", pt)
			}
			if s.Category == "" {
				t.Errorf("Suggestion %s for %s has empty Category", s.Name, pt)
			}
			if s.Purpose == "" {
				t.Errorf("Suggestion %s for %s has empty Purpose", s.Name, pt)
			}
			if s.Value == "" {
				t.Errorf("Suggestion %s for %s has empty Value", s.Name, pt)
			}
		}
	}
}

func TestPromptTokenEstimate(t *testing.T) {
	analysis := &ProjectAnalysis{
		Name:            "github.com/example/project",
		ProjectType:     "go",
		Confidence:      1.0,
		LanguageVersion: "1.22",
		DetectedTools: []ToolInfo{
			{Name: "golangci-lint", Detected: true},
			{Name: "gofmt", Detected: true},
			{Name: "go test", Detected: true},
		},
		SourceDirs:  []string{"cmd", "internal", "pkg"},
		TestDirs:    []string{"internal"},
		EntryPoints: []string{"cmd/main.go"},
	}
	recs := []CheckRecommendation{
		{ID: "fmt", Description: "Format", Command: "gofmt -l .", Category: "format", Severity: "error"},
		{ID: "vet", Description: "Vet", Command: "go vet ./...", Category: "lint", Severity: "error"},
		{ID: "lint", Description: "Lint", Command: "golangci-lint run", Category: "lint", Severity: "error"},
		{ID: "test", Description: "Test", Command: "go test ./...", Category: "test", Severity: "error"},
	}

	composer := NewComposer(analysis, recs)
	prompt := composer.Compose()

	// Rough estimate: 4 characters per token
	estimatedTokens := len(prompt) / 4

	// Target is < 5500 tokens (increased to accommodate expanded TaskSection with step-by-step guidance)
	if estimatedTokens > 5500 {
		t.Errorf("Prompt exceeds 5500 token estimate: ~%d tokens (%d chars)", estimatedTokens, len(prompt))
	}

	t.Logf("Prompt length: %d chars, ~%d tokens", len(prompt), estimatedTokens)
}

func TestAssembleSections_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		sections []PromptSection
		expected string
	}{
		{
			name:     "empty sections",
			sections: []PromptSection{},
			expected: "",
		},
		{
			name: "single section",
			sections: []PromptSection{
				{Content: "Only content"},
			},
			expected: "Only content",
		},
		{
			name: "two sections",
			sections: []PromptSection{
				{Content: "First"},
				{Content: "Second"},
			},
			expected: "First\n\n---\n\nSecond",
		},
		{
			name: "three sections",
			sections: []PromptSection{
				{Content: "First"},
				{Content: "Second"},
				{Content: "Third"},
			},
			expected: "First\n\n---\n\nSecond\n\n---\n\nThird",
		},
		{
			name: "single character sections",
			sections: []PromptSection{
				{Content: "A"},
				{Content: "B"},
			},
			expected: "A\n\n---\n\nB",
		},
		{
			name: "sections with special characters",
			sections: []PromptSection{
				{Content: "Content with\nnewline"},
				{Content: "Content with # markdown"},
			},
			expected: "Content with\nnewline\n\n---\n\nContent with # markdown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analysis := &ProjectAnalysis{Name: "test"}
			composer := NewComposer(analysis, []CheckRecommendation{})
			result := composer.assembleSections(tc.sections)

			if result != tc.expected {
				t.Errorf("assembleSections() mismatch\nExpected: %q\nGot:      %q", tc.expected, result)
			}
		})
	}
}

func TestAssembleSections_BoundaryCondition(t *testing.T) {
	analysis := &ProjectAnalysis{Name: "test"}
	composer := NewComposer(analysis, []CheckRecommendation{})

	// Test that separators are NOT added after the last section
	sections := []PromptSection{
		{Content: "Section1"},
		{Content: "Section2"},
		{Content: "Section3"},
	}

	result := composer.assembleSections(sections)

	// Should NOT end with separator
	if strings.HasSuffix(result, "\n\n---\n\n") {
		t.Error("assembleSections should not add separator after last section")
	}

	// Should contain exactly 2 separators for 3 sections
	separatorCount := strings.Count(result, "\n\n---\n\n")
	if separatorCount != 2 {
		t.Errorf("assembleSections should have 2 separators for 3 sections, got %d", separatorCount)
	}
}
