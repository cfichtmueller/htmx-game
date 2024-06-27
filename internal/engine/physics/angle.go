package physics

func NormalizeAngle(angle float64) float64 {
	for angle < Deg0 {
		angle += Deg360
	}
	for angle >= Deg360 {
		angle -= Deg360
	}
	return angle
}

func ShortesRotationDirection(current, target float64) float64 {
	delta := NormalizeAngle(target - current)
	if delta > Deg180 {
		return delta - Deg360
	}
	return delta
}
