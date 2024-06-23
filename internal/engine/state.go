package engine

import (
	"math"
	"sync"
)

type State struct {
	mu          sync.Mutex
	World       *World
	Width       float64
	Height      float64
	playerIndex map[string]Entity
}

func NewState(width, height float64) *State {
	world := NewWorld(width, height)

	world.AddSystem(NewBehaviorSystem())

	world.AddSystem(NewMovementSystem())

	collisionDetection := NewCollisionDetectionSystem()
	collisionDetection.RegisterHandler(Bullet, Tank, &BulletTankCollisionHandler{})
	collisionDetection.RegisterHandler(Bullet, Player, &BulletPlayerCollisionHandler{})
	collisionDetection.RegisterHandler(Player, Tower, &PlayerTowerCollisionHandler{})
	collisionDetection.RegisterHandler(Player, SpeedPowerUp, &PlayerPowerUpCollisionHandler{F: func(entity Entity, components *ComponentStorage) {
		components.Velocities[entity].Max += 5
	}})
	collisionDetection.RegisterHandler(Tank, Player, &TankPlayerCollisionHandler{})

	world.AddSystem(collisionDetection)
	world.AddSystem(NewHealthSystem(world))

	world.AddSystem(NewSpeedPowerUpSystem(world))

	SpawnTankShelter(world, frandom(0, world.width), frandom(0, world.height), math.Pi/2)

	return &State{
		World:       world,
		Width:       width,
		Height:      height,
		playerIndex: make(map[string]Entity),
	}
}

func (s *State) Update(dt float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.World.Update(dt)
}

func (s *State) SpawnPlayer() string {
	s.mu.Lock()
	id := randomId()
	e := SpawnPlayer(
		s.World,
		s.Width/2,
		s.Height/2,
		math.Pi,
	)
	s.playerIndex[id] = e
	s.mu.Unlock()
	return id
}

func (s *State) PlayerWithId(id string) (Entity, bool) {
	entity, ok := s.playerIndex[id]
	return entity, ok
}
