package bhv

func FailureNode() *Node {
	return &Node{
		OnTick: func(n *Node, dt float64) Status {
			return StatusFailure
		},
	}
}
