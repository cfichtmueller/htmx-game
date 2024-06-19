package client

import (
	"fmt"
	"strconv"

	"cfichtmueller.com/htmx-game/internal/state"
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

func (v *State) Update(s *state.State) {
	xFactor := v.width / s.Width
	yFactor := v.height / s.Height

	v.Cells = make([]Cell, 0, len(s.Cells))

	for _, c := range s.Cells {
		w := xFactor
		if w < 10 {
			w = 10
		}
		h := yFactor
		if h < 10 {
			h = 10
		}
		cc := Cell{
			Width:  int(w),
			Height: int(h),
			Left:   int(c.X * xFactor),
			Top:    int(c.Y * yFactor),
			Color:  c.Color,
		}
		v.Cells = append(v.Cells, cc)
	}

	for _, p := range s.Players {
		v.Cells = append(v.Cells, Cell{
			Width:  30,
			Height: 30,
			Left:   int(p.X * xFactor),
			Top:    int(p.Y * yFactor),
			Color:  p.Color,
		})
	}
}

type Cell struct {
	Width  int
	Height int
	Left   int
	Top    int
	Color  string
}
