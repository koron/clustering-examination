package main

import (
	"flag"
	"log"

	"github.com/koron/clustering-examination/internal/loader"
	"github.com/koron/clustering-examination/internal/wardsmethod"
)

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

	wardsmethod.Clustering(pp)
}
