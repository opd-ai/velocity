package combat

import (
	"testing"
)

func TestNewWeapon(t *testing.T) {
	w := NewWeapon(WeaponPrimary, 10, 0.2)

	if w == nil {
		t.Fatal("expected non-nil weapon")
	}
	if w.Type != WeaponPrimary {
		t.Error("type mismatch")
	}
	if w.Damage != 10 {
		t.Error("damage mismatch")
	}
	if w.Cooldown != 0.2 {
		t.Error("cooldown mismatch")
	}
}

func TestWeapon_CanFire(t *testing.T) {
	w := NewWeapon(WeaponPrimary, 10, 0.2)

	// Initially should be able to fire
	if !w.CanFire() {
		t.Error("expected CanFire to be true initially")
	}
}

func TestWeapon_Fire(t *testing.T) {
	w := NewWeapon(WeaponPrimary, 10, 0.2)

	w.Fire()

	// Should not be able to fire immediately after
	if w.CanFire() {
		t.Error("expected CanFire to be false after firing")
	}
}

func TestWeapon_Update(t *testing.T) {
	w := NewWeapon(WeaponPrimary, 10, 0.2)

	w.Fire()

	// Update past cooldown
	w.Update(0.3)

	if !w.CanFire() {
		t.Error("expected CanFire to be true after cooldown")
	}
}

func TestWeapon_PartialCooldown(t *testing.T) {
	w := NewWeapon(WeaponPrimary, 10, 0.2)

	w.Fire()
	w.Update(0.1) // Partial cooldown

	if w.CanFire() {
		t.Error("expected CanFire to be false during cooldown")
	}
}

func TestWeaponTypes(t *testing.T) {
	types := []WeaponType{
		WeaponPrimary,
		WeaponSecondary,
		WeaponMissile,
		WeaponBomb,
	}

	for _, wt := range types {
		w := NewWeapon(wt, 10, 0.2)
		if w.Type != wt {
			t.Errorf("expected type %d, got %d", wt, w.Type)
		}
	}
}

func TestNewStatusEffect(t *testing.T) {
	se := NewStatusEffect("slowed", 5.0)

	if se == nil {
		t.Fatal("expected non-nil status effect")
	}
	if se.Name != "slowed" {
		t.Error("name mismatch")
	}
	if se.Duration != 5.0 {
		t.Error("duration mismatch")
	}
	if !se.Active {
		t.Error("expected effect to be active")
	}
}

func TestStatusEffect_Update(t *testing.T) {
	se := NewStatusEffect("slowed", 0.5)

	se.Update(0.3)

	if !se.Active {
		t.Error("expected effect to still be active")
	}

	se.Update(0.3) // Total 0.6s, past duration

	if se.Active {
		t.Error("expected effect to expire")
	}
}

func TestStatusEffect_MultipleUpdates(t *testing.T) {
	se := NewStatusEffect("emp", 1.0)

	for i := 0; i < 100; i++ {
		se.Update(0.01)
	}

	if se.Active {
		t.Error("expected effect to expire after many updates")
	}
}

func TestWeapon_MultipleFires(t *testing.T) {
	w := NewWeapon(WeaponPrimary, 10, 0.1)

	// Fire several times with cooldown reset
	for i := 0; i < 5; i++ {
		if !w.CanFire() {
			t.Errorf("expected CanFire on iteration %d", i)
		}
		w.Fire()
		w.Update(0.2) // Past cooldown
	}
}

func BenchmarkWeapon_Update(b *testing.B) {
	w := NewWeapon(WeaponPrimary, 10, 0.1)
	w.Fire()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Update(1.0 / 60.0)
		if w.CanFire() {
			w.Fire()
		}
	}
}
