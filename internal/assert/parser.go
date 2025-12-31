package assert

import (
	"fmt"
)

// Parser parses assertion expressions into an AST.
type Parser struct {
	lexer *Lexer
	cur   Token
	peek  Token
}

// NewParser creates a new Parser for the given input.
func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	// Initialize cur and peek tokens
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances to the next token.
func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.lexer.NextToken()
}

// Parse parses the input and returns the AST.
func (p *Parser) Parse() (Expr, error) {
	return p.parseExpr(PrecLowest)
}

// Precedence levels (lowest to highest).
const (
	PrecLowest  = iota
	PrecOr      // ||
	PrecAnd     // &&
	PrecCompare // ==, !=, <, <=, >, >=
	PrecSum     // +, -
	PrecProduct // *, /
	PrecUnary   // !, -
	PrecPrimary // literals, identifiers, parentheses
)

// precedence returns the precedence level for a token type.
func precedence(t TokenType) int {
	switch t {
	case TokenOr:
		return PrecOr
	case TokenAnd:
		return PrecAnd
	case TokenEq, TokenNotEq, TokenLT, TokenLTE, TokenGT, TokenGTE:
		return PrecCompare
	case TokenPlus, TokenMinus:
		return PrecSum
	case TokenAsterisk, TokenSlash:
		return PrecProduct
	default:
		return PrecLowest
	}
}

// parseExpr parses an expression with the given precedence level.
func (p *Parser) parseExpr(prec int) (Expr, error) {
	// Parse prefix expression (unary operators or primary expressions)
	left, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}

	// Parse infix expressions (binary operators) with proper precedence
	for p.cur.Type != TokenEOF && prec < precedence(p.cur.Type) {
		left, err = p.parseInfix(left)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

// parsePrefix handles prefix expressions (unary ops and primary values).
func (p *Parser) parsePrefix() (Expr, error) {
	switch p.cur.Type {
	case TokenNumber:
		return p.parseNumber()
	case TokenString:
		return p.parseString()
	case TokenBool:
		return p.parseBool()
	case TokenIdent:
		return p.parseIdent()
	case TokenLParen:
		return p.parseParen()
	case TokenNot:
		return p.parseUnary()
	case TokenMinus:
		return p.parseUnary()
	default:
		return nil, fmt.Errorf("unexpected token %q at position %d", p.cur.Literal, p.cur.Pos)
	}
}

// parseNumber parses a numeric literal.
func (p *Parser) parseNumber() (Expr, error) {
	lit := &NumberLit{Value: p.cur.Literal}
	p.nextToken()
	return lit, nil
}

// parseString parses a string literal.
func (p *Parser) parseString() (Expr, error) {
	lit := &StringLit{Value: p.cur.Literal}
	p.nextToken()
	return lit, nil
}

// parseBool parses a boolean literal.
func (p *Parser) parseBool() (Expr, error) {
	lit := &BoolLit{Value: p.cur.Literal == "true"}
	p.nextToken()
	return lit, nil
}

// parseIdent parses an identifier.
func (p *Parser) parseIdent() (Expr, error) {
	ident := &Ident{Name: p.cur.Literal}
	p.nextToken()
	return ident, nil
}

// parseParen parses a parenthesized expression.
func (p *Parser) parseParen() (Expr, error) {
	p.nextToken() // consume '('
	inner, err := p.parseExpr(PrecLowest)
	if err != nil {
		return nil, err
	}
	if p.cur.Type != TokenRParen {
		return nil, fmt.Errorf("expected ')' at position %d, got %q", p.cur.Pos, p.cur.Literal)
	}
	p.nextToken() // consume ')'
	return &ParenExpr{Inner: inner}, nil
}

// parseUnary parses a unary expression.
func (p *Parser) parseUnary() (Expr, error) {
	op := p.cur.Type
	p.nextToken()
	right, err := p.parseExpr(PrecUnary)
	if err != nil {
		return nil, err
	}
	return &UnaryExpr{Op: op, Right: right}, nil
}

// parseInfix parses a binary (infix) expression.
func (p *Parser) parseInfix(left Expr) (Expr, error) {
	op := p.cur.Type
	prec := precedence(op)
	p.nextToken()
	right, err := p.parseExpr(prec)
	if err != nil {
		return nil, err
	}
	return &BinaryExpr{Left: left, Op: op, Right: right}, nil
}
