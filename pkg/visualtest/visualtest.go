// Package visualtest provides automated per-genre screenshot regression testing.
package visualtest

// Snapshot represents a captured frame for visual comparison.
type Snapshot struct {
	GenreID string
	Width   int
	Height  int
	Data    []byte
}

// Capture takes a snapshot of the current frame.
func Capture(genreID string, width, height int) *Snapshot {
	// Stub: will capture the frame buffer as raw pixel data.
	return &Snapshot{
		GenreID: genreID,
		Width:   width,
		Height:  height,
	}
}

// Compare returns true if two snapshots are identical.
func Compare(a, b *Snapshot) bool {
	if a.Width != b.Width || a.Height != b.Height {
		return false
	}
	if len(a.Data) != len(b.Data) {
		return false
	}
	for i := range a.Data {
		if a.Data[i] != b.Data[i] {
			return false
		}
	}
	return true
}
