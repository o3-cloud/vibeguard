// Package assist provides AI agent-assisted setup functionality for VibeGuard.
package assist

import (
	"strings"
)

// ProjectAnalysis contains all information about a project needed for prompt generation.
type ProjectAnalysis struct {
	Name            string
	ProjectType     string
	Confidence      float64
	LanguageVersion string
	DetectedTools   []ToolInfo
	SourceDirs      []string
	TestDirs        []string
	EntryPoints     []string
	BuildOutputDir  string
}

// ToolInfo represents a detected tool in the project.
type ToolInfo struct {
	Name       string
	ConfigFile string
	Version    string
	Detected   bool
}

// CheckRecommendation represents a suggested check for the project.
type CheckRecommendation struct {
	ID          string
	Description string
	Rationale   string
	Command     string
	File        string // File to read output from instead of command stdout
	Grok        []string
	Assert      string
	Severity    string
	Suggestion  string
	Requires    []string
	Category    string
}

// Composer generates AI agent setup prompts from project analysis.
type Composer struct {
	analysis        *ProjectAnalysis
	recommendations []CheckRecommendation
}

// NewComposer creates a new Composer with the given project analysis.
func NewComposer(analysis *ProjectAnalysis, recommendations []CheckRecommendation) *Composer {
	return &Composer{
		analysis:        analysis,
		recommendations: recommendations,
	}
}

// Compose generates the complete AI agent setup prompt.
func (c *Composer) Compose() string {
	sections := c.buildSections()
	return c.assembleSections(sections)
}

// buildSections creates all sections of the prompt.
func (c *Composer) buildSections() []PromptSection {
	return []PromptSection{
		HeaderSection(),
		ProjectAnalysisSection(c.analysis),
		ToolingInspectionSection(c.analysis),
		ToolingResearchSection(c.analysis.ProjectType),
		RecommendationsSection(c.recommendations),
		ConfigRequirementsSection(),
		LanguageExamplesSection(c.analysis.ProjectType),
		ValidationRulesSection(),
		TaskSection(c.analysis.Name),
	}
}

// assembleSections combines all sections into the final prompt string.
func (c *Composer) assembleSections(sections []PromptSection) string {
	var sb strings.Builder

	for i, section := range sections {
		sb.WriteString(section.Content)
		if i < len(sections)-1 {
			sb.WriteString("\n\n---\n\n")
		}
	}

	return sb.String()
}

// ComposeWithOptions allows customizing which sections to include.
func (c *Composer) ComposeWithOptions(opts ComposerOptions) string {
	var sections []PromptSection

	if opts.IncludeHeader {
		sections = append(sections, HeaderSection())
	}
	if opts.IncludeAnalysis {
		sections = append(sections, ProjectAnalysisSection(c.analysis))
	}
	if opts.IncludeToolingInspection {
		sections = append(sections, ToolingInspectionSection(c.analysis))
	}
	if opts.IncludeToolingResearch {
		sections = append(sections, ToolingResearchSection(c.analysis.ProjectType))
	}
	if opts.IncludeRecommendations {
		sections = append(sections, RecommendationsSection(c.recommendations))
	}
	if opts.IncludeRequirements {
		sections = append(sections, ConfigRequirementsSection())
	}
	if opts.IncludeExamples {
		sections = append(sections, LanguageExamplesSection(c.analysis.ProjectType))
	}
	if opts.IncludeValidation {
		sections = append(sections, ValidationRulesSection())
	}
	if opts.IncludeTask {
		sections = append(sections, TaskSection(c.analysis.Name))
	}

	return c.assembleSections(sections)
}

// ComposerOptions controls which sections to include in the prompt.
type ComposerOptions struct {
	IncludeHeader            bool
	IncludeAnalysis          bool
	IncludeToolingInspection bool
	IncludeToolingResearch   bool
	IncludeRecommendations   bool
	IncludeRequirements      bool
	IncludeExamples          bool
	IncludeValidation        bool
	IncludeTask              bool
}

// DefaultComposerOptions returns options with all sections enabled.
func DefaultComposerOptions() ComposerOptions {
	return ComposerOptions{
		IncludeHeader:            true,
		IncludeAnalysis:          true,
		IncludeToolingInspection: true,
		IncludeToolingResearch:   true,
		IncludeRecommendations:   true,
		IncludeRequirements:      true,
		IncludeExamples:          true,
		IncludeValidation:        true,
		IncludeTask:              true,
	}
}

// MinimalComposerOptions returns options for a minimal prompt (no validation details).
func MinimalComposerOptions() ComposerOptions {
	return ComposerOptions{
		IncludeHeader:            true,
		IncludeAnalysis:          true,
		IncludeToolingInspection: false,
		IncludeToolingResearch:   false,
		IncludeRecommendations:   true,
		IncludeRequirements:      true,
		IncludeExamples:          false,
		IncludeValidation:        false,
		IncludeTask:              true,
	}
}
