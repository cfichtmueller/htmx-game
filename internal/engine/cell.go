package engine

import (
	"math/rand"
)

type Cell struct {
	Agent                 *Agent
	Color                 string
	Type                  string
	HandlePlayerCollision func(c *Cell, p *Player)
}

func NewBulletCell() *Cell {
	return &Cell{
		Agent: &Agent{
			X:         0,
			Y:         0,
			Width:     10,
			Height:    10,
			Direction: rand.Float64(),
			Velocity:  10 + 100*rand.Float64(),
			Ages:      true,
			TTL:       15,
		},
		Color: "#fdc82a",
		HandlePlayerCollision: func(c *Cell, p *Player) {
			p.Die()
			c.Die()
		},
	}
}

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

func (c *Cell) Update(dt float64) {
	if c.Agent.Dead {
		return
	}
	c.Agent.Update(dt)
}

func (c *Cell) Die() {
	c.Agent.Dead = true
	c.Agent.Stop()
}
