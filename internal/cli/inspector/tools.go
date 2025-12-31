// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ToolCategory represents the category of a development tool.
type ToolCategory string

// Supported tool categories.
const (
	CategoryLinter    ToolCategory = "linter"
	CategoryFormatter ToolCategory = "formatter"
	CategoryTesting   ToolCategory = "testing"
	CategoryBuild     ToolCategory = "build"
	CategoryCI        ToolCategory = "ci"
	CategoryHooks     ToolCategory = "hooks"
	CategoryTypeCheck ToolCategory = "typecheck"
	CategorySecurity  ToolCategory = "security"
)

// ToolInfo holds information about a detected development tool.
type ToolInfo struct {
	Name       string       // Tool name (e.g., "golangci-lint", "eslint")
	Category   ToolCategory // Tool category
	Detected   bool         // Whether the tool was detected
	Version    string       // Version if detectable
	ConfigFile string       // Path to config file if found
	Confidence float64      // Confidence score 0.0-1.0
	Indicators []string     // What led to this detection
}

// ToolScanner detects development tools in a project.
type ToolScanner struct {
	root string
}

// NewToolScanner creates a new ToolScanner for the given project root.
func NewToolScanner(root string) *ToolScanner {
	return &ToolScanner{root: root}
}

// ScanAll detects all development tools in the project.
func (s *ToolScanner) ScanAll() ([]ToolInfo, error) {
	var tools []ToolInfo

	// Scan Go tools
	goTools, err := s.scanGoTools()
	if err != nil {
		return nil, err
	}
	tools = append(tools, goTools...)

	// Scan Node tools
	nodeTools, err := s.scanNodeTools()
	if err != nil {
		return nil, err
	}
	tools = append(tools, nodeTools...)

	// Scan Python tools
	pythonTools, err := s.scanPythonTools()
	if err != nil {
		return nil, err
	}
	tools = append(tools, pythonTools...)

	// Scan CI/CD
	ciTools, err := s.scanCITools()
	if err != nil {
		return nil, err
	}
	tools = append(tools, ciTools...)

	// Scan Git hooks
	hookTools, err := s.scanGitHooks()
	if err != nil {
		return nil, err
	}
	tools = append(tools, hookTools...)

	// Filter to only detected tools
	var detected []ToolInfo
	for _, tool := range tools {
		if tool.Detected {
			detected = append(detected, tool)
		}
	}

	return detected, nil
}

// ScanForProjectType scans tools relevant to a specific project type.
func (s *ToolScanner) ScanForProjectType(projectType ProjectType) ([]ToolInfo, error) {
	switch projectType {
	case Go:
		return s.scanGoTools()
	case Node:
		return s.scanNodeTools()
	case Python:
		return s.scanPythonTools()
	default:
		return s.ScanAll()
	}
}

// scanGoTools detects Go-specific development tools.
func (s *ToolScanner) scanGoTools() ([]ToolInfo, error) {
	var tools []ToolInfo

	// golangci-lint
	golangciLint := ToolInfo{
		Name:     "golangci-lint",
		Category: CategoryLinter,
	}
	if configPath := s.findFile(".golangci.yml", ".golangci.yaml", ".golangci.toml", ".golangci.json"); configPath != "" {
		golangciLint.Detected = true
		golangciLint.ConfigFile = configPath
		golangciLint.Confidence = 0.9
		golangciLint.Indicators = []string{configPath}
	}
	tools = append(tools, golangciLint)

	// gofmt (always available with Go)
	gofmt := ToolInfo{
		Name:     "gofmt",
		Category: CategoryFormatter,
	}
	if s.fileExists("go.mod") {
		gofmt.Detected = true
		gofmt.Confidence = 1.0
		gofmt.Indicators = []string{"go.mod present (gofmt included with Go)"}
	}
	tools = append(tools, gofmt)

	// go vet (always available with Go)
	govet := ToolInfo{
		Name:     "go vet",
		Category: CategoryLinter,
	}
	if s.fileExists("go.mod") {
		govet.Detected = true
		govet.Confidence = 1.0
		govet.Indicators = []string{"go.mod present (go vet included with Go)"}
	}
	tools = append(tools, govet)

	// go test (always available with Go)
	gotest := ToolInfo{
		Name:     "go test",
		Category: CategoryTesting,
	}
	if s.fileExists("go.mod") {
		gotest.Detected = true
		gotest.Confidence = 1.0
		gotest.Indicators = []string{"go.mod present (go test included with Go)"}
	}
	tools = append(tools, gotest)

	// goimports
	goimports := ToolInfo{
		Name:     "goimports",
		Category: CategoryFormatter,
	}
	// Check if goimports is mentioned in Makefile or scripts
	if s.fileContains("Makefile", "goimports") || s.fileContains("go.mod", "goimports") {
		goimports.Detected = true
		goimports.Confidence = 0.7
		goimports.Indicators = []string{"goimports mentioned in Makefile or go.mod"}
	}
	tools = append(tools, goimports)

	return tools, nil
}

// scanNodeTools detects Node.js-specific development tools.
func (s *ToolScanner) scanNodeTools() ([]ToolInfo, error) {
	var tools []ToolInfo

	// Read package.json for devDependencies
	pkgJSON, _ := s.readPackageJSON()

	// ESLint
	eslint := ToolInfo{
		Name:     "eslint",
		Category: CategoryLinter,
	}
	if configPath := s.findFile(".eslintrc", ".eslintrc.js", ".eslintrc.cjs", ".eslintrc.json", ".eslintrc.yml", ".eslintrc.yaml", "eslint.config.js", "eslint.config.mjs"); configPath != "" {
		eslint.Detected = true
		eslint.ConfigFile = configPath
		eslint.Confidence = 0.9
		eslint.Indicators = []string{configPath}
	} else if pkgJSON != nil && (pkgJSON.hasDevDep("eslint") || pkgJSON.hasField("eslintConfig")) {
		eslint.Detected = true
		eslint.Confidence = 0.8
		eslint.Indicators = []string{"eslint in package.json"}
	}
	tools = append(tools, eslint)

	// Prettier
	prettier := ToolInfo{
		Name:     "prettier",
		Category: CategoryFormatter,
	}
	if configPath := s.findFile(".prettierrc", ".prettierrc.js", ".prettierrc.cjs", ".prettierrc.json", ".prettierrc.yml", ".prettierrc.yaml", "prettier.config.js", "prettier.config.cjs"); configPath != "" {
		prettier.Detected = true
		prettier.ConfigFile = configPath
		prettier.Confidence = 0.9
		prettier.Indicators = []string{configPath}
	} else if pkgJSON != nil && (pkgJSON.hasDevDep("prettier") || pkgJSON.hasField("prettier")) {
		prettier.Detected = true
		prettier.Confidence = 0.8
		prettier.Indicators = []string{"prettier in package.json"}
	}
	tools = append(tools, prettier)

	// Jest
	jest := ToolInfo{
		Name:     "jest",
		Category: CategoryTesting,
	}
	if configPath := s.findFile("jest.config.js", "jest.config.ts", "jest.config.mjs", "jest.config.cjs", "jest.config.json"); configPath != "" {
		jest.Detected = true
		jest.ConfigFile = configPath
		jest.Confidence = 0.9
		jest.Indicators = []string{configPath}
	} else if pkgJSON != nil && (pkgJSON.hasDevDep("jest") || pkgJSON.hasField("jest")) {
		jest.Detected = true
		jest.Confidence = 0.8
		jest.Indicators = []string{"jest in package.json"}
	}
	tools = append(tools, jest)

	// Mocha
	mocha := ToolInfo{
		Name:     "mocha",
		Category: CategoryTesting,
	}
	if configPath := s.findFile(".mocharc.js", ".mocharc.json", ".mocharc.yml", ".mocharc.yaml", ".mocharc.cjs"); configPath != "" {
		mocha.Detected = true
		mocha.ConfigFile = configPath
		mocha.Confidence = 0.9
		mocha.Indicators = []string{configPath}
	} else if pkgJSON != nil && pkgJSON.hasDevDep("mocha") {
		mocha.Detected = true
		mocha.Confidence = 0.8
		mocha.Indicators = []string{"mocha in package.json"}
	}
	tools = append(tools, mocha)

	// Vitest
	vitest := ToolInfo{
		Name:     "vitest",
		Category: CategoryTesting,
	}
	if configPath := s.findFile("vitest.config.js", "vitest.config.ts", "vitest.config.mjs"); configPath != "" {
		vitest.Detected = true
		vitest.ConfigFile = configPath
		vitest.Confidence = 0.9
		vitest.Indicators = []string{configPath}
	} else if pkgJSON != nil && pkgJSON.hasDevDep("vitest") {
		vitest.Detected = true
		vitest.Confidence = 0.8
		vitest.Indicators = []string{"vitest in package.json"}
	}
	tools = append(tools, vitest)

	// TypeScript
	typescript := ToolInfo{
		Name:     "typescript",
		Category: CategoryTypeCheck,
	}
	if configPath := s.findFile("tsconfig.json"); configPath != "" {
		typescript.Detected = true
		typescript.ConfigFile = configPath
		typescript.Confidence = 0.95
		typescript.Indicators = []string{configPath}
	} else if pkgJSON != nil && (pkgJSON.hasDevDep("typescript") || pkgJSON.hasDep("typescript")) {
		typescript.Detected = true
		typescript.Confidence = 0.8
		typescript.Indicators = []string{"typescript in package.json"}
	}
	tools = append(tools, typescript)

	// npm audit (always available with npm)
	npmAudit := ToolInfo{
		Name:     "npm audit",
		Category: CategorySecurity,
	}
	if s.fileExists("package.json") {
		npmAudit.Detected = true
		npmAudit.Confidence = 1.0
		npmAudit.Indicators = []string{"package.json present (npm audit available)"}
	}
	tools = append(tools, npmAudit)

	return tools, nil
}

// scanPythonTools detects Python-specific development tools.
func (s *ToolScanner) scanPythonTools() ([]ToolInfo, error) {
	var tools []ToolInfo

	// Read pyproject.toml for tool configs
	hasPyproject := s.fileExists("pyproject.toml")
	pyprojectContent := ""
	if hasPyproject {
		content, _ := os.ReadFile(filepath.Join(s.root, "pyproject.toml"))
		pyprojectContent = string(content)
	}

	// Black
	black := ToolInfo{
		Name:     "black",
		Category: CategoryFormatter,
	}
	if hasPyproject && strings.Contains(pyprojectContent, "[tool.black]") {
		black.Detected = true
		black.ConfigFile = "pyproject.toml"
		black.Confidence = 0.9
		black.Indicators = []string{"[tool.black] in pyproject.toml"}
	} else if s.fileExists("setup.cfg") && s.fileContains("setup.cfg", "[black]") {
		black.Detected = true
		black.ConfigFile = "setup.cfg"
		black.Confidence = 0.9
		black.Indicators = []string{"[black] in setup.cfg"}
	} else if s.fileContains("requirements.txt", "black") || s.fileContains("requirements-dev.txt", "black") {
		black.Detected = true
		black.Confidence = 0.7
		black.Indicators = []string{"black in requirements"}
	}
	tools = append(tools, black)

	// Pylint
	pylint := ToolInfo{
		Name:     "pylint",
		Category: CategoryLinter,
	}
	if configPath := s.findFile(".pylintrc", "pylintrc", "pyproject.toml"); configPath != "" {
		if configPath == "pyproject.toml" && !strings.Contains(pyprojectContent, "[tool.pylint]") {
			// pyproject.toml exists but no pylint config
		} else {
			pylint.Detected = true
			pylint.ConfigFile = configPath
			pylint.Confidence = 0.9
			pylint.Indicators = []string{configPath}
		}
	}
	if !pylint.Detected && (s.fileContains("requirements.txt", "pylint") || s.fileContains("requirements-dev.txt", "pylint")) {
		pylint.Detected = true
		pylint.Confidence = 0.7
		pylint.Indicators = []string{"pylint in requirements"}
	}
	tools = append(tools, pylint)

	// Pytest
	pytest := ToolInfo{
		Name:     "pytest",
		Category: CategoryTesting,
	}
	if configPath := s.findFile("pytest.ini", "pyproject.toml", "setup.cfg"); configPath != "" {
		if configPath == "pyproject.toml" && !strings.Contains(pyprojectContent, "[tool.pytest]") {
			// Check also for pytest.ini.options
			if strings.Contains(pyprojectContent, "pytest") {
				pytest.Detected = true
				pytest.ConfigFile = configPath
				pytest.Confidence = 0.85
				pytest.Indicators = []string{"pytest config in pyproject.toml"}
			}
		} else if configPath == "setup.cfg" && s.fileContains("setup.cfg", "[pytest]") {
			pytest.Detected = true
			pytest.ConfigFile = configPath
			pytest.Confidence = 0.9
			pytest.Indicators = []string{"[pytest] in setup.cfg"}
		} else if configPath == "pytest.ini" {
			pytest.Detected = true
			pytest.ConfigFile = configPath
			pytest.Confidence = 0.95
			pytest.Indicators = []string{configPath}
		}
	}
	if !pytest.Detected && (s.fileContains("requirements.txt", "pytest") || s.fileContains("requirements-dev.txt", "pytest")) {
		pytest.Detected = true
		pytest.Confidence = 0.7
		pytest.Indicators = []string{"pytest in requirements"}
	}
	tools = append(tools, pytest)

	// Mypy
	mypy := ToolInfo{
		Name:     "mypy",
		Category: CategoryTypeCheck,
	}
	if configPath := s.findFile("mypy.ini", ".mypy.ini"); configPath != "" {
		mypy.Detected = true
		mypy.ConfigFile = configPath
		mypy.Confidence = 0.95
		mypy.Indicators = []string{configPath}
	} else if hasPyproject && strings.Contains(pyprojectContent, "[tool.mypy]") {
		mypy.Detected = true
		mypy.ConfigFile = "pyproject.toml"
		mypy.Confidence = 0.9
		mypy.Indicators = []string{"[tool.mypy] in pyproject.toml"}
	} else if s.fileContains("requirements.txt", "mypy") || s.fileContains("requirements-dev.txt", "mypy") {
		mypy.Detected = true
		mypy.Confidence = 0.7
		mypy.Indicators = []string{"mypy in requirements"}
	}
	tools = append(tools, mypy)

	// Ruff (modern Python linter)
	ruff := ToolInfo{
		Name:     "ruff",
		Category: CategoryLinter,
	}
	if configPath := s.findFile("ruff.toml", ".ruff.toml"); configPath != "" {
		ruff.Detected = true
		ruff.ConfigFile = configPath
		ruff.Confidence = 0.95
		ruff.Indicators = []string{configPath}
	} else if hasPyproject && strings.Contains(pyprojectContent, "[tool.ruff]") {
		ruff.Detected = true
		ruff.ConfigFile = "pyproject.toml"
		ruff.Confidence = 0.9
		ruff.Indicators = []string{"[tool.ruff] in pyproject.toml"}
	}
	tools = append(tools, ruff)

	// Flake8
	flake8 := ToolInfo{
		Name:     "flake8",
		Category: CategoryLinter,
	}
	if configPath := s.findFile(".flake8"); configPath != "" {
		flake8.Detected = true
		flake8.ConfigFile = configPath
		flake8.Confidence = 0.95
		flake8.Indicators = []string{configPath}
	} else if s.fileExists("setup.cfg") && s.fileContains("setup.cfg", "[flake8]") {
		flake8.Detected = true
		flake8.ConfigFile = "setup.cfg"
		flake8.Confidence = 0.9
		flake8.Indicators = []string{"[flake8] in setup.cfg"}
	}
	tools = append(tools, flake8)

	return tools, nil
}

// scanCITools detects CI/CD configurations.
func (s *ToolScanner) scanCITools() ([]ToolInfo, error) {
	var tools []ToolInfo

	// GitHub Actions
	githubActions := ToolInfo{
		Name:     "GitHub Actions",
		Category: CategoryCI,
	}
	if s.dirExists(".github/workflows") {
		// Find workflow files
		workflowFiles, _ := filepath.Glob(filepath.Join(s.root, ".github/workflows/*.yml"))
		yamlFiles, _ := filepath.Glob(filepath.Join(s.root, ".github/workflows/*.yaml"))
		workflowFiles = append(workflowFiles, yamlFiles...)
		if len(workflowFiles) > 0 {
			githubActions.Detected = true
			githubActions.ConfigFile = ".github/workflows/"
			githubActions.Confidence = 0.95
			indicators := make([]string, len(workflowFiles))
			for i, f := range workflowFiles {
				indicators[i], _ = filepath.Rel(s.root, f)
			}
			githubActions.Indicators = indicators
		}
	}
	tools = append(tools, githubActions)

	// GitLab CI
	gitlabCI := ToolInfo{
		Name:     "GitLab CI",
		Category: CategoryCI,
	}
	if s.fileExists(".gitlab-ci.yml") {
		gitlabCI.Detected = true
		gitlabCI.ConfigFile = ".gitlab-ci.yml"
		gitlabCI.Confidence = 0.95
		gitlabCI.Indicators = []string{".gitlab-ci.yml"}
	}
	tools = append(tools, gitlabCI)

	// CircleCI
	circleCI := ToolInfo{
		Name:     "CircleCI",
		Category: CategoryCI,
	}
	if configPath := s.findFile(".circleci/config.yml", ".circleci/config.yaml"); configPath != "" {
		circleCI.Detected = true
		circleCI.ConfigFile = configPath
		circleCI.Confidence = 0.95
		circleCI.Indicators = []string{configPath}
	}
	tools = append(tools, circleCI)

	// Jenkins
	jenkins := ToolInfo{
		Name:     "Jenkins",
		Category: CategoryCI,
	}
	if configPath := s.findFile("Jenkinsfile", "jenkins/Jenkinsfile"); configPath != "" {
		jenkins.Detected = true
		jenkins.ConfigFile = configPath
		jenkins.Confidence = 0.95
		jenkins.Indicators = []string{configPath}
	}
	tools = append(tools, jenkins)

	// Travis CI
	travisCI := ToolInfo{
		Name:     "Travis CI",
		Category: CategoryCI,
	}
	if s.fileExists(".travis.yml") {
		travisCI.Detected = true
		travisCI.ConfigFile = ".travis.yml"
		travisCI.Confidence = 0.95
		travisCI.Indicators = []string{".travis.yml"}
	}
	tools = append(tools, travisCI)

	return tools, nil
}

// scanGitHooks detects Git hook configurations.
func (s *ToolScanner) scanGitHooks() ([]ToolInfo, error) {
	var tools []ToolInfo

	// Pre-commit (Python-based hook manager)
	precommit := ToolInfo{
		Name:     "pre-commit",
		Category: CategoryHooks,
	}
	if s.fileExists(".pre-commit-config.yaml") || s.fileExists(".pre-commit-config.yml") {
		configPath := ".pre-commit-config.yaml"
		if s.fileExists(".pre-commit-config.yml") {
			configPath = ".pre-commit-config.yml"
		}
		precommit.Detected = true
		precommit.ConfigFile = configPath
		precommit.Confidence = 0.95
		precommit.Indicators = []string{configPath}
	}
	tools = append(tools, precommit)

	// Husky (Node.js-based hook manager)
	husky := ToolInfo{
		Name:     "husky",
		Category: CategoryHooks,
	}
	if s.dirExists(".husky") {
		husky.Detected = true
		husky.ConfigFile = ".husky/"
		husky.Confidence = 0.95
		husky.Indicators = []string{".husky/ directory"}
	} else {
		// Check package.json for husky config
		pkgJSON, _ := s.readPackageJSON()
		if pkgJSON != nil && (pkgJSON.hasDevDep("husky") || pkgJSON.hasField("husky")) {
			husky.Detected = true
			husky.Confidence = 0.8
			husky.Indicators = []string{"husky in package.json"}
		}
	}
	tools = append(tools, husky)

	// Lefthook (Go-based hook manager)
	lefthook := ToolInfo{
		Name:     "lefthook",
		Category: CategoryHooks,
	}
	if configPath := s.findFile("lefthook.yml", "lefthook.yaml", ".lefthook.yml", ".lefthook.yaml"); configPath != "" {
		lefthook.Detected = true
		lefthook.ConfigFile = configPath
		lefthook.Confidence = 0.95
		lefthook.Indicators = []string{configPath}
	}
	tools = append(tools, lefthook)

	// Raw git hooks
	rawHooks := ToolInfo{
		Name:     "git hooks",
		Category: CategoryHooks,
	}
	hookDir := filepath.Join(s.root, ".git/hooks")
	if s.dirExists(".git/hooks") {
		entries, err := os.ReadDir(hookDir)
		if err == nil {
			var activeHooks []string
			for _, entry := range entries {
				name := entry.Name()
				// Skip sample files
				if strings.HasSuffix(name, ".sample") {
					continue
				}
				// Check common hook names
				if name == "pre-commit" || name == "pre-push" || name == "commit-msg" ||
					name == "prepare-commit-msg" || name == "post-commit" {
					activeHooks = append(activeHooks, name)
				}
			}
			if len(activeHooks) > 0 {
				rawHooks.Detected = true
				rawHooks.ConfigFile = ".git/hooks/"
				rawHooks.Confidence = 0.9
				rawHooks.Indicators = activeHooks
			}
		}
	}
	tools = append(tools, rawHooks)

	return tools, nil
}

// packageJSON holds parsed package.json content.
type packageJSON struct {
	Name            string                 `json:"name"`
	Version         string                 `json:"version"`
	Dependencies    map[string]string      `json:"dependencies"`
	DevDependencies map[string]string      `json:"devDependencies"`
	Scripts         map[string]string      `json:"scripts"`
	Raw             map[string]interface{} `json:"-"`
}

func (p *packageJSON) hasDevDep(name string) bool {
	if p == nil || p.DevDependencies == nil {
		return false
	}
	_, ok := p.DevDependencies[name]
	return ok
}

func (p *packageJSON) hasDep(name string) bool {
	if p == nil || p.Dependencies == nil {
		return false
	}
	_, ok := p.Dependencies[name]
	return ok
}

func (p *packageJSON) hasField(name string) bool {
	if p == nil || p.Raw == nil {
		return false
	}
	_, ok := p.Raw[name]
	return ok
}

// readPackageJSON reads and parses package.json if it exists.
func (s *ToolScanner) readPackageJSON() (*packageJSON, error) {
	path := filepath.Join(s.root, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	// Also unmarshal raw for field checking
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err == nil {
		pkg.Raw = raw
	}

	return &pkg, nil
}

// findFile returns the first existing file from the given paths.
func (s *ToolScanner) findFile(paths ...string) string {
	for _, p := range paths {
		fullPath := filepath.Join(s.root, p)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			return p
		}
	}
	return ""
}

// fileExists checks if a file exists in the project root.
func (s *ToolScanner) fileExists(name string) bool {
	path := filepath.Join(s.root, name)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists in the project root.
func (s *ToolScanner) dirExists(name string) bool {
	path := filepath.Join(s.root, name)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// fileContains checks if a file exists and contains the given string.
func (s *ToolScanner) fileContains(name, substr string) bool {
	path := filepath.Join(s.root, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), substr)
}
