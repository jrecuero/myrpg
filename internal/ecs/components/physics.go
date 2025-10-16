package components

// ColliderComponent represents the collision properties of an entity.
// It defines the size and behavior of the collider used in collision detection.
// Solid indicates if the collider is solid (affects collision response).
// Width and Height define the size of the collider.
// OffsetX and OffsetY are offsets applied to the collider's position relative to the entity's position.
type ColliderComponent struct {
	Solid   bool // Indicates if the collider is solid (affects collision response)
	Width   int  // Width of the collider
	Height  int  // Height of the collider
	OffsetX int  // X offset of the collider relative to the entity's position
	OffsetY int  // Y offset of the collider relative to the entity's position
}

// NewColliderComponent creates a new ColliderComponent with the specified properties.
// solid indicates if the collider is solid (affects collision response).
// width and height define the size of the collider.
// offsetX and offsetY are offsets applied to the collider's position relative to the entity's position.
// returns a pointer to the newly created ColliderComponent.
func NewColliderComponent(solid bool, width, height, offsetX, offsetY int) *ColliderComponent {
	return &ColliderComponent{
		Solid:   solid,
		Width:   width,
		Height:  height,
		OffsetX: offsetX,
		OffsetY: offsetY,
	}
}
