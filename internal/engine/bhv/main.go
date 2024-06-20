package bhv

const (
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
)

type Status string

type Blackboard struct {
	data map[string]interface{}
}

func NewBlackboard() *Blackboard {
	return &Blackboard{
		data: make(map[string]interface{}),
	}
}

func (b *Blackboard) Set(key string, value interface{}) {
	b.data[key] = value
}

func (b *Blackboard) Get(key string) (interface{}, bool) {
	v, ok := b.data[key]
	return v, ok
}

func (b *Blackboard) MustGet(key string) interface{} {
	v, ok := b.data[key]
	if !ok {
		panic("didn't find key " + key + " on blackboard")
	}
	return v
}

type Tree struct {
	Root *Node
}

func NewTree(root *Node) *Tree {
	return &Tree{
		Root: root,
	}
}

func (t *Tree) Tick(bb *Blackboard) {
	if t.Root == nil {
		return
	}
	t.Root.Tick(bb)
}

type Node struct {
	Children []*Node
	Data     any
	OnTick   func(n *Node, bb *Blackboard) Status
}

func NewNode() *Node {
	return &Node{
		Children: make([]*Node, 0),
	}
}

func (n *Node) Tick(bb *Blackboard) Status {
	if n.OnTick == nil {
		return StatusFailure
	}
	return n.OnTick(n, bb)
}

func (n *Node) AddChild(child *Node) *Node {
	n.Children = append(n.Children, child)
	return n
}

func (n *Node) AddChildren(children ...*Node) *Node {
	n.Children = append(n.Children, children...)
	return n
}

func ActionNode(f func(n *Node, bb *Blackboard) Status) *Node {
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

func selectorFunc(n *Node, bb *Blackboard) Status {
	for _, c := range n.Children {
		switch c.Tick(bb) {
		case StatusRunning:
			return StatusRunning
		case StatusSuccess:
			return StatusSuccess
		}
	}
	return StatusFailure
}

func sequenceFunc(n *Node, bb *Blackboard) Status {
	for _, c := range n.Children {
		switch c.Tick(bb) {
		case StatusRunning:
			return StatusRunning
		case StatusFailure:
			return StatusFailure
		}
	}
	return StatusSuccess
}
