package rendering

import (
	"testing"

	"github.com/opd-ai/velocity/pkg/engine"
)

func TestNewViewport(t *testing.T) {
	vp := NewViewport(800, 600)

	if vp.Width != 800 || vp.Height != 600 {
		t.Errorf("expected 800x600, got %fx%f", vp.Width, vp.Height)
	}
	if vp.X != 0 || vp.Y != 0 {
		t.Error("expected viewport at origin")
	}
}

func TestViewport_SetPosition(t *testing.T) {
	vp := NewViewport(800, 600)
	vp.SetPosition(400, 300)

	// Center at (400, 300) means top-left at (0, 0)
	if vp.X != 0 || vp.Y != 0 {
		t.Errorf("expected (0,0), got (%f,%f)", vp.X, vp.Y)
	}
}

func TestViewport_Contains(t *testing.T) {
	vp := NewViewport(800, 600)

	tests := []struct {
		x, y   float64
		inside bool
	}{
		{400, 300, true},    // Center
		{0, 0, true},        // Top-left corner
		{799, 599, true},    // Just inside bottom-right
		{800, 600, false},   // Just outside
		{-1, 0, false},      // Just outside left
		{400, -1, false},    // Just outside top
	}

	for _, tt := range tests {
		result := vp.Contains(tt.x, tt.y)
		if result != tt.inside {
			t.Errorf("Contains(%f, %f) = %v, want %v", tt.x, tt.y, result, tt.inside)
		}
	}
}

func TestViewport_ContainsRect(t *testing.T) {
	vp := NewViewport(800, 600)

	tests := []struct {
		name    string
		x, y    float64
		w, h    float64
		overlap bool
	}{
		{"fully inside", 100, 100, 50, 50, true},
		{"partially left", -25, 100, 50, 50, true},
		{"partially right", 780, 100, 50, 50, true},
		{"fully outside left", -100, 100, 50, 50, false},
		{"fully outside right", 850, 100, 50, 50, false},
		{"fully outside top", 100, -100, 50, 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vp.ContainsRect(tt.x, tt.y, tt.w, tt.h)
			if result != tt.overlap {
				t.Errorf("ContainsRect(%f,%f,%f,%f) = %v, want %v",
					tt.x, tt.y, tt.w, tt.h, result, tt.overlap)
			}
		})
	}
}

func TestNewCullContext(t *testing.T) {
	vp := NewViewport(800, 600)
	cc := NewCullContext(vp, 32)

	if cc == nil {
		t.Fatal("expected non-nil CullContext")
	}
}

func TestCullContext_ShouldRender(t *testing.T) {
	vp := NewViewport(800, 600)
	cc := NewCullContext(vp, 32)

	// Entity inside viewport
	if !cc.ShouldRender(400, 300, 32, 32) {
		t.Error("expected entity at center to render")
	}

	// Entity outside viewport
	if cc.ShouldRender(-100, -100, 32, 32) {
		t.Error("expected entity far outside to be culled")
	}

	// Entity just outside but within margin
	if !cc.ShouldRender(-16, 300, 32, 32) {
		t.Error("expected entity within margin to render")
	}

	if cc.GetRenderedCount() != 2 {
		t.Errorf("expected 2 rendered, got %d", cc.GetRenderedCount())
	}
	if cc.GetCulledCount() != 1 {
		t.Errorf("expected 1 culled, got %d", cc.GetCulledCount())
	}
}

func TestCullContext_Reset(t *testing.T) {
	vp := NewViewport(800, 600)
	cc := NewCullContext(vp, 32)

	cc.ShouldRender(400, 300, 32, 32)
	cc.ShouldRender(-1000, -1000, 32, 32)

	cc.Reset()

	if cc.GetRenderedCount() != 0 || cc.GetCulledCount() != 0 {
		t.Error("expected counters to be reset")
	}
}

func TestFilterVisibleEntities(t *testing.T) {
	world := engine.NewWorld()

	// Create entities at various positions
	e1 := world.CreateEntity()
	world.AddComponent(e1, "position", &engine.Position{X: 400, Y: 300}) // Inside

	e2 := world.CreateEntity()
	world.AddComponent(e2, "position", &engine.Position{X: -500, Y: -500}) // Outside

	e3 := world.CreateEntity()
	world.AddComponent(e3, "position", &engine.Position{X: 100, Y: 100}) // Inside

	vp := NewViewport(800, 600)
	visible := FilterVisibleEntities(world, vp, 32)

	if len(visible) != 2 {
		t.Errorf("expected 2 visible entities, got %d", len(visible))
	}
}

func TestSortBatchesByRenderOrder(t *testing.T) {
	batches := []DrawBatch{
		{Type: SpriteTypeShip, Entities: []engine.Entity{1}},
		{Type: SpriteTypeProjectile, Entities: []engine.Entity{2}},
		{Type: SpriteTypeEnemy, Entities: []engine.Entity{3}},
	}

	sorted := SortBatchesByRenderOrder(batches)

	// Expected order: Projectile, Enemy, Ship
	if sorted[0].Type != SpriteTypeProjectile {
		t.Errorf("expected first to be Projectile, got %d", sorted[0].Type)
	}
	if sorted[1].Type != SpriteTypeEnemy {
		t.Errorf("expected second to be Enemy, got %d", sorted[1].Type)
	}
	if sorted[2].Type != SpriteTypeShip {
		t.Errorf("expected third to be Ship, got %d", sorted[2].Type)
	}
}

func TestRenderOrder(t *testing.T) {
	if len(RenderOrder) != 3 {
		t.Errorf("expected 3 sprite types in render order, got %d", len(RenderOrder))
	}
}

func BenchmarkFilterVisibleEntities(b *testing.B) {
	world := engine.NewWorld()

	// Create 200 entities
	for i := 0; i < 200; i++ {
		e := world.CreateEntity()
		world.AddComponent(e, "position", &engine.Position{
			X: float64(i % 100) * 10,
			Y: float64(i / 100) * 10,
		})
	}

	vp := NewViewport(800, 600)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterVisibleEntities(world, vp, 32)
	}
}
