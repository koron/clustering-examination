package main

import (
	"flag"
	"fmt"
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
	// bounding box
	min, max := point.BoundingBox(pp)
	diffLat := max.Lat - min.Lat
	fmt.Printf("Lat: %f, %f (%f)\n", min.Lat, max.Lat, diffLat)
	diffLon := max.Lon - min.Lon
	fmt.Printf("Lon: %f, %f (%f)\n", min.Lon, max.Lon, diffLon)
}
