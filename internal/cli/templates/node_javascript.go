package templates

func init() {
	Register(Template{
		Name:        "node-javascript",
		Description: "JavaScript/Node.js project with ESLint, Prettier, and testing",
		Content: `version: "1"

vars:
  source_dir: "src"
  min_coverage: "70"

checks:
  - id: format
    run: npx prettier --check .
    severity: error
    suggestion: "Run 'npx prettier --write .' to format your code"
    timeout: 60s

  - id: lint
    run: npx eslint {{.source_dir}} --max-warnings 0
    severity: error
    suggestion: "Run 'npx eslint {{.source_dir}} --fix' to fix linting issues"
    timeout: 60s

  - id: analyze
    run: npm run analyze 2>/dev/null || echo "Static analysis tool not configured (optional)"
    severity: warning
    suggestion: "Configure static analysis tool (e.g., sonarjs, code-inspector) in package.json"
    timeout: 120s
    requires:
      - lint

  - id: test
    run: npm test -- --passWithNoTests
    severity: error
    suggestion: "Run 'npm test' to diagnose test failures"
    timeout: 300s
    requires:
      - lint

  - id: coverage
    run: npm test -- --coverage --passWithNoTests
    grok:
      - "Lines\\s+:\\s+%{NUMBER:coverage}%"
    assert: "coverage >= {{.min_coverage}}"
    severity: warning
    suggestion: "Code coverage is below {{.min_coverage}}%. Increase test coverage."
    timeout: 300s
    requires:
      - test

  - id: build
    run: npm run build
    severity: error
    suggestion: "Run 'npm run build' to diagnose build errors"
    timeout: 120s
`,
	})
}
