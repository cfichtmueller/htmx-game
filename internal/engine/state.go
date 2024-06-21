package engine

import "sync"

type State struct {
	mu      sync.Mutex
	Width   float64
	Height  float64
	Cells   *CellList
	Players map[string]*Player
}

func NewState(width, height float64) *State {
	return &State{
		Width:   width,
		Height:  height,
		Cells:   NewCellList(),
		Players: map[string]*Player{},
	}
}

func (s *State) Update(dt float64) {
	s.mu.Lock()
	s.Cells.Each(func(c *Cell) { c.Update(dt) })
	s.Cells.Filter(func(c *Cell) bool { return !c.Agent.Dead })

	for _, p := range s.Players {
		p.Update(dt)
		c, ok := intersects(p, s.Cells.Cells)
		if !ok || c.HandlePlayerCollision == nil {
			continue
		}
		c.HandlePlayerCollision(c, p)
		if p.Agent.Decayed {
			delete(s.Players, p.ID)
		}

	}
	s.mu.Unlock()
}

func (s *State) AddCell(c *Cell) {
	s.mu.Lock()
	s.Cells.Add(c)
	s.mu.Unlock()
}

func (s *State) SpawnPlayer() *Player {
	s.mu.Lock()
	p := &Player{
		ID: randomId(),
		Agent: &Agent{
			X:                  s.Width / 2,
			Y:                  s.Height / 2,
			Width:              30,
			Height:             30,
			MaxVelocity:        50,
			MaxAngularVelocity: 10,
			Friction:           10,
			Decays:             true,
			DecayTTL:           10,
		},
		Color: "#3498db",
	}
	s.Players[p.ID] = p
	s.mu.Unlock()
	return p
}

func (s *State) PlayerWithId(id string) *Player {
	p, ok := s.Players[id]
	if !ok {
		return nil
	}
	return p
}
