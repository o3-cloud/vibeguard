package templates

func init() {
	Register(Template{
		Name:        "python-uv",
		Description: "Python project using uv package manager with ruff, mypy, and pytest",
		Content: `version: "1"

vars:
  source_dir: "src"
  min_coverage: "70"

checks:
  - id: format
    run: uv run ruff format --check {{.source_dir}}
    severity: error
    suggestion: "Run 'uv run ruff format {{.source_dir}}' to format your code"
    timeout: 60s

  - id: lint
    run: uv run ruff check {{.source_dir}}
    severity: error
    suggestion: "Run 'uv run ruff check --fix {{.source_dir}}' to fix linting issues"
    timeout: 60s

  - id: typecheck
    run: uv run mypy {{.source_dir}}
    severity: warning
    suggestion: "Fix type errors shown in the mypy output"
    timeout: 120s

  - id: analyze
    run: uv run pylint {{.source_dir}} --disable=all --enable=E,F 2>/dev/null || echo "pylint not installed (optional)"
    severity: warning
    suggestion: "Install pylint with uv: uv pip install pylint"
    timeout: 120s
    requires:
      - lint

  - id: test
    run: uv run pytest
    severity: error
    suggestion: "Run 'uv run pytest -v' to diagnose test failures"
    timeout: 300s
    requires:
      - lint

  - id: coverage
    run: uv run pytest --cov={{.source_dir}} --cov-report=term-missing
    grok:
      - TOTAL\s+\d+\s+\d+\s+%{NUMBER:coverage}%
    assert: "coverage >= {{.min_coverage}}"
    severity: warning
    suggestion: "Code coverage is below {{.min_coverage}}%. Increase test coverage."
    timeout: 300s
    requires:
      - test

  - id: build
    run: uv sync && uv run python -c "import {{.source_dir}}"
    severity: error
    suggestion: "Run 'uv sync' to diagnose installation errors"
    timeout: 120s
    requires:
      - lint
`,
	})
}
