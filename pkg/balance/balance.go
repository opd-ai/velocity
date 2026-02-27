// Package balance provides stat tuning tables and difficulty curves.
package balance

// StatTable holds stat values for a named entity type.
type StatTable struct {
	Name    string
	Health  float64
	Speed   float64
	Damage  float64
	Defense float64
}

// DifficultyScale returns a multiplier for the given wave number.
func DifficultyScale(wave int) float64 {
	return 1.0 + float64(wave)*0.1
}

// DefaultPlayerStats returns the base player stat table.
func DefaultPlayerStats() StatTable {
	return StatTable{
		Name:    "player",
		Health:  100,
		Speed:   200,
		Damage:  10,
		Defense: 5,
	}
}

// DefaultEnemyStats returns the base enemy stat table.
func DefaultEnemyStats() StatTable {
	return StatTable{
		Name:    "enemy",
		Health:  30,
		Speed:   100,
		Damage:  5,
		Defense: 1,
	}
}
