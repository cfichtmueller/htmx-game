package engine

type CellUpdateResult struct {
	Cells []*Cell
}

type Cell struct {
	Agent                 *Agent
	Color                 string
	Type                  string
	Data                  any
	HandleUpdate          func(c *Cell, dt float64) CellUpdateResult
	HandlePlayerCollision func(c *Cell, p *Player)
}

func (c *Cell) Update(dt float64) CellUpdateResult {
	if c.Agent.Dead {
		return CellUpdateResult{}
	}
	c.Agent.Update(dt)
	if c.HandleUpdate == nil {
		return CellUpdateResult{}
	}
	return c.HandleUpdate(c, dt)
}

func (c *Cell) Die() {
	c.Agent.Dead = true
	c.Agent.Stop()
}
