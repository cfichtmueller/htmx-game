package engine

func NewVelocityPowerUpCell(x, y float64) *Cell {
	return &Cell{
		Agent: StaticAgent(x, y, 10, 10).SetAges(true).SetTTL(30),
		Color: "#006600",
		HandlePlayerCollision: func(c *Cell, p *Player) {
			c.Die()
			p.Agent.MaxVelocity += 5
		},
	}
}
