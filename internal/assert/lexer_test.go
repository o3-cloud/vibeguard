package assert

import (
	"testing"
)

func TestLexer_Tokens(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []Token
	}{
		{
			name:  "number",
			input: "42",
			tokens: []Token{
				{Type: TokenNumber, Literal: "42", Pos: 0},
				{Type: TokenEOF, Literal: "", Pos: 2},
			},
		},
		{
			name:  "float",
			input: "3.14",
			tokens: []Token{
				{Type: TokenNumber, Literal: "3.14", Pos: 0},
				{Type: TokenEOF, Literal: "", Pos: 4},
			},
		},
		{
			name:  "identifier",
			input: "coverage",
			tokens: []Token{
				{Type: TokenIdent, Literal: "coverage", Pos: 0},
				{Type: TokenEOF, Literal: "", Pos: 8},
			},
		},
		{
			name:  "identifier with underscore",
			input: "test_coverage",
			tokens: []Token{
				{Type: TokenIdent, Literal: "test_coverage", Pos: 0},
				{Type: TokenEOF, Literal: "", Pos: 13},
			},
		},
		{
			name:  "bool true",
			input: "true",
			tokens: []Token{
				{Type: TokenBool, Literal: "true", Pos: 0},
				{Type: TokenEOF, Literal: "", Pos: 4},
			},
		},
		{
			name:  "bool false",
			input: "false",
			tokens: []Token{
				{Type: TokenBool, Literal: "false", Pos: 0},
				{Type: TokenEOF, Literal: "", Pos: 5},
			},
		},
		{
			name:  "double quoted string",
			input: `"hello world"`,
			tokens: []Token{
				{Type: TokenString, Literal: "hello world", Pos: 0},
				{Type: TokenEOF, Literal: "", Pos: 13},
			},
		},
		{
			name:  "single quoted string",
			input: `'hello world'`,
			tokens: []Token{
				{Type: TokenString, Literal: "hello world", Pos: 0},
				{Type: TokenEOF, Literal: "", Pos: 13},
			},
		},
		{
			name:  "comparison operators",
			input: "== != < <= > >=",
			tokens: []Token{
				{Type: TokenEq, Literal: "==", Pos: 0},
				{Type: TokenNotEq, Literal: "!=", Pos: 3},
				{Type: TokenLT, Literal: "<", Pos: 6},
				{Type: TokenLTE, Literal: "<=", Pos: 8},
				{Type: TokenGT, Literal: ">", Pos: 11},
				{Type: TokenGTE, Literal: ">=", Pos: 13},
				{Type: TokenEOF, Literal: "", Pos: 15},
			},
		},
		{
			name:  "arithmetic operators",
			input: "+ - * /",
			tokens: []Token{
				{Type: TokenPlus, Literal: "+", Pos: 0},
				{Type: TokenMinus, Literal: "-", Pos: 2},
				{Type: TokenAsterisk, Literal: "*", Pos: 4},
				{Type: TokenSlash, Literal: "/", Pos: 6},
				{Type: TokenEOF, Literal: "", Pos: 7},
			},
		},
		{
			name:  "logical operators",
			input: "&& || !",
			tokens: []Token{
				{Type: TokenAnd, Literal: "&&", Pos: 0},
				{Type: TokenOr, Literal: "||", Pos: 3},
				{Type: TokenNot, Literal: "!", Pos: 6},
				{Type: TokenEOF, Literal: "", Pos: 7},
			},
		},
		{
			name:  "parentheses",
			input: "()",
			tokens: []Token{
				{Type: TokenLParen, Literal: "(", Pos: 0},
				{Type: TokenRParen, Literal: ")", Pos: 1},
				{Type: TokenEOF, Literal: "", Pos: 2},
			},
		},
		{
			name:  "complex expression",
			input: "coverage >= 80 && tests > 0",
			tokens: []Token{
				{Type: TokenIdent, Literal: "coverage", Pos: 0},
				{Type: TokenGTE, Literal: ">=", Pos: 9},
				{Type: TokenNumber, Literal: "80", Pos: 12},
				{Type: TokenAnd, Literal: "&&", Pos: 15},
				{Type: TokenIdent, Literal: "tests", Pos: 18},
				{Type: TokenGT, Literal: ">", Pos: 24},
				{Type: TokenNumber, Literal: "0", Pos: 26},
				{Type: TokenEOF, Literal: "", Pos: 27},
			},
		},
		{
			name:  "whitespace handling",
			input: "  42  +  10  ",
			tokens: []Token{
				{Type: TokenNumber, Literal: "42", Pos: 2},
				{Type: TokenPlus, Literal: "+", Pos: 6},
				{Type: TokenNumber, Literal: "10", Pos: 9},
				{Type: TokenEOF, Literal: "", Pos: 13},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			for i, expected := range tt.tokens {
				tok := l.NextToken()
				if tok.Type != expected.Type {
					t.Errorf("token[%d] type: got %v, want %v", i, tok.Type, expected.Type)
				}
				if tok.Literal != expected.Literal {
					t.Errorf("token[%d] literal: got %q, want %q", i, tok.Literal, expected.Literal)
				}
				if tok.Pos != expected.Pos {
					t.Errorf("token[%d] pos: got %d, want %d", i, tok.Pos, expected.Pos)
				}
			}
		})
	}
}

func TestLexer_IllegalTokens(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"single ampersand", "&"},
		{"single pipe", "|"},
		{"single equals", "="},
		{"at symbol", "@"},
		{"hash", "#"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.NextToken()
			if tok.Type != TokenIllegal {
				t.Errorf("expected ILLEGAL token for %q, got %v", tt.input, tok.Type)
			}
		})
	}
}

func TestLexer_DigitBoundary(t *testing.T) {
	// Test digit detection at boundaries to ensure '0' through '9' are recognized
	// and characters just outside that range are not.
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantLit  string
	}{
		// Test that '9' is recognized as part of a number
		{"digit 9 recognized", "9", TokenNumber, "9"},
		{"digit 0 recognized", "0", TokenNumber, "0"},
		// Numbers with all digits 0-9
		{"all digits", "1234567890", TokenNumber, "1234567890"},
		// Colon ':' (ASCII 58, right after '9' which is 57) should not be a digit
		{"colon after number", "9:", TokenNumber, "9"},
		// Slash '/' (ASCII 47, right before '0' which is 48) should not be a digit
		{"number then slash", "0/", TokenNumber, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.NextToken()
			if tok.Type != tt.wantType {
				t.Errorf("got type %v, want %v", tok.Type, tt.wantType)
			}
			if tok.Literal != tt.wantLit {
				t.Errorf("got literal %q, want %q", tok.Literal, tt.wantLit)
			}
		})
	}
}

func TestTokenType_String(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{TokenEOF, "EOF"},
		{TokenIllegal, "ILLEGAL"},
		{TokenIdent, "IDENT"},
		{TokenNumber, "NUMBER"},
		{TokenString, "STRING"},
		{TokenBool, "BOOL"},
		{TokenPlus, "+"},
		{TokenMinus, "-"},
		{TokenAsterisk, "*"},
		{TokenSlash, "/"},
		{TokenEq, "=="},
		{TokenNotEq, "!="},
		{TokenLT, "<"},
		{TokenGT, ">"},
		{TokenLTE, "<="},
		{TokenGTE, ">="},
		{TokenAnd, "&&"},
		{TokenOr, "||"},
		{TokenNot, "!"},
		{TokenLParen, "("},
		{TokenRParen, ")"},
		{TokenType(999), "UNKNOWN"}, // Test default case
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.tokenType.String()
			if got != tt.expected {
				t.Errorf("TokenType(%d).String() = %q, want %q", tt.tokenType, got, tt.expected)
			}
		})
	}
}
