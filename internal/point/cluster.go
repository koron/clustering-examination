package point

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
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

	log.Print("setup")

	u := NewUTMatrix(len(alives))
	for i, a := range alives {
		for j, b := range alives[i+1:] {
			u.set(i, i+j+1, distance(a, b, points, nodes))
		}
	}

	log.Print("matrix")

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

	log.Print("minmax")

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for i := 0; i < num; i++ {
		fmt.Fprintf(w, "%d,%e,%e\n", i, mins[i], maxs[i])
	}

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
	//	if i%200 == 199 {
	//		log.Printf("col#%d", i)
	//	}
	//}

	// TODO:
	return nil, 0
}
