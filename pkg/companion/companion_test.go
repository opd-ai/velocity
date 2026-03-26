package companion

import "testing"

func TestNewWingman(t *testing.T) {
	w := NewWingman()
	if w == nil {
		t.Fatal("NewWingman() returned nil")
	}
	if w.genreID != "scifi" {
		t.Errorf("default genreID = %q, want %q", w.genreID, "scifi")
	}
	if w.active {
		t.Error("new wingman should not be active")
	}
}

func TestWingmanSetGenre(t *testing.T) {
	w := NewWingman()

	genres := []string{"fantasy", "horror", "cyberpunk", "postapoc", "scifi"}
	for _, genre := range genres {
		w.SetGenre(genre)
		if w.genreID != genre {
			t.Errorf("after SetGenre(%q), genreID = %q", genre, w.genreID)
		}
	}
}

func TestWingmanActivate(t *testing.T) {
	w := NewWingman()

	if w.active {
		t.Error("wingman should start inactive")
	}

	w.Activate()
	if !w.active {
		t.Error("wingman should be active after Activate()")
	}
}

func TestWingmanDeactivate(t *testing.T) {
	w := NewWingman()
	w.Activate()

	if !w.active {
		t.Fatal("wingman should be active after Activate()")
	}

	w.Deactivate()
	if w.active {
		t.Error("wingman should be inactive after Deactivate()")
	}
}

func TestWingmanUpdate(t *testing.T) {
	w := NewWingman()
	w.Activate()

	// Update should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update() panicked: %v", r)
		}
	}()

	w.Update(1.0 / 60.0) // One frame at 60 FPS
	w.Update(0.5)        // Half second
}

func TestWingmanActivateDeactivateCycle(t *testing.T) {
	w := NewWingman()

	// Multiple activations/deactivations should work
	for i := 0; i < 5; i++ {
		w.Activate()
		if !w.active {
			t.Errorf("iteration %d: wingman should be active", i)
		}
		w.Deactivate()
		if w.active {
			t.Errorf("iteration %d: wingman should be inactive", i)
		}
	}
}

func TestWingmanUpdateWhileInactive(t *testing.T) {
	w := NewWingman()
	// Don't activate

	// Update should not panic even when inactive
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update() panicked while inactive: %v", r)
		}
	}()

	w.Update(1.0 / 60.0)
}

func TestWingmanStruct(t *testing.T) {
	w := &Wingman{
		genreID: "horror",
		active:  true,
	}

	if w.genreID != "horror" {
		t.Errorf("genreID = %q, want %q", w.genreID, "horror")
	}
	if !w.active {
		t.Error("active should be true")
	}
}
