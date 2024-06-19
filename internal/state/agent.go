package state

import "math"

type Agent struct {
	X                  float64
	Y                  float64
	Width              float64
	Height             float64
	Acceleration       float64
	Velocity           float64
	Friction           float64
	MaxVelocity        float64
	AngularVelocity    float64
	MaxAngularVelocity float64
	Direction          float64
	TTL                float64
	Ages               bool
	Dead               bool
}

func StaticAgent(x, y, width, height float64) *Agent {
	return &Agent{
		X:           x,
		Y:           y,
		MaxVelocity: 0,
	}
}

func (a *Agent) Update(dt float64) {
	a.Direction += a.AngularVelocity * dt
	a.Velocity += a.Acceleration * dt
	if a.Velocity > 0 {
		a.Velocity = math.Max(0, a.Velocity-a.Friction*dt)
	} else {
		a.Velocity = math.Min(0, a.Velocity+a.Friction*dt)
	}
	a.X = a.X + dt*a.Velocity*math.Cos(a.Direction)
	a.Y = a.Y + dt*a.Velocity*math.Sin(a.Direction)
	if a.Ages {
		a.TTL = math.Max(0, a.TTL-dt)
	}
	if a.TTL == 0 {
		a.Dead = true
	}
}

func (a *Agent) SetAcceleration(acc float64) {
	a.Acceleration = acc
}

func (a *Agent) SetVelocity(v float64) {
	a.Velocity = v
	a.Acceleration = 0
}

func (a *Agent) Rotate(radians float64) {
	a.Direction += radians
}

func (a *Agent) Stop() {
	a.Velocity = 0
	a.Acceleration = 0
}

func (a *Agent) Intersects(other *Agent) bool {
	aminx := a.X
	amaxx := a.X + a.Width
	ominx := other.X
	omaxx := other.X + other.Width

	if amaxx < ominx {
		return false
	}
	if aminx > omaxx {
		return false
	}

	aminy := a.Y
	amaxy := a.Y + a.Height
	ominy := other.Y
	omaxy := other.Y + other.Height

	if amaxy < ominy {
		return false
	}

	if aminy > omaxy {
		return false
	}

	return true
}

func (a *Agent) SetAges(ages bool) *Agent {
	a.Ages = ages
	return a
}

func (a *Agent) SetTTL(ttl float64) *Agent {
	a.TTL = ttl
	return a
}
