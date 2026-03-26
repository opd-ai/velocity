//go:build !noebiten

// Package engine provides the core ECS framework, game loop, deterministic RNG,
// input handling, and camera system.
package engine

import "github.com/hajimehoshi/ebiten/v2"

// EbitenInputReader reads input from Ebiten's keyboard API.
type EbitenInputReader struct{}

// NewEbitenInputReader creates a new Ebiten-based input reader.
func NewEbitenInputReader() *EbitenInputReader {
	return &EbitenInputReader{}
}

// ReadState reads the current keyboard state based on the given bindings.
func (r *EbitenInputReader) ReadState(bindings KeyBindings) InputState {
	return InputState{
		Thrust:      ebiten.IsKeyPressed(parseKey(bindings.Thrust)),
		RotateLeft:  ebiten.IsKeyPressed(parseKey(bindings.RotateLeft)),
		RotateRight: ebiten.IsKeyPressed(parseKey(bindings.RotateRight)),
		Fire:        ebiten.IsKeyPressed(parseKey(bindings.Fire)),
		Secondary:   ebiten.IsKeyPressed(parseKey(bindings.Secondary)),
		Pause:       ebiten.IsKeyPressed(parseKey(bindings.Pause)),
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
