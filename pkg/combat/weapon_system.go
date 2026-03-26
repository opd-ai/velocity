// Package combat provides weapons, damage calculation, hit detection,
// and status effects.
package combat

import (
	"math"

	"github.com/opd-ai/velocity/pkg/engine"
)

// WeaponComponent attaches a weapon to an entity.
type WeaponComponent struct {
	Primary   *Weapon
	Secondary *Weapon
}

// NewWeaponComponent creates a weapon component with a primary weapon.
func NewWeaponComponent(primary *Weapon) *WeaponComponent {
	return &WeaponComponent{Primary: primary}
}

// FireProvider is an interface for systems that track fire input.
type FireProvider interface {
	IsFirePressed() bool
}

// WeaponSystem manages weapon firing and projectile spawning.
type WeaponSystem struct {
	world       *engine.World
	projectiles *ProjectileSystem
	input       FireProvider
}

// NewWeaponSystem creates a new weapon system.
func NewWeaponSystem(world *engine.World, projectiles *ProjectileSystem) *WeaponSystem {
	return &WeaponSystem{
		world:       world,
		projectiles: projectiles,
	}
}

// SetFireProvider sets the input provider for fire commands.
func (ws *WeaponSystem) SetFireProvider(provider FireProvider) {
	ws.input = provider
}

// Update processes weapon cooldowns and handles firing for player entity.
func (ws *WeaponSystem) Update(dt float64) {
	ws.world.ForEachEntity(func(e engine.Entity) {
		ws.updateEntityWeapon(e, dt)
	})
}

// updateEntityWeapon handles weapon logic for a single entity.
func (ws *WeaponSystem) updateEntityWeapon(e engine.Entity, dt float64) {
	weaponComp, hasWeapon := ws.world.GetComponent(e, "weapon")
	if !hasWeapon {
		return
	}

	weapon := weaponComp.(*WeaponComponent)

	// Update weapon cooldowns
	if weapon.Primary != nil {
		weapon.Primary.Update(dt)
	}
	if weapon.Secondary != nil {
		weapon.Secondary.Update(dt)
	}

	// Check if this is the player and fire is pressed
	tagComp, hasTag := ws.world.GetComponent(e, "collisiontag")
	if !hasTag {
		return
	}

	tag := tagComp.(*CollisionTag)
	if tag.Tag != "player" {
		return
	}

	// Handle player firing
	if ws.input != nil && ws.input.IsFirePressed() {
		ws.tryFire(e, weapon.Primary, "player")
	}
}

// tryFire attempts to fire a weapon from an entity.
func (ws *WeaponSystem) tryFire(e engine.Entity, weapon *Weapon, ownerType string) {
	if weapon == nil || !weapon.CanFire() {
		return
	}

	posComp, hasPos := ws.world.GetComponent(e, "position")
	rotComp, hasRot := ws.world.GetComponent(e, "rotation")

	if !hasPos || !hasRot {
		return
	}

	pos := posComp.(*engine.Position)
	rot := rotComp.(*engine.Rotation)

	// Spawn offset from entity center
	offset := 12.0
	spawnX := pos.X + math.Cos(rot.Angle)*offset
	spawnY := pos.Y + math.Sin(rot.Angle)*offset

	// Spawn projectile
	projectileSpeed := 400.0
	projectileLifetime := 2.0

	ws.projectiles.SpawnProjectile(
		spawnX, spawnY,
		rot.Angle,
		projectileSpeed,
		weapon.Damage,
		ownerType,
		projectileLifetime,
	)

	weapon.Fire()
}

// FireAtTarget fires a weapon from an entity toward a target position.
func (ws *WeaponSystem) FireAtTarget(e engine.Entity, weapon *Weapon, targetX, targetY float64, ownerType string) {
	if weapon == nil || !weapon.CanFire() {
		return
	}

	posComp, hasPos := ws.world.GetComponent(e, "position")
	if !hasPos {
		return
	}

	pos := posComp.(*engine.Position)

	// Calculate angle to target
	dx := targetX - pos.X
	dy := targetY - pos.Y
	angle := math.Atan2(dy, dx)

	// Spawn offset from entity center
	offset := 12.0
	spawnX := pos.X + math.Cos(angle)*offset
	spawnY := pos.Y + math.Sin(angle)*offset

	// Spawn projectile
	projectileSpeed := 300.0
	projectileLifetime := 2.0

	ws.projectiles.SpawnProjectile(
		spawnX, spawnY,
		angle,
		projectileSpeed,
		weapon.Damage,
		ownerType,
		projectileLifetime,
	)

	weapon.Fire()
}
