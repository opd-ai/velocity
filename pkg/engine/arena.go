// Package engine provides the core ECS framework, game loop, deterministic RNG,
// input handling, and camera system.
package engine

// ArenaMode determines how entities interact with screen boundaries.
type ArenaMode string

const (
	ArenaModeWrap    ArenaMode = "wrap"
	ArenaModeBounded ArenaMode = "bounded"
)

// ArenaSystem handles entity position clamping at screen boundaries.
type ArenaSystem struct {
	world  *World
	mode   ArenaMode
	width  float64
	height float64
}

// NewArenaSystem creates an arena system with the given dimensions and mode.
func NewArenaSystem(world *World, width, height int, mode ArenaMode) *ArenaSystem {
	return &ArenaSystem{
		world:  world,
		mode:   mode,
		width:  float64(width),
		height: float64(height),
	}
}

// SetMode changes the arena boundary behavior.
func (as *ArenaSystem) SetMode(mode ArenaMode) {
	as.mode = mode
}

// Update processes all entities and applies boundary logic.
func (as *ArenaSystem) Update(dt float64) {
	for entity := range as.world.entities {
		as.processEntity(entity)
	}
}

// processEntity applies boundary logic to a single entity.
func (as *ArenaSystem) processEntity(entity Entity) {
	posComp, hasPos := as.world.GetComponent(entity, "position")
	if !hasPos {
		return
	}

	pos := posComp.(*Position)

	switch as.mode {
	case ArenaModeWrap:
		as.applyWrap(pos)
	case ArenaModeBounded:
		as.applyBounded(entity, pos)
	}
}

// applyWrap teleports entity to opposite edge when crossing boundary.
func (as *ArenaSystem) applyWrap(pos *Position) {
	if pos.X < 0 {
		pos.X += as.width
	} else if pos.X >= as.width {
		pos.X -= as.width
	}

	if pos.Y < 0 {
		pos.Y += as.height
	} else if pos.Y >= as.height {
		pos.Y -= as.height
	}
}

// applyBounded bounces entity off boundaries by reversing velocity.
func (as *ArenaSystem) applyBounded(entity Entity, pos *Position) {
	velComp, hasVel := as.world.GetComponent(entity, "velocity")

	var vel *Velocity
	if hasVel {
		vel = velComp.(*Velocity)
	}

	// Handle X boundary
	if pos.X < 0 {
		pos.X = 0
		if vel != nil {
			vel.VX = -vel.VX
		}
	} else if pos.X >= as.width {
		pos.X = as.width - 1
		if vel != nil {
			vel.VX = -vel.VX
		}
	}

	// Handle Y boundary
	if pos.Y < 0 {
		pos.Y = 0
		if vel != nil {
			vel.VY = -vel.VY
		}
	} else if pos.Y >= as.height {
		pos.Y = as.height - 1
		if vel != nil {
			vel.VY = -vel.VY
		}
	}
}

// IsEntityInBounds returns true if the entity is within the arena.
func (as *ArenaSystem) IsEntityInBounds(entity Entity) bool {
	posComp, hasPos := as.world.GetComponent(entity, "position")
	if !hasPos {
		return false
	}

	pos := posComp.(*Position)
	return pos.X >= 0 && pos.X < as.width && pos.Y >= 0 && pos.Y < as.height
}
