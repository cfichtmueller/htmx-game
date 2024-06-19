package state

import (
	"math"
)

type Player struct {
	ID          string
	Dead        bool
	X           float64
	Y           float64
	Velocity    float64
	MaxVelocity float64
	Direction   float64
	Color       string
}

func (p *Player) Update(dt float64) {
	if p.Dead {
		return
	}
	if p.Velocity != 0 {
		p.X = p.X + dt*p.Velocity*math.Cos(p.Direction)
		p.Y = p.Y + dt*p.Velocity*math.Sin(p.Direction)
	}
	p.Velocity = math.Max(p.Velocity-2, 0)
}

func (p *Player) AcceptMoveInput(x, y float64) {
	p.Direction = math.Atan2(y, x)
	if x != 0 || y != 0 {
		p.Velocity = p.MaxVelocity
	} else {
		p.Velocity = 0
	}
}

func (p *Player) Die() {
	p.Dead = true
	p.Color = "#ff0000"
}
