package wardsmethod

import (
	"log"
	"math"
	"time"

	"gonum.org/v1/gonum/spatial/r2"
)

type Node struct {
	Left   int
	Right  int
	Weight float64
	Center r2.Vec
	Var    float64
}

type Tree []Node

func prepare(vecs []r2.Vec) ([]Node, []int) {
	nodes := make([]Node, 0, len(vecs)*2-1)
	alives := make([]int, len(vecs))
	for i, v := range vecs {
		nodes = append(nodes, Node{
			Left:   -1,
			Right:  -1,
			Weight: 1,
			Center: v,
			Var:    0,
		})
		alives[i] = i
	}
	return nodes, alives
}

func cleanup(alives []int) []int {
	w := 0
	for _, v := range alives {
		if v != -1 {
			alives[w] = v
			w++
		}
	}
	return alives[:w]
}

func centerOfGravity(a, b Node) r2.Vec {
	w := a.Weight + b.Weight
	return r2.Vec{
		X: (a.Center.X*a.Weight + b.Center.X*b.Weight) / w,
		Y: (a.Center.Y*a.Weight + b.Center.Y*b.Weight) / w,
	}
}

func enumVec(x int, nodes []Node, fn func(r2.Vec)) {
	n := nodes[x]
	if n.Left < 0 || n.Right < 0 {
		fn(n.Center)
		return
	}
	enumVec(n.Left, nodes, fn)
	enumVec(n.Right, nodes, fn)
}

func summaryStatistic1(a, b int, nodes []Node) (r2.Vec, float64) {
	center := centerOfGravity(nodes[a], nodes[b])
	var sum float64
	fn := func(v r2.Vec) {
		x := v.X - center.X
		y := v.Y - center.Y
		sum += x*x + y*y
	}
	enumVec(a, nodes, fn)
	enumVec(b, nodes, fn)

	//na, nb := nodes[a], nodes[b]
	//x := na.Center.X - nb.Center.X
	//y := na.Center.Y - nb.Center.Y
	//tmp := na.Weight*nb.Weight/(na.Weight+nb.Weight)*(x*x+y*y) + na.Var + nb.Var

	//if math.Abs(tmp-sum) > 1e-6 {
	//	log.Printf("summaryStatistic: tmp=%e sum=%e na(%d)=%+v nb(%d)=%+v", tmp, sum, a, na, b, nb)
	//	panic("end")
	//}

	return center, sum
}

func divergence(a, b Node) float64 {
	x := a.Center.X - b.Center.X
	y := a.Center.Y - b.Center.Y
	return a.Weight*b.Weight/(a.Weight+b.Weight)*(x*x+y*y)
}

func link(nodes []Node, alives []int) ([]Node, []int) {
	var (
		num  = len(alives)
		min  = math.Inf(1)
		minA = -1
		minB = -1
		minI = -1
		minJ = -1
	)
	for i := 0; i < num; i++ {
		a := alives[i]
		for j := i + 1; j < num; j++ {
			b := alives[j]
			d := divergence(nodes[a], nodes[b])
			if d < min {
				min = d
				minA = a
				minB = b
				minI = i
				minJ = j
			}
		}
	}

	next := len(nodes)
	nodes = append(nodes, Node{
		Left:   minA,
		Right:  minB,
		Weight: nodes[minA].Weight + nodes[minB].Weight,
		Center: centerOfGravity(nodes[minA], nodes[minB]),
		Var:    min,
	})

	alives[minI] = next
	copy(alives[minJ:num-1], alives[minJ+1:num])

	return nodes, alives[:num-1]
}

func Clustering(vecs []r2.Vec) Tree {
	start := time.Now()
	nodes, alives := prepare(vecs)
	for len(alives) > 1 {
		nodes, alives = link(nodes, alives)
	}
	log.Printf("Clustering elapsed %s, len(nodes)=%d", time.Since(start), len(nodes))
	return Tree(nodes)
}
