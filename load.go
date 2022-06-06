package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/koron/clustering-examination/internal/point"
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
	for i, p := range pp[:5] {
		fmt.Printf("#%d: %+v\n", i, p)
	}
}
