package wardsmethod

import (
	"sort"
)

// Top returns indexes of top n'th nodes .
func Top(tree Tree, n int) []int {
	buf := make([]int, n*2)
	prev, next := buf[:0], buf[n:n]
	prev = append(prev, len(tree)-1)
	for {
		curr := prev
		for len(curr)+len(next) < n && len(curr) > 0 {
			v := curr[0]
			curr = curr[1:]
			n := tree[v]
			if n.Left < 0 && n.Right < 0 {
				next = append(next, v)
			} else {
				next = append(next, n.Left, n.Right)
			}
		}
		if len(curr)+len(next) >= n {
			next = append(next, curr...)
			sort.Ints(next)
			return next
		}
		sort.SliceStable(next, func(i, j int) bool {
			return tree[i].Weight > tree[j].Weight
		})
		next, prev = prev[:0], next
	}
}
