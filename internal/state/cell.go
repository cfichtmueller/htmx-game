package state

import (
	"math"
	"math/rand"
)

const (
	CELL_TYPE_BULLET         = "bullet"
	CELL_TYPE_POWER_VELOCITY = "powerUpVelocity"
)

type Cell struct {
	X         float64
	Y         float64
	Direction float64
	Velocity  float64
	Dead      bool
	AgeFactor float64
	Age       float64
	Color     string
	Type      string
}

func NewBulletCell() *Cell {
	return &Cell{
		X:         0,
		Y:         0,
		Direction: rand.Float64(),
		Velocity:  10 + 100*rand.Float64(),
		AgeFactor: 1,
		Color:     "#fdc82a",
		Type:      CELL_TYPE_BULLET,
	}
}

func NewVelocityPowerUpCell(x, y float64) *Cell {
	return &Cell{
		X:         x,
		Y:         y,
		Velocity:  0,
		AgeFactor: 0.25,
		Color:     "#006600",
		Type:      CELL_TYPE_POWER_VELOCITY,
	}
}

func (c *Cell) Update(dt float64) {
	if c.Dead {
		return
	}
	if c.Velocity != 0 {
		c.X = c.X + dt*c.Velocity*math.Cos(c.Direction)
		c.Y = c.Y + dt*c.Velocity*math.Sin(c.Direction)
	}
	c.Age = c.Age + dt*c.AgeFactor
}

func (c *Cell) Die() {
	c.Dead = true
}
