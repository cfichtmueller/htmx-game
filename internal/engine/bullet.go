package engine

func SpawnBullet(world *World, x, y, direction, velocity, ttl float64) {
	entity := world.AddEntity(Bullet)

	world.Components.Positions[entity] = &Position{X: x - 5, Y: y - 5, Direction: direction}
	world.Components.Velocities[entity] = &Velocity{Current: velocity}
	world.Components.BoundingBoxes[entity] = &BoundingBox{Width: 10, Height: 10}
	world.Components.Healths[entity] = &Health{Ages: true, TTL: ttl}
}
