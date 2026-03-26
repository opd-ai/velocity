// Package combat provides weapons, damage calculation, hit detection,
// and status effects.
package combat

import (
	"math"

	"github.com/opd-ai/velocity/pkg/engine"
)

// Projectile component data for a projectile entity.
type Projectile struct {
	Damage    float64
	Speed     float64
	OwnerType string // "player" or "enemy"
	Lifetime  float64
	MaxLife   float64
}

// Health component stores entity hit points.
type Health struct {
	Current float64
	Max     float64
}

// BoundingBox represents an axis-aligned bounding box for collision.
type BoundingBox struct {
	X, Y          float64
	Width, Height float64
}

// CollisionTag marks an entity as collidable with a tag.
type CollisionTag struct {
	Tag string // "player", "enemy", "projectile"
}

// ProjectileSystem manages projectile movement, lifetime, and collision.
type ProjectileSystem struct {
	world    *engine.World
	toRemove []engine.Entity
	onHit    func(projectile, target engine.Entity, damage float64)
}

// NewProjectileSystem creates a new projectile system.
func NewProjectileSystem(world *engine.World) *ProjectileSystem {
	return &ProjectileSystem{
		world:    world,
		toRemove: make([]engine.Entity, 0, 32),
	}
}

// SetHitCallback sets the callback for when a projectile hits a target.
func (ps *ProjectileSystem) SetHitCallback(fn func(projectile, target engine.Entity, damage float64)) {
	ps.onHit = fn
}

// Update moves projectiles, checks collisions, and removes expired ones.
func (ps *ProjectileSystem) Update(dt float64) {
	ps.toRemove = ps.toRemove[:0]

	ps.world.ForEachEntity(func(e engine.Entity) {
		projComp, hasProj := ps.world.GetComponent(e, "projectile")
		if !hasProj {
			return
		}

		proj := projComp.(*Projectile)

		// Update lifetime
		proj.Lifetime -= dt
		if proj.Lifetime <= 0 {
			ps.toRemove = append(ps.toRemove, e)
			return
		}

		// Move projectile
		ps.moveProjectile(e, proj, dt)

		// Check collisions
		ps.checkCollisions(e, proj)
	})

	// Remove expired/hit projectiles
	for _, e := range ps.toRemove {
		ps.world.RemoveEntity(e)
	}
}

// moveProjectile updates the projectile position based on velocity.
func (ps *ProjectileSystem) moveProjectile(e engine.Entity, proj *Projectile, dt float64) {
	posComp, hasPos := ps.world.GetComponent(e, "position")
	velComp, hasVel := ps.world.GetComponent(e, "velocity")

	if !hasPos || !hasVel {
		return
	}

	pos := posComp.(*engine.Position)
	vel := velComp.(*engine.Velocity)

	pos.X += vel.VX * dt
	pos.Y += vel.VY * dt
}

// checkCollisions tests the projectile against all potential targets.
func (ps *ProjectileSystem) checkCollisions(projectileEntity engine.Entity, proj *Projectile) {
	projPos, projBox := ps.getProjectileBounds(projectileEntity)
	if projPos == nil {
		return
	}

	ps.world.ForEachEntity(func(target engine.Entity) {
		if target == projectileEntity {
			return
		}
		ps.checkTargetCollision(projectileEntity, proj, projPos, projBox, target)
	})
}

// getProjectileBounds retrieves position and bounding box for a projectile.
func (ps *ProjectileSystem) getProjectileBounds(e engine.Entity) (*engine.Position, *BoundingBox) {
	posComp, hasPos := ps.world.GetComponent(e, "position")
	if !hasPos {
		return nil, nil
	}
	pos := posComp.(*engine.Position)

	boxComp, hasBox := ps.world.GetComponent(e, "boundingbox")
	if hasBox {
		return pos, boxComp.(*BoundingBox)
	}
	// Default small hitbox for projectiles
	return pos, &BoundingBox{X: pos.X - 2, Y: pos.Y - 2, Width: 4, Height: 4}
}

// checkTargetCollision tests collision between a projectile and a target entity.
func (ps *ProjectileSystem) checkTargetCollision(projectileEntity engine.Entity, proj *Projectile, projPos *engine.Position, projBox *BoundingBox, target engine.Entity) {
	if !ps.isValidTarget(proj, target) {
		return
	}

	targetPos, targetBox := ps.getTargetBounds(target)
	if targetPos == nil {
		return
	}

	if ps.boxesOverlap(projPos, projBox, targetPos, targetBox) {
		ps.handleHit(projectileEntity, target, proj.Damage)
	}
}

// isValidTarget returns true if the target can be hit by the projectile.
func (ps *ProjectileSystem) isValidTarget(proj *Projectile, target engine.Entity) bool {
	tagComp, hasTag := ps.world.GetComponent(target, "collisiontag")
	if !hasTag {
		return false
	}
	tag := tagComp.(*CollisionTag)

	// Player projectiles hit enemies, enemy projectiles hit player
	if proj.OwnerType == "player" && tag.Tag != "enemy" {
		return false
	}
	if proj.OwnerType == "enemy" && tag.Tag != "player" {
		return false
	}
	return true
}

// getTargetBounds retrieves position and bounding box for a target entity.
func (ps *ProjectileSystem) getTargetBounds(target engine.Entity) (*engine.Position, *BoundingBox) {
	posComp, hasPos := ps.world.GetComponent(target, "position")
	if !hasPos {
		return nil, nil
	}
	pos := posComp.(*engine.Position)

	boxComp, hasBox := ps.world.GetComponent(target, "boundingbox")
	if hasBox {
		return pos, boxComp.(*BoundingBox)
	}
	// Default hitbox for entities
	return pos, &BoundingBox{X: pos.X - 8, Y: pos.Y - 8, Width: 16, Height: 16}
}

// boxesOverlap checks AABB collision between two entities.
func (ps *ProjectileSystem) boxesOverlap(posA *engine.Position, boxA *BoundingBox, posB *engine.Position, boxB *BoundingBox) bool {
	return CheckAABBCollision(
		posA.X+boxA.X, posA.Y+boxA.Y, boxA.Width, boxA.Height,
		posB.X+boxB.X, posB.Y+boxB.Y, boxB.Width, boxB.Height,
	)
}

// handleHit processes a collision between projectile and target.
func (ps *ProjectileSystem) handleHit(projectile, target engine.Entity, damage float64) {
	ps.toRemove = append(ps.toRemove, projectile)
	if ps.onHit != nil {
		ps.onHit(projectile, target, damage)
	}
}

// SpawnProjectile creates a new projectile entity.
func (ps *ProjectileSystem) SpawnProjectile(x, y, angle, speed, damage float64, ownerType string, lifetime float64) engine.Entity {
	e := ps.world.CreateEntity()

	ps.world.AddComponent(e, "position", &engine.Position{X: x, Y: y})
	ps.world.AddComponent(e, "velocity", &engine.Velocity{
		VX: math.Cos(angle) * speed,
		VY: math.Sin(angle) * speed,
	})
	ps.world.AddComponent(e, "projectile", &Projectile{
		Damage:    damage,
		Speed:     speed,
		OwnerType: ownerType,
		Lifetime:  lifetime,
		MaxLife:   lifetime,
	})
	ps.world.AddComponent(e, "boundingbox", &BoundingBox{
		X: -2, Y: -2, Width: 4, Height: 4,
	})

	return e
}

// CheckAABBCollision returns true if two axis-aligned boxes overlap.
func CheckAABBCollision(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax < bx+bw &&
		ax+aw > bx &&
		ay < by+bh &&
		ay+ah > by
}

// ProjectileCount returns the number of active projectiles.
func (ps *ProjectileSystem) ProjectileCount() int {
	count := 0
	ps.world.ForEachEntity(func(e engine.Entity) {
		_, hasProj := ps.world.GetComponent(e, "projectile")
		if hasProj {
			count++
		}
	})
	return count
}
