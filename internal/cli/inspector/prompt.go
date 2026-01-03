// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"github.com/vibeguard/vibeguard/internal/cli/assist"
)

// GenerateSetupPrompt creates a Claude Code-friendly setup prompt based on inspection results.
// This function uses the assist.Composer to generate the prompt.
func GenerateSetupPrompt(
	projectType *DetectionResult,
	tools []ToolInfo,
	metadata *ProjectMetadata,
	structure *ProjectStructure,
	recommendations []CheckRecommendation,
) (string, error) {
	// Convert inspector types to assist types
	analysis := convertToProjectAnalysis(projectType, tools, metadata, structure)
	assistRecs := convertToAssistRecommendations(recommendations)

	// Use the Composer to generate the prompt
	composer := assist.NewComposer(analysis, assistRecs)
	return composer.Compose(), nil
}

// convertToProjectAnalysis converts inspector types to assist.ProjectAnalysis.
func convertToProjectAnalysis(
	projectType *DetectionResult,
	tools []ToolInfo,
	metadata *ProjectMetadata,
	structure *ProjectStructure,
) *assist.ProjectAnalysis {
	// Convert tools
	assistTools := make([]assist.ToolInfo, 0, len(tools))
	for _, t := range tools {
		assistTools = append(assistTools, assist.ToolInfo{
			Name:       t.Name,
			ConfigFile: t.ConfigFile,
			Version:    t.Version,
			Detected:   t.Detected,
		})
	}

	return &assist.ProjectAnalysis{
		Name:            metadata.Name,
		ProjectType:     string(projectType.Type),
		Confidence:      projectType.Confidence,
		LanguageVersion: metadata.Version,
		DetectedTools:   assistTools,
		SourceDirs:      structure.SourceDirs,
		TestDirs:        structure.TestDirs,
		EntryPoints:     structure.EntryPoints,
		BuildOutputDir:  structure.BuildOutputDir,
	}
}

// convertToAssistRecommendations converts inspector recommendations to assist.CheckRecommendation.
func convertToAssistRecommendations(recs []CheckRecommendation) []assist.CheckRecommendation {
	result := make([]assist.CheckRecommendation, 0, len(recs))
	for _, r := range recs {
		result = append(result, assist.CheckRecommendation{
			ID:          r.ID,
			Description: r.Description,
			Rationale:   r.Rationale,
			Command:     r.Command,
			File:        r.File,
			Grok:        r.Grok,
			Assert:      r.Assert,
			Severity:    r.Severity,
			Suggestion:  r.Suggestion,
			Requires:    r.Requires,
			Category:    r.Category,
		})
	}
	return result
}
