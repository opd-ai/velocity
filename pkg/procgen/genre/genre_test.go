package genre

import (
	"testing"
)

func TestAll(t *testing.T) {
	all := All()

	if len(all) != 5 {
		t.Errorf("expected 5 genres, got %d", len(all))
	}

	expected := []string{Fantasy, SciFi, Horror, Cyberpunk, PostApoc}
	for i, g := range expected {
		if all[i] != g {
			t.Errorf("expected genre %d to be %s, got %s", i, g, all[i])
		}
	}
}

func TestGetPreset_AllGenres(t *testing.T) {
	for _, genreID := range All() {
		t.Run(genreID, func(t *testing.T) {
			preset := GetPreset(genreID)

			if preset.GenreID != genreID {
				t.Errorf("expected GenreID %s, got %s", genreID, preset.GenreID)
			}

			if len(preset.Colors) == 0 {
				t.Error("expected non-empty color palette")
			}

			// All colors should have full alpha
			for i, c := range preset.Colors {
				if c.A != 255 {
					t.Errorf("color %d has alpha %d, expected 255", i, c.A)
				}
			}
		})
	}
}

func TestGetPreset_SciFi(t *testing.T) {
	preset := GetPreset(SciFi)

	if preset.BloomScale <= 0 {
		t.Error("expected positive BloomScale for SciFi")
	}
}

func TestGetPreset_Horror(t *testing.T) {
	preset := GetPreset(Horror)

	if preset.Saturation <= 0 || preset.Saturation >= 1 {
		t.Error("expected reduced saturation for Horror")
	}
}

func TestGetPreset_Cyberpunk(t *testing.T) {
	preset := GetPreset(Cyberpunk)

	if preset.NeonGlow <= 0 {
		t.Error("expected positive NeonGlow for Cyberpunk")
	}
}

func TestGetPreset_PostApoc(t *testing.T) {
	preset := GetPreset(PostApoc)

	if preset.GrainLevel <= 0 {
		t.Error("expected positive GrainLevel for PostApoc")
	}
}

func TestGetPreset_UnknownGenre(t *testing.T) {
	preset := GetPreset("unknown")

	if preset.GenreID != "unknown" {
		t.Errorf("expected GenreID 'unknown', got %s", preset.GenreID)
	}

	// Should still have default colors
	if len(preset.Colors) == 0 {
		t.Error("expected default colors for unknown genre")
	}
}

func TestColorPalettesDistinct(t *testing.T) {
	presets := make(map[string][]byte)

	for _, genreID := range All() {
		preset := GetPreset(genreID)
		// Create a simple hash of the first color
		if len(preset.Colors) > 0 {
			c := preset.Colors[0]
			key := string([]byte{c.R, c.G, c.B})
			presets[genreID] = []byte(key)
		}
	}

	// Verify genres have different primary colors
	seen := make(map[string]string)
	for genreID, colorKey := range presets {
		key := string(colorKey)
		if existing, ok := seen[key]; ok {
			t.Errorf("genres %s and %s have the same primary color", genreID, existing)
		}
		seen[key] = genreID
	}
}
