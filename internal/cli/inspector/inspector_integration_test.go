// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"os"
	"path/filepath"
	"testing"
)

// Integration tests for the inspector package on diverse project types.
// These tests verify the inspector works correctly on various real-world scenarios.

// TestSimpleSingleToolProjects tests detection on projects with just one tool.
func TestSimpleSingleToolProjects(t *testing.T) {
	t.Run("Go project with only go.mod", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/simple\n\ngo 1.21\n")
		writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Go {
			t.Errorf("Expected Go, got %s", result.Type)
		}
		if result.Confidence < 0.8 {
			t.Errorf("Expected confidence >= 0.8, got %f", result.Confidence)
		}

		scanner := NewToolScanner(dir)
		tools, err := scanner.ScanAll()
		if err != nil {
			t.Fatalf("ScanAll failed: %v", err)
		}

		// Should detect built-in Go tools (gofmt, go vet, go test)
		foundGoFmt := false
		for _, tool := range tools {
			if tool.Name == "gofmt" && tool.Detected {
				foundGoFmt = true
			}
		}
		if !foundGoFmt {
			t.Error("Expected gofmt to be detected for Go project")
		}

		recommender := NewRecommender(Go, tools)
		recs := recommender.Recommend()
		if len(recs) == 0 {
			t.Error("Expected at least one recommendation for simple Go project")
		}
	})

	t.Run("Node project with only package.json", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{"name": "simple-node", "version": "1.0.0"}`)
		writeFile(t, dir, "index.js", "console.log('hello');\n")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Node {
			t.Errorf("Expected Node, got %s", result.Type)
		}

		scanner := NewToolScanner(dir)
		tools, err := scanner.ScanAll()
		if err != nil {
			t.Fatalf("ScanAll failed: %v", err)
		}

		// Should detect npm audit (always available with Node)
		foundNpmAudit := false
		for _, tool := range tools {
			if tool.Name == "npm audit" && tool.Detected {
				foundNpmAudit = true
			}
		}
		if !foundNpmAudit {
			t.Error("Expected npm audit to be detected for Node project")
		}
	})

	t.Run("Python project with only requirements.txt", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "requirements.txt", "flask==2.0.0\nrequests\n")
		writeFile(t, dir, "app.py", "from flask import Flask\n")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Python {
			t.Errorf("Expected Python, got %s", result.Type)
		}
		if result.Confidence < 0.3 {
			t.Errorf("Expected confidence >= 0.3 for requirements.txt, got %f", result.Confidence)
		}
	})
}

// TestComplexMultiToolProjects tests detection on projects with many tools.
func TestComplexMultiToolProjects(t *testing.T) {
	t.Run("Go project with golangci-lint and multiple configs", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/complex\n\ngo 1.21\n")
		writeFile(t, dir, "go.sum", "")
		writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")
		writeFile(t, dir, ".golangci.yml", "run:\n  timeout: 5m\n")
		writeFile(t, dir, "Makefile", "lint:\n\tgoimports -l .\n")
		mkdir(t, dir, ".github/workflows")
		writeFile(t, dir, ".github/workflows/ci.yml", "name: CI\non: push\n")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Go {
			t.Errorf("Expected Go, got %s", result.Type)
		}
		if result.Confidence < 0.9 {
			t.Errorf("Expected high confidence for full Go project, got %f", result.Confidence)
		}

		scanner := NewToolScanner(dir)
		tools, err := scanner.ScanAll()
		if err != nil {
			t.Fatalf("ScanAll failed: %v", err)
		}

		// Verify multiple tools detected
		toolNames := make(map[string]bool)
		for _, tool := range tools {
			if tool.Detected {
				toolNames[tool.Name] = true
			}
		}

		expectedTools := []string{"golangci-lint", "gofmt", "go vet", "go test", "goimports", "GitHub Actions"}
		for _, name := range expectedTools {
			if !toolNames[name] {
				t.Errorf("Expected %s to be detected", name)
			}
		}
	})

	t.Run("Node project with ESLint, Prettier, Jest, TypeScript", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{
			"name": "complex-node",
			"version": "1.0.0",
			"devDependencies": {
				"eslint": "^8.0.0",
				"prettier": "^3.0.0",
				"jest": "^29.0.0",
				"typescript": "^5.0.0"
			}
		}`)
		writeFile(t, dir, "package-lock.json", "{}")
		writeFile(t, dir, "tsconfig.json", `{"compilerOptions": {"strict": true}}`)
		writeFile(t, dir, ".eslintrc.json", `{"extends": ["eslint:recommended"]}`)
		writeFile(t, dir, ".prettierrc", `{"semi": true}`)
		writeFile(t, dir, "jest.config.js", "module.exports = {};")
		writeFile(t, dir, "src/index.ts", "export const hello = 'world';")
		mkdir(t, dir, "src")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Node {
			t.Errorf("Expected Node, got %s", result.Type)
		}

		scanner := NewToolScanner(dir)
		tools, err := scanner.ScanAll()
		if err != nil {
			t.Fatalf("ScanAll failed: %v", err)
		}

		toolNames := make(map[string]bool)
		for _, tool := range tools {
			if tool.Detected {
				toolNames[tool.Name] = true
			}
		}

		expectedTools := []string{"eslint", "prettier", "jest", "typescript", "npm audit"}
		for _, name := range expectedTools {
			if !toolNames[name] {
				t.Errorf("Expected %s to be detected", name)
			}
		}

		// Check recommendations
		recommender := NewRecommender(Node, tools)
		recs := recommender.Recommend()
		if len(recs) < 4 {
			t.Errorf("Expected at least 4 recommendations, got %d", len(recs))
		}
	})

	t.Run("Python project with modern tooling (pyproject.toml, ruff, mypy, pytest)", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pyproject.toml", `[project]
name = "complex-python"
version = "1.0.0"
description = "A complex Python project"
requires-python = ">=3.9"

[tool.ruff]
line-length = 88

[tool.mypy]
strict = true

[tool.pytest.ini_options]
testpaths = ["tests"]
`)
		writeFile(t, dir, "src/__init__.py", "")
		writeFile(t, dir, "tests/__init__.py", "")
		mkdir(t, dir, "src")
		mkdir(t, dir, "tests")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Python {
			t.Errorf("Expected Python, got %s", result.Type)
		}

		scanner := NewToolScanner(dir)
		tools, err := scanner.ScanAll()
		if err != nil {
			t.Fatalf("ScanAll failed: %v", err)
		}

		// Check that ruff and mypy are detected from pyproject.toml
		foundRuff := false
		foundMypy := false
		foundPytest := false
		for _, tool := range tools {
			if tool.Detected {
				switch tool.Name {
				case "ruff":
					foundRuff = true
				case "mypy":
					foundMypy = true
				case "pytest":
					foundPytest = true
				}
			}
		}
		if !foundRuff {
			t.Error("Expected ruff to be detected from pyproject.toml")
		}
		if !foundMypy {
			t.Error("Expected mypy to be detected from pyproject.toml")
		}
		if !foundPytest {
			t.Error("Expected pytest to be detected from pyproject.toml")
		}

		// Verify metadata extraction
		extractor := NewMetadataExtractor(dir)
		metadata, err := extractor.Extract(Python)
		if err != nil {
			t.Fatalf("Extract failed: %v", err)
		}
		if metadata.Name != "complex-python" {
			t.Errorf("Expected name 'complex-python', got %s", metadata.Name)
		}
		if metadata.Version != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got %s", metadata.Version)
		}
	})
}

// TestMinimalEdgeCaseProjects tests edge cases and minimal projects.
func TestMinimalEdgeCaseProjects(t *testing.T) {
	t.Run("Empty directory", func(t *testing.T) {
		dir := t.TempDir()

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Unknown {
			t.Errorf("Expected Unknown for empty dir, got %s", result.Type)
		}
		if result.Confidence != 0 {
			t.Errorf("Expected confidence 0 for empty dir, got %f", result.Confidence)
		}

		scanner := NewToolScanner(dir)
		tools, err := scanner.ScanAll()
		if err != nil {
			t.Fatalf("ScanAll failed: %v", err)
		}
		if len(tools) != 0 {
			t.Errorf("Expected no tools for empty dir, got %d", len(tools))
		}
	})

	t.Run("Mixed language indicators (Go + Python)", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/mixed\n\ngo 1.21\n")
		writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")
		writeFile(t, dir, "requirements.txt", "requests\n")
		writeFile(t, dir, "script.py", "import requests\n")

		detector := NewDetector(dir)
		results, err := detector.Detect()
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}

		// Should detect both languages
		foundGo := false
		foundPython := false
		for _, r := range results {
			if r.Type == Go && r.Confidence > 0 {
				foundGo = true
			}
			if r.Type == Python && r.Confidence > 0 {
				foundPython = true
			}
		}
		if !foundGo {
			t.Error("Expected Go to be detected in mixed project")
		}
		if !foundPython {
			t.Error("Expected Python to be detected in mixed project")
		}

		// Go should have higher confidence due to go.mod being stronger than requirements.txt
		primary, _ := detector.DetectPrimary()
		if primary.Type != Go {
			t.Errorf("Expected Go as primary for mixed project (stronger indicators), got %s", primary.Type)
		}
	})

	t.Run("Only source files, no config files", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")
		// No go.mod!

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		// Should still detect Go from .go files but with lower confidence
		if result.Type != Go {
			t.Errorf("Expected Go from .go files, got %s", result.Type)
		}
		if result.Confidence >= 0.5 {
			t.Errorf("Expected lower confidence without go.mod, got %f", result.Confidence)
		}
	})

	t.Run("Project with only .gitignore and README", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, ".gitignore", "*.log\nnode_modules/\n")
		writeFile(t, dir, "README.md", "# My Project\n")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Unknown {
			t.Errorf("Expected Unknown for non-code project, got %s", result.Type)
		}
	})
}

// TestUnusualProjectStructures tests projects with non-standard layouts.
func TestUnusualProjectStructures(t *testing.T) {
	t.Run("Monorepo with multiple packages", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "package.json", `{
			"name": "monorepo-root",
			"workspaces": ["packages/*"]
		}`)
		mkdir(t, dir, "packages/pkg-a")
		mkdir(t, dir, "packages/pkg-b")
		writeFile(t, dir, "packages/pkg-a/package.json", `{"name": "pkg-a"}`)
		writeFile(t, dir, "packages/pkg-b/package.json", `{"name": "pkg-b"}`)

		// Verify monorepo detection
		extractor := NewMetadataExtractor(dir)
		structure, err := extractor.ExtractStructure(Node)
		if err != nil {
			t.Fatalf("ExtractStructure failed: %v", err)
		}
		if !structure.HasMonorepo {
			t.Error("Expected monorepo to be detected")
		}
	})

	t.Run("Go monorepo with workspace", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.work", "go 1.21\n\nuse (\n\t./service-a\n\t./service-b\n)\n")
		mkdir(t, dir, "service-a")
		mkdir(t, dir, "service-b")
		writeFile(t, dir, "service-a/go.mod", "module example.com/service-a\n\ngo 1.21\n")
		writeFile(t, dir, "service-b/go.mod", "module example.com/service-b\n\ngo 1.21\n")
		writeFile(t, dir, "service-a/main.go", "package main\n\nfunc main() {}\n")
		writeFile(t, dir, "service-b/main.go", "package main\n\nfunc main() {}\n")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		// Should detect Go from the go.work or sub-modules
		if result.Type != Go {
			t.Errorf("Expected Go for Go workspace, got %s", result.Type)
		}
	})

	t.Run("Java Maven standard layout", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "pom.xml", `<?xml version="1.0"?>
<project>
  <groupId>com.example</groupId>
  <artifactId>my-app</artifactId>
  <version>1.0.0</version>
  <name>My App</name>
</project>`)
		mkdir(t, dir, "src/main/java/com/example")
		mkdir(t, dir, "src/test/java/com/example")
		writeFile(t, dir, "src/main/java/com/example/App.java", "package com.example;\npublic class App {}\n")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Java {
			t.Errorf("Expected Java, got %s", result.Type)
		}

		extractor := NewMetadataExtractor(dir)
		structure, err := extractor.ExtractStructure(Java)
		if err != nil {
			t.Fatalf("ExtractStructure failed: %v", err)
		}

		// Check Maven standard layout directories
		foundMainJava := false
		for _, d := range structure.SourceDirs {
			if d == "src/main/java" {
				foundMainJava = true
			}
		}
		if !foundMainJava {
			t.Error("Expected src/main/java in SourceDirs")
		}

		foundTestJava := false
		for _, d := range structure.TestDirs {
			if d == "src/test/java" {
				foundTestJava = true
			}
		}
		if !foundTestJava {
			t.Error("Expected src/test/java in TestDirs")
		}

		// Check metadata extraction
		metadata, err := extractor.Extract(Java)
		if err != nil {
			t.Fatalf("Extract failed: %v", err)
		}
		if metadata.Name != "My App" {
			t.Errorf("Expected name 'My App', got %s", metadata.Name)
		}
		if metadata.Version != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got %s", metadata.Version)
		}
	})

	t.Run("Rust workspace", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "Cargo.toml", `[workspace]
members = ["crates/*"]

[package]
name = "workspace-root"
version = "0.1.0"
edition = "2021"
`)
		mkdir(t, dir, "crates/lib-a")
		mkdir(t, dir, "crates/lib-b")
		writeFile(t, dir, "crates/lib-a/Cargo.toml", `[package]
name = "lib-a"
version = "0.1.0"
edition = "2021"
`)
		writeFile(t, dir, "crates/lib-b/Cargo.toml", `[package]
name = "lib-b"
version = "0.1.0"
edition = "2021"
`)
		writeFile(t, dir, "crates/lib-a/src/lib.rs", "pub fn hello() {}\n")
		mkdir(t, dir, "crates/lib-a/src")
		mkdir(t, dir, "crates/lib-b/src")

		detector := NewDetector(dir)
		result, err := detector.DetectPrimary()
		if err != nil {
			t.Fatalf("DetectPrimary failed: %v", err)
		}
		if result.Type != Rust {
			t.Errorf("Expected Rust, got %s", result.Type)
		}

		// Check for Cargo workspace detection
		extractor := NewMetadataExtractor(dir)
		structure, err := extractor.ExtractStructure(Rust)
		if err != nil {
			t.Fatalf("ExtractStructure failed: %v", err)
		}
		if !structure.HasMonorepo {
			t.Error("Expected Rust workspace to be detected as monorepo")
		}
	})

	t.Run("Pre-commit hooks configuration", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module example.com/hooks\n\ngo 1.21\n")
		writeFile(t, dir, ".pre-commit-config.yaml", `repos:
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.55.0
    hooks:
      - id: golangci-lint
`)

		scanner := NewToolScanner(dir)
		tools, err := scanner.ScanAll()
		if err != nil {
			t.Fatalf("ScanAll failed: %v", err)
		}

		foundPrecommit := false
		for _, tool := range tools {
			if tool.Name == "pre-commit" && tool.Detected {
				foundPrecommit = true
			}
		}
		if !foundPrecommit {
			t.Error("Expected pre-commit to be detected")
		}
	})
}

// TestRecommendationDeduplication verifies that duplicate recommendations are properly handled.
func TestRecommendationDeduplication(t *testing.T) {
	t.Run("Deduplicate by ID", func(t *testing.T) {
		recs := []CheckRecommendation{
			{ID: "lint", Tool: "eslint", Priority: 20},
			{ID: "lint", Tool: "tslint", Priority: 25}, // duplicate ID
			{ID: "fmt", Tool: "prettier", Priority: 10},
		}

		deduped := DeduplicateRecommendations(recs)
		if len(deduped) != 2 {
			t.Errorf("Expected 2 unique recommendations, got %d", len(deduped))
		}

		// First occurrence should be kept
		if deduped[0].ID != "lint" || deduped[0].Tool != "eslint" {
			t.Error("Expected first 'lint' rec to be eslint")
		}
	})
}

// TestGroupByCategoryIntegration verifies category grouping works correctly with full recommendations.
func TestGroupByCategoryIntegration(t *testing.T) {
	recs := []CheckRecommendation{
		{ID: "lint", Category: "lint"},
		{ID: "fmt", Category: "format"},
		{ID: "vet", Category: "lint"},
		{ID: "test", Category: "test"},
	}

	groups := GroupByCategory(recs)
	if len(groups["lint"]) != 2 {
		t.Errorf("Expected 2 lint recommendations, got %d", len(groups["lint"]))
	}
	if len(groups["format"]) != 1 {
		t.Errorf("Expected 1 format recommendation, got %d", len(groups["format"]))
	}
	if len(groups["test"]) != 1 {
		t.Errorf("Expected 1 test recommendation, got %d", len(groups["test"]))
	}
}

// Helper functions

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	// Ensure parent directory exists
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", parent, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", name, err)
	}
}

func mkdir(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", name, err)
	}
}
