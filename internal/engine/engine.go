package engine

import (
	"sync"
	"time"
)

type Engine struct {
	mu         sync.Mutex
	loopTicker *time.Ticker
	State      *State
}

func New(width, height float64) *Engine {
	return &Engine{
		loopTicker: time.NewTicker(30 * time.Millisecond),
		State:      NewState(width, height),
	}
}

func (e *Engine) Lock() {
	e.mu.Lock()
}

func (e *Engine) Unlock() {
	e.mu.Unlock()
}

func (e *Engine) Start() {
	for i := 0; i < 10; i++ {
		SpawnTower(
			e.State.World,
			frandom(50, e.State.Width-50),
			frandom(50, e.State.Height-50),
		)
	}
	go func(e *Engine) {
		last := time.Now().UnixMilli()
		for {
			<-e.loopTicker.C
			now := time.Now().UnixMilli()
			delta := float64(now-last) / 1000
			last = now
			e.mu.Lock()
			e.State.Update(delta)
			e.mu.Unlock()
		}
	}(e)
}

func (e *Engine) World() *World {
	return e.State.World
}
