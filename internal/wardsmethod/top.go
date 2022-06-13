package wardsmethod

import (
	"container/list"
	"sort"
)

// Top returns indexes of top n'th nodes .
func Top(tree Tree, n int) []int {
	nn := make([]int, 0, n)
	l := list.New()
	l.PushBack(len(tree) - 1)
	for l.Len()+len(nn) < n && l.Len() > 0 {
		f := l.Front()
		curr := f.Value.(int)
		n := tree[curr]
		if n.Left < 0 && n.Right < 0 {
			nn = append(nn, curr)
		} else {
			l.PushBack(n.Left)
			l.PushBack(n.Right)
		}
		l.Remove(f)
	}
	for e := l.Front(); e != nil; e = e.Next() {
		nn = append(nn, e.Value.(int))
	}
	sort.Ints(nn)
	return nn
}

func Top2(tree Tree, n int) []int {
	buf := make([]int, n*2)
	curr, next := buf[:0], buf[n:n]
	curr = append(curr, len(tree)-1)
	for {
		for _, v := range curr {
			if len(next) >= n {
				break
			}
			n := tree[v]
			if n.Left < 0 && n.Right < 0 {
				next = append(next, v)
			} else {
				next = append(next, n.Left, n.Right)
			}
		}
		if len(next) >= n {
			sort.Ints(next)
			return next
		}
		sort.SliceStable(next, func(i, j int) bool {
			return tree[i].Weight > tree[j].Weight
		})
		next, curr = curr[:0], next
	}
}
