// Package genre provides genre post-processing presets and asset generation helpers.
package genre

// ID constants for supported genres.
const (
	Fantasy  = "fantasy"
	SciFi    = "scifi"
	Horror   = "horror"
	Cyberpunk = "cyberpunk"
	PostApoc = "postapoc"
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
}

// GetPreset returns the visual preset for a genre.
func GetPreset(genreID string) Preset {
	switch genreID {
	case SciFi:
		return Preset{GenreID: SciFi, BloomScale: 1.2}
	case Horror:
		return Preset{GenreID: Horror, Saturation: 0.4}
	case Cyberpunk:
		return Preset{GenreID: Cyberpunk, NeonGlow: 1.5}
	case Fantasy:
		return Preset{GenreID: Fantasy, BloomScale: 0.8, Saturation: 1.1}
	case PostApoc:
		return Preset{GenreID: PostApoc, GrainLevel: 0.6}
	default:
		return Preset{GenreID: genreID}
	}
}
