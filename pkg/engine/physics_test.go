package engine

import (
	"math"
	"testing"
)

func TestPhysicsSystem_Update(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 100, Y: 100})
	world.AddComponent(entity, "velocity", &Velocity{VX: 60, VY: 0})

	config := DefaultPhysicsConfig()
	ps := NewPhysicsSystem(world, config)

	// One tick at dt = 1/60
	dt := 1.0 / 60.0
	ps.Update(dt)

	posComp, _ := world.GetComponent(entity, "position")
	pos := posComp.(*Position)

	// Position should increase by velocity * dt (after drag)
	if pos.X <= 100 {
		t.Errorf("expected position to increase, got X=%f", pos.X)
	}
}

func TestPhysicsSystem_ApplyThrust(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 0, Y: 0})
	world.AddComponent(entity, "velocity", &Velocity{VX: 0, VY: 0})
	world.AddComponent(entity, "rotation", &Rotation{Angle: 0}) // Facing right

	config := DefaultPhysicsConfig()
	ps := NewPhysicsSystem(world, config)

	dt := 1.0 / 60.0
	ps.ApplyThrust(entity, dt)

	velComp, _ := world.GetComponent(entity, "velocity")
	vel := velComp.(*Velocity)

	// Thrust along facing direction (right, angle=0)
	if vel.VX <= 0 {
		t.Errorf("expected positive VX after thrust, got %f", vel.VX)
	}
	if math.Abs(vel.VY) > 0.001 {
		t.Errorf("expected VY near zero, got %f", vel.VY)
	}
}

func TestPhysicsSystem_ApplyRotation(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "rotation", &Rotation{Angle: 0})

	config := DefaultPhysicsConfig()
	ps := NewPhysicsSystem(world, config)

	dt := 1.0 / 60.0
	ps.ApplyRotation(entity, 1.0, dt) // Rotate clockwise

	rotComp, _ := world.GetComponent(entity, "rotation")
	rot := rotComp.(*Rotation)

	if rot.Angle <= 0 {
		t.Errorf("expected angle to increase, got %f", rot.Angle)
	}
}

func TestPhysicsSystem_DragApplied(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 0, Y: 0})
	world.AddComponent(entity, "velocity", &Velocity{VX: 100, VY: 0})

	config := DefaultPhysicsConfig()
	ps := NewPhysicsSystem(world, config)

	// After update, velocity should be reduced by drag
	ps.Update(1.0 / 60.0)

	velComp, _ := world.GetComponent(entity, "velocity")
	vel := velComp.(*Velocity)

	if vel.VX >= 100 {
		t.Errorf("expected drag to reduce velocity, got VX=%f", vel.VX)
	}
}

func TestPhysicsSystem_MaxSpeedClamped(t *testing.T) {
	world := NewWorld()
	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &Position{X: 0, Y: 0})
	world.AddComponent(entity, "velocity", &Velocity{VX: 1000, VY: 0}) // Exceeds max

	config := DefaultPhysicsConfig()
	ps := NewPhysicsSystem(world, config)

	ps.Update(1.0 / 60.0)

	velComp, _ := world.GetComponent(entity, "velocity")
	vel := velComp.(*Velocity)

	speed := math.Sqrt(vel.VX*vel.VX + vel.VY*vel.VY)
	if speed > config.MaxSpeed+0.001 {
		t.Errorf("expected speed clamped to max %f, got %f", config.MaxSpeed, speed)
	}
}

func TestDefaultPhysicsConfig(t *testing.T) {
	config := DefaultPhysicsConfig()

	if config.ThrustForce <= 0 {
		t.Error("expected positive thrust force")
	}
	if config.RotationSpeed <= 0 {
		t.Error("expected positive rotation speed")
	}
	if config.DragCoeff <= 0 || config.DragCoeff > 1 {
		t.Error("expected drag coefficient in (0, 1]")
	}
	if config.MaxSpeed <= 0 {
		t.Error("expected positive max speed")
	}
}
