// Package ux provides the menu framework, HUD components, and tutorial scaffolding.
package ux

// HUD holds the state for the heads-up display.
type HUD struct {
	genreID string
	Health  float64
	Shield  float64
	Score   int64
	Wave    int
	Combo   int
}

// NewHUD creates a new HUD instance.
func NewHUD() *HUD {
	return &HUD{genreID: "scifi"}
}

// SetGenre switches the HUD visual style to match the given genre.
func (h *HUD) SetGenre(genreID string) {
	h.genreID = genreID
}

// Update refreshes HUD values.
func (h *HUD) Update(health, shield float64, score int64, wave, combo int) {
	h.Health = health
	h.Shield = shield
	h.Score = score
	h.Wave = wave
	h.Combo = combo
}

// MenuState represents the current menu screen.
type MenuState int

const (
	MenuNone MenuState = iota
	MenuMain
	MenuPause
	MenuGameOver
	MenuHighScore
)

// Menu manages menu navigation and display.
type Menu struct {
	State   MenuState
	genreID string
}

// NewMenu creates a new menu manager.
func NewMenu() *Menu {
	return &Menu{State: MenuMain, genreID: "scifi"}
}

// SetGenre switches the menu visual style to match the given genre.
func (m *Menu) SetGenre(genreID string) {
	m.genreID = genreID
}

// Tutorial manages the first-run tutorial sequence.
type Tutorial struct {
	Active bool
	Step   int
}

// NewTutorial creates a new tutorial.
func NewTutorial() *Tutorial {
	return &Tutorial{Active: true, Step: 0}
}

// Advance moves the tutorial to the next step.
func (t *Tutorial) Advance() {
	t.Step++
}

// Complete marks the tutorial as finished.
func (t *Tutorial) Complete() {
	t.Active = false
}
