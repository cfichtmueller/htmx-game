package engine

import (
	"math"
	"math/rand"

	"cfichtmueller.com/htmx-game/internal/engine/bhv"
)

type TowerCellData struct {
	behavior *bhv.Tree
}

func NewTowerCell(x, y float64) *Cell {
	return &Cell{
		Agent: &Agent{
			X:                  x,
			Y:                  y,
			Width:              30,
			Height:             30,
			MaxAngularVelocity: math.Pi,
		},
		Type:  "tower",
		Color: "#e67e22",
		Data: &TowerCellData{
			behavior: towerBehavor(),
		},
		HandleUpdate: func(c *Cell, dt float64) CellUpdateResult {
			cellList := NewCellList()
			bb := bhv.NewBlackboard()
			bb.Set("agent", c.Agent)
			bb.Set("dt", dt)
			bb.Set("cells", cellList)

			d := c.Data.(*TowerCellData)
			d.behavior.Tick(bb)
			return CellUpdateResult{Cells: cellList.Cells}
		},
		HandlePlayerCollision: func(c *Cell, p *Player) {
			c.Die()
		},
	}
}

func generateBullets(count int, x, y, direction, spread, velocity, ttl float64) []*Cell {
	res := make([]*Cell, count)
	for i := 0; i < count; i++ {
		res[i] = NewBulletCell(x, y, direction+spread*(rand.Float64()-0.5), velocity, ttl)
	}
	return res
}

func towerBehavor() *bhv.Tree {
	return bhv.NewTree(
		waitNode(
			&WaitState{TimeToWaitFn: frandomF(5, 10)},
			AimBehavior(
				&AimState{TargetDirectionFn: frandomF(0, math.Pi*2)},
				BurstBehavior(
					&BurstState{Interval: 0.3, BurstSizeFn: irandomF(3, 6)},
					bhv.ActionNode(func(n *bhv.Node, bb *bhv.Blackboard) bhv.Status {
						agent := getAgent(bb)
						bb.MustGet("cells").(*CellList).Add(generateBullets(
							1,
							agent.X+agent.Width/2,
							agent.Y+agent.Height/2,
							agent.Direction,
							0.02*rand.Float64(),
							70,
							10,
						)...)
						return bhv.StatusSuccess
					}),
				),
			),
		),
	)
}

type WaitState struct {
	TimeToWait    float64
	TimeToWaitFn  FFunc
	WaitState     bhv.Status
	timeRemaining float64
}

func waitNode(s *WaitState, child *bhv.Node) *bhv.Node {
	if s.WaitState == "" {
		s.WaitState = bhv.StatusRunning
	}
	return &bhv.Node{
		Data:     s,
		Children: []*bhv.Node{child},
		OnTick: func(n *bhv.Node, bb *bhv.Blackboard) bhv.Status {
			d := n.Data.(*WaitState)
			d.timeRemaining = math.Max(0, d.timeRemaining-getDt(bb))
			if d.timeRemaining > 0 {
				return d.WaitState
			}
			s := n.Children[0].Tick(bb)
			if s != bhv.StatusSuccess {
				return s
			}
			if d.TimeToWaitFn == nil {
				d.timeRemaining = d.TimeToWait
			} else {
				d.timeRemaining = d.TimeToWaitFn()
			}
			return bhv.StatusSuccess
		},
	}
}

type AimState struct {
	TargetDirectionFn func() float64
	isAiming          bool
	hasAimed          bool
}

func AimBehavior(s *AimState, child *bhv.Node) *bhv.Node {
	return &bhv.Node{
		Data:     s,
		Children: []*bhv.Node{child},
		OnTick: func(n *bhv.Node, bb *bhv.Blackboard) bhv.Status {
			d := n.Data.(*AimState)
			agent := getAgent(bb)

			if !d.isAiming {
				agent.SetTargetDirection(d.TargetDirectionFn())
				d.isAiming = true
			}

			if !d.hasAimed && agent.IsAutoRotating {
				return bhv.StatusRunning
			}

			d.hasAimed = true
			s := n.Children[0].Tick(bb)
			if s != bhv.StatusSuccess {
				return s
			}

			d.isAiming = false
			d.hasAimed = false
			return bhv.StatusSuccess
		},
	}
}

type BurstState struct {
	BurstSize   int
	BurstSizeFn IFunc
	Interval    float64
	remaining   int
	timeToNext  float64
}

func BurstBehavior(s *BurstState, child *bhv.Node) *bhv.Node {
	return &bhv.Node{
		Data:     s,
		Children: []*bhv.Node{child},
		OnTick: func(n *bhv.Node, bb *bhv.Blackboard) bhv.Status {
			d := n.Data.(*BurstState)
			if d.remaining == 0 {
				d.remaining = d.BurstSize
				if d.BurstSizeFn != nil {
					d.remaining = d.BurstSizeFn()
				}
				d.timeToNext = d.Interval
				return bhv.StatusSuccess
			}
			d.timeToNext = math.Max(0, d.timeToNext-getDt(bb))
			if d.timeToNext > 0 {
				return bhv.StatusRunning
			}
			d.remaining -= 1
			d.timeToNext = d.Interval
			s := n.Children[0].Tick(bb)
			if s != bhv.StatusSuccess {
				return s
			}
			return bhv.StatusRunning
		},
	}
}

func getDt(bb *bhv.Blackboard) float64 {
	return bb.MustGet("dt").(float64)
}

func getAgent(bb *bhv.Blackboard) *Agent {
	return bb.MustGet("agent").(*Agent)
}
