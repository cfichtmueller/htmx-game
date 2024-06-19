package engine

import (
	"math/rand"
	"time"
)

type Engine struct {
	loopTicker  *time.Ticker
	spawnTicker *time.Ticker
	State       *State
}

func New(width, height float64) *Engine {
	return &Engine{
		loopTicker:  time.NewTicker(30 * time.Millisecond),
		spawnTicker: time.NewTicker(100 * time.Millisecond),
		State:       NewState(width, height),
	}
}

func (e *Engine) Start() {
	go func(t *time.Ticker, s *State) {
		last := time.Now().UnixMilli()
		for {
			<-t.C
			now := time.Now().UnixMilli()
			delta := float64(now-last) / 1000
			last = now
			s.Update(delta)
		}
	}(e.loopTicker, e.State)

	go func(t *time.Ticker, s *State) {
		for {
			<-t.C
			if len(s.Cells) > 100 {
				continue
			}
			x := rand.Float64() * 10
			if x > 7 {
				s.AddCell(NewVelocityPowerUpCell(
					s.Width*rand.Float64(),
					s.Height*rand.Float64(),
				))
			} else {
				s.AddCell(NewBulletCell())
			}
		}
	}(e.spawnTicker, e.State)
}
