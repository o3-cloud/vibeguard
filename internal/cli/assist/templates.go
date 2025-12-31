// Package assist provides AI agent-assisted setup functionality for VibeGuard.
package assist

// languageExamples contains language-specific YAML examples for each supported project type.
var languageExamples = map[string]string{
	"go":      goExamples,
	"node":    nodeExamples,
	"python":  pythonExamples,
	"rust":    rustExamples,
	"ruby":    rubyExamples,
	"java":    javaExamples,
	"generic": genericExamples,
}

const goExamples = `## Go-Specific Examples

### Format Check
` + "```yaml" + `
- id: fmt
  run: test -z "$(gofmt -l .)"
  severity: error
  suggestion: "Run 'gofmt -w .' to format code"
  timeout: 5s
` + "```" + `

### Lint Check with golangci-lint
` + "```yaml" + `
- id: lint
  run: golangci-lint run {{.go_packages}}
  severity: error
  suggestion: "Fix linting issues. Run 'golangci-lint run --fix' for auto-fixes."
  timeout: 60s
  requires:
    - fmt
    - vet
` + "```" + `

### Go Vet Check
` + "```yaml" + `
- id: vet
  run: go vet {{.go_packages}}
  severity: error
  suggestion: "Fix go vet issues. These often indicate real bugs."
  timeout: 30s
` + "```" + `

### Test with Coverage
` + "```yaml" + `
- id: test
  run: go test -race {{.go_packages}}
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s

- id: coverage
  run: go test ./... -coverprofile cover.out && go tool cover -func cover.out
  grok:
    - total:.*\(statements\)\s+%{NUMBER:coverage}%
  assert: "coverage >= {{.coverage_threshold}}"
  severity: warning
  suggestion: "Coverage is {{.coverage}}%, target is {{.coverage_threshold}}%. Add tests."
  requires:
    - test
  timeout: 300s
` + "```" + `

### Build Check
` + "```yaml" + `
- id: build
  run: go build -o /dev/null {{.go_packages}}
  severity: error
  suggestion: "Fix compilation errors before committing"
  timeout: 60s
  requires:
    - vet
` + "```" + ``

const nodeExamples = `## Node.js-Specific Examples

### Format Check with Prettier
` + "```yaml" + `
- id: format
  run: npx prettier --check .
  severity: error
  suggestion: "Run 'npx prettier --write .' to format code"
  timeout: 30s
` + "```" + `

### Lint Check with ESLint
` + "```yaml" + `
- id: lint
  run: npx eslint .
  severity: error
  suggestion: "Fix ESLint errors. Run 'npx eslint . --fix' for auto-fixes."
  timeout: 60s
  requires:
    - format
` + "```" + `

### TypeScript Type Checking
` + "```yaml" + `
- id: typecheck
  run: npx tsc --noEmit
  severity: error
  suggestion: "Fix TypeScript type errors before committing"
  timeout: 60s
` + "```" + `

### Jest Tests with Coverage
` + "```yaml" + `
- id: test
  run: npm test
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s
  requires:
    - typecheck

- id: coverage
  run: npm test -- --coverage --coverageReporters=text-summary
  grok:
    - Statements\s+:\s+%{NUMBER:coverage}%
  assert: "coverage >= {{.coverage_threshold}}"
  severity: warning
  suggestion: "Coverage is {{.coverage}}%, target is {{.coverage_threshold}}%."
  requires:
    - test
  timeout: 300s
` + "```" + `

### Security Audit
` + "```yaml" + `
- id: audit
  run: npm audit --audit-level={{.audit_level}}
  severity: warning
  suggestion: "Run 'npm audit fix' or update vulnerable dependencies"
  timeout: 60s
` + "```" + `

### Build Check
` + "```yaml" + `
- id: build
  run: npm run build
  severity: error
  suggestion: "Fix build errors before committing"
  timeout: 120s
  requires:
    - lint
    - typecheck
` + "```" + ``

const pythonExamples = `## Python-Specific Examples

### Format Check with Black
` + "```yaml" + `
- id: format
  run: black --check .
  severity: error
  suggestion: "Run 'black .' to format code"
  timeout: 30s
` + "```" + `

### Import Sorting with isort
` + "```yaml" + `
- id: imports
  run: isort --check-only --diff .
  severity: error
  suggestion: "Run 'isort .' to fix import ordering"
  timeout: 30s
` + "```" + `

### Lint Check with Ruff
` + "```yaml" + `
- id: lint
  run: ruff check .
  severity: error
  suggestion: "Fix Ruff errors. Run 'ruff check --fix .' for auto-fixes."
  timeout: 30s
  requires:
    - format
    - imports
` + "```" + `

### Type Checking with mypy
` + "```yaml" + `
- id: typecheck
  run: mypy .
  severity: error
  suggestion: "Fix mypy type errors. Consider adding type hints."
  timeout: 60s
` + "```" + `

### Pytest Tests with Coverage
` + "```yaml" + `
- id: test
  run: pytest
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s

- id: coverage
  run: pytest --cov --cov-report=term-missing 2>&1 | grep 'TOTAL'
  grok:
    - TOTAL\s+\d+\s+\d+\s+%{NUMBER:coverage}%
  assert: "coverage >= {{.coverage_threshold}}"
  severity: warning
  suggestion: "Coverage is {{.coverage}}%, target is {{.coverage_threshold}}%. Add tests."
  requires:
    - test
  timeout: 300s
` + "```" + `

### Security Audit with pip-audit
` + "```yaml" + `
- id: audit
  run: pip-audit
  severity: warning
  suggestion: "Run 'pip-audit --fix' or update vulnerable packages"
  timeout: 120s
` + "```" + ``

const rustExamples = `## Rust-Specific Examples

### Format Check
` + "```yaml" + `
- id: format
  run: cargo fmt -- --check
  severity: error
  suggestion: "Run 'cargo fmt' to format code"
  timeout: 30s
` + "```" + `

### Lint Check with Clippy
` + "```yaml" + `
- id: lint
  run: cargo clippy -- -D warnings
  severity: error
  suggestion: "Fix Clippy warnings. Run 'cargo clippy --fix' for auto-fixes."
  timeout: 60s
  requires:
    - format
` + "```" + `

### Tests
` + "```yaml" + `
- id: test
  run: cargo test
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s
  requires:
    - lint
` + "```" + `

### Build Check
` + "```yaml" + `
- id: build
  run: cargo build --release
  severity: error
  suggestion: "Fix compilation errors before committing"
  timeout: 120s
` + "```" + `

### Security Audit
` + "```yaml" + `
- id: audit
  run: cargo audit
  severity: warning
  suggestion: "Update vulnerable dependencies"
  timeout: 60s
` + "```" + ``

const rubyExamples = `## Ruby-Specific Examples

### Lint Check with RuboCop
` + "```yaml" + `
- id: lint
  run: bundle exec rubocop
  severity: error
  suggestion: "Fix RuboCop issues. Run 'bundle exec rubocop -a' for auto-fixes."
  timeout: 60s
` + "```" + `

### RSpec Tests
` + "```yaml" + `
- id: test
  run: bundle exec rspec
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s
  requires:
    - lint
` + "```" + `

### Security Audit with Bundler Audit
` + "```yaml" + `
- id: audit
  run: bundle exec bundle-audit check --update
  severity: warning
  suggestion: "Update vulnerable gems"
  timeout: 60s
` + "```" + ``

const javaExamples = `## Java-Specific Examples

### Build with Maven
` + "```yaml" + `
- id: build
  run: mvn compile -q
  severity: error
  suggestion: "Fix compilation errors"
  timeout: 120s
` + "```" + `

### Tests with Maven
` + "```yaml" + `
- id: test
  run: mvn test -q
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s
  requires:
    - build
` + "```" + `

### Checkstyle
` + "```yaml" + `
- id: lint
  run: mvn checkstyle:check -q
  severity: error
  suggestion: "Fix Checkstyle violations"
  timeout: 60s
` + "```" + `

### Build with Gradle
` + "```yaml" + `
- id: build
  run: ./gradlew build -x test
  severity: error
  suggestion: "Fix compilation errors"
  timeout: 120s
` + "```" + `

### Tests with Gradle
` + "```yaml" + `
- id: test
  run: ./gradlew test
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s
  requires:
    - build
` + "```" + ``

const genericExamples = `## Generic Examples

### Shell Script Lint Check
` + "```yaml" + `
- id: shellcheck
  run: shellcheck **/*.sh
  severity: warning
  suggestion: "Fix ShellCheck warnings for better script quality"
  timeout: 30s
` + "```" + `

### Markdown Lint Check
` + "```yaml" + `
- id: markdown
  run: markdownlint '**/*.md'
  severity: warning
  suggestion: "Fix markdown formatting issues"
  timeout: 30s
` + "```" + `

### YAML Lint Check
` + "```yaml" + `
- id: yamllint
  run: yamllint .
  severity: warning
  suggestion: "Fix YAML formatting issues"
  timeout: 30s
` + "```" + `

### Secret Detection
` + "```yaml" + `
- id: secrets
  run: gitleaks detect --no-git
  severity: error
  suggestion: "Remove detected secrets before committing"
  timeout: 60s
` + "```" + ``
