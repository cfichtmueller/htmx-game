package engine

import (
	"cfichtmueller.com/htmx-game/internal/engine/bhv"
	"cfichtmueller.com/htmx-game/internal/engine/physics"
)

type ComponentStorage struct {
	Accelerations map[Entity]*Acceleration
	AutoMove      map[Entity]*AutoMove
	Behaviors     map[Entity]*Behavior
	BoundingBoxes map[Entity]*physics.Rectangle
	Frictions     map[Entity]*Friction
	Healths       map[Entity]*Health
	Positions     map[Entity]*physics.Position
	Sensings      map[Entity]*Sensing
	Velocities    map[Entity]*Velocity
	EntityTypes   map[Entity]*EntityTypeComponent
}

func NewComponentStorage() *ComponentStorage {
	return &ComponentStorage{
		Accelerations: make(map[Entity]*Acceleration),
		AutoMove:      make(map[Entity]*AutoMove),
		Behaviors:     make(map[Entity]*Behavior),
		BoundingBoxes: make(map[Entity]*physics.Rectangle),
		Frictions:     make(map[Entity]*Friction),
		Healths:       make(map[Entity]*Health),
		Positions:     make(map[Entity]*physics.Position),
		Sensings:      make(map[Entity]*Sensing),
		Velocities:    make(map[Entity]*Velocity),
		EntityTypes:   make(map[Entity]*EntityTypeComponent),
	}
}

func (s *ComponentStorage) RemoveEntity(entity Entity) {
	delete(s.Accelerations, entity)
	delete(s.AutoMove, entity)
	delete(s.Behaviors, entity)
	delete(s.BoundingBoxes, entity)
	delete(s.Frictions, entity)
	delete(s.Healths, entity)
	delete(s.Positions, entity)
	delete(s.Sensings, entity)
	delete(s.Velocities, entity)
	delete(s.EntityTypes, entity)
}

type Acceleration struct {
	Current        float64
	Max            float64
	AngularCurrent float64
	AngularMax     float64
}

type AutoMove struct {
	TargetDirection       float64
	TargetDirectionActive bool
}

func (a *AutoMove) SetTargetDirection(d float64) {
	a.TargetDirection = d
	a.TargetDirectionActive = true
}

type Behavior struct {
	Tree *bhv.Tree
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

type SensedEntity struct {
	Entity   Entity
	Type     EntityType
	Position *physics.Position
}

type Sensing struct {
	SensedEntities []SensedEntity
	Ranges         map[EntityType]float64
}

func NewSensing() *Sensing {
	return &Sensing{
		Ranges: make(map[EntityType]float64),
	}
}

func (s *Sensing) SetRange(t EntityType, d float64) *Sensing {
	s.Ranges[t] = d
	return s
}

type Velocity struct {
	Current        float64
	Max            float64
	AngularCurrent float64
	AngularMax     float64
}
