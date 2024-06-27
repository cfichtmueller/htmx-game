package physics

import "math"

func DirectionTo(p1, p2 *Position) float64 {
	deltaX := p2.X - p1.X
	deltaY := p2.Y - p1.Y

	return math.Atan2(deltaY, deltaX)
}
