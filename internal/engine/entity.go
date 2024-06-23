package engine

type EntityType int

const (
	Player EntityType = iota
	Bullet
	Tank
	TankShelter
	Tower
	SpeedPowerUp
)

type Entity int64
