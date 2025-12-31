package templates

func init() {
	Register(Template{
		Name:        "rust-cargo",
		Description: "Rust project using Cargo with clippy, formatting, and testing",
		Content: `version: "1"

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

  - id: test
    run: cargo test
    severity: error
    suggestion: "Run 'cargo test' to diagnose test failures"
    timeout: 300s
    requires:
      - clippy

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
