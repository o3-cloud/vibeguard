package templates

func init() {
	Register(Template{
		Name:        "go-minimal",
		Description: "Minimal Go project with basic formatting, vetting, and testing",
		Content: `version: "1"

vars:
  go_packages: "./..."

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

  - id: test
    run: go test {{.go_packages}}
    severity: error
    suggestion: "Run 'go test ./...' to diagnose test failures"
    timeout: 300s
`,
	})
}
