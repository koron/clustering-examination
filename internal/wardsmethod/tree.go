package wardsmethod

import "gonum.org/v1/gonum/spatial/r2"

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

func (tr Tree) Weights(indexes []int) []float64 {
	weights := make([]float64, len(indexes))
	for i, x := range indexes {
		weights[i] = tr[x].Weight
	}
	return weights
}

func (tr Tree) Deltas(indexes []int) []float64 {
	deltas := make([]float64, len(indexes))
	for i, x := range indexes {
		deltas[i] = tr[x].Delta
	}
	return deltas
}
