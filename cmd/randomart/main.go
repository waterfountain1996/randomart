package main

import (
	cryptorand "crypto/rand"
	"flag"
	"fmt"
	"image"
	"image/png"
	"math/big"
	"math/rand/v2"
	"os"

	"github.com/waterfountain1996/randomart"
	"github.com/waterfountain1996/randomart/ast"
)

func main() {
	outfile := flag.String("out", "randomart.png", "Image output file")
	flag.Parse()

	var seed [32]byte
	if _, err := cryptorand.Read(seed[:]); err != nil {
		die(fmt.Errorf("error seeding the RNG: %w", err))
	}
	src := rand.NewChaCha8(seed)

	expr, err := randomart.Fuzz(src, exampleGrammar(), randomart.DefaultDepth)
	if err != nil {
		die(fmt.Errorf("error generating an expression tree: %w", err))
	}

	im, err := randomart.FromAST(src, image.Rect(0, 0, 600, 600), expr)
	if err != nil {
		die(err)
	}

	if err := writeImage(*outfile, im); err != nil {
		die(fmt.Errorf("error saving the image: %w", err))
	}
}

func exampleGrammar() *ast.Rule {
	a := &ast.Rule{
		Name: "A",
		Branches: []*ast.Branch{
			{
				Node:        ast.Symbol("x"),
				Probability: big.NewRat(1, 3),
			},
			{
				Node:        ast.Symbol("y"),
				Probability: big.NewRat(1, 3),
			},
			{
				Node:        ast.Symbol("rnd"),
				Probability: big.NewRat(1, 3),
			},
		},
	}

	c := &ast.Rule{
		Name:     "C",
		Branches: make([]*ast.Branch, 0, 3),
	}
	c.Branches = append(c.Branches, &ast.Branch{
		Node:        a,
		Probability: big.NewRat(1, 4),
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.BinOp{
			Op:  "add",
			Lhs: c,
			Rhs: c,
		},
		Probability: big.NewRat(3, 8),
	})
	c.Branches = append(c.Branches, &ast.Branch{
		Node: &ast.BinOp{
			Op:  "mul",
			Lhs: c,
			Rhs: c,
		},
		Probability: big.NewRat(3, 8),
	})

	entry := &ast.Rule{
		Name: "E",
		Branches: []*ast.Branch{
			{
				Node:        &ast.Triple{A: c, B: c, C: c},
				Probability: big.NewRat(1, 1),
			},
		},
	}
	return entry
}

func grayscale() ast.Node {
	return ast.Symbol("x")
}

func cool() ast.Node {
	return &ast.IfStmt{
		Cond: &ast.BinOp{
			Op: "gte",
			Lhs: &ast.BinOp{
				Op:  "mul",
				Lhs: ast.Symbol("x"),
				Rhs: ast.Symbol("y"),
			},
			Rhs: ast.Number(0.0),
		},
		Then: &ast.Triple{
			A: ast.Symbol("x"),
			B: ast.Symbol("y"),
			C: ast.Number(1.0),
		},
		Else: &ast.Triple{
			A: &ast.BinOp{Op: "mod", Lhs: ast.Symbol("x"), Rhs: ast.Symbol("y")},
			B: &ast.BinOp{Op: "mod", Lhs: ast.Symbol("x"), Rhs: ast.Symbol("y")},
			C: &ast.BinOp{Op: "mod", Lhs: ast.Symbol("x"), Rhs: ast.Symbol("y")},
		},
	}
}

func writeImage(filename string, im image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	if err := png.Encode(f, im); err != nil {
		return err
	}

	return f.Sync()
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "randomart: %s\n", err)
	os.Exit(1)
}
