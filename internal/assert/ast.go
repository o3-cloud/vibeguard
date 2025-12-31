package assert

// Node represents a node in the AST.
type Node interface {
	node()
}

// Expr represents an expression node.
type Expr interface {
	Node
	expr()
}

// NumberLit represents a numeric literal.
type NumberLit struct {
	Value string
}

func (*NumberLit) node() {}
func (*NumberLit) expr() {}

// StringLit represents a string literal.
type StringLit struct {
	Value string
}

func (*StringLit) node() {}
func (*StringLit) expr() {}

// BoolLit represents a boolean literal.
type BoolLit struct {
	Value bool
}

func (*BoolLit) node() {}
func (*BoolLit) expr() {}

// Ident represents an identifier (variable reference).
type Ident struct {
	Name string
}

func (*Ident) node() {}
func (*Ident) expr() {}

// UnaryExpr represents a unary expression (e.g., !x, -x).
type UnaryExpr struct {
	Op    TokenType
	Right Expr
}

func (*UnaryExpr) node() {}
func (*UnaryExpr) expr() {}

// BinaryExpr represents a binary expression (e.g., a + b, x == y).
type BinaryExpr struct {
	Left  Expr
	Op    TokenType
	Right Expr
}

func (*BinaryExpr) node() {}
func (*BinaryExpr) expr() {}

// ParenExpr represents a parenthesized expression.
type ParenExpr struct {
	Inner Expr
}

func (*ParenExpr) node() {}
func (*ParenExpr) expr() {}
