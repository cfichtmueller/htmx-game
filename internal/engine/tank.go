package engine

import (
	"math"

	"cfichtmueller.com/htmx-game/internal/engine/bhv"
)

func SpawnTankShelter(world *World, x, y, direction float64) {
	entity := world.AddEntity(TankShelter)
	world.Components.Positions[entity] = &Position{X: 0, Y: 2 * world.height / 3}
	world.Components.BoundingBoxes[entity] = &BoundingBox{Width: 30, Height: 30}
	world.Components.Behaviors[entity] = &Behavior{
		Tree: bhv.NewTree(
			waitNode(
				&WaitState{TimeToWaitFn: frandomF(2, 4)},
				&bhv.Node{OnTick: func(n *bhv.Node, bb *bhv.Blackboard) bhv.Status {
					SpawnTank(
						world,
						x,
						y,
						direction,
					)
					return bhv.StatusSuccess
				}},
			),
		),
		BbFunc: func(dt float64) *bhv.Blackboard {
			bb := bhv.NewBlackboard()
			bb.Set("dt", dt)
			return bb
		},
	}
}

func SpawnTank(world *World, x, y, direction float64) {
	entity := world.AddEntity(Tank)

	world.Components.Healths[entity] = &Health{Ages: true, TTL: 30, Decays: true, DecayTTL: 15}
	world.Components.Positions[entity] = &Position{X: 0, Y: 2 * world.height / 3}
	world.Components.BoundingBoxes[entity] = &BoundingBox{Width: 30, Height: 30}
	world.Components.Velocities[entity] = &Velocity{Current: 30}
	world.Components.Behaviors[entity] = &Behavior{
		Tree: bhv.NewTree(
			bhv.SequenceNode(
				&bhv.Node{
					OnTick: func(n *bhv.Node, bb *bhv.Blackboard) bhv.Status {
						if world.Components.Healths[entity].Dead {
							return bhv.StatusFailure
						}
						return bhv.StatusSuccess
					},
				},
				waitNode(
					&WaitState{TimeToWaitFn: frandomF(2, 4), InitialWait: 5},
					&bhv.Node{
						OnTick: func(n *bhv.Node, bb *bhv.Blackboard) bhv.Status {
							world.Components.Positions[entity].Direction += frandom(-math.Pi/2, math.Pi/2)
							return bhv.StatusSuccess
						}},
				),
			),
		),
		BbFunc: func(dt float64) *bhv.Blackboard {
			bb := bhv.NewBlackboard()
			bb.Set("dt", dt)
			return bb
		},
	}
}
