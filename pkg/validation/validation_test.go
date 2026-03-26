package validation

import (
	"testing"
)

func TestValidateGenre_Valid(t *testing.T) {
	validGenres := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc"}

	for _, genre := range validGenres {
		t.Run(genre, func(t *testing.T) {
			err := ValidateGenre(genre)
			if err != nil {
				t.Errorf("ValidateGenre(%q) returned error: %v", genre, err)
			}
		})
	}
}

func TestValidateGenre_Invalid(t *testing.T) {
	invalidGenres := []string{"", "sci-fi", "SCIFI", "unknown", "space", "medieval"}

	for _, genre := range invalidGenres {
		t.Run(genre, func(t *testing.T) {
			err := ValidateGenre(genre)
			if err == nil {
				t.Errorf("ValidateGenre(%q) should return error", genre)
			}
		})
	}
}

func TestValidateArenaMode_Valid(t *testing.T) {
	validModes := []string{"wrap", "bounded"}

	for _, mode := range validModes {
		t.Run(mode, func(t *testing.T) {
			err := ValidateArenaMode(mode)
			if err != nil {
				t.Errorf("ValidateArenaMode(%q) returned error: %v", mode, err)
			}
		})
	}
}

func TestValidateArenaMode_Invalid(t *testing.T) {
	invalidModes := []string{"", "WRAP", "Bounded", "infinite", "scroll"}

	for _, mode := range invalidModes {
		t.Run(mode, func(t *testing.T) {
			err := ValidateArenaMode(mode)
			if err == nil {
				t.Errorf("ValidateArenaMode(%q) should return error", mode)
			}
		})
	}
}

func TestValidatePort_Valid(t *testing.T) {
	validPorts := []int{1, 80, 443, 8080, 27015, 65535}

	for _, port := range validPorts {
		t.Run("", func(t *testing.T) {
			err := ValidatePort(port)
			if err != nil {
				t.Errorf("ValidatePort(%d) returned error: %v", port, err)
			}
		})
	}
}

func TestValidatePort_Invalid(t *testing.T) {
	invalidPorts := []int{0, -1, -100, 65536, 100000}

	for _, port := range invalidPorts {
		t.Run("", func(t *testing.T) {
			err := ValidatePort(port)
			if err == nil {
				t.Errorf("ValidatePort(%d) should return error", port)
			}
		})
	}
}

func TestValidGenres_Contains(t *testing.T) {
	expected := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc"}

	if len(ValidGenres) != len(expected) {
		t.Errorf("expected %d genres, got %d", len(expected), len(ValidGenres))
	}

	for _, genre := range expected {
		if !ValidGenres[genre] {
			t.Errorf("expected %q to be in ValidGenres", genre)
		}
	}
}

func TestValidArenaModes_Contains(t *testing.T) {
	expected := []string{"wrap", "bounded"}

	if len(ValidArenaModes) != len(expected) {
		t.Errorf("expected %d modes, got %d", len(expected), len(ValidArenaModes))
	}

	for _, mode := range expected {
		if !ValidArenaModes[mode] {
			t.Errorf("expected %q to be in ValidArenaModes", mode)
		}
	}
}

func TestValidateGenre_ErrorMessage(t *testing.T) {
	err := ValidateGenre("invalid")
	if err == nil {
		t.Fatal("expected error")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestValidateArenaMode_ErrorMessage(t *testing.T) {
	err := ValidateArenaMode("invalid")
	if err == nil {
		t.Fatal("expected error")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestValidatePort_ErrorMessage(t *testing.T) {
	err := ValidatePort(-1)
	if err == nil {
		t.Fatal("expected error")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("expected non-empty error message")
	}
}

func BenchmarkValidateGenre(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateGenre("scifi")
	}
}

func BenchmarkValidateArenaMode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateArenaMode("wrap")
	}
}

func BenchmarkValidatePort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidatePort(8080)
	}
}
