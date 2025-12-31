// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestRealWorldProjects tests the inspector on popular real-world projects.
// This test clones actual repositories and runs the inspector on them.
// Skip with: go test -short
func TestRealWorldProjects(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real-world project tests in short mode")
	}

	testCases := []struct {
		name         string
		repo         string
		expectedType ProjectType
		expectedMin  int // Minimum expected tools to find
	}{
		// Go projects
		{
			name:         "go/cli",
			repo:         "https://github.com/cli/cli",
			expectedType: Go,
			expectedMin:  3,
		},
		// Node projects
		{
			name:         "node/express",
			repo:         "https://github.com/expressjs/express",
			expectedType: Node,
			expectedMin:  2,
		},
		// Python projects
		{
			name:         "python/requests",
			repo:         "https://github.com/psf/requests",
			expectedType: Python,
			expectedMin:  2,
		},
	}

	// Create a temporary directory for all clones
	tmpRoot := t.TempDir()

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Clone with shallow depth to speed up test
			dir := filepath.Join(tmpRoot, strings.ReplaceAll(tc.name, "/", "-"))
			cmd := exec.Command("git", "clone", "--depth=1", tc.repo, dir)
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				t.Skipf("Failed to clone %s (requires network): %v", tc.repo, err)
			}

			// Run detector
			detector := NewDetector(dir)
			result, err := detector.DetectPrimary()
			if err != nil {
				t.Fatalf("DetectPrimary failed: %v", err)
			}
			if result.Type != tc.expectedType {
				t.Errorf("Expected %s, got %s (confidence: %.2f, indicators: %v)",
					tc.expectedType, result.Type, result.Confidence, result.Indicators)
			}

			// Run tool scanner
			scanner := NewToolScanner(dir)
			tools, err := scanner.ScanAll()
			if err != nil {
				t.Fatalf("ScanAll failed: %v", err)
			}
			if len(tools) < tc.expectedMin {
				t.Errorf("Expected at least %d tools, got %d: %v",
					tc.expectedMin, len(tools), toolNames(tools))
			}

			// Run metadata extractor
			extractor := NewMetadataExtractor(dir)
			metadata, err := extractor.Extract(result.Type)
			if err != nil {
				t.Fatalf("Extract failed: %v", err)
			}
			t.Logf("Project: %s, Version: %s", metadata.Name, metadata.Version)

			// Run recommender
			recommender := NewRecommender(result.Type, tools)
			recs := recommender.Recommend()
			t.Logf("Recommendations (%d):", len(recs))
			for _, rec := range recs {
				t.Logf("  - [%s] %s: %s", rec.Category, rec.ID, rec.Description)
			}
		})
	}
}

// TestVibeguardSelfInspection tests the inspector on the vibeguard project itself.
func TestVibeguardSelfInspection(t *testing.T) {
	// Find the project root (go up from the test directory)
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Go up until we find go.mod
	for dir != "/" {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		dir = filepath.Dir(dir)
	}

	if dir == "/" {
		t.Skip("Could not find project root")
	}

	t.Logf("Testing vibeguard at: %s", dir)

	// Run detector
	detector := NewDetector(dir)
	result, err := detector.DetectPrimary()
	if err != nil {
		t.Fatalf("DetectPrimary failed: %v", err)
	}
	if result.Type != Go {
		t.Errorf("Expected Go for vibeguard, got %s", result.Type)
	}
	t.Logf("Detection: Type=%s, Confidence=%.2f, Indicators=%v",
		result.Type, result.Confidence, result.Indicators)

	// Run tool scanner
	scanner := NewToolScanner(dir)
	tools, err := scanner.ScanAll()
	if err != nil {
		t.Fatalf("ScanAll failed: %v", err)
	}
	t.Logf("Detected tools (%d):", len(tools))
	for _, tool := range tools {
		t.Logf("  - [%s] %s (confidence: %.2f, indicators: %v)",
			tool.Category, tool.Name, tool.Confidence, tool.Indicators)
	}

	// Verify expected tools for vibeguard
	// Note: vibeguard doesn't have .golangci.yml config (uses defaults) or .pre-commit-config.yaml
	// but it does have raw git hooks
	expectedTools := map[string]bool{
		"gofmt":          false,
		"go vet":         false,
		"go test":        false,
		"GitHub Actions": false,
		"git hooks":      false,
	}

	for _, tool := range tools {
		if _, ok := expectedTools[tool.Name]; ok {
			expectedTools[tool.Name] = true
		}
	}

	for name, found := range expectedTools {
		if !found {
			t.Errorf("Expected tool not detected: %s", name)
		}
	}

	// Run metadata extractor
	extractor := NewMetadataExtractor(dir)
	metadata, err := extractor.Extract(Go)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}
	if metadata.Name == "" {
		t.Error("Expected module name to be extracted")
	}
	t.Logf("Metadata: Name=%s, Version=%s, Go=%s",
		metadata.Name, metadata.Version, metadata.Extra["go_version"])

	// Run structure extractor
	structure, err := extractor.ExtractStructure(Go)
	if err != nil {
		t.Fatalf("ExtractStructure failed: %v", err)
	}
	t.Logf("Structure:")
	t.Logf("  EntryPoints: %v", structure.EntryPoints)
	t.Logf("  SourceDirs: %v", structure.SourceDirs)
	t.Logf("  TestDirs: %v", structure.TestDirs)
	t.Logf("  ConfigFiles: %v", structure.ConfigFiles)
	t.Logf("  HasMonorepo: %v", structure.HasMonorepo)
	t.Logf("  BuildOutputDir: %s", structure.BuildOutputDir)

	// Verify expected structure
	if len(structure.ConfigFiles) == 0 {
		t.Error("Expected config files to be found")
	}

	// Run recommender
	recommender := NewRecommender(Go, tools)
	recs := recommender.Recommend()
	t.Logf("Recommendations (%d):", len(recs))
	for _, rec := range recs {
		t.Logf("  - [%s] %s (priority: %d): %s",
			rec.Category, rec.ID, rec.Priority, rec.Description)
		if rec.Command != "" {
			t.Logf("      Command: %s", rec.Command)
		}
	}

	// Verify we get reasonable recommendations
	if len(recs) < 3 {
		t.Errorf("Expected at least 3 recommendations for vibeguard, got %d", len(recs))
	}

	// Check for essential Go recommendations
	// Note: Without golangci-lint config, we get "vet" instead of "lint" for linting
	foundVet := false
	foundFmt := false
	foundTest := false
	for _, rec := range recs {
		switch rec.ID {
		case "vet":
			foundVet = true
		case "fmt":
			foundFmt = true
		case "test":
			foundTest = true
		}
	}

	if !foundVet {
		t.Error("Expected vet recommendation")
	}
	if !foundFmt {
		t.Error("Expected fmt recommendation")
	}
	if !foundTest {
		t.Error("Expected test recommendation")
	}
}

// TestInspectorEndToEnd runs a full end-to-end test simulating what init --assist would do.
func TestInspectorEndToEnd(t *testing.T) {
	testCases := []struct {
		name         string
		setup        func(dir string)
		expectedType ProjectType
		minTools     int
		minRecs      int
	}{
		{
			name: "complete_go_project",
			setup: func(dir string) {
				mustWrite(dir, "go.mod", "module example.com/complete\n\ngo 1.21\n")
				mustWrite(dir, "go.sum", "")
				mustWrite(dir, "main.go", "package main\n\nfunc main() {}\n")
				mustMkdir(dir, "cmd/myapp")
				mustWrite(dir, "cmd/myapp/main.go", "package main\n\nfunc main() {}\n")
				mustMkdir(dir, "internal/pkg")
				mustWrite(dir, "internal/pkg/lib.go", "package pkg\n")
				mustWrite(dir, ".golangci.yml", "run:\n  timeout: 5m\n")
				mustWrite(dir, "Makefile", ".PHONY: test\ntest:\n\tgo test ./...\n")
				mustMkdir(dir, ".github/workflows")
				mustWrite(dir, ".github/workflows/ci.yml", "name: CI\non: push\n")
				mustWrite(dir, ".pre-commit-config.yaml", "repos: []\n")
			},
			expectedType: Go,
			minTools:     5,
			minRecs:      3,
		},
		{
			name: "complete_node_project",
			setup: func(dir string) {
				mustWrite(dir, "package.json", `{
					"name": "test-project",
					"version": "1.0.0",
					"devDependencies": {
						"eslint": "^8.0.0",
						"prettier": "^3.0.0",
						"jest": "^29.0.0",
						"typescript": "^5.0.0"
					},
					"scripts": {
						"test": "jest",
						"lint": "eslint .",
						"format": "prettier --check ."
					}
				}`)
				mustWrite(dir, "package-lock.json", "{}")
				mustWrite(dir, "tsconfig.json", `{"compilerOptions": {"strict": true}}`)
				mustWrite(dir, ".eslintrc.json", `{"extends": ["eslint:recommended"]}`)
				mustWrite(dir, ".prettierrc", `{"semi": true}`)
				mustMkdir(dir, "src")
				mustWrite(dir, "src/index.ts", "export const hello = 'world';")
				mustMkdir(dir, "__tests__")
				mustWrite(dir, "__tests__/index.test.ts", "test('hello', () => {});")
				mustMkdir(dir, ".husky")
				mustWrite(dir, ".husky/pre-commit", "#!/bin/sh\nnpm test\n")
			},
			expectedType: Node,
			minTools:     6,
			minRecs:      4,
		},
		{
			name: "complete_python_project",
			setup: func(dir string) {
				mustWrite(dir, "pyproject.toml", `[project]
name = "test-project"
version = "1.0.0"
description = "A test project"
requires-python = ">=3.9"

[tool.ruff]
line-length = 88

[tool.mypy]
strict = true

[tool.pytest.ini_options]
testpaths = ["tests"]

[tool.black]
line-length = 88
`)
				mustMkdir(dir, "src/mypackage")
				mustWrite(dir, "src/mypackage/__init__.py", "")
				mustWrite(dir, "src/mypackage/main.py", "def main():\n    pass\n")
				mustMkdir(dir, "tests")
				mustWrite(dir, "tests/__init__.py", "")
				mustWrite(dir, "tests/test_main.py", "def test_main():\n    pass\n")
				mustMkdir(dir, ".github/workflows")
				mustWrite(dir, ".github/workflows/ci.yml", "name: CI\non: push\n")
			},
			expectedType: Python,
			minTools:     4,
			minRecs:      3,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			tc.setup(dir)

			// Phase 1: Detection
			detector := NewDetector(dir)
			result, err := detector.DetectPrimary()
			if err != nil {
				t.Fatalf("Detection failed: %v", err)
			}
			if result.Type != tc.expectedType {
				t.Errorf("Expected %s, got %s", tc.expectedType, result.Type)
			}

			// Phase 2: Tool Scanning
			scanner := NewToolScanner(dir)
			tools, err := scanner.ScanAll()
			if err != nil {
				t.Fatalf("Tool scanning failed: %v", err)
			}
			if len(tools) < tc.minTools {
				t.Errorf("Expected at least %d tools, got %d: %v",
					tc.minTools, len(tools), toolNames(tools))
			}

			// Phase 3: Metadata Extraction
			extractor := NewMetadataExtractor(dir)
			metadata, err := extractor.Extract(result.Type)
			if err != nil {
				t.Fatalf("Metadata extraction failed: %v", err)
			}
			if metadata.Name == "" {
				t.Error("Expected project name to be extracted")
			}

			// Phase 4: Structure Extraction
			structure, err := extractor.ExtractStructure(result.Type)
			if err != nil {
				t.Fatalf("Structure extraction failed: %v", err)
			}
			if len(structure.ConfigFiles) == 0 {
				t.Error("Expected config files to be found")
			}

			// Phase 5: Recommendations
			recommender := NewRecommender(result.Type, tools)
			recs := recommender.Recommend()
			if len(recs) < tc.minRecs {
				t.Errorf("Expected at least %d recommendations, got %d",
					tc.minRecs, len(recs))
			}

			// Log summary
			t.Logf("Summary for %s:", tc.name)
			t.Logf("  Type: %s (confidence: %.2f)", result.Type, result.Confidence)
			t.Logf("  Name: %s, Version: %s", metadata.Name, metadata.Version)
			t.Logf("  Tools: %v", toolNames(tools))
			t.Logf("  Recommendations: %d", len(recs))
		})
	}
}

// Helper functions

func toolNames(tools []ToolInfo) []string {
	names := make([]string, len(tools))
	for i, tool := range tools {
		names[i] = fmt.Sprintf("%s (%.0f%%)", tool.Name, tool.Confidence*100)
	}
	return names
}

func mustWrite(dir, name, content string) {
	path := filepath.Join(dir, name)
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		panic(err)
	}
}

func mustMkdir(dir, name string) {
	if err := os.MkdirAll(filepath.Join(dir, name), 0755); err != nil {
		panic(err)
	}
}
