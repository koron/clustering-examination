package point

type Point struct {
	Lat float64
	Lon float64

	Name string
}

func BoundingBox(pp []Point) (min, max Point) {
	min = pp[0]
	min.Name = "MIN"
	max = pp[0]
	max.Name = "MAX"
	for _, p := range pp[1:] {
		if p.Lat < min.Lat {
			min.Lat = p.Lat
		} else if p.Lat > max.Lat {
			max.Lat = p.Lat
		}
		if p.Lon < min.Lon {
			min.Lon = p.Lon
		} else if p.Lon > max.Lon {
			max.Lon = p.Lon
		}
	}
	return min, max
}
