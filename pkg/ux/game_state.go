// Package ux provides the menu framework, HUD components, and tutorial scaffolding.
package ux

// GameState represents the overall game state.
type GameState int

const (
	// StateMainMenu shows the main menu
	StateMainMenu GameState = iota
	// StatePlaying is the active gameplay state
	StatePlaying
	// StatePaused shows the pause menu
	StatePaused
	// StateGameOver shows the game over screen
	StateGameOver
)

// GameStateManager handles state transitions and game flow.
type GameStateManager struct {
	state         GameState
	previousState GameState
	finalScore    int64
	finalWave     int
	onStateChange func(from, to GameState)
}

// NewGameStateManager creates a new game state manager.
func NewGameStateManager() *GameStateManager {
	return &GameStateManager{
		state: StateMainMenu,
	}
}

// State returns the current game state.
func (gsm *GameStateManager) State() GameState {
	return gsm.state
}

// SetStateChangeCallback sets the callback for state transitions.
func (gsm *GameStateManager) SetStateChangeCallback(fn func(from, to GameState)) {
	gsm.onStateChange = fn
}

// SetFinalScore sets the score to display on game over.
func (gsm *GameStateManager) SetFinalScore(score int64) {
	gsm.finalScore = score
}

// SetFinalWave sets the wave reached for game over display.
func (gsm *GameStateManager) SetFinalWave(wave int) {
	gsm.finalWave = wave
}

// FinalScore returns the score achieved when game ended.
func (gsm *GameStateManager) FinalScore() int64 {
	return gsm.finalScore
}

// FinalWave returns the wave reached when game ended.
func (gsm *GameStateManager) FinalWave() int {
	return gsm.finalWave
}

// transition changes to a new state with callback.
func (gsm *GameStateManager) transition(to GameState) {
	if gsm.state == to {
		return
	}
	from := gsm.state
	gsm.previousState = from
	gsm.state = to
	if gsm.onStateChange != nil {
		gsm.onStateChange(from, to)
	}
}

// StartGame transitions from main menu to playing.
func (gsm *GameStateManager) StartGame() {
	if gsm.state == StateMainMenu {
		gsm.transition(StatePlaying)
	}
}

// PauseGame transitions from playing to paused.
func (gsm *GameStateManager) PauseGame() {
	if gsm.state == StatePlaying {
		gsm.transition(StatePaused)
	}
}

// ResumeGame transitions from paused back to playing.
func (gsm *GameStateManager) ResumeGame() {
	if gsm.state == StatePaused {
		gsm.transition(StatePlaying)
	}
}

// GameOver transitions to the game over state.
func (gsm *GameStateManager) GameOver(score int64, wave int) {
	if gsm.state == StatePlaying {
		gsm.finalScore = score
		gsm.finalWave = wave
		gsm.transition(StateGameOver)
	}
}

// ReturnToMainMenu transitions to main menu from any state.
func (gsm *GameStateManager) ReturnToMainMenu() {
	gsm.transition(StateMainMenu)
}

// IsPlaying returns true if game is in playing state.
func (gsm *GameStateManager) IsPlaying() bool {
	return gsm.state == StatePlaying
}

// IsPaused returns true if game is paused.
func (gsm *GameStateManager) IsPaused() bool {
	return gsm.state == StatePaused
}

// IsMenuActive returns true if any menu is showing.
func (gsm *GameStateManager) IsMenuActive() bool {
	return gsm.state != StatePlaying
}

// MenuItem represents a selectable menu option.
type MenuItem struct {
	Label    string
	Action   string
	Selected bool
}

// MenuItems holds the available options for each menu state.
type MenuItems struct {
	MainMenu  []MenuItem
	PauseMenu []MenuItem
	GameOver  []MenuItem
}

// DefaultMenuItems returns the default menu configuration.
func DefaultMenuItems() MenuItems {
	return MenuItems{
		MainMenu: []MenuItem{
			{Label: "Start", Action: "start"},
			{Label: "Settings", Action: "settings"},
			{Label: "Quit", Action: "quit"},
		},
		PauseMenu: []MenuItem{
			{Label: "Resume", Action: "resume"},
			{Label: "Settings", Action: "settings"},
			{Label: "Quit to Menu", Action: "quit_menu"},
		},
		GameOver: []MenuItem{
			{Label: "Retry", Action: "retry"},
			{Label: "Main Menu", Action: "main_menu"},
		},
	}
}

// MenuController handles menu navigation and selection.
type MenuController struct {
	items       MenuItems
	selectedIdx int
	stateManager *GameStateManager
	onAction    func(action string)
}

// NewMenuController creates a new menu controller.
func NewMenuController(stateManager *GameStateManager) *MenuController {
	return &MenuController{
		items:        DefaultMenuItems(),
		selectedIdx:  0,
		stateManager: stateManager,
	}
}

// SetActionCallback sets the callback for menu actions.
func (mc *MenuController) SetActionCallback(fn func(action string)) {
	mc.onAction = fn
}

// GetCurrentItems returns the menu items for the current state.
func (mc *MenuController) GetCurrentItems() []MenuItem {
	switch mc.stateManager.State() {
	case StateMainMenu:
		return mc.items.MainMenu
	case StatePaused:
		return mc.items.PauseMenu
	case StateGameOver:
		return mc.items.GameOver
	default:
		return nil
	}
}

// SelectionIndex returns the currently selected item index.
func (mc *MenuController) SelectionIndex() int {
	return mc.selectedIdx
}

// MoveUp moves selection up (wrapping).
func (mc *MenuController) MoveUp() {
	items := mc.GetCurrentItems()
	if len(items) == 0 {
		return
	}
	mc.selectedIdx--
	if mc.selectedIdx < 0 {
		mc.selectedIdx = len(items) - 1
	}
}

// MoveDown moves selection down (wrapping).
func (mc *MenuController) MoveDown() {
	items := mc.GetCurrentItems()
	if len(items) == 0 {
		return
	}
	mc.selectedIdx++
	if mc.selectedIdx >= len(items) {
		mc.selectedIdx = 0
	}
}

// Select executes the currently selected menu item.
func (mc *MenuController) Select() {
	items := mc.GetCurrentItems()
	if len(items) == 0 || mc.selectedIdx >= len(items) {
		return
	}

	action := items[mc.selectedIdx].Action
	mc.handleAction(action)
}

// handleAction processes a menu action.
func (mc *MenuController) handleAction(action string) {
	switch action {
	case "start":
		mc.stateManager.StartGame()
		mc.selectedIdx = 0
	case "resume":
		mc.stateManager.ResumeGame()
	case "retry":
		mc.stateManager.StartGame()
		mc.selectedIdx = 0
	case "main_menu", "quit_menu":
		mc.stateManager.ReturnToMainMenu()
		mc.selectedIdx = 0
	}

	if mc.onAction != nil {
		mc.onAction(action)
	}
}

// ResetSelection resets the selection to the first item.
func (mc *MenuController) ResetSelection() {
	mc.selectedIdx = 0
}

// AddContinueOption adds a "Continue" option to the main menu.
func (mc *MenuController) AddContinueOption() {
	continueItem := MenuItem{Label: "Continue", Action: "continue"}
	// Insert at the beginning of the main menu
	mc.items.MainMenu = append([]MenuItem{continueItem}, mc.items.MainMenu...)
}
