// Package engine provides the core ECS framework, game loop, deterministic RNG,
// input handling, and camera system.
package engine

import "math/rand"

// GenreSetter is the interface that all output-producing systems must implement.
type GenreSetter interface {
	SetGenre(genreID string)
}

// Entity represents a unique game object identifier.
type Entity uint64

// Component is the interface for all ECS components.
type Component interface{}

// System is the interface for all ECS systems.
type System interface {
	Update(dt float64)
}

// World holds all entities and their components.
type World struct {
	nextID   Entity
	entities map[Entity]map[string]Component
	systems  []System
}

// NewWorld creates a new empty ECS world.
func NewWorld() *World {
	return &World{
		entities: make(map[Entity]map[string]Component),
	}
}

// CreateEntity creates a new entity and returns its ID.
func (w *World) CreateEntity() Entity {
	w.nextID++
	w.entities[w.nextID] = make(map[string]Component)
	return w.nextID
}

// AddComponent adds a named component to an entity.
func (w *World) AddComponent(e Entity, name string, c Component) {
	if comps, ok := w.entities[e]; ok {
		comps[name] = c
	}
}

// GetComponent returns a component by name for an entity.
func (w *World) GetComponent(e Entity, name string) (Component, bool) {
	if comps, ok := w.entities[e]; ok {
		c, exists := comps[name]
		return c, exists
	}
	return nil, false
}

// RemoveEntity removes an entity and all its components.
func (w *World) RemoveEntity(e Entity) {
	delete(w.entities, e)
}

// AddSystem registers a system with the world.
func (w *World) AddSystem(s System) {
	w.systems = append(w.systems, s)
}

// Update runs all registered systems.
func (w *World) Update(dt float64) {
	for _, s := range w.systems {
		s.Update(dt)
	}
}

// DeterministicRNG returns a seeded random source for reproducible runs.
func DeterministicRNG(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// InputState tracks the current state of player input.
type InputState struct {
	Thrust      bool
	RotateLeft  bool
	RotateRight bool
	Fire        bool
	Secondary   bool
	Pause       bool
}

// Camera tracks the viewport position and screen-shake state.
type Camera struct {
	X, Y          float64
	ShakeAmount   float64
	ShakeDuration float64
}

// NewCamera creates a new camera at the origin.
func NewCamera() *Camera {
	return &Camera{}
}

// Shake applies a screen-shake effect.
func (c *Camera) Shake(amount, duration float64) {
	c.ShakeAmount = amount
	c.ShakeDuration = duration
}

// Update advances the camera state by dt seconds.
func (c *Camera) Update(dt float64) {
	if c.ShakeDuration > 0 {
		c.ShakeDuration -= dt
		if c.ShakeDuration <= 0 {
			c.ShakeAmount = 0
			c.ShakeDuration = 0
		}
	}
}
