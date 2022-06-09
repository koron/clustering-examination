package main

import (
	"flag"
	"log"

	"github.com/koron/clustering-examination/internal/point"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func run(in, out string) error {
	pp, err := point.LoadTSVFile(in)
	if err != nil {
		return err
	}

	xys := make(plotter.XYs, 0, len(pp))
	for _, p := range pp {
		xys = append(xys, plotter.XY{X: p.Lon, Y: p.Lat})
	}
	sc, err := plotter.NewScatter(xys)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Add(sc)
	err = p.Save(1000, 1000, "tmp/out.svg")
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

	err := run(flag.Arg(0), "out.png")
	if err != nil {
		log.Fatal(err)
	}
}
