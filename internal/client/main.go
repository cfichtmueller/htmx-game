package client

import (
	"fmt"
	"math"
	"strconv"

	"cfichtmueller.com/htmx-game/internal/engine"
)

type State struct {
	Screen *Screen
	Map    Map
	Cells  []Cell
}

func NewState(w, h string) (*State, error) {
	width, err := strconv.Atoi(w)
	if err != nil {
		return nil, fmt.Errorf("unable to parse width: %v", err)
	}
	height, err := strconv.Atoi(h)
	if err != nil {
		return nil, fmt.Errorf("unable to parse height: %v", err)
	}

	return &State{
		Screen: &Screen{
			Width:  float64(width),
			Height: float64(height),
			Factor: 1,
		},
	}, nil
}

func (v *State) Update(s *engine.State) {
	v.Screen.ScaleTo(s.Width, s.Height)

	v.Cells = make([]Cell, 0, s.Cells.Length())

	s.Cells.Each(func(c *engine.Cell) {
		v.Cells = append(v.Cells, cellFromAgent(
			v.Screen,
			c.Color,
			c.Agent,
			c.Type,
		))
	})

	for _, p := range s.Players {
		v.Cells = append(v.Cells, cellFromAgent(
			v.Screen,
			p.Color,
			p.Agent,
			"player",
		))
	}

	v.Map = Map{
		Width:  v.Screen.MapLength(s.Width, 1),
		Height: v.Screen.MapLength(s.Height, 1),
		Left:   v.Screen.MapX(0),
		Top:    v.Screen.MapY(0),
	}
}

type Screen struct {
	Width   float64
	Height  float64
	Factor  float64
	xOffset float64
	yOffset float64
}

func (s *Screen) ScaleTo(width, height float64) {
	oAspect := width / height
	myAspect := s.Width / s.Height
	xFactor := s.Width / width
	yFactor := s.Height / height
	s.Factor = math.Min(xFactor, yFactor)
	if oAspect > myAspect {
		s.xOffset = 0
		s.yOffset = s.Height/2 - s.Factor*height/2
	} else {
		s.xOffset = s.Width/2 - s.Factor*width/2
		s.yOffset = 0
	}
}

func (s *Screen) MapLength(l float64, min float64) int {
	return int(math.Max(min, s.Factor*l))
}

func (s *Screen) MapX(x float64) int {
	return int(x*s.Factor + s.xOffset)
}

func (s *Screen) MapY(y float64) int {
	return int(y*s.Factor + s.yOffset)
}

type Map struct {
	Width  int
	Height int
	Left   int
	Top    int
}

type Cell struct {
	Width    int
	Height   int
	Left     int
	Top      int
	Rotation float64
	Color    string
	Type     string
}

func cellFromAgent(s *Screen, color string, a *engine.Agent, t string) Cell {
	return Cell{
		Width:    s.MapLength(a.Width, 10),
		Height:   s.MapLength(a.Height, 10),
		Left:     s.MapX(a.X),
		Top:      s.MapY(a.Y),
		Rotation: a.Direction + math.Pi/2,
		Color:    color,
		Type:     t,
	}
}
