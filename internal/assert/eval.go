// Package assert provides assertion expression evaluation.
package assert

import (
	"fmt"
	"strconv"
	"strings"
)

// Value represents a runtime value during evaluation.
type Value struct {
	raw string // original string representation
}

// NewValue creates a new Value from a string.
func NewValue(s string) Value {
	return Value{raw: s}
}

// String returns the string representation of the value.
func (v Value) String() string {
	return v.raw
}

// AsBool interprets the value as a boolean.
// "true" -> true, "false" -> false
// Non-empty string -> true, empty string -> false
// Numbers: 0 -> false, non-zero -> true
func (v Value) AsBool() bool {
	if v.raw == "true" {
		return true
	}
	if v.raw == "false" || v.raw == "" {
		return false
	}
	// Try as number: 0 is false, non-zero is true
	if f, err := strconv.ParseFloat(v.raw, 64); err == nil {
		return f != 0
	}
	// Non-empty string is truthy
	return true
}

// AsFloat attempts to parse the value as a float64.
func (v Value) AsFloat() (float64, bool) {
	f, err := strconv.ParseFloat(v.raw, 64)
	return f, err == nil
}

// IsNumeric returns true if the value can be parsed as a number.
func (v Value) IsNumeric() bool {
	_, ok := v.AsFloat()
	return ok
}

// Evaluator evaluates assertion expressions.
type Evaluator struct{}

// New creates a new Evaluator.
func New() *Evaluator {
	return &Evaluator{}
}

// Eval evaluates an assertion expression with the given variables.
// Returns (true, nil) if the assertion passes, (false, nil) if it fails,
// or (false, error) if there's an evaluation error.
// If expr is empty, returns true (no assertion = always pass).
func (e *Evaluator) Eval(expr string, vars map[string]string) (bool, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return true, nil
	}

	parser := NewParser(expr)
	ast, err := parser.Parse()
	if err != nil {
		return false, fmt.Errorf("parse error in assertion %q: %w", expr, err)
	}

	result, err := e.eval(ast, vars)
	if err != nil {
		return false, fmt.Errorf("eval error in assertion %q: %w", expr, err)
	}

	return result.AsBool(), nil
}

// eval recursively evaluates an AST node.
func (e *Evaluator) eval(node Expr, vars map[string]string) (Value, error) {
	switch n := node.(type) {
	case *NumberLit:
		return NewValue(n.Value), nil

	case *StringLit:
		return NewValue(n.Value), nil

	case *BoolLit:
		if n.Value {
			return NewValue("true"), nil
		}
		return NewValue("false"), nil

	case *Ident:
		val, ok := vars[n.Name]
		if !ok {
			// Return empty string for undefined variables
			return NewValue(""), nil
		}
		return NewValue(val), nil

	case *ParenExpr:
		return e.eval(n.Inner, vars)

	case *UnaryExpr:
		return e.evalUnary(n, vars)

	case *BinaryExpr:
		return e.evalBinary(n, vars)

	default:
		return Value{}, fmt.Errorf("unknown node type: %T", node)
	}
}

// evalUnary evaluates a unary expression.
func (e *Evaluator) evalUnary(expr *UnaryExpr, vars map[string]string) (Value, error) {
	right, err := e.eval(expr.Right, vars)
	if err != nil {
		return Value{}, err
	}

	switch expr.Op {
	case TokenNot:
		if right.AsBool() {
			return NewValue("false"), nil
		}
		return NewValue("true"), nil

	case TokenMinus:
		f, ok := right.AsFloat()
		if !ok {
			return Value{}, fmt.Errorf("cannot negate non-numeric value %q (operator: -%s)", right.raw, right.raw)
		}
		return NewValue(formatFloat(-f)), nil

	default:
		return Value{}, fmt.Errorf("unknown unary operator: %v", expr.Op)
	}
}

// evalBinary evaluates a binary expression.
func (e *Evaluator) evalBinary(expr *BinaryExpr, vars map[string]string) (Value, error) {
	left, err := e.eval(expr.Left, vars)
	if err != nil {
		return Value{}, err
	}

	// Short-circuit evaluation for logical operators
	switch expr.Op {
	case TokenAnd:
		if !left.AsBool() {
			return NewValue("false"), nil
		}
		right, err := e.eval(expr.Right, vars)
		if err != nil {
			return Value{}, err
		}
		if right.AsBool() {
			return NewValue("true"), nil
		}
		return NewValue("false"), nil

	case TokenOr:
		if left.AsBool() {
			return NewValue("true"), nil
		}
		right, err := e.eval(expr.Right, vars)
		if err != nil {
			return Value{}, err
		}
		if right.AsBool() {
			return NewValue("true"), nil
		}
		return NewValue("false"), nil
	}

	// Evaluate right side for all other operators
	right, err := e.eval(expr.Right, vars)
	if err != nil {
		return Value{}, err
	}

	switch expr.Op {
	// Arithmetic operators
	case TokenPlus:
		return e.evalArithmetic(left, right, func(a, b float64) float64 { return a + b })
	case TokenMinus:
		return e.evalArithmetic(left, right, func(a, b float64) float64 { return a - b })
	case TokenAsterisk:
		return e.evalArithmetic(left, right, func(a, b float64) float64 { return a * b })
	case TokenSlash:
		return e.evalArithmetic(left, right, func(a, b float64) float64 {
			if b == 0 {
				return 0 // Division by zero returns 0
			}
			return a / b
		})

	// Comparison operators
	case TokenEq:
		return e.evalComparison(left, right, func(cmp int) bool { return cmp == 0 })
	case TokenNotEq:
		return e.evalComparison(left, right, func(cmp int) bool { return cmp != 0 })
	case TokenLT:
		return e.evalComparison(left, right, func(cmp int) bool { return cmp < 0 })
	case TokenLTE:
		return e.evalComparison(left, right, func(cmp int) bool { return cmp <= 0 })
	case TokenGT:
		return e.evalComparison(left, right, func(cmp int) bool { return cmp > 0 })
	case TokenGTE:
		return e.evalComparison(left, right, func(cmp int) bool { return cmp >= 0 })

	default:
		return Value{}, fmt.Errorf("unknown binary operator: %v", expr.Op)
	}
}

// evalArithmetic evaluates an arithmetic expression.
func (e *Evaluator) evalArithmetic(left, right Value, op func(a, b float64) float64) (Value, error) {
	lf, lok := left.AsFloat()
	rf, rok := right.AsFloat()

	if !lok || !rok {
		return Value{}, fmt.Errorf("arithmetic requires numeric operands, got left=%q right=%q", left.raw, right.raw)
	}

	result := op(lf, rf)
	return NewValue(formatFloat(result)), nil
}

// evalComparison evaluates a comparison expression.
// Returns true/false value based on the comparison.
// Auto-coerces numeric strings when both sides are numeric.
func (e *Evaluator) evalComparison(left, right Value, compare func(cmp int) bool) (Value, error) {
	// Try numeric comparison first (auto-coerce)
	lf, lok := left.AsFloat()
	rf, rok := right.AsFloat()

	var cmp int
	if lok && rok {
		// Both are numeric - compare as numbers
		if lf < rf {
			cmp = -1
		} else if lf > rf {
			cmp = 1
		} else {
			cmp = 0
		}
	} else {
		// String comparison
		cmp = strings.Compare(left.raw, right.raw)
	}

	if compare(cmp) {
		return NewValue("true"), nil
	}
	return NewValue("false"), nil
}

// formatFloat formats a float64 as a string, removing trailing zeros.
func formatFloat(f float64) string {
	// Check if it's a whole number
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
