package client

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"cfichtmueller.com/htmx-game/internal/engine"
)

const (
	Z_BASE int = iota
	Z_GROUND
	Z_BUILDING
	Z_ABOVE_GROUND
)

var (
	renders = map[engine.EntityType]EntityRender{
		engine.Bullet: {
			Alive: "bullet.png",
			Z:     Z_ABOVE_GROUND,
		},
		engine.Player: {
			Alive: "player.png",
			Dead:  "player_dead.png",
			Z:     Z_GROUND,
		},
		engine.Tank: {
			Alive: "tank.png",
			Dead:  "tank_dead.png",
			Z:     Z_GROUND,
		},
		engine.TankShelter: {
			Alive: "shelter.png",
			Z:     Z_BUILDING,
		},
		engine.Tower: {
			Alive: "tower.png",
			Dead:  "tower_dead.png",
			Z:     Z_GROUND,
		},
		engine.SpeedPowerUp: {
			Alive: "powerup_speed.png",
			Z:     Z_GROUND,
		},
	}
)

type EntityRender struct {
	Alive string
	Dead  string
	Z     int
}

type State struct {
	mu     sync.Mutex
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
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Screen.ScaleTo(s.Width, s.Height)

	for _, entity := range s.World.Entities {
		t := s.World.Components.EntityTypes[entity]
		position, hasPosition := s.World.Components.Positions[entity]
		bb, hasBb := s.World.Components.BoundingBoxes[entity]
		if !hasPosition || !hasBb {
			continue
		}

		health, hasHealth := s.World.Components.Healths[entity]
		isDead := false
		if hasHealth && health.Dead {
			isDead = true
		}

		r := renders[t.Type]
		image := r.Alive
		if isDead && r.Dead != "" {
			image = r.Dead

		}

		v.Cells = append(v.Cells, Cell{
			Width:    v.Screen.MapLength(bb.Width, 10),
			Height:   v.Screen.MapLength(bb.Height, 10),
			Left:     v.Screen.MapX(position.X),
			Top:      v.Screen.MapY(position.Y),
			Rotation: position.Direction + math.Pi/2,
			Image:    image,
			ZIndex:   r.Z,
		})
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
	Image    string
	ZIndex   int
}
