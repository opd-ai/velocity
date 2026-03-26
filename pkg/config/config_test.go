package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Change to temp directory without config file
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Check display defaults
	if cfg.Display.Width != 800 {
		t.Errorf("expected width 800, got %d", cfg.Display.Width)
	}
	if cfg.Display.Height != 600 {
		t.Errorf("expected height 600, got %d", cfg.Display.Height)
	}
	if cfg.Display.Fullscreen != false {
		t.Error("expected fullscreen false")
	}
	if cfg.Display.VSync != true {
		t.Error("expected vsync true")
	}

	// Check audio defaults
	if cfg.Audio.MasterVolume != 0.8 {
		t.Errorf("expected master_volume 0.8, got %f", cfg.Audio.MasterVolume)
	}
	if cfg.Audio.MusicVolume != 0.6 {
		t.Errorf("expected music_volume 0.6, got %f", cfg.Audio.MusicVolume)
	}
	if cfg.Audio.SFXVolume != 0.8 {
		t.Errorf("expected sfx_volume 0.8, got %f", cfg.Audio.SFXVolume)
	}

	// Check gameplay defaults
	if cfg.Gameplay.Genre != "scifi" {
		t.Errorf("expected genre 'scifi', got %s", cfg.Gameplay.Genre)
	}
	if cfg.Gameplay.ArenaMode != "wrap" {
		t.Errorf("expected arena_mode 'wrap', got %s", cfg.Gameplay.ArenaMode)
	}
	if cfg.Gameplay.Seed != 0 {
		t.Errorf("expected seed 0, got %d", cfg.Gameplay.Seed)
	}

	// Check controls defaults
	if cfg.Controls.Thrust != "W" {
		t.Errorf("expected thrust 'W', got %s", cfg.Controls.Thrust)
	}
	if cfg.Controls.Fire != "Space" {
		t.Errorf("expected fire 'Space', got %s", cfg.Controls.Fire)
	}
	if cfg.Controls.Pause != "Escape" {
		t.Errorf("expected pause 'Escape', got %s", cfg.Controls.Pause)
	}
}

func TestLoad_CustomConfig(t *testing.T) {
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	// Write custom config
	configContent := `
display:
  width: 1920
  height: 1080
  fullscreen: true
  vsync: false

audio:
  master_volume: 0.5
  music_volume: 0.3
  sfx_volume: 1.0

gameplay:
  genre: fantasy
  arena_mode: bounded
  seed: 12345

controls:
  thrust: Up
  fire: Enter
`
	err = os.WriteFile(filepath.Join(tmpDir, "config.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Check custom display values
	if cfg.Display.Width != 1920 {
		t.Errorf("expected width 1920, got %d", cfg.Display.Width)
	}
	if cfg.Display.Height != 1080 {
		t.Errorf("expected height 1080, got %d", cfg.Display.Height)
	}
	if cfg.Display.Fullscreen != true {
		t.Error("expected fullscreen true")
	}

	// Check custom audio values
	if cfg.Audio.MasterVolume != 0.5 {
		t.Errorf("expected master_volume 0.5, got %f", cfg.Audio.MasterVolume)
	}

	// Check custom gameplay values
	if cfg.Gameplay.Genre != "fantasy" {
		t.Errorf("expected genre 'fantasy', got %s", cfg.Gameplay.Genre)
	}
	if cfg.Gameplay.ArenaMode != "bounded" {
		t.Errorf("expected arena_mode 'bounded', got %s", cfg.Gameplay.ArenaMode)
	}
	if cfg.Gameplay.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", cfg.Gameplay.Seed)
	}

	// Check custom controls
	if cfg.Controls.Thrust != "Up" {
		t.Errorf("expected thrust 'Up', got %s", cfg.Controls.Thrust)
	}
	if cfg.Controls.Fire != "Enter" {
		t.Errorf("expected fire 'Enter', got %s", cfg.Controls.Fire)
	}
}

func TestDisplayConfig_Fields(t *testing.T) {
	dc := DisplayConfig{
		Width:      1024,
		Height:     768,
		Fullscreen: true,
		VSync:      false,
	}

	if dc.Width != 1024 {
		t.Error("width mismatch")
	}
	if dc.Height != 768 {
		t.Error("height mismatch")
	}
	if !dc.Fullscreen {
		t.Error("fullscreen mismatch")
	}
	if dc.VSync {
		t.Error("vsync mismatch")
	}
}

func TestAudioConfig_Fields(t *testing.T) {
	ac := AudioConfig{
		MasterVolume: 0.7,
		MusicVolume:  0.5,
		SFXVolume:    0.9,
	}

	if ac.MasterVolume != 0.7 {
		t.Error("master_volume mismatch")
	}
	if ac.MusicVolume != 0.5 {
		t.Error("music_volume mismatch")
	}
	if ac.SFXVolume != 0.9 {
		t.Error("sfx_volume mismatch")
	}
}

func TestGameplayConfig_Fields(t *testing.T) {
	gc := GameplayConfig{
		Genre:     "horror",
		ArenaMode: "bounded",
		Seed:      42,
	}

	if gc.Genre != "horror" {
		t.Error("genre mismatch")
	}
	if gc.ArenaMode != "bounded" {
		t.Error("arena_mode mismatch")
	}
	if gc.Seed != 42 {
		t.Error("seed mismatch")
	}
}

func TestControlsConfig_Fields(t *testing.T) {
	cc := ControlsConfig{
		Thrust:      "W",
		RotateLeft:  "A",
		RotateRight: "D",
		Fire:        "Space",
		Secondary:   "Shift",
		Pause:       "Escape",
	}

	if cc.Thrust != "W" {
		t.Error("thrust mismatch")
	}
	if cc.Fire != "Space" {
		t.Error("fire mismatch")
	}
	if cc.Pause != "Escape" {
		t.Error("pause mismatch")
	}
}

func TestConfig_AllFieldsPopulated(t *testing.T) {
	cfg := Config{
		Display: DisplayConfig{
			Width:      800,
			Height:     600,
			Fullscreen: false,
			VSync:      true,
		},
		Audio: AudioConfig{
			MasterVolume: 0.8,
			MusicVolume:  0.6,
			SFXVolume:    0.8,
		},
		Gameplay: GameplayConfig{
			Genre:     "scifi",
			ArenaMode: "wrap",
			Seed:      0,
		},
		Controls: ControlsConfig{
			Thrust:      "W",
			RotateLeft:  "A",
			RotateRight: "D",
			Fire:        "Space",
			Secondary:   "Shift",
			Pause:       "Escape",
		},
	}

	if cfg.Display.Width == 0 {
		t.Error("display width not set")
	}
	if cfg.Audio.MasterVolume == 0 {
		t.Error("audio master not set")
	}
	if cfg.Gameplay.Genre == "" {
		t.Error("gameplay genre not set")
	}
	if cfg.Controls.Thrust == "" {
		t.Error("controls thrust not set")
	}
}
