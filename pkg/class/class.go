// Package class provides ship hull classes and skill/upgrade trees.
package class

// HullClass represents a ship hull type with base stats.
type HullClass struct {
	Name    string
	Health  float64
	Speed   float64
	Armor   float64
	Slots   int
}

// DefaultHulls returns the set of base hull classes.
func DefaultHulls() []HullClass {
	return []HullClass{
		{Name: "scout", Health: 60, Speed: 300, Armor: 2, Slots: 2},
		{Name: "interceptor", Health: 80, Speed: 260, Armor: 3, Slots: 3},
		{Name: "gunship", Health: 120, Speed: 180, Armor: 6, Slots: 4},
		{Name: "carrier", Health: 200, Speed: 120, Armor: 10, Slots: 5},
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
