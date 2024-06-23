package engine

import (
	"math"

	"cfichtmueller.com/htmx-game/internal/engine/bhv"
)

type System interface {
	Update(entities []Entity, components *ComponentStorage, dt float64)
}

type BehaviorSystem struct{}

func NewBehaviorSystem() *BehaviorSystem {
	return &BehaviorSystem{}
}

func (s *BehaviorSystem) Update(entities []Entity, components *ComponentStorage, dt float64) {
	for _, entity := range entities {
		behavior, hasBehavior := components.Behaviors[entity]
		if !hasBehavior {
			continue
		}
		behavior.Tree.Tick(dt)
	}
}

type Collision struct {
	EntityA, EntityB Entity
}

type CollisionHandler interface {
	HandleCollision(entityA, entityB Entity, components *ComponentStorage)
}

type CollisionDetectionSystem struct {
	collisions []Collision
	handlers   map[EntityType]map[EntityType]CollisionHandler
}

func NewCollisionDetectionSystem() *CollisionDetectionSystem {
	return &CollisionDetectionSystem{
		collisions: make([]Collision, 0),
		handlers:   make(map[EntityType]map[EntityType]CollisionHandler),
	}
}

func (s *CollisionDetectionSystem) RegisterHandler(typeA, typeB EntityType, handler CollisionHandler) {
	if s.handlers[typeA] == nil {
		s.handlers[typeA] = make(map[EntityType]CollisionHandler)
	}
	s.handlers[typeA][typeB] = handler
}

func (s *CollisionDetectionSystem) Update(entities []Entity, components *ComponentStorage, dt float64) {
	s.collisions = s.collisions[:0]

	for a := 0; a < len(entities); a++ {
		for b := a + 1; b < len(entities); b++ {
			entityA := entities[a]
			entityB := entities[b]

			posA, hasPosA := components.Positions[entityA]
			bbA, hasBBA := components.BoundingBoxes[entityA]
			posB, hasPosB := components.Positions[entityB]
			bbB, hasBBB := components.BoundingBoxes[entityB]

			if hasPosA &&
				hasBBA &&
				hasPosB &&
				hasBBB &&
				posA.X < posB.X+bbB.Width &&
				posA.X+bbA.Width > posB.X &&
				posA.Y < posB.Y+bbB.Height &&
				posA.Y+bbA.Height > posB.Y {
				s.collisions = append(s.collisions, Collision{EntityA: entityA, EntityB: entityB})
				s.handleCollision(entityA, entityB, components)
			}
		}
	}
}

func (s *CollisionDetectionSystem) handleCollision(entityA, entityB Entity, components *ComponentStorage) {
	typeAComp, hasTypeA := components.EntityTypes[entityA]
	typeBComp, hasTypeB := components.EntityTypes[entityB]

	if !hasTypeA || !hasTypeB {
		return
	}

	typeA := typeAComp.Type
	typeB := typeBComp.Type

	if handler, ok := s.handlers[typeA][typeB]; ok {
		handler.HandleCollision(entityA, entityB, components)
	} else if handler, ok := s.handlers[typeB][typeA]; ok {
		handler.HandleCollision(entityB, entityA, components)
	}
}

func (s *CollisionDetectionSystem) Collisions() []Collision {
	return s.collisions
}

type HealthSystem struct {
	world *World
}

func NewHealthSystem(world *World) *HealthSystem {
	return &HealthSystem{world: world}
}

func (s *HealthSystem) Update(entities []Entity, components *ComponentStorage, dt float64) {
	for _, entity := range entities {
		health, hasHealth := components.Healths[entity]
		if !hasHealth {
			continue
		}
		if health.Ages {
			health.TTL = math.Max(0, health.TTL-dt)
		}
		if health.Ages && health.TTL == 0 {
			health.Dead = true
		}
		if health.Dead && health.Decays {
			health.DecayTTL = math.Max(0, health.DecayTTL-dt)
		}
		if health.Dead && health.Decays && health.DecayTTL == 0 {
			health.Decayed = true
		}

		if health.Dead && (!health.Decays || health.Decayed) {
			s.world.RemoveEntity(entity)
		}
	}
}

type MovementSystem struct{}

func NewMovementSystem() *MovementSystem {
	return &MovementSystem{}
}

func (s *MovementSystem) Update(entities []Entity, components *ComponentStorage, dt float64) {
	for _, entity := range entities {
		pos, hasPos := components.Positions[entity]
		acceleration, hasAcceleration := components.Accelerations[entity]
		friction, hasFriction := components.Frictions[entity]
		velocity, hasVelocity := components.Velocities[entity]

		if !hasPos || !hasVelocity {
			continue
		}

		if hasAcceleration {
			velocity.Current = math.Min(velocity.Max, velocity.Current+acceleration.Current*dt)
			velocity.AngularCurrent = math.Min(velocity.AngularMax, velocity.AngularCurrent+acceleration.AngularCurrent*dt)
		}

		if hasFriction {
			velocity.Current = math.Max(0, velocity.Current-friction.Current*dt)
			velocity.AngularCurrent = math.Max(0, velocity.AngularCurrent-friction.AngularCurrent*dt)
		}

		if hasPos && hasVelocity {
			pos.Direction += velocity.AngularCurrent
			pos.X += dt * velocity.Current * math.Cos(pos.Direction)
			pos.Y += dt * velocity.Current * math.Sin(pos.Direction)
		}
	}
}

type SpeedPowerUpSystem struct {
	behavior *bhv.Tree
}

func NewSpeedPowerUpSystem(world *World) *SpeedPowerUpSystem {
	return &SpeedPowerUpSystem{
		behavior: bhv.NewTree(
			bhv.SequenceNode(
				waitNode(&WaitState{TimeToWaitFn: frandomF(5, 7)},
					&bhv.Node{
						OnTick: func(n *bhv.Node, dt float64) bhv.Status {
							entity := world.AddEntity(SpeedPowerUp)
							world.Components.Positions[entity] = &Position{
								X:         frandom(70, world.width-70),
								Y:         frandom(70, world.height-70),
								Direction: -math.Pi / 2,
							}
							world.Components.BoundingBoxes[entity] = &BoundingBox{
								Width:  20,
								Height: 20,
							}
							world.Components.Healths[entity] = &Health{
								Ages: true,
								TTL:  30,
							}
							return bhv.StatusSuccess
						},
					},
				),
			),
		),
	}
}

func (s *SpeedPowerUpSystem) Update(entities []Entity, components *ComponentStorage, dt float64) {
	s.behavior.Tick(dt)
}
