package engine

import "cfichtmueller.com/htmx-game/internal/engine/bhv"

type ComponentStorage struct {
	Accelerations map[Entity]*Acceleration
	Behaviors     map[Entity]*Behavior
	BoundingBoxes map[Entity]*BoundingBox
	Frictions     map[Entity]*Friction
	Healths       map[Entity]*Health
	Positions     map[Entity]*Position
	Velocities    map[Entity]*Velocity
	EntityTypes   map[Entity]*EntityTypeComponent
}

func NewComponentStorage() *ComponentStorage {
	return &ComponentStorage{
		Accelerations: make(map[Entity]*Acceleration),
		Behaviors:     make(map[Entity]*Behavior),
		BoundingBoxes: make(map[Entity]*BoundingBox),
		Frictions:     make(map[Entity]*Friction),
		Healths:       make(map[Entity]*Health),
		Positions:     make(map[Entity]*Position),
		Velocities:    make(map[Entity]*Velocity),
		EntityTypes:   make(map[Entity]*EntityTypeComponent),
	}
}

func (s *ComponentStorage) RemoveEntity(entity Entity) {
	delete(s.Accelerations, entity)
	delete(s.Behaviors, entity)
	delete(s.BoundingBoxes, entity)
	delete(s.Frictions, entity)
	delete(s.Healths, entity)
	delete(s.Positions, entity)
	delete(s.Velocities, entity)
	delete(s.EntityTypes, entity)
}

type Acceleration struct {
	Current        float64
	Max            float64
	AngularCurrent float64
	AngularMax     float64
}

type Behavior struct {
	Tree *bhv.Tree
}

type BoundingBox struct {
	Width, Height float64
}

type EntityTypeComponent struct {
	Type EntityType
}

type Friction struct {
	Current        float64
	AngularCurrent float64
}

type Health struct {
	Current  int
	TTL      float64
	Ages     bool
	Dead     bool
	Decays   bool
	DecayTTL float64
	Decayed  bool
}

type Position struct {
	X, Y, Direction float64
}

type Velocity struct {
	Current        float64
	Max            float64
	AngularCurrent float64
	AngularMax     float64
}
