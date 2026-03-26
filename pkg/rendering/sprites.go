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
	const defaultFillChance = 0.45
	defaultPalette := []color.RGBA{{R: 200, G: 200, B: 200, A: 255}}

	img := generateSymmetricBody(rng, genreID, size, defaultFillChance, defaultPalette)
	addCenterDetail(img, rng, size, getPalette(genreID, defaultPalette))
	return img
}

// GenerateEnemySprite creates a procedurally generated enemy sprite.
func GenerateEnemySprite(rng *rand.Rand, genreID string, size int) *image.RGBA {
	const enemyFillChance = 0.55
	defaultPalette := []color.RGBA{{R: 255, G: 100, B: 100, A: 255}}

	img := generateSymmetricBody(rng, genreID, size, enemyFillChance, defaultPalette)
	addHostileAccent(img, rng, size, getPalette(genreID, defaultPalette))
	return img
}

// generateSymmetricBody creates the common symmetric pixel pattern for sprites.
func generateSymmetricBody(rng *rand.Rand, genreID string, size int, fillChance float64, defaultPalette []color.RGBA) *image.RGBA {
	if size < 8 {
		size = 8
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	palette := getPalette(genreID, defaultPalette)

	fillSymmetricHalf(img, rng, size, fillChance, palette)
	fillCenterColumn(img, rng, size, fillChance, palette)

	return img
}

// getPalette returns the genre palette or the provided default.
func getPalette(genreID string, defaultPalette []color.RGBA) []color.RGBA {
	preset := genre.GetPreset(genreID)
	if len(preset.Colors) > 0 {
		return preset.Colors
	}
	return defaultPalette
}

// fillSymmetricHalf generates the left half with symmetric mirroring.
func fillSymmetricHalf(img *image.RGBA, rng *rand.Rand, size int, fillChance float64, palette []color.RGBA) {
	halfWidth := size / 2
	for y := 0; y < size; y++ {
		for x := 0; x < halfWidth; x++ {
			if shouldFillPixel(rng, x, y, size, fillChance) {
				c := palette[rng.Intn(len(palette))]
				setPixel(img, x, y, c)
				setPixel(img, size-1-x, y, c)
			}
		}
	}
}

// fillCenterColumn fills the center column for odd-width sprites.
func fillCenterColumn(img *image.RGBA, rng *rand.Rand, size int, fillChance float64, palette []color.RGBA) {
	if size%2 == 0 {
		return
	}
	centerX := size / 2
	for y := 0; y < size; y++ {
		if rng.Float64() < fillChance {
			c := palette[rng.Intn(len(palette))]
			setPixel(img, centerX, y, c)
		}
	}
}

// GenerateProjectileSprite creates a procedurally generated projectile sprite.
func GenerateProjectileSprite(rng *rand.Rand, genreID string, size int) *image.RGBA {
	size = clampProjectileSize(size)
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	c := selectProjectileColor(rng, genreID)
	drawDiamondPattern(img, size, c)
	return img
}

// clampProjectileSize ensures minimum projectile size of 4 pixels.
func clampProjectileSize(size int) int {
	if size < 4 {
		return 4
	}
	return size
}

// selectProjectileColor picks a color from the genre palette or uses default.
func selectProjectileColor(rng *rand.Rand, genreID string) color.RGBA {
	preset := genre.GetPreset(genreID)
	palette := preset.Colors
	if len(palette) == 0 {
		return color.RGBA{R: 255, G: 255, B: 100, A: 255}
	}
	return palette[rng.Intn(len(palette))]
}

// drawDiamondPattern draws a diamond/cross pattern for the projectile.
func drawDiamondPattern(img *image.RGBA, size int, c color.RGBA) {
	centerX := size / 2
	centerY := size / 2
	halfSize := size / 2
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if isDiamondPixel(x, y, centerX, centerY, halfSize) {
				setPixel(img, x, y, c)
			}
		}
	}
}

// isDiamondPixel checks if a pixel falls within the diamond pattern.
func isDiamondPixel(x, y, centerX, centerY, halfSize int) bool {
	dx := abs(x - centerX)
	dy := abs(y - centerY)
	return dx+dy <= halfSize
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
	halfWidth := size / 2
	c := brightenColor(palette[0], 50)

	// Draw a 3x3 symmetric cockpit area centered at (halfWidth, centerY)
	for dy := -1; dy <= 1; dy++ {
		y := centerY + dy
		if y < 0 || y >= size {
			continue
		}
		drawSymmetricRow(img, y, halfWidth, size, c)
	}
}

// brightenColor increases RGB values by the given amount, clamping at 255.
func brightenColor(c color.RGBA, amount uint8) color.RGBA {
	return color.RGBA{
		R: min(c.R+amount, 255),
		G: min(c.G+amount, 255),
		B: min(c.B+amount, 255),
		A: c.A,
	}
}

// drawSymmetricRow draws a symmetric row of pixels centered at halfWidth.
func drawSymmetricRow(img *image.RGBA, y, halfWidth, size int, c color.RGBA) {
	// Draw center pixel for odd-sized sprites
	if size%2 == 1 {
		setPixel(img, halfWidth, y, c)
	}

	// Draw symmetric pairs at offsets ±1
	for offset := 1; offset <= 1; offset++ {
		leftX := halfWidth - offset
		rightX := halfWidth - 1 + offset
		if size%2 == 1 {
			rightX = halfWidth + offset
		}
		if leftX >= 0 && rightX < size {
			setPixel(img, leftX, y, c)
			setPixel(img, rightX, y, c)
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
