package state

type Player struct {
	ID       string
	Dead     bool
	X        float64
	Y        float64
	Velocity float64
	Color    string
}

func (p *Player) Die() {
	p.Dead = true
	p.Color = "#ff0000"
}
