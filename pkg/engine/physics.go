// Package engine provides the core ECS framework, game loop, deterministic RNG,
// input handling, and camera system.
package engine

import "math"

// Position component stores an entity's 2D location.
type Position struct {
	X, Y float64
}

// Velocity component stores an entity's 2D movement vector.
type Velocity struct {
	VX, VY float64
}

// Rotation component stores an entity's facing angle in radians.
type Rotation struct {
	Angle float64
}

// PhysicsConfig holds physics tuning parameters.
type PhysicsConfig struct {
	ThrustForce   float64
	RotationSpeed float64
	DragCoeff     float64
	MaxSpeed      float64
}

// DefaultPhysicsConfig returns the default physics tuning parameters.
func DefaultPhysicsConfig() PhysicsConfig {
	return PhysicsConfig{
		ThrustForce:   200.0,
		RotationSpeed: 4.0,
		DragCoeff:     0.98,
		MaxSpeed:      400.0,
	}
}

// PhysicsSystem applies Newtonian 2D flight physics to entities.
type PhysicsSystem struct {
	world  *World
	config PhysicsConfig
}

// NewPhysicsSystem creates a physics system attached to the given world.
func NewPhysicsSystem(world *World, config PhysicsConfig) *PhysicsSystem {
	return &PhysicsSystem{
		world:  world,
		config: config,
	}
}

// Update applies physics to all entities with position and velocity components.
func (ps *PhysicsSystem) Update(dt float64) {
	for entity := range ps.world.entities {
		ps.updateEntity(entity, dt)
	}
}

// updateEntity processes physics for a single entity.
func (ps *PhysicsSystem) updateEntity(entity Entity, dt float64) {
	posComp, hasPos := ps.world.GetComponent(entity, "position")
	velComp, hasVel := ps.world.GetComponent(entity, "velocity")

	if !hasPos || !hasVel {
		return
	}

	pos := posComp.(*Position)
	vel := velComp.(*Velocity)

	// Apply drag
	vel.VX *= ps.config.DragCoeff
	vel.VY *= ps.config.DragCoeff

	// Clamp to max speed
	speed := math.Sqrt(vel.VX*vel.VX + vel.VY*vel.VY)
	if speed > ps.config.MaxSpeed {
		scale := ps.config.MaxSpeed / speed
		vel.VX *= scale
		vel.VY *= scale
	}

	// Update position from velocity
	pos.X += vel.VX * dt
	pos.Y += vel.VY * dt
}

// ApplyThrust applies thrust acceleration to an entity along its rotation.
func (ps *PhysicsSystem) ApplyThrust(entity Entity, dt float64) {
	velComp, hasVel := ps.world.GetComponent(entity, "velocity")
	rotComp, hasRot := ps.world.GetComponent(entity, "rotation")

	if !hasVel || !hasRot {
		return
	}

	vel := velComp.(*Velocity)
	rot := rotComp.(*Rotation)

	// Thrust applies acceleration along ship facing vector
	accelX := math.Cos(rot.Angle) * ps.config.ThrustForce * dt
	accelY := math.Sin(rot.Angle) * ps.config.ThrustForce * dt

	vel.VX += accelX
	vel.VY += accelY
}

// ApplyRotation rotates an entity by the given direction.
func (ps *PhysicsSystem) ApplyRotation(entity Entity, direction, dt float64) {
	rotComp, hasRot := ps.world.GetComponent(entity, "rotation")
	if !hasRot {
		return
	}

	rot := rotComp.(*Rotation)
	rot.Angle += direction * ps.config.RotationSpeed * dt

	// Normalize angle to [0, 2π)
	for rot.Angle < 0 {
		rot.Angle += 2 * math.Pi
	}
	for rot.Angle >= 2*math.Pi {
		rot.Angle -= 2 * math.Pi
	}
}
