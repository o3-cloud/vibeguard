// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"os"
	"path/filepath"
)

// ProjectType represents the detected programming language/framework of a project.
type ProjectType string

// Supported project types.
const (
	Go      ProjectType = "go"
	Node    ProjectType = "node"
	Python  ProjectType = "python"
	Ruby    ProjectType = "ruby"
	Rust    ProjectType = "rust"
	Java    ProjectType = "java"
	Unknown ProjectType = "unknown"
)

// DetectionResult holds the result of project type detection.
type DetectionResult struct {
	Type       ProjectType // Detected project type
	Confidence float64     // Confidence score 0.0-1.0
	Indicators []string    // Files/patterns that led to this detection
}

// Detector detects project types based on file patterns and project structure.
type Detector struct {
	root string
}

// NewDetector creates a new Detector for the given project root directory.
func NewDetector(root string) *Detector {
	return &Detector{root: root}
}

// Detect analyzes the project and returns detection results for all detected types.
// Returns results sorted by confidence (highest first).
func (d *Detector) Detect() ([]DetectionResult, error) {
	var results []DetectionResult

	// Check each project type
	detectors := []func() (*DetectionResult, error){
		d.detectGo,
		d.detectNode,
		d.detectPython,
		d.detectRuby,
		d.detectRust,
		d.detectJava,
	}

	for _, detect := range detectors {
		result, err := detect()
		if err != nil {
			return nil, err
		}
		if result != nil && result.Confidence > 0 {
			results = append(results, *result)
		}
	}

	// Sort by confidence (highest first)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Confidence > results[i].Confidence {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// If no results, return Unknown
	if len(results) == 0 {
		results = append(results, DetectionResult{
			Type:       Unknown,
			Confidence: 0,
			Indicators: []string{},
		})
	}

	return results, nil
}

// DetectPrimary returns the most likely project type based on highest confidence.
func (d *Detector) DetectPrimary() (*DetectionResult, error) {
	results, err := d.Detect()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &DetectionResult{Type: Unknown, Confidence: 0}, nil
	}
	return &results[0], nil
}

// detectGo checks for Go project indicators.
func (d *Detector) detectGo() (*DetectionResult, error) {
	result := &DetectionResult{
		Type:       Go,
		Confidence: 0,
		Indicators: []string{},
	}

	// Check for go.mod (strongest indicator - 0.6)
	if d.fileExists("go.mod") {
		result.Confidence += 0.6
		result.Indicators = append(result.Indicators, "go.mod")
	}

	// Check for go.sum (0.2)
	if d.fileExists("go.sum") {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "go.sum")
	}

	// Check for .go files (0.2 if any found)
	goFiles, err := d.findFiles("*.go", 3)
	if err != nil {
		return nil, err
	}
	if len(goFiles) > 0 {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "*.go files")
	}

	// Cap confidence at 1.0
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}

	return result, nil
}

// detectNode checks for Node.js project indicators.
func (d *Detector) detectNode() (*DetectionResult, error) {
	result := &DetectionResult{
		Type:       Node,
		Confidence: 0,
		Indicators: []string{},
	}

	// Check for package.json (strongest indicator - 0.6)
	if d.fileExists("package.json") {
		result.Confidence += 0.6
		result.Indicators = append(result.Indicators, "package.json")
	}

	// Check for package-lock.json or yarn.lock (0.2)
	if d.fileExists("package-lock.json") {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "package-lock.json")
	} else if d.fileExists("yarn.lock") {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "yarn.lock")
	} else if d.fileExists("pnpm-lock.yaml") {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "pnpm-lock.yaml")
	}

	// Check for node_modules directory (0.2)
	if d.dirExists("node_modules") {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "node_modules/")
	}

	// Cap confidence at 1.0
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}

	return result, nil
}

// detectPython checks for Python project indicators.
func (d *Detector) detectPython() (*DetectionResult, error) {
	result := &DetectionResult{
		Type:       Python,
		Confidence: 0,
		Indicators: []string{},
	}

	// Check for pyproject.toml (strongest indicator - 0.5)
	if d.fileExists("pyproject.toml") {
		result.Confidence += 0.5
		result.Indicators = append(result.Indicators, "pyproject.toml")
	}

	// Check for setup.py (0.4)
	if d.fileExists("setup.py") {
		result.Confidence += 0.4
		result.Indicators = append(result.Indicators, "setup.py")
	}

	// Check for requirements.txt (0.3)
	if d.fileExists("requirements.txt") {
		result.Confidence += 0.3
		result.Indicators = append(result.Indicators, "requirements.txt")
	}

	// Check for Pipfile (0.3)
	if d.fileExists("Pipfile") {
		result.Confidence += 0.3
		result.Indicators = append(result.Indicators, "Pipfile")
	}

	// Check for .py files (0.2 if any found)
	pyFiles, err := d.findFiles("*.py", 3)
	if err != nil {
		return nil, err
	}
	if len(pyFiles) > 0 {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "*.py files")
	}

	// Cap confidence at 1.0
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}

	return result, nil
}

// detectRuby checks for Ruby project indicators.
func (d *Detector) detectRuby() (*DetectionResult, error) {
	result := &DetectionResult{
		Type:       Ruby,
		Confidence: 0,
		Indicators: []string{},
	}

	// Check for Gemfile (strongest indicator - 0.6)
	if d.fileExists("Gemfile") {
		result.Confidence += 0.6
		result.Indicators = append(result.Indicators, "Gemfile")
	}

	// Check for Gemfile.lock (0.2)
	if d.fileExists("Gemfile.lock") {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "Gemfile.lock")
	}

	// Check for .rb files (0.2 if any found)
	rbFiles, err := d.findFiles("*.rb", 3)
	if err != nil {
		return nil, err
	}
	if len(rbFiles) > 0 {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "*.rb files")
	}

	// Cap confidence at 1.0
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}

	return result, nil
}

// detectRust checks for Rust project indicators.
func (d *Detector) detectRust() (*DetectionResult, error) {
	result := &DetectionResult{
		Type:       Rust,
		Confidence: 0,
		Indicators: []string{},
	}

	// Check for Cargo.toml (strongest indicator - 0.7)
	if d.fileExists("Cargo.toml") {
		result.Confidence += 0.7
		result.Indicators = append(result.Indicators, "Cargo.toml")
	}

	// Check for Cargo.lock (0.2)
	if d.fileExists("Cargo.lock") {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "Cargo.lock")
	}

	// Check for .rs files (0.1 if any found)
	rsFiles, err := d.findFiles("*.rs", 3)
	if err != nil {
		return nil, err
	}
	if len(rsFiles) > 0 {
		result.Confidence += 0.1
		result.Indicators = append(result.Indicators, "*.rs files")
	}

	// Cap confidence at 1.0
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}

	return result, nil
}

// detectJava checks for Java project indicators.
func (d *Detector) detectJava() (*DetectionResult, error) {
	result := &DetectionResult{
		Type:       Java,
		Confidence: 0,
		Indicators: []string{},
	}

	// Check for pom.xml (Maven - strongest indicator - 0.6)
	if d.fileExists("pom.xml") {
		result.Confidence += 0.6
		result.Indicators = append(result.Indicators, "pom.xml")
	}

	// Check for build.gradle or build.gradle.kts (Gradle - 0.6)
	if d.fileExists("build.gradle") {
		result.Confidence += 0.6
		result.Indicators = append(result.Indicators, "build.gradle")
	} else if d.fileExists("build.gradle.kts") {
		result.Confidence += 0.6
		result.Indicators = append(result.Indicators, "build.gradle.kts")
	}

	// Check for .java files (0.2 if any found)
	// Use depth 5 for Java since standard Maven/Gradle layout is src/main/java/...
	javaFiles, err := d.findFiles("*.java", 5)
	if err != nil {
		return nil, err
	}
	if len(javaFiles) > 0 {
		result.Confidence += 0.2
		result.Indicators = append(result.Indicators, "*.java files")
	}

	// Cap confidence at 1.0
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}

	return result, nil
}

// fileExists checks if a file exists in the project root.
func (d *Detector) fileExists(name string) bool {
	path := filepath.Join(d.root, name)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists in the project root.
func (d *Detector) dirExists(name string) bool {
	path := filepath.Join(d.root, name)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// findFiles searches for files matching the pattern in the project.
// It limits the search to avoid scanning large directories.
// maxDepth limits how deep to recurse (0 = root only, -1 = unlimited).
func (d *Detector) findFiles(pattern string, maxDepth int) ([]string, error) {
	var matches []string
	maxResults := 10 // Limit results to avoid scanning entire codebase

	err := filepath.WalkDir(d.root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}

		// Calculate depth
		relPath, _ := filepath.Rel(d.root, path)
		depth := 0
		if relPath != "." {
			depth = len(filepath.SplitList(relPath))
			// Count path separators for depth
			for _, c := range relPath {
				if c == filepath.Separator {
					depth++
				}
			}
			depth++ // Add 1 for the file itself
		}

		// Skip if too deep
		if maxDepth >= 0 && depth > maxDepth {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip common non-source directories
		if entry.IsDir() {
			name := entry.Name()
			if name == "node_modules" || name == "vendor" || name == ".git" ||
				name == "__pycache__" || name == ".venv" || name == "venv" ||
				name == "target" || name == "build" || name == "dist" {
				return filepath.SkipDir
			}
		}

		// Check if file matches pattern
		if !entry.IsDir() {
			matched, err := filepath.Match(pattern, entry.Name())
			if err != nil {
				return nil
			}
			if matched {
				matches = append(matches, path)
				if len(matches) >= maxResults {
					return filepath.SkipAll
				}
			}
		}

		return nil
	})

	return matches, err
}
