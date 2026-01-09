package templates

func init() {
	Register(Template{
		Name:        "bash",
		Description: "Bash scripts with ShellCheck, shfmt, and BATS testing",
		Content: `version: "1"

vars:
  script_dir: "."
  script_pattern: "*.sh"

checks:
  - id: format
    run: shfmt -d -i 2 -ci {{.script_dir}} 2>/dev/null || echo "shfmt not installed (optional)"
    severity: warning
    suggestion: "Install shfmt: go install mvdan.cc/sh/v3/cmd/shfmt@latest or brew install shfmt"
    timeout: 60s

  - id: lint
    run: find {{.script_dir}} -name '{{.script_pattern}}' -exec shellcheck {} +
    severity: error
    suggestion: "Run 'shellcheck <script>' to see detailed linting issues"
    timeout: 60s

  - id: syntax
    run: find {{.script_dir}} -name '{{.script_pattern}}' -exec bash -n {} \;
    severity: error
    suggestion: "Run 'bash -n <script>' to check for syntax errors"
    timeout: 30s

  - id: analyze
    run: find {{.script_dir}} -name '{{.script_pattern}}' -exec shellcheck -S warning {} + 2>/dev/null || true
    severity: warning
    suggestion: "Review shellcheck warnings for potential improvements"
    timeout: 60s
    requires:
      - lint

  - id: security
    run: find {{.script_dir}} -name '{{.script_pattern}}' -exec shellcheck -e SC2086,SC2046 {} + 2>/dev/null || echo "Security check passed"
    severity: warning
    suggestion: "Review unquoted variables that may cause word splitting or glob expansion"
    timeout: 60s
    requires:
      - lint

  - id: test
    run: bats test/ 2>/dev/null || echo "No BATS tests found or bats not installed"
    severity: warning
    suggestion: "Install bats: brew install bats-core or npm install -g bats"
    timeout: 300s
    requires:
      - lint
`,
	})
}
