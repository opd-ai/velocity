// Package world provides quest objectives, weather systems, environmental
// storytelling, loot, economy, and mini-game subsystems.
package world

// Weather represents a space weather phenomenon.
type Weather struct {
	Name       string
	GenreID    string
	Active     bool
	Intensity  float64
}

// NewWeather creates a new weather instance.
func NewWeather(name, genreID string) *Weather {
	return &Weather{Name: name, GenreID: genreID}
}

// SetGenre switches weather visuals to match the given genre.
func (w *Weather) SetGenre(genreID string) {
	w.GenreID = genreID
}

// Objective represents a wave objective.
type Objective struct {
	Description string
	Completed   bool
}

// NewObjective creates a new wave objective.
func NewObjective(description string) *Objective {
	return &Objective{Description: description}
}

// Complete marks the objective as finished.
func (o *Objective) Complete() {
	o.Completed = true
}

// LootDrop represents a pickup dropped by a destroyed enemy.
type LootDrop struct {
	Type string
	X, Y float64
}

// Economy manages the between-wave shop.
type Economy struct {
	Credits int64
}

// NewEconomy creates a new economy with zero credits.
func NewEconomy() *Economy {
	return &Economy{}
}

// AddCredits adds credits to the player's balance.
func (e *Economy) AddCredits(amount int64) {
	e.Credits += amount
}

// Spend deducts credits if sufficient balance exists.
func (e *Economy) Spend(amount int64) bool {
	if e.Credits >= amount {
		e.Credits -= amount
		return true
	}
	return false
}
