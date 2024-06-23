package engine

func SpawnPlayer(world *World, x, y, direction float64) Entity {
	entity := world.AddEntity(Player)
	world.Components.Positions[entity] = &Position{X: x, Y: y, Direction: direction}
	world.Components.Velocities[entity] = &Velocity{Max: 50, AngularMax: 10}
	world.Components.Frictions[entity] = &Friction{Current: 30}
	world.Components.BoundingBoxes[entity] = &BoundingBox{Width: 30, Height: 30}
	world.Components.Healths[entity] = &Health{Decays: true, DecayTTL: 10}
	return entity
}
