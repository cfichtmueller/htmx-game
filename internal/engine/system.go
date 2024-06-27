package engine

import (
	"math"

	"cfichtmueller.com/htmx-game/internal/engine/bhv"
	"cfichtmueller.com/htmx-game/internal/engine/physics"
)

type System interface {
	Update(entities []Entity, components *ComponentStorage, dt float64)
}

type AutoMoveSystem struct{}

func NewAutoMoveSystem() *AutoMoveSystem {
	return &AutoMoveSystem{}
}

func (s *AutoMoveSystem) Update(entities []Entity, components *ComponentStorage, dt float64) {
	for _, entity := range entities {
		autoMove, hasAutoMove := components.AutoMove[entity]
		if !hasAutoMove || !autoMove.TargetDirectionActive {
			continue
		}
		position, hasPosition := components.Positions[entity]
		acceleration, hasAcceleration := components.Accelerations[entity]
		velocity, hasVelocity := components.Velocities[entity]
		if !hasPosition || !hasVelocity {
			continue
		}

		delta := physics.ShortesRotationDirection(position.Direction, autoMove.TargetDirection)
		if math.Abs(delta) > physics.Deg1*5 {
			if hasAcceleration {
				acceleration.AngularCurrent = math.Copysign(acceleration.AngularMax, delta)
			} else {
				velocity.AngularCurrent = math.Copysign(velocity.AngularMax, delta)
			}
		} else {
			if hasAcceleration {
				acceleration.AngularCurrent = 0
			}
			velocity.AngularCurrent = 0
			position.Direction = autoMove.TargetDirection
			autoMove.TargetDirectionActive = false
		}
	}
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
	HandleCollision(entityA, entityB Entity, components *ComponentStorage, dt float64)
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
				physics.Collides(posA, bbA, posB, bbB) {
				s.collisions = append(s.collisions, Collision{EntityA: entityA, EntityB: entityB})
				s.handleCollision(entityA, entityB, components, dt)
			}
		}
	}
}

func (s *CollisionDetectionSystem) handleCollision(entityA, entityB Entity, components *ComponentStorage, dt float64) {
	typeAComp, hasTypeA := components.EntityTypes[entityA]
	typeBComp, hasTypeB := components.EntityTypes[entityB]

	if !hasTypeA || !hasTypeB {
		return
	}

	typeA := typeAComp.Type
	typeB := typeBComp.Type

	if handler, ok := s.handlers[typeA][typeB]; ok {
		handler.HandleCollision(entityA, entityB, components, dt)
	} else if handler, ok := s.handlers[typeB][typeA]; ok {
		handler.HandleCollision(entityB, entityA, components, dt)
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
		health, hasHealth := components.Healths[entity]

		if !hasPos || !hasVelocity {
			continue
		}

		if hasHealth && health.Dead {
			if hasAcceleration {
				acceleration.Current = 0
				acceleration.AngularCurrent = 0
			}
			velocity.Current = 0
			velocity.AngularCurrent = 0
		}

		if hasAcceleration {
			velocity.Current = physics.Accelerate(velocity.Current, velocity.Max, acceleration.Current, dt)
			velocity.AngularCurrent = physics.Accelerate(velocity.AngularCurrent, velocity.AngularMax, acceleration.AngularCurrent, dt)
		}

		if hasFriction {
			velocity.Current = math.Max(0, velocity.Current-friction.Current*dt)
			velocity.AngularCurrent = math.Max(0, velocity.AngularCurrent-friction.AngularCurrent*dt)
		}

		if hasPos && hasVelocity {
			pos.Direction += velocity.AngularCurrent * dt
			pos.X, pos.Y = physics.Move(pos.X, pos.Y, pos.Direction, velocity.Current, dt)
		}
	}
}

type SensingSystem struct{}

func NewSensingSystem() *SensingSystem {
	return &SensingSystem{}
}

func (s *SensingSystem) Update(entities []Entity, components *ComponentStorage, dt float64) {
	for _, entity := range entities {
		pos, hasPos := components.Positions[entity]
		sensing, hasSensing := components.Sensings[entity]

		if !hasPos || !hasSensing {
			continue
		}

		sensing.SensedEntities = []SensedEntity{}

		for _, otherEntity := range entities {
			if entity == otherEntity {
				continue
			}

			otherPos, hasOtherPos := components.Positions[otherEntity]
			otherType, hasOtherType := components.EntityTypes[otherEntity]

			if !hasOtherPos || !hasOtherType {
				continue
			}

			sensingRange, hasSensingRange := sensing.Ranges[otherType.Type]
			if !hasSensingRange {
				continue
			}

			distance := physics.Distance(pos, otherPos)
			if distance <= sensingRange {
				sensing.SensedEntities = append(sensing.SensedEntities, SensedEntity{
					Entity:   otherEntity,
					Type:     otherType.Type,
					Position: otherPos,
				})
			}
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
				bhv.WaitNode(&bhv.WaitState{TimeToWaitFn: frandomF(5, 7)},
					&bhv.Node{
						OnTick: func(n *bhv.Node, dt float64) bhv.Status {
							entity := world.AddEntity(SpeedPowerUp)
							world.Components.Positions[entity] = &physics.Position{
								X:         frandom(70, world.Width-70),
								Y:         frandom(70, world.Height-70),
								Direction: -physics.Deg90,
							}
							world.Components.BoundingBoxes[entity] = &physics.Rectangle{W: 20, H: 20}
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
