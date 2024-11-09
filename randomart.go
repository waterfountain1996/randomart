package randomart

import (
	"github.com/waterfountain1996/randomart/ast"
)

// Default grammar used for generating images.
var Grammar = defaultGrammar()

func defaultGrammar() *ast.Rule {
	a := &ast.Rule{
		Name: "A",
		Branches: []*ast.Branch{
			{
				Node: ast.Symbol("x"),
			},
			{
				Node: ast.Symbol("y"),
			},
			{
				Node: ast.Symbol("rnd"),
			},
		},
	}

	c := &ast.Rule{
		Name:     "C",
		Branches: make([]*ast.Branch, 0, 16),
	}
	c.Branches = append(c.Branches, &ast.Branch{
		Node: a,
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.IfStmt{
			Cond: &ast.BinOp{
				Op:  "gt",
				Lhs: ast.Symbol("x"),
				Rhs: ast.Symbol("rnd"),
			},
			Then: c,
			Else: c,
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.IfStmt{
			Cond: &ast.BinOp{
				Op:  "gt",
				Lhs: ast.Symbol("y"),
				Rhs: c,
			},
			Then: c,
			Else: c,
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.BinOp{
			Op:  "add",
			Lhs: c,
			Rhs: c,
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.BinOp{
			Op:  "mul",
			Lhs: c,
			Rhs: c,
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.BinOp{
			Op:  "mod",
			Lhs: c,
			Rhs: c,
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.BinOp{
			Op:  "div",
			Lhs: c,
			Rhs: ast.Number(-0.112),
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.Func1{
			Name: "sin",
			Arg: &ast.BinOp{
				Op:  "mul",
				Lhs: c,
				Rhs: a,
			},
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.Func1{
			Name: "exp",
			Arg:  c,
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.Func1{
			Name: "sqrt",
			Arg:  c,
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.Func1{
			Name: "log",
			Arg:  ast.Number(0.7741),
		},
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.IfStmt{
			Cond: &ast.BinOp{
				Op: "gte",
				Lhs: &ast.BinOp{
					Op:  "mul",
					Lhs: a,
					Rhs: ast.Number(-0.048),
				},
				Rhs: ast.Number(0.1134),
			},
			Then: c,
			Else: c,
		},
	})

	entry := &ast.Rule{
		Name: "E",
		Branches: []*ast.Branch{
			{
				Node: &ast.Triple{A: c, B: c, C: c},
			},
		},
	}
	return entry
}
