//go:build noebiten

// Velocity — stub entry point for headless test builds.
// This file is used when building with -tags=noebiten to avoid
// Ebitengine GLFW initialization in CI environments.
package main

// main is a no-op stub for test builds.
func main() {}

// savePath returns the path to the save file.
// Duplicated here for test access without Ebitengine dependencies.
func savePath() string {
	return "velocity_save.json"
}
