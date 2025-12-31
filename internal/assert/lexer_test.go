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
