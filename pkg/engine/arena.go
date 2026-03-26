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
	vel := as.getVelocity(entity)
	as.bounceOnXBoundary(pos, vel)
	as.bounceOnYBoundary(pos, vel)
}

// getVelocity retrieves the velocity component for an entity, or nil.
func (as *ArenaSystem) getVelocity(entity Entity) *Velocity {
	velComp, hasVel := as.world.GetComponent(entity, "velocity")
	if !hasVel {
		return nil
	}
	return velComp.(*Velocity)
}

// bounceOnXBoundary handles left/right edge collision.
func (as *ArenaSystem) bounceOnXBoundary(pos *Position, vel *Velocity) {
	if pos.X < 0 {
		pos.X = 0
		reverseVX(vel)
	} else if pos.X >= as.width {
		pos.X = as.width - 1
		reverseVX(vel)
	}
}

// bounceOnYBoundary handles top/bottom edge collision.
func (as *ArenaSystem) bounceOnYBoundary(pos *Position, vel *Velocity) {
	if pos.Y < 0 {
		pos.Y = 0
		reverseVY(vel)
	} else if pos.Y >= as.height {
		pos.Y = as.height - 1
		reverseVY(vel)
	}
}

// reverseVX negates the X velocity component if vel is non-nil.
func reverseVX(vel *Velocity) {
	if vel != nil {
		vel.VX = -vel.VX
	}
}

// reverseVY negates the Y velocity component if vel is non-nil.
func reverseVY(vel *Velocity) {
	if vel != nil {
		vel.VY = -vel.VY
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
