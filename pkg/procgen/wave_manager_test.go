package procgen

import (
	"testing"

	"github.com/opd-ai/velocity/pkg/engine"
)

func TestNewWaveManager(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	if wm.world != world {
		t.Error("World not set correctly")
	}
	if wm.currentWave != 0 {
		t.Errorf("Initial wave should be 0, got %d", wm.currentWave)
	}
	if wm.WaveInProgress() {
		t.Error("Wave should not be in progress initially")
	}
}

func TestWaveManager_CurrentWave(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	if wm.CurrentWave() != 0 {
		t.Errorf("Expected wave 0, got %d", wm.CurrentWave())
	}
}

func TestWaveManager_StartNextWave(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	wm.StartNextWave()

	if wm.CurrentWave() != 1 {
		t.Errorf("Expected wave 1 after starting, got %d", wm.CurrentWave())
	}
	if !wm.WaveInProgress() {
		t.Error("Wave should be in progress after starting")
	}

	// Should have spawned 3 enemies (wave 1: 1+2)
	count := ai.CountEnemies()
	if count != 3 {
		t.Errorf("Expected 3 enemies, got %d", count)
	}
}

func TestWaveManager_WaveStartCallback(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	startedWave := 0
	wm.SetWaveStartCallback(func(wave int) {
		startedWave = wave
	})

	wm.StartNextWave()

	if startedWave != 1 {
		t.Errorf("Callback should receive wave 1, got %d", startedWave)
	}
}

func TestWaveManager_WaveCompleteCallback(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	completedWave := 0
	wm.SetWaveCompleteCallback(func(wave int) {
		completedWave = wave
	})

	wm.StartNextWave()

	// Remove all enemies to complete the wave
	removeAllEnemies(world)

	wm.Update(1.0 / 60.0)

	if completedWave != 1 {
		t.Errorf("Callback should receive wave 1, got %d", completedWave)
	}
	if wm.WaveInProgress() {
		t.Error("Wave should no longer be in progress")
	}
}

func TestWaveManager_OnEnemyKilled(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	if wm.TotalKills() != 0 {
		t.Errorf("Initial kills should be 0, got %d", wm.TotalKills())
	}

	wm.OnEnemyKilled()
	wm.OnEnemyKilled()

	if wm.TotalKills() != 2 {
		t.Errorf("Expected 2 kills, got %d", wm.TotalKills())
	}
}

func TestWaveManager_DifficultyMultiplier(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	// Wave 0: multiplier should be 0.9 (1.0 + (0-1)*0.1)
	// Wave 1: multiplier should be 1.0
	// Wave 5: multiplier should be 1.4

	wm.StartNextWave() // Wave 1
	if wm.DifficultyMultiplier() != 1.0 {
		t.Errorf("Wave 1 multiplier expected 1.0, got %f", wm.DifficultyMultiplier())
	}

	// Simulate completing waves and starting new ones
	removeAllEnemies(world)
	wm.Update(1.0 / 60.0)
	wm.StartNextWave() // Wave 2

	expected := 1.1
	actual := wm.DifficultyMultiplier()
	if actual != expected {
		t.Errorf("Wave 2 multiplier expected %f, got %f", expected, actual)
	}
}

func TestWaveManager_GetWaveStats(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	wm.StartNextWave()

	stats := wm.GetWaveStats()

	if stats.CurrentWave != 1 {
		t.Errorf("Expected current wave 1, got %d", stats.CurrentWave)
	}
	if stats.EnemyCount != 3 {
		t.Errorf("Expected 3 enemies, got %d", stats.EnemyCount)
	}
	// Health for wave 1: 10 + 1*5 = 15
	if stats.ExpectedHealth != 15.0 {
		t.Errorf("Expected health 15, got %f", stats.ExpectedHealth)
	}
	// Speed for wave 1: 50 + 1*5 = 55
	if stats.ExpectedSpeed != 55.0 {
		t.Errorf("Expected speed 55, got %f", stats.ExpectedSpeed)
	}
}

func TestWaveManager_Reset(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	wm.StartNextWave()
	wm.OnEnemyKilled()

	wm.Reset()

	if wm.CurrentWave() != 0 {
		t.Errorf("Wave should reset to 0, got %d", wm.CurrentWave())
	}
	if wm.TotalKills() != 0 {
		t.Errorf("Kills should reset to 0, got %d", wm.TotalKills())
	}
	if wm.WaveInProgress() {
		t.Error("Wave should not be in progress after reset")
	}
}

func TestWaveManager_MultipleWaves(t *testing.T) {
	world := engine.NewWorld()
	gen := NewGenerator(12345)
	spawner := NewWaveSpawner(world, gen, 800, 600)
	ai := NewEnemyAISystem(world)
	wm := NewWaveManager(world, spawner, ai)

	// Progress through 3 waves
	for wave := 1; wave <= 3; wave++ {
		wm.StartNextWave()

		if wm.CurrentWave() != wave {
			t.Errorf("Expected wave %d, got %d", wave, wm.CurrentWave())
		}

		expectedEnemies := wave + 2
		if ai.CountEnemies() != expectedEnemies {
			t.Errorf("Wave %d: expected %d enemies, got %d", wave, expectedEnemies, ai.CountEnemies())
		}

		removeAllEnemies(world)
		wm.Update(1.0 / 60.0)
	}
}

// removeAllEnemies helper function to clear enemies from the world.
func removeAllEnemies(world *engine.World) {
	toRemove := []engine.Entity{}
	world.ForEachEntity(func(e engine.Entity) {
		if _, ok := world.GetComponent(e, "enemy"); ok {
			toRemove = append(toRemove, e)
		}
	})
	for _, e := range toRemove {
		world.RemoveEntity(e)
	}
}
