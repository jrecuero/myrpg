// Package ecs provides an Entity-Component-System (ECS) framework for game development.
// It defines core components such as Transform, SpriteComponent, and ColliderComponent
// that can be attached to entities to define their properties and behaviors.
// The ECS architecture allows for flexible and modular game object management.
// Each component struct includes fields relevant to its purpose, along with constructor functions
// for easy instantiation. The Transform component handles position and size, the SpriteComponent
// manages visual representation, and the ColliderComponent defines collision properties.
package ecs

import "sync/atomic"

type Entity struct {
	ID         int64                  // Unique identifier for the entity
	components map[string]interface{} // Components associated with the entity
}

var (
	maxEntityID int64 = 0 // Global counter for generating unique entity IDs
)

// NewEntity creates a new entity with a unique ID and initializes its components map.
// returns a pointer to the newly created Entity.
func NewEntity() *Entity {
	id := atomic.AddInt64(&maxEntityID, 1)
	return &Entity{
		ID:         id,
		components: make(map[string]interface{}),
	}
}

// AddComponent adds a component to the entity.
// name is the name of the component.
// component is the component data.
// returns nothing.
func (e *Entity) AddComponent(name string, component interface{}) {
	e.components[name] = component
}

// RemoveComponent removes a component from the entity.
// name is the name of the component to remove.
// returns nothing.
func (e *Entity) RemoveComponent(name string) {
	delete(e.components, name)
}

// GetComponent retrieves a component from the entity.
// name is the name of the component to retrieve.
// returns the component and a boolean indicating if it was found.
func (e *Entity) GetComponent(name string) (interface{}, bool) {
	component, exists := e.components[name]
	return component, exists
}

// HasComponent checks if the entity has a specific component.
// name is the name of the component to check.
// returns true if the component exists, false otherwise.
func (e *Entity) HasComponent(name string) bool {
	_, exists := e.components[name]
	return exists
}

// Transform retrieves the Transform component from the entity.
// returns a pointer to the Transform component or nil if not found.
func (e *Entity) Transform() *Transform {
	if comp, exists := e.GetComponent(ComponentTransform); exists {
		if transform, ok := comp.(*Transform); ok {
			return transform
		}
	}
	return nil
}

// Sprite retrieves the SpriteComponent from the entity.
// returns a pointer to the SpriteComponent or nil if not found.
func (e *Entity) Sprite() *SpriteComponent {
	if comp, exists := e.GetComponent(ComponentSprite); exists {
		if sprite, ok := comp.(*SpriteComponent); ok {
			return sprite
		}
	}
	return nil
}

// Collider retrieves the ColliderComponent from the entity.
// returns a pointer to the ColliderComponent or nil if not found.
func (e *Entity) Collider() *ColliderComponent {
	if comp, exists := e.GetComponent(ComponentCollider); exists {
		if collider, ok := comp.(*ColliderComponent); ok {
			return collider
		}
	}
	return nil
}
