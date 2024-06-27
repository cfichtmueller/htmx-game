package engine

import (
	"cfichtmueller.com/htmx-game/internal/engine/bhv"
	"cfichtmueller.com/htmx-game/internal/engine/physics"
)

func SpawnTower(world *World, x, y float64) {
	entity := world.AddEntity(Tower)

	world.Components.AutoMove[entity] = &AutoMove{}
	world.Components.Healths[entity] = &Health{Decays: true, DecayTTL: 30}
	world.Components.Positions[entity] = &physics.Position{X: x, Y: y}
	world.Components.BoundingBoxes[entity] = &physics.Rectangle{W: 30, H: 30}
	world.Components.Velocities[entity] = &Velocity{AngularMax: physics.Deg90}
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
						TargetDirectionFn: frandomF(physics.Deg0, physics.Deg360),
					},
					bhv.BurstBehavior(
						&bhv.BurstState{Interval: 0.3, BurstSizeFn: irandomF(3, 6)},
						bhv.ActionNode(func(n *bhv.Node, dt float64) bhv.Status {
							towerPos := world.Components.Positions[entity]
							towerBb := world.Components.BoundingBoxes[entity]
							spread := frandom(-0.02, 0.02)
							SpawnBullet(
								world,
								towerPos.X+towerBb.W/2,
								towerPos.Y+towerBb.H/2,
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
}

func AimBehavior(s *AimState, child *bhv.Node) *bhv.Node {
	return &bhv.Node{
		Data:     s,
		Children: []*bhv.Node{child},
		OnTick: func(n *bhv.Node, dt float64) bhv.Status {
			d := n.Data.(*AimState)
			autoMove := d.World.Components.AutoMove[d.Entity]

			if !d.isAiming && !d.hasAimed {
				autoMove.SetTargetDirection(d.TargetDirectionFn())
				d.isAiming = true
			}

			if d.isAiming {
				if autoMove.TargetDirectionActive {
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
