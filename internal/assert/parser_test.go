package assert

import (
	"strings"
	"testing"
)

func TestParser_NewParser(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple identifier", "x"},
		{"simple number", "42"},
		{"empty string", ""},
		{"long expression", "x && y || z > 10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			if p.input != tt.input {
				t.Errorf("NewParser input mismatch: got %q, want %q", p.input, tt.input)
			}
		})
	}
}

func TestParser_FormatErrorBoundary(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		pos        int
		msg        string
		checkEqual bool
	}{
		{
			name:       "error at position 1 (boundary)",
			input:      "abc",
			pos:        1,
			msg:        "test error",
			checkEqual: true,
		},
		{
			name:       "error at position 2",
			input:      "abc",
			pos:        2,
			msg:        "test error",
			checkEqual: true,
		},
		{
			name:       "error at position beyond input length",
			input:      "ab",
			pos:        10,
			msg:        "test error",
			checkEqual: true,
		},
		{
			name:       "error at position 0",
			input:      "test",
			pos:        0,
			msg:        "error at start",
			checkEqual: true,
		},
		{
			name:       "empty input with error",
			input:      "",
			pos:        1,
			msg:        "error on empty",
			checkEqual: true,
		},
		{
			name:       "long input with error at end",
			input:      "very_long_expression_here",
			pos:        26,
			msg:        "error at end",
			checkEqual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			result := p.formatError(tt.pos, tt.msg)

			// Check that result contains the message
			if !strings.Contains(result, tt.msg) {
				t.Errorf("formatError did not contain message: got %q, want to contain %q", result, tt.msg)
			}

			// Check that result contains the input
			if !strings.Contains(result, tt.input) {
				t.Errorf("formatError did not contain input: got %q, want to contain %q", result, tt.input)
			}

			// Check that result contains the pointer line with "^"
			if !strings.Contains(result, "^") {
				t.Errorf("formatError did not contain pointer: got %q", result)
			}

			// Verify the format has newlines separating components
			parts := strings.Split(result, "\n")
			if len(parts) < 3 {
				t.Errorf("formatError should have at least 3 lines separated by newlines, got %d lines: %q", len(parts), result)
			}
		})
	}
}

func TestParser_FormatErrorPointerAccuracy(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		pos       int
		minSpaces int
		maxSpaces int
	}{
		{
			name:      "position 1 - exact start",
			input:     "abcdef",
			pos:       1,
			minSpaces: 2,
			maxSpaces: 2,
		},
		{
			name:      "position 2 - one space before",
			input:     "abcdef",
			pos:       2,
			minSpaces: 3,
			maxSpaces: 3,
		},
		{
			name:      "position 3 - two spaces before",
			input:     "abcdef",
			pos:       3,
			minSpaces: 4,
			maxSpaces: 4,
		},
		{
			name:      "position beyond input",
			input:     "abc",
			pos:       10,
			minSpaces: 5,
			maxSpaces: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			result := p.formatError(tt.pos, "test")

			// Extract the pointer line (last line)
			lines := strings.Split(result, "\n")
			if len(lines) < 3 {
				t.Fatalf("formatError should have 3+ lines, got %d", len(lines))
			}
			pointerLine := lines[len(lines)-1]

			// Count leading spaces before the caret
			spaceCount := 0
			for _, ch := range pointerLine {
				if ch == ' ' {
					spaceCount++
				} else {
					break
				}
			}

			// Verify caret position
			if !strings.Contains(pointerLine, "^") {
				t.Errorf("pointer line missing caret: %q", pointerLine)
			}

			if spaceCount < tt.minSpaces || spaceCount > tt.maxSpaces {
				t.Errorf("pointer spacing wrong: got %d spaces, want %d-%d, line: %q", spaceCount, tt.minSpaces, tt.maxSpaces, pointerLine)
			}
		})
	}
}

func TestParser_ParseSimpleExpression(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"single number", "42", false},
		{"single identifier", "x", false},
		{"single string", `"hello"`, false},
		{"true literal", "true", false},
		{"false literal", "false", false},
		{"parenthesized expr", "(x)", false},
		{"unary not", "!x", false},
		{"unary minus", "-42", false},
		{"invalid token at start", "@", true},
		{"missing closing paren", "(x", true},
		{"unexpected token after paren", "(x]", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			expr, err := p.Parse()

			if tt.shouldErr && err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if !tt.shouldErr && expr == nil {
				t.Errorf("Parse(%q) returned nil expression", tt.input)
			}
		})
	}
}

func TestParser_BinaryOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"or operator", "x || y", false},
		{"and operator", "x && y", false},
		{"comparison equal", "x == y", false},
		{"comparison not equal", "x != y", false},
		{"comparison less than", "x < y", false},
		{"comparison less equal", "x <= y", false},
		{"comparison greater", "x > y", false},
		{"comparison greater equal", "x >= y", false},
		{"addition", "x + y", false},
		{"subtraction", "x - y", false},
		{"multiplication", "x * y", false},
		{"division", "x / y", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			expr, err := p.Parse()

			if tt.shouldErr && err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if !tt.shouldErr && expr == nil {
				t.Errorf("Parse(%q) returned nil expression", tt.input)
			}

			// Check that binary expressions are parsed correctly
			if !tt.shouldErr {
				if binExpr, ok := expr.(*BinaryExpr); ok {
					if binExpr.Left == nil || binExpr.Right == nil {
						t.Errorf("Parse(%q) created incomplete binary expression", tt.input)
					}
				}
			}
		})
	}
}

func TestParser_PrecedenceLoopBoundary(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"single operator", "x + y", false},
		{"two operators same precedence", "x + y + z", false},
		{"two operators different precedence", "x + y * z", false},
		{"multiple operators", "a || b && c == d < e + f * g", false},
		{"operators at boundary", "x + y", false},
		{"chain of comparisons", "x > y < z", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			expr, err := p.Parse()

			if tt.shouldErr && err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if !tt.shouldErr && expr == nil {
				t.Errorf("Parse(%q) returned nil expression", tt.input)
			}
		})
	}
}

func TestParser_PrecedenceFunction(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  int
	}{
		{TokenOr, PrecOr},
		{TokenAnd, PrecAnd},
		{TokenEq, PrecCompare},
		{TokenNotEq, PrecCompare},
		{TokenLT, PrecCompare},
		{TokenLTE, PrecCompare},
		{TokenGT, PrecCompare},
		{TokenGTE, PrecCompare},
		{TokenPlus, PrecSum},
		{TokenMinus, PrecSum},
		{TokenAsterisk, PrecProduct},
		{TokenSlash, PrecProduct},
		{TokenNumber, PrecLowest},
		{TokenString, PrecLowest},
		{TokenIdent, PrecLowest},
		{TokenEOF, PrecLowest},
	}

	for _, tt := range tests {
		t.Run(tt.tokenType.String(), func(t *testing.T) {
			result := precedence(tt.tokenType)
			if result != tt.expected {
				t.Errorf("precedence(%v) = %d, want %d", tt.tokenType, result, tt.expected)
			}
		})
	}
}

func TestParser_PrecedenceOrderingBoundary(t *testing.T) {
	// Verify strict precedence ordering (each should be distinct and ordered)
	precedences := []int{
		PrecLowest,
		PrecOr,
		PrecAnd,
		PrecCompare,
		PrecSum,
		PrecProduct,
		PrecUnary,
		PrecPrimary,
	}

	// Check that precedences are strictly increasing
	for i := 1; i < len(precedences); i++ {
		if precedences[i] <= precedences[i-1] {
			t.Errorf("precedence ordering violated at index %d: %d not > %d", i, precedences[i], precedences[i-1])
		}
	}
}

func TestParser_NextTokenBoundary(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens int
	}{
		{"empty string", "", 1},         // EOF only
		{"single token", "42", 2},       // 42, EOF
		{"two tokens", "x y", 2},        // x (ident), EOF (y is part of ident "xy" or separate)
		{"operator", "!", 2},            // !, EOF
		{"full expression", "x + y", 4}, // x, +, y, EOF
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			count := 0
			for p.cur.Type != TokenEOF && count < 20 {
				p.nextToken()
				count++
			}
			if p.cur.Type != TokenEOF {
				t.Errorf("nextToken did not reach EOF for input %q", tt.input)
			}
		})
	}
}

func TestParser_ParseParenthesisBoundary(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"simple paren", "(x)", false},
		{"nested paren", "((x))", false},
		{"triple nested", "(((x)))", false},
		{"paren with binary", "(x + y)", false},
		{"multiple parens", "(x) && (y)", false},
		{"unclosed paren", "(x", true},
		{"wrong closing bracket", "(x]", true},
		{"nested unclosed", "((x)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			expr, err := p.Parse()

			if tt.shouldErr && err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if !tt.shouldErr && expr == nil {
				t.Errorf("Parse(%q) returned nil expression", tt.input)
			}
		})
	}
}

func TestParser_UnaryOperatorsBoundary(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"unary not", "!x", false},
		{"unary minus", "-42", false},
		{"double unary", "!!x", false},
		{"unary mixed", "!-x", false},
		{"minus mixed", "-!x", false},
		{"unary in binary", "!x && y", false},
		{"unary with paren", "!(x)", false},
		{"unary without operand", "!", true},
		{"unary without operand minus", "-", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			_, err := p.Parse()

			if tt.shouldErr && err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestParser_LiteralParsingBoundary(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		exprType  string
	}{
		{"number int", "42", false, "*assert.NumberLit"},
		{"number zero", "0", false, "*assert.NumberLit"},
		{"number large", "999999999", false, "*assert.NumberLit"},
		{"string simple", `"hello"`, false, "*assert.StringLit"},
		{"string empty", `""`, false, "*assert.StringLit"},
		{"bool true", "true", false, "*assert.BoolLit"},
		{"bool false", "false", false, "*assert.BoolLit"},
		{"identifier simple", "x", false, "*assert.Ident"},
		{"identifier long", "very_long_identifier_name", false, "*assert.Ident"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			_, err := p.Parse()

			if tt.shouldErr && err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

func TestParser_ComplexExpressionBoundary(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"complex chain", "a || b && c == d < e + f * g / h", false},
		{"many operators", "a + b - c * d / e < f > g == h != i && j || k", false},
		{"deeply nested", "((((a))))", false},
		{"complex with unary", "!a || -b && c == d", false},
		{"all binary operators", "a || b && c == d != e < f <= g > h >= i + j - k * l / m", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			expr, err := p.Parse()

			if tt.shouldErr && err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if !tt.shouldErr && expr == nil {
				t.Errorf("Parse(%q) returned nil expression", tt.input)
			}
		})
	}
}

func TestParser_ErrorMessageQuality(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:          "invalid token error",
			input:         "@",
			shouldContain: []string{"unexpected token", "@", "position"},
		},
		{
			name:          "missing paren error",
			input:         "(x",
			shouldContain: []string{"expected", ")", "position"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			_, err := p.Parse()

			if err == nil {
				t.Fatalf("Parse(%q) expected error", tt.input)
			}

			errMsg := err.Error()

			for _, substr := range tt.shouldContain {
				if !strings.Contains(errMsg, substr) {
					t.Errorf("error message should contain %q, got: %q", substr, errMsg)
				}
			}

			for _, substr := range tt.shouldNotContain {
				if strings.Contains(errMsg, substr) {
					t.Errorf("error message should not contain %q, got: %q", substr, errMsg)
				}
			}
		})
	}
}
