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
)

const exampleUsage = `Usage of randomart:
	$ randomart -out x.png -depth 9 -width 400 -height 400
	Generate a 400x400 image at x.png with a expression tree depth of 9

	$ randomart -seed foobar -out x.png
	Generate an image with the given seed. Only the first 32 bytes of the seed string are used.

`

const advancedUsage = `Advanced options:
	-depth int
	      Expression tree depth (default 8)

	-out string
	      Image output file (default "randomart.png")

	-seed string
	      Seed for the random number generator

	-width int
	      Image width (default 600)

	-height int
	      Image height (default 600)

`

func main() {
	var (
		helpFlag = flag.Bool("help", false, "Show advanced options")
		outfile  = flag.String("out", "randomart.png", "Image output file")
		depth    = flag.Int("depth", randomart.DefaultDepth, "Expression tree depth")
		seedFlag = flag.String("seed", "", "Seed for the random number generator")
		width    = flag.Int("width", 600, "Image width")
		height   = flag.Int("height", 600, "Image height")
	)
	fs := flag.CommandLine
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), exampleUsage)
		fmt.Fprintln(fs.Output(), "For more options, run \"randomart -help\"")
	}
	flag.Parse()

	if *helpFlag {
		fmt.Print(exampleUsage)
		fmt.Print(advancedUsage)
		os.Exit(2)
	}

	if *width <= 0 {
		die(fmt.Errorf("width must be greater than 0"))
	}

	if *height <= 0 {
		die(fmt.Errorf("height must be greater than 0"))
	}

	var seed [32]byte
	if *seedFlag != "" {
		copy(seed[:], []byte(*seedFlag))
	} else {
		if _, err := cryptorand.Read(seed[:]); err != nil {
			die(fmt.Errorf("error seeding the RNG: %w", err))
		}
	}
	src := rand.NewChaCha8(seed)

	expr, err := randomart.Fuzz(src, randomart.Grammar, *depth)
	if err != nil {
		die(fmt.Errorf("error generating an expression tree: %w", err))
	}

	im, err := randomart.FromAST(src, image.Rect(0, 0, *width, *height), expr)
	if err != nil {
		die(err)
	}

	if err := writeImage(*outfile, im); err != nil {
		die(fmt.Errorf("error saving the image: %w", err))
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
