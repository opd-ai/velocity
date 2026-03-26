package combat

import (
	"math"
	"testing"

	"github.com/opd-ai/velocity/pkg/engine"
)

func TestNewProjectileSystem(t *testing.T) {
	world := engine.NewWorld()
	ps := NewProjectileSystem(world)

	if ps == nil {
		t.Fatal("expected non-nil ProjectileSystem")
	}
}

func TestProjectileSystem_SpawnProjectile(t *testing.T) {
	world := engine.NewWorld()
	ps := NewProjectileSystem(world)

	e := ps.SpawnProjectile(100, 200, 0, 500, 10, "player", 2.0)

	if e == 0 {
		t.Fatal("expected valid entity ID")
	}

	// Verify components
	posComp, hasPos := world.GetComponent(e, "position")
	if !hasPos {
		t.Fatal("expected position component")
	}
	pos := posComp.(*engine.Position)
	if pos.X != 100 || pos.Y != 200 {
		t.Errorf("expected position (100,200), got (%f,%f)", pos.X, pos.Y)
	}

	projComp, hasProj := world.GetComponent(e, "projectile")
	if !hasProj {
		t.Fatal("expected projectile component")
	}
	proj := projComp.(*Projectile)
	if proj.Damage != 10 {
		t.Errorf("expected damage 10, got %f", proj.Damage)
	}
	if proj.OwnerType != "player" {
		t.Errorf("expected owner 'player', got %s", proj.OwnerType)
	}
}

func TestProjectileSystem_Update_Movement(t *testing.T) {
	world := engine.NewWorld()
	ps := NewProjectileSystem(world)

	e := ps.SpawnProjectile(0, 0, 0, 100, 10, "player", 2.0) // Angle 0 = right

	dt := 0.1
	ps.Update(dt)

	posComp, _ := world.GetComponent(e, "position")
	pos := posComp.(*engine.Position)

	// Should have moved right
	if pos.X <= 0 {
		t.Errorf("expected X > 0 after update, got %f", pos.X)
	}
}

func TestProjectileSystem_Update_Expiry(t *testing.T) {
	world := engine.NewWorld()
	ps := NewProjectileSystem(world)

	_ = ps.SpawnProjectile(0, 0, 0, 100, 10, "player", 0.1) // Short lifetime

	if ps.ProjectileCount() != 1 {
		t.Fatal("expected 1 projectile")
	}

	// Update past lifetime
	ps.Update(0.2)

	if ps.ProjectileCount() != 0 {
		t.Errorf("expected projectile to expire, count=%d", ps.ProjectileCount())
	}
}

func TestProjectileSystem_Collision(t *testing.T) {
	world := engine.NewWorld()
	ps := NewProjectileSystem(world)

	// Create an enemy target
	enemy := world.CreateEntity()
	world.AddComponent(enemy, "position", &engine.Position{X: 50, Y: 0})
	world.AddComponent(enemy, "collisiontag", &CollisionTag{Tag: "enemy"})
	world.AddComponent(enemy, "boundingbox", &BoundingBox{X: -16, Y: -16, Width: 32, Height: 32})

	// Track hit
	hitCount := 0
	ps.SetHitCallback(func(proj, target engine.Entity, damage float64) {
		hitCount++
		if target != enemy {
			t.Error("expected enemy to be hit")
		}
		if damage != 10 {
			t.Errorf("expected damage 10, got %f", damage)
		}
	})

	// Spawn projectile heading toward enemy
	ps.SpawnProjectile(0, 0, 0, 100, 10, "player", 2.0)

	// Update until collision
	for i := 0; i < 10; i++ {
		ps.Update(0.1)
	}

	if hitCount != 1 {
		t.Errorf("expected 1 hit, got %d", hitCount)
	}
}

func TestProjectileSystem_NoSelfCollision(t *testing.T) {
	world := engine.NewWorld()
	ps := NewProjectileSystem(world)

	// Create player
	player := world.CreateEntity()
	world.AddComponent(player, "position", &engine.Position{X: 10, Y: 0})
	world.AddComponent(player, "collisiontag", &CollisionTag{Tag: "player"})
	world.AddComponent(player, "boundingbox", &BoundingBox{X: -16, Y: -16, Width: 32, Height: 32})

	hitCount := 0
	ps.SetHitCallback(func(proj, target engine.Entity, damage float64) {
		hitCount++
	})

	// Player projectile should not hit player
	ps.SpawnProjectile(0, 0, 0, 100, 10, "player", 2.0)

	for i := 0; i < 10; i++ {
		ps.Update(0.1)
	}

	if hitCount != 0 {
		t.Errorf("expected no hits on same team, got %d", hitCount)
	}
}

func TestCheckAABBCollision(t *testing.T) {
	tests := []struct {
		name    string
		ax, ay  float64
		aw, ah  float64
		bx, by  float64
		bw, bh  float64
		collide bool
	}{
		{"overlapping", 0, 0, 10, 10, 5, 5, 10, 10, true},
		{"touching", 0, 0, 10, 10, 10, 0, 10, 10, false}, // Touching edge, not overlapping
		{"separate x", 0, 0, 10, 10, 20, 0, 10, 10, false},
		{"separate y", 0, 0, 10, 10, 0, 20, 10, 10, false},
		{"contained", 5, 5, 5, 5, 0, 0, 20, 20, true},
		{"same position", 0, 0, 10, 10, 0, 0, 10, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckAABBCollision(tt.ax, tt.ay, tt.aw, tt.ah, tt.bx, tt.by, tt.bw, tt.bh)
			if result != tt.collide {
				t.Errorf("expected %v, got %v", tt.collide, result)
			}
		})
	}
}

func TestProjectile_Fields(t *testing.T) {
	p := Projectile{
		Damage:    25,
		Speed:     300,
		OwnerType: "enemy",
		Lifetime:  1.5,
		MaxLife:   2.0,
	}

	if p.Damage != 25 {
		t.Error("damage mismatch")
	}
	if p.OwnerType != "enemy" {
		t.Error("owner mismatch")
	}
}

func TestHealth_Fields(t *testing.T) {
	h := Health{Current: 80, Max: 100}

	if h.Current != 80 || h.Max != 100 {
		t.Error("health mismatch")
	}
}

func TestBoundingBox_Fields(t *testing.T) {
	bb := BoundingBox{X: -10, Y: -10, Width: 20, Height: 20}

	if bb.Width != 20 || bb.Height != 20 {
		t.Error("bounding box mismatch")
	}
}

func TestCollisionTag_Fields(t *testing.T) {
	ct := CollisionTag{Tag: "player"}

	if ct.Tag != "player" {
		t.Error("tag mismatch")
	}
}

func TestProjectileSystem_ProjectileCount(t *testing.T) {
	world := engine.NewWorld()
	ps := NewProjectileSystem(world)

	if ps.ProjectileCount() != 0 {
		t.Error("expected 0 projectiles")
	}

	ps.SpawnProjectile(0, 0, 0, 100, 10, "player", 2.0)
	ps.SpawnProjectile(0, 0, math.Pi, 100, 10, "player", 2.0)

	if ps.ProjectileCount() != 2 {
		t.Errorf("expected 2 projectiles, got %d", ps.ProjectileCount())
	}
}

func BenchmarkProjectileSystem_Update(b *testing.B) {
	world := engine.NewWorld()
	ps := NewProjectileSystem(world)

	// Spawn many projectiles
	for i := 0; i < 100; i++ {
		ps.SpawnProjectile(float64(i*10), 0, 0, 100, 10, "player", 10.0)
	}

	// Create some enemies
	for i := 0; i < 20; i++ {
		e := world.CreateEntity()
		world.AddComponent(e, "position", &engine.Position{X: float64(i * 50), Y: 0})
		world.AddComponent(e, "collisiontag", &CollisionTag{Tag: "enemy"})
		world.AddComponent(e, "boundingbox", &BoundingBox{X: -8, Y: -8, Width: 16, Height: 16})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.Update(1.0 / 60.0)
	}
}
