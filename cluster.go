package main

import (
	"flag"
	"log"

	"github.com/koron/exam-clustering/internal/point"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("require a TSV file")
	}
	flag.Arg(0)

	pp, err := point.LoadTSVFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	point.Clustering(pp)
}
