package ecs

// World represents the game world containing all entities.
// It manages the collection of entities and provides methods to add, remove, and retrieve entities.
// entities is a slice of pointers to Entity, representing all entities in the world.
type World struct {
	entities []*Entity // All entities in the world
}

// NewWorld creates a new, empty game world.
// returns a pointer to the newly created World.
func NewWorld() *World {
	return &World{
		entities: make([]*Entity, 0),
	}
}

// AddEntity adds a new entity to the world.
// entity is a pointer to the Entity to be added.
// returns nothing.
func (w *World) AddEntity(entity *Entity) {
	w.entities = append(w.entities, entity)
}

// RemoveEntity removes an entity from the world.
// entity is a pointer to the Entity to be removed.
// returns nothing.
func (w *World) RemoveEntity(entity *Entity) {
	for i, e := range w.entities {
		if e == entity {
			w.entities = append(w.entities[:i], w.entities[i+1:]...)
			break
		}
	}
}

// RemoveByName removes an entity from the world by its name.
// name is the name of the entity to be removed.
// returns true if an entity was found and removed, false otherwise.
func (w *World) RemoveByName(name string) bool {
	for i, e := range w.entities {
		if e.Name == name {
			w.entities = append(w.entities[:i], w.entities[i+1:]...)
			return true
		}
	}
	return false
}

// GetEntities returns all entities in the world.
// returns a slice of pointers to Entity representing all entities in the world.
func (w *World) GetEntities() []*Entity {
	return w.entities
}

// FindWithComponent finds all entities that have a specific component.
// name is the name of the component to search for.
// returns a slice of pointers to Entity that have the specified component.
func (w *World) FindWithComponent(name string) []*Entity {
	var result []*Entity
	for _, e := range w.entities {
		if e.HasComponent(name) {
			result = append(result, e)
		}
	}
	return result
}

// FindByID finds an entity by its unique ID.
// id is the unique identifier of the entity to search for.
// returns a pointer to the Entity if found, or nil if not found.
func (w *World) FindByID(id int64) *Entity {
	for _, e := range w.entities {
		if e.ID == id {
			return e
		}
	}
	return nil
}

// FindByName finds an entity by its name.
// name is the name of the entity to search for.
// returns a pointer to the Entity if found, or nil if not found.
func (w *World) FindByName(name string) *Entity {
	for _, e := range w.entities {
		if e.Name == name {
			return e
		}
	}
	return nil
}

// Clear removes all entities from the world.
// returns nothing.
func (w *World) Clear() {
	w.entities = make([]*Entity, 0)
}
