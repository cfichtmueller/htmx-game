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

type Placement struct {
	x, y int
}

func PlaceInRaster(count, width, height int, generator func(x, y int)) {
	placements := make([]Placement, count)
	placed := 0
	for placed < count {
		x := irandom(0, width)
		y := irandom(0, height)

		conflict := false
		for _, placement := range placements[:placed] {
			if placement.x == x && placement.y == y {
				conflict = true
				break
			}
		}

		if !conflict {
			placements[placed] = Placement{x: x, y: y}
			placed++
			generator(x, y)
		}
	}
}
