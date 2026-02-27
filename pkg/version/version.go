// Package version provides embedded version information and save-file migration.
package version

// Version is the current build version, set at compile time.
var Version = "0.0.0-dev"

// SaveVersion is the current save-file format version.
const SaveVersion = 1

// GetVersion returns the current build version string.
func GetVersion() string {
	return Version
}

// GetSaveVersion returns the current save-file format version.
func GetSaveVersion() int {
	return SaveVersion
}
