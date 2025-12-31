// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateSetupPrompt tests prompt generation for the vibeguard project itself.
// This is part of Phase 5 testing: verifying generated prompts work with Claude Code.
func TestGenerateSetupPrompt(t *testing.T) {
	// Get vibeguard root (this test is in internal/cli/inspector)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	vibeguardRoot := filepath.Join(cwd, "..", "..", "..")

	// Run full inspection
	detector := NewDetector(vibeguardRoot)
	projectType, err := detector.DetectPrimary()
	if err != nil {
		t.Fatalf("failed to detect project type: %v", err)
	}

	scanner := NewToolScanner(vibeguardRoot)
	tools, err := scanner.ScanForProjectType(projectType.Type)
	if err != nil {
		t.Fatalf("failed to scan tools: %v", err)
	}

	extractor := NewMetadataExtractor(vibeguardRoot)
	metadata, err := extractor.Extract(projectType.Type)
	if err != nil {
		t.Fatalf("failed to extract metadata: %v", err)
	}

	structure, err := extractor.ExtractStructure(projectType.Type)
	if err != nil {
		t.Fatalf("failed to extract structure: %v", err)
	}

	recommender := NewRecommender(projectType.Type, tools)
	recommendations := recommender.Recommend()
	recommendations = DeduplicateRecommendations(recommendations)

	// Generate the prompt
	prompt, err := GenerateSetupPrompt(projectType, tools, metadata, structure, recommendations)
	if err != nil {
		t.Fatalf("failed to generate prompt: %v", err)
	}

	// Log prompt length for analysis
	t.Logf("Generated prompt length: %d characters, ~%d tokens", len(prompt), len(prompt)/4)

	// Verify prompt contains expected sections
	requiredSections := []string{
		"# VibeGuard AI Agent Setup Guide",
		"## Project Analysis",
		"## Configuration Requirements",
		"## Validation Rules",
		"## Your Task",
	}

	for _, section := range requiredSections {
		if !strings.Contains(prompt, section) {
			t.Errorf("prompt missing required section: %s", section)
		}
	}

	// Write prompt to file for manual testing
	promptFile := filepath.Join(vibeguardRoot, "docs", "log", "generated-setup-prompt.md")
	if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
		t.Logf("warning: could not write prompt file: %v", err)
	} else {
		t.Logf("Prompt written to: %s", promptFile)
	}
}
