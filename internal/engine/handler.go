package engine

type BulletTankCollisionHandler struct{}

func (h *BulletTankCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage) {
	components.RemoveEntity(entityA)
	components.Healths[entityB].Dead = true
	components.Velocities[entityB].Current = 0
}

type BulletPlayerCollisionHandler struct{}

func (h *BulletPlayerCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage) {
	components.RemoveEntity(entityA)
	components.Healths[entityB].Dead = true
}

type PlayerTowerCollisionHandler struct{}

func (h *PlayerTowerCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage) {
	playerHealth := components.Healths[entityA]
	if playerHealth.Dead {
		return
	}
	components.Healths[entityB].Dead = true
}

type PlayerPowerUpCollisionHandler struct {
	F func(entity Entity, components *ComponentStorage)
}

func (h *PlayerPowerUpCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage) {
	h.F(entityA, components)
	components.RemoveEntity(entityB)
}

type TankPlayerCollisionHandler struct{}

func (h *TankPlayerCollisionHandler) HandleCollision(entityA, entityB Entity, components *ComponentStorage) {
	tankHealth := components.Healths[entityA]
	if tankHealth.Dead {
		return
	}
	health := components.Healths[entityB]
	health.Dead = true
}
