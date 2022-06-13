package main

import (
	"bufio"
	"flag"
	"log"
	"os"

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

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("require a TSV file")
	}
	flag.Arg(0)

	pp, err := loader.LoadTSVFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	tree := wardsmethod.Clustering(pp)
	err = writeTree("tmp/tree.txt", tree)
	if err != nil {
		log.Fatal("failed to write tree:", err)
	}
}
