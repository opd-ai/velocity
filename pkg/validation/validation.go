// Package validation provides input and configuration validation.
package validation

import "fmt"

// ValidGenres contains the set of supported genre identifiers.
var ValidGenres = map[string]bool{
	"fantasy":  true,
	"scifi":    true,
	"horror":   true,
	"cyberpunk": true,
	"postapoc": true,
}

// ValidArenaModes contains the set of supported arena modes.
var ValidArenaModes = map[string]bool{
	"wrap":    true,
	"bounded": true,
}

// ValidateGenre returns an error if the genre is not supported.
func ValidateGenre(genre string) error {
	if !ValidGenres[genre] {
		return fmt.Errorf("invalid genre %q", genre)
	}
	return nil
}

// ValidateArenaMode returns an error if the arena mode is not supported.
func ValidateArenaMode(mode string) error {
	if !ValidArenaModes[mode] {
		return fmt.Errorf("invalid arena mode %q", mode)
	}
	return nil
}

// ValidatePort returns an error if the port is out of valid range.
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %d", port)
	}
	return nil
}
