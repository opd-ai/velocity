// Package procgen provides procedural content generation systems.
package procgen

import "github.com/opd-ai/velocity/pkg/engine"

// WaveManager handles wave progression and difficulty ramping.
type WaveManager struct {
	world        *engine.World
	spawner      *WaveSpawner
	aiSystem     *EnemyAISystem
	currentWave  int
	totalKills   int
	waveKills    int
	waveInProgress bool
	onWaveComplete func(waveNumber int)
	onWaveStart    func(waveNumber int)
}

// NewWaveManager creates a new wave manager.
func NewWaveManager(world *engine.World, spawner *WaveSpawner, aiSystem *EnemyAISystem) *WaveManager {
	return &WaveManager{
		world:       world,
		spawner:     spawner,
		aiSystem:    aiSystem,
		currentWave: 0,
	}
}

// SetWaveCompleteCallback sets the callback for wave completion.
func (wm *WaveManager) SetWaveCompleteCallback(fn func(waveNumber int)) {
	wm.onWaveComplete = fn
}

// SetWaveStartCallback sets the callback for wave start.
func (wm *WaveManager) SetWaveStartCallback(fn func(waveNumber int)) {
	wm.onWaveStart = fn
}

// CurrentWave returns the current wave number.
func (wm *WaveManager) CurrentWave() int {
	return wm.currentWave
}

// TotalKills returns the total number of enemies killed.
func (wm *WaveManager) TotalKills() int {
	return wm.totalKills
}

// WaveInProgress returns true if a wave is currently active.
func (wm *WaveManager) WaveInProgress() bool {
	return wm.waveInProgress
}

// StartNextWave begins the next wave.
func (wm *WaveManager) StartNextWave() {
	wm.currentWave++
	wm.waveKills = 0
	wm.waveInProgress = true

	wm.spawner.SpawnWave(wm.currentWave)

	if wm.onWaveStart != nil {
		wm.onWaveStart(wm.currentWave)
	}
}

// OnEnemyKilled should be called when an enemy is destroyed.
func (wm *WaveManager) OnEnemyKilled() {
	wm.totalKills++
	wm.waveKills++
}

// Update checks wave completion and manages transitions.
func (wm *WaveManager) Update(dt float64) {
	if !wm.waveInProgress {
		return
	}

	// Check if all enemies are destroyed
	enemyCount := wm.aiSystem.CountEnemies()
	if enemyCount == 0 {
		wm.waveInProgress = false

		if wm.onWaveComplete != nil {
			wm.onWaveComplete(wm.currentWave)
		}
	}
}

// DifficultyMultiplier returns a scaling factor for the current wave.
func (wm *WaveManager) DifficultyMultiplier() float64 {
	// Linear difficulty ramp: 1.0 at wave 1, increasing 0.1 per wave
	return 1.0 + float64(wm.currentWave-1)*0.1
}

// WaveStats returns statistics about wave progression.
type WaveStats struct {
	CurrentWave    int
	EnemyCount     int
	ExpectedHealth float64
	ExpectedSpeed  float64
	TotalKills     int
}

// GetWaveStats returns current wave statistics.
func (wm *WaveManager) GetWaveStats() WaveStats {
	return WaveStats{
		CurrentWave:    wm.currentWave,
		EnemyCount:     wm.aiSystem.CountEnemies(),
		ExpectedHealth: 10.0 + float64(wm.currentWave)*5.0,
		ExpectedSpeed:  50.0 + float64(wm.currentWave)*5.0,
		TotalKills:     wm.totalKills,
	}
}

// Reset resets the wave manager to initial state.
func (wm *WaveManager) Reset() {
	wm.currentWave = 0
	wm.totalKills = 0
	wm.waveKills = 0
	wm.waveInProgress = false
}
