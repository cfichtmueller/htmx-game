package engine

type Player struct {
	ID    string
	Agent *Agent
	Color string
}

func (p *Player) Update(dt float64) {
	p.Agent.Update(dt)
}

func (p *Player) Die() {
	p.Agent.Dead = true
	p.Color = "#2c3e50"
	p.Agent.Stop()
}
