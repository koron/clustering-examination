package loader

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"

	"gonum.org/v1/gonum/spatial/r2"
)

func LoadTSV(r io.Reader) ([]r2.Vec, error) {
	rr := csv.NewReader(r)
	rr.Comma = '\t'
	rr.Comment = '#'
	//rr.FieldsPerRecord = 3
	vecs := make([]r2.Vec, 0, 4096)
	for {
		ff, err := rr.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		lat, err := strconv.ParseFloat(ff[0], 64)
		if err != nil {
			return nil, err
		}
		lon, err := strconv.ParseFloat(ff[1], 64)
		if err != nil {
			return nil, err
		}
		vecs = append(vecs, r2.Vec{X: lon, Y: lat})
	}
	return vecs, nil
}

func LoadTSVFile(name string) ([]r2.Vec, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadTSV(f)
}
