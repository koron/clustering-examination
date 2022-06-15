package wardsmethod

import (
	"sort"
)

// Top returns indexes of top n'th nodes .
func Top(tree Tree, n int) []int {
	buf := make([]int, n*2)
	prev, next := buf[:0], buf[n:n]
	prev = append(prev, tree.Root())
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
			return tree[next[i]].Weight > tree[next[j]].Weight
		})
		next, prev = prev[:0], next
	}
}

func Mean(tree Tree, n int) []int {
	ave := tree[tree.Root()].Weight / float64(n)
	small := make([]int, 0, n)
	curr := make([]int, 0, n*2)
	curr = append(curr, tree.Root())
	for len(small)+len(curr) < n && len(curr) > 0 {
		v := curr[0]
		curr = curr[1:]
		n := tree[v]
		if n.Weight < ave || n.Left < 0 && n.Right < 0 {
			small = append(small, v)
			continue
		}
		curr = append(curr, n.Left, n.Right)
		//fmt.Printf("#%d w:%d -> %d,%d\n", v, int(n.Weight), int(tree[n.Left].Weight), int(tree[n.Right].Weight))
		sort.SliceStable(curr, func(i, j int) bool {
			return tree[curr[i]].Weight > tree[curr[j]].Weight
		})
	}
	curr = append(curr, small...)
	sort.Ints(curr)
	return curr
}
