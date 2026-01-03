package assert

import (
	"strings"
	"testing"
)

func TestEvaluator_EmptyExpr(t *testing.T) {
	e := New()
	result, err := e.Eval("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Error("empty expression should return true")
	}
}

func TestEvaluator_Literals(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want bool
	}{
		{"true literal", "true", true},
		{"false literal", "false", false},
		{"number non-zero", "42", true},
		{"number zero", "0", false},
		{"float non-zero", "3.14", true},
		{"float zero", "0.0", false},
		{"string non-empty", `"hello"`, true},
		{"string empty", `""`, false},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, result, tt.want)
			}
		})
	}
}

func TestEvaluator_Variables(t *testing.T) {
	tests := []struct {
		name string
		expr string
		vars map[string]string
		want bool
	}{
		{
			name: "variable exists and true",
			expr: "coverage",
			vars: map[string]string{"coverage": "80"},
			want: true,
		},
		{
			name: "variable exists and zero",
			expr: "count",
			vars: map[string]string{"count": "0"},
			want: false,
		},
		{
			name: "undefined variable",
			expr: "undefined",
			vars: map[string]string{},
			want: false,
		},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, tt.vars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, result, tt.want)
			}
		})
	}
}

func TestEvaluator_Comparisons(t *testing.T) {
	tests := []struct {
		name string
		expr string
		vars map[string]string
		want bool
	}{
		// Numeric comparisons
		{"num eq true", "10 == 10", nil, true},
		{"num eq false", "10 == 20", nil, false},
		{"num neq true", "10 != 20", nil, true},
		{"num neq false", "10 != 10", nil, false},
		{"num lt true", "10 < 20", nil, true},
		{"num lt false", "20 < 10", nil, false},
		{"num lte true eq", "10 <= 10", nil, true},
		{"num lte true lt", "10 <= 20", nil, true},
		{"num lte false", "20 <= 10", nil, false},
		{"num gt true", "20 > 10", nil, true},
		{"num gt false", "10 > 20", nil, false},
		{"num gte true eq", "10 >= 10", nil, true},
		{"num gte true gt", "20 >= 10", nil, true},
		{"num gte false", "10 >= 20", nil, false},

		// Float comparisons
		{"float comparison", "3.14 > 3.0", nil, true},
		{"float eq", "3.14 == 3.14", nil, true},

		// String comparisons
		{"string eq", `"hello" == "hello"`, nil, true},
		{"string neq", `"hello" != "world"`, nil, true},
		{"string lt", `"abc" < "def"`, nil, true},
		{"string gt", `"xyz" > "abc"`, nil, true},

		// Variables with auto-coercion
		{
			name: "variable numeric comparison",
			expr: "coverage >= 80",
			vars: map[string]string{"coverage": "85.5"},
			want: true,
		},
		{
			name: "variable numeric comparison fail",
			expr: "coverage >= 80",
			vars: map[string]string{"coverage": "75"},
			want: false,
		},
		{
			name: "string variable comparison",
			expr: `status == "ok"`,
			vars: map[string]string{"status": "ok"},
			want: true,
		},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, tt.vars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, result, tt.want)
			}
		})
	}
}

func TestEvaluator_Arithmetic(t *testing.T) {
	tests := []struct {
		name string
		expr string
		vars map[string]string
		want bool
	}{
		{"addition", "10 + 5 == 15", nil, true},
		{"subtraction", "10 - 5 == 5", nil, true},
		{"multiplication", "10 * 5 == 50", nil, true},
		{"division", "10 / 5 == 2", nil, true},
		{"complex expression", "(10 + 5) * 2 == 30", nil, true},
		{"precedence", "10 + 5 * 2 == 20", nil, true},
		{
			name: "with variables",
			expr: "total / count >= 80",
			vars: map[string]string{"total": "160", "count": "2"},
			want: true,
		},
		{
			name: "coverage calculation",
			expr: "(covered / total) * 100 >= 80",
			vars: map[string]string{"covered": "85", "total": "100"},
			want: true,
		},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, tt.vars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, result, tt.want)
			}
		})
	}
}

func TestEvaluator_LogicalOperators(t *testing.T) {
	tests := []struct {
		name string
		expr string
		vars map[string]string
		want bool
	}{
		// AND operator
		{"and both true", "true && true", nil, true},
		{"and left false", "false && true", nil, false},
		{"and right false", "true && false", nil, false},
		{"and both false", "false && false", nil, false},

		// OR operator
		{"or both true", "true || true", nil, true},
		{"or left true", "true || false", nil, true},
		{"or right true", "false || true", nil, true},
		{"or both false", "false || false", nil, false},

		// NOT operator
		{"not true", "!true", nil, false},
		{"not false", "!false", nil, true},

		// Complex expressions
		{
			name: "and with comparisons",
			expr: "coverage >= 80 && tests > 0",
			vars: map[string]string{"coverage": "85", "tests": "10"},
			want: true,
		},
		{
			name: "or with comparisons",
			expr: "coverage >= 80 || warnings == 0",
			vars: map[string]string{"coverage": "70", "warnings": "0"},
			want: true,
		},
		{
			name: "not with comparison",
			expr: "!(errors > 0)",
			vars: map[string]string{"errors": "0"},
			want: true,
		},
		{
			name: "complex logical",
			expr: "(coverage >= 80 && tests > 0) || force",
			vars: map[string]string{"coverage": "70", "tests": "10", "force": "true"},
			want: true,
		},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, tt.vars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, result, tt.want)
			}
		})
	}
}

func TestEvaluator_UnaryMinus(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want bool
	}{
		{"negative number", "-5 < 0", true},
		{"negative in expr", "10 + -5 == 5", true},
		{"double negative", "--5 == 5", true},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, result, tt.want)
			}
		})
	}
}

func TestEvaluator_Parentheses(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want bool
	}{
		{"override precedence", "(10 + 5) * 2 == 30", true},
		{"nested parens", "((10 + 5) * 2) == 30", true},
		{"logical grouping", "(true || false) && true", true},
		{"comparison grouping", "(10 > 5) == true", true},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, result, tt.want)
			}
		})
	}
}

func TestEvaluator_ParseErrors(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"missing operand", "10 +"},
		{"unclosed paren", "(10 + 5"},
		{"missing paren content", "()"},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.Eval(tt.expr, nil)
			if err == nil {
				t.Errorf("expected error for %q, got nil", tt.expr)
			}
		})
	}
}

func TestEvaluator_ParseErrorContext(t *testing.T) {
	// Tests that verify error messages include proper context with pointer
	tests := []struct {
		name        string
		expr        string
		wantContain string
	}{
		{
			name:        "error at position 1",
			expr:        "@",
			wantContain: "^",
		},
		{
			name:        "error shows input",
			expr:        "10 + @",
			wantContain: "10 + @",
		},
		{
			name:        "long expression error",
			expr:        "a && b && c && @",
			wantContain: "a && b && c && @",
		},
		{
			name:        "error at start",
			expr:        ")",
			wantContain: "^",
		},
		{
			name:        "error in middle",
			expr:        "10 + + 20",
			wantContain: "^",
		},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.Eval(tt.expr, nil)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tt.expr)
			}
			errStr := err.Error()
			if !strings.Contains(errStr, tt.wantContain) {
				t.Errorf("error %q should contain %q", errStr, tt.wantContain)
			}
		})
	}
}

// TestParser_FormatErrorPointerPosition kills mutations at parser.go:37
// that change the loop conditions for building the error pointer.
// The pointer should point to the exact position of the error.
// The loop builds `pos-1` spaces then adds ^, so position N has N-1 spaces.
func TestParser_FormatErrorPointerPosition(t *testing.T) {
	tests := []struct {
		name           string
		expr           string
		wantPointerPos int // position of ^ relative to start of expression line (pos-1)
	}{
		{
			name:           "error at position 1",
			expr:           "@",
			wantPointerPos: 0, // pos=1, so 0 spaces before ^
		},
		{
			name:           "error at position 6",
			expr:           "10 + @",
			wantPointerPos: 4, // pos=5, so 4 spaces before ^ (pointing at @)
		},
		{
			name:           "error in middle of long expr",
			expr:           "a && b && @",
			wantPointerPos: 9, // pos=10, so 9 spaces before ^
		},
		{
			name:           "error at position 3",
			expr:           "a +@",
			wantPointerPos: 2, // pos=3, so 2 spaces before ^ (pointing at @)
		},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.Eval(tt.expr, nil)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tt.expr)
			}
			errStr := err.Error()

			// Find the line with the pointer (^)
			lines := strings.Split(errStr, "\n")
			var pointerLine string
			for _, line := range lines {
				if strings.Contains(line, "^") && !strings.Contains(line, "unexpected") {
					pointerLine = line
					break
				}
			}
			if pointerLine == "" {
				t.Fatalf("no pointer line found in error: %q", errStr)
			}

			// The pointer line is indented with "  " (2 spaces), then spaces, then ^
			// Count spaces before ^ (excluding the leading "  " indent)
			pointerLine = strings.TrimPrefix(pointerLine, "  ")
			pointerPos := strings.Index(pointerLine, "^")
			if pointerPos != tt.wantPointerPos {
				t.Errorf("pointer position = %d, want %d\nerror: %s", pointerPos, tt.wantPointerPos, errStr)
			}
		})
	}
}

func TestEvaluator_ShortCircuit(t *testing.T) {
	// Test that short-circuit evaluation works properly
	// If short-circuit is working, undefined variable access won't cause issues
	// when the result is already determined

	e := New()

	// AND short-circuits on false left operand
	result, err := e.Eval("false && undefined_var", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Error("expected false for short-circuit AND")
	}

	// OR short-circuits on true left operand
	result, err = e.Eval("true || undefined_var", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Error("expected true for short-circuit OR")
	}
}

func TestEvaluator_RealWorldExamples(t *testing.T) {
	tests := []struct {
		name string
		expr string
		vars map[string]string
		want bool
	}{
		{
			name: "test coverage check",
			expr: "coverage >= 80",
			vars: map[string]string{"coverage": "85.5"},
			want: true,
		},
		{
			name: "go vet clean",
			expr: "exit_code == 0",
			vars: map[string]string{"exit_code": "0"},
			want: true,
		},
		{
			name: "linter warnings below threshold",
			expr: "warnings < 10",
			vars: map[string]string{"warnings": "5"},
			want: true,
		},
		{
			name: "binary size under limit",
			expr: "size_mb <= 50",
			vars: map[string]string{"size_mb": "42.3"},
			want: true,
		},
		{
			name: "complex ci check",
			expr: "(coverage >= 80 && tests_passed > 0) && lint_errors == 0",
			vars: map[string]string{
				"coverage":     "85",
				"tests_passed": "42",
				"lint_errors":  "0",
			},
			want: true,
		},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, tt.vars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, result, tt.want)
			}
		})
	}
}

func TestEvaluator_DivisionByZero(t *testing.T) {
	e := New()

	// Division by zero should return 0, not error
	result, err := e.Eval("10 / 0 == 0", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Error("expected division by zero to return 0")
	}
}

func TestEvaluator_SingleQuoteStrings(t *testing.T) {
	e := New()

	result, err := e.Eval(`'hello' == 'hello'`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Error("expected single-quoted string comparison to work")
	}
}

// TestEvaluator_LessThanBoundary kills mutation at eval.go:215
// that changes `cmp < 0` to `cmp <= 0` in the less-than comparison.
// This test ensures that equal values return false for <.
func TestEvaluator_LessThanBoundary(t *testing.T) {
	e := New()

	// When values are equal, < should return false
	result, err := e.Eval("10 < 10", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Error("10 < 10 should be false, not true (boundary condition)")
	}

	// Also test with floats to be thorough
	result, err = e.Eval("3.14 < 3.14", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Error("3.14 < 3.14 should be false, not true (boundary condition)")
	}

	// Test with variables
	result, err = e.Eval("x < x", map[string]string{"x": "42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Error("x < x should be false when x equals x (boundary condition)")
	}
}

// TestFormatFloat_Precision kills mutations at eval.go:276
// that change the precision parameter from -1 to other values.
// The -1 precision means "smallest number of digits necessary".
func TestFormatFloat_Precision(t *testing.T) {
	// Test that formatFloat preserves full precision for fractional numbers
	// If precision is changed from -1 to 1, "3.14" would become "3.1"
	tests := []struct {
		name     string
		expr     string
		vars     map[string]string
		wantTrue bool
	}{
		{
			name:     "preserves multiple decimal places",
			expr:     "x == 3.14159",
			vars:     map[string]string{"x": "3.14159"},
			wantTrue: true,
		},
		{
			name:     "arithmetic preserves precision",
			expr:     "1.5 + 1.25 == 2.75",
			wantTrue: true,
		},
		{
			name:     "division result precision",
			expr:     "7 / 4 == 1.75",
			wantTrue: true,
		},
		{
			name:     "negative with precision",
			expr:     "-3.14159 == -3.14159",
			wantTrue: true,
		},
	}

	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Eval(tt.expr, tt.vars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.wantTrue {
				t.Errorf("Eval(%q) = %v, want %v (precision test)", tt.expr, result, tt.wantTrue)
			}
		})
	}
}

// Value type tests
func TestValue_String(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"42", "42"},
		{"", ""},
		{"3.14", "3.14"},
		{"true", "true"},
	}

	for _, tt := range tests {
		v := NewValue(tt.input)
		if got := v.String(); got != tt.expected {
			t.Errorf("Value(%q).String() = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestValue_IsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"42", true},
		{"3.14", true},
		{"-10", true},
		{"0", true},
		{"1e10", true},
		{"hello", false},
		{"", false},
		{"true", false},
		{"12abc", false},
	}

	for _, tt := range tests {
		v := NewValue(tt.input)
		if got := v.IsNumeric(); got != tt.expected {
			t.Errorf("Value(%q).IsNumeric() = %v, want %v", tt.input, got, tt.expected)
		}
	}
}
