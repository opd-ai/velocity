// Package procgen provides procedural content generation systems.
package procgen

// WaveConfig describes a single procedural wave of enemies.
type WaveConfig struct {
	WaveNumber int
	EnemyCount int
	Seed       int64
}

// Generator produces procedural content from a seed.
type Generator struct {
	seed    int64
	genreID string
}

// NewGenerator creates a new procedural content generator.
func NewGenerator(seed int64) *Generator {
	return &Generator{seed: seed, genreID: "scifi"}
}

// SetGenre switches procedural generation to match the given genre.
func (g *Generator) SetGenre(genreID string) {
	g.genreID = genreID
}

// GenerateWave produces a wave configuration for the given wave number.
func (g *Generator) GenerateWave(waveNumber int) WaveConfig {
	return WaveConfig{
		WaveNumber: waveNumber,
		EnemyCount: waveNumber + 2,
		Seed:       g.seed + int64(waveNumber),
	}
}
