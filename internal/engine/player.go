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

type PlayerList struct {
	Players []*Player
	idIndex map[string]*Player
}

func NewPlayerList() *PlayerList {
	return &PlayerList{
		Players: make([]*Player, 0),
		idIndex: make(map[string]*Player),
	}
}

func (l *PlayerList) Add(p *Player) {
	l.idIndex[p.ID] = p
	l.Players = append(l.Players, p)
}

func (l *PlayerList) Each(f func(p *Player)) {
	for _, p := range l.Players {
		f(p)
	}
}

func (l *PlayerList) Filter(f func(p *Player) bool) {
	newPlayers := make([]*Player, 0, len(l.Players))
	for _, p := range l.Players {
		if f(p) {
			newPlayers = append(newPlayers, p)
		} else {
			delete(l.idIndex, p.ID)
		}
	}
	l.Players = newPlayers
}

func (l *PlayerList) FindWithId(id string) *Player {
	return l.idIndex[id]
}
