package assert

import (
	"unicode"
)

// Lexer tokenizes assertion expressions.
type Lexer struct {
	input   string
	pos     int  // current position in input
	readPos int  // current reading position (after current char)
	ch      byte // current char under examination
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position.
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // ASCII NUL signifies EOF
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

// peekChar returns the next character without advancing the position.
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token
	tok.Pos = l.pos

	switch l.ch {
	case '+':
		tok = Token{Type: TokenPlus, Literal: string(l.ch), Pos: l.pos}
	case '-':
		tok = Token{Type: TokenMinus, Literal: string(l.ch), Pos: l.pos}
	case '*':
		tok = Token{Type: TokenAsterisk, Literal: string(l.ch), Pos: l.pos}
	case '/':
		tok = Token{Type: TokenSlash, Literal: string(l.ch), Pos: l.pos}
	case '(':
		tok = Token{Type: TokenLParen, Literal: string(l.ch), Pos: l.pos}
	case ')':
		tok = Token{Type: TokenRParen, Literal: string(l.ch), Pos: l.pos}
	case '=':
		if l.peekChar() == '=' {
			pos := l.pos
			l.readChar()
			tok = Token{Type: TokenEq, Literal: "==", Pos: pos}
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch), Pos: l.pos}
		}
	case '!':
		if l.peekChar() == '=' {
			pos := l.pos
			l.readChar()
			tok = Token{Type: TokenNotEq, Literal: "!=", Pos: pos}
		} else {
			tok = Token{Type: TokenNot, Literal: string(l.ch), Pos: l.pos}
		}
	case '<':
		if l.peekChar() == '=' {
			pos := l.pos
			l.readChar()
			tok = Token{Type: TokenLTE, Literal: "<=", Pos: pos}
		} else {
			tok = Token{Type: TokenLT, Literal: string(l.ch), Pos: l.pos}
		}
	case '>':
		if l.peekChar() == '=' {
			pos := l.pos
			l.readChar()
			tok = Token{Type: TokenGTE, Literal: ">=", Pos: pos}
		} else {
			tok = Token{Type: TokenGT, Literal: string(l.ch), Pos: l.pos}
		}
	case '&':
		if l.peekChar() == '&' {
			pos := l.pos
			l.readChar()
			tok = Token{Type: TokenAnd, Literal: "&&", Pos: pos}
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch), Pos: l.pos}
		}
	case '|':
		if l.peekChar() == '|' {
			pos := l.pos
			l.readChar()
			tok = Token{Type: TokenOr, Literal: "||", Pos: pos}
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch), Pos: l.pos}
		}
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
		return tok
	case '\'':
		tok.Type = TokenString
		tok.Literal = l.readStringSingleQuote()
		return tok
	case 0:
		tok = Token{Type: TokenEOF, Literal: "", Pos: l.pos}
	default:
		if isLetter(l.ch) {
			tok.Pos = l.pos
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Pos = l.pos
			tok.Literal = l.readNumber()
			tok.Type = TokenNumber
			return tok
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch), Pos: l.pos}
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace consumes whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier (variable name).
func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// readNumber reads a numeric literal (int or float).
func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	// Handle decimal point for floats
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[pos:l.pos]
}

// readString reads a double-quoted string literal.
func (l *Lexer) readString() string {
	pos := l.pos + 1 // skip opening quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	str := l.input[pos:l.pos]
	l.readChar() // consume closing quote
	return str
}

// readStringSingleQuote reads a single-quoted string literal.
func (l *Lexer) readStringSingleQuote() string {
	pos := l.pos + 1 // skip opening quote
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	str := l.input[pos:l.pos]
	l.readChar() // consume closing quote
	return str
}

// lookupIdent checks if an identifier is a keyword (true/false).
func lookupIdent(ident string) TokenType {
	switch ident {
	case "true", "false":
		return TokenBool
	default:
		return TokenIdent
	}
}

// isLetter returns true if ch is a letter.
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

// isDigit returns true if ch is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
