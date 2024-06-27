package physics

func Collides(p1 *Position, s1 *Rectangle, p2 *Position, s2 *Rectangle) bool {
	return p1.X < p2.X+s2.W &&
		p1.X+s1.W > p2.X &&
		p1.Y < p2.Y+s2.H &&
		p1.Y+s1.H > p2.Y
}
