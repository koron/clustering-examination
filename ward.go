package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/koron/clustering-examination/internal/loader"
	"github.com/koron/clustering-examination/internal/wardsmethod"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func writeTree(name string, tree wardsmethod.Tree) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	defer b.Flush()
	return wardsmethod.DumpTree(b, tree)
}

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

func drawHist(name, title string, data []float64, n int) (*plotter.Histogram, error) {
	h, err := plotter.NewHist(plotter.Values(data), n)
	if err != nil {
		return nil, err
	}
	p := plot.New()
	p.Title.Text = title
	p.Add(h)
	err = p.Save(1024, 1024, name)
	if err != nil {
		return nil, err
	}
	return h, nil
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

	// draw histograms of weight, delta of nodes
	//drawHist("tmp/ward-mid-weigts.png", "Middle Weight", midTree.Weights(midAlives), 10)
	//drawHist("tmp/ward-mid-deltas.png", "Middle Delta", midTree.Deltas(midAlives), 10)
	//drawHist("tmp/ward-last-weigts.png", "Last Weigt", tree.Weights(tops), 10)
	//drawHist("tmp/ward-last-deltas.png", "Last Delta", tree.Deltas(tops), 10)
	//drawHist("tmp/ward-mean-weigts.png", "Mean Weigt", tree.Weights(means), 10)
	//drawHist("tmp/ward-mean-deltas.png", "Mean Delta", tree.Deltas(means), 10)

	//fmt.Println()
	//wardsmethod.Dump(os.Stdout, midTree, midAlives)
	//fmt.Println()
	//wardsmethod.Dump(os.Stdout, tree, tops)
	//fmt.Println()
	//wardsmethod.Dump(os.Stdout, tree, means)

	//err = writeTree("tmp/tree.txt", tree)
	//if err != nil {
	//	return err
	//}

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
