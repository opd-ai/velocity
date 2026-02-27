// Package companion provides wingman AI behaviour trees.
package companion

// Wingman represents an AI-controlled companion ship.
type Wingman struct {
	genreID string
	active  bool
}

// NewWingman creates a new wingman AI.
func NewWingman() *Wingman {
	return &Wingman{genreID: "scifi"}
}

// SetGenre switches the wingman's visual style to match the given genre.
func (w *Wingman) SetGenre(genreID string) {
	w.genreID = genreID
}

// Activate enables the wingman.
func (w *Wingman) Activate() {
	w.active = true
}

// Deactivate disables the wingman.
func (w *Wingman) Deactivate() {
	w.active = false
}

// Update advances the wingman AI by dt seconds.
func (w *Wingman) Update(dt float64) {
	// Stub: will run behaviour tree for escort, evade, and attack patterns.
}
