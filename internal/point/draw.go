package point

import (
	"fmt"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func drawHist(name string, data []float64, n int) (*plotter.Histogram, error) {
	h, err := plotter.NewHist(plotter.Values(data), n)
	if err != nil {
		return nil, err
	}
	p := plot.New()
	p.Add(h)
	err = p.Save(1024, 1024, name)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func drawHistograms(nameFormat string, curr []float64) {
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf(nameFormat, i)
		h, err := drawHist(name, curr, 10)
		if err != nil {
			log.Printf("drawHist failed: %s", err)
			return
		}
		b1 := h.Bins[0]
		if b1.Weight == 0 {
			break
		}
		next := make([]float64, 0, int(b1.Weight))
		for _, v := range curr {
			if v < b1.Max {
				next = append(next, v)
			}
		}
		curr = next
	}
}

func drawBubbles(name string, points []Point, alives []int, mins []float64) {
	data := make(plotter.XYZs, 0, len(points))
	for _, n := range alives {
		data = append(data, plotter.XYZ{
			X: points[n].Lon,
			Y: points[n].Lat,
			Z: math.Sqrt(mins[n]),
		})
	}
	sc, err := plotter.NewScatter(data)
	if err != nil {
		log.Printf("NewScatter failed: %s", err)
		return
	}
	//c := color.RGBA{}
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		if i%100 == 0 {
			fmt.Printf("#%d z=%e\n", i, data[i].Z)
		}
		return draw.GlyphStyle{
			Radius: vg.Length(data[i].Z),
			//Radius: vg.Length(0.1),
			Shape: draw.CircleGlyph{},
		}
	}

	p := plot.New()
	p.Add(sc)
	err = p.Save(1000, 1000, name)
	if err != nil {
		log.Printf("Plot.Save failed: %s", err)
		return
	}
}
