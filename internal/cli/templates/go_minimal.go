package templates

func init() {
	Register(Template{
		Name:        "go-minimal",
		Description: "Minimal Go project with basic formatting, vetting, and testing",
		Content: `version: "1"

vars:
  go_packages: "./..."
  min_coverage: "70"

checks:
  - id: fmt
    run: test -z "$(gofmt -l .)"
    severity: error
    suggestion: "Run 'gofmt -w .' to format your code"
    timeout: 30s

  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    suggestion: "Run 'go vet ./...' and fix reported issues"
    timeout: 60s

  - id: analyze
    run: staticcheck {{.go_packages}} 2>/dev/null || echo "staticcheck not installed (optional)"
    severity: warning
    suggestion: "Install staticcheck: go install honnef.co/go/tools/cmd/staticcheck@latest"
    timeout: 120s
    requires:
      - vet

  - id: security
    run: gosec {{.go_packages}} 2>/dev/null || echo "gosec not installed (optional)"
    severity: warning
    suggestion: "Install gosec: go install github.com/securego/gosec/v2/cmd/gosec@latest"
    timeout: 120s
    requires:
      - analyze

  - id: test
    run: go test {{.go_packages}}
    severity: error
    suggestion: "Run 'go test ./...' to diagnose test failures"
    timeout: 300s

  - id: coverage
    run: go test {{.go_packages}} -coverprofile=cover.out && go tool cover -func=cover.out
    grok:
      - total:.*\(statements\)\s+%{NUMBER:coverage}%
    assert: "coverage >= {{.min_coverage}}"
    severity: warning
    suggestion: "Code coverage is below {{.min_coverage}}%. Increase test coverage."
    timeout: 300s
    requires:
      - test

  - id: build
    run: go build {{.go_packages}}
    severity: error
    suggestion: "Run 'go build ./...' to diagnose build errors"
    timeout: 120s
    requires:
      - vet
`,
	})
}
