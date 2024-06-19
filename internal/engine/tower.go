package engine

import (
	"math"
	"math/rand"
)

type TowerCellData struct {
	timeToNextBurst    float64
	timeToNextBullet   float64
	burstDirection     float64
	bulletsLeftInBurst int
}

func NewTowerCell(x, y float64) *Cell {
	return &Cell{
		Agent: &Agent{
			X:      x,
			Y:      y,
			Width:  30,
			Height: 30,
		},
		Type:  "tower",
		Color: "#000000",
		Data: &TowerCellData{
			timeToNextBurst: 5 * rand.Float64(),
		},
		HandleUpdate: func(c *Cell, dt float64) CellUpdateResult {
			d := c.Data.(*TowerCellData)
			d.timeToNextBurst -= dt
			d.timeToNextBullet -= dt
			if d.timeToNextBurst > 0 || d.timeToNextBullet > 0 {
				return CellUpdateResult{}
			}
			if d.bulletsLeftInBurst == 0 {
				d.timeToNextBurst = 5 + 10*rand.Float64()
				d.bulletsLeftInBurst = 3 + int(rand.Float64()*3)
				d.burstDirection = math.Pi * 2 * rand.Float64()
				c.Agent.Direction = d.burstDirection
				return CellUpdateResult{}
			}
			d.bulletsLeftInBurst -= 1
			d.timeToNextBullet = 0.3
			return CellUpdateResult{
				Cells: generateBullets(
					1,
					c.Agent.X,
					c.Agent.Y,
					d.burstDirection,
					0.02*rand.Float64(),
					70,
					10,
				),
			}
		},
		HandlePlayerCollision: func(c *Cell, p *Player) {
			c.Die()
		},
	}
}

func generateBullets(count int, x, y, direction, spread, velocity, ttl float64) []*Cell {
	res := make([]*Cell, count)
	for i := 0; i < count; i++ {
		res[i] = NewBulletCell(x, y, direction+spread*(rand.Float64()-0.5), velocity, ttl)
	}
	return res
}
