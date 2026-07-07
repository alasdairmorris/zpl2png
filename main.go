// Command zpl2png reads ZPL (Zebra Programming Language) from stdin and writes
// the rendered label as a PNG to stdout.
//
// A ZPL stream may contain more than one label (each delimited by ^XA ... ^XZ).
// By default the first label is rendered; use -index to select another, or
// -index=-1 to list how many labels were found without rendering.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ingridhq/zebrash"
	"github.com/ingridhq/zebrash/drawers"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "zpl2png:", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		index     = flag.Int("index", 0, "index of the label to render (0-based); -1 lists label count without rendering")
		dpmm      = flag.Int("dpmm", 8, "print density in dots per millimetre (6=152dpi, 8=203dpi, 12=300dpi, 24=600dpi)")
		widthMm   = flag.Float64("width", 0, "label width in millimetres (0 = library default, 101.6mm)")
		heightMm  = flag.Float64("height", 0, "label height in millimetres (0 = library default, 203.2mm)")
		grayscale = flag.Bool("grayscale", false, "output 8-bit grayscale PNG preserving anti-aliasing instead of binary monochrome")
		inverted  = flag.Bool("inverted", false, "render inverted-orientation labels upside-down")
		outPath   = flag.String("o", "-", "output file path (\"-\" for stdout)")
	)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: zpl2png [flags] < input.zpl > output.png\n\nReads ZPL from stdin and writes a PNG to stdout.\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	zplData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}
	if len(bytes.TrimSpace(zplData)) == 0 {
		return fmt.Errorf("no ZPL received on stdin")
	}

	labels, err := zebrash.NewParser().Parse(zplData)
	if err != nil {
		return fmt.Errorf("parsing ZPL: %w", err)
	}
	if len(labels) == 0 {
		return fmt.Errorf("no labels found in input (expected at least one ^XA ... ^XZ block)")
	}

	if *index == -1 {
		fmt.Fprintf(os.Stderr, "found %d label(s)\n", len(labels))
		return nil
	}
	if *index < 0 || *index >= len(labels) {
		return fmt.Errorf("label index %d out of range: input contains %d label(s)", *index, len(labels))
	}

	options := drawers.DrawerOptions{
		LabelWidthMm:         *widthMm,
		LabelHeightMm:        *heightMm,
		Dpmm:                 *dpmm,
		GrayscaleOutput:      *grayscale,
		EnableInvertedLabels: *inverted,
	}.WithDefaults()

	var buf bytes.Buffer
	if err := zebrash.NewDrawer().DrawLabelAsPng(labels[*index], &buf, options); err != nil {
		return fmt.Errorf("rendering label %d to PNG: %w", *index, err)
	}

	out := os.Stdout
	if *outPath != "-" {
		f, err := os.Create(*outPath)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		out = f
	}
	if _, err := out.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("writing PNG output: %w", err)
	}
	return nil
}
