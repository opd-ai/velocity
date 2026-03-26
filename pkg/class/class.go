// Package class provides ship hull classes and skill/upgrade trees.
package class

// Ship hull class base stats - these values define the core balance between
// ship types and should be tuned via playtesting.
const (
	// Scout hull - fast, fragile reconnaissance ship
	ScoutHealth = 60.0
	ScoutSpeed  = 300.0
	ScoutArmor  = 2.0
	ScoutSlots  = 2

	// Interceptor hull - balanced fighter
	InterceptorHealth = 80.0
	InterceptorSpeed  = 260.0
	InterceptorArmor  = 3.0
	InterceptorSlots  = 3

	// Gunship hull - heavy assault ship
	GunshipHealth = 120.0
	GunshipSpeed  = 180.0
	GunshipArmor  = 6.0
	GunshipSlots  = 4

	// Carrier hull - slow capital ship
	CarrierHealth = 200.0
	CarrierSpeed  = 120.0
	CarrierArmor  = 10.0
	CarrierSlots  = 5
)

// HullClass represents a ship hull type with base stats.
type HullClass struct {
	Name   string
	Health float64
	Speed  float64
	Armor  float64
	Slots  int
}

// DefaultHulls returns the set of base hull classes.
func DefaultHulls() []HullClass {
	return []HullClass{
		{Name: "scout", Health: ScoutHealth, Speed: ScoutSpeed, Armor: ScoutArmor, Slots: ScoutSlots},
		{Name: "interceptor", Health: InterceptorHealth, Speed: InterceptorSpeed, Armor: InterceptorArmor, Slots: InterceptorSlots},
		{Name: "gunship", Health: GunshipHealth, Speed: GunshipSpeed, Armor: GunshipArmor, Slots: GunshipSlots},
		{Name: "carrier", Health: CarrierHealth, Speed: CarrierSpeed, Armor: CarrierArmor, Slots: CarrierSlots},
	}
}

// UpgradeNode represents a single node in the ship upgrade tree.
type UpgradeNode struct {
	Name     string
	Level    int
	MaxLevel int
}

// CanUpgrade returns true if the node has not reached max level.
func (u *UpgradeNode) CanUpgrade() bool {
	return u.Level < u.MaxLevel
}

// Upgrade increments the node level.
func (u *UpgradeNode) Upgrade() {
	if u.CanUpgrade() {
		u.Level++
	}
}
