// Package genre provides genre post-processing presets and asset generation helpers.
package genre

import "image/color"

// ID constants for supported genres.
const (
	Fantasy   = "fantasy"
	SciFi     = "scifi"
	Horror    = "horror"
	Cyberpunk = "cyberpunk"
	PostApoc  = "postapoc"
)

// All returns a list of all supported genre IDs.
func All() []string {
	return []string{Fantasy, SciFi, Horror, Cyberpunk, PostApoc}
}

// Preset holds genre-specific visual post-processing parameters.
type Preset struct {
	GenreID    string
	BloomScale float64
	Saturation float64
	NeonGlow   float64
	GrainLevel float64
	Colors     []color.RGBA
}

// GetPreset returns the visual preset for a genre.
func GetPreset(genreID string) Preset {
	switch genreID {
	case SciFi:
		return Preset{
			GenreID:    SciFi,
			BloomScale: 1.2,
			Colors:     sciFiColors(),
		}
	case Horror:
		return Preset{
			GenreID:    Horror,
			Saturation: 0.4,
			Colors:     horrorColors(),
		}
	case Cyberpunk:
		return Preset{
			GenreID:  Cyberpunk,
			NeonGlow: 1.5,
			Colors:   cyberpunkColors(),
		}
	case Fantasy:
		return Preset{
			GenreID:    Fantasy,
			BloomScale: 0.8,
			Saturation: 1.1,
			Colors:     fantasyColors(),
		}
	case PostApoc:
		return Preset{
			GenreID:    PostApoc,
			GrainLevel: 0.6,
			Colors:     postApocColors(),
		}
	default:
		return Preset{GenreID: genreID, Colors: sciFiColors()}
	}
}

// sciFiColors returns the sci-fi genre color palette.
func sciFiColors() []color.RGBA {
	return []color.RGBA{
		{R: 80, G: 120, B: 200, A: 255},  // Steel blue
		{R: 100, G: 180, B: 255, A: 255}, // Cyan
		{R: 200, G: 200, B: 220, A: 255}, // Silver
		{R: 50, G: 80, B: 120, A: 255},   // Dark blue
		{R: 255, G: 100, B: 100, A: 255}, // Engine red
	}
}

// horrorColors returns the horror genre color palette.
func horrorColors() []color.RGBA {
	return []color.RGBA{
		{R: 80, G: 60, B: 80, A: 255},   // Dark purple
		{R: 120, G: 80, B: 100, A: 255}, // Dried blood
		{R: 60, G: 80, B: 60, A: 255},   // Organic green
		{R: 40, G: 40, B: 50, A: 255},   // Void
		{R: 150, G: 100, B: 80, A: 255}, // Bone
	}
}

// cyberpunkColors returns the cyberpunk genre color palette.
func cyberpunkColors() []color.RGBA {
	return []color.RGBA{
		{R: 255, G: 0, B: 200, A: 255}, // Hot pink
		{R: 0, G: 255, B: 255, A: 255}, // Cyan neon
		{R: 255, G: 255, B: 0, A: 255}, // Yellow neon
		{R: 100, G: 0, B: 150, A: 255}, // Purple
		{R: 50, G: 50, B: 80, A: 255},  // Dark chrome
	}
}

// fantasyColors returns the fantasy genre color palette.
func fantasyColors() []color.RGBA {
	return []color.RGBA{
		{R: 200, G: 180, B: 100, A: 255}, // Gold
		{R: 100, G: 150, B: 200, A: 255}, // Mystic blue
		{R: 180, G: 100, B: 180, A: 255}, // Arcane purple
		{R: 255, G: 220, B: 150, A: 255}, // Enchanted glow
		{R: 120, G: 80, B: 60, A: 255},   // Wood
	}
}

// postApocColors returns the post-apocalyptic genre color palette.
func postApocColors() []color.RGBA {
	return []color.RGBA{
		{R: 160, G: 120, B: 80, A: 255},  // Rust
		{R: 100, G: 100, B: 90, A: 255},  // Gunmetal
		{R: 80, G: 60, B: 40, A: 255},    // Dirt
		{R: 140, G: 140, B: 120, A: 255}, // Weathered steel
		{R: 200, G: 160, B: 100, A: 255}, // Sand
	}
}
