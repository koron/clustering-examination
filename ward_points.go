package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"sort"
	"time"

	"github.com/koron/clustering-examination/internal/loader"
	"github.com/koron/clustering-examination/internal/wardsmethod"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

func statistics(label string, nodes []wardsmethod.Node, alives []int) {
	num := float64(len(alives))
	var sumW, sumD float64
	for _, v := range alives {
		n := nodes[v]
		sumW += n.Weight
		sumD += n.Delta
	}
	meanW := sumW / num
	meanD := sumD / num
	var varW, varD float64
	for _, v := range alives {
		n := nodes[v]
		varW += math.Pow(n.Weight-meanW, 2)
		varD += math.Pow(n.Delta-meanD, 2)
	}
	varW /= num
	varD /= num
	fmt.Printf("%6s: weight=%e±%e delta=%e±%e sumW=%e\n", label, meanW, varW, meanD, varD, sumW)
}

func wards(name string) error {
	pp, err := loader.LoadTSVFile(name)
	if err != nil {
		return err
	}

	var (
		midTree   wardsmethod.Tree
		midAlives []int
	)

	start := time.Now()
	tree := wardsmethod.Clustering(pp, wardsmethod.MonitorFunc(func(nodes []wardsmethod.Node, alives []int) {
		if len(alives) == 45 {
			midTree = wardsmethod.Tree(nodes)
			midAlives = make([]int, len(alives))
			copy(midAlives, alives)
		}
	}))
	log.Printf("Clustering elapsed %s, len(nodes)=%d", time.Since(start), len(tree))

	statistics("mid", midTree, midAlives)
	tops := wardsmethod.Top(tree, 45)
	statistics("last", tree, tops)
	means := wardsmethod.Mean(tree, 45)
	statistics("mean", tree, means)

	p := plot.New()
	drawPointsAscending(p, "mid", midTree, midAlives, color.Gray{0x66})
	//drawPointsAscending(p, "tops", tree, tops)
	drawPointsAscending(p, "means", tree, means, plotutil.DarkColors[0])
	p.Save(1000, 1000, "tmp/compare.png")

	return nil
}

func drawPointsAscending(p *plot.Plot, name string, tree wardsmethod.Tree, indexes []int, c color.Color) error {
	xys := make(plotter.XYs, len(indexes))
	for i, idx := range indexes {
		xys[i] = plotter.XY{Y: tree[idx].Weight}
	}
	sort.Slice(xys, func(i, j int) bool {
		return xys[i].Y < xys[j].Y
	})
	for i := range xys {
		xys[i].X = float64(i)
	}

	l, err := plotter.NewLine(xys)
	if err != nil {
		return err
	}
	l.Color = c
	p.Add(l)
	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("require a TSV file")
	}
	err := wards(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
}
