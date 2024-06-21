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

type CellList struct {
	Cells []*Cell
}

func NewCellList() *CellList {
	return &CellList{
		Cells: make([]*Cell, 0),
	}
}

func NewCellListWithCells(cells ...*Cell) *CellList {
	return &CellList{Cells: cells}
}

func (l *CellList) Add(cells ...*Cell) {
	l.Cells = append(l.Cells, cells...)
}

func (l *CellList) Each(f func(c *Cell)) {
	for _, c := range l.Cells {
		f(c)
	}
}

func (l *CellList) Filter(f func(c *Cell) bool) {
	newCells := make([]*Cell, 0, len(l.Cells))
	for _, c := range l.Cells {
		if f(c) {
			newCells = append(newCells, c)
		}
	}
	l.Cells = newCells
}

func (l *CellList) Length() int {
	return len(l.Cells)
}
