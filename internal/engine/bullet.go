package engine

import "cfichtmueller.com/htmx-game/internal/engine/physics"

func SpawnBullet(world *World, x, y, direction, velocity, ttl float64) {
	entity := world.AddEntity(Bullet)

	world.Components.Positions[entity] = &physics.Position{X: x - 5, Y: y - 5, Direction: direction}
	world.Components.Velocities[entity] = &Velocity{Current: velocity}
	world.Components.BoundingBoxes[entity] = &physics.Rectangle{W: 10, H: 10}
	world.Components.Healths[entity] = &Health{Ages: true, TTL: ttl}
}
