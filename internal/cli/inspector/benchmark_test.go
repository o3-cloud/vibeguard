package inspector

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// createLargeProject creates a simulated large project for benchmarking.
// It creates a project with many nested directories and files to stress-test
// the file scanning logic.
func createLargeProject(b *testing.B, projectType string, numDirs, filesPerDir int) string {
	b.Helper()
	tmpDir := b.TempDir()

	// Create base project files
	switch projectType {
	case "go":
		benchWriteFile(b, tmpDir, "go.mod", "module example.com/largeproject\n\ngo 1.21")
		benchWriteFile(b, tmpDir, "go.sum", "# checksums")
		benchWriteFile(b, tmpDir, ".golangci.yml", "linters:\n  enable:\n    - gofmt")
		benchWriteFile(b, tmpDir, "main.go", "package main\n\nfunc main() {}")
	case "node":
		benchWriteFile(b, tmpDir, "package.json", `{"name": "large-project", "devDependencies": {"eslint": "^8.0.0"}}`)
		benchWriteFile(b, tmpDir, "package-lock.json", "{}")
		benchWriteFile(b, tmpDir, ".eslintrc.json", "{}")
		_ = os.MkdirAll(filepath.Join(tmpDir, "node_modules"), 0755)
	case "python":
		benchWriteFile(b, tmpDir, "pyproject.toml", "[tool.pytest]\n[tool.black]")
		benchWriteFile(b, tmpDir, "requirements.txt", "flask\nrequests")
		benchWriteFile(b, tmpDir, "setup.py", "from setuptools import setup")
	}

	// Create nested directory structure
	extensions := map[string]string{"go": ".go", "node": ".js", "python": ".py"}
	ext := extensions[projectType]

	for i := 0; i < numDirs; i++ {
		dirPath := filepath.Join(tmpDir, fmt.Sprintf("pkg%d", i))
		_ = os.MkdirAll(dirPath, 0755)
		for j := 0; j < filesPerDir; j++ {
			filename := fmt.Sprintf("file%d%s", j, ext)
			content := fmt.Sprintf("// file %d in dir %d", j, i)
			benchWriteFile(b, dirPath, filename, content)
		}
	}

	// Create some deeply nested paths
	deepPath := filepath.Join(tmpDir, "a", "b", "c", "d", "e")
	_ = os.MkdirAll(deepPath, 0755)
	for i := 0; i < filesPerDir; i++ {
		benchWriteFile(b, deepPath, fmt.Sprintf("deep%d%s", i, ext), "// deep file")
	}

	return tmpDir
}

func benchWriteFile(b *testing.B, dir, name, content string) {
	b.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		b.Fatalf("failed to write %s: %v", path, err)
	}
}

// BenchmarkDetector_LargeGoProject tests detector performance on a large Go project.
func BenchmarkDetector_LargeGoProject(b *testing.B) {
	tmpDir := createLargeProject(b, "go", 50, 20) // 50 dirs * 20 files = 1000+ files
	detector := NewDetector(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.Detect()
		if err != nil {
			b.Fatalf("Detect() error = %v", err)
		}
	}
}

// BenchmarkDetector_LargeNodeProject tests detector performance on a large Node project.
func BenchmarkDetector_LargeNodeProject(b *testing.B) {
	tmpDir := createLargeProject(b, "node", 50, 20)
	detector := NewDetector(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.Detect()
		if err != nil {
			b.Fatalf("Detect() error = %v", err)
		}
	}
}

// BenchmarkDetector_LargePythonProject tests detector performance on a large Python project.
func BenchmarkDetector_LargePythonProject(b *testing.B) {
	tmpDir := createLargeProject(b, "python", 50, 20)
	detector := NewDetector(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.Detect()
		if err != nil {
			b.Fatalf("Detect() error = %v", err)
		}
	}
}

// BenchmarkToolScanner_LargeProject tests tool scanner on a large project with many configs.
func BenchmarkToolScanner_LargeProject(b *testing.B) {
	tmpDir := createLargeProject(b, "go", 50, 20)

	// Add more tool configs
	benchWriteFile(b, tmpDir, ".pre-commit-config.yaml", "repos: []")
	benchWriteFile(b, tmpDir, "Makefile", "all:\n\tgo build")

	// Create GitHub Actions workflow
	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	_ = os.MkdirAll(workflowDir, 0755)
	benchWriteFile(b, workflowDir, "ci.yml", "name: CI\non: push")
	benchWriteFile(b, workflowDir, "release.yml", "name: Release")

	scanner := NewToolScanner(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.ScanAll()
		if err != nil {
			b.Fatalf("ScanAll() error = %v", err)
		}
	}
}

// BenchmarkMetadataExtractor_LargeGoProject tests metadata extraction performance.
func BenchmarkMetadataExtractor_LargeGoProject(b *testing.B) {
	tmpDir := createLargeProject(b, "go", 50, 20)

	// Add realistic go.mod content
	goModContent := `module example.com/largeproject

go 1.21

require (
	github.com/stretchr/testify v1.8.0
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
)
`
	benchWriteFile(b, tmpDir, "go.mod", goModContent)

	extractor := NewMetadataExtractor(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := extractor.Extract(Go)
		if err != nil {
			b.Fatalf("Extract() error = %v", err)
		}
		_, err = extractor.ExtractStructure(Go)
		if err != nil {
			b.Fatalf("ExtractStructure() error = %v", err)
		}
	}
}

// BenchmarkFullInspection tests the full inspection flow.
func BenchmarkFullInspection(b *testing.B) {
	tmpDir := createLargeProject(b, "go", 30, 15)

	// Add more configs for realistic scenario
	benchWriteFile(b, tmpDir, ".golangci.yml", "linters:\n  enable:\n    - gofmt")
	benchWriteFile(b, tmpDir, "Makefile", "build:\n\tgo build")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate full inspection flow
		detector := NewDetector(tmpDir)
		result, err := detector.DetectPrimary()
		if err != nil {
			b.Fatalf("DetectPrimary() error = %v", err)
		}

		scanner := NewToolScanner(tmpDir)
		_, err = scanner.ScanForProjectType(result.Type)
		if err != nil {
			b.Fatalf("ScanForProjectType() error = %v", err)
		}

		extractor := NewMetadataExtractor(tmpDir)
		_, err = extractor.Extract(result.Type)
		if err != nil {
			b.Fatalf("Extract() error = %v", err)
		}
		_, err = extractor.ExtractStructure(result.Type)
		if err != nil {
			b.Fatalf("ExtractStructure() error = %v", err)
		}
	}
}

// BenchmarkFileCache tests the file cache performance improvement.
func BenchmarkFileCache(b *testing.B) {
	tmpDir := b.TempDir()

	// Create test files
	files := []string{"go.mod", "go.sum", "main.go", ".golangci.yml", "Makefile"}
	for _, f := range files {
		benchWriteFile(b, tmpDir, f, "content")
	}

	b.Run("WithoutCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, f := range files {
				path := filepath.Join(tmpDir, f)
				_, _ = os.Stat(path)
			}
		}
	})

	b.Run("WithCache", func(b *testing.B) {
		cache := NewFileCache(tmpDir)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, f := range files {
				cache.FileExists(f)
			}
		}
	})
}

// BenchmarkFileCache_ReadFile tests cached file reading.
func BenchmarkFileCache_ReadFile(b *testing.B) {
	tmpDir := b.TempDir()

	// Create a larger config file
	content := `linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign

linters-settings:
  gofmt:
    simplify: true
  govet:
    check-shadowing: true
`
	benchWriteFile(b, tmpDir, ".golangci.yml", content)

	b.Run("WithoutCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			path := filepath.Join(tmpDir, ".golangci.yml")
			_, _ = os.ReadFile(path)
		}
	})

	b.Run("WithCache", func(b *testing.B) {
		cache := NewFileCache(tmpDir)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.ReadFile(".golangci.yml")
		}
	})
}

// TestPerformanceTargets verifies that the performance targets are met.
// This is a test rather than a benchmark to provide pass/fail results.
func TestPerformanceTargets(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	tmpDir := createTestProject(t, map[string]string{
		"go.mod":            "module test\n\ngo 1.21",
		"go.sum":            "...",
		"main.go":           "package main",
		".golangci.yml":     "linters: {}",
		"Makefile":          "all: build",
		"cmd/app/main.go":   "package main",
		"internal/pkg/x.go": "package pkg",
	}, []string{".github/workflows"})

	// Create workflow file
	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte("name: CI"), 0644); err != nil {
		t.Fatal(err)
	}

	// Use direct timing for fast sanity checks (not testing.Benchmark which runs until stable)
	const iterations = 10

	// Test Detector
	detector := NewDetector(tmpDir)
	var detectElapsed time.Duration
	for i := 0; i < iterations; i++ {
		start := time.Now()
		_, _ = detector.DetectPrimary()
		detectElapsed += time.Since(start)
	}
	avgDetect := detectElapsed / iterations
	if avgDetect > 50*time.Millisecond {
		t.Errorf("DetectPrimary too slow: %v per op (target: <50ms)", avgDetect)
	}
	t.Logf("DetectPrimary: %v per op", avgDetect)

	// Test ToolScanner
	scanner := NewToolScanner(tmpDir)
	var scanElapsed time.Duration
	for i := 0; i < iterations; i++ {
		start := time.Now()
		_, _ = scanner.ScanAll()
		scanElapsed += time.Since(start)
	}
	avgScan := scanElapsed / iterations
	if avgScan > 50*time.Millisecond {
		t.Errorf("ScanAll too slow: %v per op (target: <50ms)", avgScan)
	}
	t.Logf("ScanAll: %v per op", avgScan)

	// Test full flow
	var fullElapsed time.Duration
	for i := 0; i < iterations; i++ {
		start := time.Now()
		d := NewDetector(tmpDir)
		result, _ := d.DetectPrimary()
		s := NewToolScanner(tmpDir)
		_, _ = s.ScanForProjectType(result.Type)
		e := NewMetadataExtractor(tmpDir)
		_, _ = e.Extract(result.Type)
		_, _ = e.ExtractStructure(result.Type)
		fullElapsed += time.Since(start)
	}
	avgFull := fullElapsed / iterations
	if avgFull > 100*time.Millisecond {
		t.Errorf("Full inspection flow too slow: %v per op (target: <100ms)", avgFull)
	}
	t.Logf("Full flow: %v per op", avgFull)
}
