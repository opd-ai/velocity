package world

import "testing"

func TestNewWeather(t *testing.T) {
	w := NewWeather("solar flare", "scifi")
	if w == nil {
		t.Fatal("NewWeather() returned nil")
	}
	if w.Name != "solar flare" {
		t.Errorf("Name = %q, want %q", w.Name, "solar flare")
	}
	if w.GenreID != "scifi" {
		t.Errorf("GenreID = %q, want %q", w.GenreID, "scifi")
	}
	if w.Active {
		t.Error("new weather should not be active")
	}
	if w.Intensity != 0 {
		t.Errorf("Intensity = %f, want 0", w.Intensity)
	}
}

func TestWeatherSetGenre(t *testing.T) {
	w := NewWeather("ion storm", "scifi")

	genres := []string{"fantasy", "horror", "cyberpunk", "postapoc"}
	for _, genre := range genres {
		w.SetGenre(genre)
		if w.GenreID != genre {
			t.Errorf("after SetGenre(%q), GenreID = %q", genre, w.GenreID)
		}
	}
}

func TestWeatherStruct(t *testing.T) {
	w := &Weather{
		Name:      "void fog",
		GenreID:   "horror",
		Active:    true,
		Intensity: 0.75,
	}

	if w.Name != "void fog" {
		t.Errorf("Name = %q, want %q", w.Name, "void fog")
	}
	if w.GenreID != "horror" {
		t.Errorf("GenreID = %q, want %q", w.GenreID, "horror")
	}
	if !w.Active {
		t.Error("Active should be true")
	}
	if w.Intensity != 0.75 {
		t.Errorf("Intensity = %f, want 0.75", w.Intensity)
	}
}

func TestNewObjective(t *testing.T) {
	o := NewObjective("Destroy all enemies")
	if o == nil {
		t.Fatal("NewObjective() returned nil")
	}
	if o.Description != "Destroy all enemies" {
		t.Errorf("Description = %q, want %q", o.Description, "Destroy all enemies")
	}
	if o.Completed {
		t.Error("new objective should not be completed")
	}
}

func TestObjectiveComplete(t *testing.T) {
	o := NewObjective("Survive for 60 seconds")

	if o.Completed {
		t.Error("objective should start incomplete")
	}

	o.Complete()

	if !o.Completed {
		t.Error("objective should be completed after Complete()")
	}
}

func TestObjectiveStruct(t *testing.T) {
	o := &Objective{
		Description: "Custom objective",
		Completed:   true,
	}

	if o.Description != "Custom objective" {
		t.Errorf("Description = %q, want %q", o.Description, "Custom objective")
	}
	if !o.Completed {
		t.Error("Completed should be true")
	}
}

func TestLootDropStruct(t *testing.T) {
	loot := LootDrop{
		Type: "health",
		X:    100.5,
		Y:    200.25,
	}

	if loot.Type != "health" {
		t.Errorf("Type = %q, want %q", loot.Type, "health")
	}
	if loot.X != 100.5 {
		t.Errorf("X = %f, want 100.5", loot.X)
	}
	if loot.Y != 200.25 {
		t.Errorf("Y = %f, want 200.25", loot.Y)
	}
}

func TestNewEconomy(t *testing.T) {
	e := NewEconomy()
	if e == nil {
		t.Fatal("NewEconomy() returned nil")
	}
	if e.Credits != 0 {
		t.Errorf("new economy should have 0 credits, got %d", e.Credits)
	}
}

func TestEconomyAddCredits(t *testing.T) {
	e := NewEconomy()

	e.AddCredits(100)
	if e.Credits != 100 {
		t.Errorf("Credits = %d, want 100", e.Credits)
	}

	e.AddCredits(50)
	if e.Credits != 150 {
		t.Errorf("Credits = %d, want 150", e.Credits)
	}
}

func TestEconomyAddCreditsIgnoresNonPositive(t *testing.T) {
	e := NewEconomy()
	e.AddCredits(100)

	e.AddCredits(0)
	if e.Credits != 100 {
		t.Errorf("Credits = %d, want 100 (0 should be ignored)", e.Credits)
	}

	e.AddCredits(-50)
	if e.Credits != 100 {
		t.Errorf("Credits = %d, want 100 (negative should be ignored)", e.Credits)
	}
}

func TestEconomySpend(t *testing.T) {
	e := NewEconomy()
	e.AddCredits(100)

	ok := e.Spend(30)
	if !ok {
		t.Error("Spend(30) should succeed with 100 credits")
	}
	if e.Credits != 70 {
		t.Errorf("Credits = %d, want 70", e.Credits)
	}

	ok = e.Spend(70)
	if !ok {
		t.Error("Spend(70) should succeed with 70 credits")
	}
	if e.Credits != 0 {
		t.Errorf("Credits = %d, want 0", e.Credits)
	}
}

func TestEconomySpendInsufficientFunds(t *testing.T) {
	e := NewEconomy()
	e.AddCredits(50)

	ok := e.Spend(100)
	if ok {
		t.Error("Spend(100) should fail with 50 credits")
	}
	if e.Credits != 50 {
		t.Errorf("Credits = %d, want 50 (unchanged)", e.Credits)
	}
}

func TestEconomySpendNonPositive(t *testing.T) {
	e := NewEconomy()
	e.AddCredits(100)

	ok := e.Spend(0)
	if ok {
		t.Error("Spend(0) should return false")
	}
	if e.Credits != 100 {
		t.Errorf("Credits = %d, want 100 (unchanged)", e.Credits)
	}

	ok = e.Spend(-50)
	if ok {
		t.Error("Spend(-50) should return false")
	}
	if e.Credits != 100 {
		t.Errorf("Credits = %d, want 100 (unchanged)", e.Credits)
	}
}

func TestEconomyStruct(t *testing.T) {
	e := &Economy{Credits: 500}

	if e.Credits != 500 {
		t.Errorf("Credits = %d, want 500", e.Credits)
	}
}

func TestWeatherTypes(t *testing.T) {
	weatherTypes := []struct {
		name    string
		genreID string
	}{
		{"solar flare", "scifi"},
		{"ion storm", "scifi"},
		{"nebula interference", "scifi"},
		{"debris field", "postapoc"},
		{"void fog", "horror"},
		{"data storm", "cyberpunk"},
		{"arcane tempest", "fantasy"},
	}

	for _, wt := range weatherTypes {
		w := NewWeather(wt.name, wt.genreID)
		if w.Name != wt.name {
			t.Errorf("NewWeather(%q, %q).Name = %q", wt.name, wt.genreID, w.Name)
		}
		if w.GenreID != wt.genreID {
			t.Errorf("NewWeather(%q, %q).GenreID = %q", wt.name, wt.genreID, w.GenreID)
		}
	}
}
