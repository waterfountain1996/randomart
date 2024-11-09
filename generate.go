package randomart

import (
	"fmt"
	"image"
	"image/color"
	"math/rand/v2"

	"github.com/waterfountain1996/randomart/ast"
)

// FromAST generates an Image from an AST expression.
func FromAST(src rand.Source, bounds image.Rectangle, expr ast.Node) (image.Image, error) {
	var (
		rng = rand.New(src)
		im  = image.NewNRGBA(bounds)
	)
	for y := range bounds.Max.Y {
		ny := normalize(y, bounds.Max.Y)
		for x := range bounds.Max.X {
			nx := normalize(x, bounds.Max.X)

			node, err := eval(rng, nx, ny, expr)
			if err != nil {
				return nil, fmt.Errorf("eval at %d,%d: %w", x, y, err)
			}

			var c color.Color
			switch tt := node.(type) {
			case ast.Number:
				c = color.Gray{Y: upscale(float64(tt))}
			case *ast.Triple:
				r, ok := tt.A.(ast.Number)
				if !ok {
					return nil, fmt.Errorf("eval: invalid color triple: %T, %T, %T,", tt.A, tt.B, tt.C)
				}
				g, ok := tt.B.(ast.Number)
				if !ok {
					return nil, fmt.Errorf("eval: invalid color triple: %T, %T, %T,", tt.A, tt.B, tt.C)
				}
				b, ok := tt.C.(ast.Number)
				if !ok {
					return nil, fmt.Errorf("eval: invalid color triple: %T, %T, %T,", tt.A, tt.B, tt.C)
				}
				c = color.RGBA{
					R: upscale(float64(r)),
					G: upscale(float64(g)),
					B: upscale(float64(b)),
					A: 255,
				}
			default:
				return nil, fmt.Errorf("eval: invalid expression evaluated to %T", tt)
			}

			im.Set(x, y, c)
		}
	}

	return im, nil
}

// Scales v from [-1..1] to [0.255].
func upscale(v float64) uint8 {
	return uint8((v + 1) / 2 * 255)
}

// Scales v from [0..bounds] to [-1..1].
func normalize(v int, bounds int) float64 {
	return float64(v)/float64(bounds)*2 - 1
}
