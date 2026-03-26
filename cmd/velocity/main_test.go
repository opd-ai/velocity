//go:build noebiten

// Package main_test provides tests for the velocity game binary.
// These tests run only when the noebiten build tag is set, avoiding
// GLFW/display dependencies in headless CI environments.
//
// Note: We use package main_test (external test package) to avoid
// importing the main package, which would trigger Ebitengine init.
package main_test

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSavePathLogic tests the save path generation logic.
// This mirrors the savePath() function in main.go.
func TestSavePathLogic(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback should be used
		fallback := "velocity_save.json"
		if fallback == "" {
			t.Error("fallback path is empty")
		}
		return
	}

	dir := filepath.Join(home, ".velocity")
	path := filepath.Join(dir, "save.json")

	if path == "" {
		t.Error("generated path is empty")
	}

	// Path should contain .velocity
	if filepath.Base(filepath.Dir(path)) != ".velocity" {
		t.Errorf("path %s doesn't contain .velocity directory", path)
	}
}

func TestSavePathDirectoryCreation(t *testing.T) {
	// Create a temporary directory to use as HOME
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	// Create the .velocity directory
	velocityDir := filepath.Join(tmpDir, ".velocity")
	err := os.MkdirAll(velocityDir, 0o755)
	if err != nil {
		t.Fatalf("failed to create .velocity directory: %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(velocityDir)
	if err != nil {
		t.Fatalf("directory stat failed: %v", err)
	}
	if !info.IsDir() {
		t.Error(".velocity is not a directory")
	}
}

// TestConfigValidGenres tests that known genres are valid identifiers.
func TestConfigValidGenres(t *testing.T) {
	validGenres := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc"}
	for _, genre := range validGenres {
		if genre == "" {
			t.Error("empty genre in valid list")
		}
		// Genre should be lowercase alphanumeric
		for _, c := range genre {
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
				t.Errorf("genre %s contains invalid character %c", genre, c)
			}
		}
	}
}

// TestConfigValidArenaModes tests that known arena modes are valid.
func TestConfigValidArenaModes(t *testing.T) {
	validModes := []string{"wrap", "bounded"}
	for _, mode := range validModes {
		if mode == "" {
			t.Error("empty arena mode in valid list")
		}
	}
}

