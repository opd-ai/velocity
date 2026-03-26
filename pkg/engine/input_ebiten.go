//go:build !noebiten

// Package engine provides the core ECS framework, game loop, deterministic RNG,
// input handling, and camera system.
package engine

import "github.com/hajimehoshi/ebiten/v2"

// EbitenInputReader reads input from Ebiten's keyboard and gamepad APIs.
type EbitenInputReader struct {
	// deadzone threshold for analog sticks
	deadzone float64
	// Screen dimensions for touch region calculation
	screenWidth  int
	screenHeight int
}

// NewEbitenInputReader creates a new Ebiten-based input reader.
func NewEbitenInputReader() *EbitenInputReader {
	return &EbitenInputReader{
		deadzone:     0.25,
		screenWidth:  800, // Default, can be updated
		screenHeight: 600,
	}
}

// SetScreenSize updates the screen dimensions for touch region calculation.
func (r *EbitenInputReader) SetScreenSize(width, height int) {
	r.screenWidth = width
	r.screenHeight = height
}

// ReadState reads the current input state from keyboard, gamepads, and touch.
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

	// Merge touch input
	r.mergeTouchInput(&state)

	return state
}

// mergeTouchInput reads from touch screen and maps regions to virtual buttons.
// Touch regions:
//   - Left 1/3 of screen: rotate left
//   - Right 1/3 of screen: rotate right
//   - Bottom center: thrust
//   - Top center: fire
func (r *EbitenInputReader) mergeTouchInput(state *InputState) {
	touchIDs := ebiten.AppendTouchIDs(nil)
	for _, id := range touchIDs {
		x, y := ebiten.TouchPosition(id)
		r.mapTouchToRotation(x, state)
		r.mapTouchToAction(x, y, state)
	}
}

// mapTouchToRotation sets rotation state based on horizontal touch position.
func (r *EbitenInputReader) mapTouchToRotation(x int, state *InputState) {
	leftThird := r.screenWidth / 3
	rightThird := r.screenWidth * 2 / 3

	if x < leftThird {
		state.RotateLeft = true
	} else if x > rightThird {
		state.RotateRight = true
	}
}

// mapTouchToAction sets thrust or fire based on vertical touch in center column.
func (r *EbitenInputReader) mapTouchToAction(x, y int, state *InputState) {
	leftThird := r.screenWidth / 3
	rightThird := r.screenWidth * 2 / 3
	topHalf := r.screenHeight / 2

	if x >= leftThird && x <= rightThird {
		if y > topHalf {
			state.Thrust = true
		} else {
			state.Fire = true
		}
	}
}

// mergeGamepadInput reads from all connected gamepads and merges into state.
func (r *EbitenInputReader) mergeGamepadInput(state *InputState) {
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {
		r.readGamepadRotation(id, state)
		r.readGamepadThrust(id, state)
		r.readGamepadButtons(id, state)
	}
}

// readGamepadRotation reads rotation input from D-pad and left stick.
func (r *EbitenInputReader) readGamepadRotation(id ebiten.GamepadID, state *InputState) {
	leftX := ebiten.GamepadAxisValue(id, 0) // Left stick X
	if leftX < -r.deadzone || ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton13) {
		state.RotateLeft = true
	}
	if leftX > r.deadzone || ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton14) {
		state.RotateRight = true
	}
}

// readGamepadThrust reads thrust input from D-pad, left stick, and triggers.
func (r *EbitenInputReader) readGamepadThrust(id ebiten.GamepadID, state *InputState) {
	leftY := ebiten.GamepadAxisValue(id, 1) // Left stick Y
	if leftY < -r.deadzone || ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton11) {
		state.Thrust = true
	}
	rightTrigger := ebiten.GamepadAxisValue(id, 5)
	if rightTrigger > r.deadzone {
		state.Thrust = true
	}
}

// readGamepadButtons reads face button input for fire, secondary, and pause.
func (r *EbitenInputReader) readGamepadButtons(id ebiten.GamepadID, state *InputState) {
	if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton0) {
		state.Fire = true
	}
	if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton1) {
		state.Secondary = true
	}
	if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton7) ||
		ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton9) {
		state.Pause = true
	}
}

// keyNameMap maps string key names to ebiten.Key values.
var keyNameMap = map[string]ebiten.Key{
	"W":      ebiten.KeyW,
	"A":      ebiten.KeyA,
	"S":      ebiten.KeyS,
	"D":      ebiten.KeyD,
	"Space":  ebiten.KeySpace,
	"Shift":  ebiten.KeyShiftLeft,
	"Escape": ebiten.KeyEscape,
	"Enter":  ebiten.KeyEnter,
	"Up":     ebiten.KeyUp,
	"Down":   ebiten.KeyDown,
	"Left":   ebiten.KeyLeft,
	"Right":  ebiten.KeyRight,
}

// parseKey converts a string key name to an ebiten.Key.
func parseKey(name string) ebiten.Key {
	if key, ok := keyNameMap[name]; ok {
		return key
	}
	return ebiten.KeyW
}
