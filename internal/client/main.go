package client

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"cfichtmueller.com/htmx-game/internal/engine"
	"cfichtmueller.com/htmx-game/internal/engine/physics"
)

const (
	Z_BASE         = 0
	Z_GROUND       = 10
	Z_POWER_UP     = 20
	Z_GROUND_UNIT  = 30
	Z_BUILDING     = 40
	Z_ABOVE_GROUND = 50
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
			Z:     Z_GROUND_UNIT,
		},
		engine.Tank: {
			Alive: "tank.png",
			Dead:  "tank_dead.png",
			Z:     Z_GROUND_UNIT,
		},
		engine.TankShelter: {
			Alive: "shelter.png",
			Z:     Z_BUILDING,
		},
		engine.Tower: {
			Alive: "tower.png",
			Dead:  "tower_dead.png",
			Z:     Z_POWER_UP,
		},
		engine.SpeedPowerUp: {
			Alive: "powerup_speed.png",
			Z:     Z_POWER_UP,
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
	Masks  []Mask
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
		Masks: make([]Mask, 0),
	}, nil
}

func (v *State) Update(e *engine.Engine) {
	v.mu.Lock()
	defer v.mu.Unlock()
	screen := v.Screen
	screen.ScaleTo(e.World.Width, e.World.Height)

	if screen.xOffset > 0 {
		v.Masks = []Mask{
			{
				Left:   0,
				Top:    0,
				Width:  int(screen.xOffset),
				Height: int(screen.Height),
			},
			{
				Left:   int(v.Screen.xOffset) + screen.MapLength(e.World.Width, 0),
				Top:    0,
				Width:  int(v.Screen.xOffset),
				Height: int(v.Screen.Height),
			},
		}
	} else if v.Screen.yOffset > 0 {
		v.Masks = []Mask{
			{
				Left:   0,
				Top:    0,
				Width:  int(v.Screen.Width),
				Height: int(v.Screen.yOffset),
			},
			{
				Left:   0,
				Top:    int(screen.yOffset) + screen.MapLength(e.World.Height, 0),
				Width:  int(v.Screen.Width),
				Height: int(v.Screen.yOffset),
			},
		}
	}

	for _, entity := range e.World.Entities {
		t := e.World.Components.EntityTypes[entity]
		position, hasPosition := e.World.Components.Positions[entity]
		bb, hasBb := e.World.Components.BoundingBoxes[entity]
		if !hasPosition || !hasBb {
			continue
		}

		if position.X < 0 || position.Y < 0 || position.X > e.World.Width || position.Y > e.World.Height {
			continue
		}

		health, hasHealth := e.World.Components.Healths[entity]
		isDead := false
		if hasHealth && health.Dead {
			isDead = true
		}

		r := renders[t.Type]
		zIndex := r.Z
		image := r.Alive
		if isDead {
			if r.Dead != "" {
				image = r.Dead
			}
			zIndex--
		}

		v.Cells = append(v.Cells, Cell{
			Width:    v.Screen.MapLength(bb.W, 10),
			Height:   v.Screen.MapLength(bb.H, 10),
			Left:     v.Screen.MapX(position.X),
			Top:      v.Screen.MapY(position.Y),
			Rotation: position.Direction + physics.Deg90,
			Image:    image,
			ZIndex:   zIndex,
		})
	}

	v.Map = Map{
		Width:  v.Screen.MapLength(e.World.Width, 1),
		Height: v.Screen.MapLength(e.World.Height, 1),
		Left:   v.Screen.MapX(0),
		Top:    v.Screen.MapY(0),
	}
}

type Mask struct {
	Width  int
	Height int
	Left   int
	Top    int
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
