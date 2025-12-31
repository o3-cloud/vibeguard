package templates

func init() {
	Register(Template{
		Name:        "node-javascript",
		Description: "JavaScript/Node.js project with ESLint, Prettier, and testing",
		Content: `version: "1"

vars:
  source_dir: "src"

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

  - id: test
    run: npm test -- --passWithNoTests
    severity: error
    suggestion: "Run 'npm test' to diagnose test failures"
    timeout: 300s
    requires:
      - lint

  - id: build
    run: npm run build
    severity: error
    suggestion: "Run 'npm run build' to diagnose build errors"
    timeout: 120s
`,
	})
}
