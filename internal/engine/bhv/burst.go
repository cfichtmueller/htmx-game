package bhv

import "math"

type BurstState struct {
	BurstSize   int
	BurstSizeFn func() int
	Interval    float64
	remaining   int
	timeToNext  float64
}

func BurstBehavior(s *BurstState, child *Node) *Node {
	return &Node{
		Data:     s,
		Children: []*Node{child},
		OnTick: func(n *Node, dt float64) Status {
			d := n.Data.(*BurstState)
			if d.remaining == 0 {
				d.remaining = d.BurstSize
				if d.BurstSizeFn != nil {
					d.remaining = d.BurstSizeFn()
				}
				d.timeToNext = d.Interval
				return StatusSuccess
			}
			d.timeToNext = math.Max(0, d.timeToNext-dt)
			if d.timeToNext > 0 {
				return StatusRunning
			}
			d.remaining -= 1
			d.timeToNext = d.Interval
			s := n.Children[0].Tick(dt)
			if s != StatusSuccess {
				return s
			}
			return StatusRunning
		},
	}
}
