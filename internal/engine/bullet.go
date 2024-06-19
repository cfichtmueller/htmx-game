package engine

func NewBulletCell(x, y, direction, velocity, ttl float64) *Cell {
	return &Cell{
		Agent: &Agent{
			X:         x,
			Y:         y,
			Width:     10,
			Height:    10,
			Direction: direction,
			Velocity:  velocity,
			Ages:      true,
			TTL:       ttl,
		},
		Color: "#fdc82a",
		Type:  "bullet",
		HandlePlayerCollision: func(c *Cell, p *Player) {
			p.Die()
			c.Die()
		},
	}
}
