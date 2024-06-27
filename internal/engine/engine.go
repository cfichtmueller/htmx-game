package engine

import (
	"sync"
	"time"

	"cfichtmueller.com/htmx-game/internal/engine/physics"
)

type Engine struct {
	mu          sync.Mutex
	World       *World
	Fps         int
	loopTicker  *time.Ticker
	playerIndex map[string]Entity
}

func New(width, height float64) *Engine {
	world := NewWorld(width, height)

	// beware - the order of systems is important

	world.AddSystem(NewAutoMoveSystem())
	world.AddSystem(NewMovementSystem())

	collisionDetection := NewCollisionDetectionSystem()
	collisionDetection.RegisterHandler(Bullet, Tank, NewBulletPlayerCollisionHandler(world))
	collisionDetection.RegisterHandler(Bullet, Player, NewBulletPlayerCollisionHandler(world))
	collisionDetection.RegisterHandler(Player, Tower, &PlayerTowerCollisionHandler{})
	collisionDetection.RegisterHandler(Tank, Tower, &TankTowerCollisionHandler{})
	collisionDetection.RegisterHandler(Player, SpeedPowerUp, NewPlayerPowerUpCollisionHandler(world, func(entity Entity, components *ComponentStorage) {
		components.Velocities[entity].Max += 5
	}))
	collisionDetection.RegisterHandler(Tank, Player, &TankPlayerCollisionHandler{})

	world.AddSystem(collisionDetection)
	world.AddSystem(NewHealthSystem(world))

	world.AddSystem(NewSensingSystem())
	world.AddSystem(NewSpeedPowerUpSystem(world))
	world.AddSystem(NewBehaviorSystem())

	return &Engine{
		World:       world,
		loopTicker:  time.NewTicker(30 * time.Millisecond),
		playerIndex: make(map[string]Entity),
	}
}

func (e *Engine) Lock() {
	e.mu.Lock()
}

func (e *Engine) Unlock() {
	e.mu.Unlock()
}

func (e *Engine) Start() {
	PlaceInRaster(1, int((e.World.Width-100)/50), int((e.World.Height-100)/50), func(x, y int) {
		SpawnTankShelter(
			e.World,
			float64(50+x*50),
			float64(50+y*50),
			physics.Deg90,
		)
	})
	PlaceInRaster(10, int((e.World.Width-100)/50), int((e.World.Height-100)/50), func(x, y int) {
		SpawnTower(
			e.World,
			float64(50+x*50),
			float64(50+y*50),
		)
	})

	go func(e *Engine) {
		last := time.Now().UnixMilli()
		for {
			<-e.loopTicker.C
			now := time.Now().UnixMilli()
			delta := float64(now-last) / 1000
			last = now
			e.mu.Lock()
			if (now - last) != 0 {
				e.Fps = int(1000 / (now - last))
			}
			e.World.Update(delta)
			e.mu.Unlock()
		}
	}(e)
}

func (e *Engine) SpawnPlayer() string {
	id := randomId()
	p := SpawnPlayer(
		e.World,
		e.World.Width/2,
		e.World.Height/2,
		physics.Deg0,
	)
	e.playerIndex[id] = p
	return id
}

func (e *Engine) PlayerWithId(id string) (Entity, bool) {
	entity, ok := e.playerIndex[id]
	return entity, ok
}
