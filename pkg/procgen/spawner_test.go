package procgen

import (
	"testing"

	"github.com/opd-ai/velocity/pkg/combat"
	"github.com/opd-ai/velocity/pkg/engine"
)

func TestNewWaveSpawner(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)

	if spawner.world != world {
		t.Error("World not set correctly")
	}
	if spawner.generator != gen {
		t.Error("Generator not set correctly")
	}
	if spawner.screenWidth != 800 {
		t.Errorf("Screen width expected 800, got %f", spawner.screenWidth)
	}
	if spawner.screenHeight != 600 {
		t.Errorf("Screen height expected 600, got %f", spawner.screenHeight)
	}
}

func TestWaveSpawner_SpawnWave(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)

	// Wave 1 should spawn 3 enemies (N + 2)
	enemies := spawner.SpawnWave(1)

	if len(enemies) != 3 {
		t.Errorf("Wave 1 expected 3 enemies, got %d", len(enemies))
	}

	// Verify each enemy has required components
	for i, e := range enemies {
		if _, ok := world.GetComponent(e, "position"); !ok {
			t.Errorf("Enemy %d missing position component", i)
		}
		if _, ok := world.GetComponent(e, "velocity"); !ok {
			t.Errorf("Enemy %d missing velocity component", i)
		}
		if _, ok := world.GetComponent(e, "health"); !ok {
			t.Errorf("Enemy %d missing health component", i)
		}
		if _, ok := world.GetComponent(e, "enemy"); !ok {
			t.Errorf("Enemy %d missing enemy AI component", i)
		}
		if _, ok := world.GetComponent(e, "collisiontag"); !ok {
			t.Errorf("Enemy %d missing collision tag", i)
		}
		if _, ok := world.GetComponent(e, "sprite"); !ok {
			t.Errorf("Enemy %d missing sprite component", i)
		}
	}
}

func TestWaveSpawner_SpawnWave_CountFormula(t *testing.T) {
	tests := []struct {
		wave     int
		expected int
	}{
		{1, 3},   // 1 + 2 = 3
		{5, 7},   // 5 + 2 = 7
		{10, 12}, // 10 + 2 = 12
	}

	for _, tt := range tests {
		world := engine.NewWorld()
		gen := NewGenerator(12345)
		spawner := NewWaveSpawner(world, gen, 800, 600)

		enemies := spawner.SpawnWave(tt.wave)
		if len(enemies) != tt.expected {
			t.Errorf("Wave %d: expected %d enemies, got %d", tt.wave, tt.expected, len(enemies))
		}
	}
}

func TestWaveSpawner_EnemyStats(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)

	// Verify wave 1 stats (health = 10 + 1*5 = 15)
	enemies := spawner.SpawnWave(1)
	e := enemies[0]

	healthComp, _ := world.GetComponent(e, "health")
	h := healthComp.(*combat.Health)
	if h.Current != 15.0 {
		t.Errorf("Wave 1 enemy health expected 15, got %f", h.Current)
	}

	aiComp, _ := world.GetComponent(e, "enemy")
	ai := aiComp.(*EnemyAI)
	// Speed = 50 + 1*5 = 55
	if ai.Speed != 55.0 {
		t.Errorf("Wave 1 enemy speed expected 55, got %f", ai.Speed)
	}
}

func TestWaveSpawner_OffscreenPositions(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)

	enemies := spawner.SpawnWave(1)

	for i, e := range enemies {
		posComp, _ := world.GetComponent(e, "position")
		pos := posComp.(*engine.Position)

		// Verify position is outside screen bounds
		isOffscreen := pos.X < 0 || pos.X > 800 || pos.Y < 0 || pos.Y > 600
		if !isOffscreen {
			t.Errorf("Enemy %d at (%f, %f) should be offscreen", i, pos.X, pos.Y)
		}
	}
}

func TestNewEnemyAISystem(t *testing.T) {
	world := engine.NewWorld()
	ais := NewEnemyAISystem(world)

	if ais.world != world {
		t.Error("World not set correctly")
	}
}

func TestEnemyAISystem_SetPlayerEntity(t *testing.T) {
	world := engine.NewWorld()
	ais := NewEnemyAISystem(world)

	player := world.CreateEntity()
	ais.SetPlayerEntity(player)

	if ais.playerEntity != player {
		t.Error("Player entity not set correctly")
	}
}

func TestEnemyAISystem_Update_MovesTowardPlayer(t *testing.T) {
	world := engine.NewWorld()
	ais := NewEnemyAISystem(world)

	// Create player at center
	player := world.CreateEntity()
	world.AddComponent(player, "position", &engine.Position{X: 400, Y: 300})
	ais.SetPlayerEntity(player)

	// Create enemy to the left
	enemy := world.CreateEntity()
	world.AddComponent(enemy, "position", &engine.Position{X: 100, Y: 300})
	world.AddComponent(enemy, "velocity", &engine.Velocity{VX: 0, VY: 0})
	world.AddComponent(enemy, "rotation", &engine.Rotation{Angle: 0})
	world.AddComponent(enemy, "enemy", &EnemyAI{
		State:  EnemyStateApproach,
		Speed:  50.0,
		Damage: 10.0,
	})

	// Update should set velocity toward player
	ais.Update(1.0 / 60.0)

	velComp, _ := world.GetComponent(enemy, "velocity")
	vel := velComp.(*engine.Velocity)

	// Velocity should point right (toward player at x=400)
	if vel.VX <= 0 {
		t.Errorf("Enemy should move toward player (VX > 0), got VX=%f", vel.VX)
	}
}

func TestEnemyAISystem_Update_NoPlayerNoMove(t *testing.T) {
	world := engine.NewWorld()
	ais := NewEnemyAISystem(world)
	// Don't set player entity

	enemy := world.CreateEntity()
	world.AddComponent(enemy, "position", &engine.Position{X: 100, Y: 300})
	world.AddComponent(enemy, "velocity", &engine.Velocity{VX: 0, VY: 0})
	world.AddComponent(enemy, "enemy", &EnemyAI{
		State: EnemyStateApproach,
		Speed: 50.0,
	})

	// Update should not crash and velocity should remain 0
	ais.Update(1.0 / 60.0)

	velComp, _ := world.GetComponent(enemy, "velocity")
	vel := velComp.(*engine.Velocity)

	if vel.VX != 0 || vel.VY != 0 {
		t.Errorf("Velocity should be 0 without player, got VX=%f, VY=%f", vel.VX, vel.VY)
	}
}

func TestEnemyAISystem_CountEnemies(t *testing.T) {
	world := engine.NewWorld()
	ais := NewEnemyAISystem(world)

	// No enemies initially
	if count := ais.CountEnemies(); count != 0 {
		t.Errorf("Expected 0 enemies, got %d", count)
	}

	// Add an enemy
	e1 := world.CreateEntity()
	world.AddComponent(e1, "enemy", &EnemyAI{})
	if count := ais.CountEnemies(); count != 1 {
		t.Errorf("Expected 1 enemy, got %d", count)
	}

	// Add another
	e2 := world.CreateEntity()
	world.AddComponent(e2, "enemy", &EnemyAI{})
	if count := ais.CountEnemies(); count != 2 {
		t.Errorf("Expected 2 enemies, got %d", count)
	}

	// Non-enemy entity shouldn't count
	nonEnemy := world.CreateEntity()
	world.AddComponent(nonEnemy, "position", &engine.Position{})
	if count := ais.CountEnemies(); count != 2 {
		t.Errorf("Expected 2 enemies (non-enemy shouldn't count), got %d", count)
	}
}

func TestEnemyAI_States(t *testing.T) {
	// Verify state constants are distinct
	states := []EnemyState{EnemyStateApproach, EnemyStateAttack, EnemyStateRetreat}
	seen := make(map[EnemyState]bool)

	for _, s := range states {
		if seen[s] {
			t.Errorf("Duplicate state value: %v", s)
		}
		seen[s] = true
	}
}
