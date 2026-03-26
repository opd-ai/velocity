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
	Active        bool
	Step          int
	stepsComplete map[string]bool
	prompts       []TutorialPrompt
}

// TutorialPrompt describes a tutorial instruction.
type TutorialPrompt struct {
	Step    int
	Action  string // "thrust", "rotate", "fire", "kill"
	Text    string
	KeyHint string
}

// TutorialStepThrust is the first tutorial step.
const TutorialStepThrust = 0

// TutorialStepRotate is the second tutorial step.
const TutorialStepRotate = 1

// TutorialStepFire is the third tutorial step.
const TutorialStepFire = 2

// TutorialStepKill is the final tutorial step.
const TutorialStepKill = 3

// NewTutorial creates a new tutorial.
func NewTutorial() *Tutorial {
	return &Tutorial{
		Active:        true,
		Step:          TutorialStepThrust,
		stepsComplete: make(map[string]bool),
		prompts: []TutorialPrompt{
			{Step: TutorialStepThrust, Action: "thrust", Text: "THRUST FORWARD", KeyHint: "W / Left Stick Up"},
			{Step: TutorialStepRotate, Action: "rotate", Text: "ROTATE YOUR SHIP", KeyHint: "A/D / Left Stick"},
			{Step: TutorialStepFire, Action: "fire", Text: "FIRE WEAPONS", KeyHint: "SPACE / A Button"},
			{Step: TutorialStepKill, Action: "kill", Text: "DESTROY AN ENEMY", KeyHint: ""},
		},
	}
}

// Advance moves the tutorial to the next step.
func (t *Tutorial) Advance() {
	t.Step++
	if t.Step >= len(t.prompts) {
		t.Complete()
	}
}

// Complete marks the tutorial as finished.
func (t *Tutorial) Complete() {
	t.Active = false
}

// MarkAction records that the player performed an action.
func (t *Tutorial) MarkAction(action string) bool {
	if !t.Active {
		return false
	}

	// Check if this action completes the current step
	if t.Step < len(t.prompts) && t.prompts[t.Step].Action == action {
		t.stepsComplete[action] = true
		t.Advance()
		return true
	}

	return false
}

// CurrentPrompt returns the current tutorial prompt.
func (t *Tutorial) CurrentPrompt() *TutorialPrompt {
	if !t.Active || t.Step >= len(t.prompts) {
		return nil
	}
	return &t.prompts[t.Step]
}

// Progress returns the tutorial progress as a fraction (0.0 to 1.0).
func (t *Tutorial) Progress() float64 {
	if len(t.prompts) == 0 {
		return 1.0
	}
	return float64(t.Step) / float64(len(t.prompts))
}
