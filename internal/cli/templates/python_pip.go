package templates

func init() {
	Register(Template{
		Name:        "python-pip",
		Description: "Python project using pip with ruff, mypy, and pytest",
		Content: `version: "1"

vars:
  source_dir: "src"
  min_coverage: "70"

checks:
  - id: format
    run: ruff format --check {{.source_dir}}
    severity: error
    suggestion: "Run 'ruff format {{.source_dir}}' to format your code"
    timeout: 60s

  - id: lint
    run: ruff check {{.source_dir}}
    severity: error
    suggestion: "Run 'ruff check --fix {{.source_dir}}' to fix linting issues"
    timeout: 60s

  - id: typecheck
    run: mypy {{.source_dir}}
    severity: warning
    suggestion: "Fix type errors shown in the mypy output"
    timeout: 120s

  - id: analyze
    run: pylint {{.source_dir}} --disable=all --enable=E,F 2>/dev/null || echo "pylint not installed (optional)"
    severity: warning
    suggestion: "Install pylint: pip install pylint"
    timeout: 120s
    requires:
      - lint

  - id: test
    run: pytest
    severity: error
    suggestion: "Run 'pytest -v' to diagnose test failures"
    timeout: 300s
    requires:
      - lint

  - id: coverage
    run: pytest --cov={{.source_dir}} --cov-report=term-missing
    grok:
      - TOTAL\s+\d+\s+\d+\s+%{NUMBER:coverage}%
    assert: "coverage >= {{.min_coverage}}"
    severity: warning
    suggestion: "Code coverage is below {{.min_coverage}}%. Increase test coverage."
    timeout: 300s
    requires:
      - test

  - id: build
    run: pip install -e . && python -c "import {{.source_dir}}"
    severity: error
    suggestion: "Run 'pip install -e .' to diagnose installation errors"
    timeout: 120s
    requires:
      - lint
`,
	})
}
