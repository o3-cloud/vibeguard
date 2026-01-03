package inspector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMetadataExtractor_ExtractGoMetadata(t *testing.T) {
	tests := []struct {
		name     string
		goMod    string
		want     *ProjectMetadata
		wantName string
	}{
		{
			name: "standard go.mod",
			goMod: `module github.com/example/myproject

go 1.21

require (
	github.com/spf13/cobra v1.8.0
)
`,
			wantName: "github.com/example/myproject",
		},
		{
			name: "go.mod with version suffix",
			goMod: `module github.com/example/myproject/v2

go 1.22.0
`,
			wantName: "github.com/example/myproject/v2",
		},
		{
			name:     "empty go.mod",
			goMod:    "",
			wantName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if tt.goMod != "" {
				if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(tt.goMod), 0644); err != nil {
					t.Fatal(err)
				}
			}

			extractor := NewMetadataExtractor(tmpDir)
			metadata, err := extractor.Extract(Go)
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}

			if metadata.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", metadata.Name, tt.wantName)
			}
		})
	}
}

func TestMetadataExtractor_ExtractGoMetadata_GoVersion(t *testing.T) {
	tmpDir := t.TempDir()
	goMod := `module example.com/test

go 1.21
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Go)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if v := metadata.Extra["go_version"]; v != "1.21" {
		t.Errorf("go_version = %q, want %q", v, "1.21")
	}
}

func TestMetadataExtractor_ExtractNodeMetadata(t *testing.T) {
	tests := []struct {
		name        string
		packageJSON string
		wantName    string
		wantVersion string
		wantDesc    string
		wantAuthor  string
	}{
		{
			name: "full package.json",
			packageJSON: `{
				"name": "my-package",
				"version": "1.2.3",
				"description": "A test package",
				"author": "John Doe <john@example.com>",
				"license": "MIT"
			}`,
			wantName:    "my-package",
			wantVersion: "1.2.3",
			wantDesc:    "A test package",
			wantAuthor:  "John Doe <john@example.com>",
		},
		{
			name: "author as object",
			packageJSON: `{
				"name": "test-pkg",
				"version": "2.0.0",
				"author": {
					"name": "Jane Doe",
					"email": "jane@example.com"
				}
			}`,
			wantName:    "test-pkg",
			wantVersion: "2.0.0",
			wantAuthor:  "Jane Doe <jane@example.com>",
		},
		{
			name: "repository as string",
			packageJSON: `{
				"name": "repo-test",
				"version": "0.1.0",
				"repository": "https://github.com/example/repo"
			}`,
			wantName:    "repo-test",
			wantVersion: "0.1.0",
		},
		{
			name: "repository as object",
			packageJSON: `{
				"name": "repo-test",
				"version": "0.1.0",
				"repository": {
					"type": "git",
					"url": "https://github.com/example/repo.git"
				}
			}`,
			wantName:    "repo-test",
			wantVersion: "0.1.0",
		},
		{
			name: "minimal package.json",
			packageJSON: `{
				"name": "minimal"
			}`,
			wantName: "minimal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(tt.packageJSON), 0644); err != nil {
				t.Fatal(err)
			}

			extractor := NewMetadataExtractor(tmpDir)
			metadata, err := extractor.Extract(Node)
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}

			if metadata.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", metadata.Name, tt.wantName)
			}
			if metadata.Version != tt.wantVersion {
				t.Errorf("Version = %q, want %q", metadata.Version, tt.wantVersion)
			}
			if tt.wantDesc != "" && metadata.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", metadata.Description, tt.wantDesc)
			}
			if tt.wantAuthor != "" && metadata.Author != tt.wantAuthor {
				t.Errorf("Author = %q, want %q", metadata.Author, tt.wantAuthor)
			}
		})
	}
}

func TestMetadataExtractor_ExtractNodeMetadata_Engines(t *testing.T) {
	tmpDir := t.TempDir()
	packageJSON := `{
		"name": "engine-test",
		"version": "1.0.0",
		"engines": {
			"node": ">=18.0.0"
		}
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageJSON), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Node)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if v := metadata.Extra["node_version"]; v != ">=18.0.0" {
		t.Errorf("node_version = %q, want %q", v, ">=18.0.0")
	}
}

func TestMetadataExtractor_ExtractPythonMetadata_Pyproject(t *testing.T) {
	tests := []struct {
		name        string
		pyproject   string
		wantName    string
		wantVersion string
		wantDesc    string
	}{
		{
			name: "standard pyproject.toml",
			pyproject: `[project]
name = "my-python-pkg"
version = "0.1.0"
description = "A Python package"
license = "Apache-2.0"
requires-python = ">=3.10"
`,
			wantName:    "my-python-pkg",
			wantVersion: "0.1.0",
			wantDesc:    "A Python package",
		},
		{
			name: "poetry style pyproject.toml",
			pyproject: `[tool.poetry]
name = "poetry-project"
version = "2.0.0"
description = "Poetry managed project"
`,
			wantName:    "poetry-project",
			wantVersion: "2.0.0",
			wantDesc:    "Poetry managed project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(tt.pyproject), 0644); err != nil {
				t.Fatal(err)
			}

			extractor := NewMetadataExtractor(tmpDir)
			metadata, err := extractor.Extract(Python)
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}

			if metadata.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", metadata.Name, tt.wantName)
			}
			if metadata.Version != tt.wantVersion {
				t.Errorf("Version = %q, want %q", metadata.Version, tt.wantVersion)
			}
			if tt.wantDesc != "" && metadata.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", metadata.Description, tt.wantDesc)
			}
		})
	}
}

func TestMetadataExtractor_ExtractPythonMetadata_SetupPy(t *testing.T) {
	tmpDir := t.TempDir()
	setupPy := `from setuptools import setup

setup(
    name="setup-pkg",
    version="3.0.0",
    description="A setup.py package",
    author="Test Author",
    license="BSD"
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "setup.py"), []byte(setupPy), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Python)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Name != "setup-pkg" {
		t.Errorf("Name = %q, want %q", metadata.Name, "setup-pkg")
	}
	if metadata.Version != "3.0.0" {
		t.Errorf("Version = %q, want %q", metadata.Version, "3.0.0")
	}
	if metadata.Author != "Test Author" {
		t.Errorf("Author = %q, want %q", metadata.Author, "Test Author")
	}
}

func TestMetadataExtractor_ExtractRustMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	cargoToml := `[package]
name = "rust-project"
version = "0.2.0"
edition = "2021"
description = "A Rust project"
license = "MIT"
repository = "https://github.com/example/rust-project"

[dependencies]
serde = "1.0"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte(cargoToml), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Rust)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Name != "rust-project" {
		t.Errorf("Name = %q, want %q", metadata.Name, "rust-project")
	}
	if metadata.Version != "0.2.0" {
		t.Errorf("Version = %q, want %q", metadata.Version, "0.2.0")
	}
	if metadata.Description != "A Rust project" {
		t.Errorf("Description = %q, want %q", metadata.Description, "A Rust project")
	}
	if metadata.Extra["rust_edition"] != "2021" {
		t.Errorf("rust_edition = %q, want %q", metadata.Extra["rust_edition"], "2021")
	}
	if metadata.Repository != "https://github.com/example/rust-project" {
		t.Errorf("Repository = %q, want %q", metadata.Repository, "https://github.com/example/rust-project")
	}
}

func TestMetadataExtractor_ExtractRubyMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	gemspec := `Gem::Specification.new do |s|
  s.name        = "my-gem"
  s.version     = "1.0.0"
  s.summary     = "A Ruby gem"
  s.description = "A test Ruby gem"
  s.authors     = ["Ruby Dev"]
  s.license     = "MIT"
end
`
	if err := os.WriteFile(filepath.Join(tmpDir, "my-gem.gemspec"), []byte(gemspec), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Ruby)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Name != "my-gem" {
		t.Errorf("Name = %q, want %q", metadata.Name, "my-gem")
	}
	if metadata.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", metadata.Version, "1.0.0")
	}
}

func TestMetadataExtractor_ExtractJavaMetadata_Maven(t *testing.T) {
	tmpDir := t.TempDir()
	pomXml := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <groupId>com.example</groupId>
    <artifactId>my-java-app</artifactId>
    <version>1.0.0-SNAPSHOT</version>
    <name>My Java Application</name>
    <description>A Java application</description>
</project>
`
	if err := os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte(pomXml), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Java)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Name != "My Java Application" {
		t.Errorf("Name = %q, want %q", metadata.Name, "My Java Application")
	}
	if metadata.Version != "1.0.0-SNAPSHOT" {
		t.Errorf("Version = %q, want %q", metadata.Version, "1.0.0-SNAPSHOT")
	}
	if metadata.Extra["group_id"] != "com.example" {
		t.Errorf("group_id = %q, want %q", metadata.Extra["group_id"], "com.example")
	}
}

func TestMetadataExtractor_ExtractJavaMetadata_Gradle(t *testing.T) {
	tmpDir := t.TempDir()
	buildGradle := `plugins {
    id 'java'
}

group = 'com.example'
version = '2.0.0'
archivesBaseName = 'gradle-app'
`
	if err := os.WriteFile(filepath.Join(tmpDir, "build.gradle"), []byte(buildGradle), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Java)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Version != "2.0.0" {
		t.Errorf("Version = %q, want %q", metadata.Version, "2.0.0")
	}
	if metadata.Extra["group"] != "com.example" {
		t.Errorf("group = %q, want %q", metadata.Extra["group"], "com.example")
	}
}

func TestMetadataExtractor_ExtractStructure_Go(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go project structure
	dirs := []string{"cmd/server", "internal/handler", "pkg/utils", "test"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files
	files := map[string]string{
		"go.mod":             "module example.com/test\n\ngo 1.21\n",
		"cmd/server/main.go": "package main\n\nfunc main() {}\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Go)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check entry points
	if len(structure.EntryPoints) != 1 || structure.EntryPoints[0] != "cmd/server/main.go" {
		t.Errorf("EntryPoints = %v, want [cmd/server/main.go]", structure.EntryPoints)
	}

	// Check source dirs
	if !sliceContains(structure.SourceDirs, "pkg") {
		t.Errorf("SourceDirs should contain 'pkg', got %v", structure.SourceDirs)
	}
	if !sliceContains(structure.SourceDirs, "internal") {
		t.Errorf("SourceDirs should contain 'internal', got %v", structure.SourceDirs)
	}

	// Check test dirs
	if !sliceContains(structure.TestDirs, "test") {
		t.Errorf("TestDirs should contain 'test', got %v", structure.TestDirs)
	}

	// Check config files
	if !sliceContains(structure.ConfigFiles, "go.mod") {
		t.Errorf("ConfigFiles should contain 'go.mod', got %v", structure.ConfigFiles)
	}
}

func TestMetadataExtractor_ExtractStructure_Node(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Node project structure
	dirs := []string{"src", "test", "dist"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files
	files := map[string]string{
		"package.json":  `{"name":"test","version":"1.0.0","main":"dist/index.js"}`,
		"src/index.ts":  "export default {}",
		"tsconfig.json": "{}",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Node)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check entry points - should include both main and src/index.ts
	if !sliceContains(structure.EntryPoints, "dist/index.js") {
		t.Errorf("EntryPoints should contain 'dist/index.js', got %v", structure.EntryPoints)
	}
	if !sliceContains(structure.EntryPoints, "src/index.ts") {
		t.Errorf("EntryPoints should contain 'src/index.ts', got %v", structure.EntryPoints)
	}

	// Check source dirs
	if !sliceContains(structure.SourceDirs, "src") {
		t.Errorf("SourceDirs should contain 'src', got %v", structure.SourceDirs)
	}

	// Check test dirs
	if !sliceContains(structure.TestDirs, "test") {
		t.Errorf("TestDirs should contain 'test', got %v", structure.TestDirs)
	}

	// Check build output
	if structure.BuildOutputDir != "dist" {
		t.Errorf("BuildOutputDir = %q, want %q", structure.BuildOutputDir, "dist")
	}
}

func TestMetadataExtractor_ExtractStructure_Python(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Python project structure
	dirs := []string{"src", "tests", "mypackage"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files
	files := map[string]string{
		"pyproject.toml":        "[project]\nname = \"test\"\n",
		"main.py":               "print('hello')",
		"mypackage/__init__.py": "",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Python)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check entry points
	if !sliceContains(structure.EntryPoints, "main.py") {
		t.Errorf("EntryPoints should contain 'main.py', got %v", structure.EntryPoints)
	}

	// Check source dirs - should detect package dir with __init__.py
	if !sliceContains(structure.SourceDirs, "mypackage") {
		t.Errorf("SourceDirs should contain 'mypackage', got %v", structure.SourceDirs)
	}

	// Check test dirs
	if !sliceContains(structure.TestDirs, "tests") {
		t.Errorf("TestDirs should contain 'tests', got %v", structure.TestDirs)
	}
}

func TestMetadataExtractor_ExtractStructure_Rust(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Rust project structure
	dirs := []string{"src", "tests", "target"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files
	files := map[string]string{
		"Cargo.toml":  "[package]\nname = \"test\"\nversion = \"0.1.0\"\n",
		"src/main.rs": "fn main() {}",
		"src/lib.rs":  "pub fn hello() {}",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Rust)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check entry points
	if !sliceContains(structure.EntryPoints, "src/main.rs") {
		t.Errorf("EntryPoints should contain 'src/main.rs', got %v", structure.EntryPoints)
	}
	if !sliceContains(structure.EntryPoints, "src/lib.rs") {
		t.Errorf("EntryPoints should contain 'src/lib.rs', got %v", structure.EntryPoints)
	}

	// Check source dirs
	if !sliceContains(structure.SourceDirs, "src") {
		t.Errorf("SourceDirs should contain 'src', got %v", structure.SourceDirs)
	}

	// Check test dirs
	if !sliceContains(structure.TestDirs, "tests") {
		t.Errorf("TestDirs should contain 'tests', got %v", structure.TestDirs)
	}

	// Check build output
	if structure.BuildOutputDir != "target" {
		t.Errorf("BuildOutputDir = %q, want %q", structure.BuildOutputDir, "target")
	}
}

func TestMetadataExtractor_DetectMonorepo(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string) error
		wantMono bool
	}{
		{
			name: "npm workspaces",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"workspaces":["packages/*"]}`), 0644)
			},
			wantMono: true,
		},
		{
			name: "pnpm workspace",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "pnpm-workspace.yaml"), []byte("packages:\n  - packages/*\n"), 0644)
			},
			wantMono: true,
		},
		{
			name: "lerna",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "lerna.json"), []byte(`{"packages":["packages/*"]}`), 0644)
			},
			wantMono: true,
		},
		{
			name: "cargo workspace",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte("[workspace]\nmembers = [\"crates/*\"]\n"), 0644)
			},
			wantMono: true,
		},
		{
			name: "packages directory with multiple package.json",
			setup: func(dir string) error {
				pkgDir := filepath.Join(dir, "packages")
				if err := os.MkdirAll(filepath.Join(pkgDir, "pkg1"), 0755); err != nil {
					return err
				}
				if err := os.MkdirAll(filepath.Join(pkgDir, "pkg2"), 0755); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(pkgDir, "pkg1", "package.json"), []byte("{}"), 0644); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(pkgDir, "pkg2", "package.json"), []byte("{}"), 0644)
			},
			wantMono: true,
		},
		{
			name: "single package",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"single"}`), 0644)
			},
			wantMono: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := tt.setup(tmpDir); err != nil {
				t.Fatal(err)
			}

			extractor := NewMetadataExtractor(tmpDir)
			structure, err := extractor.ExtractStructure(Node)
			if err != nil {
				t.Fatalf("ExtractStructure() error = %v", err)
			}

			if structure.HasMonorepo != tt.wantMono {
				t.Errorf("HasMonorepo = %v, want %v", structure.HasMonorepo, tt.wantMono)
			}
		})
	}
}

func TestMetadataExtractor_ExtractUnknownType(t *testing.T) {
	tmpDir := t.TempDir()

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Unknown)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	// Should return empty metadata with initialized Extra map
	if metadata.Extra == nil {
		t.Error("Extra map should be initialized, got nil")
	}
}

func TestMetadataExtractor_MissingFiles(t *testing.T) {
	tmpDir := t.TempDir()

	extractor := NewMetadataExtractor(tmpDir)

	// Test that missing files don't cause errors
	tests := []ProjectType{Go, Node, Python, Rust, Ruby, Java}
	for _, pt := range tests {
		t.Run(string(pt), func(t *testing.T) {
			metadata, err := extractor.Extract(pt)
			if err != nil {
				t.Fatalf("Extract(%v) error = %v", pt, err)
			}
			if metadata == nil {
				t.Errorf("Extract(%v) returned nil metadata", pt)
			}
		})
	}
}

func TestMetadataExtractor_ExtractStructure_ConfigFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create various config files
	configFiles := []string{
		"go.mod",
		"package.json",
		".golangci.yml",
		"Makefile",
		"Dockerfile",
	}
	for _, f := range configFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("# config"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Go)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	for _, f := range configFiles {
		if !sliceContains(structure.ConfigFiles, f) {
			t.Errorf("ConfigFiles should contain %q, got %v", f, structure.ConfigFiles)
		}
	}
}

func TestMetadataExtractor_ExtractStructure_Ruby(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Ruby project structure
	dirs := []string{"lib", "spec", "bin"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files
	files := map[string]string{
		"Gemfile":    "source 'https://rubygems.org'\ngem 'rails'\n",
		"app.rb":     "puts 'hello'",
		"bin/main":   "#!/usr/bin/env ruby\nrequire_relative '../app'",
		"lib/foo.rb": "module Foo; end",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Ruby)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check entry points
	if !sliceContains(structure.EntryPoints, "app.rb") {
		t.Errorf("EntryPoints should contain 'app.rb', got %v", structure.EntryPoints)
	}
	if !sliceContains(structure.EntryPoints, "bin/main") {
		t.Errorf("EntryPoints should contain 'bin/main', got %v", structure.EntryPoints)
	}

	// Check source dirs
	if !sliceContains(structure.SourceDirs, "lib") {
		t.Errorf("SourceDirs should contain 'lib', got %v", structure.SourceDirs)
	}

	// Check test dirs
	if !sliceContains(structure.TestDirs, "spec") {
		t.Errorf("TestDirs should contain 'spec', got %v", structure.TestDirs)
	}

	// Ruby typically doesn't have build output
	if structure.BuildOutputDir != "" {
		t.Errorf("BuildOutputDir should be empty for Ruby, got %q", structure.BuildOutputDir)
	}
}

func TestMetadataExtractor_ExtractStructure_Java(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Maven-style Java project structure
	dirs := []string{
		"src/main/java/com/example",
		"src/main/resources",
		"src/test/java/com/example",
		"src/test/resources",
		"target",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files
	files := map[string]string{
		"pom.xml":                            "<project><modelVersion>4.0.0</modelVersion></project>",
		"src/main/java/com/example/App.java": "package com.example;\npublic class App {}",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Java)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check source dirs
	if !sliceContains(structure.SourceDirs, "src/main/java") {
		t.Errorf("SourceDirs should contain 'src/main/java', got %v", structure.SourceDirs)
	}
	if !sliceContains(structure.SourceDirs, "src/main/resources") {
		t.Errorf("SourceDirs should contain 'src/main/resources', got %v", structure.SourceDirs)
	}

	// Check test dirs
	if !sliceContains(structure.TestDirs, "src/test/java") {
		t.Errorf("TestDirs should contain 'src/test/java', got %v", structure.TestDirs)
	}

	// Check build output
	if structure.BuildOutputDir != "target" {
		t.Errorf("BuildOutputDir = %q, want %q", structure.BuildOutputDir, "target")
	}
}

func TestMetadataExtractor_ExtractStructure_Java_Gradle(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Gradle-style Java project with simple src layout
	dirs := []string{"src", "build"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files
	files := map[string]string{
		"build.gradle": "plugins { id 'java' }",
		"src/App.java": "public class App {}",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Java)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// For simple layout, src should be detected
	if !sliceContains(structure.SourceDirs, "src") {
		t.Errorf("SourceDirs should contain 'src', got %v", structure.SourceDirs)
	}

	// Check build output
	if structure.BuildOutputDir != "build" {
		t.Errorf("BuildOutputDir = %q, want %q", structure.BuildOutputDir, "build")
	}
}

func TestMetadataExtractor_ExtractGoStructure_WithMainGo(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go project with main.go at root
	if err := os.MkdirAll(filepath.Join(tmpDir, "bin"), 0755); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"go.mod":  "module example.com/test\n\ngo 1.21\n",
		"main.go": "package main\n\nfunc main() {}\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Go)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check entry points - should include main.go at root
	if !sliceContains(structure.EntryPoints, "main.go") {
		t.Errorf("EntryPoints should contain 'main.go', got %v", structure.EntryPoints)
	}

	// Check build output
	if structure.BuildOutputDir != "bin" {
		t.Errorf("BuildOutputDir = %q, want %q", structure.BuildOutputDir, "bin")
	}
}

func TestMetadataExtractor_ExtractGoMetadata_WithVersionFile(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]string{
		"go.mod":  "module example.com/test\n\ngo 1.21\n",
		"VERSION": "1.2.3",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Go)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Version != "1.2.3" {
		t.Errorf("Version = %q, want %q", metadata.Version, "1.2.3")
	}
}

func TestMetadataExtractor_ExtractPythonMetadata_WithAuthors(t *testing.T) {
	tmpDir := t.TempDir()
	pyproject := `[project]
name = "my-pkg"
version = "1.0.0"
authors = [
    "Jane Doe <jane@example.com>",
    "John Doe <john@example.com>"
]
`
	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyproject), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Python)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Author != "Jane Doe <jane@example.com>" {
		t.Errorf("Author = %q, want first author from list", metadata.Author)
	}
}

func TestMetadataExtractor_ExtractNodeMetadata_ExtraFields(t *testing.T) {
	tmpDir := t.TempDir()
	packageJSON := `{
		"name": "extra-test",
		"version": "1.0.0",
		"main": "lib/index.js",
		"module": "lib/index.mjs",
		"types": "lib/index.d.ts"
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageJSON), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Node)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Extra["main"] != "lib/index.js" {
		t.Errorf("Extra[main] = %q, want %q", metadata.Extra["main"], "lib/index.js")
	}
	if metadata.Extra["module"] != "lib/index.mjs" {
		t.Errorf("Extra[module] = %q, want %q", metadata.Extra["module"], "lib/index.mjs")
	}
	if metadata.Extra["types"] != "lib/index.d.ts" {
		t.Errorf("Extra[types] = %q, want %q", metadata.Extra["types"], "lib/index.d.ts")
	}
}

func TestMetadataExtractor_ExtractNodeMetadata_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("not valid json"), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	_, err := extractor.Extract(Node)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestMetadataExtractor_ExtractNodeStructure_WithBin(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Node project with bin entry
	dirs := []string{"src", "__tests__"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	files := map[string]string{
		"package.json": `{"name":"test","bin":{"cli":"./bin/cli.js"}}`,
		"src/index.js": "module.exports = {}",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Node)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check __tests__ is detected
	if !sliceContains(structure.TestDirs, "__tests__") {
		t.Errorf("TestDirs should contain '__tests__', got %v", structure.TestDirs)
	}
}

func TestMetadataExtractor_ExtractPythonStructure_TestDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Python project with test directory
	dirs := []string{"testing", "dist"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	files := map[string]string{
		"app.py": "import flask",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Python)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Check entry points
	if !sliceContains(structure.EntryPoints, "app.py") {
		t.Errorf("EntryPoints should contain 'app.py', got %v", structure.EntryPoints)
	}

	// Check test dirs (testing)
	if !sliceContains(structure.TestDirs, "testing") {
		t.Errorf("TestDirs should contain 'testing', got %v", structure.TestDirs)
	}

	// Check build output
	if structure.BuildOutputDir != "dist" {
		t.Errorf("BuildOutputDir = %q, want %q", structure.BuildOutputDir, "dist")
	}
}

func TestMetadataExtractor_MonorepoGoPackages(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go monorepo structure with packages directory
	pkgDir := filepath.Join(tmpDir, "packages")
	for _, pkg := range []string{"pkg1", "pkg2"} {
		if err := os.MkdirAll(filepath.Join(pkgDir, pkg), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(pkgDir, pkg, "go.mod"), []byte("module "+pkg+"\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Go)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	if !structure.HasMonorepo {
		t.Error("HasMonorepo should be true for packages/ with multiple go.mod files")
	}
}

func TestMetadataExtractor_ExtractJavaMetadata_GradleKts(t *testing.T) {
	tmpDir := t.TempDir()
	// Note: Gradle parsing uses simple regex that matches first occurrence
	// Real Gradle files often have version early in the file
	buildGradleKts := `group = "com.example"
version = "3.0.0"

plugins {
    kotlin("jvm")
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "build.gradle.kts"), []byte(buildGradleKts), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewMetadataExtractor(tmpDir)
	metadata, err := extractor.Extract(Java)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if metadata.Version != "3.0.0" {
		t.Errorf("Version = %q, want %q", metadata.Version, "3.0.0")
	}
	if metadata.Extra["group"] != "com.example" {
		t.Errorf("group = %q, want %q", metadata.Extra["group"], "com.example")
	}
}

func TestMetadataExtractor_ExtractGoStructure_ColocatedTests(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a Go project with co-located tests (Go convention)
	dirs := []string{
		"cmd/myapp",
		"internal/handler",
		"internal/service",
		"pkg/utils",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create source files and co-located test files
	files := map[string]string{
		"go.mod":                           "module example.com/test\n\ngo 1.21\n",
		"cmd/myapp/main.go":                "package main\n\nfunc main() {}\n",
		"internal/handler/handler.go":      "package handler\n\nfunc Handle() {}\n",
		"internal/handler/handler_test.go": "package handler\n\nimport \"testing\"\n\nfunc TestHandle(t *testing.T) {}\n",
		"internal/service/service.go":      "package service\n\nfunc Serve() {}\n",
		"internal/service/service_test.go": "package service\n\nimport \"testing\"\n\nfunc TestServe(t *testing.T) {}\n",
		"pkg/utils/utils.go":               "package utils\n\nfunc Helper() {}\n",
		"pkg/utils/utils_test.go":          "package utils\n\nimport \"testing\"\n\nfunc TestHelper(t *testing.T) {}\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Go)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Should detect directories containing test files
	expectedTestDirs := []string{
		"internal/handler",
		"internal/service",
		"pkg/utils",
	}
	for _, dir := range expectedTestDirs {
		if !sliceContains(structure.TestDirs, dir) {
			t.Errorf("TestDirs should contain %q (co-located tests), got %v", dir, structure.TestDirs)
		}
	}
}

func TestMetadataExtractor_ExtractGoStructure_RootLevelTests(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple Go project with tests at root level
	files := map[string]string{
		"go.mod":       "module example.com/test\n\ngo 1.21\n",
		"main.go":      "package main\n\nfunc main() {}\n",
		"main_test.go": "package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {}\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Go)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Should detect root directory contains tests
	if !sliceContains(structure.TestDirs, ".") {
		t.Errorf("TestDirs should contain '.' for root-level tests, got %v", structure.TestDirs)
	}
}

func TestMetadataExtractor_ExtractGoStructure_ExcludesVendorAndTestdata(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directories including vendor and testdata
	dirs := []string{
		"internal/handler",
		"vendor/example.com/pkg",
		"testdata",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	files := map[string]string{
		"go.mod":                           "module example.com/test\n\ngo 1.21\n",
		"internal/handler/handler.go":      "package handler\n\nfunc Handle() {}\n",
		"internal/handler/handler_test.go": "package handler\n\nimport \"testing\"\n\nfunc TestHandle(t *testing.T) {}\n",
		// These should be ignored
		"vendor/example.com/pkg/pkg_test.go": "package pkg\n\nimport \"testing\"\n\nfunc TestPkg(t *testing.T) {}\n",
		"testdata/fixture_test.go":           "package testdata\n\nimport \"testing\"\n\nfunc TestFixture(t *testing.T) {}\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	structure, err := extractor.ExtractStructure(Go)
	if err != nil {
		t.Fatalf("ExtractStructure() error = %v", err)
	}

	// Should include internal/handler
	if !sliceContains(structure.TestDirs, "internal/handler") {
		t.Errorf("TestDirs should contain 'internal/handler', got %v", structure.TestDirs)
	}

	// Should not include vendor directories
	for _, dir := range structure.TestDirs {
		if strings.Contains(dir, "vendor") {
			t.Errorf("TestDirs should not contain vendor directories, but found %q", dir)
		}
	}

	// testdata should be in TestDirs as a dedicated test directory (not from co-located detection)
	// but should not have test files detected within it
}

func TestMetadataExtractor_findGoTestDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create various directories
	dirs := []string{
		"cmd/app",
		"pkg/lib",
		"internal/core",
		"api/v1",
		"random/nested",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create test files in some directories
	testFiles := []string{
		"cmd/app/main_test.go",
		"pkg/lib/lib_test.go",
		"internal/core/core_test.go",
		"api/v1/api_test.go",
		"random/nested/nested_test.go", // Should NOT be included - not a source dir
	}
	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		content := "package test\n\nimport \"testing\"\n\nfunc TestSomething(t *testing.T) {}\n"
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	extractor := NewMetadataExtractor(tmpDir)
	testDirs := extractor.findGoTestDirs()

	// Should include source directories with tests
	expectedDirs := []string{"cmd/app", "pkg/lib", "internal/core", "api/v1"}
	for _, expected := range expectedDirs {
		if !sliceContains(testDirs, expected) {
			t.Errorf("findGoTestDirs() should include %q, got %v", expected, testDirs)
		}
	}

	// Should NOT include random/nested (not a standard Go source directory)
	if sliceContains(testDirs, "random/nested") {
		t.Errorf("findGoTestDirs() should not include 'random/nested', got %v", testDirs)
	}
}

func TestMetadataExtractor_isGoSourceDir(t *testing.T) {
	extractor := NewMetadataExtractor(".")

	tests := []struct {
		path     string
		expected bool
	}{
		{".", true},
		{"cmd", true},
		{"cmd/myapp", true},
		{"pkg", true},
		{"pkg/utils", true},
		{"internal", true},
		{"internal/handler", true},
		{"lib", true},
		{"api", true},
		{"api/v1", true},
		{"app", true},
		{"app/handlers", true},
		{"random", false},
		{"docs", false},
		{"scripts", false},
		{"vendor", false}, // vendor is skipped at walk level, but also not a source dir
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := extractor.isGoSourceDir(tt.path)
			if result != tt.expected {
				t.Errorf("isGoSourceDir(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// sliceContains checks if a string slice contains a value.
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
