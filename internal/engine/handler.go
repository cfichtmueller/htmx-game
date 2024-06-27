package engine

import (
	"cfichtmueller.com/htmx-game/internal/engine/physics"
)

type BulletTankCollisionHandler struct {
	world *World
}

func NewBulletTankCollisionHandler(world *World) *BulletPlayerCollisionHandler {
	return &BulletPlayerCollisionHandler{world: world}
}

func (h *BulletTankCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage, dt float64) {
	h.world.RemoveEntity(entityA)
	components.Healths[entityB].Dead = true
	components.Velocities[entityB].Current = 0
}

type BulletPlayerCollisionHandler struct {
	world *World
}

func NewBulletPlayerCollisionHandler(world *World) *BulletPlayerCollisionHandler {
	return &BulletPlayerCollisionHandler{world: world}
}

func (h *BulletPlayerCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage, dt float64) {
	h.world.RemoveEntity(entityA)
	components.Healths[entityB].Dead = true
}

type PlayerTowerCollisionHandler struct{}

func (h *PlayerTowerCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage, dt float64) {
	playerHealth := components.Healths[entityA]
	if playerHealth.Dead {
		return
	}
	components.Healths[entityB].Dead = true
}

type PlayerPowerUpCollisionHandler struct {
	world *World
	f     func(entity Entity, components *ComponentStorage)
}

func NewPlayerPowerUpCollisionHandler(world *World, f func(entity Entity, components *ComponentStorage)) *PlayerPowerUpCollisionHandler {
	return &PlayerPowerUpCollisionHandler{
		world: world,
		f:     f,
	}
}

func (h *PlayerPowerUpCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage, dt float64) {
	h.f(entityA, components)
	h.world.RemoveEntity(entityB)
}

type TankPlayerCollisionHandler struct{}

func (h *TankPlayerCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage, dt float64) {
	tankHealth := components.Healths[entityA]
	if tankHealth.Dead {
		return
	}
	health := components.Healths[entityB]
	health.Dead = true
}

type TankTowerCollisionHandler struct{}

func (h *TankTowerCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage, dt float64) {
	position := components.Positions[entityA]
	velocity := components.Velocities[entityA]

	physics.Move2(position, -velocity.Current, dt)
	position.Direction += physics.Deg180
}
