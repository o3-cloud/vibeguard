package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// TestEventValue_UnmarshalYAML_StringInlineContent tests unmarshaling string (inline content)
func TestEventValue_UnmarshalYAML_StringInlineContent(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    EventValue
		wantErr bool
	}{
		{
			name:    "simple_string",
			yaml:    `"This is inline content"`,
			want:    EventValue{Content: "This is inline content", IsInline: true},
			wantErr: false,
		},
		{
			name:    "empty_string",
			yaml:    `""`,
			want:    EventValue{Content: "", IsInline: true},
			wantErr: false,
		},
		{
			name:    "string_that_looks_like_id",
			yaml:    `"init"`,
			want:    EventValue{Content: "init", IsInline: true},
			wantErr: false,
		},
		{
			name:    "multiline_string",
			yaml:    `"This is a\nmultiline\nstring"`,
			want:    EventValue{Content: "This is a\nmultiline\nstring", IsInline: true},
			wantErr: false,
		},
		{
			name:    "string_with_special_chars",
			yaml:    `"Check {{.variable}} for details"`,
			want:    EventValue{Content: "Check {{.variable}} for details", IsInline: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ev EventValue
			err := yaml.Unmarshal([]byte(tt.yaml), &ev)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if ev.Content != tt.want.Content || ev.IsInline != tt.want.IsInline {
				t.Errorf("UnmarshalYAML() = %+v, want %+v", ev, tt.want)
			}
		})
	}
}

// TestEventValue_UnmarshalYAML_ArrayPromptIDs tests unmarshaling arrays (prompt IDs)
func TestEventValue_UnmarshalYAML_ArrayPromptIDs(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    EventValue
		wantErr bool
	}{
		{
			name:    "single_id_in_array",
			yaml:    `["init"]`,
			want:    EventValue{IDs: []string{"init"}, IsInline: false},
			wantErr: false,
		},
		{
			name:    "multiple_ids",
			yaml:    `["init", "code-review", "security-audit"]`,
			want:    EventValue{IDs: []string{"init", "code-review", "security-audit"}, IsInline: false},
			wantErr: false,
		},
		{
			name:    "empty_array",
			yaml:    `[]`,
			want:    EventValue{IDs: []string{}, IsInline: false},
			wantErr: false,
		},
		{
			name: "array_flow_style",
			yaml: `[init, code-review]`,
			want: EventValue{IDs: []string{"init", "code-review"}, IsInline: false},
		},
		{
			name: "array_with_hyphens",
			yaml: `["test-generator", "security-audit"]`,
			want: EventValue{IDs: []string{"test-generator", "security-audit"}, IsInline: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ev EventValue
			err := yaml.Unmarshal([]byte(tt.yaml), &ev)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(ev.IDs) != len(tt.want.IDs) {
				t.Errorf("UnmarshalYAML() IDs length = %d, want %d", len(ev.IDs), len(tt.want.IDs))
				return
			}
			for i, id := range ev.IDs {
				if id != tt.want.IDs[i] {
					t.Errorf("UnmarshalYAML() IDs[%d] = %s, want %s", i, id, tt.want.IDs[i])
				}
			}
			if ev.IsInline != tt.want.IsInline {
				t.Errorf("UnmarshalYAML() IsInline = %v, want %v", ev.IsInline, tt.want.IsInline)
			}
		})
	}
}

// TestEventValue_MarshalYAML tests marshaling EventValue back to YAML
func TestEventValue_MarshalYAML(t *testing.T) {
	tests := []struct {
		name  string
		ev    EventValue
		want  interface{}
		check func(interface{}) bool
	}{
		{
			name: "inline_content",
			ev:   EventValue{Content: "This is inline", IsInline: true},
			want: "This is inline",
			check: func(v interface{}) bool {
				s, ok := v.(string)
				return ok && s == "This is inline"
			},
		},
		{
			name: "prompt_ids",
			ev:   EventValue{IDs: []string{"init", "code-review"}, IsInline: false},
			check: func(v interface{}) bool {
				ids, ok := v.([]string)
				return ok && len(ids) == 2 && ids[0] == "init" && ids[1] == "code-review"
			},
		},
		{
			name: "empty_ids",
			ev:   EventValue{IDs: []string{}, IsInline: false},
			check: func(v interface{}) bool {
				ids, ok := v.([]string)
				return ok && len(ids) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := tt.ev.MarshalYAML()
			if err != nil {
				t.Errorf("MarshalYAML() error = %v", err)
				return
			}
			if !tt.check(v) {
				t.Errorf("MarshalYAML() = %v, check failed", v)
			}
		})
	}
}

// TestEventHandler_UnmarshalYAML_AllEvents tests unmarshaling complete EventHandler
func TestEventHandler_UnmarshalYAML_AllEvents(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		check   func(EventHandler) bool
	}{
		{
			name: "all_events_present",
			yaml: `
on:
  success: [code-review]
  failure: [init, security-audit]
  timeout: "Check timed out"
`,
			check: func(eh EventHandler) bool {
				return len(eh.Success.IDs) == 1 && eh.Success.IDs[0] == "code-review" &&
					len(eh.Failure.IDs) == 2 && eh.Failure.IDs[0] == "init" &&
					eh.Timeout.Content == "Check timed out" && eh.Timeout.IsInline
			},
		},
		{
			name: "only_failure",
			yaml: `
on:
  failure: [init]
`,
			check: func(eh EventHandler) bool {
				return len(eh.Failure.IDs) == 1 && eh.Failure.IDs[0] == "init" &&
					len(eh.Success.IDs) == 0 && eh.Timeout.Content == ""
			},
		},
		{
			name: "mixed_inline_and_ids",
			yaml: `
on:
  success: [code-review]
  failure: "Fix the issue"
  timeout: [init]
`,
			check: func(eh EventHandler) bool {
				return len(eh.Success.IDs) == 1 && eh.Success.IDs[0] == "code-review" &&
					eh.Failure.Content == "Fix the issue" && eh.Failure.IsInline &&
					len(eh.Timeout.IDs) == 1 && eh.Timeout.IDs[0] == "init"
			},
		},
		{
			name: "empty_on_section",
			yaml: `
on:
`,
			check: func(eh EventHandler) bool {
				return len(eh.Success.IDs) == 0 && len(eh.Failure.IDs) == 0 && len(eh.Timeout.IDs) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var check Check
			fullYAML := `id: test
run: "true"
` + tt.yaml
			err := yaml.Unmarshal([]byte(fullYAML), &check)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !tt.check(check.On) {
				t.Errorf("UnmarshalYAML() check failed for %+v", check.On)
			}
		})
	}
}

// TestEventHandler_InCheckContext tests EventHandler integration in Check struct
func TestEventHandler_InCheckContext(t *testing.T) {
	yamlContent := `
checks:
  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    on:
      success: [code-review]
      failure: [init, security-audit]
      timeout: "Check timed out. Try again."
  - id: fmt
    run: gofmt -l .
    severity: error
`

	var cfg struct {
		Checks []Check
	}

	err := yaml.Unmarshal([]byte(yamlContent), &cfg)
	if err != nil {
		t.Fatalf("UnmarshalYAML() error = %v", err)
	}

	if len(cfg.Checks) != 2 {
		t.Errorf("Expected 2 checks, got %d", len(cfg.Checks))
	}

	// Check first check with event handler
	if cfg.Checks[0].ID != "vet" {
		t.Errorf("First check ID = %s, want vet", cfg.Checks[0].ID)
	}
	if len(cfg.Checks[0].On.Success.IDs) != 1 || cfg.Checks[0].On.Success.IDs[0] != "code-review" {
		t.Errorf("Success event IDs = %v, want [code-review]", cfg.Checks[0].On.Success.IDs)
	}
	if len(cfg.Checks[0].On.Failure.IDs) != 2 {
		t.Errorf("Failure event IDs length = %d, want 2", len(cfg.Checks[0].On.Failure.IDs))
	}
	if cfg.Checks[0].On.Timeout.Content != "Check timed out. Try again." {
		t.Errorf("Timeout event content = %q, want \"Check timed out. Try again.\"", cfg.Checks[0].On.Timeout.Content)
	}

	// Check second check without event handler
	if cfg.Checks[1].ID != "fmt" {
		t.Errorf("Second check ID = %s, want fmt", cfg.Checks[1].ID)
	}
	if len(cfg.Checks[1].On.Success.IDs) != 0 || cfg.Checks[1].On.Success.Content != "" {
		t.Errorf("Second check should have empty On handler")
	}
}

// TestEventHandler_EdgeCases tests edge cases and boundary conditions
func TestEventHandler_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		check   func(EventValue) bool
	}{
		{
			name: "null_value",
			yaml: `null`,
			// Null should unmarshal to zero EventValue
			check: func(ev EventValue) bool {
				return ev.Content == "" && len(ev.IDs) == 0
			},
		},
		{
			name: "numeric_string_treated_as_string",
			yaml: `"123"`,
			check: func(ev EventValue) bool {
				return ev.Content == "123" && ev.IsInline
			},
		},
		{
			name: "boolean_like_string",
			yaml: `"true"`,
			check: func(ev EventValue) bool {
				return ev.Content == "true" && ev.IsInline
			},
		},
		{
			name: "array_with_null_elements",
			yaml: `[init, null, code-review]`,
			// This tests YAML parser behavior - null becomes empty string
			check: func(ev EventValue) bool {
				return !ev.IsInline && len(ev.IDs) >= 1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ev EventValue
			err := yaml.Unmarshal([]byte(tt.yaml), &ev)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !tt.check(ev) {
				t.Errorf("UnmarshalYAML() check failed for %+v", ev)
			}
		})
	}
}

// TestEventHandler_PromptIDVsInlineDistinction tests the critical distinction between IDs and inline
func TestEventHandler_PromptIDVsInlineDistinction(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		isInline bool
		value    string
		ids      []string
	}{
		{
			name:     "bare_string_is_inline",
			yaml:     `failure: init`,
			isInline: true,
			value:    "init",
			ids:      nil,
		},
		{
			name:     "quoted_string_is_inline",
			yaml:     `failure: "init"`,
			isInline: true,
			value:    "init",
			ids:      nil,
		},
		{
			name:     "single_element_array_is_id",
			yaml:     `failure: [init]`,
			isInline: false,
			value:    "",
			ids:      []string{"init"},
		},
		{
			name:     "flow_style_single_is_id",
			yaml:     `failure: [init]`,
			isInline: false,
			value:    "",
			ids:      []string{"init"},
		},
		{
			name: "block_style_array_is_id",
			yaml: `failure:
  - init`,
			isInline: false,
			value:    "",
			ids:      []string{"init"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var check Check
			fullYAML := `id: test
run: "true"
on:
  ` + tt.yaml
			err := yaml.Unmarshal([]byte(fullYAML), &check)
			if err != nil {
				t.Fatalf("UnmarshalYAML() error = %v", err)
			}

			if check.On.Failure.IsInline != tt.isInline {
				t.Errorf("IsInline = %v, want %v", check.On.Failure.IsInline, tt.isInline)
			}
			if tt.isInline && check.On.Failure.Content != tt.value {
				t.Errorf("Content = %q, want %q", check.On.Failure.Content, tt.value)
			}
			if !tt.isInline {
				if len(check.On.Failure.IDs) != len(tt.ids) {
					t.Errorf("IDs length = %d, want %d", len(check.On.Failure.IDs), len(tt.ids))
				}
				for i, id := range check.On.Failure.IDs {
					if i >= len(tt.ids) || id != tt.ids[i] {
						t.Errorf("IDs[%d] = %q, want %q", i, id, tt.ids[i])
					}
				}
			}
		})
	}
}
