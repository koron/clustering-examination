package point

import (
	"math"
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

func (c *Cluster) forPoints(start int, fn func(Point)) {
	n := c.Nodes[start]
	if n.A >= 0 {
		c.forPoints(n.A, fn)
	}
	if n.B >= 0 {
		c.forPoints(n.B, fn)
	}
	if n.P >= 0 {
		fn(c.Points[n.P])
	}
}

func (c *Cluster) wardDistance(a, b int) float64 {
	n := 0
	var center Point
	c.forPoints(a, func(p Point) {
		n++
		center.Lat += p.Lat
		center.Lon += p.Lon
	})
	c.forPoints(b, func(p Point) {
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
	c.forPoints(a, func(p Point) {
		lat := p.Lat - center.Lat
		lon := p.Lon - center.Lon
		sum += lat*lat + lon*lon
	})
	c.forPoints(b, func(p Point) {
		lat := p.Lat - center.Lat
		lon := p.Lon - center.Lon
		sum += lat*lat + lon*lon
	})
	return sum - c.Nodes[a].Dist - c.Nodes[b].Dist
}

func Clustering(points []Point) (c *Cluster, root int) {
	num := len(points)
	nodes := make([]Node, 0, num*2)
	alives := make([]int, num)
	for i := range points {
		alives[i] = i
		nodes = append(nodes, Node{A: -1, B: -1, P: i, Dist: 0})
	}

	c = &Cluster{Points: points, Nodes: nodes}

	u := NewUTMatrix(len(alives))
	for i := 0; i < len(alives); i++ {
		for j := i + 1; j < len(alives); j++ {
			u.Set(i, j, c.wardDistance(alives[i], alives[j]))
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
			v := u.Get(i, j)
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

	drawHistograms("tmp/hist_min%d.png", mins)

	//drawBubbles("tmp/bubbles.svg", points, alives, mins)

	// TODO:
	return nil, 0
}
