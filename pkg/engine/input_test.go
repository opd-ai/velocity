package engine

import (
	"testing"
)

// mockInputReader is a test double for InputReader.
type mockInputReader struct {
	state InputState
}

func (m *mockInputReader) ReadState(_ KeyBindings) InputState {
	return m.state
}

func TestDefaultKeyBindings(t *testing.T) {
	kb := DefaultKeyBindings()

	if kb.Thrust == "" {
		t.Error("expected Thrust binding to be set")
	}
	if kb.RotateLeft == "" {
		t.Error("expected RotateLeft binding to be set")
	}
	if kb.RotateRight == "" {
		t.Error("expected RotateRight binding to be set")
	}
	if kb.Fire == "" {
		t.Error("expected Fire binding to be set")
	}
}

func TestNewInputSystem(t *testing.T) {
	world := NewWorld()
	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()

	is := NewInputSystem(world, physics, bindings, nil)

	if is == nil {
		t.Fatal("expected non-nil InputSystem")
	}
}

func TestInputSystem_SetPlayerEntity(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 0, Y: 0})
	world.AddComponent(entity, "velocity", &Velocity{VX: 0, VY: 0})
	world.AddComponent(entity, "rotation", &Rotation{Angle: 0})

	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()
	is := NewInputSystem(world, physics, bindings, nil)

	is.SetPlayerEntity(entity)

	// Should not panic on Update even though no keys pressed
	is.Update(1.0 / 60.0)
}

func TestInputSystem_GetState(t *testing.T) {
	world := NewWorld()
	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()
	is := NewInputSystem(world, physics, bindings, nil)

	state := is.GetState()

	// Initial state should be all false
	if state.Thrust || state.Fire || state.Pause {
		t.Error("expected initial state to be all false")
	}
}

func TestInputSystem_SetState(t *testing.T) {
	world := NewWorld()
	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()
	is := NewInputSystem(world, physics, bindings, nil)

	is.SetState(InputState{Fire: true})

	if !is.IsFirePressed() {
		t.Error("expected Fire to be true after SetState")
	}
}

func TestInputSystem_ApplyThrust(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 0, Y: 0})
	world.AddComponent(entity, "velocity", &Velocity{VX: 0, VY: 0})
	world.AddComponent(entity, "rotation", &Rotation{Angle: 0})

	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()

	reader := &mockInputReader{state: InputState{Thrust: true}}
	is := NewInputSystem(world, physics, bindings, reader)
	is.SetPlayerEntity(entity)

	is.Update(1.0 / 60.0)

	velComp, _ := world.GetComponent(entity, "velocity")
	vel := velComp.(*Velocity)

	if vel.VX <= 0 {
		t.Errorf("expected positive velocity after thrust, got VX=%f", vel.VX)
	}
}

func TestInputSystem_ApplyRotation(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 0, Y: 0})
	world.AddComponent(entity, "velocity", &Velocity{VX: 0, VY: 0})
	world.AddComponent(entity, "rotation", &Rotation{Angle: 0})

	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()

	reader := &mockInputReader{state: InputState{RotateRight: true}}
	is := NewInputSystem(world, physics, bindings, reader)
	is.SetPlayerEntity(entity)

	is.Update(1.0 / 60.0)

	rotComp, _ := world.GetComponent(entity, "rotation")
	rot := rotComp.(*Rotation)

	if rot.Angle <= 0 {
		t.Errorf("expected positive angle after rotate right, got %f", rot.Angle)
	}
}

func TestInputSystem_UpdateWithNoPlayer(t *testing.T) {
	world := NewWorld()
	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()
	is := NewInputSystem(world, physics, bindings, nil)

	// Should not panic when no player entity is set
	is.Update(1.0 / 60.0)
}

func TestInputSystem_IsFirePressed(t *testing.T) {
	world := NewWorld()
	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()
	is := NewInputSystem(world, physics, bindings, nil)

	// Fire should initially be false
	if is.IsFirePressed() {
		t.Error("expected fire to be false initially")
	}
}

func TestInputSystem_IsPausePressed(t *testing.T) {
	world := NewWorld()
	config := DefaultPhysicsConfig()
	physics := NewPhysicsSystem(world, config)
	bindings := DefaultKeyBindings()
	is := NewInputSystem(world, physics, bindings, nil)

	// Pause should initially be false
	if is.IsPausePressed() {
		t.Error("expected pause to be false initially")
	}
}
