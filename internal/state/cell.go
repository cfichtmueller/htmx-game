package state

import (
	"math/rand"
)

const (
	CELL_TYPE_BULLET         = "bullet"
	CELL_TYPE_POWER_VELOCITY = "powerUpVelocity"
)

type Cell struct {
	Agent *Agent
	Color string
	Type  string
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
		Type:  CELL_TYPE_BULLET,
	}
}

func NewVelocityPowerUpCell(x, y float64) *Cell {
	return &Cell{
		Agent: StaticAgent(x, y, 10, 10).SetAges(true).SetTTL(30),
		Color: "#006600",
		Type:  CELL_TYPE_POWER_VELOCITY,
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
