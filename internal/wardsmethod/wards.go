package wardsmethod

import (
	"math"

	"gonum.org/v1/gonum/spatial/r2"
)

type Node struct {
	Left   int
	Right  int
	Weight float64
	Center r2.Vec
	Delta  float64
}

type Tree []Node

func (tr Tree) Root() int {
	return len(tr) - 1
}

func prepare(vecs []r2.Vec) ([]Node, []int) {
	nodes := make([]Node, 0, len(vecs)*2-1)
	alives := make([]int, len(vecs))
	for i, v := range vecs {
		nodes = append(nodes, Node{
			Left:   -1,
			Right:  -1,
			Weight: 1,
			Center: v,
			Delta:  0,
		})
		alives[i] = i
	}
	return nodes, alives
}

func centerOfGravity(a, b Node) r2.Vec {
	w := a.Weight + b.Weight
	return r2.Vec{
		X: (a.Center.X*a.Weight + b.Center.X*b.Weight) / w,
		Y: (a.Center.Y*a.Weight + b.Center.Y*b.Weight) / w,
	}
}

func delta(a, b Node) float64 {
	x := a.Center.X - b.Center.X
	y := a.Center.Y - b.Center.Y
	return a.Weight * b.Weight / (a.Weight + b.Weight) * (x*x + y*y)
}

func link(nodes []Node, alives []int) ([]Node, []int) {
	var (
		num  = len(alives)
		minD = math.Inf(1)
		minA = -1
		minB = -1
		minI = -1
		minJ = -1
	)
	for i := 0; i < num; i++ {
		a := alives[i]
		for j := i + 1; j < num; j++ {
			b := alives[j]
			d := delta(nodes[a], nodes[b])
			if d < minD {
				minD = d
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
		Delta:  minD,
	})

	alives[minI] = next
	copy(alives[minJ:num-1], alives[minJ+1:num])

	return nodes, alives[:num-1]
}

func Clustering(vecs []r2.Vec, m Monitor) Tree {
	nodes, alives := prepare(vecs)
	for len(alives) > 1 {
		if m != nil {
			m.Monitor(nodes, alives)
		}
		nodes, alives = link(nodes, alives)
	}
	return Tree(nodes)
}
