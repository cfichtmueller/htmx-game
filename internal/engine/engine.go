package engine

import (
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
	for i := 0; i < 10; i++ {
		e.State.AddCell(NewTowerCell(
			frandom(0, e.State.Width),
			frandom(0, e.State.Height),
		))
	}
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
			x := frandom(0, 10)
			if x > 7 {
				s.AddCell(NewVelocityPowerUpCell(
					frandom(0, s.Width),
					frandom(0, s.Height),
				))
			} else {
				s.AddCell(NewBulletCell(0, 0, frandom(0, 1), frandom(10, 110), 15))
			}
		}
	}(e.spawnTicker, e.State)
}
