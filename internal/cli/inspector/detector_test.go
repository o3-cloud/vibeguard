package inspector

import (
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
