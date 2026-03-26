package engine

import (
	"testing"
)

func TestNewArenaSystem(t *testing.T) {
	world := NewWorld()
	as := NewArenaSystem(world, 800, 600, ArenaModeWrap)

	if as == nil {
		t.Fatal("expected non-nil ArenaSystem")
	}
}

func TestArenaSystem_WrapMode_LeftEdge(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: -10, Y: 300})

	as := NewArenaSystem(world, 800, 600, ArenaModeWrap)
	as.Update(0)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	if pos.X < 0 || pos.X >= 800 {
		t.Errorf("expected wrapped X position, got %f", pos.X)
	}
}

func TestArenaSystem_WrapMode_RightEdge(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 810, Y: 300})

	as := NewArenaSystem(world, 800, 600, ArenaModeWrap)
	as.Update(0)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	if pos.X < 0 || pos.X >= 800 {
		t.Errorf("expected wrapped X position, got %f", pos.X)
	}
}

func TestArenaSystem_WrapMode_TopEdge(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 400, Y: -10})

	as := NewArenaSystem(world, 800, 600, ArenaModeWrap)
	as.Update(0)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	if pos.Y < 0 || pos.Y >= 600 {
		t.Errorf("expected wrapped Y position, got %f", pos.Y)
	}
}

func TestArenaSystem_WrapMode_BottomEdge(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 400, Y: 610})

	as := NewArenaSystem(world, 800, 600, ArenaModeWrap)
	as.Update(0)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	if pos.Y < 0 || pos.Y >= 600 {
		t.Errorf("expected wrapped Y position, got %f", pos.Y)
	}
}

func TestArenaSystem_BoundedMode_LeftEdge(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: -10, Y: 300})
	world.AddComponent(entity, "velocity", &Velocity{VX: -50, VY: 0})

	as := NewArenaSystem(world, 800, 600, ArenaModeBounded)
	as.Update(0)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	velComp, _ := world.GetComponent(entity, "velocity")
	vel := velComp.(*Velocity)

	if pos.X < 0 {
		t.Errorf("expected position clamped to 0, got %f", pos.X)
	}
	if vel.VX <= 0 {
		t.Errorf("expected velocity reversed, got VX=%f", vel.VX)
	}
}

func TestArenaSystem_BoundedMode_RightEdge(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 810, Y: 300})
	world.AddComponent(entity, "velocity", &Velocity{VX: 50, VY: 0})

	as := NewArenaSystem(world, 800, 600, ArenaModeBounded)
	as.Update(0)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	velComp, _ := world.GetComponent(entity, "velocity")
	vel := velComp.(*Velocity)

	if pos.X >= 800 {
		t.Errorf("expected position clamped to <800, got %f", pos.X)
	}
	if vel.VX >= 0 {
		t.Errorf("expected velocity reversed, got VX=%f", vel.VX)
	}
}

func TestArenaSystem_BoundedMode_NoVelocity(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: -10, Y: 300})

	as := NewArenaSystem(world, 800, 600, ArenaModeBounded)

	// Should not panic without velocity component
	as.Update(0)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	if pos.X < 0 {
		t.Errorf("expected position clamped to 0, got %f", pos.X)
	}
}

func TestArenaSystem_SetMode(t *testing.T) {
	world := NewWorld()
	as := NewArenaSystem(world, 800, 600, ArenaModeWrap)

	as.SetMode(ArenaModeBounded)

	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: -10, Y: 300})
	world.AddComponent(entity, "velocity", &Velocity{VX: -50, VY: 0})

	as.Update(0)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	// Should clamp rather than wrap since we changed mode
	if pos.X != 0 {
		t.Errorf("expected position clamped to 0 in bounded mode, got %f", pos.X)
	}
}

func TestArenaSystem_IsEntityInBounds(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 400, Y: 300})

	as := NewArenaSystem(world, 800, 600, ArenaModeWrap)

	if !as.IsEntityInBounds(entity) {
		t.Error("expected entity to be in bounds")
	}

	// Move out of bounds
	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)
	pos.X = -10

	if as.IsEntityInBounds(entity) {
		t.Error("expected entity to be out of bounds")
	}
}

func TestArenaSystem_UpdateNoPosition(t *testing.T) {
	world := NewWorld()
	_ = world.CreateEntity()
	// No position component

	as := NewArenaSystem(world, 800, 600, ArenaModeWrap)

	// Should not panic
	as.Update(0)
}
