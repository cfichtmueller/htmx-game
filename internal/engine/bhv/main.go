package bhv

const (
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
)

type Status string

type Tree struct {
	Root *Node
}

func NewTree(root *Node) *Tree {
	return &Tree{
		Root: root,
	}
}

func (t *Tree) Tick(dt float64) {
	if t.Root == nil {
		return
	}
	t.Root.Tick(dt)
}

type Node struct {
	Children []*Node
	Data     any
	OnTick   func(n *Node, dt float64) Status
}

func NewNode() *Node {
	return &Node{
		Children: make([]*Node, 0),
	}
}

func (n *Node) Tick(dt float64) Status {
	if n.OnTick == nil {
		return StatusFailure
	}
	return n.OnTick(n, dt)
}

func (n *Node) AddChild(child *Node) *Node {
	n.Children = append(n.Children, child)
	return n
}

func (n *Node) AddChildren(children ...*Node) *Node {
	n.Children = append(n.Children, children...)
	return n
}

func ActionNode(f func(n *Node, dt float64) Status) *Node {
	return &Node{
		OnTick: f,
	}
}

func SelectorNode(children ...*Node) *Node {
	n := NewNode().AddChildren(children...)
	n.OnTick = selectorFunc
	return n
}

func SequenceNode(children ...*Node) *Node {
	n := NewNode().AddChildren(children...)
	n.OnTick = sequenceFunc
	return n
}

func selectorFunc(n *Node, dt float64) Status {
	for _, c := range n.Children {
		switch c.Tick(dt) {
		case StatusRunning:
			return StatusRunning
		case StatusSuccess:
			return StatusSuccess
		}
	}
	return StatusFailure
}

func sequenceFunc(n *Node, dt float64) Status {
	for _, c := range n.Children {
		switch c.Tick(dt) {
		case StatusRunning:
			return StatusRunning
		case StatusFailure:
			return StatusFailure
		}
	}
	return StatusSuccess
}
