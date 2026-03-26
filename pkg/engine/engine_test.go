package engine

import (
	"testing"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld()
	if w == nil {
		t.Fatal("expected non-nil World")
	}
}

func TestWorld_CreateEntity(t *testing.T) {
	w := NewWorld()
	e1 := w.CreateEntity()
	e2 := w.CreateEntity()

	if e1 == e2 {
		t.Error("expected unique entity IDs")
	}
}

func TestWorld_AddGetComponent(t *testing.T) {
	w := NewWorld()
	e := w.CreateEntity()

	pos := &Position{X: 100, Y: 200}
	w.AddComponent(e, "position", pos)

	comp, ok := w.GetComponent(e, "position")
	if !ok {
		t.Fatal("expected component to exist")
	}

	retrieved := comp.(*Position)
	if retrieved.X != 100 || retrieved.Y != 200 {
		t.Error("position mismatch")
	}
}

func TestWorld_RemoveEntity(t *testing.T) {
	w := NewWorld()
	e := w.CreateEntity()
	w.AddComponent(e, "position", &Position{X: 0, Y: 0})

	w.RemoveEntity(e)

	_, ok := w.GetComponent(e, "position")
	if ok {
		t.Error("expected component to be removed")
	}
}

func TestWorld_ForEachEntity(t *testing.T) {
	w := NewWorld()
	w.CreateEntity()
	w.CreateEntity()
	w.CreateEntity()

	count := 0
	w.ForEachEntity(func(e Entity) {
		count++
	})

	if count != 3 {
		t.Errorf("expected 3 entities, got %d", count)
	}
}

func TestWorld_EntityCount(t *testing.T) {
	w := NewWorld()

	if w.EntityCount() != 0 {
		t.Error("expected empty world")
	}

	w.CreateEntity()
	w.CreateEntity()

	if w.EntityCount() != 2 {
		t.Errorf("expected 2 entities, got %d", w.EntityCount())
	}
}

func TestWorld_AddSystem(t *testing.T) {
	w := NewWorld()

	mockSys := &mockSystem{}
	w.AddSystem(mockSys)

	w.Update(1.0 / 60.0)

	if mockSys.updateCount != 1 {
		t.Errorf("expected 1 update call, got %d", mockSys.updateCount)
	}
}

type mockSystem struct {
	updateCount int
}

func (m *mockSystem) Update(dt float64) {
	m.updateCount++
}

func TestDeterministicRNG(t *testing.T) {
	rng1 := DeterministicRNG(12345)
	rng2 := DeterministicRNG(12345)

	for i := 0; i < 10; i++ {
		v1 := rng1.Intn(1000)
		v2 := rng2.Intn(1000)
		if v1 != v2 {
			t.Errorf("expected deterministic values, got %d vs %d", v1, v2)
		}
	}
}

func TestInputState(t *testing.T) {
	state := InputState{
		Thrust:      true,
		RotateLeft:  false,
		RotateRight: true,
		Fire:        true,
		Secondary:   false,
		Pause:       false,
	}

	if !state.Thrust {
		t.Error("expected Thrust to be true")
	}
	if state.Pause {
		t.Error("expected Pause to be false")
	}
}

func TestNewCamera(t *testing.T) {
	c := NewCamera()
	if c == nil {
		t.Fatal("expected non-nil Camera")
	}
	if c.X != 0 || c.Y != 0 {
		t.Error("expected camera at origin")
	}
}

func TestCamera_Shake(t *testing.T) {
	c := NewCamera()
	c.Shake(10.0, 0.5)

	if c.ShakeAmount != 10.0 {
		t.Errorf("expected ShakeAmount 10.0, got %f", c.ShakeAmount)
	}
	if c.ShakeDuration != 0.5 {
		t.Errorf("expected ShakeDuration 0.5, got %f", c.ShakeDuration)
	}
}

func TestCamera_Update(t *testing.T) {
	c := NewCamera()
	c.Shake(10.0, 0.2)

	// Update past shake duration
	c.Update(0.3)

	if c.ShakeAmount != 0 {
		t.Errorf("expected ShakeAmount 0 after duration, got %f", c.ShakeAmount)
	}
	if c.ShakeDuration != 0 {
		t.Errorf("expected ShakeDuration 0 after duration, got %f", c.ShakeDuration)
	}
}

func TestCamera_UpdatePartial(t *testing.T) {
	c := NewCamera()
	c.Shake(10.0, 0.5)

	// Update partially
	c.Update(0.2)

	if c.ShakeDuration <= 0 {
		t.Error("expected shake to still be active")
	}
	if c.ShakeAmount != 10.0 {
		t.Error("expected shake amount unchanged during shake")
	}
}

func TestWorld_GetComponentNotFound(t *testing.T) {
	w := NewWorld()
	e := w.CreateEntity()

	_, ok := w.GetComponent(e, "nonexistent")
	if ok {
		t.Error("expected component not found")
	}
}

func TestWorld_GetComponentInvalidEntity(t *testing.T) {
	w := NewWorld()

	_, ok := w.GetComponent(Entity(999), "position")
	if ok {
		t.Error("expected component not found for invalid entity")
	}
}
