// Package rendering provides sprite generation, animation, particle systems,
// dynamic lighting, draw batching, and viewport culling.
package rendering

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/opd-ai/velocity/pkg/procgen/genre"
)

// SpriteType identifies the kind of sprite to generate.
type SpriteType int

const (
	SpriteTypeShip SpriteType = iota
	SpriteTypeEnemy
	SpriteTypeProjectile
)

// SpriteKey uniquely identifies a cached sprite.
type SpriteKey struct {
	GenreID string
	Type    SpriteType
	Variant int
}

// GenerateShipSprite creates a procedurally generated ship sprite.
func GenerateShipSprite(rng *rand.Rand, genreID string, size int) *image.RGBA {
	if size < 8 {
		size = 8
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	preset := genre.GetPreset(genreID)
	palette := preset.Colors

	if len(palette) == 0 {
		palette = []color.RGBA{{R: 200, G: 200, B: 200, A: 255}}
	}

	// Generate left half with symmetric mirroring
	fillChance := 0.45
	halfWidth := size / 2
	for y := 0; y < size; y++ {
		for x := 0; x < halfWidth; x++ {
			if shouldFillPixel(rng, x, y, size, fillChance) {
				c := palette[rng.Intn(len(palette))]
				setPixel(img, x, y, c)
				// Mirror to right half
				setPixel(img, size-1-x, y, c)
			}
		}
	}

	// Fill center column for odd-width sprites
	if size%2 == 1 {
		centerX := halfWidth
		for y := 0; y < size; y++ {
			if rng.Float64() < fillChance {
				c := palette[rng.Intn(len(palette))]
				setPixel(img, centerX, y, c)
			}
		}
	}

	// Add a cockpit/engine detail in center
	addCenterDetail(img, rng, size, palette)

	return img
}

// GenerateEnemySprite creates a procedurally generated enemy sprite.
func GenerateEnemySprite(rng *rand.Rand, genreID string, size int) *image.RGBA {
	if size < 8 {
		size = 8
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	preset := genre.GetPreset(genreID)
	palette := preset.Colors

	if len(palette) == 0 {
		palette = []color.RGBA{{R: 255, G: 100, B: 100, A: 255}}
	}

	// Enemies have a more aggressive fill pattern
	fillChance := 0.55
	halfWidth := size / 2
	for y := 0; y < size; y++ {
		for x := 0; x < halfWidth; x++ {
			if shouldFillPixel(rng, x, y, size, fillChance) {
				c := palette[rng.Intn(len(palette))]
				setPixel(img, x, y, c)
				// Mirror to right half
				setPixel(img, size-1-x, y, c)
			}
		}
	}

	// Fill center column for odd-width sprites
	if size%2 == 1 {
		centerX := halfWidth
		for y := 0; y < size; y++ {
			if rng.Float64() < fillChance {
				c := palette[rng.Intn(len(palette))]
				setPixel(img, centerX, y, c)
			}
		}
	}

	// Add hostile accent
	addHostileAccent(img, rng, size, palette)

	return img
}

// GenerateProjectileSprite creates a procedurally generated projectile sprite.
func GenerateProjectileSprite(rng *rand.Rand, genreID string, size int) *image.RGBA {
	if size < 4 {
		size = 4
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	preset := genre.GetPreset(genreID)
	palette := preset.Colors

	if len(palette) == 0 {
		palette = []color.RGBA{{R: 255, G: 255, B: 100, A: 255}}
	}

	// Projectiles are simple symmetric patterns
	c := palette[rng.Intn(len(palette))]
	centerX := size / 2
	centerY := size / 2

	// Simple diamond/cross pattern
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := abs(x - centerX)
			dy := abs(y - centerY)
			if dx+dy <= size/2 {
				setPixel(img, x, y, c)
			}
		}
	}

	return img
}

// shouldFillPixel determines if a pixel should be filled based on position.
func shouldFillPixel(rng *rand.Rand, x, y, size int, baseChance float64) bool {
	// Higher fill probability toward center
	centerDist := float64(abs(x-size/4)+abs(y-size/2)) / float64(size)
	adjustedChance := baseChance * (1.2 - centerDist*0.4)

	return rng.Float64() < adjustedChance
}

// addCenterDetail adds a cockpit or engine detail to the sprite center.
func addCenterDetail(img *image.RGBA, rng *rand.Rand, size int, palette []color.RGBA) {
	centerY := size / 3

	// Small bright cockpit area - draw symmetrically
	c := palette[0]
	c.R = min(c.R+50, 255)
	c.G = min(c.G+50, 255)
	c.B = min(c.B+50, 255)

	halfWidth := size / 2
	for dy := -1; dy <= 1; dy++ {
		y := centerY + dy
		if y < 0 || y >= size {
			continue
		}
		// Draw center column
		for dx := -1; dx <= 1; dx++ {
			// Draw on both sides symmetrically
			if dx < 0 {
				leftX := halfWidth + dx
				rightX := halfWidth - 1 - dx
				if leftX >= 0 && rightX < size {
					setPixel(img, leftX, y, c)
					setPixel(img, rightX, y, c)
				}
			} else if dx == 0 && size%2 == 1 {
				// Center pixel for odd-sized sprites
				setPixel(img, halfWidth, y, c)
			} else if dx > 0 {
				leftX := halfWidth - dx
				rightX := halfWidth - 1 + dx
				if leftX >= 0 && rightX < size {
					setPixel(img, leftX, y, c)
					setPixel(img, rightX, y, c)
				}
			}
		}
	}
}

// addHostileAccent adds menacing details to enemy sprites.
func addHostileAccent(img *image.RGBA, rng *rand.Rand, size int, palette []color.RGBA) {
	// Red accent for enemies
	accent := color.RGBA{R: 255, G: 50, B: 50, A: 255}
	if len(palette) > 0 {
		accent = palette[0]
		accent.R = min(accent.R+100, 255)
	}

	// Add eye-like details
	eyeY := size / 3
	eyeX1 := size / 4
	eyeX2 := size - 1 - size/4

	setPixel(img, eyeX1, eyeY, accent)
	setPixel(img, eyeX2, eyeY, accent)
}

// setPixel safely sets a pixel in the image.
func setPixel(img *image.RGBA, x, y int, c color.RGBA) {
	bounds := img.Bounds()
	if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
		img.SetRGBA(x, y, c)
	}
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// min returns the minimum of two uint8 values.
func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}
