package procgen

import (
	"testing"

	"github.com/opd-ai/velocity/pkg/procgen/genre"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator(12345)

	if g == nil {
		t.Fatal("expected non-nil Generator")
	}
}

func TestGenerator_GenerateWave_Formula(t *testing.T) {
	g := NewGenerator(12345)

	tests := []struct {
		waveNumber int
		wantCount  int
	}{
		{1, 3},  // 1 + 2 = 3
		{5, 7},  // 5 + 2 = 7
		{10, 12}, // 10 + 2 = 12
		{0, 2},  // 0 + 2 = 2
	}

	for _, tt := range tests {
		t.Run(string(rune('0'+tt.waveNumber)), func(t *testing.T) {
			config := g.GenerateWave(tt.waveNumber)

			if config.EnemyCount != tt.wantCount {
				t.Errorf("wave %d: expected EnemyCount %d, got %d",
					tt.waveNumber, tt.wantCount, config.EnemyCount)
			}

			if config.WaveNumber != tt.waveNumber {
				t.Errorf("expected WaveNumber %d, got %d", tt.waveNumber, config.WaveNumber)
			}
		})
	}
}

func TestGenerator_SetGenre(t *testing.T) {
	g := NewGenerator(12345)

	for _, genreID := range genre.All() {
		t.Run(genreID, func(t *testing.T) {
			g.SetGenre(genreID)
			// Should not panic
		})
	}
}

func TestGenerator_Determinism(t *testing.T) {
	seed := int64(42)

	g1 := NewGenerator(seed)
	g2 := NewGenerator(seed)

	config1 := g1.GenerateWave(5)
	config2 := g2.GenerateWave(5)

	if config1.EnemyCount != config2.EnemyCount {
		t.Error("expected deterministic enemy count")
	}

	if config1.Seed != config2.Seed {
		t.Error("expected deterministic wave seed")
	}
}

func TestGenerator_WaveSeedDiffers(t *testing.T) {
	g := NewGenerator(12345)

	config1 := g.GenerateWave(1)
	config2 := g.GenerateWave(2)

	if config1.Seed == config2.Seed {
		t.Error("expected different seeds for different waves")
	}
}

func TestWaveConfig(t *testing.T) {
	config := WaveConfig{
		WaveNumber: 5,
		EnemyCount: 7,
		Seed:       12345,
	}

	if config.WaveNumber != 5 {
		t.Errorf("expected WaveNumber 5, got %d", config.WaveNumber)
	}
	if config.EnemyCount != 7 {
		t.Errorf("expected EnemyCount 7, got %d", config.EnemyCount)
	}
	if config.Seed != 12345 {
		t.Errorf("expected Seed 12345, got %d", config.Seed)
	}
}
