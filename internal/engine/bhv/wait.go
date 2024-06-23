package bhv

import "math"

type WaitState struct {
	InitialWait   float64
	TimeToWait    float64
	TimeToWaitFn  func() float64
	WaitState     Status
	timeRemaining float64
}

func WaitNode(s *WaitState, child *Node) *Node {
	if s.WaitState == "" {
		s.WaitState = StatusRunning
	}
	if s.InitialWait > 0 {
		s.timeRemaining = s.InitialWait
	}
	return &Node{
		Data:     s,
		Children: []*Node{child},
		OnTick: func(n *Node, dt float64) Status {
			d := n.Data.(*WaitState)
			d.timeRemaining = math.Max(0, d.timeRemaining-dt)
			if d.timeRemaining > 0 {
				return d.WaitState
			}
			s := n.Children[0].Tick(dt)
			if s != StatusSuccess {
				return s
			}
			if d.TimeToWaitFn == nil {
				d.timeRemaining = d.TimeToWait
			} else {
				d.timeRemaining = d.TimeToWaitFn()
			}
			return StatusSuccess
		},
	}
}
