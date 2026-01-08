package templates

func init() {
	Register(Template{
		Name:        "rust-cargo",
		Description: "Rust project using Cargo with clippy, formatting, and testing",
		Content: `version: "1"

vars:
  min_coverage: "70"

checks:
  - id: fmt
    run: cargo fmt -- --check
    severity: error
    suggestion: "Run 'cargo fmt' to format your code"
    timeout: 60s

  - id: clippy
    run: cargo clippy -- -D warnings
    severity: error
    suggestion: "Run 'cargo clippy --fix' to fix linting issues"
    timeout: 120s

  - id: analyze
    run: cargo deny check 2>/dev/null || echo "cargo-deny not installed (optional)"
    severity: warning
    suggestion: "Install cargo-deny: cargo install cargo-deny"
    timeout: 120s
    requires:
      - clippy

  - id: security
    run: cargo audit 2>/dev/null || echo "cargo-audit not installed (optional)"
    severity: warning
    suggestion: "Install cargo-audit: cargo install cargo-audit"
    timeout: 120s
    requires:
      - analyze

  - id: test
    run: cargo test
    severity: error
    suggestion: "Run 'cargo test' to diagnose test failures"
    timeout: 300s
    requires:
      - clippy

  - id: coverage
    run: cargo tarpaulin --out Stdout --timeout 300
    grok:
      - Coverage:\s+%{NUMBER:coverage}%
    assert: "coverage >= {{.min_coverage}}"
    severity: warning
    suggestion: "Code coverage is below {{.min_coverage}}%. Increase test coverage."
    timeout: 600s
    requires:
      - test

  - id: build
    run: cargo build --release
    severity: error
    suggestion: "Run 'cargo build' to diagnose build errors"
    timeout: 300s
    requires:
      - clippy
`,
	})
}
