package rendering

import (
	"image"
	"image/color"
	"math/rand"
	"testing"

	"github.com/opd-ai/velocity/pkg/engine"
)

func TestNewRenderer(t *testing.T) {
	r := NewRenderer()

	if r == nil {
		t.Fatal("expected non-nil Renderer")
	}

	if r.GetGenre() != "scifi" {
		t.Errorf("expected default genre 'scifi', got %s", r.GetGenre())
	}
}

func TestRenderer_SetGenre(t *testing.T) {
	r := NewRenderer()
	r.SetGenre("fantasy")

	if r.GetGenre() != "fantasy" {
		t.Errorf("expected genre 'fantasy', got %s", r.GetGenre())
	}
}

func TestRenderer_GetOrCreateShipSprite(t *testing.T) {
	r := NewRenderer()
	r.SetSeed(12345)

	sprite := r.GetOrCreateShipSprite(0, 16)
	if sprite == nil {
		t.Fatal("expected non-nil sprite")
	}

	// Second call should return cached sprite
	sprite2 := r.GetOrCreateShipSprite(0, 16)
	if sprite != sprite2 {
		t.Error("expected cached sprite to be returned")
	}
}

func TestRenderer_ClearCache(t *testing.T) {
	r := NewRenderer()
	r.SetSeed(12345)

	_ = r.GetOrCreateShipSprite(0, 16)
	if r.cache.Size() == 0 {
		t.Error("expected cache to have entries")
	}

	r.ClearCache()
	if r.cache.Size() != 0 {
		t.Error("expected cache to be empty after clear")
	}
}

func TestSpriteCache_GetSet(t *testing.T) {
	cache := NewSpriteCache()
	key := SpriteKey{GenreID: "scifi", Type: SpriteTypeShip, Variant: 0}

	_, ok := cache.Get(key)
	if ok {
		t.Error("expected cache miss")
	}

	rng := rand.New(rand.NewSource(12345))
	sprite := GenerateShipSprite(rng, "scifi", 16)
	cache.Set(key, sprite)

	retrieved, ok := cache.Get(key)
	if !ok {
		t.Error("expected cache hit")
	}
	if retrieved != sprite {
		t.Error("expected same sprite")
	}
}

func TestSpriteCache_GetOrCreate(t *testing.T) {
	cache := NewSpriteCache()
	key := SpriteKey{GenreID: "scifi", Type: SpriteTypeShip, Variant: 0}

	callCount := 0
	generator := func() *image.RGBA {
		callCount++
		rng := rand.New(rand.NewSource(12345))
		return GenerateShipSprite(rng, "scifi", 16)
	}

	// First call should invoke generator
	sprite1 := cache.GetOrCreate(key, generator)
	if sprite1 == nil {
		t.Fatal("expected non-nil sprite")
	}
	if callCount != 1 {
		t.Errorf("expected 1 generator call, got %d", callCount)
	}

	// Second call should use cache
	sprite2 := cache.GetOrCreate(key, generator)
	if sprite2 != sprite1 {
		t.Error("expected cached sprite")
	}
	if callCount != 1 {
		t.Errorf("expected 1 generator call (cached), got %d", callCount)
	}
}

func TestNewParticleSystem(t *testing.T) {
	ps := NewParticleSystem()
	if ps == nil {
		t.Fatal("expected non-nil ParticleSystem")
	}
	if ps.Count() != 0 {
		t.Error("expected empty particle system")
	}
}

func TestParticleSystem_Emit(t *testing.T) {
	ps := NewParticleSystem()
	ps.SetSeed(12345)

	ps.Emit(100, 100, 10)

	if ps.Count() != 10 {
		t.Errorf("expected 10 particles, got %d", ps.Count())
	}
}

func TestParticleSystem_Update(t *testing.T) {
	ps := NewParticleSystem()
	ps.SetSeed(12345)

	ps.Emit(100, 100, 5)
	initialCount := ps.Count()

	// Update many times to expire particles
	for i := 0; i < 100; i++ {
		ps.Update(0.1)
	}

	if ps.Count() >= initialCount {
		t.Error("expected particles to expire")
	}
}

func TestParticleSystem_GetParticles(t *testing.T) {
	ps := NewParticleSystem()
	ps.SetSeed(12345)

	ps.Emit(100, 100, 5)

	particles := ps.GetParticles()
	if len(particles) != 5 {
		t.Errorf("expected 5 particles, got %d", len(particles))
	}

	// Verify it's a copy by modifying
	particles[0].X = 999
	newParticles := ps.GetParticles()
	if newParticles[0].X == 999 {
		t.Error("expected GetParticles to return a copy")
	}
}

func TestParticleSystem_SetGenre(t *testing.T) {
	ps := NewParticleSystem()

	for _, genre := range []string{"scifi", "fantasy", "horror", "cyberpunk", "postapoc"} {
		ps.SetGenre(genre)
		ps.Emit(0, 0, 1)
		// Should not panic
	}
}

func TestParticle_Fields(t *testing.T) {
	p := Particle{
		X:       100,
		Y:       200,
		VX:      10,
		VY:      -5,
		Life:    1.0,
		MaxLife: 1.0,
		Color:   color.RGBA{R: 255, G: 0, B: 0, A: 255},
		Size:    3.0,
	}

	if p.X != 100 || p.Y != 200 {
		t.Error("position mismatch")
	}
	if p.VX != 10 || p.VY != -5 {
		t.Error("velocity mismatch")
	}
	if p.Color.R != 255 {
		t.Error("color mismatch")
	}
}

func TestCreateDrawBatches(t *testing.T) {
	world := engine.NewWorld()

	// Create entities with sprites
	e1 := world.CreateEntity()
	world.AddComponent(e1, "sprite", &SpriteComponent{Type: SpriteTypeShip, Variant: 0, Size: 16})

	e2 := world.CreateEntity()
	world.AddComponent(e2, "sprite", &SpriteComponent{Type: SpriteTypeEnemy, Variant: 0, Size: 16})

	e3 := world.CreateEntity()
	world.AddComponent(e3, "sprite", &SpriteComponent{Type: SpriteTypeShip, Variant: 1, Size: 16})

	batches := CreateDrawBatches(world)

	if len(batches) == 0 {
		t.Fatal("expected non-empty batches")
	}

	// Count ships and enemies
	shipCount := 0
	enemyCount := 0
	for _, batch := range batches {
		if batch.Type == SpriteTypeShip {
			shipCount += len(batch.Entities)
		}
		if batch.Type == SpriteTypeEnemy {
			enemyCount += len(batch.Entities)
		}
	}

	if shipCount != 2 {
		t.Errorf("expected 2 ships in batch, got %d", shipCount)
	}
	if enemyCount != 1 {
		t.Errorf("expected 1 enemy in batch, got %d", enemyCount)
	}
}

func TestSpriteComponent(t *testing.T) {
	sc := SpriteComponent{
		Type:    SpriteTypeEnemy,
		Variant: 5,
		Size:    32,
	}

	if sc.Type != SpriteTypeEnemy {
		t.Error("type mismatch")
	}
	if sc.Variant != 5 {
		t.Error("variant mismatch")
	}
	if sc.Size != 32 {
		t.Error("size mismatch")
	}
}

func BenchmarkParticleSystem_Update(b *testing.B) {
	ps := NewParticleSystem()
	ps.SetSeed(12345)

	// Add many particles
	ps.Emit(100, 100, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.Update(1.0 / 60.0)
		// Re-emit to maintain particle count
		if ps.Count() < 50 {
			ps.Emit(100, 100, 50)
		}
	}
}
