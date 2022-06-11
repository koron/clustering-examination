package point

import (
	"fmt"
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

func (c *Cluster) center(nodeNums ...int) Point {
	var center Point
	var weight float64
	add := func(p Point) {
		weight++
		center.Lat += p.Lat
		center.Lon += p.Lon
	}
	for _, n := range nodeNums {
		c.forPoints(n, add)
	}
	if weight == 0 {
		return Point{}
	}
	center.Lat /= weight
	center.Lon /= weight
	return center
}

func (c *Cluster) sumOfSquaresOfDistances(center Point, nodeNums ...int) float64 {
	var sum float64
	sumDist := func(p Point) {
		lat := p.Lat - center.Lat
		lon := p.Lon - center.Lon
		sum += lat*lat + lon*lon
	}
	for _, n := range nodeNums {
		c.forPoints(n, sumDist)
	}
	return sum
}

func (c *Cluster) wardDistance(a, b int) float64 {
	// calculate the center of gravity
	center := c.center(a, b)

	// calculate sum of squares of distances
	sum := c.sumOfSquaresOfDistances(center, a, b)

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

	smallOrBig(2e-6, mins)

	drawHistograms("tmp/hist_min%d.png", mins)

	//drawBubbles("tmp/bubbles.svg", points, alives, mins)

	// TODO:
	return nil, 0
}

func smallOrBig(threshold float64, data []float64) {
	smaller, bigger := 0, 0
	for _, v := range data {
		if v < threshold {
			smaller++
		} else {
			bigger++
		}
	}
	fmt.Printf("threshold=%e smaller=%d bigger=%d smaller/total=%f\n", threshold, smaller, bigger, float64(smaller)/float64(smaller+bigger))
}
