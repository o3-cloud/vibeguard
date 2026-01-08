package assist

import (
	"strings"
	"testing"
)

func TestProjectAnalysisSection_EmptyAnalysis(t *testing.T) {
	// Test boundary: empty project analysis
	analysis := &ProjectAnalysis{
		Name:        "empty-project",
		ProjectType: "unknown",
		Confidence:  0,
	}

	section := ProjectAnalysisSection(analysis)

	if !strings.Contains(section.Content, "**Project Name:** empty-project") {
		t.Error("should contain project name")
	}

	// New format: should guide agent to analyze the project
	if !strings.Contains(section.Content, "analyze the project structure") {
		t.Error("should contain guidance to analyze project")
	}
	if !strings.Contains(section.Content, "Project Type & Language") {
		t.Error("should ask agent to identify project type")
	}
}

func TestProjectAnalysisSection_WithTools(t *testing.T) {
	// Test boundary: len(DetectedTools) > 0
	// Note: New implementation delegates tool detection to the agent
	analysis := &ProjectAnalysis{
		Name:        "tool-project",
		ProjectType: "go",
		Confidence:  1.0,
		DetectedTools: []ToolInfo{
			{Name: "golangci-lint", ConfigFile: ".golangci.yml", Detected: true},
			{Name: "gofmt", ConfigFile: "", Detected: true},
		},
	}

	section := ProjectAnalysisSection(analysis)

	// New format: should ask agent to analyze tools
	if !strings.Contains(section.Content, "Detected Tools") {
		t.Error("should ask agent to identify detected tools")
	}
	if !strings.Contains(section.Content, "configuration files") {
		t.Error("should mention looking for configuration files")
	}
}

func TestProjectAnalysisSection_WithDirectories(t *testing.T) {
	// Test boundaries: len(SourceDirs) > 0, len(TestDirs) > 0, len(EntryPoints) > 0
	// Note: New implementation delegates directory analysis to the agent
	analysis := &ProjectAnalysis{
		Name:           "dir-project",
		ProjectType:    "node",
		SourceDirs:     []string{"src", "lib"},
		TestDirs:       []string{"test", "spec"},
		EntryPoints:    []string{"src/index.js"},
		BuildOutputDir: "dist",
	}

	section := ProjectAnalysisSection(analysis)

	// New format: should ask agent to identify directories
	if !strings.Contains(section.Content, "Project Structure") {
		t.Error("should ask agent to identify project structure")
	}
	if !strings.Contains(section.Content, "directory structure") {
		t.Error("should mention examining directory structure")
	}
}

func TestProjectAnalysisSection_OnlyTestDirs(t *testing.T) {
	// Test boundary: len(TestDirs) > 0 in isolation
	// Note: New implementation delegates directory analysis to the agent
	analysis := &ProjectAnalysis{
		Name:        "test-only",
		ProjectType: "python",
		TestDirs:    []string{"tests"},
	}

	section := ProjectAnalysisSection(analysis)

	// New format: should ask agent to identify all directories
	if !strings.Contains(section.Content, "Initial Analysis Instructions") {
		t.Error("should ask agent for initial analysis")
	}
}

func TestProjectAnalysisSection_OnlyEntryPoints(t *testing.T) {
	// Test boundary: len(EntryPoints) > 0 in isolation
	// Note: New implementation delegates directory analysis to the agent
	analysis := &ProjectAnalysis{
		Name:        "entry-only",
		ProjectType: "rust",
		EntryPoints: []string{"src/main.rs"},
	}

	section := ProjectAnalysisSection(analysis)

	// New format: should ask agent to identify entry points
	if !strings.Contains(section.Content, "entry points") {
		t.Error("should ask agent to identify entry points")
	}
}

func TestRecommendationsSection_MinimalCheck(t *testing.T) {
	// Test boundary: empty optional fields
	recs := []CheckRecommendation{
		{
			ID:          "basic-check",
			Description: "Basic check",
			Rationale:   "For testing",
			Command:     "true",
			Severity:    "error",
			Category:    "test",
			// File, Grok, Assert, Suggestion, Requires are empty
		},
	}

	section := RecommendationsSection(recs)

	if !strings.Contains(section.Content, "### basic-check (test)") {
		t.Error("should include check ID and category")
	}
	if !strings.Contains(section.Content, "**Description:** Basic check") {
		t.Error("should include description")
	}

	// Should not include optional fields
	if strings.Contains(section.Content, "**File:**") {
		t.Error("should not include File when empty")
	}
	if strings.Contains(section.Content, "**Grok Patterns:**") {
		t.Error("should not include Grok when empty")
	}
	if strings.Contains(section.Content, "**Assertion:**") {
		t.Error("should not include Assertion when empty")
	}
	if strings.Contains(section.Content, "**Suggestion on failure:**") {
		t.Error("should not include Suggestion when empty")
	}
	if strings.Contains(section.Content, "**Requires:**") {
		t.Error("should not include Requires when empty")
	}
}

func TestRecommendationsSection_FullCheck(t *testing.T) {
	// Test boundaries: rec.File != "", len(rec.Grok) > 0, rec.Assert != "", rec.Suggestion != "", len(rec.Requires) > 0
	recs := []CheckRecommendation{
		{
			ID:          "full-check",
			Description: "Full check",
			Rationale:   "Complete example",
			Command:     "test-cmd",
			Category:    "test",
			Severity:    "error",
			File:        "test-report.txt",
			Grok:        []string{"pattern1", "pattern2"},
			Assert:      "value > 0",
			Suggestion:  "Fix the issue",
			Requires:    []string{"dep1", "dep2"},
		},
	}

	section := RecommendationsSection(recs)

	if !strings.Contains(section.Content, "**File:** `test-report.txt`") {
		t.Error("should include File when present")
	}
	if !strings.Contains(section.Content, "**Grok Patterns:**") {
		t.Error("should include Grok when present")
	}
	if !strings.Contains(section.Content, "**Assertion:** `value > 0`") {
		t.Error("should include Assertion when present")
	}
	if !strings.Contains(section.Content, "**Suggestion on failure:** Fix the issue") {
		t.Error("should include Suggestion when present")
	}
	if !strings.Contains(section.Content, "**Requires:** dep1, dep2") {
		t.Error("should include Requires when present")
	}
}

func TestRecommendationsSection_MultipleChecks(t *testing.T) {
	// Test with multiple checks to ensure proper formatting
	recs := []CheckRecommendation{
		{
			ID:          "check1",
			Description: "First check",
			Rationale:   "R1",
			Command:     "cmd1",
			Category:    "cat1",
			Severity:    "error",
		},
		{
			ID:          "check2",
			Description: "Second check",
			Rationale:   "R2",
			Command:     "cmd2",
			Category:    "cat2",
			Severity:    "warning",
		},
	}

	section := RecommendationsSection(recs)

	if !strings.Contains(section.Content, "### check1 (cat1)") {
		t.Error("should include first check")
	}
	if !strings.Contains(section.Content, "### check2 (cat2)") {
		t.Error("should include second check")
	}
}

func TestValidationRulesSection_Content(t *testing.T) {
	// Test ValidationRulesSection
	section := ValidationRulesSection()

	// Should be valid HTML/markdown
	if len(section.Content) == 0 {
		t.Error("ValidationRulesSection should not be empty")
	}

	if !strings.Contains(section.Content, "YAML") {
		t.Error("should contain YAML references")
	}
}

func TestConfigRequirementsSection_Content(t *testing.T) {
	// Test ConfigRequirementsSection
	section := ConfigRequirementsSection()

	if len(section.Content) == 0 {
		t.Error("ConfigRequirementsSection should not be empty")
	}

	if !strings.Contains(section.Content, "Configuration") {
		t.Error("should contain Configuration references")
	}
}

func TestTaskSection_WithProjectName(t *testing.T) {
	testCases := []struct {
		projectName string
	}{
		{"github.com/example/myapp"},
		{"my-project"},
		{""},
	}

	for _, tc := range testCases {
		t.Run(tc.projectName, func(t *testing.T) {
			task := TaskSection(tc.projectName)

			if len(task.Content) == 0 {
				t.Error("TaskSection should not be empty")
			}

			if !strings.Contains(task.Content, "Task") {
				t.Error("should contain task heading")
			}
		})
	}
}

func TestLanguageExamplesSection_Coverage(t *testing.T) {
	// Test various language types to ensure boundary conditions
	testCases := []struct {
		projectType   string
		shouldContain string
	}{
		{"go", "Go"},
		{"node", "Node.js"},
		{"python", "Python"},
		{"rust", "Rust"},
		{"ruby", "Ruby"},
		{"java", "Java"},
		{"unknown-lang", "Generic"},
		{"", "Generic"},
	}

	for _, tc := range testCases {
		t.Run(tc.projectType, func(t *testing.T) {
			section := LanguageExamplesSection(tc.projectType)

			if !strings.Contains(section.Content, tc.shouldContain) {
				t.Errorf("expected %q in section for project type %q", tc.shouldContain, tc.projectType)
			}
		})
	}
}

func TestToolingResearchSection_WithSuggestions(t *testing.T) {
	// Test boundary: len(suggestions) > 0
	// Test with project types that have suggestions
	testCases := []string{"go", "node", "python"}

	for _, pt := range testCases {
		t.Run(pt, func(t *testing.T) {
			section := ToolingResearchSection(pt)

			// For known languages, should have suggestions
			if len(section.Content) == 0 {
				t.Error("should have content for known project type")
			}
			// Most will have "Suggested Tools" section if they have suggestions
			if strings.Contains(section.Content, "Suggested") {
				// Should contain tool suggestion details
				if !strings.Contains(section.Content, "Purpose") && !strings.Contains(section.Content, "Value") {
					t.Error("should contain tool suggestion details")
				}
			}
		})
	}
}

func TestRecommendationsSection_CommandField(t *testing.T) {
	// Test boundary: suggestion.Command != "" for individual commands
	// Note: This is in ToolingResearchSection context, testing command presence/absence
	recs := []CheckRecommendation{
		{
			ID:          "with-command",
			Description: "Has command",
			Rationale:   "Test",
			Command:     "example-command",
			Category:    "test",
			Severity:    "error",
		},
	}

	section := RecommendationsSection(recs)

	if !strings.Contains(section.Content, "**Command:** `example-command`") {
		t.Error("should include command when provided")
	}
}
