package templates

func init() {
	Register(Template{
		Name:        "python-poetry",
		Description: "Python project using Poetry with ruff, mypy, and pytest",
		Content: `version: "1"

vars:
  source_dir: "src"
  min_coverage: "70"

checks:
  - id: format
    run: poetry run ruff format --check {{.source_dir}}
    severity: error
    suggestion: "Run 'poetry run ruff format {{.source_dir}}' to format your code"
    timeout: 60s

  - id: lint
    run: poetry run ruff check {{.source_dir}}
    severity: error
    suggestion: "Run 'poetry run ruff check --fix {{.source_dir}}' to fix linting issues"
    timeout: 60s

  - id: typecheck
    run: poetry run mypy {{.source_dir}}
    severity: warning
    suggestion: "Fix type errors shown in the mypy output"
    timeout: 120s

  - id: analyze
    run: poetry run pylint {{.source_dir}} --disable=all --enable=E,F 2>/dev/null || echo "pylint not installed (optional)"
    severity: warning
    suggestion: "Install pylint with poetry: poetry add --group dev pylint"
    timeout: 120s
    requires:
      - lint

  - id: security
    run: poetry run pip-audit 2>/dev/null || echo "pip-audit not installed (optional)"
    severity: warning
    suggestion: "Install pip-audit: poetry add --group dev pip-audit"
    timeout: 120s
    requires:
      - analyze

  - id: test
    run: poetry run pytest
    severity: error
    suggestion: "Run 'poetry run pytest -v' to diagnose test failures"
    timeout: 300s
    requires:
      - lint

  - id: coverage
    run: poetry run pytest --cov={{.source_dir}} --cov-report=term-missing
    grok:
      - TOTAL\s+\d+\s+\d+\s+%{NUMBER:coverage}%
    assert: "coverage >= {{.min_coverage}}"
    severity: warning
    suggestion: "Code coverage is below {{.min_coverage}}%. Increase test coverage."
    timeout: 300s
    requires:
      - test

  - id: build
    run: poetry install && poetry run python -c "import {{.source_dir}}"
    severity: error
    suggestion: "Run 'poetry install' to diagnose installation errors"
    timeout: 120s
    requires:
      - lint
`,
	})
}
