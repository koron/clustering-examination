package graphmatrix

import "gonum.org/v1/gonum/graph"

type Undirect struct {
	n int
	d []bool
}

func NewUndirect(n int) Undirect {
	size := n * (n + 1) / 2
	return Undirect{
		n: n,
		d: make([]bool, size),
	}
}

func (u Undirect) index(a, b int) int {
	if a > b {
		a, b = b, a
	}
	return ((2*u.n-a-1)*a)/2 + b
}

func (u Undirect) Get(a, b int) bool {
	return u.d[u.index(a, b)]
}

func (u Undirect) Set(a, b int, v bool) {
	u.d[u.index(a, b)] = v
}

func (u Undirect) Node(id int64) graph.Node {
	if id < 0 || id >= int64(u.n) {
		return nil
	}
	return undirectNode(id)
}

func (u Undirect) Nodes() graph.Nodes {
	nodes := make([]int, u.n)
	for i := 0; i < u.n; i++ {
		nodes[i] = i
	}
	return newUndirectNodes(nodes)
}

func (u Undirect) From(id int64) graph.Nodes {
	from := int(id)
	nodes := make([]int, 0, u.n-1)
	for i := 0; i < u.n; i++ {
		if i == from {
			continue
		}
		if u.Get(from, i) {
			nodes = append(nodes, i)
		}
	}
	return newUndirectNodes(nodes)
}

func (u Undirect) HasEdgeBetween(xid, yid int64) bool {
	a, b := int(xid), int(yid)
	return !u.Get(a, b)
}

func (u Undirect) Edge(uid, vid int64) graph.Edge {
	return u.EdgeBetween(uid, vid)
}

func (u Undirect) EdgeBetween(xid, yid int64) graph.Edge {
	a, b := int(xid), int(yid)
	if !u.Get(a, b) {
		return nil
	}
	return undirectEdge{a, b}
}

var _ graph.Undirected = Undirect{}

type undirectEdge []int

func (e undirectEdge) From() graph.Node {
	return undirectNode(e[0])
}

func (e undirectEdge) To() graph.Node {
	return undirectNode(e[1])
}

func (e undirectEdge) ReversedEdge() graph.Edge {
	return undirectEdge{e[1], e[0]}
}

var _ graph.Edge = undirectEdge{}

type undirectNode int

func (n undirectNode) ID() int64 {
	return int64(n)
}

var _ graph.Node = undirectNode(0)

type undirectNodes struct {
	nodes []int
	curr  int
}

func newUndirectNodes(nodes []int) *undirectNodes {
	return &undirectNodes{
		nodes: nodes,
		curr:  -1,
	}
}

func (i *undirectNodes) Node() graph.Node {
	if i.curr < 0 || i.curr >= len(i.nodes) {
		return nil
	}
	return undirectNode(i.nodes[i.curr])
}

func (i *undirectNodes) Next() bool {
	if i.curr < len(i.nodes) {
		i.curr++
	}
	return i.curr < len(i.nodes)
}

func (i *undirectNodes) Len() int {
	return len(i.nodes) - i.curr
}

func (i *undirectNodes) Reset() {
	i.curr = -1
}

var _ graph.Nodes = (*undirectNodes)(nil)
