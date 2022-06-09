package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/koron/clustering-examination/internal/plotter"
	"github.com/koron/clustering-examination/internal/point"
)

func plot(in, out string) error {
	pp, err := point.LoadTSVFile(in)
	if err != nil {
		return err
	}

	min, max := point.BoundingBox(pp)
	p := plotter.New(plotter.Rectangle{
		Min: plotter.Point{X: 12, Y: 12},
		Max: plotter.Point{X: 1012, Y: 1012},
	}, plotter.RectangleF{
		Min: plotter.PointF{X: min.Lon, Y: min.Lat},
		Max: plotter.PointF{X: max.Lon, Y: max.Lat},
	})

	// TODO:
	img := image.NewNRGBA(image.Rectangle{Max: image.Point{X: 1024, Y: 1024}})
	for i := range img.Pix {
		img.Pix[i] = 0xff
	}

	fg := color.NRGBA{A: 0xff} // black
	for _, src := range pp {
		dst, in := p.Plot(plotter.PointF{X: src.Lon, Y: src.Lat})
		if !in {
			continue
		}
		img.SetNRGBA(dst.X, dst.Y, fg)
		img.SetNRGBA(dst.X, dst.Y+1, fg)
		img.SetNRGBA(dst.X+1, dst.Y, fg)
		img.SetNRGBA(dst.X+1, dst.Y+1, fg)
	}

	// write an image to PNG file.
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil

}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("require a TSV file")
	}
	flag.Arg(0)

	err := plot(flag.Arg(0), "tmp/out.png")
	if err != nil {
		log.Fatal(err)
	}
}
