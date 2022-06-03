package point

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

func LoadTSV(r io.Reader) ([]Point, error) {
	r2 := csv.NewReader(r)
	r2.Comma = '\t'
	r2.Comment = '#'
	//r2.FieldsPerRecord = 3
	pp := make([]Point, 0, 4096)
	for {
		ff, err := r2.Read()
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
		pp = append(pp, Point{Lat: lat, Lon: lon, Name: ff[2]})
	}
	return pp, nil
}

func LoadTSVFile(name string) ([]Point, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadTSV(f)
}
