package randomart

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/waterfountain1996/randomart/ast"
)

// eval evalutes expr to it's simplest form: either an ast.Number or *ast.Triple.
func eval(rng *rand.Rand, x, y float64, expr ast.Node) (ast.Node, error) {
	switch node := expr.(type) {
	case ast.Number, ast.Bool:
		return node, nil
	case ast.Symbol:
		switch node {
		case "x":
			return ast.Number(x), nil
		case "y":
			return ast.Number(y), nil
		case "rnd":
			return ast.Number(rng.Float64()*2 - 1), nil
		default:
			return nil, fmt.Errorf("unknown symbol: '%s'", node)
		}
	case *ast.Triple:
		a, err := eval(rng, x, y, node.A)
		if err != nil {
			return nil, fmt.Errorf("triple: a: %w", err)
		}
		b, err := eval(rng, x, y, node.B)
		if err != nil {
			return nil, fmt.Errorf("triple: b: %w", err)
		}
		c, err := eval(rng, x, y, node.C)
		if err != nil {
			return nil, fmt.Errorf("triple: c: %w", err)
		}
		return &ast.Triple{A: a, B: b, C: c}, nil
	case *ast.Func1:
		return evalFunc1(rng, x, y, node)
	case *ast.BinOp:
		return evalBinOp(rng, x, y, node)
	case *ast.IfStmt:
		return evalIfStmt(rng, x, y, node)
	}
	return nil, fmt.Errorf("unknown node: %T", expr)
}

// evalFunc1 evalues a function with an arity of 1.
func evalFunc1(rng *rand.Rand, x, y float64, expr *ast.Func1) (ast.Node, error) {
	arg, err := eval(rng, x, y, expr.Arg)
	if err != nil {
		return nil, fmt.Errorf("func1: arg: %w", err)
	}
	a, ok := arg.(ast.Number)
	if !ok {
		// if trip, ok := arg.(*ast.Triple); ok {
		// 	return eval(rng, x, y, &ast.Triple{
		// 		A: &ast.Func1{Name: expr.Name, Arg: trip.A},
		// 		B: &ast.Func1{Name: expr.Name, Arg: trip.B},
		// 		C: &ast.Func1{Name: expr.Name, Arg: trip.C},
		// 	})
		// }
		return nil, fmt.Errorf("func1: expected arg to be Number, got %T", arg)
	}

	switch expr.Name {
	case "sin":
		return ast.Number(math.Sin(float64(a))), nil
	case "cos":
		return ast.Number(math.Cos(float64(a))), nil
	case "exp":
		return ast.Number(math.Exp(float64(a))), nil
	case "log":
		return ast.Number(math.Log(float64(a))), nil
	case "log2":
		return ast.Number(math.Log2(float64(a))), nil
	case "log10":
		return ast.Number(math.Log10(float64(a))), nil
	case "log1p":
		return ast.Number(math.Log1p(float64(a))), nil
	case "sqrt":
		return ast.Number(math.Sqrt(float64(a))), nil
	default:
		return nil, fmt.Errorf("func1: unknown function: %s/1", expr.Name)
	}
}

// evalBinOp evalues a binary operation (or a function with an arity of 2).
func evalBinOp(rng *rand.Rand, x, y float64, expr *ast.BinOp) (ast.Node, error) {
	lhs, err := eval(rng, x, y, expr.Lhs)
	if err != nil {
		return nil, fmt.Errorf("binop: lhs: %w", err)
	}
	// The fractions 0.299, 0.587 and 0.114 are used to convert RGB to grayscale.
	a, ok := lhs.(ast.Number)
	if !ok {
		// if trip, ok := lhs.(*ast.Triple); ok {
		// 	return eval(rng, x, y, &ast.BinOp{
		// 		Op: expr.Op,
		// 		Lhs: &ast.BinOp{
		// 			Op:  "add",
		// 			Lhs: &ast.BinOp{Op: "mul", Lhs: ast.Number(0.299), Rhs: trip.A},
		// 			Rhs: &ast.BinOp{
		// 				Op:  "add",
		// 				Lhs: &ast.BinOp{Op: "mul", Lhs: ast.Number(0.587), Rhs: trip.B},
		// 				Rhs: &ast.BinOp{Op: "mul", Lhs: ast.Number(0.114), Rhs: trip.C},
		// 			},
		// 		},
		// 		Rhs: expr.Rhs,
		// 	})
		// }
		return nil, fmt.Errorf("binop: expected lhs to be Number, got %T", lhs)
	}

	rhs, err := eval(rng, x, y, expr.Rhs)
	if err != nil {
		return nil, fmt.Errorf("binop: rhs: %w", err)
	}
	b, ok := rhs.(ast.Number)
	if !ok {
		// if trip, ok := rhs.(*ast.Triple); ok {
		// 	return eval(rng, x, y, &ast.BinOp{
		// 		Op:  expr.Op,
		// 		Lhs: expr.Lhs,
		// 		Rhs: &ast.BinOp{
		// 			Op:  "add",
		// 			Lhs: &ast.BinOp{Op: "mul", Lhs: ast.Number(0.299), Rhs: trip.A},
		// 			Rhs: &ast.BinOp{
		// 				Op:  "add",
		// 				Lhs: &ast.BinOp{Op: "mul", Lhs: ast.Number(0.587), Rhs: trip.B},
		// 				Rhs: &ast.BinOp{Op: "mul", Lhs: ast.Number(0.114), Rhs: trip.C},
		// 			},
		// 		},
		// 	})
		// }
		return nil, fmt.Errorf("binop: expected rhs to be Number, got %T", rhs)
	}

	switch expr.Op {
	case "+", "add":
		return ast.Number((a + b) / 2), nil
	case "-", "sub":
		return ast.Number((a - b) / 2), nil
	case "*", "mul":
		return ast.Number(a * b), nil
	case "/", "div":
		return ast.Number(a / b), nil
	case "%", "mod":
		return ast.Number(math.Mod(float64(a), float64(b))), nil
	case ">", "gt":
		return ast.Bool(a > b), nil
	case ">=", "gte":
		return ast.Bool(a >= b), nil
	case "<", "lt":
		return ast.Bool(a < b), nil
	case "<=", "lte":
		return ast.Bool(a <= b), nil
	}

	return nil, fmt.Errorf("binop: unknown operator: %s", expr.Op)
}

// evalIfStmt evalutes an if statement. stmt.Cond must evaluate to ast.Bool.
func evalIfStmt(rng *rand.Rand, x, y float64, stmt *ast.IfStmt) (ast.Node, error) {
	cond, err := eval(rng, x, y, stmt.Cond)
	if err != nil {
		return nil, fmt.Errorf("if: cond: %w", err)
	}

	truth, ok := cond.(ast.Bool)
	if !ok {
		return nil, fmt.Errorf("if: expected cond to be Bool, got %T", cond)
	}

	if truth {
		res, err := eval(rng, x, y, stmt.Then)
		if err != nil {
			return nil, fmt.Errorf("if: then: %w", err)
		}
		return res, nil
	}

	res, err := eval(rng, x, y, stmt.Else)
	if err != nil {
		return nil, fmt.Errorf("if: else: %w", err)
	}
	return res, nil
}
