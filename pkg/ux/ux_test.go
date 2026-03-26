package ux

import "testing"

func TestNewGameStateManager(t *testing.T) {
	gsm := NewGameStateManager()

	if gsm.State() != StateMainMenu {
		t.Errorf("Initial state should be StateMainMenu, got %v", gsm.State())
	}
}

func TestGameStateManager_StartGame(t *testing.T) {
	gsm := NewGameStateManager()

	gsm.StartGame()

	if gsm.State() != StatePlaying {
		t.Errorf("State should be StatePlaying after StartGame, got %v", gsm.State())
	}
}

func TestGameStateManager_PauseGame(t *testing.T) {
	gsm := NewGameStateManager()
	gsm.StartGame()

	gsm.PauseGame()

	if gsm.State() != StatePaused {
		t.Errorf("State should be StatePaused after PauseGame, got %v", gsm.State())
	}
}

func TestGameStateManager_ResumeGame(t *testing.T) {
	gsm := NewGameStateManager()
	gsm.StartGame()
	gsm.PauseGame()

	gsm.ResumeGame()

	if gsm.State() != StatePlaying {
		t.Errorf("State should be StatePlaying after ResumeGame, got %v", gsm.State())
	}
}

func TestGameStateManager_GameOver(t *testing.T) {
	gsm := NewGameStateManager()
	gsm.StartGame()

	gsm.GameOver(1000, 5)

	if gsm.State() != StateGameOver {
		t.Errorf("State should be StateGameOver, got %v", gsm.State())
	}
	if gsm.FinalScore() != 1000 {
		t.Errorf("FinalScore should be 1000, got %d", gsm.FinalScore())
	}
	if gsm.FinalWave() != 5 {
		t.Errorf("FinalWave should be 5, got %d", gsm.FinalWave())
	}
}

func TestGameStateManager_ReturnToMainMenu(t *testing.T) {
	gsm := NewGameStateManager()
	gsm.StartGame()

	gsm.ReturnToMainMenu()

	if gsm.State() != StateMainMenu {
		t.Errorf("State should be StateMainMenu, got %v", gsm.State())
	}
}

func TestGameStateManager_StateChangeCallback(t *testing.T) {
	gsm := NewGameStateManager()

	var fromState, toState GameState
	gsm.SetStateChangeCallback(func(from, to GameState) {
		fromState = from
		toState = to
	})

	gsm.StartGame()

	if fromState != StateMainMenu {
		t.Errorf("Callback should receive from=StateMainMenu, got %v", fromState)
	}
	if toState != StatePlaying {
		t.Errorf("Callback should receive to=StatePlaying, got %v", toState)
	}
}

func TestGameStateManager_IsPlaying(t *testing.T) {
	gsm := NewGameStateManager()

	if gsm.IsPlaying() {
		t.Error("Should not be playing initially")
	}

	gsm.StartGame()

	if !gsm.IsPlaying() {
		t.Error("Should be playing after StartGame")
	}
}

func TestGameStateManager_IsPaused(t *testing.T) {
	gsm := NewGameStateManager()
	gsm.StartGame()

	if gsm.IsPaused() {
		t.Error("Should not be paused")
	}

	gsm.PauseGame()

	if !gsm.IsPaused() {
		t.Error("Should be paused")
	}
}

func TestGameStateManager_IsMenuActive(t *testing.T) {
	gsm := NewGameStateManager()

	if !gsm.IsMenuActive() {
		t.Error("Menu should be active in MainMenu")
	}

	gsm.StartGame()

	if gsm.IsMenuActive() {
		t.Error("Menu should not be active while playing")
	}
}

func TestGameStateManager_InvalidTransitions(t *testing.T) {
	gsm := NewGameStateManager()

	// Can't pause from main menu
	gsm.PauseGame()
	if gsm.State() != StateMainMenu {
		t.Error("Should remain in MainMenu, can't pause from menu")
	}

	// Can't resume from main menu
	gsm.ResumeGame()
	if gsm.State() != StateMainMenu {
		t.Error("Should remain in MainMenu, can't resume from menu")
	}

	// Can't game over from main menu
	gsm.GameOver(100, 1)
	if gsm.State() != StateMainMenu {
		t.Error("Should remain in MainMenu, can't game over from menu")
	}
}

func TestDefaultMenuItems(t *testing.T) {
	items := DefaultMenuItems()

	if len(items.MainMenu) != 3 {
		t.Errorf("MainMenu should have 3 items, got %d", len(items.MainMenu))
	}
	if len(items.PauseMenu) != 3 {
		t.Errorf("PauseMenu should have 3 items, got %d", len(items.PauseMenu))
	}
	if len(items.GameOver) != 2 {
		t.Errorf("GameOver should have 2 items, got %d", len(items.GameOver))
	}
}

func TestNewMenuController(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	if mc.SelectionIndex() != 0 {
		t.Errorf("Initial selection should be 0, got %d", mc.SelectionIndex())
	}
}

func TestMenuController_GetCurrentItems(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	// Main menu state
	items := mc.GetCurrentItems()
	if len(items) != 3 {
		t.Errorf("Expected 3 main menu items, got %d", len(items))
	}

	// Playing state - no menu
	gsm.StartGame()
	items = mc.GetCurrentItems()
	if items != nil {
		t.Errorf("Expected nil items while playing, got %v", items)
	}

	// Paused state
	gsm.PauseGame()
	items = mc.GetCurrentItems()
	if len(items) != 3 {
		t.Errorf("Expected 3 pause menu items, got %d", len(items))
	}
}

func TestMenuController_MoveUp(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	// Move up from 0 should wrap to last item
	mc.MoveUp()
	if mc.SelectionIndex() != 2 {
		t.Errorf("Selection should wrap to 2, got %d", mc.SelectionIndex())
	}

	// Move up again
	mc.MoveUp()
	if mc.SelectionIndex() != 1 {
		t.Errorf("Selection should be 1, got %d", mc.SelectionIndex())
	}
}

func TestMenuController_MoveDown(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	mc.MoveDown()
	if mc.SelectionIndex() != 1 {
		t.Errorf("Selection should be 1, got %d", mc.SelectionIndex())
	}

	mc.MoveDown()
	if mc.SelectionIndex() != 2 {
		t.Errorf("Selection should be 2, got %d", mc.SelectionIndex())
	}

	// Wrap to 0
	mc.MoveDown()
	if mc.SelectionIndex() != 0 {
		t.Errorf("Selection should wrap to 0, got %d", mc.SelectionIndex())
	}
}

func TestMenuController_Select_Start(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	// "Start" is the first item
	mc.Select()

	if gsm.State() != StatePlaying {
		t.Errorf("Expected StatePlaying after selecting Start, got %v", gsm.State())
	}
}

func TestMenuController_Select_Resume(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	gsm.StartGame()
	gsm.PauseGame()

	// "Resume" is the first item in pause menu
	mc.ResetSelection()
	mc.Select()

	if gsm.State() != StatePlaying {
		t.Errorf("Expected StatePlaying after selecting Resume, got %v", gsm.State())
	}
}

func TestMenuController_Select_QuitToMenu(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	gsm.StartGame()
	gsm.PauseGame()

	// Navigate to "Quit to Menu" (third item, index 2)
	mc.ResetSelection()
	mc.MoveDown()
	mc.MoveDown()
	mc.Select()

	if gsm.State() != StateMainMenu {
		t.Errorf("Expected StateMainMenu after Quit to Menu, got %v", gsm.State())
	}
}

func TestMenuController_ActionCallback(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	var receivedAction string
	mc.SetActionCallback(func(action string) {
		receivedAction = action
	})

	mc.Select() // Select "Start"

	if receivedAction != "start" {
		t.Errorf("Expected action 'start', got '%s'", receivedAction)
	}
}

func TestMenuController_ResetSelection(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	mc.MoveDown()
	mc.MoveDown()

	mc.ResetSelection()

	if mc.SelectionIndex() != 0 {
		t.Errorf("Selection should be reset to 0, got %d", mc.SelectionIndex())
	}
}

func TestMenuController_SelectResetsIndexOnStart(t *testing.T) {
	gsm := NewGameStateManager()
	mc := NewMenuController(gsm)

	mc.MoveDown() // Settings
	mc.MoveUp()   // Back to Start
	mc.Select()   // Start game

	// After transitioning, selection should reset
	if mc.SelectionIndex() != 0 {
		t.Errorf("Selection should reset to 0 on state change, got %d", mc.SelectionIndex())
	}
}

func TestHUD_Update(t *testing.T) {
	hud := NewHUD()

	hud.Update(75.0, 50.0, 1500, 3, 5)

	if hud.Health != 75.0 {
		t.Errorf("Health expected 75, got %f", hud.Health)
	}
	if hud.Shield != 50.0 {
		t.Errorf("Shield expected 50, got %f", hud.Shield)
	}
	if hud.Score != 1500 {
		t.Errorf("Score expected 1500, got %d", hud.Score)
	}
	if hud.Wave != 3 {
		t.Errorf("Wave expected 3, got %d", hud.Wave)
	}
	if hud.Combo != 5 {
		t.Errorf("Combo expected 5, got %d", hud.Combo)
	}
}

func TestHUD_SetGenre(t *testing.T) {
	hud := NewHUD()

	hud.SetGenre("fantasy")

	// Genre is set (internal field verification via behavior would happen in render tests)
}

func TestMenu_SetGenre(t *testing.T) {
	menu := NewMenu()

	menu.SetGenre("horror")

	// Genre is set (internal field verification via behavior would happen in render tests)
}

func TestTutorial_Advance(t *testing.T) {
	tut := NewTutorial()

	if tut.Step != 0 {
		t.Errorf("Initial step should be 0, got %d", tut.Step)
	}

	tut.Advance()

	if tut.Step != 1 {
		t.Errorf("Step should be 1 after advance, got %d", tut.Step)
	}
}

func TestTutorial_Complete(t *testing.T) {
	tut := NewTutorial()

	if !tut.Active {
		t.Error("Tutorial should be active initially")
	}

	tut.Complete()

	if tut.Active {
		t.Error("Tutorial should be inactive after Complete")
	}
}
