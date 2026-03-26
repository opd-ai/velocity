package version

import "testing"

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	if v == "" {
		t.Error("GetVersion() returned empty string")
	}
	// Default value is "0.0.0-dev"
	if v != "0.0.0-dev" {
		t.Logf("Version was overridden to: %s", v)
	}
}

func TestGetSaveVersion(t *testing.T) {
	sv := GetSaveVersion()
	if sv != SaveVersion {
		t.Errorf("GetSaveVersion() = %d, want %d", sv, SaveVersion)
	}
	if sv < 1 {
		t.Errorf("SaveVersion should be at least 1, got %d", sv)
	}
}

func TestVersionVariable(t *testing.T) {
	// Version can be set at compile time
	original := Version
	Version = "1.2.3-test"
	defer func() { Version = original }()

	if GetVersion() != "1.2.3-test" {
		t.Errorf("GetVersion() = %s, want 1.2.3-test", GetVersion())
	}
}

func TestSaveVersionConstant(t *testing.T) {
	// SaveVersion should be a positive integer
	if SaveVersion <= 0 {
		t.Errorf("SaveVersion = %d, should be positive", SaveVersion)
	}
}
