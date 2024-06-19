package client

import (
	"fmt"
	"math"
	"strconv"

	"cfichtmueller.com/htmx-game/internal/engine"
)

type State struct {
	width  float64
	height float64
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
		width:  float64(width),
		height: float64(height),
	}, nil
}

func (v *State) Update(s *engine.State) {
	xFactor := v.width / s.Width
	yFactor := v.height / s.Height
	factor := math.Min(xFactor, yFactor)

	v.Cells = make([]Cell, 0, len(s.Cells))

	for _, c := range s.Cells {
		v.Cells = append(v.Cells, cellFromAgent(
			factor,
			c.Color,
			c.Agent,
		))
	}

	for _, p := range s.Players {
		v.Cells = append(v.Cells, cellFromAgent(
			factor,
			p.Color,
			p.Agent,
		))
	}
}

type Cell struct {
	Width    int
	Height   int
	Left     int
	Top      int
	Rotation float64
	Color    string
}

func cellFromAgent(factor float64, color string, a *engine.Agent) Cell {
	w := math.Max(10, factor*a.Width)
	h := math.Max(10, factor*a.Height)
	return Cell{
		Width:    int(w),
		Height:   int(h),
		Left:     int(a.X * factor),
		Top:      int(a.Y * factor),
		Rotation: a.Direction,
		Color:    color,
	}
}
