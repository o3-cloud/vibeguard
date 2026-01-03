// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

// CheckRecommendation represents a suggested check based on detected tools.
type CheckRecommendation struct {
	ID          string   // Unique check identifier
	Description string   // Human-readable description of what the check does
	Rationale   string   // Why this check is recommended
	Command     string   // Shell command to execute
	File        string   // File to read output from instead of command stdout
	Grok        []string // Optional grok patterns for output extraction
	Assert      string   // Optional assertion expression
	Severity    string   // "error" or "warning"
	Suggestion  string   // Guidance shown on failure (supports templating)
	Requires    []string // Dependencies on other check IDs
	Timeout     string   // Optional timeout duration
	Category    string   // Category: lint, format, test, build, security, etc.
	Tool        string   // The tool this check uses
	Priority    int      // Ordering priority (lower = higher priority)
}

// Recommender generates check recommendations based on project analysis.
type Recommender struct {
	projectType ProjectType
	tools       []ToolInfo
}

// NewRecommender creates a new Recommender with detected project info.
func NewRecommender(projectType ProjectType, tools []ToolInfo) *Recommender {
	return &Recommender{
		projectType: projectType,
		tools:       tools,
	}
}

// Recommend generates check recommendations based on detected tools.
func (r *Recommender) Recommend() []CheckRecommendation {
	var recommendations []CheckRecommendation

	// Get recommendations for each detected tool
	for _, tool := range r.tools {
		if !tool.Detected {
			continue
		}

		toolRecs := r.recommendationsForTool(tool)
		recommendations = append(recommendations, toolRecs...)
	}

	// Add project-type specific recommendations that aren't tool-based
	projectRecs := r.projectTypeRecommendations()
	recommendations = append(recommendations, projectRecs...)

	// Sort by priority
	sortRecommendations(recommendations)

	return recommendations
}

// RecommendForCategory returns recommendations filtered by category.
func (r *Recommender) RecommendForCategory(category string) []CheckRecommendation {
	all := r.Recommend()
	var filtered []CheckRecommendation
	for _, rec := range all {
		if rec.Category == category {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}

// recommendationsForTool returns recommendations for a specific tool.
func (r *Recommender) recommendationsForTool(tool ToolInfo) []CheckRecommendation {
	switch tool.Name {
	// Go tools
	case "golangci-lint":
		return r.golangciLintRecommendations(tool)
	case "gofmt":
		return r.gofmtRecommendations(tool)
	case "go vet":
		return r.goVetRecommendations(tool)
	case "go test":
		return r.goTestRecommendations(tool)
	case "goimports":
		return r.goimportsRecommendations(tool)

	// Node.js tools
	case "eslint":
		return r.eslintRecommendations(tool)
	case "prettier":
		return r.prettierRecommendations(tool)
	case "jest":
		return r.jestRecommendations(tool)
	case "mocha":
		return r.mochaRecommendations(tool)
	case "vitest":
		return r.vitestRecommendations(tool)
	case "typescript":
		return r.typescriptRecommendations(tool)
	case "npm audit":
		return r.npmAuditRecommendations(tool)

	// Python tools
	case "black":
		return r.blackRecommendations(tool)
	case "pylint":
		return r.pylintRecommendations(tool)
	case "pytest":
		return r.pytestRecommendations(tool)
	case "mypy":
		return r.mypyRecommendations(tool)
	case "ruff":
		return r.ruffRecommendations(tool)
	case "flake8":
		return r.flake8Recommendations(tool)
	case "isort":
		return r.isortRecommendations(tool)
	case "pip-audit":
		return r.pipAuditRecommendations(tool)

	// Git hooks
	case "pre-commit":
		return r.precommitRecommendations(tool)
	case "husky":
		return r.huskyRecommendations(tool)
	case "lefthook":
		return r.lefthookRecommendations(tool)

	default:
		return nil
	}
}

// projectTypeRecommendations returns project-type specific recommendations.
func (r *Recommender) projectTypeRecommendations() []CheckRecommendation {
	switch r.projectType {
	case Go:
		return r.goProjectRecommendations()
	case Node:
		return r.nodeProjectRecommendations()
	case Python:
		return r.pythonProjectRecommendations()
	default:
		return nil
	}
}

// Go tool recommendations

func (r *Recommender) golangciLintRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "lint",
			Description: "Run golangci-lint to check for code quality issues",
			Rationale:   "golangci-lint aggregates multiple Go linters and provides comprehensive code analysis",
			Command:     "golangci-lint run ./...",
			Severity:    "error",
			Suggestion:  "Fix the linting issues reported above. Run 'golangci-lint run --fix' to auto-fix some issues.",
			Category:    "lint",
			Tool:        "golangci-lint",
			Priority:    20,
		},
	}
}

func (r *Recommender) gofmtRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "fmt",
			Description: "Check Go code formatting with gofmt",
			Rationale:   "Consistent formatting improves readability and reduces diffs",
			Command:     "test -z \"$(gofmt -l .)\"",
			Severity:    "error",
			Suggestion:  "Run 'gofmt -w .' to format your Go code.",
			Category:    "format",
			Tool:        "gofmt",
			Priority:    10,
		},
	}
}

func (r *Recommender) goVetRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "vet",
			Description: "Run go vet to detect suspicious constructs",
			Rationale:   "go vet finds bugs that the compiler doesn't catch, like incorrect printf format strings",
			Command:     "go vet ./...",
			Severity:    "error",
			Suggestion:  "Fix the issues reported by go vet. These are often real bugs.",
			Category:    "lint",
			Tool:        "go vet",
			Priority:    15,
		},
	}
}

func (r *Recommender) goTestRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "test",
			Description: "Run Go tests",
			Rationale:   "Tests verify that code behaves as expected",
			Command:     "go test ./...",
			Severity:    "error",
			Suggestion:  "Fix failing tests before committing.",
			Category:    "test",
			Tool:        "go test",
			Priority:    30,
		},
		{
			ID:          "coverage",
			Description: "Check test coverage meets minimum threshold",
			Rationale:   "Code coverage helps identify untested code paths",
			Command:     "go test -cover ./... 2>&1 | tail -1",
			Grok:        []string{"coverage: %{NUMBER:coverage}%"},
			Assert:      "coverage >= 70",
			Severity:    "warning",
			Suggestion:  "Coverage is {{.coverage}}%, target is 70%. Add tests to improve coverage.",
			Requires:    []string{"test"},
			Category:    "test",
			Tool:        "go test",
			Priority:    35,
		},
	}
}

func (r *Recommender) goimportsRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "imports",
			Description: "Check import organization with goimports",
			Rationale:   "goimports ensures consistent import grouping and removes unused imports",
			Command:     "test -z \"$(goimports -l .)\"",
			Severity:    "error",
			Suggestion:  "Run 'goimports -w .' to fix import organization.",
			Category:    "format",
			Tool:        "goimports",
			Priority:    11,
		},
	}
}

// Node.js tool recommendations

func (r *Recommender) eslintRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "lint",
			Description: "Run ESLint to check for code quality issues",
			Rationale:   "ESLint catches common JavaScript/TypeScript errors and enforces coding standards",
			Command:     "npx eslint .",
			Severity:    "error",
			Suggestion:  "Fix ESLint errors. Run 'npx eslint . --fix' to auto-fix some issues.",
			Category:    "lint",
			Tool:        "eslint",
			Priority:    20,
		},
	}
}

func (r *Recommender) prettierRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "fmt",
			Description: "Check code formatting with Prettier",
			Rationale:   "Consistent formatting improves readability and reduces diffs",
			Command:     "npx prettier --check .",
			Severity:    "error",
			Suggestion:  "Run 'npx prettier --write .' to format your code.",
			Category:    "format",
			Tool:        "prettier",
			Priority:    10,
		},
	}
}

func (r *Recommender) jestRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "test",
			Description: "Run Jest tests",
			Rationale:   "Tests verify that code behaves as expected",
			Command:     "npx jest",
			Severity:    "error",
			Suggestion:  "Fix failing tests before committing.",
			Category:    "test",
			Tool:        "jest",
			Priority:    30,
		},
		{
			ID:          "coverage",
			Description: "Check test coverage meets minimum threshold",
			Rationale:   "Code coverage helps identify untested code paths",
			Command:     "npx jest --coverage --coverageReporters=text-summary 2>&1 | grep 'Lines'",
			Grok:        []string{"Lines\\s*:\\s*%{NUMBER:coverage}%"},
			Assert:      "coverage >= 70",
			Severity:    "warning",
			Suggestion:  "Coverage is {{.coverage}}%, target is 70%. Add tests to improve coverage.",
			Requires:    []string{"test"},
			Category:    "test",
			Tool:        "jest",
			Priority:    35,
		},
	}
}

func (r *Recommender) mochaRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "test",
			Description: "Run Mocha tests",
			Rationale:   "Tests verify that code behaves as expected",
			Command:     "npx mocha",
			Severity:    "error",
			Suggestion:  "Fix failing tests before committing.",
			Category:    "test",
			Tool:        "mocha",
			Priority:    30,
		},
	}
}

func (r *Recommender) vitestRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "test",
			Description: "Run Vitest tests",
			Rationale:   "Tests verify that code behaves as expected",
			Command:     "npx vitest run",
			Severity:    "error",
			Suggestion:  "Fix failing tests before committing.",
			Category:    "test",
			Tool:        "vitest",
			Priority:    30,
		},
	}
}

func (r *Recommender) typescriptRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "typecheck",
			Description: "Run TypeScript type checking",
			Rationale:   "TypeScript catches type errors at compile time",
			Command:     "npx tsc --noEmit",
			Severity:    "error",
			Suggestion:  "Fix TypeScript type errors before committing.",
			Category:    "typecheck",
			Tool:        "typescript",
			Priority:    25,
		},
	}
}

func (r *Recommender) npmAuditRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "security",
			Description: "Check for known security vulnerabilities in dependencies",
			Rationale:   "npm audit identifies packages with known security issues",
			Command:     "npm audit --audit-level=high",
			Severity:    "warning",
			Suggestion:  "Run 'npm audit fix' to attempt automatic fixes, or manually update vulnerable packages.",
			Category:    "security",
			Tool:        "npm audit",
			Priority:    50,
		},
	}
}

// Python tool recommendations

func (r *Recommender) blackRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "fmt",
			Description: "Check Python code formatting with Black",
			Rationale:   "Black provides uncompromising, consistent code formatting",
			Command:     "black --check .",
			Severity:    "error",
			Suggestion:  "Run 'black .' to format your Python code.",
			Category:    "format",
			Tool:        "black",
			Priority:    10,
		},
	}
}

func (r *Recommender) pylintRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "lint",
			Description: "Run Pylint to check for code quality issues",
			Rationale:   "Pylint catches bugs, code smells, and style issues in Python code",
			Command:     "pylint **/*.py",
			Severity:    "error",
			Suggestion:  "Fix Pylint errors and warnings. Consider adjusting .pylintrc for project-specific rules.",
			Category:    "lint",
			Tool:        "pylint",
			Priority:    20,
		},
	}
}

func (r *Recommender) pytestRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "test",
			Description: "Run pytest tests",
			Rationale:   "Tests verify that code behaves as expected",
			Command:     "pytest",
			Severity:    "error",
			Suggestion:  "Fix failing tests before committing.",
			Category:    "test",
			Tool:        "pytest",
			Priority:    30,
		},
		{
			ID:          "coverage",
			Description: "Check test coverage meets minimum threshold",
			Rationale:   "Code coverage helps identify untested code paths",
			Command:     "pytest --cov --cov-report=term-missing 2>&1 | grep 'TOTAL'",
			Grok:        []string{"TOTAL\\s+\\d+\\s+\\d+\\s+%{NUMBER:coverage}%"},
			Assert:      "coverage >= 70",
			Severity:    "warning",
			Suggestion:  "Coverage is {{.coverage}}%, target is 70%. Add tests to improve coverage.",
			Requires:    []string{"test"},
			Category:    "test",
			Tool:        "pytest",
			Priority:    35,
		},
	}
}

func (r *Recommender) mypyRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "typecheck",
			Description: "Run mypy for static type checking",
			Rationale:   "mypy catches type errors before runtime",
			Command:     "mypy .",
			Severity:    "error",
			Suggestion:  "Fix mypy type errors. Consider adding type hints to improve coverage.",
			Category:    "typecheck",
			Tool:        "mypy",
			Priority:    25,
		},
	}
}

func (r *Recommender) ruffRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "lint",
			Description: "Run Ruff for fast Python linting",
			Rationale:   "Ruff is an extremely fast Python linter that replaces multiple tools",
			Command:     "ruff check .",
			Severity:    "error",
			Suggestion:  "Fix Ruff errors. Run 'ruff check --fix .' to auto-fix some issues.",
			Category:    "lint",
			Tool:        "ruff",
			Priority:    20,
		},
	}
}

func (r *Recommender) flake8Recommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "lint",
			Description: "Run Flake8 for Python style checking",
			Rationale:   "Flake8 enforces PEP 8 style guidelines and catches common issues",
			Command:     "flake8 .",
			Severity:    "error",
			Suggestion:  "Fix Flake8 errors to comply with Python style guidelines.",
			Category:    "lint",
			Tool:        "flake8",
			Priority:    20,
		},
	}
}

func (r *Recommender) isortRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "imports",
			Description: "Check Python import sorting with isort",
			Rationale:   "isort ensures consistent import organization and grouping",
			Command:     "isort --check-only --diff .",
			Severity:    "error",
			Suggestion:  "Run 'isort .' to fix import ordering.",
			Category:    "format",
			Tool:        "isort",
			Priority:    11,
		},
	}
}

func (r *Recommender) pipAuditRecommendations(tool ToolInfo) []CheckRecommendation {
	return []CheckRecommendation{
		{
			ID:          "security",
			Description: "Check for known security vulnerabilities in Python dependencies",
			Rationale:   "pip-audit identifies packages with known security issues from PyPI advisory database",
			Command:     "pip-audit",
			Severity:    "warning",
			Suggestion:  "Update vulnerable packages or review security advisories. Run 'pip-audit --fix' to attempt automatic fixes.",
			Category:    "security",
			Tool:        "pip-audit",
			Priority:    50,
		},
	}
}

// Git hooks tool recommendations (minimal - these are usually run manually)

func (r *Recommender) precommitRecommendations(tool ToolInfo) []CheckRecommendation {
	// pre-commit is itself a hook runner, so we don't recommend additional checks
	return nil
}

func (r *Recommender) huskyRecommendations(tool ToolInfo) []CheckRecommendation {
	// husky is a hook manager, not a check tool
	return nil
}

func (r *Recommender) lefthookRecommendations(tool ToolInfo) []CheckRecommendation {
	// lefthook is a hook manager, not a check tool
	return nil
}

// Project type specific recommendations (for tools not explicitly detected)

func (r *Recommender) goProjectRecommendations() []CheckRecommendation {
	var recs []CheckRecommendation

	// Check if we already have these tools detected
	hasGoBuild := false
	for _, tool := range r.tools {
		if tool.Name == "go build" && tool.Detected {
			hasGoBuild = true
		}
	}

	// Add build check if not already covered
	if !hasGoBuild {
		recs = append(recs, CheckRecommendation{
			ID:          "build",
			Description: "Verify Go code compiles successfully",
			Rationale:   "Catch compilation errors before they reach CI",
			Command:     "go build ./...",
			Severity:    "error",
			Suggestion:  "Fix compilation errors before committing.",
			Category:    "build",
			Tool:        "go build",
			Priority:    5,
		})
	}

	return recs
}

func (r *Recommender) nodeProjectRecommendations() []CheckRecommendation {
	var recs []CheckRecommendation

	// Check if we already have a build step
	hasBuild := false
	for _, tool := range r.tools {
		if (tool.Name == "typescript" || tool.Name == "webpack" || tool.Name == "vite") && tool.Detected {
			hasBuild = true
		}
	}

	// Suggest build check if there's a build script but no explicit build tool detected
	if !hasBuild {
		recs = append(recs, CheckRecommendation{
			ID:          "build",
			Description: "Run npm build script",
			Rationale:   "Verify the project builds successfully",
			Command:     "npm run build --if-present",
			Severity:    "error",
			Suggestion:  "Fix build errors before committing.",
			Category:    "build",
			Tool:        "npm",
			Priority:    40,
		})
	}

	return recs
}

func (r *Recommender) pythonProjectRecommendations() []CheckRecommendation {
	// No additional project-wide recommendations beyond detected tools
	return nil
}

// sortRecommendations sorts recommendations by priority (lower = higher priority).
func sortRecommendations(recs []CheckRecommendation) {
	// Simple bubble sort for small lists
	for i := 0; i < len(recs); i++ {
		for j := i + 1; j < len(recs); j++ {
			if recs[j].Priority < recs[i].Priority {
				recs[i], recs[j] = recs[j], recs[i]
			}
		}
	}
}

// DeduplicateRecommendations removes duplicate recommendations by ID,
// keeping the first occurrence of each ID.
func DeduplicateRecommendations(recs []CheckRecommendation) []CheckRecommendation {
	seen := make(map[string]bool)
	var result []CheckRecommendation

	for _, rec := range recs {
		if !seen[rec.ID] {
			seen[rec.ID] = true
			result = append(result, rec)
		}
	}

	return result
}

// FilterByTools returns only recommendations that use the specified tools.
func FilterByTools(recs []CheckRecommendation, tools []string) []CheckRecommendation {
	toolSet := make(map[string]bool)
	for _, t := range tools {
		toolSet[t] = true
	}

	var filtered []CheckRecommendation
	for _, rec := range recs {
		if toolSet[rec.Tool] {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}

// GroupByCategory groups recommendations by their category.
func GroupByCategory(recs []CheckRecommendation) map[string][]CheckRecommendation {
	groups := make(map[string][]CheckRecommendation)
	for _, rec := range recs {
		groups[rec.Category] = append(groups[rec.Category], rec)
	}
	return groups
}
