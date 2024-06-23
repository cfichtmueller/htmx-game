package engine

func SetEntityDirection(world *World, entity Entity, d float64) {
	world.Components.Positions[entity].Direction = d
}

func KillEntity(world *World, entity Entity) {
	world.Components.Healths[entity].Dead = true
}

func IsEntityDead(world *World, entity Entity) bool {
	return world.Components.Healths[entity].Dead
}

func SetEntityVelocity(world *World, entity Entity, v float64) {
	velocity := world.Components.Velocities[entity]
	velocity.Current = velocity.Max * v
}
