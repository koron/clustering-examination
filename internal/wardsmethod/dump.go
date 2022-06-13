package wardsmethod

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func dump(w io.Writer, nodes []Node, alives []int) {
	x := make([]int, len(alives))
	copy(x, alives)
	sort.Ints(x)
	for _, n := range x {
		fmt.Fprintf(w, "node#%d: w=%-4d d=%f\n", n, int(nodes[n].Weight), nodes[n].Delta)
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
