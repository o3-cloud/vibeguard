package config

import (
	"testing"
)

func TestInterpolateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		vars     map[string]string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			vars:     map[string]string{"foo": "bar"},
			expected: "",
		},
		{
			name:     "no placeholders",
			input:    "hello world",
			vars:     map[string]string{"foo": "bar"},
			expected: "hello world",
		},
		{
			name:     "single placeholder",
			input:    "hello {{.name}}",
			vars:     map[string]string{"name": "world"},
			expected: "hello world",
		},
		{
			name:     "multiple same placeholders",
			input:    "{{.x}} and {{.x}}",
			vars:     map[string]string{"x": "value"},
			expected: "value and value",
		},
		{
			name:     "multiple different placeholders",
			input:    "{{.a}} {{.b}} {{.c}}",
			vars:     map[string]string{"a": "1", "b": "2", "c": "3"},
			expected: "1 2 3",
		},
		{
			name:     "undefined placeholder remains",
			input:    "hello {{.undefined}}",
			vars:     map[string]string{"foo": "bar"},
			expected: "hello {{.undefined}}",
		},
		{
			name:     "placeholder at start",
			input:    "{{.prefix}}suffix",
			vars:     map[string]string{"prefix": "pre-"},
			expected: "pre-suffix",
		},
		{
			name:     "placeholder at end",
			input:    "prefix{{.suffix}}",
			vars:     map[string]string{"suffix": "-end"},
			expected: "prefix-end",
		},
		{
			name:     "placeholder only",
			input:    "{{.value}}",
			vars:     map[string]string{"value": "result"},
			expected: "result",
		},
		{
			name:     "empty vars map",
			input:    "{{.foo}} stays",
			vars:     map[string]string{},
			expected: "{{.foo}} stays",
		},
		{
			name:     "nil vars map",
			input:    "{{.foo}} stays",
			vars:     nil,
			expected: "{{.foo}} stays",
		},
		{
			name:     "value with special characters",
			input:    "path: {{.path}}",
			vars:     map[string]string{"path": "./..."},
			expected: "path: ./...",
		},
		{
			name:     "underscore in variable name",
			input:    "{{.go_packages}}",
			vars:     map[string]string{"go_packages": "./cmd/..."},
			expected: "./cmd/...",
		},
		{
			name:     "command with variable",
			input:    "go test -race {{.packages}}",
			vars:     map[string]string{"packages": "./..."},
			expected: "go test -race ./...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Vars: tt.vars}
			result := cfg.interpolateString(tt.input)
			if result != tt.expected {
				t.Errorf("interpolateString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInterpolate(t *testing.T) {
	cfg := &Config{
		Vars: map[string]string{
			"packages": "./...",
			"timeout":  "30s",
		},
		Checks: []Check{
			{
				ID:         "test-check",
				Run:        "go test {{.packages}}",
				Assert:     "{{.packages}} passed",
				Suggestion: "Run with {{.packages}} to fix",
				File:       "{{.packages}}/output.txt",
				Grok:       GrokSpec{"pattern {{.packages}}", "another {{.timeout}}"},
			},
			{
				ID:  "no-vars",
				Run: "echo hello",
			},
		},
	}

	cfg.Interpolate()

	// Check first check
	if cfg.Checks[0].Run != "go test ./..." {
		t.Errorf("Run not interpolated: got %q", cfg.Checks[0].Run)
	}
	if cfg.Checks[0].Assert != "./... passed" {
		t.Errorf("Assert not interpolated: got %q", cfg.Checks[0].Assert)
	}
	if cfg.Checks[0].Suggestion != "Run with ./... to fix" {
		t.Errorf("Suggestion not interpolated: got %q", cfg.Checks[0].Suggestion)
	}
	if cfg.Checks[0].File != "./.../output.txt" {
		t.Errorf("File not interpolated: got %q", cfg.Checks[0].File)
	}
	if cfg.Checks[0].Grok[0] != "pattern ./..." {
		t.Errorf("Grok[0] not interpolated: got %q", cfg.Checks[0].Grok[0])
	}
	if cfg.Checks[0].Grok[1] != "another 30s" {
		t.Errorf("Grok[1] not interpolated: got %q", cfg.Checks[0].Grok[1])
	}

	// Check second check unchanged
	if cfg.Checks[1].Run != "echo hello" {
		t.Errorf("Run changed unexpectedly: got %q", cfg.Checks[1].Run)
	}
}

func TestInterpolateWithExtracted(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		vars      map[string]string
		extracted map[string]string
		expected  string
	}{
		{
			name:      "empty template",
			template:  "",
			vars:      map[string]string{"a": "1"},
			extracted: map[string]string{"b": "2"},
			expected:  "",
		},
		{
			name:      "config vars only",
			template:  "value is {{.a}}",
			vars:      map[string]string{"a": "config-val"},
			extracted: nil,
			expected:  "value is config-val",
		},
		{
			name:      "extracted vars only",
			template:  "captured: {{.error_msg}}",
			vars:      nil,
			extracted: map[string]string{"error_msg": "file not found"},
			expected:  "captured: file not found",
		},
		{
			name:      "both config and extracted",
			template:  "{{.tool}}: {{.error_msg}}",
			vars:      map[string]string{"tool": "golangci-lint"},
			extracted: map[string]string{"error_msg": "unused variable"},
			expected:  "golangci-lint: unused variable",
		},
		{
			name:      "config vars applied first then extracted",
			template:  "value: {{.key}}",
			vars:      map[string]string{"key": "from-config"},
			extracted: map[string]string{"key": "from-extracted"},
			expected:  "value: from-config", // config vars are applied first and win
		},
		{
			name:      "complex suggestion template",
			template:  "Fix {{.error_type}} error in {{.file}}:{{.line}} - {{.message}}",
			vars:      map[string]string{},
			extracted: map[string]string{"error_type": "syntax", "file": "main.go", "line": "42", "message": "unexpected EOF"},
			expected:  "Fix syntax error in main.go:42 - unexpected EOF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InterpolateWithExtracted(tt.template, tt.vars, tt.extracted)
			if result != tt.expected {
				t.Errorf("InterpolateWithExtracted() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestInterpolateEmptyChecks(t *testing.T) {
	cfg := &Config{
		Vars:   map[string]string{"foo": "bar"},
		Checks: []Check{},
	}

	// Should not panic
	cfg.Interpolate()

	if len(cfg.Checks) != 0 {
		t.Errorf("Checks should still be empty")
	}
}

func TestInterpolateNilVars(t *testing.T) {
	cfg := &Config{
		Vars: nil,
		Checks: []Check{
			{
				ID:  "test",
				Run: "{{.undefined}}",
			},
		},
	}

	// Should not panic
	cfg.Interpolate()

	// Undefined vars should remain as-is
	if cfg.Checks[0].Run != "{{.undefined}}" {
		t.Errorf("Run should remain unchanged with nil vars: got %q", cfg.Checks[0].Run)
	}
}
