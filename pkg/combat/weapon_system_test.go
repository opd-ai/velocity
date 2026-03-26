package combat

import (
	"testing"

	"github.com/opd-ai/velocity/pkg/engine"
)

// mockFireProvider implements FireProvider for testing.
type mockFireProvider struct {
	firePressed bool
}

func (m *mockFireProvider) IsFirePressed() bool {
	return m.firePressed
}

func TestNewWeaponComponent(t *testing.T) {
	weapon := NewWeapon(WeaponPrimary, 10.0, 0.2)
	wc := NewWeaponComponent(weapon)

	if wc.Primary != weapon {
		t.Error("Primary weapon not set correctly")
	}
	if wc.Secondary != nil {
		t.Error("Secondary weapon should be nil")
	}
}

func TestWeaponSystemCreation(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	if ws.world != world {
		t.Error("World not set correctly")
	}
	if ws.projectiles != projSys {
		t.Error("Projectile system not set correctly")
	}
}

func TestWeaponSystemSetFireProvider(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	provider := &mockFireProvider{}
	ws.SetFireProvider(provider)

	if ws.input != provider {
		t.Error("Fire provider not set correctly")
	}
}

func TestWeaponSystemUpdate_CooldownTicks(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	// Create entity with weapon
	e := world.CreateEntity()
	weapon := NewWeapon(WeaponPrimary, 10.0, 1.0)
	weapon.Fire() // Start cooldown
	world.AddComponent(e, "weapon", NewWeaponComponent(weapon))

	if weapon.CanFire() {
		t.Error("Weapon should be on cooldown")
	}

	// Update should tick cooldown
	ws.Update(0.5)
	if weapon.CanFire() {
		t.Error("Weapon should still be on cooldown after 0.5s")
	}

	ws.Update(0.6)
	if !weapon.CanFire() {
		t.Error("Weapon should be off cooldown after 1.1s")
	}
}

func TestWeaponSystemUpdate_PlayerFiring(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	provider := &mockFireProvider{firePressed: true}
	ws.SetFireProvider(provider)

	// Create player entity with weapon and required components
	player := world.CreateEntity()
	world.AddComponent(player, "position", &engine.Position{X: 100, Y: 100})
	world.AddComponent(player, "rotation", &engine.Rotation{Angle: 0})
	world.AddComponent(player, "collisiontag", &CollisionTag{Tag: "player"})
	world.AddComponent(player, "weapon", NewWeaponComponent(NewWeapon(WeaponPrimary, 10.0, 0.2)))

	// Should spawn projectile when fire is pressed
	ws.Update(1.0 / 60.0)

	count := projSys.ProjectileCount()
	if count != 1 {
		t.Errorf("Expected 1 projectile, got %d", count)
	}
}

func TestWeaponSystemUpdate_NoFireWhenNotPressed(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	provider := &mockFireProvider{firePressed: false}
	ws.SetFireProvider(provider)

	// Create player entity
	player := world.CreateEntity()
	world.AddComponent(player, "position", &engine.Position{X: 100, Y: 100})
	world.AddComponent(player, "rotation", &engine.Rotation{Angle: 0})
	world.AddComponent(player, "collisiontag", &CollisionTag{Tag: "player"})
	world.AddComponent(player, "weapon", NewWeaponComponent(NewWeapon(WeaponPrimary, 10.0, 0.2)))

	ws.Update(1.0 / 60.0)

	count := projSys.ProjectileCount()
	if count != 0 {
		t.Errorf("Expected 0 projectiles (fire not pressed), got %d", count)
	}
}

func TestWeaponSystemUpdate_CooldownPreventsRapidFire(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	provider := &mockFireProvider{firePressed: true}
	ws.SetFireProvider(provider)

	// Create player entity with slow-firing weapon
	player := world.CreateEntity()
	world.AddComponent(player, "position", &engine.Position{X: 100, Y: 100})
	world.AddComponent(player, "rotation", &engine.Rotation{Angle: 0})
	world.AddComponent(player, "collisiontag", &CollisionTag{Tag: "player"})
	world.AddComponent(player, "weapon", NewWeaponComponent(NewWeapon(WeaponPrimary, 10.0, 1.0))) // 1 second cooldown

	// First shot
	ws.Update(1.0 / 60.0)
	count1 := projSys.ProjectileCount()

	// Immediate second shot should be blocked by cooldown
	ws.Update(1.0 / 60.0)
	count2 := projSys.ProjectileCount()

	if count1 != 1 || count2 != 1 {
		t.Errorf("Expected 1 projectile with cooldown blocking second, got %d then %d", count1, count2)
	}
}

func TestWeaponSystemUpdate_NonPlayerDoesNotFire(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	provider := &mockFireProvider{firePressed: true}
	ws.SetFireProvider(provider)

	// Create enemy entity with weapon
	enemy := world.CreateEntity()
	world.AddComponent(enemy, "position", &engine.Position{X: 100, Y: 100})
	world.AddComponent(enemy, "rotation", &engine.Rotation{Angle: 0})
	world.AddComponent(enemy, "collisiontag", &CollisionTag{Tag: "enemy"})
	world.AddComponent(enemy, "weapon", NewWeaponComponent(NewWeapon(WeaponPrimary, 10.0, 0.2)))

	// Enemy should not fire from player input
	ws.Update(1.0 / 60.0)

	count := projSys.ProjectileCount()
	if count != 0 {
		t.Errorf("Expected 0 projectiles (enemy should not fire), got %d", count)
	}
}

func TestFireAtTarget(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	// Create enemy entity
	enemy := world.CreateEntity()
	world.AddComponent(enemy, "position", &engine.Position{X: 100, Y: 100})

	weapon := NewWeapon(WeaponPrimary, 10.0, 0.5)

	// Fire at target
	ws.FireAtTarget(enemy, weapon, 200, 100, "enemy")

	count := projSys.ProjectileCount()
	if count != 1 {
		t.Errorf("Expected 1 projectile, got %d", count)
	}

	// Verify projectile direction (should be to the right)
	world.ForEachEntity(func(e engine.Entity) {
		velComp, hasVel := world.GetComponent(e, "velocity")
		if !hasVel {
			return
		}
		vel := velComp.(*engine.Velocity)
		if vel.VX <= 0 {
			t.Errorf("Projectile should move right (VX > 0), got VX=%f", vel.VX)
		}
	})
}

func TestFireAtTarget_CooldownRespected(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	enemy := world.CreateEntity()
	world.AddComponent(enemy, "position", &engine.Position{X: 100, Y: 100})

	weapon := NewWeapon(WeaponPrimary, 10.0, 1.0)

	// First shot
	ws.FireAtTarget(enemy, weapon, 200, 100, "enemy")
	count1 := projSys.ProjectileCount()

	// Second shot immediately - should be blocked
	ws.FireAtTarget(enemy, weapon, 200, 100, "enemy")
	count2 := projSys.ProjectileCount()

	if count1 != 1 || count2 != 1 {
		t.Errorf("Cooldown not respected: got %d then %d projectiles", count1, count2)
	}
}

func TestFireAtTarget_NilWeapon(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	enemy := world.CreateEntity()
	world.AddComponent(enemy, "position", &engine.Position{X: 100, Y: 100})

	// Should not panic with nil weapon
	ws.FireAtTarget(enemy, nil, 200, 100, "enemy")

	count := projSys.ProjectileCount()
	if count != 0 {
		t.Errorf("Expected 0 projectiles with nil weapon, got %d", count)
	}
}

func TestFireAtTarget_MissingPosition(t *testing.T) {
	world := engine.NewWorld()
	projSys := NewProjectileSystem(world)
	ws := NewWeaponSystem(world, projSys)

	enemy := world.CreateEntity()
	// No position component

	weapon := NewWeapon(WeaponPrimary, 10.0, 0.5)

	// Should not panic without position
	ws.FireAtTarget(enemy, weapon, 200, 100, "enemy")

	count := projSys.ProjectileCount()
	if count != 0 {
		t.Errorf("Expected 0 projectiles without position, got %d", count)
	}
}
