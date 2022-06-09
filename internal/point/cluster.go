package point

import (
	"fmt"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Node struct {
	A int
	B int
	P int

	Dist float64
}

type Cluster struct {
	Points []Point
	Nodes  []Node
}

func enumeratePoint(start int, points []Point, nodes []Node, fn func(Point)) {
	n := nodes[start]
	if n.A >= 0 {
		enumeratePoint(n.A, points, nodes, fn)
	}
	if n.B >= 0 {
		enumeratePoint(n.B, points, nodes, fn)
	}
	if n.P >= 0 {
		fn(points[n.P])
	}
}

func distance(a, b int, points []Point, nodes []Node) float64 {
	n := 0
	var center Point
	enumeratePoint(a, points, nodes, func(p Point) {
		n++
		center.Lat += p.Lat
		center.Lon += p.Lon
	})
	enumeratePoint(b, points, nodes, func(p Point) {
		n++
		center.Lat += p.Lat
		center.Lon += p.Lon
	})
	if n == 0 {
		return 0
	}
	center.Lat /= float64(n)
	center.Lon /= float64(n)

	var sum float64
	enumeratePoint(a, points, nodes, func(p Point) {
		lat := p.Lat - center.Lat
		lon := p.Lon - center.Lon
		sum += lat*lat + lon*lon
	})
	return sum - nodes[a].Dist - nodes[b].Dist
}

func Clustering(points []Point) (c *Cluster, root int) {
	num := len(points)
	nodes := make([]Node, 0, num*2)
	alives := make([]int, num)
	for i := range points {
		alives[i] = i
		nodes = append(nodes, Node{A: -1, B: -1, P: i, Dist: 0})
	}

	u := NewUTMatrix(len(alives))
	for i, a := range alives {
		for j, b := range alives[i+1:] {
			u.set(i, i+j+1, distance(a, b, points, nodes))
		}
	}

	mins := make([]float64, num)
	maxs := make([]float64, num)
	for i := 0; i < num; i++ {
		min, max := math.Inf(0), math.Inf(-1)
		for j := 0; j < num; j++ {
			if i == j {
				continue
			}
			v := u.get(i, j)
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
		mins[i], maxs[i] = min, max
	}

	// dump min and max
	//w := bufio.NewWriter(os.Stdout)
	//defer w.Flush()
	//for i := 0; i < num; i++ {
	//	fmt.Fprintf(w, "%d,%e,%e\n", i, mins[i], maxs[i])
	//}

	//w := bufio.NewWriter(os.Stdout)
	//defer w.Flush()
	//for i := range alives {
	//	for j := range alives {
	//		if j != 0 {
	//			fmt.Fprint(w, ",")
	//		}
	//		fmt.Fprintf(w, "%f", u.get(i, j))
	//	}
	//	fmt.Fprintln(w)
	//}

	curr := mins
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("tmp/hist_min%d.png", i)
		h, err := drawHist(name, curr, 10)
		if err != nil {
			log.Printf("drawHist failed: %s", err)
			return nil, 0
		}

		b1 := h.Bins[0]
		if b1.Weight == 0 {
			break
		}
		next := make([]float64, 0, int(b1.Weight))
		for _, v := range mins {
			if v < b1.Max {
				next = append(next, v)
			}
		}
		curr = next
	}

	// TODO:
	return nil, 0
}

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