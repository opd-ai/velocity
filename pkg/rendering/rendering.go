// Package rendering provides sprite generation, animation, particle systems,
// dynamic lighting, draw batching, and viewport culling.
package rendering

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sync"

	"github.com/opd-ai/velocity/pkg/engine"
)

// Renderer handles all drawing operations.
type Renderer struct {
	genreID string
	cache   *SpriteCache
	rng     *rand.Rand
}

// NewRenderer creates a new renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		genreID: "scifi",
		cache:   NewSpriteCache(),
		rng:     rand.New(rand.NewSource(0)),
	}
}

// SetGenre switches rendering assets to match the given genre.
func (r *Renderer) SetGenre(genreID string) {
	r.genreID = genreID
}

// SetSeed sets the RNG seed for sprite generation.
func (r *Renderer) SetSeed(seed int64) {
	r.rng = rand.New(rand.NewSource(seed))
}

// GetGenre returns the current genre ID.
func (r *Renderer) GetGenre() string {
	return r.genreID
}

// GetOrCreateShipSprite returns a cached ship sprite or generates a new one.
func (r *Renderer) GetOrCreateShipSprite(variant int, size int) *image.RGBA {
	key := SpriteKey{GenreID: r.genreID, Type: SpriteTypeShip, Variant: variant}
	return r.cache.GetOrCreate(key, func() *image.RGBA {
		return GenerateShipSprite(r.rng, r.genreID, size)
	})
}

// GetOrCreateEnemySprite returns a cached enemy sprite or generates a new one.
func (r *Renderer) GetOrCreateEnemySprite(variant int, size int) *image.RGBA {
	key := SpriteKey{GenreID: r.genreID, Type: SpriteTypeEnemy, Variant: variant}
	return r.cache.GetOrCreate(key, func() *image.RGBA {
		return GenerateEnemySprite(r.rng, r.genreID, size)
	})
}

// GetOrCreateProjectileSprite returns a cached projectile sprite or generates a new one.
func (r *Renderer) GetOrCreateProjectileSprite(variant int, size int) *image.RGBA {
	key := SpriteKey{GenreID: r.genreID, Type: SpriteTypeProjectile, Variant: variant}
	return r.cache.GetOrCreate(key, func() *image.RGBA {
		return GenerateProjectileSprite(r.rng, r.genreID, size)
	})
}

// ClearCache clears the sprite cache (e.g., after genre change).
func (r *Renderer) ClearCache() {
	r.cache.Clear()
}

// SpriteCache caches generated sprite bitmaps keyed by genre and variant.
type SpriteCache struct {
	entries map[string]*image.RGBA
	mu      sync.RWMutex
}

// NewSpriteCache creates an empty sprite cache.
func NewSpriteCache() *SpriteCache {
	return &SpriteCache{entries: make(map[string]*image.RGBA)}
}

// keyString converts a SpriteKey to a cache key string.
func keyString(key SpriteKey) string {
	return fmt.Sprintf("%s:%d:%d", key.GenreID, key.Type, key.Variant)
}

// Get retrieves a sprite from the cache.
func (sc *SpriteCache) Get(key SpriteKey) (*image.RGBA, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	img, ok := sc.entries[keyString(key)]
	return img, ok
}

// Set stores a sprite in the cache.
func (sc *SpriteCache) Set(key SpriteKey, img *image.RGBA) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.entries[keyString(key)] = img
}

// GetOrCreate retrieves a sprite from cache or creates it using the generator.
func (sc *SpriteCache) GetOrCreate(key SpriteKey, gen func() *image.RGBA) *image.RGBA {
	if img, ok := sc.Get(key); ok {
		return img
	}
	img := gen()
	sc.Set(key, img)
	return img
}

// Clear removes all entries from the cache.
func (sc *SpriteCache) Clear() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.entries = make(map[string]*image.RGBA)
}

// Size returns the number of cached sprites.
func (sc *SpriteCache) Size() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return len(sc.entries)
}

// SpriteComponent holds sprite data for an entity.
type SpriteComponent struct {
	Type    SpriteType
	Variant int
	Size    int
}

// Particle represents a single particle in the particle system.
type Particle struct {
	X, Y    float64
	VX, VY  float64
	Life    float64
	MaxLife float64
	Color   color.RGBA
	Size    float64
}

// ParticleSystem manages a collection of particles.
type ParticleSystem struct {
	particles []Particle
	genreID   string
	rng       *rand.Rand
	mu        sync.Mutex
}

// NewParticleSystem creates an empty particle system.
func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		particles: make([]Particle, 0, 256),
		genreID:   "scifi",
		rng:       rand.New(rand.NewSource(0)),
	}
}

// SetGenre switches particle colors to match the given genre.
func (ps *ParticleSystem) SetGenre(genreID string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.genreID = genreID
}

// SetSeed sets the RNG seed for particle randomization.
func (ps *ParticleSystem) SetSeed(seed int64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.rng = rand.New(rand.NewSource(seed))
}

// Emit spawns new particles at the given position.
func (ps *ParticleSystem) Emit(x, y float64, count int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i := 0; i < count; i++ {
		angle := ps.rng.Float64() * 6.28318 // 2*PI
		speed := 20.0 + ps.rng.Float64()*80.0
		life := 0.3 + ps.rng.Float64()*0.5

		p := Particle{
			X:       x,
			Y:       y,
			VX:      speed * cosApprox(angle),
			VY:      speed * sinApprox(angle),
			Life:    life,
			MaxLife: life,
			Color:   ps.getParticleColor(),
			Size:    2.0 + ps.rng.Float64()*3.0,
		}
		ps.particles = append(ps.particles, p)
	}
}

// EmitDirectional spawns particles moving in a specific direction.
func (ps *ParticleSystem) EmitDirectional(x, y, angle, spread float64, count int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i := 0; i < count; i++ {
		particleAngle := angle + (ps.rng.Float64()-0.5)*spread
		speed := 30.0 + ps.rng.Float64()*60.0
		life := 0.2 + ps.rng.Float64()*0.3

		p := Particle{
			X:       x,
			Y:       y,
			VX:      speed * cosApprox(particleAngle),
			VY:      speed * sinApprox(particleAngle),
			Life:    life,
			MaxLife: life,
			Color:   ps.getParticleColor(),
			Size:    1.5 + ps.rng.Float64()*2.0,
		}
		ps.particles = append(ps.particles, p)
	}
}

// Update advances all particles by dt seconds and removes expired ones.
func (ps *ParticleSystem) Update(dt float64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	alive := ps.particles[:0]
	for i := range ps.particles {
		p := &ps.particles[i]
		p.X += p.VX * dt
		p.Y += p.VY * dt
		p.Life -= dt

		// Apply drag
		p.VX *= 0.98
		p.VY *= 0.98

		if p.Life > 0 {
			alive = append(alive, *p)
		}
	}
	ps.particles = alive
}

// GetParticles returns a copy of the current particles for rendering.
func (ps *ParticleSystem) GetParticles() []Particle {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	result := make([]Particle, len(ps.particles))
	copy(result, ps.particles)
	return result
}

// Count returns the number of active particles.
func (ps *ParticleSystem) Count() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.particles)
}

// getParticleColor returns a color appropriate for the current genre.
func (ps *ParticleSystem) getParticleColor() color.RGBA {
	switch ps.genreID {
	case "scifi":
		return color.RGBA{R: 100, G: 180, B: 255, A: 255}
	case "fantasy":
		return color.RGBA{R: 255, G: 200, B: 100, A: 255}
	case "horror":
		return color.RGBA{R: 150, G: 80, B: 80, A: 255}
	case "cyberpunk":
		return color.RGBA{R: 255, G: 0, B: 200, A: 255}
	case "postapoc":
		return color.RGBA{R: 200, G: 150, B: 80, A: 255}
	default:
		return color.RGBA{R: 255, G: 200, B: 100, A: 255}
	}
}

// cosApprox is a fast cosine approximation.
func cosApprox(x float64) float64 {
	// Taylor series approximation
	x2 := x * x
	return 1 - x2/2 + x2*x2/24
}

// sinApprox is a fast sine approximation.
func sinApprox(x float64) float64 {
	// Taylor series approximation
	x2 := x * x
	return x - x*x2/6 + x*x2*x2/120
}

// DrawBatch represents a batch of entities to draw together.
type DrawBatch struct {
	Type     SpriteType
	Entities []engine.Entity
}

// CreateDrawBatches groups entities by sprite type for batched rendering.
func CreateDrawBatches(world *engine.World) []DrawBatch {
	batches := make(map[SpriteType][]engine.Entity)

	world.ForEachEntity(func(e engine.Entity) {
		spriteComp, hasSprite := world.GetComponent(e, "sprite")
		if !hasSprite {
			return
		}
		sprite := spriteComp.(*SpriteComponent)
		batches[sprite.Type] = append(batches[sprite.Type], e)
	})

	result := make([]DrawBatch, 0, len(batches))
	for spriteType, entities := range batches {
		result = append(result, DrawBatch{
			Type:     spriteType,
			Entities: entities,
		})
	}
	return result
}
