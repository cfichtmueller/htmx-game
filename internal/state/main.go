package state

import (
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
)

var (
	idChars = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	idCap   = len(idChars) - 1
)

type State struct {
	mu      sync.Mutex
	Width   float64
	Height  float64
	Cells   []*Cell
	Players map[string]*Player
}

func New(width, height float64) *State {
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
		c.Update(dt)
		if !c.Agent.Dead {
			newCells = append(newCells, c)
		}
	}
	s.Cells = newCells
	for _, p := range s.Players {
		p.Update(dt)
		c, ok := intersects(p, s.Cells)
		if !ok {
			continue
		}
		switch c.Type {
		case CELL_TYPE_BULLET:
			p.Die()
			c.Die()
		case CELL_TYPE_POWER_VELOCITY:
			c.Die()
			p.Agent.MaxVelocity += 5
		}
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

func intersects(p *Player, cells []*Cell) (*Cell, bool) {
	for _, c := range cells {
		if p.Agent.Intersects(c.Agent) {
			return c, true
		}
	}
	return nil, false
}

func randomId() string {
	result := strings.Builder{}
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(fmt.Errorf("couldn't create random: %v", err))
	}
	for i := 0; i < len(bytes); i++ {
		index := int(bytes[i]) % idCap
		result.WriteString(idChars[index])
	}
	return result.String()
}
