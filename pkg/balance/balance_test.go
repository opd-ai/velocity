package balance

import "testing"

func TestDifficultyScale(t *testing.T) {
	tests := []struct {
		wave     int
		expected float64
	}{
		{0, 1.0},
		{1, 1.1},
		{5, 1.5},
		{10, 2.0},
		{20, 3.0},
		{100, 11.0},
	}

	for _, tt := range tests {
		got := DifficultyScale(tt.wave)
		if got != tt.expected {
			t.Errorf("DifficultyScale(%d) = %f, want %f", tt.wave, got, tt.expected)
		}
	}
}

func TestDefaultPlayerStats(t *testing.T) {
	stats := DefaultPlayerStats()

	if stats.Name != "player" {
		t.Errorf("Name = %q, want %q", stats.Name, "player")
	}
	if stats.Health != 100 {
		t.Errorf("Health = %f, want 100", stats.Health)
	}
	if stats.Speed != 200 {
		t.Errorf("Speed = %f, want 200", stats.Speed)
	}
	if stats.Damage != 10 {
		t.Errorf("Damage = %f, want 10", stats.Damage)
	}
	if stats.Defense != 5 {
		t.Errorf("Defense = %f, want 5", stats.Defense)
	}
}

func TestDefaultEnemyStats(t *testing.T) {
	stats := DefaultEnemyStats()

	if stats.Name != "enemy" {
		t.Errorf("Name = %q, want %q", stats.Name, "enemy")
	}
	if stats.Health != 30 {
		t.Errorf("Health = %f, want 30", stats.Health)
	}
	if stats.Speed != 100 {
		t.Errorf("Speed = %f, want 100", stats.Speed)
	}
	if stats.Damage != 5 {
		t.Errorf("Damage = %f, want 5", stats.Damage)
	}
	if stats.Defense != 1 {
		t.Errorf("Defense = %f, want 1", stats.Defense)
	}
}

func TestStatTableStruct(t *testing.T) {
	st := StatTable{
		Name:    "custom",
		Health:  50,
		Speed:   150,
		Damage:  15,
		Defense: 3,
	}

	if st.Name != "custom" {
		t.Errorf("Name = %q, want %q", st.Name, "custom")
	}
	if st.Health != 50 {
		t.Errorf("Health = %f, want 50", st.Health)
	}
	if st.Speed != 150 {
		t.Errorf("Speed = %f, want 150", st.Speed)
	}
	if st.Damage != 15 {
		t.Errorf("Damage = %f, want 15", st.Damage)
	}
	if st.Defense != 3 {
		t.Errorf("Defense = %f, want 3", st.Defense)
	}
}

func TestDifficultyScaleNegativeWave(t *testing.T) {
	// Negative waves should still work (though unusual)
	got := DifficultyScale(-5)
	expected := 0.5 // 1.0 + (-5)*0.1

	if got != expected {
		t.Errorf("DifficultyScale(-5) = %f, want %f", got, expected)
	}
}

func TestPlayerStatsReasonableValues(t *testing.T) {
	stats := DefaultPlayerStats()

	// Player should be stronger than enemies
	enemy := DefaultEnemyStats()

	if stats.Health <= enemy.Health {
		t.Error("Player health should exceed enemy health")
	}
	if stats.Speed <= enemy.Speed {
		t.Error("Player speed should exceed enemy speed")
	}
	if stats.Damage <= enemy.Damage {
		t.Error("Player damage should exceed enemy damage")
	}
}
