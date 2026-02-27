// Package combat provides weapons, damage calculation, hit detection,
// and status effects.
package combat

// WeaponType identifies a weapon class.
type WeaponType int

const (
	WeaponPrimary WeaponType = iota
	WeaponSecondary
	WeaponMissile
	WeaponBomb
)

// Weapon represents a ship weapon.
type Weapon struct {
	Type     WeaponType
	Damage   float64
	Cooldown float64
	timer    float64
}

// NewWeapon creates a new weapon of the given type.
func NewWeapon(wt WeaponType, damage, cooldown float64) *Weapon {
	return &Weapon{Type: wt, Damage: damage, Cooldown: cooldown}
}

// CanFire returns true if the weapon cooldown has elapsed.
func (w *Weapon) CanFire() bool {
	return w.timer <= 0
}

// Fire triggers the weapon and resets the cooldown.
func (w *Weapon) Fire() {
	if w.CanFire() {
		w.timer = w.Cooldown
	}
}

// Update advances the weapon cooldown by dt seconds.
func (w *Weapon) Update(dt float64) {
	if w.timer > 0 {
		w.timer -= dt
	}
}

// StatusEffect represents a debuff applied to a ship.
type StatusEffect struct {
	Name     string
	Duration float64
	Active   bool
}

// NewStatusEffect creates a new status effect.
func NewStatusEffect(name string, duration float64) *StatusEffect {
	return &StatusEffect{Name: name, Duration: duration, Active: true}
}

// Update advances the effect timer.
func (se *StatusEffect) Update(dt float64) {
	if se.Active {
		se.Duration -= dt
		if se.Duration <= 0 {
			se.Active = false
		}
	}
}
