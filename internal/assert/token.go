// Package assert provides assertion expression evaluation.
package assert

// TokenType represents the type of a token.
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenIllegal

	// Literals
	TokenIdent  // variable name
	TokenNumber // numeric literal (int or float)
	TokenString // quoted string literal
	TokenBool   // true or false

	// Operators
	TokenPlus     // +
	TokenMinus    // -
	TokenAsterisk // *
	TokenSlash    // /

	// Comparison operators
	TokenEq    // ==
	TokenNotEq // !=
	TokenLT    // <
	TokenLTE   // <=
	TokenGT    // >
	TokenGTE   // >=

	// Logical operators
	TokenAnd // &&
	TokenOr  // ||
	TokenNot // !

	// Delimiters
	TokenLParen // (
	TokenRParen // )
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Pos     int // position in input
}

// String returns a string representation of the token type.
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenIllegal:
		return "ILLEGAL"
	case TokenIdent:
		return "IDENT"
	case TokenNumber:
		return "NUMBER"
	case TokenString:
		return "STRING"
	case TokenBool:
		return "BOOL"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenAsterisk:
		return "*"
	case TokenSlash:
		return "/"
	case TokenEq:
		return "=="
	case TokenNotEq:
		return "!="
	case TokenLT:
		return "<"
	case TokenLTE:
		return "<="
	case TokenGT:
		return ">"
	case TokenGTE:
		return ">="
	case TokenAnd:
		return "&&"
	case TokenOr:
		return "||"
	case TokenNot:
		return "!"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	default:
		return "UNKNOWN"
	}
}
