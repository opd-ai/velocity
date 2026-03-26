// Package combat provides weapons, damage calculation, hit detection,
// and status effects.
package combat

import (
	"github.com/opd-ai/velocity/pkg/engine"
)

// DamageEvent represents damage dealt to an entity.
type DamageEvent struct {
	Target     engine.Entity
	Amount     float64
	Source     engine.Entity
	SourceType string
}

// DeathEvent represents an entity being destroyed.
type DeathEvent struct {
	Entity   engine.Entity
	Position engine.Position
	Score    int
}

// DamageSystem handles damage application and entity destruction.
type DamageSystem struct {
	world         *engine.World
	pendingDamage []DamageEvent
	pendingDeaths []DeathEvent
	onDeath       func(DeathEvent)
}

// NewDamageSystem creates a new damage system.
func NewDamageSystem(world *engine.World) *DamageSystem {
	return &DamageSystem{
		world:         world,
		pendingDamage: make([]DamageEvent, 0, 16),
		pendingDeaths: make([]DeathEvent, 0, 8),
	}
}

// SetDeathCallback sets the callback for when an entity is destroyed.
func (ds *DamageSystem) SetDeathCallback(fn func(DeathEvent)) {
	ds.onDeath = fn
}

// QueueDamage queues damage to be applied on next update.
func (ds *DamageSystem) QueueDamage(target, source engine.Entity, amount float64, sourceType string) {
	ds.pendingDamage = append(ds.pendingDamage, DamageEvent{
		Target:     target,
		Amount:     amount,
		Source:     source,
		SourceType: sourceType,
	})
}

// ApplyDamage immediately applies damage to an entity.
func (ds *DamageSystem) ApplyDamage(target engine.Entity, amount float64) bool {
	healthComp, hasHealth := ds.world.GetComponent(target, "health")
	if !hasHealth {
		return false
	}

	health := healthComp.(*Health)
	health.Current -= amount

	if health.Current <= 0 {
		health.Current = 0
		return true // Entity died
	}

	return false
}

// Update processes all pending damage and removes dead entities.
func (ds *DamageSystem) Update(dt float64) {
	ds.pendingDeaths = ds.pendingDeaths[:0]

	// Process pending damage
	for _, event := range ds.pendingDamage {
		if ds.ApplyDamage(event.Target, event.Amount) {
			// Entity died
			ds.handleDeath(event.Target)
		}
	}

	// Clear pending damage
	ds.pendingDamage = ds.pendingDamage[:0]

	// Process deaths
	for _, death := range ds.pendingDeaths {
		if ds.onDeath != nil {
			ds.onDeath(death)
		}
		ds.world.RemoveEntity(death.Entity)
	}
}

// handleDeath prepares an entity for removal.
func (ds *DamageSystem) handleDeath(entity engine.Entity) {
	var pos engine.Position

	posComp, hasPos := ds.world.GetComponent(entity, "position")
	if hasPos {
		pos = *posComp.(*engine.Position)
	}

	// Calculate score based on entity type
	score := 100 // Default score
	tagComp, hasTag := ds.world.GetComponent(entity, "collisiontag")
	if hasTag {
		tag := tagComp.(*CollisionTag)
		if tag.Tag == "enemy" {
			score = 100
		}
	}

	ds.pendingDeaths = append(ds.pendingDeaths, DeathEvent{
		Entity:   entity,
		Position: pos,
		Score:    score,
	})
}

// IsEntityAlive returns true if the entity has health > 0.
func (ds *DamageSystem) IsEntityAlive(entity engine.Entity) bool {
	healthComp, hasHealth := ds.world.GetComponent(entity, "health")
	if !hasHealth {
		return true // No health component = invulnerable
	}

	health := healthComp.(*Health)
	return health.Current > 0
}

// GetHealthPercent returns the health percentage (0-1) for an entity.
func (ds *DamageSystem) GetHealthPercent(entity engine.Entity) float64 {
	healthComp, hasHealth := ds.world.GetComponent(entity, "health")
	if !hasHealth {
		return 1.0
	}

	health := healthComp.(*Health)
	if health.Max <= 0 {
		return 1.0
	}

	return health.Current / health.Max
}

// Heal restores health to an entity.
func (ds *DamageSystem) Heal(entity engine.Entity, amount float64) {
	healthComp, hasHealth := ds.world.GetComponent(entity, "health")
	if !hasHealth {
		return
	}

	health := healthComp.(*Health)
	health.Current += amount
	if health.Current > health.Max {
		health.Current = health.Max
	}
}
