// Package ecs provides an Entity-Component-System (ECS) framework for game development.
// It defines core components such as Transform, SpriteComponent, and ColliderComponent
// that can be attached to entities to define their properties and behaviors.
// The ECS architecture allows for flexible and modular game object management.
// Each component struct includes fields relevant to its purpose, along with constructor functions
// for easy instantiation. The Transform component handles position and size, the SpriteComponent
// manages visual representation, and the ColliderComponent defines collision properties.
package ecs

import (
	"fmt"
	"sync/atomic"

	"github.com/jrecuero/myrpg/internal/ecs/components"
)

type Entity struct {
	ID         int64                  // Unique identifier for the entity
	Name       string                 // Name of the entity
	components map[string]interface{} // Components associated with the entity
	tags       map[string]bool        // Tags associated with the entity
}

var (
	maxEntityID int64 = 0 // Global counter for generating unique entity IDs
)

// NewEntity creates a new entity with a unique ID and name, and initializes its components map.
// name is the name to assign to the entity.
// returns a pointer to the newly created Entity.
func NewEntity(name string) *Entity {
	id := atomic.AddInt64(&maxEntityID, 1)
	return &Entity{
		ID:         id,
		Name:       name,
		components: make(map[string]interface{}),
		tags:       make(map[string]bool),
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
func (e *Entity) Transform() *components.Transform {
	if comp, exists := e.GetComponent(ComponentTransform); exists {
		if transform, ok := comp.(*components.Transform); ok {
			return transform
		}
	}
	return nil
}

// Sprite retrieves the SpriteComponent from the entity.
// returns a pointer to the SpriteComponent or nil if not found.
func (e *Entity) Sprite() *components.SpriteComponent {
	if comp, exists := e.GetComponent(ComponentSprite); exists {
		if sprite, ok := comp.(*components.SpriteComponent); ok {
			return sprite
		}
	}
	return nil
}

// Collider retrieves the ColliderComponent from the entity.
// returns a pointer to the ColliderComponent or nil if not found.
func (e *Entity) Collider() *components.ColliderComponent {
	if comp, exists := e.GetComponent(ComponentCollider); exists {
		if collider, ok := comp.(*components.ColliderComponent); ok {
			return collider
		}
	}
	return nil
}

// RPGStats retrieves the RPGStatsComponent from the entity.
// returns a pointer to the RPGStatsComponent or nil if not found.
func (e *Entity) RPGStats() *components.RPGStatsComponent {
	if comp, exists := e.GetComponent(ComponentRPGStats); exists {
		if stats, ok := comp.(*components.RPGStatsComponent); ok {
			return stats
		}
	}
	return nil
}

// Animation retrieves the AnimationComponent from the entity.
// returns a pointer to the AnimationComponent or nil if not found.
func (e *Entity) Animation() *components.AnimationComponent {
	if comp, exists := e.GetComponent(ComponentAnimation); exists {
		if animation, ok := comp.(*components.AnimationComponent); ok {
			return animation
		}
	}
	return nil
}

// AddTag adds a tag to the entity.
// tag is the tag string to add.
func (e *Entity) AddTag(tag string) {
	e.tags[tag] = true
}

// RemoveTag removes a tag from the entity.
// tag is the tag string to remove.
func (e *Entity) RemoveTag(tag string) {
	delete(e.tags, tag)
}

// HasTag checks if the entity has a specific tag.
// tag is the tag string to check.
// returns true if the entity has the tag, false otherwise.
func (e *Entity) HasTag(tag string) bool {
	_, exists := e.tags[tag]
	return exists
}

// GetTags returns all tags associated with the entity.
// returns a slice of tag strings.
func (e *Entity) GetTags() []string {
	tags := make([]string, 0, len(e.tags))
	for tag := range e.tags {
		tags = append(tags, tag)
	}
	return tags
}

// GetID returns the unique ID of the entity as a string
func (e *Entity) GetID() string {
	return fmt.Sprintf("%d", e.ID)
}
