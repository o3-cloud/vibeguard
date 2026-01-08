package templates

func init() {
	Register(Template{
		Name:        "bun-typescript",
		Description: "TypeScript/Bun project with ESLint, Prettier, type checking, and testing",
		Content: `version: "1"

vars:
  source_dir: "src"
  min_coverage: "70"

checks:
  - id: format
    run: bunx prettier --check .
    severity: error
    suggestion: "Run 'bunx prettier --write .' to format your code"
    timeout: 60s

  - id: lint
    run: bunx eslint {{.source_dir}} --max-warnings 0
    severity: error
    suggestion: "Run 'bunx eslint {{.source_dir}} --fix' to fix linting issues"
    timeout: 60s

  - id: typecheck
    run: bunx tsc --noEmit
    severity: error
    suggestion: "Fix TypeScript type errors shown in the output"
    timeout: 120s

  - id: analyze
    run: bun run analyze 2>/dev/null || echo "Static analysis tool not configured (optional)"
    severity: warning
    suggestion: "Configure static analysis tool (e.g., sonarjs, code-inspector) in package.json"
    timeout: 120s
    requires:
      - lint

  - id: security
    run: bun audit 2>/dev/null || echo "bun audit failed - check for vulnerabilities"
    severity: warning
    suggestion: "Run 'bun audit' to see vulnerability details and 'bun update' to update packages"
    timeout: 120s
    requires:
      - analyze

  - id: test
    run: bun test --passWithNoTests 2>/dev/null || bun run test -- --passWithNoTests
    severity: error
    suggestion: "Run 'bun test' to diagnose test failures"
    timeout: 300s
    requires:
      - lint
      - typecheck

  - id: build
    run: bun run build
    severity: error
    suggestion: "Run 'bun run build' to diagnose build errors"
    timeout: 120s
    requires:
      - typecheck
`,
	})
}
