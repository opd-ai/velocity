package combat

import (
	"testing"

	"github.com/opd-ai/velocity/pkg/engine"
)

func TestNewDamageSystem(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	if ds == nil {
		t.Fatal("expected non-nil DamageSystem")
	}
}

func TestDamageSystem_ApplyDamage(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "health", &Health{Current: 100, Max: 100})

	died := ds.ApplyDamage(entity, 30)

	if died {
		t.Error("expected entity to survive 30 damage")
	}

	healthComp, _ := world.GetComponent(entity, "health")
	health := healthComp.(*Health)

	if health.Current != 70 {
		t.Errorf("expected 70 health, got %f", health.Current)
	}
}

func TestDamageSystem_ApplyDamage_Kill(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "health", &Health{Current: 50, Max: 100})

	died := ds.ApplyDamage(entity, 100)

	if !died {
		t.Error("expected entity to die from 100 damage")
	}

	healthComp, _ := world.GetComponent(entity, "health")
	health := healthComp.(*Health)

	if health.Current != 0 {
		t.Errorf("expected 0 health, got %f", health.Current)
	}
}

func TestDamageSystem_ApplyDamage_NoHealth(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	// No health component

	died := ds.ApplyDamage(entity, 100)

	if died {
		t.Error("expected no death without health component")
	}
}

func TestDamageSystem_QueueDamage(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "health", &Health{Current: 100, Max: 100})

	ds.QueueDamage(entity, 0, 25, "projectile")
	ds.QueueDamage(entity, 0, 25, "projectile")

	// Damage not applied yet
	healthComp, _ := world.GetComponent(entity, "health")
	health := healthComp.(*Health)
	if health.Current != 100 {
		t.Error("expected damage to be queued, not applied")
	}

	// Process damage
	ds.Update(0)

	healthComp, _ = world.GetComponent(entity, "health")
	health = healthComp.(*Health)
	if health.Current != 50 {
		t.Errorf("expected 50 health after queued damage, got %f", health.Current)
	}
}

func TestDamageSystem_DeathCallback(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "position", &engine.Position{X: 100, Y: 200})
	world.AddComponent(entity, "health", &Health{Current: 10, Max: 100})
	world.AddComponent(entity, "collisiontag", &CollisionTag{Tag: "enemy"})

	deathCount := 0
	var lastDeath DeathEvent
	ds.SetDeathCallback(func(event DeathEvent) {
		deathCount++
		lastDeath = event
	})

	ds.QueueDamage(entity, 0, 100, "projectile")
	ds.Update(0)

	if deathCount != 1 {
		t.Errorf("expected 1 death callback, got %d", deathCount)
	}

	if lastDeath.Position.X != 100 || lastDeath.Position.Y != 200 {
		t.Error("expected death position to match entity position")
	}

	if lastDeath.Score != 100 {
		t.Errorf("expected score 100, got %d", lastDeath.Score)
	}
}

func TestDamageSystem_EntityRemoved(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "health", &Health{Current: 10, Max: 100})

	ds.QueueDamage(entity, 0, 100, "projectile")
	ds.Update(0)

	// Entity should be removed
	_, hasHealth := world.GetComponent(entity, "health")
	if hasHealth {
		t.Error("expected entity to be removed after death")
	}
}

func TestDamageSystem_IsEntityAlive(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "health", &Health{Current: 50, Max: 100})

	if !ds.IsEntityAlive(entity) {
		t.Error("expected entity to be alive")
	}

	ds.ApplyDamage(entity, 50)

	if ds.IsEntityAlive(entity) {
		t.Error("expected entity to be dead")
	}
}

func TestDamageSystem_IsEntityAlive_NoHealth(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	// No health component = invulnerable

	if !ds.IsEntityAlive(entity) {
		t.Error("expected entity without health to be considered alive")
	}
}

func TestDamageSystem_GetHealthPercent(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "health", &Health{Current: 75, Max: 100})

	percent := ds.GetHealthPercent(entity)

	if percent != 0.75 {
		t.Errorf("expected 0.75 health percent, got %f", percent)
	}
}

func TestDamageSystem_GetHealthPercent_NoHealth(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()

	percent := ds.GetHealthPercent(entity)

	if percent != 1.0 {
		t.Errorf("expected 1.0 for entity without health, got %f", percent)
	}
}

func TestDamageSystem_Heal(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "health", &Health{Current: 50, Max: 100})

	ds.Heal(entity, 30)

	healthComp, _ := world.GetComponent(entity, "health")
	health := healthComp.(*Health)

	if health.Current != 80 {
		t.Errorf("expected 80 health after heal, got %f", health.Current)
	}
}

func TestDamageSystem_Heal_Cap(t *testing.T) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	entity := world.CreateEntity()
	world.AddComponent(entity, "health", &Health{Current: 90, Max: 100})

	ds.Heal(entity, 50)

	healthComp, _ := world.GetComponent(entity, "health")
	health := healthComp.(*Health)

	if health.Current != 100 {
		t.Errorf("expected health capped at 100, got %f", health.Current)
	}
}

func TestDamageEvent_Fields(t *testing.T) {
	event := DamageEvent{
		Target:     engine.Entity(1),
		Amount:     50,
		Source:     engine.Entity(2),
		SourceType: "explosion",
	}

	if event.Amount != 50 {
		t.Error("amount mismatch")
	}
	if event.SourceType != "explosion" {
		t.Error("source type mismatch")
	}
}

func TestDeathEvent_Fields(t *testing.T) {
	event := DeathEvent{
		Entity:   engine.Entity(1),
		Position: engine.Position{X: 100, Y: 200},
		Score:    150,
	}

	if event.Score != 150 {
		t.Error("score mismatch")
	}
	if event.Position.X != 100 {
		t.Error("position mismatch")
	}
}

func BenchmarkDamageSystem_Update(b *testing.B) {
	world := engine.NewWorld()
	ds := NewDamageSystem(world)

	// Create many entities
	for i := 0; i < 100; i++ {
		e := world.CreateEntity()
		world.AddComponent(e, "health", &Health{Current: 100, Max: 100})
		world.AddComponent(e, "position", &engine.Position{X: float64(i * 10), Y: 0})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Queue some damage
		ds.QueueDamage(engine.Entity(1), 0, 1, "test")
		ds.Update(1.0 / 60.0)
	}
}
