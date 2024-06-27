package physics

import "math"

// Accelerate computes the new velocity based on initial velocity, max velocity, acceleration and passed time
func Accelerate(v, vmax, a, dt float64) float64 {
	return math.Min(vmax, v+a*dt)
}

// Move computes new x and y coordinates based on direction, speed and passed time
func Move(x, y, direction, v, dt float64) (float64, float64) {
	return x + dt*v*math.Cos(direction), y + dt*v*math.Sin(direction)
}

func Move2(p *Position, v, dt float64) {
	p.X += dt * v * math.Cos(p.Direction)
	p.Y += dt * v * math.Sin(p.Direction)
}
