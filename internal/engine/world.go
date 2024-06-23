package engine

type World struct {
	nextEntity       int64
	Entities         []Entity
	entitiesToRemove map[Entity]bool
	Components       *ComponentStorage
	systems          []System
	width            float64
	height           float64
}

func NewWorld(width, height float64) *World {
	return &World{
		Entities:         make([]Entity, 0),
		entitiesToRemove: make(map[Entity]bool),
		Components:       NewComponentStorage(),
		systems:          make([]System, 0),
		width:            width,
		height:           height,
	}
}

func (w *World) AddEntity(entityType EntityType) Entity {
	entity := Entity(w.nextEntity)
	w.nextEntity++
	w.Entities = append(w.Entities, entity)
	w.Components.EntityTypes[entity] = &EntityTypeComponent{
		Type: entityType,
	}
	return entity
}

func (w *World) RemoveEntity(entity Entity) {
	w.entitiesToRemove[entity] = true
}

func (w *World) AddSystem(sytem System) {
	w.systems = append(w.systems, sytem)
}

func (w *World) Update(dt float64) {
	for _, system := range w.systems {
		system.Update(w.Entities, w.Components, dt)
	}
	w.cleanupEntities()
}

func (w *World) cleanupEntities() {
	for entity := range w.entitiesToRemove {
		w.Components.RemoveEntity(entity)
		for i, e := range w.Entities {
			if e == entity {
				w.Entities = append(w.Entities[:i], w.Entities[i+1:]...)
				break
			}
		}
	}
	w.entitiesToRemove = make(map[Entity]bool)
}

func (w *World) SetVelocity(entity Entity, v float64) {
	w.Components.Velocities[entity].Current = v
}
