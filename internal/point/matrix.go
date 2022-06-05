package point

type UTMatrix struct {
	n int
	d []float64
}

func NewUTMatrix(n int) UTMatrix {
	size := n * (n + 1) / 2
	return UTMatrix{
		n: n,
		d: make([]float64, size),
	}
}

func (m UTMatrix) index(a, b int) int {
	if a > b {
		a, b = b, a
	}
	n := m.n - a
	return len(m.d) - (n * (n + 1) / 2) + (b - a)
}

func (m UTMatrix) get(a, b int) float64 {
	return m.d[m.index(a, b)]
}

func (m UTMatrix) set(a, b int, v float64) {
	m.d[m.index(a, b)] = v
}
