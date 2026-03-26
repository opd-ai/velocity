package visualtest

import "testing"

func TestCapture(t *testing.T) {
	snap := Capture("scifi", 800, 600)
	if snap == nil {
		t.Fatal("Capture() returned nil")
	}
	if snap.GenreID != "scifi" {
		t.Errorf("GenreID = %q, want %q", snap.GenreID, "scifi")
	}
	if snap.Width != 800 {
		t.Errorf("Width = %d, want 800", snap.Width)
	}
	if snap.Height != 600 {
		t.Errorf("Height = %d, want 600", snap.Height)
	}
}

func TestCaptureDifferentGenres(t *testing.T) {
	genres := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc"}

	for _, genre := range genres {
		snap := Capture(genre, 1920, 1080)
		if snap.GenreID != genre {
			t.Errorf("Capture(%q, ...).GenreID = %q", genre, snap.GenreID)
		}
	}
}

func TestCaptureDifferentResolutions(t *testing.T) {
	tests := []struct {
		width, height int
	}{
		{640, 480},
		{800, 600},
		{1280, 720},
		{1920, 1080},
		{2560, 1440},
	}

	for _, tt := range tests {
		snap := Capture("scifi", tt.width, tt.height)
		if snap.Width != tt.width {
			t.Errorf("Width = %d, want %d", snap.Width, tt.width)
		}
		if snap.Height != tt.height {
			t.Errorf("Height = %d, want %d", snap.Height, tt.height)
		}
	}
}

func TestCompareIdentical(t *testing.T) {
	a := &Snapshot{
		GenreID: "scifi",
		Width:   100,
		Height:  100,
		Data:    []byte{0, 1, 2, 3, 4, 5},
	}
	b := &Snapshot{
		GenreID: "scifi",
		Width:   100,
		Height:  100,
		Data:    []byte{0, 1, 2, 3, 4, 5},
	}

	if !Compare(a, b) {
		t.Error("Compare() should return true for identical snapshots")
	}
}

func TestCompareDifferentWidth(t *testing.T) {
	a := &Snapshot{Width: 100, Height: 100, Data: []byte{0, 1, 2}}
	b := &Snapshot{Width: 200, Height: 100, Data: []byte{0, 1, 2}}

	if Compare(a, b) {
		t.Error("Compare() should return false for different widths")
	}
}

func TestCompareDifferentHeight(t *testing.T) {
	a := &Snapshot{Width: 100, Height: 100, Data: []byte{0, 1, 2}}
	b := &Snapshot{Width: 100, Height: 200, Data: []byte{0, 1, 2}}

	if Compare(a, b) {
		t.Error("Compare() should return false for different heights")
	}
}

func TestCompareDifferentDataLength(t *testing.T) {
	a := &Snapshot{Width: 100, Height: 100, Data: []byte{0, 1, 2}}
	b := &Snapshot{Width: 100, Height: 100, Data: []byte{0, 1, 2, 3}}

	if Compare(a, b) {
		t.Error("Compare() should return false for different data lengths")
	}
}

func TestCompareDifferentDataContent(t *testing.T) {
	a := &Snapshot{Width: 100, Height: 100, Data: []byte{0, 1, 2, 3}}
	b := &Snapshot{Width: 100, Height: 100, Data: []byte{0, 1, 9, 3}}

	if Compare(a, b) {
		t.Error("Compare() should return false for different data content")
	}
}

func TestCompareEmptyData(t *testing.T) {
	a := &Snapshot{Width: 100, Height: 100, Data: []byte{}}
	b := &Snapshot{Width: 100, Height: 100, Data: []byte{}}

	if !Compare(a, b) {
		t.Error("Compare() should return true for both empty data")
	}
}

func TestCompareNilData(t *testing.T) {
	a := &Snapshot{Width: 100, Height: 100, Data: nil}
	b := &Snapshot{Width: 100, Height: 100, Data: nil}

	if !Compare(a, b) {
		t.Error("Compare() should return true for both nil data")
	}
}

func TestCompareNilVsEmpty(t *testing.T) {
	a := &Snapshot{Width: 100, Height: 100, Data: nil}
	b := &Snapshot{Width: 100, Height: 100, Data: []byte{}}

	// Both have length 0, so they should be equal
	if !Compare(a, b) {
		t.Error("Compare() should return true for nil vs empty (both len 0)")
	}
}

func TestSnapshotStruct(t *testing.T) {
	snap := &Snapshot{
		GenreID: "horror",
		Width:   1024,
		Height:  768,
		Data:    []byte{255, 128, 64, 32},
	}

	if snap.GenreID != "horror" {
		t.Errorf("GenreID = %q, want %q", snap.GenreID, "horror")
	}
	if snap.Width != 1024 {
		t.Errorf("Width = %d, want 1024", snap.Width)
	}
	if snap.Height != 768 {
		t.Errorf("Height = %d, want 768", snap.Height)
	}
	if len(snap.Data) != 4 {
		t.Errorf("len(Data) = %d, want 4", len(snap.Data))
	}
}

func TestCompareLargeData(t *testing.T) {
	// Test with larger data to ensure loop works correctly
	data1 := make([]byte, 1000)
	data2 := make([]byte, 1000)
	for i := range data1 {
		data1[i] = byte(i % 256)
		data2[i] = byte(i % 256)
	}

	a := &Snapshot{Width: 100, Height: 100, Data: data1}
	b := &Snapshot{Width: 100, Height: 100, Data: data2}

	if !Compare(a, b) {
		t.Error("Compare() should return true for identical large data")
	}

	// Change one byte
	data2[500] = 99
	if Compare(a, b) {
		t.Error("Compare() should return false when one byte differs")
	}
}
