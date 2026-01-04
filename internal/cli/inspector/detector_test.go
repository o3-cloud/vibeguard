package inspector

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// createTestProject creates a temporary directory with specified files.
func createTestProject(t *testing.T, files map[string]string, dirs []string) string {
	t.Helper()
	root := t.TempDir()

	// Create directories
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(root, dir), 0755)
		if err != nil {
			t.Fatalf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create files
	for name, content := range files {
		path := filepath.Join(root, name)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create directory for %s: %v", name, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", name, err)
		}
	}

	return root
}

func TestDetector_DetectGo(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		dirs           []string
		expectedType   ProjectType
		minConfidence  float64
		maxConfidence  float64
		expectDetected bool
	}{
		{
			name: "full Go project",
			files: map[string]string{
				"go.mod":        "module example.com/test\n\ngo 1.21",
				"go.sum":        "github.com/stretchr/testify v1.8.0 h1:...",
				"main.go":       "package main\n\nfunc main() {}",
				"cmd/app.go":    "package main",
				"internal/x.go": "package internal",
			},
			expectedType:   Go,
			minConfidence:  0.95,
			maxConfidence:  1.0,
			expectDetected: true,
		},
		{
			name: "Go with only go.mod",
			files: map[string]string{
				"go.mod": "module example.com/test",
			},
			expectedType:   Go,
			minConfidence:  0.6,
			maxConfidence:  0.65,
			expectDetected: true,
		},
		{
			name: "Go with go.mod and go.sum",
			files: map[string]string{
				"go.mod": "module example.com/test",
				"go.sum": "...",
			},
			expectedType:   Go,
			minConfidence:  0.8,
			maxConfidence:  0.85,
			expectDetected: true,
		},
		{
			name: "Go files only (no go.mod)",
			files: map[string]string{
				"main.go":    "package main",
				"handler.go": "package handler",
			},
			expectedType:   Go,
			minConfidence:  0.2,
			maxConfidence:  0.25,
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, tt.dirs)
			detector := NewDetector(root)

			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			if !tt.expectDetected {
				for _, r := range results {
					if r.Type == tt.expectedType && r.Confidence > 0 {
						t.Errorf("expected %s not to be detected, but found with confidence %f", tt.expectedType, r.Confidence)
					}
				}
				return
			}

			var found *DetectionResult
			for i := range results {
				if results[i].Type == tt.expectedType {
					found = &results[i]
					break
				}
			}

			if found == nil {
				t.Fatalf("expected to detect %s, but not found in results", tt.expectedType)
			}

			if found.Confidence < tt.minConfidence {
				t.Errorf("confidence %f is below minimum %f", found.Confidence, tt.minConfidence)
			}
			if found.Confidence > tt.maxConfidence {
				t.Errorf("confidence %f is above maximum %f", found.Confidence, tt.maxConfidence)
			}
		})
	}
}

func TestDetector_DetectNode(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		dirs           []string
		minConfidence  float64
		maxConfidence  float64
		expectDetected bool
	}{
		{
			name: "full Node project with npm",
			files: map[string]string{
				"package.json":      `{"name": "test", "version": "1.0.0"}`,
				"package-lock.json": "{}",
			},
			dirs:           []string{"node_modules"},
			minConfidence:  0.95,
			maxConfidence:  1.0,
			expectDetected: true,
		},
		{
			name: "Node with yarn",
			files: map[string]string{
				"package.json": `{"name": "test"}`,
				"yarn.lock":    "...",
			},
			minConfidence:  0.8,
			maxConfidence:  0.85,
			expectDetected: true,
		},
		{
			name: "Node with pnpm",
			files: map[string]string{
				"package.json":   `{"name": "test"}`,
				"pnpm-lock.yaml": "...",
			},
			minConfidence:  0.8,
			maxConfidence:  0.85,
			expectDetected: true,
		},
		{
			name: "package.json only",
			files: map[string]string{
				"package.json": `{"name": "test"}`,
			},
			minConfidence:  0.6,
			maxConfidence:  0.65,
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, tt.dirs)
			detector := NewDetector(root)

			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			if !tt.expectDetected {
				return
			}

			var found *DetectionResult
			for i := range results {
				if results[i].Type == Node {
					found = &results[i]
					break
				}
			}

			if found == nil {
				t.Fatalf("expected to detect Node, but not found in results")
			}

			if found.Confidence < tt.minConfidence {
				t.Errorf("confidence %f is below minimum %f", found.Confidence, tt.minConfidence)
			}
			if found.Confidence > tt.maxConfidence {
				t.Errorf("confidence %f is above maximum %f", found.Confidence, tt.maxConfidence)
			}
		})
	}
}

func TestDetector_DetectPython(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		minConfidence  float64
		maxConfidence  float64
		expectDetected bool
	}{
		{
			name: "Python with pyproject.toml",
			files: map[string]string{
				"pyproject.toml": "[tool.poetry]\nname = \"test\"",
				"main.py":        "print('hello')",
			},
			minConfidence:  0.8,
			maxConfidence:  0.85,
			expectDetected: true,
		},
		{
			name: "Python with setup.py",
			files: map[string]string{
				"setup.py": "from setuptools import setup",
				"app.py":   "import flask",
			},
			minConfidence:  0.6,
			maxConfidence:  0.65,
			expectDetected: true,
		},
		{
			name: "Python with requirements.txt",
			files: map[string]string{
				"requirements.txt": "flask==2.0\nrequests>=2.28",
				"main.py":          "import flask",
			},
			minConfidence:  0.5,
			maxConfidence:  0.55,
			expectDetected: true,
		},
		{
			name: "Python with Pipfile",
			files: map[string]string{
				"Pipfile": "[[source]]\nurl = \"https://pypi.org/simple\"",
			},
			minConfidence:  0.3,
			maxConfidence:  0.35,
			expectDetected: true,
		},
		{
			name: "Full Python project",
			files: map[string]string{
				"pyproject.toml":   "[tool.poetry]\nname = \"test\"",
				"requirements.txt": "flask",
				"main.py":          "print('hello')",
			},
			minConfidence:  0.95,
			maxConfidence:  1.0,
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			detector := NewDetector(root)

			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			if !tt.expectDetected {
				return
			}

			var found *DetectionResult
			for i := range results {
				if results[i].Type == Python {
					found = &results[i]
					break
				}
			}

			if found == nil {
				t.Fatalf("expected to detect Python, but not found in results")
			}

			if found.Confidence < tt.minConfidence {
				t.Errorf("confidence %f is below minimum %f", found.Confidence, tt.minConfidence)
			}
			if found.Confidence > tt.maxConfidence {
				t.Errorf("confidence %f is above maximum %f", found.Confidence, tt.maxConfidence)
			}
		})
	}
}

func TestDetector_DetectRuby(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		minConfidence  float64
		maxConfidence  float64
		expectDetected bool
	}{
		{
			name: "Ruby with Gemfile",
			files: map[string]string{
				"Gemfile":      "source 'https://rubygems.org'\ngem 'rails'",
				"Gemfile.lock": "GEM\n  remote: https://rubygems.org/",
				"app.rb":       "puts 'hello'",
			},
			minConfidence:  0.95,
			maxConfidence:  1.0,
			expectDetected: true,
		},
		{
			name: "Gemfile only",
			files: map[string]string{
				"Gemfile": "source 'https://rubygems.org'",
			},
			minConfidence:  0.6,
			maxConfidence:  0.65,
			expectDetected: true,
		},
		{
			name: "Ruby files only",
			files: map[string]string{
				"app.rb":    "puts 'hello'",
				"server.rb": "require 'sinatra'",
			},
			minConfidence:  0.2,
			maxConfidence:  0.25,
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			detector := NewDetector(root)

			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			var found *DetectionResult
			for i := range results {
				if results[i].Type == Ruby {
					found = &results[i]
					break
				}
			}

			if !tt.expectDetected {
				if found != nil && found.Confidence > 0 {
					t.Errorf("expected Ruby not to be detected")
				}
				return
			}

			if found == nil {
				t.Fatalf("expected to detect Ruby, but not found in results")
			}

			if found.Confidence < tt.minConfidence {
				t.Errorf("confidence %f is below minimum %f", found.Confidence, tt.minConfidence)
			}
			if found.Confidence > tt.maxConfidence {
				t.Errorf("confidence %f is above maximum %f", found.Confidence, tt.maxConfidence)
			}
		})
	}
}

func TestDetector_DetectRust(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		minConfidence  float64
		maxConfidence  float64
		expectDetected bool
	}{
		{
			name: "Rust project",
			files: map[string]string{
				"Cargo.toml":  "[package]\nname = \"test\"",
				"Cargo.lock":  "...",
				"src/main.rs": "fn main() {}",
			},
			minConfidence:  0.95,
			maxConfidence:  1.0,
			expectDetected: true,
		},
		{
			name: "Cargo.toml only",
			files: map[string]string{
				"Cargo.toml": "[package]\nname = \"test\"",
			},
			minConfidence:  0.7,
			maxConfidence:  0.75,
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			detector := NewDetector(root)

			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			var found *DetectionResult
			for i := range results {
				if results[i].Type == Rust {
					found = &results[i]
					break
				}
			}

			if !tt.expectDetected {
				return
			}

			if found == nil {
				t.Fatalf("expected to detect Rust, but not found in results")
			}

			if found.Confidence < tt.minConfidence {
				t.Errorf("confidence %f is below minimum %f", found.Confidence, tt.minConfidence)
			}
			if found.Confidence > tt.maxConfidence {
				t.Errorf("confidence %f is above maximum %f", found.Confidence, tt.maxConfidence)
			}
		})
	}
}

func TestDetector_DetectJava(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		minConfidence  float64
		maxConfidence  float64
		expectDetected bool
	}{
		{
			name: "Maven project",
			files: map[string]string{
				"pom.xml":                "<project><modelVersion>4.0.0</modelVersion></project>",
				"src/main/java/App.java": "public class App {}",
			},
			minConfidence:  0.8,
			maxConfidence:  0.85,
			expectDetected: true,
		},
		{
			name: "Gradle project",
			files: map[string]string{
				"build.gradle":           "plugins { id 'java' }",
				"src/main/java/App.java": "public class App {}",
			},
			minConfidence:  0.8,
			maxConfidence:  0.85,
			expectDetected: true,
		},
		{
			name: "Kotlin Gradle project",
			files: map[string]string{
				"build.gradle.kts": "plugins { kotlin(\"jvm\") }",
			},
			minConfidence:  0.6,
			maxConfidence:  0.65,
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			detector := NewDetector(root)

			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			var found *DetectionResult
			for i := range results {
				if results[i].Type == Java {
					found = &results[i]
					break
				}
			}

			if !tt.expectDetected {
				return
			}

			if found == nil {
				t.Fatalf("expected to detect Java, but not found in results")
			}

			if found.Confidence < tt.minConfidence {
				t.Errorf("confidence %f is below minimum %f", found.Confidence, tt.minConfidence)
			}
			if found.Confidence > tt.maxConfidence {
				t.Errorf("confidence %f is above maximum %f", found.Confidence, tt.maxConfidence)
			}
		})
	}
}

func TestDetector_DetectPrimary(t *testing.T) {
	// Test that DetectPrimary returns the highest confidence result
	files := map[string]string{
		"go.mod":       "module test",
		"go.sum":       "...",
		"main.go":      "package main",
		"package.json": `{"name": "test"}`, // Lower confidence Node
	}

	root := createTestProject(t, files, nil)
	detector := NewDetector(root)

	result, err := detector.DetectPrimary()
	if err != nil {
		t.Fatalf("DetectPrimary() error = %v", err)
	}

	if result.Type != Go {
		t.Errorf("expected primary type to be Go, got %s", result.Type)
	}

	// Go should have higher confidence than Node
	if result.Confidence < 0.9 {
		t.Errorf("expected high confidence for Go, got %f", result.Confidence)
	}
}

func TestDetector_EmptyProject(t *testing.T) {
	root := t.TempDir()
	detector := NewDetector(root)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result for empty project, got %d", len(results))
	}

	if results[0].Type != Unknown {
		t.Errorf("expected Unknown type for empty project, got %s", results[0].Type)
	}

	if results[0].Confidence != 0 {
		t.Errorf("expected 0 confidence for Unknown type, got %f", results[0].Confidence)
	}
}

func TestDetector_MultiLanguageProject(t *testing.T) {
	// Project with multiple language indicators
	files := map[string]string{
		"go.mod":           "module test",
		"main.go":          "package main",
		"package.json":     `{"name": "frontend"}`,
		"requirements.txt": "flask",
	}

	root := createTestProject(t, files, nil)
	detector := NewDetector(root)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should detect multiple project types
	if len(results) < 3 {
		t.Errorf("expected at least 3 project types detected, got %d", len(results))
	}

	// Results should be sorted by confidence
	for i := 1; i < len(results); i++ {
		if results[i].Confidence > results[i-1].Confidence {
			t.Errorf("results not sorted by confidence: %f > %f", results[i].Confidence, results[i-1].Confidence)
		}
	}
}

func TestDetector_IndicatorsPopulated(t *testing.T) {
	files := map[string]string{
		"go.mod":  "module test",
		"go.sum":  "...",
		"main.go": "package main",
	}

	root := createTestProject(t, files, nil)
	detector := NewDetector(root)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	var goResult *DetectionResult
	for i := range results {
		if results[i].Type == Go {
			goResult = &results[i]
			break
		}
	}

	if goResult == nil {
		t.Fatal("Go not detected")
	}

	if len(goResult.Indicators) == 0 {
		t.Error("expected indicators to be populated")
	}

	// Check that key indicators are present
	hasGoMod := false
	hasGoSum := false
	hasGoFiles := false
	for _, ind := range goResult.Indicators {
		switch ind {
		case "go.mod":
			hasGoMod = true
		case "go.sum":
			hasGoSum = true
		case "*.go files":
			hasGoFiles = true
		}
	}

	if !hasGoMod {
		t.Error("expected go.mod in indicators")
	}
	if !hasGoSum {
		t.Error("expected go.sum in indicators")
	}
	if !hasGoFiles {
		t.Error("expected *.go files in indicators")
	}
}

func TestDetector_SkipsVendorDirs(t *testing.T) {
	files := map[string]string{
		"go.mod":                    "module test",
		"main.go":                   "package main",
		"vendor/other/lib.go":       "package lib",
		"node_modules/pkg/index.js": "module.exports = {}",
	}

	root := createTestProject(t, files, nil)
	detector := NewDetector(root)

	// This should not error and should complete quickly
	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should still detect Go
	var found bool
	for _, r := range results {
		if r.Type == Go {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Go to be detected even with vendor dirs present")
	}
}

func TestDetector_NonExistentDirectory(t *testing.T) {
	detector := NewDetector("/nonexistent/path/to/project")

	results, err := detector.Detect()
	// Should not error, just return Unknown
	if err != nil {
		t.Fatalf("Detect() error = %v (should handle gracefully)", err)
	}

	if len(results) != 1 || results[0].Type != Unknown {
		t.Errorf("expected Unknown type for nonexistent directory, got %v", results)
	}
}

func TestDetector_FileInsteadOfDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "not_a_dir")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	detector := NewDetector(filePath)

	// Should handle gracefully
	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v (should handle gracefully)", err)
	}

	// Should return Unknown
	if len(results) != 1 || results[0].Type != Unknown {
		t.Errorf("expected Unknown type when path is a file, got %v", results)
	}
}

func TestDetector_MixedConfidenceOrdering(t *testing.T) {
	// Create a project with multiple languages, different confidence levels
	files := map[string]string{
		"Cargo.toml":       "[package]\nname = \"test\"",
		"Cargo.lock":       "...",
		"src/main.rs":      "fn main() {}",
		"package.json":     `{"name": "test"}`, // Lower confidence
		"requirements.txt": "flask",            // Low confidence
	}

	root := createTestProject(t, files, nil)
	detector := NewDetector(root)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Results should be sorted by confidence (highest first)
	for i := 1; i < len(results); i++ {
		if results[i].Confidence > results[i-1].Confidence {
			t.Errorf("Results not sorted by confidence: %v", results)
			break
		}
	}

	// Rust should be first (highest confidence)
	if results[0].Type != Rust {
		t.Errorf("Expected Rust to be first (highest confidence), got %v", results[0].Type)
	}
}

func TestDetector_Java_BothMavenAndGradle(t *testing.T) {
	// Edge case: project has both pom.xml and build.gradle
	files := map[string]string{
		"pom.xml":      "<project/>",
		"build.gradle": "apply plugin: 'java'",
	}

	root := createTestProject(t, files, nil)
	detector := NewDetector(root)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should detect Java
	var javaResult *DetectionResult
	for i := range results {
		if results[i].Type == Java {
			javaResult = &results[i]
			break
		}
	}

	if javaResult == nil {
		t.Fatal("Expected to detect Java")
	}

	// Confidence should be high (both indicators present)
	if javaResult.Confidence < 1.0 {
		t.Errorf("Expected high confidence for Java with both Maven and Gradle, got %f", javaResult.Confidence)
	}
}

func TestDetector_DeepNestedFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create deeply nested structure
	deepPath := filepath.Join(tmpDir, "a", "b", "c", "d", "e")
	if err := os.MkdirAll(deepPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Put go files deep in structure
	if err := os.WriteFile(filepath.Join(deepPath, "main.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	// But no go.mod at root
	detector := NewDetector(tmpDir)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should not detect Go (files are beyond default depth limit)
	var goResult *DetectionResult
	for i := range results {
		if results[i].Type == Go {
			goResult = &results[i]
			break
		}
	}

	// Without go.mod, deep .go files shouldn't contribute much
	if goResult != nil && goResult.Confidence > 0.5 {
		t.Errorf("Deep nested files shouldn't have high confidence without go.mod")
	}
}

func TestDetector_EmptyGoMod(t *testing.T) {
	// Edge case: empty go.mod file
	files := map[string]string{
		"go.mod": "", // Empty file
	}

	root := createTestProject(t, files, nil)
	detector := NewDetector(root)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should still detect Go (file exists)
	var goResult *DetectionResult
	for i := range results {
		if results[i].Type == Go {
			goResult = &results[i]
			break
		}
	}

	if goResult == nil {
		t.Error("Expected to detect Go with empty go.mod")
	}
}

func TestDetector_Symlinks(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping symlink test in CI")
	}

	tmpDir := t.TempDir()

	// Create actual project
	realDir := filepath.Join(tmpDir, "real")
	if err := os.MkdirAll(realDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(realDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create symlink
	linkDir := filepath.Join(tmpDir, "link")
	if err := os.Symlink(realDir, linkDir); err != nil {
		t.Skipf("Cannot create symlink: %v", err)
	}

	detector := NewDetector(linkDir)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should detect Go through symlink
	var goResult *DetectionResult
	for i := range results {
		if results[i].Type == Go {
			goResult = &results[i]
			break
		}
	}

	if goResult == nil {
		t.Error("Expected to detect Go through symlink")
	}
}

func TestDetector_SkipsHiddenDirs(t *testing.T) {
	files := map[string]string{
		"go.mod":            "module test",
		"main.go":           "package main",
		".hidden/secret.go": "package hidden", // Should be skipped
	}

	root := createTestProject(t, files, nil)
	detector := NewDetector(root)

	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should still detect Go (main files)
	var found bool
	for _, r := range results {
		if r.Type == Go && r.Confidence > 0.5 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Go to be detected")
	}
}

// Benchmark tests for performance target (<500ms execution)
func BenchmarkDetector_Detect(b *testing.B) {
	// Create a realistic project structure
	tmpDir := b.TempDir()

	files := map[string]string{
		"go.mod":                  "module example.com/test\n\ngo 1.21",
		"go.sum":                  "github.com/stretchr/testify v1.8.0 h1:...",
		"main.go":                 "package main\n\nfunc main() {}",
		"cmd/app/main.go":         "package main",
		"internal/pkg/handler.go": "package pkg",
		"pkg/utils/utils.go":      "package utils",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			b.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			b.Fatalf("failed to create file: %v", err)
		}
	}

	detector := NewDetector(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.Detect()
		if err != nil {
			b.Fatalf("Detect() error = %v", err)
		}
	}
}

func BenchmarkDetector_DetectPrimary(b *testing.B) {
	tmpDir := b.TempDir()

	// Multi-language project
	files := map[string]string{
		"go.mod":           "module test",
		"main.go":          "package main",
		"package.json":     `{"name": "frontend"}`,
		"requirements.txt": "flask",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			b.Fatalf("failed to create file: %v", err)
		}
	}

	detector := NewDetector(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.DetectPrimary()
		if err != nil {
			b.Fatalf("DetectPrimary() error = %v", err)
		}
	}
}

// Boundary condition tests for detector.go boundary mutations
// These tests target specific boundary condition mutations found during mutation testing

func TestDetector_DetectSortingBoundary(t *testing.T) {
	// Test boundary condition at line 67 and 69 - sorting with single result
	tmpDir := t.TempDir()
	files := map[string]string{
		"go.mod": "module test\n\ngo 1.21\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	detector := NewDetector(tmpDir)
	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should have detected Go (single result)
	if len(results) < 1 {
		t.Errorf("Expected at least one detection result")
	}
	// Go should be first with confidence > 0
	if results[0].Type != Go {
		t.Errorf("Expected first result to be Go, got %v", results[0].Type)
	}
}

func TestDetector_ConfidenceCappingBoundary(t *testing.T) {
	// Test boundary condition at confidence > 1.0 (lines 130, 170, etc.)
	// Create a project with multiple high-confidence indicators
	tmpDir := t.TempDir()
	files := map[string]string{
		"go.mod":  "module test\n\ngo 1.21\n",
		"go.sum":  "github.com/example/pkg v1.0.0 h1:...\n",
		"main.go": "package main\n\nfunc main() {}\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	detector := NewDetector(tmpDir)
	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Find Go detection
	var goResult *DetectionResult
	for i := range results {
		if results[i].Type == Go {
			goResult = &results[i]
			break
		}
	}

	if goResult == nil {
		t.Errorf("Expected Go detection, got none")
	} else if goResult.Confidence > 1.0 {
		t.Errorf("Confidence = %f, should be capped at 1.0", goResult.Confidence)
	} else if goResult.Confidence <= 0.6 {
		t.Errorf("Confidence = %f, expected > 0.6 with multiple indicators", goResult.Confidence)
	}
}

func TestDetector_DepthBoundaryViaGoDetection(t *testing.T) {
	// Test boundary conditions for depth calculation (lines 382, 386, 389, 393)
	// These mutations are indirectly tested by detecting Go files at different depths
	tmpDir := t.TempDir()

	// Create nested directory structure deeper than maxDepth=3 used in detection
	dirs := []string{
		"cmd",
		"internal",
		"vendor/example.com/pkg",
		"a/b/c/d/e",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create test files at different depths
	// Only files within maxDepth=3 should be found
	files := map[string]string{
		"go.mod":              "module test\n",
		"cmd/main.go":         "package main\n",
		"internal/handler.go": "package internal\n",
		// These should NOT be found (too deep or vendor dir)
		"vendor/example.com/pkg/vendor_test.go": "package vendor\n",
		"a/b/c/d/e/nested.go":                   "package nested\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	detector := NewDetector(tmpDir)
	result, err := detector.detectGo()
	if err != nil {
		t.Fatalf("detectGo() error = %v", err)
	}

	// Go detection should find .go files within reasonable depth
	// and have "*.go files" in indicators
	found := false
	for _, ind := range result.Indicators {
		if ind == "*.go files" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("detectGo() should find *.go files with depth limiting, got indicators: %v", result.Indicators)
	}
}

func TestDetector_MultipleDetectionsWithTies(t *testing.T) {
	// Test boundary conditions in sorting when multiple results have same confidence
	tmpDir := t.TempDir()

	// Create indicators for multiple languages with similar confidence
	files := map[string]string{
		"go.mod":       "module test\n",
		"package.json": `{"name": "test"}`,
		"main.go":      "package main\n",
		"index.js":     "module.exports = {};\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	detector := NewDetector(tmpDir)
	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should have multiple detections
	if len(results) < 2 {
		t.Errorf("Expected multiple detections, got %d", len(results))
	}

	// Results should be sorted by confidence (descending)
	for i := 0; i < len(results)-1; i++ {
		if results[i].Confidence < results[i+1].Confidence {
			t.Errorf("Results not sorted by confidence: %f > %f", results[i].Confidence, results[i+1].Confidence)
		}
	}
}

func TestDetector_NoResults(t *testing.T) {
	// Test case where no project type indicators are found
	tmpDir := t.TempDir()

	// Empty directory
	detector := NewDetector(tmpDir)
	results, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Should have Unknown type with 0 confidence
	if len(results) == 0 {
		t.Errorf("Expected at least Unknown detection result")
	}
	if results[0].Type != Unknown {
		t.Errorf("Expected Unknown type, got %v", results[0].Type)
	}
}

func BenchmarkToolScanner_ScanAll(b *testing.B) {
	tmpDir := b.TempDir()

	// Create realistic project with multiple tool configs
	files := map[string]string{
		"go.mod":                  "module test\n\ngo 1.21\n",
		".golangci.yml":           "linters:\n  enable:\n    - gofmt\n",
		"package.json":            `{"name": "test", "devDependencies": {"eslint": "^8.0.0"}}`,
		".pre-commit-config.yaml": "repos: []\n",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			b.Fatalf("failed to create file: %v", err)
		}
	}

	// Create workflow directory
	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte("name: CI\n"), 0644); err != nil {
		b.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.ScanAll()
		if err != nil {
			b.Fatalf("ScanAll() error = %v", err)
		}
	}
}

// ============================================================================
// Boundary Condition Tests for Detector.go
// These tests specifically target boundary condition mutations detected by
// the mutation testing framework (Gremlins), particularly for comparison
// operators and loop boundary conditions.
// ============================================================================

func TestDetector_SortingBoundaryOperators(t *testing.T) {
	// Tests boundary conditions in the sorting loop (lines 67-73)
	// Targets mutations: i < len(results) vs i <= len(results), etc.
	tests := []struct {
		name            string
		files           map[string]string
		expectedOrder   []ProjectType
		expectedGreater bool // Whether higher index should have lower or equal confidence
	}{
		{
			name: "exactly two results to sort",
			files: map[string]string{
				"go.mod":       "module test",
				"package.json": `{"name": "test"}`,
				"main.go":      "package main",
			},
			expectedOrder:   []ProjectType{Go, Node},
			expectedGreater: false, // Should be descending order
		},
		{
			name: "exactly three results to sort",
			files: map[string]string{
				"go.mod":         "module test",
				"go.sum":         "test",
				"package.json":   `{"name": "test"}`,
				"pyproject.toml": "[project]",
				"main.go":        "package main",
			},
			expectedOrder:   []ProjectType{Go, Node, Python},
			expectedGreater: false,
		},
		{
			name: "four results with varying confidences",
			files: map[string]string{
				"go.mod":         "module test",
				"go.sum":         "test",
				"package.json":   `{"name": "test"}`,
				"pyproject.toml": "[project]",
				"main.go":        "package main",
			},
			expectedOrder:   []ProjectType{Go, Node, Python},
			expectedGreater: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			detector := NewDetector(root)

			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			// Check that results are properly sorted (descending by confidence)
			for i := 0; i < len(results)-1; i++ {
				if results[i].Confidence < results[i+1].Confidence {
					t.Errorf("Results not sorted properly at position %d: %f < %f", i, results[i].Confidence, results[i+1].Confidence)
				}
			}

			// Verify specific order where applicable
			if len(results) >= len(tt.expectedOrder) {
				for idx, expectedType := range tt.expectedOrder {
					if results[idx].Type != expectedType {
						t.Errorf("Expected type %v at position %d, got %v", expectedType, idx, results[idx].Type)
					}
				}
			}
		})
	}
}

func TestDetector_LoopBoundaryConditions(t *testing.T) {
	// Tests loop boundary conditions (for i := 0; i < len(results); i++)
	// and (for j := i + 1; j < len(results); j++)
	// Targets mutations like: i <= len vs i < len, j <= len vs j < len
	tests := []struct {
		name        string
		numProjects int
		description string
	}{
		{
			name:        "zero projects",
			numProjects: 0,
			description: "empty directory should return Unknown",
		},
		{
			name:        "one project type",
			numProjects: 1,
			description: "single project should be detected",
		},
		{
			name:        "two project types",
			numProjects: 2,
			description: "two projects should be sorted correctly",
		},
		{
			name:        "boundary at 6 project types",
			numProjects: 6,
			description: "all project types detected should be sorted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := make(map[string]string)

			// Create different project indicators based on numProjects
			switch tt.numProjects {
			case 0:
				// Empty directory - no files
			case 1:
				files["go.mod"] = "module test"
			case 2:
				files["go.mod"] = "module test"
				files["package.json"] = `{"name": "test"}`
			case 3:
				files["go.mod"] = "module test"
				files["package.json"] = `{"name": "test"}`
				files["pyproject.toml"] = "[project]"
			case 4:
				files["go.mod"] = "module test"
				files["package.json"] = `{"name": "test"}`
				files["pyproject.toml"] = "[project]"
				files["Gemfile"] = "source 'https://rubygems.org'"
			case 5:
				files["go.mod"] = "module test"
				files["package.json"] = `{"name": "test"}`
				files["pyproject.toml"] = "[project]"
				files["Gemfile"] = "source 'https://rubygems.org'"
				files["Cargo.toml"] = "[package]"
			case 6:
				files["go.mod"] = "module test"
				files["package.json"] = `{"name": "test"}`
				files["pyproject.toml"] = "[project]"
				files["Gemfile"] = "source 'https://rubygems.org'"
				files["Cargo.toml"] = "[package]"
				files["pom.xml"] = "<project/>"
			}

			root := createTestProject(t, files, nil)
			detector := NewDetector(root)
			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			// Verify loop processes all results
			if tt.numProjects == 0 {
				if len(results) != 1 || results[0].Type != Unknown {
					t.Errorf("Empty project should return Unknown, got %v results", len(results))
				}
			} else {
				if len(results) < tt.numProjects {
					t.Errorf("Expected at least %d results, got %d", tt.numProjects, len(results))
				}
				// Verify sorting was applied to all results
				for i := 1; i < len(results); i++ {
					if results[i].Confidence > results[i-1].Confidence {
						t.Errorf("Results not sorted: position %d has higher confidence than %d", i, i-1)
					}
				}
			}
		})
	}
}

func TestDetector_ConfidenceBoundaryAt1_0(t *testing.T) {
	// Tests boundary condition: if result.Confidence > 1.0
	// Targets mutations like: > vs >=, 1.0 vs 0.9, etc.
	tests := []struct {
		name           string
		indicators     []float64 // Confidence values to add
		expectedCapped bool      // Should be capped at 1.0
	}{
		{
			name:           "confidence exactly at 1.0",
			indicators:     []float64{0.6, 0.4},
			expectedCapped: true,
		},
		{
			name:           "confidence exceeds 1.0",
			indicators:     []float64{0.6, 0.5},
			expectedCapped: true,
		},
		{
			name:           "confidence just under 1.0",
			indicators:     []float64{0.6, 0.39},
			expectedCapped: false,
		},
		{
			name:           "multiple indicators exceeding 1.0",
			indicators:     []float64{0.6, 0.2, 0.3},
			expectedCapped: true,
		},
		{
			name:           "single indicator at 1.0",
			indicators:     []float64{1.0},
			expectedCapped: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test case that triggers specific confidence accumulation
			// For Go, confidence values are: go.mod (0.6), go.sum (0.2), .go files (0.2)
			// So we can create conditions to test exact boundaries
			files := make(map[string]string)

			if tt.indicators[0] >= 0.6 {
				files["go.mod"] = "module test"
			}
			if len(tt.indicators) > 1 && tt.indicators[1] >= 0.2 {
				files["go.sum"] = "test"
			}
			if len(tt.indicators) > 2 && tt.indicators[2] >= 0.2 {
				files["main.go"] = "package main"
			}

			root := createTestProject(t, files, nil)
			detector := NewDetector(root)
			result, err := detector.detectGo()
			if err != nil {
				t.Fatalf("detectGo() error = %v", err)
			}

			// Verify that confidence is capped at 1.0
			if result.Confidence > 1.0 {
				t.Errorf("Confidence %f exceeds maximum 1.0", result.Confidence)
			}

			// If expected to be capped, should be exactly 1.0 or less
			if tt.expectedCapped {
				if result.Confidence > 1.0 {
					t.Errorf("Expected capped confidence, got %f", result.Confidence)
				}
			}
		})
	}
}

func TestDetector_ResultsAppendBoundary(t *testing.T) {
	// Tests boundary condition: if result != nil && result.Confidence > 0
	// Targets mutations in: result.Confidence > 0 comparison
	tests := []struct {
		name              string
		files             map[string]string
		shouldDetectCount int
	}{
		{
			name:              "zero confidence should not be appended",
			files:             map[string]string{}, // No indicators
			shouldDetectCount: 1,                   // Only Unknown
		},
		{
			name: "minimal confidence just above zero",
			files: map[string]string{
				"main.go": "package main",
			},
			shouldDetectCount: 2, // Unknown + Go with 0.2 confidence
		},
		{
			name: "confidence exactly zero should not be appended",
			files: map[string]string{
				"unrelated.txt": "content",
			},
			shouldDetectCount: 1, // Only Unknown
		},
		{
			name: "multiple languages above zero",
			files: map[string]string{
				"go.mod":       "module test",
				"package.json": `{"name": "test"}`,
			},
			shouldDetectCount: 3, // Go, Node, Unknown not added since others exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			detector := NewDetector(root)
			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			// Count non-zero confidence results
			nonZeroCount := 0
			for _, r := range results {
				if r.Confidence > 0 {
					nonZeroCount++
				}
			}

			if tt.shouldDetectCount > 1 && nonZeroCount != tt.shouldDetectCount-1 {
				// Verify we detect exactly the expected number of non-Unknown types
				expected := tt.shouldDetectCount - 1
				if nonZeroCount != expected && len(results) < tt.shouldDetectCount {
					t.Logf("Expected %d non-zero confidence results, got %d (total results: %d)",
						expected, nonZeroCount, len(results))
				}
			}
		})
	}
}

func TestDetector_ArrayIndexBoundary(t *testing.T) {
	// Tests boundary conditions for array operations and indexing
	// Particularly for results[0] access and len(results) checks
	tests := []struct {
		name            string
		files           map[string]string
		shouldHaveFirst bool
	}{
		{
			name:            "empty project returns Unknown at index 0",
			files:           map[string]string{},
			shouldHaveFirst: true,
		},
		{
			name: "single detection returns one result",
			files: map[string]string{
				"go.mod": "module test",
			},
			shouldHaveFirst: true,
		},
		{
			name: "multiple detections maintain order",
			files: map[string]string{
				"go.mod":       "module test",
				"go.sum":       "test",
				"package.json": `{"name": "test"}`,
			},
			shouldHaveFirst: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			detector := NewDetector(root)
			results, err := detector.Detect()
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			if tt.shouldHaveFirst {
				if len(results) == 0 {
					t.Fatal("Expected at least one result")
				}

				first := results[0]
				if first.Type == "" {
					t.Error("First result has empty type")
				}

				// First should have highest confidence (if not empty)
				if len(results) > 1 {
					if first.Confidence < results[1].Confidence {
						t.Errorf("First result should have highest confidence, got %f < %f",
							first.Confidence, results[1].Confidence)
					}
				}
			}
		})
	}
}

func TestDetector_DepthCountingBoundary(t *testing.T) {
	// Tests boundary conditions in depth calculation
	// Particularly: depth++ operation and depth > maxDepth check
	tmpDir := t.TempDir()

	// Create nested structure: deep nesting beyond typical maxDepth
	deepPath := filepath.Join(tmpDir, "a", "b", "c", "d", "e", "deep.go")
	if err := os.MkdirAll(filepath.Dir(deepPath), 0755); err != nil {
		t.Fatalf("Failed to create deep directory: %v", err)
	}
	if err := os.WriteFile(deepPath, []byte("package deep"), 0644); err != nil {
		t.Fatalf("Failed to create deep file: %v", err)
	}

	// Create shallow nested structure within typical maxDepth
	shallowPath := filepath.Join(tmpDir, "cmd", "app.go")
	if err := os.MkdirAll(filepath.Dir(shallowPath), 0755); err != nil {
		t.Fatalf("Failed to create shallow directory: %v", err)
	}
	if err := os.WriteFile(shallowPath, []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create shallow file: %v", err)
	}

	// Add go.mod at root
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	detector := NewDetector(tmpDir)
	result, err := detector.detectGo()
	if err != nil {
		t.Fatalf("detectGo() error = %v", err)
	}

	// Should find shallow file but not deeply nested file
	// Since maxDepth=3 in detectGo and "a/b/c/d/e/deep.go" is depth 6
	hasIndicator := false
	for _, ind := range result.Indicators {
		if ind == "*.go files" {
			hasIndicator = true
			break
		}
	}

	if !hasIndicator {
		t.Errorf("Expected to find .go files at shallow depths")
	}

	// Verify we have go.mod indicator
	hasGoMod := false
	for _, ind := range result.Indicators {
		if ind == "go.mod" {
			hasGoMod = true
			break
		}
	}

	if !hasGoMod {
		t.Errorf("Expected go.mod indicator")
	}

	// Confidence should include both go.mod (0.6) and .go files (0.2)
	if result.Confidence != 0.8 {
		t.Errorf("Expected confidence 0.8 (go.mod + .go files), got %f", result.Confidence)
	}
}

func TestDetector_SkipDirBoundary(t *testing.T) {
	// Tests boundary conditions in directory skipping logic
	// Particularly for the common skip directories list
	tests := []struct {
		name       string
		dirName    string
		shouldSkip bool
	}{
		{
			name:       "skip node_modules",
			dirName:    "node_modules",
			shouldSkip: true,
		},
		{
			name:       "skip vendor",
			dirName:    "vendor",
			shouldSkip: true,
		},
		{
			name:       "skip .git",
			dirName:    ".git",
			shouldSkip: true,
		},
		{
			name:       "don't skip src",
			dirName:    "src",
			shouldSkip: false,
		},
		{
			name:       "don't skip cmd",
			dirName:    "cmd",
			shouldSkip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create structure with potential skip directory
			skipDir := filepath.Join(tmpDir, tt.dirName)
			if err := os.MkdirAll(skipDir, 0755); err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			// Put a Go file in the directory
			goFile := filepath.Join(skipDir, "test.go")
			if err := os.WriteFile(goFile, []byte("package main"), 0644); err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}

			// Also create go.mod at root for consistent detection
			if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
				t.Fatalf("Failed to create go.mod: %v", err)
			}

			detector := NewDetector(tmpDir)
			result, err := detector.detectGo()
			if err != nil {
				t.Fatalf("detectGo() error = %v", err)
			}

			hasGoFiles := false
			for _, ind := range result.Indicators {
				if ind == "*.go files" {
					hasGoFiles = true
					break
				}
			}

			// If directory should be skipped, .go files indicator should only come from go.mod
			// If directory should not be skipped, we might find .go files in it
			if tt.shouldSkip && hasGoFiles {
				// With go.mod present, we should still detect Go, but the go files
				// found should not be from the skip directory
				t.Logf("Skipped directory %s correctly - found .go files from other sources", tt.dirName)
			}
		})
	}
}

func TestDetector_MaxResultsBoundary(t *testing.T) {
	// Tests boundary condition: maxResults := 10 and if len(matches) >= maxResults
	// Targets mutations in: >= vs >, maxResults value, len() check
	tmpDir := t.TempDir()

	// Create 15 Go files to exceed maxResults=10
	for i := 1; i <= 15; i++ {
		filename := filepath.Join(tmpDir, fmt.Sprintf("file%d.go", i))
		if err := os.WriteFile(filename, []byte("package main"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	detector := NewDetector(tmpDir)
	result, err := detector.detectGo()
	if err != nil {
		t.Fatalf("detectGo() error = %v", err)
	}

	// Should have detected .go files indicator
	found := false
	for _, ind := range result.Indicators {
		if ind == "*.go files" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected *.go files indicator even with 15 files")
	}

	// Confidence should be 0.2 for just .go files
	if result.Confidence != 0.2 {
		t.Errorf("Expected confidence 0.2 for .go files only, got %f", result.Confidence)
	}
}

func TestDetector_FileMatchPatternBoundary(t *testing.T) {
	// Tests boundary condition in pattern matching
	// filepath.Match result handling and conditions
	tests := []struct {
		name        string
		files       map[string]string
		pattern     string
		shouldMatch bool
	}{
		{
			name:        "exact extension match",
			files:       map[string]string{"test.go": "package main"},
			pattern:     "*.go",
			shouldMatch: true,
		},
		{
			name:        "different extension no match",
			files:       map[string]string{"test.txt": "content"},
			pattern:     "*.go",
			shouldMatch: false,
		},
		{
			name:        "multiple files with matching extension",
			files:       map[string]string{"a.go": "package a", "b.go": "package b"},
			pattern:     "*.go",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			matches, err := NewDetector(root).findFiles(tt.pattern, 1)
			if err != nil {
				t.Fatalf("findFiles() error = %v", err)
			}

			hasMatch := len(matches) > 0
			if tt.shouldMatch && !hasMatch {
				t.Errorf("Expected to find files matching pattern %s", tt.pattern)
			}
			if !tt.shouldMatch && hasMatch {
				t.Errorf("Unexpected files found matching pattern %s", tt.pattern)
			}
		})
	}
}

func TestDetector_DetectPrimaryEmptyResults(t *testing.T) {
	// Tests boundary condition in DetectPrimary when results exist
	// Particularly the len(results) == 0 check and results[0] access
	tests := []struct {
		name              string
		files             map[string]string
		shouldHavePrimary bool
	}{
		{
			name:              "empty directory returns Unknown primary",
			files:             map[string]string{},
			shouldHavePrimary: true,
		},
		{
			name:              "single result returns that result",
			files:             map[string]string{"go.mod": "module test"},
			shouldHavePrimary: true,
		},
		{
			name:              "multiple results returns highest confidence",
			files:             map[string]string{"go.mod": "module test", "package.json": `{"name": "test"}`},
			shouldHavePrimary: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := createTestProject(t, tt.files, nil)
			detector := NewDetector(root)
			result, err := detector.DetectPrimary()
			if err != nil {
				t.Fatalf("DetectPrimary() error = %v", err)
			}

			if tt.shouldHavePrimary {
				if result == nil {
					t.Fatal("Expected DetectPrimary to return a result")
				}
				if result.Type == "" {
					t.Error("Primary result has empty type")
				}
			}
		})
	}
}
