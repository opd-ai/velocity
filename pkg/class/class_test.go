package class

import "testing"

func TestDefaultHulls(t *testing.T) {
	hulls := DefaultHulls()

	if len(hulls) != 4 {
		t.Fatalf("DefaultHulls() returned %d hulls, want 4", len(hulls))
	}

	expectedNames := []string{"scout", "interceptor", "gunship", "carrier"}
	for i, name := range expectedNames {
		if hulls[i].Name != name {
			t.Errorf("hulls[%d].Name = %q, want %q", i, hulls[i].Name, name)
		}
	}
}

func TestHullClassScout(t *testing.T) {
	hulls := DefaultHulls()
	scout := hulls[0]

	if scout.Name != "scout" {
		t.Errorf("Name = %q, want %q", scout.Name, "scout")
	}
	if scout.Health != 60 {
		t.Errorf("Health = %f, want 60", scout.Health)
	}
	if scout.Speed != 300 {
		t.Errorf("Speed = %f, want 300", scout.Speed)
	}
	if scout.Armor != 2 {
		t.Errorf("Armor = %f, want 2", scout.Armor)
	}
	if scout.Slots != 2 {
		t.Errorf("Slots = %d, want 2", scout.Slots)
	}
}

func TestHullClassCarrier(t *testing.T) {
	hulls := DefaultHulls()
	carrier := hulls[3]

	if carrier.Name != "carrier" {
		t.Errorf("Name = %q, want %q", carrier.Name, "carrier")
	}
	if carrier.Health != 200 {
		t.Errorf("Health = %f, want 200", carrier.Health)
	}
	if carrier.Speed != 120 {
		t.Errorf("Speed = %f, want 120", carrier.Speed)
	}
	if carrier.Armor != 10 {
		t.Errorf("Armor = %f, want 10", carrier.Armor)
	}
	if carrier.Slots != 5 {
		t.Errorf("Slots = %d, want 5", carrier.Slots)
	}
}

func TestHullClassTradeoffs(t *testing.T) {
	hulls := DefaultHulls()

	// Scout should be fastest but weakest
	scout := hulls[0]
	carrier := hulls[3]

	if scout.Speed <= carrier.Speed {
		t.Error("Scout should be faster than carrier")
	}
	if scout.Health >= carrier.Health {
		t.Error("Scout should have less health than carrier")
	}
	if scout.Armor >= carrier.Armor {
		t.Error("Scout should have less armor than carrier")
	}
	if scout.Slots >= carrier.Slots {
		t.Error("Scout should have fewer slots than carrier")
	}
}

func TestUpgradeNodeCanUpgrade(t *testing.T) {
	tests := []struct {
		level    int
		maxLevel int
		expected bool
	}{
		{0, 5, true},
		{4, 5, true},
		{5, 5, false},
		{0, 0, false},
		{10, 5, false}, // Over max
	}

	for _, tt := range tests {
		node := UpgradeNode{Level: tt.level, MaxLevel: tt.maxLevel}
		if got := node.CanUpgrade(); got != tt.expected {
			t.Errorf("UpgradeNode{Level: %d, MaxLevel: %d}.CanUpgrade() = %v, want %v",
				tt.level, tt.maxLevel, got, tt.expected)
		}
	}
}

func TestUpgradeNodeUpgrade(t *testing.T) {
	node := UpgradeNode{Name: "weapons", Level: 0, MaxLevel: 3}

	node.Upgrade()
	if node.Level != 1 {
		t.Errorf("Level = %d after first upgrade, want 1", node.Level)
	}

	node.Upgrade()
	node.Upgrade()
	if node.Level != 3 {
		t.Errorf("Level = %d after 3 upgrades, want 3", node.Level)
	}

	// Should not exceed max
	node.Upgrade()
	if node.Level != 3 {
		t.Errorf("Level = %d after upgrade at max, want 3 (no change)", node.Level)
	}
}

func TestUpgradeNodeStruct(t *testing.T) {
	node := UpgradeNode{
		Name:     "shields",
		Level:    2,
		MaxLevel: 5,
	}

	if node.Name != "shields" {
		t.Errorf("Name = %q, want %q", node.Name, "shields")
	}
	if node.Level != 2 {
		t.Errorf("Level = %d, want 2", node.Level)
	}
	if node.MaxLevel != 5 {
		t.Errorf("MaxLevel = %d, want 5", node.MaxLevel)
	}
}

func TestHullClassStruct(t *testing.T) {
	hull := HullClass{
		Name:   "custom",
		Health: 150,
		Speed:  220,
		Armor:  5,
		Slots:  3,
	}

	if hull.Name != "custom" {
		t.Errorf("Name = %q, want %q", hull.Name, "custom")
	}
	if hull.Health != 150 {
		t.Errorf("Health = %f, want 150", hull.Health)
	}
	if hull.Speed != 220 {
		t.Errorf("Speed = %f, want 220", hull.Speed)
	}
	if hull.Armor != 5 {
		t.Errorf("Armor = %f, want 5", hull.Armor)
	}
	if hull.Slots != 3 {
		t.Errorf("Slots = %d, want 3", hull.Slots)
	}
}
