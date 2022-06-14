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

func wards(name string) error {
	pp, err := loader.LoadTSVFile(name)
	if err != nil {
		return err
	}

	var midNodes []wardsmethod.Node
	var midAlives []int

	start := time.Now()
	tree := wardsmethod.Clustering(pp, wardsmethod.MonitorFunc(func(nodes []wardsmethod.Node, alives []int) {
		if len(alives) == 45 {
			midNodes = nodes
			midAlives = make([]int, len(alives))
			copy(midAlives, alives)
		}
	}))
	log.Printf("Clustering elapsed %s, len(nodes)=%d", time.Since(start), len(tree))

	statistics("mid", midNodes, midAlives)
	tops := wardsmethod.Top(tree, 45)
	statistics("last", tree, tops)

	//fmt.Println()
	//wardsmethod.Dump(os.Stdout, midNodes, midAlives)
	//fmt.Println()
	//wardsmethod.Dump(os.Stdout, tree, tops)

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
