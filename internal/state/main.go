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
	newCells := make([]*Cell, 0, len(s.Cells))
	for _, c := range s.Cells {
		c.Update(dt)
		if c.Age < 15 && !c.Dead {
			newCells = append(newCells, c)
		}
	}
	s.Cells = newCells
	for _, p := range s.Players {
		c, ok := intersects(p, s.Cells)
		if !ok {
			return
		}
		switch c.Type {
		case CELL_TYPE_BULLET:
			p.Die()
			c.Die()
		case CELL_TYPE_POWER_VELOCITY:
			c.Die()
			p.Velocity = p.Velocity + 2
		}
	}
}

func (s *State) AddCell(c *Cell) {
	s.mu.Lock()
	s.Cells = append(s.Cells, c)
	s.mu.Unlock()
}

func (s *State) MoveCell(c *Cell, x, y float64) {
	nx := c.X + x
	ny := c.Y + y
	if nx > s.Width {
		nx = nx - s.Width
	}
	if ny > s.Height {
		ny = ny - s.Height
	}
	c.X = nx
	c.Y = ny
}

func (s *State) SpawnPlayer() *Player {
	s.mu.Lock()
	p := &Player{
		ID:       randomId(),
		X:        s.Width / 2,
		Y:        s.Height / 2,
		Velocity: 5,
		Color:    "#ff00ff",
	}
	s.Players[p.ID] = p
	s.mu.Unlock()
	return p
}

func (s *State) PlayerWithId(id string) *Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.Players[id]
	if !ok {
		return nil
	}
	return p
}

func (s *State) MovePlayer(p *Player, x, y float64) {
	if p.Dead {
		return
	}
	nx := p.X + x*p.Velocity
	ny := p.Y + y*p.Velocity
	if nx > s.Width {
		nx = s.Width
	}
	if ny > s.Height {
		ny = s.Height
	}
	p.X = nx
	p.Y = ny
}

func intersects(p *Player, cells []*Cell) (*Cell, bool) {
	minX := p.X - 20
	maxX := p.X + 20
	minY := p.Y - 20
	maxY := p.Y + 20
	for _, c := range cells {
		if c.X > minX && c.X < maxX && c.Y > minY && c.Y < maxY {
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
