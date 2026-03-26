// Package procgen provides procedural content generation systems.
package procgen

import (
	"math"
	"math/rand"

	"github.com/opd-ai/velocity/pkg/combat"
	"github.com/opd-ai/velocity/pkg/engine"
	"github.com/opd-ai/velocity/pkg/rendering"
)

// EnemyConfig describes an enemy to spawn.
type EnemyConfig struct {
	Health float64
	Speed  float64
	Damage float64
}

// WaveSpawner handles spawning enemies for each wave.
type WaveSpawner struct {
	world        *engine.World
	generator    *Generator
	screenWidth  float64
	screenHeight float64
}

// NewWaveSpawner creates a new wave spawner.
func NewWaveSpawner(world *engine.World, generator *Generator, width, height int) *WaveSpawner {
	return &WaveSpawner{
		world:        world,
		generator:    generator,
		screenWidth:  float64(width),
		screenHeight: float64(height),
	}
}

// SpawnWave creates all enemies for the given wave number.
func (ws *WaveSpawner) SpawnWave(waveNumber int) []engine.Entity {
	config := ws.generator.GenerateWave(waveNumber)
	rng := engine.DeterministicRNG(config.Seed)

	enemies := make([]engine.Entity, 0, config.EnemyCount)

	// Calculate enemy stats based on wave formula
	enemyConfig := ws.calculateEnemyStats(waveNumber)

	for i := 0; i < config.EnemyCount; i++ {
		x, y := ws.randomOffscreenPosition(rng)
		e := ws.spawnEnemy(x, y, enemyConfig, i)
		enemies = append(enemies, e)
	}

	return enemies
}

// calculateEnemyStats returns enemy stats for the given wave.
// Formula: health = 10 + wave*5, speed = 1.0 + wave*0.1
func (ws *WaveSpawner) calculateEnemyStats(waveNumber int) EnemyConfig {
	return EnemyConfig{
		Health: 10.0 + float64(waveNumber)*5.0,
		Speed:  50.0 + float64(waveNumber)*5.0, // Scaled for screen units
		Damage: 10.0,
	}
}

// randomOffscreenPosition returns a position just outside the screen.
func (ws *WaveSpawner) randomOffscreenPosition(rng *rand.Rand) (float64, float64) {
	margin := 50.0

	// Pick a random edge: 0=top, 1=right, 2=bottom, 3=left
	edge := rng.Intn(4)

	var x, y float64
	switch edge {
	case 0: // Top
		x = rng.Float64() * ws.screenWidth
		y = -margin
	case 1: // Right
		x = ws.screenWidth + margin
		y = rng.Float64() * ws.screenHeight
	case 2: // Bottom
		x = rng.Float64() * ws.screenWidth
		y = ws.screenHeight + margin
	case 3: // Left
		x = -margin
		y = rng.Float64() * ws.screenHeight
	}

	return x, y
}

// spawnEnemy creates a single enemy entity.
func (ws *WaveSpawner) spawnEnemy(x, y float64, config EnemyConfig, variant int) engine.Entity {
	e := ws.world.CreateEntity()

	ws.world.AddComponent(e, "position", &engine.Position{X: x, Y: y})
	ws.world.AddComponent(e, "velocity", &engine.Velocity{VX: 0, VY: 0})
	ws.world.AddComponent(e, "rotation", &engine.Rotation{Angle: 0})

	ws.world.AddComponent(e, "health", &combat.Health{
		Current: config.Health,
		Max:     config.Health,
	})

	ws.world.AddComponent(e, "collisiontag", &combat.CollisionTag{Tag: "enemy"})
	ws.world.AddComponent(e, "boundingbox", &combat.BoundingBox{
		X: -8, Y: -8, Width: 16, Height: 16,
	})

	ws.world.AddComponent(e, "sprite", &rendering.SpriteComponent{
		Type:    rendering.SpriteTypeEnemy,
		Variant: variant % 5, // 5 enemy variants
		Size:    16,
	})

	ws.world.AddComponent(e, "enemy", &EnemyAI{
		State:  EnemyStateApproach,
		Speed:  config.Speed,
		Damage: config.Damage,
	})

	return e
}

// EnemyState represents the AI state of an enemy.
type EnemyState int

const (
	EnemyStateApproach EnemyState = iota
	EnemyStateAttack
	EnemyStateRetreat
)

// EnemyAI holds enemy behavior data.
type EnemyAI struct {
	State  EnemyState
	Speed  float64
	Damage float64
	Target engine.Entity
}

// EnemyAISystem handles enemy movement and behavior.
type EnemyAISystem struct {
	world        *engine.World
	playerEntity engine.Entity
}

// NewEnemyAISystem creates a new enemy AI system.
func NewEnemyAISystem(world *engine.World) *EnemyAISystem {
	return &EnemyAISystem{world: world}
}

// SetPlayerEntity sets the player entity for AI targeting.
func (ais *EnemyAISystem) SetPlayerEntity(player engine.Entity) {
	ais.playerEntity = player
}

// Update moves all enemies according to their AI state.
func (ais *EnemyAISystem) Update(dt float64) {
	if ais.playerEntity == 0 {
		return
	}

	playerPos, hasPlayerPos := ais.world.GetComponent(ais.playerEntity, "position")
	if !hasPlayerPos {
		return
	}

	targetPos := playerPos.(*engine.Position)

	ais.world.ForEachEntity(func(e engine.Entity) {
		aiComp, hasAI := ais.world.GetComponent(e, "enemy")
		if !hasAI {
			return
		}

		ai := aiComp.(*EnemyAI)
		ais.updateEnemy(e, ai, targetPos, dt)
	})
}

// updateEnemy updates a single enemy's movement.
func (ais *EnemyAISystem) updateEnemy(e engine.Entity, ai *EnemyAI, target *engine.Position, dt float64) {
	posComp, hasPos := ais.world.GetComponent(e, "position")
	velComp, hasVel := ais.world.GetComponent(e, "velocity")
	rotComp, hasRot := ais.world.GetComponent(e, "rotation")

	if !hasPos || !hasVel {
		return
	}

	pos := posComp.(*engine.Position)
	vel := velComp.(*engine.Velocity)

	// Calculate direction to player
	dx := target.X - pos.X
	dy := target.Y - pos.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 0 {
		// Normalize direction
		dx /= dist
		dy /= dist

		// Set velocity toward player
		vel.VX = dx * ai.Speed
		vel.VY = dy * ai.Speed

		// Update rotation to face player
		if hasRot {
			rot := rotComp.(*engine.Rotation)
			rot.Angle = math.Atan2(dy, dx)
		}
	}
}

// CountEnemies returns the number of active enemies.
func (ais *EnemyAISystem) CountEnemies() int {
	count := 0
	ais.world.ForEachEntity(func(e engine.Entity) {
		_, hasAI := ais.world.GetComponent(e, "enemy")
		if hasAI {
			count++
		}
	})
	return count
}
