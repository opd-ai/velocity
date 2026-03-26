package rendering

import (
	"math/rand"
	"testing"

	"github.com/opd-ai/velocity/pkg/procgen/genre"
)

func TestGenerateShipSprite(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))

	tests := []struct {
		name    string
		genreID string
		size    int
	}{
		{"scifi 16x16", genre.SciFi, 16},
		{"fantasy 16x16", genre.Fantasy, 16},
		{"horror 16x16", genre.Horror, 16},
		{"cyberpunk 16x16", genre.Cyberpunk, 16},
		{"postapoc 16x16", genre.PostApoc, 16},
		{"scifi small", genre.SciFi, 8},
		{"scifi large", genre.SciFi, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := GenerateShipSprite(rng, tt.genreID, tt.size)

			if img == nil {
				t.Fatal("expected non-nil image")
			}

			bounds := img.Bounds()
			expectedSize := tt.size
			if tt.size < 8 {
				expectedSize = 8
			}
			if bounds.Dx() != expectedSize || bounds.Dy() != expectedSize {
				t.Errorf("expected %dx%d image, got %dx%d", expectedSize, expectedSize, bounds.Dx(), bounds.Dy())
			}

			// Verify some pixels are filled
			hasPixels := false
			for y := 0; y < bounds.Dy(); y++ {
				for x := 0; x < bounds.Dx(); x++ {
					if img.RGBAAt(x, y).A > 0 {
						hasPixels = true
						break
					}
				}
			}
			if !hasPixels {
				t.Error("expected some non-transparent pixels")
			}
		})
	}
}

func TestGenerateEnemySprite(t *testing.T) {
	rng := rand.New(rand.NewSource(54321))

	for _, genreID := range genre.All() {
		t.Run(genreID, func(t *testing.T) {
			img := GenerateEnemySprite(rng, genreID, 16)

			if img == nil {
				t.Fatal("expected non-nil image")
			}

			bounds := img.Bounds()
			if bounds.Dx() != 16 || bounds.Dy() != 16 {
				t.Errorf("expected 16x16 image, got %dx%d", bounds.Dx(), bounds.Dy())
			}
		})
	}
}

func TestGenerateProjectileSprite(t *testing.T) {
	rng := rand.New(rand.NewSource(99999))

	tests := []struct {
		name string
		size int
	}{
		{"small", 4},
		{"medium", 8},
		{"tiny", 2}, // Should be clamped to 4
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := GenerateProjectileSprite(rng, genre.SciFi, tt.size)

			if img == nil {
				t.Fatal("expected non-nil image")
			}

			expectedSize := tt.size
			if expectedSize < 4 {
				expectedSize = 4
			}

			bounds := img.Bounds()
			if bounds.Dx() != expectedSize {
				t.Errorf("expected %d width, got %d", expectedSize, bounds.Dx())
			}
		})
	}
}

func TestSpriteDeterminism(t *testing.T) {
	seed := int64(42)

	rng1 := rand.New(rand.NewSource(seed))
	img1 := GenerateShipSprite(rng1, genre.SciFi, 16)

	rng2 := rand.New(rand.NewSource(seed))
	img2 := GenerateShipSprite(rng2, genre.SciFi, 16)

	// Sprites should be identical with same seed
	bounds := img1.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			c1 := img1.RGBAAt(x, y)
			c2 := img2.RGBAAt(x, y)
			if c1 != c2 {
				t.Errorf("pixel mismatch at (%d,%d): %v vs %v", x, y, c1, c2)
				return
			}
		}
	}
}

func TestSpriteSymmetry(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	img := GenerateShipSprite(rng, genre.SciFi, 16)

	bounds := img.Bounds()
	size := bounds.Dx()

	// Check horizontal symmetry
	for y := 0; y < size; y++ {
		for x := 0; x < size/2; x++ {
			left := img.RGBAAt(x, y)
			right := img.RGBAAt(size-1-x, y)
			if left != right {
				t.Errorf("symmetry broken at y=%d: left(%d)=%v, right(%d)=%v",
					y, x, left, size-1-x, right)
				return
			}
		}
	}
}

func TestSpriteKey(t *testing.T) {
	key1 := SpriteKey{GenreID: genre.SciFi, Type: SpriteTypeShip, Variant: 0}
	key2 := SpriteKey{GenreID: genre.SciFi, Type: SpriteTypeShip, Variant: 0}
	key3 := SpriteKey{GenreID: genre.Horror, Type: SpriteTypeShip, Variant: 0}

	if key1 != key2 {
		t.Error("expected equal keys to be equal")
	}
	if key1 == key3 {
		t.Error("expected different genre keys to differ")
	}
}

func BenchmarkGenerateShipSprite(b *testing.B) {
	rng := rand.New(rand.NewSource(12345))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateShipSprite(rng, genre.SciFi, 16)
	}
}

func BenchmarkGenerateEnemySprite(b *testing.B) {
	rng := rand.New(rand.NewSource(12345))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateEnemySprite(rng, genre.SciFi, 16)
	}
}
