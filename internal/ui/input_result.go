package ui

// InputResult represents the result of input handling by UI widgets
type InputResult struct {
	EscConsumed   bool // Whether the ESC key was consumed
	MouseConsumed bool // Whether mouse input was consumed
}

// NewInputResult creates a new InputResult with default values
func NewInputResult() InputResult {
	return InputResult{
		EscConsumed:   false,
		MouseConsumed: false,
	}
}

// Combine merges this InputResult with another, using OR logic
func (ir *InputResult) Combine(other InputResult) {
	ir.EscConsumed = ir.EscConsumed || other.EscConsumed
	ir.MouseConsumed = ir.MouseConsumed || other.MouseConsumed
}

// HasAnyConsumption returns true if any input was consumed
func (ir *InputResult) HasAnyConsumption() bool {
	return ir.EscConsumed || ir.MouseConsumed
}
