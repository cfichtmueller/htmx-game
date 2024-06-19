package engine

import "sync"

type State struct {
	mu      sync.Mutex
	Width   float64
	Height  float64
	Cells   []*Cell
	Players map[string]*Player
}

func NewState(width, height float64) *State {
	return &State{
		Width:   width,
		Height:  height,
		Cells:   make([]*Cell, 0),
		Players: map[string]*Player{},
	}
}

func (s *State) Update(dt float64) {
	s.mu.Lock()
	newCells := make([]*Cell, 0, len(s.Cells))
	for _, c := range s.Cells {
		r := c.Update(dt)
		newCells = append(newCells, r.Cells...)
		if !c.Agent.Dead {
			newCells = append(newCells, c)
		}
	}
	s.Cells = newCells
	for _, p := range s.Players {
		p.Update(dt)
		c, ok := intersects(p, s.Cells)
		if !ok || c.HandlePlayerCollision == nil {
			continue
		}
		c.HandlePlayerCollision(c, p)
	}
	s.mu.Unlock()
}

func (s *State) AddCell(c *Cell) {
	s.mu.Lock()
	s.Cells = append(s.Cells, c)
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
		},
		Color: "#ff00ff",
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
