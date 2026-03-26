// Package rendering provides sprite generation, animation, particle systems,
// dynamic lighting, draw batching, and viewport culling.
package rendering

import "github.com/opd-ai/velocity/pkg/engine"

// Viewport represents the visible area of the game world.
type Viewport struct {
	X, Y          float64
	Width, Height float64
}

// NewViewport creates a viewport with the given dimensions.
func NewViewport(width, height int) *Viewport {
	return &Viewport{
		Width:  float64(width),
		Height: float64(height),
	}
}

// SetPosition updates the viewport center position.
func (v *Viewport) SetPosition(x, y float64) {
	v.X = x - v.Width/2
	v.Y = y - v.Height/2
}

// Contains returns true if the point is within the viewport.
func (v *Viewport) Contains(x, y float64) bool {
	return x >= v.X && x < v.X+v.Width &&
		y >= v.Y && y < v.Y+v.Height
}

// ContainsRect returns true if the rectangle overlaps the viewport.
func (v *Viewport) ContainsRect(x, y, width, height float64) bool {
	return x+width >= v.X && x < v.X+v.Width &&
		y+height >= v.Y && y < v.Y+v.Height
}

// CullContext holds culling state for a frame.
type CullContext struct {
	viewport      *Viewport
	margin        float64
	culledCount   int
	renderedCount int
}

// NewCullContext creates a new culling context for the given viewport.
func NewCullContext(viewport *Viewport, margin float64) *CullContext {
	return &CullContext{
		viewport: viewport,
		margin:   margin,
	}
}

// ShouldRender returns true if the entity should be rendered.
func (cc *CullContext) ShouldRender(x, y, width, height float64) bool {
	// Expand viewport by margin for entities partially off-screen
	expanded := &Viewport{
		X:      cc.viewport.X - cc.margin,
		Y:      cc.viewport.Y - cc.margin,
		Width:  cc.viewport.Width + 2*cc.margin,
		Height: cc.viewport.Height + 2*cc.margin,
	}

	if expanded.ContainsRect(x-width/2, y-height/2, width, height) {
		cc.renderedCount++
		return true
	}

	cc.culledCount++
	return false
}

// GetCulledCount returns the number of entities culled this frame.
func (cc *CullContext) GetCulledCount() int {
	return cc.culledCount
}

// GetRenderedCount returns the number of entities rendered this frame.
func (cc *CullContext) GetRenderedCount() int {
	return cc.renderedCount
}

// Reset clears the frame counters.
func (cc *CullContext) Reset() {
	cc.culledCount = 0
	cc.renderedCount = 0
}

// FilterVisibleEntities returns only entities within the viewport.
func FilterVisibleEntities(world *engine.World, viewport *Viewport, margin float64) []engine.Entity {
	cc := NewCullContext(viewport, margin)
	var visible []engine.Entity

	world.ForEachEntity(func(e engine.Entity) {
		posComp, hasPos := world.GetComponent(e, "position")
		if !hasPos {
			return
		}

		pos := posComp.(*engine.Position)

		// Default sprite size for culling; could be fetched from sprite component
		spriteSize := 32.0

		if cc.ShouldRender(pos.X, pos.Y, spriteSize, spriteSize) {
			visible = append(visible, e)
		}
	})

	return visible
}

// RenderOrder defines the order in which sprite types are rendered.
var RenderOrder = []SpriteType{
	SpriteTypeProjectile,
	SpriteTypeEnemy,
	SpriteTypeShip,
}

// SortBatchesByRenderOrder sorts batches according to render order.
func SortBatchesByRenderOrder(batches []DrawBatch) []DrawBatch {
	orderMap := make(map[SpriteType]int)
	for i, st := range RenderOrder {
		orderMap[st] = i
	}

	sorted := make([]DrawBatch, len(batches))
	copy(sorted, batches)

	// Simple insertion sort for small slice
	for i := 1; i < len(sorted); i++ {
		j := i
		for j > 0 && orderMap[sorted[j].Type] < orderMap[sorted[j-1].Type] {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
			j--
		}
	}

	return sorted
}
