package engine

import (
	"cfichtmueller.com/htmx-game/internal/engine/bhv"
	"cfichtmueller.com/htmx-game/internal/engine/physics"
)

func SpawnTankShelter(world *World, x, y, direction float64) {
	entity := world.AddEntity(TankShelter)
	world.Components.Positions[entity] = &physics.Position{X: x, Y: y, Direction: direction}
	world.Components.BoundingBoxes[entity] = &physics.Rectangle{W: 30, H: 30}
	world.Components.Behaviors[entity] = &Behavior{
		Tree: bhv.NewTree(
			bhv.WaitNode(
				&bhv.WaitState{TimeToWaitFn: frandomF(2, 4)},
				&bhv.Node{OnTick: func(n *bhv.Node, dt float64) bhv.Status {
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
	}
}

func SpawnTank(world *World, x, y, direction float64) {
	entity := world.AddEntity(Tank)

	world.Components.AutoMove[entity] = &AutoMove{}
	world.Components.Healths[entity] = &Health{Ages: true, TTL: 30, Decays: true, DecayTTL: 15}
	world.Components.Positions[entity] = &physics.Position{X: x, Y: y, Direction: direction}
	world.Components.BoundingBoxes[entity] = &physics.Rectangle{W: 30, H: 30}
	world.Components.Sensings[entity] = NewSensing().SetRange(Player, 150)
	world.Components.Velocities[entity] = &Velocity{Current: 30, AngularMax: physics.Deg180}
	world.Components.Behaviors[entity] = &Behavior{
		Tree: bhv.NewTree(
			bhv.SelectorNode(
				deadTankBehavior(world, entity),
				tankAvoidWorldBoundariesBehavior(world, entity),
				tankChasePlayerBehavior(world, entity),
				tankMoveRandomlyBehavior(world, entity),
			),
		),
	}
}

func deadTankBehavior(world *World, entity Entity) *bhv.Node {
	return &bhv.Node{
		OnTick: func(n *bhv.Node, dt float64) bhv.Status {
			if world.Components.Healths[entity].Dead {
				return bhv.StatusSuccess
			}
			return bhv.StatusFailure
		},
	}
}

func tankAvoidWorldBoundariesBehavior(world *World, entity Entity) *bhv.Node {
	return &bhv.Node{OnTick: func(n *bhv.Node, dt float64) bhv.Status {
		autoMove := world.Components.AutoMove[entity]
		pos := world.Components.Positions[entity]
		if pos.X < 50 {
			if pos.Y < 50 {
				autoMove.SetTargetDirection(physics.Deg135)
			} else {
				autoMove.SetTargetDirection(physics.Deg0)
			}
			return bhv.StatusSuccess
		}
		if pos.X > world.Width-50 {
			if pos.Y > world.Height-50 {
				autoMove.SetTargetDirection(physics.Deg315)
			} else {
				autoMove.SetTargetDirection(physics.Deg180)
			}
			return bhv.StatusSuccess
		}
		if pos.Y < 50 {
			autoMove.SetTargetDirection(physics.Deg90)
			return bhv.StatusSuccess
		}
		if pos.Y > world.Height-50 {
			autoMove.SetTargetDirection(physics.Deg270)
			return bhv.StatusSuccess
		}
		return bhv.StatusFailure
	}}
}

func tankChasePlayerBehavior(world *World, entity Entity) *bhv.Node {
	return bhv.SelectorNode(
		&bhv.Node{OnTick: func(n *bhv.Node, dt float64) bhv.Status {
			sensings := world.Components.Sensings[entity]
			for _, other := range sensings.SensedEntities {
				if other.Type != Player {
					continue
				}
				pHealth := world.Components.Healths[other.Entity]
				if pHealth.Dead {
					continue
				}
				position := world.Components.Positions[entity]
				autoMove := world.Components.AutoMove[entity]
				autoMove.SetTargetDirection(physics.DirectionTo(position, other.Position))
				return bhv.StatusSuccess
			}
			return bhv.StatusFailure
		}},
	)
}

func tankMoveRandomlyBehavior(world *World, entity Entity) *bhv.Node {
	return bhv.WaitNode(
		&bhv.WaitState{TimeToWaitFn: frandomF(2, 4), InitialWait: 5},
		&bhv.Node{
			OnTick: func(n *bhv.Node, dt float64) bhv.Status {
				position := world.Components.Positions[entity]
				autoMove := world.Components.AutoMove[entity]
				autoMove.SetTargetDirection(position.Direction + frandom(-physics.Deg45, physics.Deg45))
				return bhv.StatusSuccess
			}},
	)
}
