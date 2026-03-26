//go:build !noebiten

// Package engine provides the core ECS framework, game loop, deterministic RNG,
// input handling, and camera system.
package engine

import "github.com/hajimehoshi/ebiten/v2"

// EbitenInputReader reads input from Ebiten's keyboard and gamepad APIs.
type EbitenInputReader struct {
	// deadzone threshold for analog sticks
	deadzone float64
}

// NewEbitenInputReader creates a new Ebiten-based input reader.
func NewEbitenInputReader() *EbitenInputReader {
	return &EbitenInputReader{
		deadzone: 0.25,
	}
}

// ReadState reads the current input state from keyboard and gamepads.
func (r *EbitenInputReader) ReadState(bindings KeyBindings) InputState {
	state := InputState{
		Thrust:      ebiten.IsKeyPressed(parseKey(bindings.Thrust)),
		RotateLeft:  ebiten.IsKeyPressed(parseKey(bindings.RotateLeft)),
		RotateRight: ebiten.IsKeyPressed(parseKey(bindings.RotateRight)),
		Fire:        ebiten.IsKeyPressed(parseKey(bindings.Fire)),
		Secondary:   ebiten.IsKeyPressed(parseKey(bindings.Secondary)),
		Pause:       ebiten.IsKeyPressed(parseKey(bindings.Pause)),
	}

	// Merge gamepad input (any connected gamepad)
	r.mergeGamepadInput(&state)

	return state
}

// mergeGamepadInput reads from all connected gamepads and merges into state.
func (r *EbitenInputReader) mergeGamepadInput(state *InputState) {
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {
		// D-pad or left stick for rotation
		leftX := ebiten.GamepadAxisValue(id, 0) // Left stick X
		if leftX < -r.deadzone || ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton13) { // Left D-pad
			state.RotateLeft = true
		}
		if leftX > r.deadzone || ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton14) { // Right D-pad
			state.RotateRight = true
		}

		// Left stick Y or D-pad up for thrust
		leftY := ebiten.GamepadAxisValue(id, 1) // Left stick Y
		if leftY < -r.deadzone || ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton11) { // Up D-pad
			state.Thrust = true
		}

		// Face buttons
		// A/Cross (button 0) - Fire
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton0) {
			state.Fire = true
		}
		// B/Circle (button 1) - Secondary
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton1) {
			state.Secondary = true
		}
		// Start (button 9) or + (button 7) - Pause
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton7) ||
			ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton9) {
			state.Pause = true
		}

		// Right trigger (axis 5) for thrust as alternative
		rightTrigger := ebiten.GamepadAxisValue(id, 5)
		if rightTrigger > r.deadzone {
			state.Thrust = true
		}
	}
}

// parseKey converts a string key name to an ebiten.Key.
func parseKey(name string) ebiten.Key {
	switch name {
	case "W":
		return ebiten.KeyW
	case "A":
		return ebiten.KeyA
	case "S":
		return ebiten.KeyS
	case "D":
		return ebiten.KeyD
	case "Space":
		return ebiten.KeySpace
	case "Shift":
		return ebiten.KeyShiftLeft
	case "Escape":
		return ebiten.KeyEscape
	case "Enter":
		return ebiten.KeyEnter
	case "Up":
		return ebiten.KeyUp
	case "Down":
		return ebiten.KeyDown
	case "Left":
		return ebiten.KeyLeft
	case "Right":
		return ebiten.KeyRight
	default:
		return ebiten.KeyW
	}
}
