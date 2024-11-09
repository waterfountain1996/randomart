package main

import (
	cryptorand "crypto/rand"
	"flag"
	"fmt"
	"image"
	"image/png"
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

	expr := grayscale()

	im, err := randomart.FromAST(src, image.Rect(0, 0, 600, 600), expr)
	if err != nil {
		die(err)
	}

	if err := writeImage(*outfile, im); err != nil {
		die(fmt.Errorf("error saving the image: %w", err))
	}
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
