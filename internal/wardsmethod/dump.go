package wardsmethod

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func Dump(w io.Writer, nodes []Node, alives []int) {
	x := make([]int, len(alives))
	copy(x, alives)
	sort.Ints(x)
	for _, i := range x {
		n := nodes[i]
		l := nodes[n.Left]
		r := nodes[n.Right]
		fmt.Fprintf(w, "node#%d: w=%-4d d=%f  - L(w=%-4d d=%f) - R(w=%-4d d=%f)\n", i, int(n.Weight), n.Delta, int(l.Weight), l.Delta, int(r.Weight), r.Delta)
	}
}

func dumpNode(w io.Writer, tree Tree, n int, depth int) error {
	node := tree[n]
	if node.Left >= 0 {
		err := dumpNode(w, tree, node.Left, depth+1)
		if err != nil {
			return err
		}
	}
	if node.Delta > 0 {
		_, err := fmt.Fprintf(w, "%s#%d w=%d d=%f\n", strings.Repeat("  ", depth), n, int(node.Weight), node.Delta)
		if err != nil {
			return err
		}
	}
	if node.Right >= 0 {
		err := dumpNode(w, tree, node.Right, depth+1)
		if err != nil {
			return err
		}
	}
	return nil
}

func DumpTree(w io.Writer, tree Tree) error {
	return dumpNode(w, tree, len(tree)-1, 0)
}
