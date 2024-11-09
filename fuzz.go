package randomart

import (
	"fmt"
	"math/rand/v2"

	"github.com/waterfountain1996/randomart/ast"
)

// Default expression tree depth.
const DefaultDepth = 8

// Fuzz takes a grammar and produces a randomly generated expression tree of given depth.
func Fuzz(src rand.Source, entry *ast.Rule, depth int) (ast.Node, error) {
	return evalGrammar(rand.New(src), entry, depth+1)
}

func evalGrammar(rng *rand.Rand, expr ast.Node, depth int) (ast.Node, error) {
	switch node := expr.(type) {
	case ast.Number, ast.Bool, ast.Symbol:
		return node, nil
	case *ast.Rule:
		if depth <= 0 {
			return evalGrammar(rng, node.Branches[0].Node, depth-1)
		}
		idx := rng.IntN(len(node.Branches))
		br := node.Branches[idx]
		return evalGrammar(rng, br.Node, depth-1)
	case *ast.Triple:
		a, err := evalGrammar(rng, node.A, depth)
		if err != nil {
			return nil, fmt.Errorf("triple: a: %w", err)
		}
		b, err := evalGrammar(rng, node.B, depth)
		if err != nil {
			return nil, fmt.Errorf("triple: b: %w", err)
		}
		c, err := evalGrammar(rng, node.C, depth)
		if err != nil {
			return nil, fmt.Errorf("triple: c: %w", err)
		}
		return &ast.Triple{A: a, B: b, C: c}, nil
	case *ast.Func1:
		return evalGrammarFunc1(rng, node, depth)
	case *ast.BinOp:
		return evalGrammarBinOp(rng, node, depth)
	case *ast.IfStmt:
		return evalGrammarIfStmt(rng, node, depth)
	}
	return nil, fmt.Errorf("unknown node: %T", expr)
}

func evalGrammarFunc1(rng *rand.Rand, expr *ast.Func1, depth int) (ast.Node, error) {
	arg, err := evalGrammar(rng, expr.Arg, depth)
	if err != nil {
		return nil, fmt.Errorf("func1: arg: %w", err)
	}

	return &ast.Func1{
		Name: expr.Name,
		Arg:  arg,
	}, nil
}

func evalGrammarBinOp(rng *rand.Rand, expr *ast.BinOp, depth int) (ast.Node, error) {
	lhs, err := evalGrammar(rng, expr.Lhs, depth)
	if err != nil {
		return nil, fmt.Errorf("binop: lhs: %w", err)
	}

	rhs, err := evalGrammar(rng, expr.Rhs, depth)
	if err != nil {
		return nil, fmt.Errorf("binop: rhs: %w", err)
	}

	return &ast.BinOp{
		Op:  expr.Op,
		Lhs: lhs,
		Rhs: rhs,
	}, nil
}

func evalGrammarIfStmt(rng *rand.Rand, stmt *ast.IfStmt, depth int) (ast.Node, error) {
	cond, err := evalGrammar(rng, stmt.Cond, depth)
	if err != nil {
		return nil, fmt.Errorf("if: cond: %w", err)
	}

	thenStmt, err := evalGrammar(rng, stmt.Then, depth)
	if err != nil {
		return nil, fmt.Errorf("if: then: %w", err)
	}

	elseStmt, err := evalGrammar(rng, stmt.Else, depth)
	if err != nil {
		return nil, fmt.Errorf("if: else: %w", err)
	}

	return &ast.IfStmt{
		Cond: cond,
		Then: thenStmt,
		Else: elseStmt,
	}, nil
}
