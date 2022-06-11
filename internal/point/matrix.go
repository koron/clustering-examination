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
	// The formula before transform:
	//		x := m.n - a
	//		len(m.d) - (x * (x + 1) / 2) + (b - a)
	// This is same with:
	//		m.n*a - a*(a+1)/2 + b
	return ((2*m.n-a-1)*a)/2 + b
}

func (m UTMatrix) Get(a, b int) float64 {
	return m.d[m.index(a, b)]
}

func (m UTMatrix) Set(a, b int, v float64) {
	m.d[m.index(a, b)] = v
}
