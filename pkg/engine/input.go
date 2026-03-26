// Package engine provides the core ECS framework, game loop, deterministic RNG,
// input handling, and camera system.
package engine

// InputReader defines how to read raw input. Implemented by EbitenInputReader.
type InputReader interface {
	ReadState(bindings KeyBindings) InputState
}

// KeyBindings holds the string key names for player controls.
type KeyBindings struct {
	Thrust      string
	RotateLeft  string
	RotateRight string
	Fire        string
	Secondary   string
	Pause       string
}

// DefaultKeyBindings returns the default key bindings.
func DefaultKeyBindings() KeyBindings {
	return KeyBindings{
		Thrust:      "W",
		RotateLeft:  "A",
		RotateRight: "D",
		Fire:        "Space",
		Secondary:   "Shift",
		Pause:       "Escape",
	}
}

// InputSystem reads player input and applies it to entities.
type InputSystem struct {
	world        *World
	physics      *PhysicsSystem
	bindings     KeyBindings
	reader       InputReader
	playerEntity Entity
	state        InputState
}

// NewInputSystem creates an input system attached to the world.
func NewInputSystem(world *World, physics *PhysicsSystem, bindings KeyBindings, reader InputReader) *InputSystem {
	return &InputSystem{
		world:    world,
		physics:  physics,
		bindings: bindings,
		reader:   reader,
	}
}

// SetPlayerEntity sets which entity receives input commands.
func (is *InputSystem) SetPlayerEntity(entity Entity) {
	is.playerEntity = entity
}

// GetState returns the current input state.
func (is *InputSystem) GetState() InputState {
	return is.state
}

// SetState sets the input state directly (for testing or replay).
func (is *InputSystem) SetState(state InputState) {
	is.state = state
}

// Update reads input and applies it to the player entity.
func (is *InputSystem) Update(dt float64) {
	if is.reader != nil {
		is.state = is.reader.ReadState(is.bindings)
	}
	is.applyToPlayer(dt)
}

// applyToPlayer applies the input state to the player entity.
func (is *InputSystem) applyToPlayer(dt float64) {
	if is.playerEntity == 0 || is.physics == nil {
		return
	}

	if is.state.Thrust {
		is.physics.ApplyThrust(is.playerEntity, dt)
	}

	if is.state.RotateLeft {
		is.physics.ApplyRotation(is.playerEntity, -1.0, dt)
	}

	if is.state.RotateRight {
		is.physics.ApplyRotation(is.playerEntity, 1.0, dt)
	}
}

// IsFirePressed returns true if the fire button is pressed.
func (is *InputSystem) IsFirePressed() bool {
	return is.state.Fire
}

// IsPausePressed returns true if the pause button is pressed.
func (is *InputSystem) IsPausePressed() bool {
	return is.state.Pause
}
