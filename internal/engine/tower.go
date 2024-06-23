package engine

import (
	"math"

	"cfichtmueller.com/htmx-game/internal/engine/bhv"
)

func SpawnTower(world *World, x, y float64) {
	entity := world.AddEntity(Tower)

	world.Components.Healths[entity] = &Health{Decays: true, DecayTTL: 30}
	world.Components.Positions[entity] = &Position{X: x, Y: y}
	world.Components.BoundingBoxes[entity] = &BoundingBox{Width: 30, Height: 30}
	world.Components.Behaviors[entity] = &Behavior{
		Tree: towerBehavior(world, entity),
	}
}

func towerBehavior(world *World, entity Entity) *bhv.Tree {
	return bhv.NewTree(
		bhv.SequenceNode(
			&bhv.Node{
				OnTick: func(n *bhv.Node, dt float64) bhv.Status {
					health := world.Components.Healths[entity]
					if health.Dead {
						return bhv.StatusFailure
					}
					return bhv.StatusSuccess
				},
			},
			bhv.WaitNode(
				&bhv.WaitState{TimeToWaitFn: frandomF(5, 10)},
				AimBehavior(
					&AimState{
						World:             world,
						Entity:            entity,
						TargetDirectionFn: frandomF(0, math.Pi*2),
					},
					bhv.BurstBehavior(
						&bhv.BurstState{Interval: 0.3, BurstSizeFn: irandomF(3, 6)},
						bhv.ActionNode(func(n *bhv.Node, dt float64) bhv.Status {
							towerPos := world.Components.Positions[entity]
							towerBb := world.Components.BoundingBoxes[entity]
							spread := frandom(-0.02, 0.02)
							SpawnBullet(
								world,
								towerPos.X+towerBb.Width/2,
								towerPos.Y+towerBb.Height/2,
								towerPos.Direction+spread,
								70,
								10,
							)
							return bhv.StatusSuccess
						}),
					),
				),
			),
		),
	)
}

type AimState struct {
	World             *World
	Entity            Entity
	TargetDirectionFn func() float64
	isAiming          bool
	hasAimed          bool
	targetDirection   float64
}

func AimBehavior(s *AimState, child *bhv.Node) *bhv.Node {
	return &bhv.Node{
		Data:     s,
		Children: []*bhv.Node{child},
		OnTick: func(n *bhv.Node, dt float64) bhv.Status {
			d := n.Data.(*AimState)

			if !d.isAiming && !d.hasAimed {
				d.targetDirection = d.TargetDirectionFn()
				d.isAiming = true
			}

			if d.isAiming {
				position := d.World.Components.Positions[d.Entity]
				dirDiff := d.targetDirection - position.Direction
				position.Direction += math.Min(dirDiff, math.Pi/6)
				if dirDiff != 0 {
					return bhv.StatusRunning
				}
			}

			d.hasAimed = true
			s := n.Children[0].Tick(dt)
			if s != bhv.StatusSuccess {
				return s
			}

			d.isAiming = false
			d.hasAimed = false
			return bhv.StatusSuccess
		},
	}
}
