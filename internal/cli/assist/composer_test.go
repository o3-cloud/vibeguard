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
		"## Project Analysis",
		"**Project Name:** github.com/example/myapp",
		"**Project Type:** go",
		"**Detection Confidence:** 100%",
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
		if !strings.Contains(prompt, "## Project Analysis") {
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

	expectedContents := []string{
		"**Project Name:** my-app",
		"**Project Type:** python",
		"**Detection Confidence:** 85%",
		"**Language Version:** 3.11",
		"- black (config: pyproject.toml)",
		"- pytest",
		"- Source directories: src",
		"- Test directories: tests",
		"- Entry points: src/main.py",
		"- Build output: dist",
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
	if !strings.Contains(section.Content, "# vibeguard.yaml for my-project") {
		t.Error("Task section should mention project name")
	}
	if !strings.Contains(section.Content, "version: \"1\"") {
		t.Error("Task section should show version requirement")
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

	// Target is < 4000 tokens per spec
	if estimatedTokens > 4000 {
		t.Errorf("Prompt exceeds 4000 token estimate: ~%d tokens (%d chars)", estimatedTokens, len(prompt))
	}

	t.Logf("Prompt length: %d chars, ~%d tokens", len(prompt), estimatedTokens)
}
