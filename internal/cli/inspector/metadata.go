// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProjectMetadata holds extracted metadata from project configuration files.
type ProjectMetadata struct {
	Name        string            // Project name
	Version     string            // Project version
	Description string            // Project description
	License     string            // License identifier
	Repository  string            // Repository URL
	Author      string            // Author name or email
	Keywords    []string          // Project keywords/tags
	Extra       map[string]string // Additional metadata fields
}

// ProjectStructure describes the layout of a project.
type ProjectStructure struct {
	EntryPoints    []string // Main entry points (e.g., cmd/main.go, src/index.js)
	SourceDirs     []string // Main source directories (e.g., src, lib, pkg)
	TestDirs       []string // Test directories (e.g., tests, test, __tests__)
	ConfigFiles    []string // Configuration files found
	HasMonorepo    bool     // Whether this appears to be a monorepo
	BuildOutputDir string   // Build output directory (e.g., dist, build, bin)
}

// MetadataExtractor extracts metadata from project configuration files.
type MetadataExtractor struct {
	root string
}

// NewMetadataExtractor creates a new MetadataExtractor for the given project root.
func NewMetadataExtractor(root string) *MetadataExtractor {
	return &MetadataExtractor{root: root}
}

// Extract extracts metadata based on the detected project type.
func (m *MetadataExtractor) Extract(projectType ProjectType) (*ProjectMetadata, error) {
	switch projectType {
	case Go:
		return m.extractGoMetadata()
	case Node:
		return m.extractNodeMetadata()
	case Python:
		return m.extractPythonMetadata()
	case Rust:
		return m.extractRustMetadata()
	case Ruby:
		return m.extractRubyMetadata()
	case Java:
		return m.extractJavaMetadata()
	default:
		return &ProjectMetadata{Extra: make(map[string]string)}, nil
	}
}

// ExtractStructure analyzes the project structure.
func (m *MetadataExtractor) ExtractStructure(projectType ProjectType) (*ProjectStructure, error) {
	structure := &ProjectStructure{
		EntryPoints:    []string{},
		SourceDirs:     []string{},
		TestDirs:       []string{},
		ConfigFiles:    []string{},
		HasMonorepo:    false,
		BuildOutputDir: "",
	}

	// Find config files that exist
	configCandidates := []string{
		"go.mod", "go.sum",
		"package.json", "package-lock.json", "yarn.lock", "pnpm-lock.yaml",
		"pyproject.toml", "setup.py", "setup.cfg", "requirements.txt",
		"Cargo.toml", "Cargo.lock",
		"Gemfile", "Gemfile.lock",
		"pom.xml", "build.gradle", "build.gradle.kts",
		".golangci.yml", ".eslintrc.json", ".prettierrc",
		"tsconfig.json", "jest.config.js", "vitest.config.ts",
		"Makefile", "Dockerfile", "docker-compose.yml",
	}

	for _, cfg := range configCandidates {
		if m.fileExists(cfg) {
			structure.ConfigFiles = append(structure.ConfigFiles, cfg)
		}
	}

	// Extract structure based on project type
	switch projectType {
	case Go:
		m.extractGoStructure(structure)
	case Node:
		m.extractNodeStructure(structure)
	case Python:
		m.extractPythonStructure(structure)
	case Rust:
		m.extractRustStructure(structure)
	case Ruby:
		m.extractRubyStructure(structure)
	case Java:
		m.extractJavaStructure(structure)
	}

	// Detect monorepo patterns
	structure.HasMonorepo = m.detectMonorepo()

	return structure, nil
}

// extractGoMetadata extracts metadata from go.mod.
func (m *MetadataExtractor) extractGoMetadata() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	goModPath := filepath.Join(m.root, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return metadata, nil // go.mod not found, return empty metadata
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	moduleRegex := regexp.MustCompile(`^module\s+(.+)$`)
	goVersionRegex := regexp.MustCompile(`^go\s+(\d+\.\d+(?:\.\d+)?)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Extract module name
		if matches := moduleRegex.FindStringSubmatch(line); len(matches) > 1 {
			metadata.Name = strings.TrimSpace(matches[1])
		}

		// Extract Go version
		if matches := goVersionRegex.FindStringSubmatch(line); len(matches) > 1 {
			metadata.Extra["go_version"] = matches[1]
		}
	}

	// Try to extract version from git tags or VERSION file
	if version := m.extractVersionFromFile("VERSION"); version != "" {
		metadata.Version = version
	}

	return metadata, scanner.Err()
}

// extractNodeMetadata extracts metadata from package.json.
func (m *MetadataExtractor) extractNodeMetadata() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	pkgJSONPath := filepath.Join(m.root, "package.json")
	data, err := os.ReadFile(pkgJSONPath)
	if err != nil {
		return metadata, nil // package.json not found
	}

	var pkg struct {
		Name        string   `json:"name"`
		Version     string   `json:"version"`
		Description string   `json:"description"`
		License     string   `json:"license"`
		Author      any      `json:"author"` // Can be string or object
		Keywords    []string `json:"keywords"`
		Repository  any      `json:"repository"` // Can be string or object
		Main        string   `json:"main"`
		Module      string   `json:"module"`
		Types       string   `json:"types"`
		Engines     struct {
			Node string `json:"node"`
			Npm  string `json:"npm"`
		} `json:"engines"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return metadata, err
	}

	metadata.Name = pkg.Name
	metadata.Version = pkg.Version
	metadata.Description = pkg.Description
	metadata.License = pkg.License
	metadata.Keywords = pkg.Keywords

	// Handle author field (can be string or object)
	switch author := pkg.Author.(type) {
	case string:
		metadata.Author = author
	case map[string]interface{}:
		if name, ok := author["name"].(string); ok {
			metadata.Author = name
			if email, ok := author["email"].(string); ok {
				metadata.Author = name + " <" + email + ">"
			}
		}
	}

	// Handle repository field (can be string or object)
	switch repo := pkg.Repository.(type) {
	case string:
		metadata.Repository = repo
	case map[string]interface{}:
		if url, ok := repo["url"].(string); ok {
			metadata.Repository = url
		}
	}

	// Store extra fields
	if pkg.Main != "" {
		metadata.Extra["main"] = pkg.Main
	}
	if pkg.Module != "" {
		metadata.Extra["module"] = pkg.Module
	}
	if pkg.Types != "" {
		metadata.Extra["types"] = pkg.Types
	}
	if pkg.Engines.Node != "" {
		metadata.Extra["node_version"] = pkg.Engines.Node
	}

	return metadata, nil
}

// extractPythonMetadata extracts metadata from pyproject.toml or setup.py.
func (m *MetadataExtractor) extractPythonMetadata() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	// Try pyproject.toml first
	if m.fileExists("pyproject.toml") {
		return m.extractPyprojectToml()
	}

	// Fall back to setup.py
	if m.fileExists("setup.py") {
		return m.extractSetupPy()
	}

	return metadata, nil
}

// extractPyprojectToml extracts metadata from pyproject.toml.
func (m *MetadataExtractor) extractPyprojectToml() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	content, err := os.ReadFile(filepath.Join(m.root, "pyproject.toml"))
	if err != nil {
		return metadata, nil
	}

	lines := strings.Split(string(content), "\n")
	inProjectSection := false
	inAuthorsSection := false

	// Simple TOML parsing for common fields
	nameRegex := regexp.MustCompile(`^name\s*=\s*"([^"]+)"`)
	versionRegex := regexp.MustCompile(`^version\s*=\s*"([^"]+)"`)
	descRegex := regexp.MustCompile(`^description\s*=\s*"([^"]+)"`)
	licenseRegex := regexp.MustCompile(`^license\s*=\s*"([^"]+)"`)
	pythonRegex := regexp.MustCompile(`^requires-python\s*=\s*"([^"]+)"`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track section
		if trimmed == "[project]" || trimmed == "[tool.poetry]" {
			inProjectSection = true
			inAuthorsSection = false
			continue
		}
		if strings.HasPrefix(trimmed, "[") {
			inProjectSection = false
			inAuthorsSection = false
			continue
		}
		if trimmed == "authors = [" {
			inAuthorsSection = true
			continue
		}
		if inAuthorsSection && trimmed == "]" {
			inAuthorsSection = false
			continue
		}

		if !inProjectSection {
			continue
		}

		// Extract fields
		if matches := nameRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Name = matches[1]
		}
		if matches := versionRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Version = matches[1]
		}
		if matches := descRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Description = matches[1]
		}
		if matches := licenseRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.License = matches[1]
		}
		if matches := pythonRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Extra["python_version"] = matches[1]
		}

		// Extract first author
		if inAuthorsSection && metadata.Author == "" {
			authorRegex := regexp.MustCompile(`"([^"]+)"`)
			if matches := authorRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
				metadata.Author = matches[1]
			}
		}
	}

	return metadata, nil
}

// extractSetupPy extracts metadata from setup.py using regex (limited parsing).
func (m *MetadataExtractor) extractSetupPy() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	content, err := os.ReadFile(filepath.Join(m.root, "setup.py"))
	if err != nil {
		return metadata, nil
	}

	text := string(content)

	// Simple regex extraction for common setup() arguments
	patterns := map[string]*regexp.Regexp{
		"name":        regexp.MustCompile(`name\s*=\s*["']([^"']+)["']`),
		"version":     regexp.MustCompile(`version\s*=\s*["']([^"']+)["']`),
		"description": regexp.MustCompile(`description\s*=\s*["']([^"']+)["']`),
		"author":      regexp.MustCompile(`author\s*=\s*["']([^"']+)["']`),
		"license":     regexp.MustCompile(`license\s*=\s*["']([^"']+)["']`),
	}

	for field, regex := range patterns {
		if matches := regex.FindStringSubmatch(text); len(matches) > 1 {
			switch field {
			case "name":
				metadata.Name = matches[1]
			case "version":
				metadata.Version = matches[1]
			case "description":
				metadata.Description = matches[1]
			case "author":
				metadata.Author = matches[1]
			case "license":
				metadata.License = matches[1]
			}
		}
	}

	return metadata, nil
}

// extractRustMetadata extracts metadata from Cargo.toml.
func (m *MetadataExtractor) extractRustMetadata() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	content, err := os.ReadFile(filepath.Join(m.root, "Cargo.toml"))
	if err != nil {
		return metadata, nil
	}

	lines := strings.Split(string(content), "\n")
	inPackageSection := false

	// Simple TOML parsing
	nameRegex := regexp.MustCompile(`^name\s*=\s*"([^"]+)"`)
	versionRegex := regexp.MustCompile(`^version\s*=\s*"([^"]+)"`)
	descRegex := regexp.MustCompile(`^description\s*=\s*"([^"]+)"`)
	licenseRegex := regexp.MustCompile(`^license\s*=\s*"([^"]+)"`)
	editionRegex := regexp.MustCompile(`^edition\s*=\s*"([^"]+)"`)
	repoRegex := regexp.MustCompile(`^repository\s*=\s*"([^"]+)"`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "[package]" {
			inPackageSection = true
			continue
		}
		if strings.HasPrefix(trimmed, "[") && trimmed != "[package]" {
			inPackageSection = false
			continue
		}

		if !inPackageSection {
			continue
		}

		if matches := nameRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Name = matches[1]
		}
		if matches := versionRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Version = matches[1]
		}
		if matches := descRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Description = matches[1]
		}
		if matches := licenseRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.License = matches[1]
		}
		if matches := editionRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Extra["rust_edition"] = matches[1]
		}
		if matches := repoRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			metadata.Repository = matches[1]
		}
	}

	return metadata, nil
}

// extractRubyMetadata extracts metadata from Gemfile or .gemspec.
func (m *MetadataExtractor) extractRubyMetadata() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	// Look for gemspec files
	entries, err := os.ReadDir(m.root)
	if err != nil {
		return metadata, nil
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".gemspec") {
			return m.extractGemspec(entry.Name())
		}
	}

	return metadata, nil
}

// extractGemspec extracts metadata from a .gemspec file.
func (m *MetadataExtractor) extractGemspec(filename string) (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	content, err := os.ReadFile(filepath.Join(m.root, filename))
	if err != nil {
		return metadata, nil
	}

	text := string(content)

	// Ruby gemspec regex patterns
	patterns := map[string]*regexp.Regexp{
		"name":        regexp.MustCompile(`\.name\s*=\s*["']([^"']+)["']`),
		"version":     regexp.MustCompile(`\.version\s*=\s*["']([^"']+)["']`),
		"description": regexp.MustCompile(`\.(?:description|summary)\s*=\s*["']([^"']+)["']`),
		"author":      regexp.MustCompile(`\.authors?\s*=\s*\[?\s*["']([^"']+)["']`),
		"license":     regexp.MustCompile(`\.license\s*=\s*["']([^"']+)["']`),
	}

	for field, regex := range patterns {
		if matches := regex.FindStringSubmatch(text); len(matches) > 1 {
			switch field {
			case "name":
				metadata.Name = matches[1]
			case "version":
				metadata.Version = matches[1]
			case "description":
				metadata.Description = matches[1]
			case "author":
				metadata.Author = matches[1]
			case "license":
				metadata.License = matches[1]
			}
		}
	}

	return metadata, nil
}

// extractJavaMetadata extracts metadata from pom.xml or build.gradle.
func (m *MetadataExtractor) extractJavaMetadata() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	// Try pom.xml first (Maven)
	if m.fileExists("pom.xml") {
		return m.extractPomXml()
	}

	// Fall back to build.gradle (Gradle)
	if m.fileExists("build.gradle") || m.fileExists("build.gradle.kts") {
		return m.extractBuildGradle()
	}

	return metadata, nil
}

// extractPomXml extracts metadata from Maven pom.xml.
func (m *MetadataExtractor) extractPomXml() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	content, err := os.ReadFile(filepath.Join(m.root, "pom.xml"))
	if err != nil {
		return metadata, nil
	}

	text := string(content)

	// Simple XML regex extraction (not a full XML parser)
	patterns := map[string]*regexp.Regexp{
		"groupId":     regexp.MustCompile(`<groupId>([^<]+)</groupId>`),
		"artifactId":  regexp.MustCompile(`<artifactId>([^<]+)</artifactId>`),
		"version":     regexp.MustCompile(`<version>([^<]+)</version>`),
		"name":        regexp.MustCompile(`<name>([^<]+)</name>`),
		"description": regexp.MustCompile(`<description>([^<]+)</description>`),
	}

	for field, regex := range patterns {
		if matches := regex.FindStringSubmatch(text); len(matches) > 1 {
			switch field {
			case "groupId":
				metadata.Extra["group_id"] = matches[1]
			case "artifactId":
				if metadata.Name == "" {
					metadata.Name = matches[1]
				}
			case "version":
				if metadata.Version == "" {
					metadata.Version = matches[1]
				}
			case "name":
				metadata.Name = matches[1]
			case "description":
				metadata.Description = matches[1]
			}
		}
	}

	return metadata, nil
}

// extractBuildGradle extracts metadata from build.gradle or build.gradle.kts.
func (m *MetadataExtractor) extractBuildGradle() (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{
		Extra: make(map[string]string),
	}

	filename := "build.gradle"
	if !m.fileExists(filename) {
		filename = "build.gradle.kts"
	}

	content, err := os.ReadFile(filepath.Join(m.root, filename))
	if err != nil {
		return metadata, nil
	}

	text := string(content)

	// Gradle patterns (both Groovy and Kotlin DSL)
	patterns := map[string]*regexp.Regexp{
		"group":       regexp.MustCompile(`group\s*=?\s*["']([^"']+)["']`),
		"version":     regexp.MustCompile(`version\s*=?\s*["']([^"']+)["']`),
		"archiveBase": regexp.MustCompile(`archivesBaseName\s*=?\s*["']([^"']+)["']`),
	}

	for field, regex := range patterns {
		if matches := regex.FindStringSubmatch(text); len(matches) > 1 {
			switch field {
			case "group":
				metadata.Extra["group"] = matches[1]
			case "version":
				metadata.Version = matches[1]
			case "archiveBase":
				metadata.Name = matches[1]
			}
		}
	}

	return metadata, nil
}

// extractGoStructure extracts Go project structure.
func (m *MetadataExtractor) extractGoStructure(s *ProjectStructure) {
	// Common Go entry points
	if m.dirExists("cmd") {
		entries, _ := os.ReadDir(filepath.Join(m.root, "cmd"))
		for _, entry := range entries {
			if entry.IsDir() {
				mainFile := filepath.Join("cmd", entry.Name(), "main.go")
				if m.fileExists(mainFile) {
					s.EntryPoints = append(s.EntryPoints, mainFile)
				}
			}
		}
	}
	if m.fileExists("main.go") {
		s.EntryPoints = append(s.EntryPoints, "main.go")
	}

	// Common Go source directories
	for _, dir := range []string{"pkg", "internal", "lib"} {
		if m.dirExists(dir) {
			s.SourceDirs = append(s.SourceDirs, dir)
		}
	}

	// Test directories
	for _, dir := range []string{"test", "tests", "testdata"} {
		if m.dirExists(dir) {
			s.TestDirs = append(s.TestDirs, dir)
		}
	}

	// Build output
	for _, dir := range []string{"bin", "dist", "build"} {
		if m.dirExists(dir) {
			s.BuildOutputDir = dir
			break
		}
	}
}

// extractNodeStructure extracts Node.js project structure.
func (m *MetadataExtractor) extractNodeStructure(s *ProjectStructure) {
	// Read package.json for entry points
	if pkgData, err := os.ReadFile(filepath.Join(m.root, "package.json")); err == nil {
		var pkg struct {
			Main   string `json:"main"`
			Module string `json:"module"`
			Bin    any    `json:"bin"`
		}
		if json.Unmarshal(pkgData, &pkg) == nil {
			if pkg.Main != "" {
				s.EntryPoints = append(s.EntryPoints, pkg.Main)
			}
			if pkg.Module != "" && pkg.Module != pkg.Main {
				s.EntryPoints = append(s.EntryPoints, pkg.Module)
			}
		}
	}

	// Common entry points
	for _, entry := range []string{"src/index.js", "src/index.ts", "index.js", "index.ts", "src/main.js", "src/main.ts"} {
		if m.fileExists(entry) && !contains(s.EntryPoints, entry) {
			s.EntryPoints = append(s.EntryPoints, entry)
		}
	}

	// Source directories
	for _, dir := range []string{"src", "lib", "app"} {
		if m.dirExists(dir) {
			s.SourceDirs = append(s.SourceDirs, dir)
		}
	}

	// Test directories
	for _, dir := range []string{"test", "tests", "__tests__", "spec"} {
		if m.dirExists(dir) {
			s.TestDirs = append(s.TestDirs, dir)
		}
	}

	// Build output
	for _, dir := range []string{"dist", "build", "out", ".next"} {
		if m.dirExists(dir) {
			s.BuildOutputDir = dir
			break
		}
	}
}

// extractPythonStructure extracts Python project structure.
func (m *MetadataExtractor) extractPythonStructure(s *ProjectStructure) {
	// Common Python entry points
	for _, entry := range []string{"main.py", "app.py", "__main__.py", "src/__main__.py"} {
		if m.fileExists(entry) {
			s.EntryPoints = append(s.EntryPoints, entry)
		}
	}

	// Source directories
	for _, dir := range []string{"src", "lib", "app"} {
		if m.dirExists(dir) {
			s.SourceDirs = append(s.SourceDirs, dir)
		}
	}

	// Look for package directories (containing __init__.py)
	entries, _ := os.ReadDir(m.root)
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") && !strings.HasPrefix(entry.Name(), "_") {
			initFile := filepath.Join(m.root, entry.Name(), "__init__.py")
			if _, err := os.Stat(initFile); err == nil {
				if !contains(s.SourceDirs, entry.Name()) {
					s.SourceDirs = append(s.SourceDirs, entry.Name())
				}
			}
		}
	}

	// Test directories
	for _, dir := range []string{"tests", "test", "testing"} {
		if m.dirExists(dir) {
			s.TestDirs = append(s.TestDirs, dir)
		}
	}

	// Build output
	for _, dir := range []string{"dist", "build", ".tox"} {
		if m.dirExists(dir) {
			s.BuildOutputDir = dir
			break
		}
	}
}

// extractRustStructure extracts Rust project structure.
func (m *MetadataExtractor) extractRustStructure(s *ProjectStructure) {
	// Standard Rust entry points
	if m.fileExists("src/main.rs") {
		s.EntryPoints = append(s.EntryPoints, "src/main.rs")
	}
	if m.fileExists("src/lib.rs") {
		s.EntryPoints = append(s.EntryPoints, "src/lib.rs")
	}

	// Source directory
	if m.dirExists("src") {
		s.SourceDirs = append(s.SourceDirs, "src")
	}

	// Test directories
	if m.dirExists("tests") {
		s.TestDirs = append(s.TestDirs, "tests")
	}

	// Build output
	s.BuildOutputDir = "target"
}

// extractRubyStructure extracts Ruby project structure.
func (m *MetadataExtractor) extractRubyStructure(s *ProjectStructure) {
	// Common Ruby entry points
	for _, entry := range []string{"app.rb", "main.rb", "bin/main"} {
		if m.fileExists(entry) {
			s.EntryPoints = append(s.EntryPoints, entry)
		}
	}

	// Source directories
	for _, dir := range []string{"lib", "app", "src"} {
		if m.dirExists(dir) {
			s.SourceDirs = append(s.SourceDirs, dir)
		}
	}

	// Test directories
	for _, dir := range []string{"test", "tests", "spec"} {
		if m.dirExists(dir) {
			s.TestDirs = append(s.TestDirs, dir)
		}
	}

	// Build output (Ruby typically doesn't have build output)
	s.BuildOutputDir = ""
}

// extractJavaStructure extracts Java project structure.
func (m *MetadataExtractor) extractJavaStructure(s *ProjectStructure) {
	// Maven standard layout
	if m.dirExists("src/main/java") {
		s.SourceDirs = append(s.SourceDirs, "src/main/java")
	}
	if m.dirExists("src/main/resources") {
		s.SourceDirs = append(s.SourceDirs, "src/main/resources")
	}

	// Test directories
	if m.dirExists("src/test/java") {
		s.TestDirs = append(s.TestDirs, "src/test/java")
	}
	if m.dirExists("src/test/resources") {
		s.TestDirs = append(s.TestDirs, "src/test/resources")
	}

	// Alternative layouts
	if m.dirExists("src") && !m.dirExists("src/main") {
		s.SourceDirs = append(s.SourceDirs, "src")
	}

	// Build output
	for _, dir := range []string{"target", "build", "out"} {
		if m.dirExists(dir) {
			s.BuildOutputDir = dir
			break
		}
	}
}

// detectMonorepo checks for common monorepo patterns.
func (m *MetadataExtractor) detectMonorepo() bool {
	// Check for workspaces in package.json
	if pkgData, err := os.ReadFile(filepath.Join(m.root, "package.json")); err == nil {
		var pkg struct {
			Workspaces any `json:"workspaces"`
		}
		if json.Unmarshal(pkgData, &pkg) == nil && pkg.Workspaces != nil {
			return true
		}
	}

	// Check for pnpm-workspace.yaml
	if m.fileExists("pnpm-workspace.yaml") {
		return true
	}

	// Check for lerna.json
	if m.fileExists("lerna.json") {
		return true
	}

	// Check for Cargo workspace
	if cargoData, err := os.ReadFile(filepath.Join(m.root, "Cargo.toml")); err == nil {
		if strings.Contains(string(cargoData), "[workspace]") {
			return true
		}
	}

	// Check for packages/ or apps/ directory with multiple package.json/go.mod files
	for _, dir := range []string{"packages", "apps", "libs", "modules"} {
		if m.dirExists(dir) {
			entries, _ := os.ReadDir(filepath.Join(m.root, dir))
			count := 0
			for _, entry := range entries {
				if entry.IsDir() {
					subdir := filepath.Join(m.root, dir, entry.Name())
					if _, err := os.Stat(filepath.Join(subdir, "package.json")); err == nil {
						count++
					} else if _, err := os.Stat(filepath.Join(subdir, "go.mod")); err == nil {
						count++
					}
				}
			}
			if count >= 2 {
				return true
			}
		}
	}

	return false
}

// extractVersionFromFile tries to read version from a VERSION file.
func (m *MetadataExtractor) extractVersionFromFile(filename string) string {
	content, err := os.ReadFile(filepath.Join(m.root, filename))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}

// fileExists checks if a file exists in the project root.
func (m *MetadataExtractor) fileExists(name string) bool {
	path := filepath.Join(m.root, name)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists in the project root.
func (m *MetadataExtractor) dirExists(name string) bool {
	path := filepath.Join(m.root, name)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// contains checks if a string slice contains a value.
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
