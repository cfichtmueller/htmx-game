package engine

type Player struct {
	ID    string
	Dead  bool
	Agent *Agent
	Color string
}

func (p *Player) Update(dt float64) {
	if p.Dead {
		return
	}
	p.Agent.Update(dt)
}

func (p *Player) Die() {
	p.Dead = true
	p.Color = "#2c3e50"
	p.Agent.Stop()
}
