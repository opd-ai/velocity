//go:build noebiten

// Package engine provides the core ECS framework, game loop, deterministic RNG,
// input handling, and camera system.
package engine

// StubInputReader is a no-op input reader for headless/test environments.
type StubInputReader struct{}

// NewEbitenInputReader returns a stub input reader when built without ebiten.
func NewEbitenInputReader() *StubInputReader {
	return &StubInputReader{}
}

// SetScreenSize is a no-op for the stub reader.
func (r *StubInputReader) SetScreenSize(width, height int) {}

// ReadState returns an empty input state (no keys pressed).
func (r *StubInputReader) ReadState(bindings KeyBindings) InputState {
	return InputState{}
}
