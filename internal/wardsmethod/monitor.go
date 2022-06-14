package wardsmethod

type Monitor interface {
	Monitor(nodes []Node, alives []int)
}

type MonitorFunc func([]Node, []int)

func (f MonitorFunc) Monitor(nodes []Node, alives []int) {
	f(nodes, alives)
}
