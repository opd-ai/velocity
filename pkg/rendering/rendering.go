// Package rendering provides sprite generation, animation, particle systems,
// dynamic lighting, draw batching, and viewport culling.
package rendering

import "image/color"

// Renderer handles all drawing operations.
type Renderer struct {
	genreID string
}

// NewRenderer creates a new renderer.
func NewRenderer() *Renderer {
	return &Renderer{genreID: "scifi"}
}

// SetGenre switches rendering assets to match the given genre.
func (r *Renderer) SetGenre(genreID string) {
	r.genreID = genreID
}

// SpriteCache caches generated sprite bitmaps keyed by genre and variant.
type SpriteCache struct {
	entries map[string]interface{}
}

// NewSpriteCache creates an empty sprite cache.
func NewSpriteCache() *SpriteCache {
	return &SpriteCache{entries: make(map[string]interface{})}
}

// Particle represents a single particle in the particle system.
type Particle struct {
	X, Y     float64
	VX, VY   float64
	Life     float64
	MaxLife  float64
	Color    color.RGBA
}

// ParticleSystem manages a collection of particles.
type ParticleSystem struct {
	particles []Particle
}

// NewParticleSystem creates an empty particle system.
func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{}
}

// Emit spawns new particles at the given position.
func (ps *ParticleSystem) Emit(x, y float64, count int) {
	// Stub: will create particles with randomized velocities.
}

// Update advances all particles by dt seconds and removes expired ones.
func (ps *ParticleSystem) Update(dt float64) {
	// Stub: will update particle positions and lifetimes.
}
