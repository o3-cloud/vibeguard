package inspector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToolScanner_ScanGoTools(t *testing.T) {
	// Create temp directory with Go project structure
	tmpDir := t.TempDir()

	// Write go.mod
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Write golangci-lint config
	if err := os.WriteFile(filepath.Join(tmpDir, ".golangci.yml"), []byte("linters:\n  enable:\n    - gofmt\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGoTools()
	if err != nil {
		t.Fatalf("scanGoTools failed: %v", err)
	}

	// Check golangci-lint detected
	var golangciLint *ToolInfo
	var gofmt *ToolInfo
	var govet *ToolInfo
	var gotest *ToolInfo

	for i := range tools {
		switch tools[i].Name {
		case "golangci-lint":
			golangciLint = &tools[i]
		case "gofmt":
			gofmt = &tools[i]
		case "go vet":
			govet = &tools[i]
		case "go test":
			gotest = &tools[i]
		}
	}

	if golangciLint == nil || !golangciLint.Detected {
		t.Error("golangci-lint should be detected")
	} else {
		if golangciLint.ConfigFile != ".golangci.yml" {
			t.Errorf("golangci-lint config file should be .golangci.yml, got %s", golangciLint.ConfigFile)
		}
		if golangciLint.Confidence < 0.8 {
			t.Errorf("golangci-lint confidence should be >= 0.8, got %f", golangciLint.Confidence)
		}
	}

	if gofmt == nil || !gofmt.Detected {
		t.Error("gofmt should be detected (included with Go)")
	}

	if govet == nil || !govet.Detected {
		t.Error("go vet should be detected (included with Go)")
	}

	if gotest == nil || !gotest.Detected {
		t.Error("go test should be detected (included with Go)")
	}
}

func TestToolScanner_ScanGoTools_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGoTools()
	if err != nil {
		t.Fatalf("scanGoTools failed: %v", err)
	}

	// Go builtin tools should not be detected without go.mod
	for _, tool := range tools {
		if tool.Name == "gofmt" && tool.Detected {
			t.Error("gofmt should not be detected without go.mod")
		}
		if tool.Name == "go vet" && tool.Detected {
			t.Error("go vet should not be detected without go.mod")
		}
		if tool.Name == "go test" && tool.Detected {
			t.Error("go test should not be detected without go.mod")
		}
	}
}

func TestToolScanner_ScanNodeTools(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with devDependencies
	pkgJSON := `{
  "name": "test-project",
  "version": "1.0.0",
  "devDependencies": {
    "eslint": "^8.0.0",
    "prettier": "^3.0.0",
    "jest": "^29.0.0"
  }
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// Write ESLint config
	if err := os.WriteFile(filepath.Join(tmpDir, ".eslintrc.json"), []byte(`{"extends": "eslint:recommended"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Write Jest config
	if err := os.WriteFile(filepath.Join(tmpDir, "jest.config.js"), []byte("module.exports = {}"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	toolMap := make(map[string]*ToolInfo)
	for i := range tools {
		toolMap[tools[i].Name] = &tools[i]
	}

	// Check ESLint
	if eslint, ok := toolMap["eslint"]; !ok || !eslint.Detected {
		t.Error("eslint should be detected")
	} else {
		if eslint.ConfigFile != ".eslintrc.json" {
			t.Errorf("eslint config should be .eslintrc.json, got %s", eslint.ConfigFile)
		}
	}

	// Check Prettier (in devDeps but no config file)
	if prettier, ok := toolMap["prettier"]; !ok || !prettier.Detected {
		t.Error("prettier should be detected from package.json")
	}

	// Check Jest
	if jest, ok := toolMap["jest"]; !ok || !jest.Detected {
		t.Error("jest should be detected")
	} else {
		if jest.ConfigFile != "jest.config.js" {
			t.Errorf("jest config should be jest.config.js, got %s", jest.ConfigFile)
		}
	}

	// Check npm audit (always available with package.json)
	if npmAudit, ok := toolMap["npm audit"]; !ok || !npmAudit.Detected {
		t.Error("npm audit should be detected with package.json")
	}
}

func TestToolScanner_ScanNodeTools_TypeScript(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json
	pkgJSON := `{"name": "test", "devDependencies": {"typescript": "^5.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// Write tsconfig.json
	if err := os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte(`{"compilerOptions": {}}`), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	var typescript *ToolInfo
	for i := range tools {
		if tools[i].Name == "typescript" {
			typescript = &tools[i]
			break
		}
	}

	if typescript == nil || !typescript.Detected {
		t.Error("typescript should be detected")
	} else {
		if typescript.ConfigFile != "tsconfig.json" {
			t.Errorf("typescript config should be tsconfig.json, got %s", typescript.ConfigFile)
		}
		if typescript.Confidence < 0.9 {
			t.Errorf("typescript confidence should be >= 0.9 with tsconfig.json, got %f", typescript.Confidence)
		}
	}
}

func TestToolScanner_ScanPythonTools(t *testing.T) {
	tmpDir := t.TempDir()

	// Write pyproject.toml with tool configs
	pyproject := `[project]
name = "test"
version = "1.0.0"

[tool.black]
line-length = 88

[tool.mypy]
strict = true

[tool.pytest.ini_options]
testpaths = ["tests"]
`
	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyproject), 0644); err != nil {
		t.Fatal(err)
	}

	// Write requirements-dev.txt with pylint
	if err := os.WriteFile(filepath.Join(tmpDir, "requirements-dev.txt"), []byte("pylint==2.17.0\npytest==7.0.0\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	toolMap := make(map[string]*ToolInfo)
	for i := range tools {
		toolMap[tools[i].Name] = &tools[i]
	}

	// Check Black
	if black, ok := toolMap["black"]; !ok || !black.Detected {
		t.Error("black should be detected from pyproject.toml")
	} else {
		if black.ConfigFile != "pyproject.toml" {
			t.Errorf("black config should be pyproject.toml, got %s", black.ConfigFile)
		}
	}

	// Check Mypy
	if mypy, ok := toolMap["mypy"]; !ok || !mypy.Detected {
		t.Error("mypy should be detected from pyproject.toml")
	}

	// Check Pylint (from requirements-dev.txt)
	if pylint, ok := toolMap["pylint"]; !ok || !pylint.Detected {
		t.Error("pylint should be detected from requirements-dev.txt")
	}

	// Check Pytest
	if pytest, ok := toolMap["pytest"]; !ok || !pytest.Detected {
		t.Error("pytest should be detected")
	}
}

func TestToolScanner_ScanPythonTools_Ruff(t *testing.T) {
	tmpDir := t.TempDir()

	// Write ruff.toml
	if err := os.WriteFile(filepath.Join(tmpDir, "ruff.toml"), []byte("line-length = 88\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var ruff *ToolInfo
	for i := range tools {
		if tools[i].Name == "ruff" {
			ruff = &tools[i]
			break
		}
	}

	if ruff == nil || !ruff.Detected {
		t.Error("ruff should be detected")
	} else {
		if ruff.ConfigFile != "ruff.toml" {
			t.Errorf("ruff config should be ruff.toml, got %s", ruff.ConfigFile)
		}
	}
}

func TestToolScanner_ScanCITools_GitHubActions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .github/workflows directory
	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write workflow file
	workflow := `name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`
	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte(workflow), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanCITools()
	if err != nil {
		t.Fatalf("scanCITools failed: %v", err)
	}

	var githubActions *ToolInfo
	for i := range tools {
		if tools[i].Name == "GitHub Actions" {
			githubActions = &tools[i]
			break
		}
	}

	if githubActions == nil || !githubActions.Detected {
		t.Error("GitHub Actions should be detected")
	} else {
		if githubActions.ConfigFile != ".github/workflows/" {
			t.Errorf("GitHub Actions config should be .github/workflows/, got %s", githubActions.ConfigFile)
		}
		if len(githubActions.Indicators) == 0 {
			t.Error("GitHub Actions should have indicators")
		}
	}
}

func TestToolScanner_ScanCITools_GitLabCI(t *testing.T) {
	tmpDir := t.TempDir()

	// Write .gitlab-ci.yml
	gitlabCI := `stages:
  - test
test:
  script:
    - go test ./...
`
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitlab-ci.yml"), []byte(gitlabCI), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanCITools()
	if err != nil {
		t.Fatalf("scanCITools failed: %v", err)
	}

	var gitlabCITool *ToolInfo
	for i := range tools {
		if tools[i].Name == "GitLab CI" {
			gitlabCITool = &tools[i]
			break
		}
	}

	if gitlabCITool == nil || !gitlabCITool.Detected {
		t.Error("GitLab CI should be detected")
	}
}

func TestToolScanner_ScanCITools_CircleCI(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .circleci directory
	circleDir := filepath.Join(tmpDir, ".circleci")
	if err := os.MkdirAll(circleDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write config.yml
	if err := os.WriteFile(filepath.Join(circleDir, "config.yml"), []byte("version: 2.1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanCITools()
	if err != nil {
		t.Fatalf("scanCITools failed: %v", err)
	}

	var circleCITool *ToolInfo
	for i := range tools {
		if tools[i].Name == "CircleCI" {
			circleCITool = &tools[i]
			break
		}
	}

	if circleCITool == nil || !circleCITool.Detected {
		t.Error("CircleCI should be detected")
	}
}

func TestToolScanner_ScanGitHooks_PreCommit(t *testing.T) {
	tmpDir := t.TempDir()

	// Write .pre-commit-config.yaml
	precommitConfig := `repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
`
	if err := os.WriteFile(filepath.Join(tmpDir, ".pre-commit-config.yaml"), []byte(precommitConfig), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGitHooks()
	if err != nil {
		t.Fatalf("scanGitHooks failed: %v", err)
	}

	var precommit *ToolInfo
	for i := range tools {
		if tools[i].Name == "pre-commit" {
			precommit = &tools[i]
			break
		}
	}

	if precommit == nil || !precommit.Detected {
		t.Error("pre-commit should be detected")
	} else {
		if precommit.ConfigFile != ".pre-commit-config.yaml" {
			t.Errorf("pre-commit config should be .pre-commit-config.yaml, got %s", precommit.ConfigFile)
		}
	}
}

func TestToolScanner_ScanGitHooks_Husky(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .husky directory
	huskyDir := filepath.Join(tmpDir, ".husky")
	if err := os.MkdirAll(huskyDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write pre-commit hook
	if err := os.WriteFile(filepath.Join(huskyDir, "pre-commit"), []byte("#!/bin/sh\nnpm test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGitHooks()
	if err != nil {
		t.Fatalf("scanGitHooks failed: %v", err)
	}

	var husky *ToolInfo
	for i := range tools {
		if tools[i].Name == "husky" {
			husky = &tools[i]
			break
		}
	}

	if husky == nil || !husky.Detected {
		t.Error("husky should be detected")
	} else {
		if husky.ConfigFile != ".husky/" {
			t.Errorf("husky config should be .husky/, got %s", husky.ConfigFile)
		}
	}
}

func TestToolScanner_ScanGitHooks_Lefthook(t *testing.T) {
	tmpDir := t.TempDir()

	// Write lefthook.yml
	if err := os.WriteFile(filepath.Join(tmpDir, "lefthook.yml"), []byte("pre-commit:\n  commands:\n    lint:\n      run: golangci-lint run\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGitHooks()
	if err != nil {
		t.Fatalf("scanGitHooks failed: %v", err)
	}

	var lefthook *ToolInfo
	for i := range tools {
		if tools[i].Name == "lefthook" {
			lefthook = &tools[i]
			break
		}
	}

	if lefthook == nil || !lefthook.Detected {
		t.Error("lefthook should be detected")
	}
}

func TestToolScanner_ScanAll(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mixed project with Go and GitHub Actions
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, ".golangci.yml"), []byte("linters:\n  enable:\n    - gofmt\n"), 0644); err != nil {
		t.Fatal(err)
	}

	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte("name: CI\non: push\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.ScanAll()
	if err != nil {
		t.Fatalf("ScanAll failed: %v", err)
	}

	// Should only include detected tools
	for _, tool := range tools {
		if !tool.Detected {
			t.Errorf("ScanAll returned undetected tool: %s", tool.Name)
		}
	}

	// Should have Go tools and GitHub Actions
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	expectedTools := []string{"golangci-lint", "gofmt", "go vet", "go test", "GitHub Actions"}
	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("expected tool %s not found in ScanAll results", expected)
		}
	}
}

func TestToolScanner_ScanForProjectType(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go project
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)

	// Scan for Go project type
	tools, err := scanner.ScanForProjectType(Go)
	if err != nil {
		t.Fatalf("ScanForProjectType failed: %v", err)
	}

	// Should contain Go-specific tools
	hasGoTool := false
	for _, tool := range tools {
		if tool.Name == "gofmt" || tool.Name == "go vet" || tool.Name == "go test" {
			hasGoTool = true
			break
		}
	}
	if !hasGoTool {
		t.Error("ScanForProjectType(Go) should return Go tools")
	}
}

func TestToolScanner_EmptyProject(t *testing.T) {
	tmpDir := t.TempDir()

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.ScanAll()
	if err != nil {
		t.Fatalf("ScanAll failed: %v", err)
	}

	// Empty project should return no detected tools
	if len(tools) != 0 {
		t.Errorf("empty project should have no detected tools, got %d", len(tools))
	}
}

func TestToolInfo_Category(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go project
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, ".golangci.yml"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.ScanAll()
	if err != nil {
		t.Fatalf("ScanAll failed: %v", err)
	}

	categoryMap := make(map[string]ToolCategory)
	for _, tool := range tools {
		categoryMap[tool.Name] = tool.Category
	}

	// Check categories are set correctly
	if cat, ok := categoryMap["golangci-lint"]; ok && cat != CategoryLinter {
		t.Errorf("golangci-lint should be CategoryLinter, got %s", cat)
	}
	if cat, ok := categoryMap["gofmt"]; ok && cat != CategoryFormatter {
		t.Errorf("gofmt should be CategoryFormatter, got %s", cat)
	}
	if cat, ok := categoryMap["go test"]; ok && cat != CategoryTesting {
		t.Errorf("go test should be CategoryTesting, got %s", cat)
	}
}

func TestToolScanner_ScanForProjectType_Node(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Node project
	pkgJSON := `{"name": "test", "devDependencies": {"eslint": "^8.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.ScanForProjectType(Node)
	if err != nil {
		t.Fatalf("ScanForProjectType failed: %v", err)
	}

	// Should contain Node-specific tools
	hasNodeTool := false
	for _, tool := range tools {
		if tool.Name == "eslint" || tool.Name == "npm audit" {
			hasNodeTool = true
			break
		}
	}
	if !hasNodeTool {
		t.Error("ScanForProjectType(Node) should return Node tools")
	}
}

func TestToolScanner_ScanForProjectType_Python(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Python project
	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte("[tool.black]\nline-length = 88\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.ScanForProjectType(Python)
	if err != nil {
		t.Fatalf("ScanForProjectType failed: %v", err)
	}

	// Should contain Python-specific tools
	hasPythonTool := false
	for _, tool := range tools {
		if tool.Name == "black" {
			hasPythonTool = true
			break
		}
	}
	if !hasPythonTool {
		t.Error("ScanForProjectType(Python) should return Python tools")
	}
}

func TestToolScanner_ScanForProjectType_Unknown(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go project to have something to detect
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.ScanForProjectType(Unknown)
	if err != nil {
		t.Fatalf("ScanForProjectType failed: %v", err)
	}

	// Should call ScanAll for unknown types
	if len(tools) == 0 {
		t.Error("ScanForProjectType(Unknown) should return tools from ScanAll")
	}
}

func TestToolScanner_ScanNodeTools_Mocha(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with mocha
	pkgJSON := `{"name": "test", "devDependencies": {"mocha": "^10.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// Write mocharc config
	if err := os.WriteFile(filepath.Join(tmpDir, ".mocharc.json"), []byte(`{"timeout": 5000}`), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	var mocha *ToolInfo
	for i := range tools {
		if tools[i].Name == "mocha" {
			mocha = &tools[i]
			break
		}
	}

	if mocha == nil || !mocha.Detected {
		t.Error("mocha should be detected")
	} else {
		if mocha.ConfigFile != ".mocharc.json" {
			t.Errorf("mocha config should be .mocharc.json, got %s", mocha.ConfigFile)
		}
	}
}

func TestToolScanner_ScanNodeTools_Vitest(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with vitest
	pkgJSON := `{"name": "test", "devDependencies": {"vitest": "^1.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// Write vitest config
	if err := os.WriteFile(filepath.Join(tmpDir, "vitest.config.ts"), []byte("export default {}"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	var vitest *ToolInfo
	for i := range tools {
		if tools[i].Name == "vitest" {
			vitest = &tools[i]
			break
		}
	}

	if vitest == nil || !vitest.Detected {
		t.Error("vitest should be detected")
	}
}

func TestToolScanner_ScanNodeTools_TypeScriptFromDeps(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with typescript in dependencies (not devDeps)
	pkgJSON := `{"name": "test", "dependencies": {"typescript": "^5.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	var typescript *ToolInfo
	for i := range tools {
		if tools[i].Name == "typescript" {
			typescript = &tools[i]
			break
		}
	}

	if typescript == nil || !typescript.Detected {
		t.Error("typescript should be detected from dependencies")
	}
}

func TestToolScanner_ScanNodeTools_EslintConfigVariants(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
	}{
		{"eslint.config.js", "eslint.config.js"},
		{"eslint.config.mjs", "eslint.config.mjs"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Write package.json
			if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name": "test"}`), 0644); err != nil {
				t.Fatal(err)
			}

			// Write config file
			if err := os.WriteFile(filepath.Join(tmpDir, tc.configFile), []byte("export default {}"), 0644); err != nil {
				t.Fatal(err)
			}

			scanner := NewToolScanner(tmpDir)
			tools, err := scanner.scanNodeTools()
			if err != nil {
				t.Fatalf("scanNodeTools failed: %v", err)
			}

			var eslint *ToolInfo
			for i := range tools {
				if tools[i].Name == "eslint" {
					eslint = &tools[i]
					break
				}
			}

			if eslint == nil || !eslint.Detected {
				t.Errorf("eslint should be detected with %s", tc.configFile)
			}
		})
	}
}

func TestToolScanner_ScanPythonTools_Flake8(t *testing.T) {
	tmpDir := t.TempDir()

	// Write .flake8
	if err := os.WriteFile(filepath.Join(tmpDir, ".flake8"), []byte("[flake8]\nmax-line-length = 88\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var flake8 *ToolInfo
	for i := range tools {
		if tools[i].Name == "flake8" {
			flake8 = &tools[i]
			break
		}
	}

	if flake8 == nil || !flake8.Detected {
		t.Error("flake8 should be detected")
	} else {
		if flake8.ConfigFile != ".flake8" {
			t.Errorf("flake8 config should be .flake8, got %s", flake8.ConfigFile)
		}
	}
}

func TestToolScanner_ScanPythonTools_Flake8InSetupCfg(t *testing.T) {
	tmpDir := t.TempDir()

	// Write setup.cfg with flake8 section
	setupCfg := `[flake8]
max-line-length = 88
`
	if err := os.WriteFile(filepath.Join(tmpDir, "setup.cfg"), []byte(setupCfg), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var flake8 *ToolInfo
	for i := range tools {
		if tools[i].Name == "flake8" {
			flake8 = &tools[i]
			break
		}
	}

	if flake8 == nil || !flake8.Detected {
		t.Error("flake8 should be detected from setup.cfg")
	}
}

func TestToolScanner_ScanPythonTools_PytestIni(t *testing.T) {
	tmpDir := t.TempDir()

	// Write pytest.ini
	if err := os.WriteFile(filepath.Join(tmpDir, "pytest.ini"), []byte("[pytest]\ntestpaths = tests\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var pytest *ToolInfo
	for i := range tools {
		if tools[i].Name == "pytest" {
			pytest = &tools[i]
			break
		}
	}

	if pytest == nil || !pytest.Detected {
		t.Error("pytest should be detected from pytest.ini")
	} else {
		if pytest.ConfigFile != "pytest.ini" {
			t.Errorf("pytest config should be pytest.ini, got %s", pytest.ConfigFile)
		}
	}
}

func TestToolScanner_ScanPythonTools_PylintRc(t *testing.T) {
	tmpDir := t.TempDir()

	// Write .pylintrc
	if err := os.WriteFile(filepath.Join(tmpDir, ".pylintrc"), []byte("[MASTER]\njobs=4\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var pylint *ToolInfo
	for i := range tools {
		if tools[i].Name == "pylint" {
			pylint = &tools[i]
			break
		}
	}

	if pylint == nil || !pylint.Detected {
		t.Error("pylint should be detected from .pylintrc")
	}
}

func TestToolScanner_ScanPythonTools_MypyIni(t *testing.T) {
	tmpDir := t.TempDir()

	// Write mypy.ini
	if err := os.WriteFile(filepath.Join(tmpDir, "mypy.ini"), []byte("[mypy]\nstrict = true\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var mypy *ToolInfo
	for i := range tools {
		if tools[i].Name == "mypy" {
			mypy = &tools[i]
			break
		}
	}

	if mypy == nil || !mypy.Detected {
		t.Error("mypy should be detected from mypy.ini")
	} else {
		if mypy.ConfigFile != "mypy.ini" {
			t.Errorf("mypy config should be mypy.ini, got %s", mypy.ConfigFile)
		}
	}
}

func TestToolScanner_ScanPythonTools_BlackInSetupCfg(t *testing.T) {
	tmpDir := t.TempDir()

	// Write setup.cfg with black section
	setupCfg := `[black]
line-length = 88
`
	if err := os.WriteFile(filepath.Join(tmpDir, "setup.cfg"), []byte(setupCfg), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var black *ToolInfo
	for i := range tools {
		if tools[i].Name == "black" {
			black = &tools[i]
			break
		}
	}

	if black == nil || !black.Detected {
		t.Error("black should be detected from setup.cfg")
	}
}

func TestToolScanner_ScanPythonTools_BlackInRequirements(t *testing.T) {
	tmpDir := t.TempDir()

	// Write requirements.txt with black
	if err := os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte("black==23.0.0\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var black *ToolInfo
	for i := range tools {
		if tools[i].Name == "black" {
			black = &tools[i]
			break
		}
	}

	if black == nil || !black.Detected {
		t.Error("black should be detected from requirements.txt")
	}
}

func TestToolScanner_ScanCITools_Jenkins(t *testing.T) {
	tmpDir := t.TempDir()

	// Write Jenkinsfile
	if err := os.WriteFile(filepath.Join(tmpDir, "Jenkinsfile"), []byte("pipeline {\n  agent any\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanCITools()
	if err != nil {
		t.Fatalf("scanCITools failed: %v", err)
	}

	var jenkins *ToolInfo
	for i := range tools {
		if tools[i].Name == "Jenkins" {
			jenkins = &tools[i]
			break
		}
	}

	if jenkins == nil || !jenkins.Detected {
		t.Error("Jenkins should be detected")
	}
}

func TestToolScanner_ScanCITools_TravisCI(t *testing.T) {
	tmpDir := t.TempDir()

	// Write .travis.yml
	if err := os.WriteFile(filepath.Join(tmpDir, ".travis.yml"), []byte("language: go\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanCITools()
	if err != nil {
		t.Fatalf("scanCITools failed: %v", err)
	}

	var travis *ToolInfo
	for i := range tools {
		if tools[i].Name == "Travis CI" {
			travis = &tools[i]
			break
		}
	}

	if travis == nil || !travis.Detected {
		t.Error("Travis CI should be detected")
	}
}

func TestToolScanner_ScanGitHooks_RawHooks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .git/hooks directory
	hooksDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write pre-commit hook (not .sample)
	if err := os.WriteFile(filepath.Join(hooksDir, "pre-commit"), []byte("#!/bin/sh\necho 'pre-commit'\n"), 0755); err != nil {
		t.Fatal(err)
	}

	// Write pre-push hook
	if err := os.WriteFile(filepath.Join(hooksDir, "pre-push"), []byte("#!/bin/sh\necho 'pre-push'\n"), 0755); err != nil {
		t.Fatal(err)
	}

	// Write a sample file (should be ignored)
	if err := os.WriteFile(filepath.Join(hooksDir, "commit-msg.sample"), []byte("#!/bin/sh\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGitHooks()
	if err != nil {
		t.Fatalf("scanGitHooks failed: %v", err)
	}

	var rawHooks *ToolInfo
	for i := range tools {
		if tools[i].Name == "git hooks" {
			rawHooks = &tools[i]
			break
		}
	}

	if rawHooks == nil || !rawHooks.Detected {
		t.Error("raw git hooks should be detected")
	} else {
		if len(rawHooks.Indicators) != 2 {
			t.Errorf("expected 2 hook indicators, got %d: %v", len(rawHooks.Indicators), rawHooks.Indicators)
		}
	}
}

func TestToolScanner_ScanGitHooks_HuskyInPackageJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with husky in devDeps
	pkgJSON := `{"name": "test", "devDependencies": {"husky": "^8.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGitHooks()
	if err != nil {
		t.Fatalf("scanGitHooks failed: %v", err)
	}

	var husky *ToolInfo
	for i := range tools {
		if tools[i].Name == "husky" {
			husky = &tools[i]
			break
		}
	}

	if husky == nil || !husky.Detected {
		t.Error("husky should be detected from package.json")
	}
}

func TestToolScanner_ScanGitHooks_PreCommitYml(t *testing.T) {
	tmpDir := t.TempDir()

	// Write .pre-commit-config.yml (alternate extension)
	if err := os.WriteFile(filepath.Join(tmpDir, ".pre-commit-config.yml"), []byte("repos: []\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGitHooks()
	if err != nil {
		t.Fatalf("scanGitHooks failed: %v", err)
	}

	var precommit *ToolInfo
	for i := range tools {
		if tools[i].Name == "pre-commit" {
			precommit = &tools[i]
			break
		}
	}

	if precommit == nil || !precommit.Detected {
		t.Error("pre-commit should be detected with .yml extension")
	} else {
		if precommit.ConfigFile != ".pre-commit-config.yml" {
			t.Errorf("pre-commit config should be .pre-commit-config.yml, got %s", precommit.ConfigFile)
		}
	}
}

func TestToolScanner_ScanGoTools_Goimports(t *testing.T) {
	tmpDir := t.TempDir()

	// Write go.mod
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Write Makefile mentioning goimports
	if err := os.WriteFile(filepath.Join(tmpDir, "Makefile"), []byte("fmt:\n\tgoimports -w .\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGoTools()
	if err != nil {
		t.Fatalf("scanGoTools failed: %v", err)
	}

	var goimports *ToolInfo
	for i := range tools {
		if tools[i].Name == "goimports" {
			goimports = &tools[i]
			break
		}
	}

	if goimports == nil || !goimports.Detected {
		t.Error("goimports should be detected from Makefile")
	}
}

func TestToolScanner_ReadPackageJSON_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	// Write invalid JSON
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("not valid json"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	pkg, err := scanner.readPackageJSON()

	// Should return error for invalid JSON
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if pkg != nil {
		t.Error("expected nil package for invalid JSON")
	}
}

func TestToolScanner_ReadPackageJSON_NotExists(t *testing.T) {
	tmpDir := t.TempDir()

	scanner := NewToolScanner(tmpDir)
	pkg, err := scanner.readPackageJSON()

	// Should return error when file doesn't exist
	if err == nil {
		t.Error("expected error for missing file")
	}
	if pkg != nil {
		t.Error("expected nil package for missing file")
	}
}

func TestPackageJSON_HasMethods(t *testing.T) {
	// Test with nil package
	var nilPkg *packageJSON
	if nilPkg.hasDevDep("test") {
		t.Error("nil package should return false for hasDevDep")
	}
	if nilPkg.hasDep("test") {
		t.Error("nil package should return false for hasDep")
	}
	if nilPkg.hasField("test") {
		t.Error("nil package should return false for hasField")
	}

	// Test with empty package
	emptyPkg := &packageJSON{}
	if emptyPkg.hasDevDep("test") {
		t.Error("empty package should return false for hasDevDep")
	}
	if emptyPkg.hasDep("test") {
		t.Error("empty package should return false for hasDep")
	}
	if emptyPkg.hasField("test") {
		t.Error("empty package should return false for hasField")
	}
}

func TestToolScanner_ScanNodeTools_EslintInPackageJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with eslintConfig field (no separate config file)
	pkgJSON := `{
		"name": "test",
		"eslintConfig": {
			"extends": "eslint:recommended"
		}
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	var eslint *ToolInfo
	for i := range tools {
		if tools[i].Name == "eslint" {
			eslint = &tools[i]
			break
		}
	}

	if eslint == nil || !eslint.Detected {
		t.Error("eslint should be detected from package.json eslintConfig field")
	}
}

func TestToolScanner_ScanNodeTools_PrettierInPackageJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with prettier field
	pkgJSON := `{
		"name": "test",
		"prettier": {
			"semi": false
		}
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	var prettier *ToolInfo
	for i := range tools {
		if tools[i].Name == "prettier" {
			prettier = &tools[i]
			break
		}
	}

	if prettier == nil || !prettier.Detected {
		t.Error("prettier should be detected from package.json prettier field")
	}
}

func TestToolScanner_ScanNodeTools_JestInPackageJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with jest field
	pkgJSON := `{
		"name": "test",
		"jest": {
			"testEnvironment": "node"
		}
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	var jest *ToolInfo
	for i := range tools {
		if tools[i].Name == "jest" {
			jest = &tools[i]
			break
		}
	}

	if jest == nil || !jest.Detected {
		t.Error("jest should be detected from package.json jest field")
	}
}

func TestToolScanner_ScanPythonTools_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Empty directory - no Python tools should be detected
	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	for _, tool := range tools {
		if tool.Detected {
			t.Errorf("No tools should be detected in empty directory, but found: %s", tool.Name)
		}
	}
}

func TestToolScanner_ScanPythonTools_PytestInSetupCfg(t *testing.T) {
	tmpDir := t.TempDir()

	// Write setup.cfg with pytest section
	setupCfg := `[pytest]
testpaths = tests
`
	if err := os.WriteFile(filepath.Join(tmpDir, "setup.cfg"), []byte(setupCfg), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var pytest *ToolInfo
	for i := range tools {
		if tools[i].Name == "pytest" {
			pytest = &tools[i]
			break
		}
	}

	if pytest == nil || !pytest.Detected {
		t.Error("pytest should be detected from setup.cfg")
	}
}

func TestToolScanner_ScanPythonTools_RuffInPyproject(t *testing.T) {
	tmpDir := t.TempDir()

	// Write pyproject.toml with ruff section
	pyproject := `[tool.ruff]
line-length = 100
`
	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyproject), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var ruff *ToolInfo
	for i := range tools {
		if tools[i].Name == "ruff" {
			ruff = &tools[i]
			break
		}
	}

	if ruff == nil || !ruff.Detected {
		t.Error("ruff should be detected from pyproject.toml")
	} else {
		if ruff.ConfigFile != "pyproject.toml" {
			t.Errorf("ruff config should be pyproject.toml, got %s", ruff.ConfigFile)
		}
	}
}

func TestToolScanner_ScanPythonTools_PylintInPyproject(t *testing.T) {
	tmpDir := t.TempDir()

	// Write pyproject.toml with pylint section
	pyproject := `[tool.pylint]
max-line-length = 120
`
	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyproject), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var pylint *ToolInfo
	for i := range tools {
		if tools[i].Name == "pylint" {
			pylint = &tools[i]
			break
		}
	}

	if pylint == nil || !pylint.Detected {
		t.Error("pylint should be detected from pyproject.toml")
	}
}

func TestToolScanner_ScanAll_WithErrors(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a realistic project with multiple tool types
	files := map[string]string{
		"go.mod":                  "module test\n\ngo 1.21\n",
		".golangci.yml":           "linters:\n  enable:\n    - gofmt\n",
		"package.json":            `{"name": "test", "devDependencies": {"eslint": "^8.0.0"}}`,
		".pre-commit-config.yaml": "repos: []\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create .github/workflows directory
	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte("name: CI\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.ScanAll()
	if err != nil {
		t.Fatalf("ScanAll failed: %v", err)
	}

	// Should detect tools from multiple categories
	categories := make(map[ToolCategory]bool)
	for _, tool := range tools {
		categories[tool.Category] = true
	}

	if !categories[CategoryLinter] {
		t.Error("Should detect linter tools")
	}
	if !categories[CategoryHooks] {
		t.Error("Should detect hook tools")
	}
	if !categories[CategoryCI] {
		t.Error("Should detect CI tools")
	}
}

func TestToolScanner_ScanNodeTools_NoPackageJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Only ESLint config file, no package.json
	if err := os.WriteFile(filepath.Join(tmpDir, ".eslintrc.json"), []byte(`{"extends": "eslint:recommended"}`), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanNodeTools()
	if err != nil {
		t.Fatalf("scanNodeTools failed: %v", err)
	}

	var eslint *ToolInfo
	for i := range tools {
		if tools[i].Name == "eslint" {
			eslint = &tools[i]
			break
		}
	}

	// ESLint should still be detected from config file
	if eslint == nil || !eslint.Detected {
		t.Error("eslint should be detected from .eslintrc.json even without package.json")
	}
}

func TestToolScanner_ScanGitHooks_HuskyField(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with husky field (old style config)
	pkgJSON := `{
		"name": "test",
		"husky": {
			"hooks": {
				"pre-commit": "npm test"
			}
		}
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanGitHooks()
	if err != nil {
		t.Fatalf("scanGitHooks failed: %v", err)
	}

	var husky *ToolInfo
	for i := range tools {
		if tools[i].Name == "husky" {
			husky = &tools[i]
			break
		}
	}

	if husky == nil || !husky.Detected {
		t.Error("husky should be detected from package.json husky field")
	}
}

func TestToolScanner_ScanCITools_MultipleWorkflows(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .github/workflows directory with multiple workflow files
	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create multiple workflow files
	workflows := []string{"ci.yml", "release.yaml", "test.yml"}
	for _, wf := range workflows {
		if err := os.WriteFile(filepath.Join(workflowDir, wf), []byte("name: "+wf+"\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanCITools()
	if err != nil {
		t.Fatalf("scanCITools failed: %v", err)
	}

	var githubActions *ToolInfo
	for i := range tools {
		if tools[i].Name == "GitHub Actions" {
			githubActions = &tools[i]
			break
		}
	}

	if githubActions == nil || !githubActions.Detected {
		t.Error("GitHub Actions should be detected")
	} else {
		// Should have indicators for all workflow files
		if len(githubActions.Indicators) < 3 {
			t.Errorf("Should have indicators for all workflows, got %d: %v", len(githubActions.Indicators), githubActions.Indicators)
		}
	}
}

func TestToolScanner_ScanPythonTools_MypyFromRequirements(t *testing.T) {
	tmpDir := t.TempDir()

	// Write requirements-dev.txt with mypy
	if err := os.WriteFile(filepath.Join(tmpDir, "requirements-dev.txt"), []byte("mypy==1.0.0\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var mypy *ToolInfo
	for i := range tools {
		if tools[i].Name == "mypy" {
			mypy = &tools[i]
			break
		}
	}

	if mypy == nil || !mypy.Detected {
		t.Error("mypy should be detected from requirements-dev.txt")
	}
}

func TestToolScanner_ScanPythonTools_PytestFromRequirements(t *testing.T) {
	tmpDir := t.TempDir()

	// Write requirements.txt with pytest
	if err := os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte("pytest>=7.0.0\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewToolScanner(tmpDir)
	tools, err := scanner.scanPythonTools()
	if err != nil {
		t.Fatalf("scanPythonTools failed: %v", err)
	}

	var pytest *ToolInfo
	for i := range tools {
		if tools[i].Name == "pytest" {
			pytest = &tools[i]
			break
		}
	}

	if pytest == nil || !pytest.Detected {
		t.Error("pytest should be detected from requirements.txt")
	}
}

func TestToolScanner_ScanNodeTools_PrettierConfigVariants(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
	}{
		{"prettierrc", ".prettierrc"},
		{"prettierrc.js", ".prettierrc.js"},
		{"prettierrc.json", ".prettierrc.json"},
		{"prettier.config.js", "prettier.config.js"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Write package.json (required for npm audit detection)
			if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name": "test"}`), 0644); err != nil {
				t.Fatal(err)
			}

			// Write config file
			if err := os.WriteFile(filepath.Join(tmpDir, tc.configFile), []byte("{}"), 0644); err != nil {
				t.Fatal(err)
			}

			scanner := NewToolScanner(tmpDir)
			tools, err := scanner.scanNodeTools()
			if err != nil {
				t.Fatalf("scanNodeTools failed: %v", err)
			}

			var prettier *ToolInfo
			for i := range tools {
				if tools[i].Name == "prettier" {
					prettier = &tools[i]
					break
				}
			}

			if prettier == nil || !prettier.Detected {
				t.Errorf("prettier should be detected with %s", tc.configFile)
			}
		})
	}
}

func TestToolScanner_ScanGoTools_GolangciLintVariants(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
	}{
		{"yml", ".golangci.yml"},
		{"yaml", ".golangci.yaml"},
		{"toml", ".golangci.toml"},
		{"json", ".golangci.json"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Write go.mod
			if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0644); err != nil {
				t.Fatal(err)
			}

			// Write config file
			if err := os.WriteFile(filepath.Join(tmpDir, tc.configFile), []byte(""), 0644); err != nil {
				t.Fatal(err)
			}

			scanner := NewToolScanner(tmpDir)
			tools, err := scanner.scanGoTools()
			if err != nil {
				t.Fatalf("scanGoTools failed: %v", err)
			}

			var golangciLint *ToolInfo
			for i := range tools {
				if tools[i].Name == "golangci-lint" {
					golangciLint = &tools[i]
					break
				}
			}

			if golangciLint == nil || !golangciLint.Detected {
				t.Errorf("golangci-lint should be detected with %s", tc.configFile)
			} else {
				if golangciLint.ConfigFile != tc.configFile {
					t.Errorf("golangci-lint config should be %s, got %s", tc.configFile, golangciLint.ConfigFile)
				}
			}
		})
	}
}
