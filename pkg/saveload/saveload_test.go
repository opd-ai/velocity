package saveload

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoad_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test_save.json")

	original := &RunState{
		Version:    1,
		Seed:       12345,
		Genre:      "scifi",
		Wave:       5,
		Score:      10000,
		PlayerData: []byte{1, 2, 3, 4},
	}

	// Save
	if err := Save(path, original); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("save file not created")
	}

	// Load
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Verify fields
	if loaded.Version != original.Version {
		t.Errorf("version mismatch: got %d, want %d", loaded.Version, original.Version)
	}
	if loaded.Seed != original.Seed {
		t.Errorf("seed mismatch: got %d, want %d", loaded.Seed, original.Seed)
	}
	if loaded.Genre != original.Genre {
		t.Errorf("genre mismatch: got %s, want %s", loaded.Genre, original.Genre)
	}
	if loaded.Wave != original.Wave {
		t.Errorf("wave mismatch: got %d, want %d", loaded.Wave, original.Wave)
	}
	if loaded.Score != original.Score {
		t.Errorf("score mismatch: got %d, want %d", loaded.Score, original.Score)
	}
	if len(loaded.PlayerData) != len(original.PlayerData) {
		t.Errorf("player_data length mismatch: got %d, want %d", len(loaded.PlayerData), len(original.PlayerData))
	}
}

func TestSave_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "new_save.json")

	state := &RunState{
		Version: 1,
		Seed:    42,
		Genre:   "fantasy",
		Wave:    1,
		Score:   0,
	}

	if err := Save(path, state); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Check file was created
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Fatal("save file not created")
	}
	if info.Size() == 0 {
		t.Error("save file is empty")
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/save.json")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "invalid.json")

	// Write invalid JSON
	if err := os.WriteFile(path, []byte("not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestRunState_Fields(t *testing.T) {
	rs := RunState{
		Version:    2,
		Seed:       99999,
		Genre:      "horror",
		Wave:       10,
		Score:      50000,
		PlayerData: []byte("test data"),
	}

	if rs.Version != 2 {
		t.Error("version mismatch")
	}
	if rs.Seed != 99999 {
		t.Error("seed mismatch")
	}
	if rs.Genre != "horror" {
		t.Error("genre mismatch")
	}
	if rs.Wave != 10 {
		t.Error("wave mismatch")
	}
	if rs.Score != 50000 {
		t.Error("score mismatch")
	}
	if string(rs.PlayerData) != "test data" {
		t.Error("player_data mismatch")
	}
}

func TestRunState_EmptyPlayerData(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "no_player_data.json")

	original := &RunState{
		Version: 1,
		Seed:    1,
		Genre:   "scifi",
		Wave:    1,
		Score:   100,
		// PlayerData intentionally nil
	}

	if err := Save(path, original); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(loaded.PlayerData) != 0 {
		t.Errorf("expected empty player_data, got %d bytes", len(loaded.PlayerData))
	}
}

func TestSave_Overwrite(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "overwrite.json")

	// First save
	state1 := &RunState{Version: 1, Score: 100}
	if err := Save(path, state1); err != nil {
		t.Fatal(err)
	}

	// Overwrite with new data
	state2 := &RunState{Version: 1, Score: 200}
	if err := Save(path, state2); err != nil {
		t.Fatal(err)
	}

	// Load and verify overwrite
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Score != 200 {
		t.Errorf("expected score 200 after overwrite, got %d", loaded.Score)
	}
}

func TestSave_InvalidPath(t *testing.T) {
	state := &RunState{Version: 1}
	err := Save("/nonexistent/directory/save.json", state)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestRunState_AllGenres(t *testing.T) {
	genres := []string{"scifi", "fantasy", "horror", "cyberpunk", "postapoc"}
	tmpDir := t.TempDir()

	for _, genre := range genres {
		t.Run(genre, func(t *testing.T) {
			path := filepath.Join(tmpDir, genre+"_save.json")
			state := &RunState{Version: 1, Genre: genre}

			if err := Save(path, state); err != nil {
				t.Fatalf("Save() error: %v", err)
			}

			loaded, err := Load(path)
			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}

			if loaded.Genre != genre {
				t.Errorf("genre mismatch: got %s, want %s", loaded.Genre, genre)
			}
		})
	}
}

func BenchmarkSave(b *testing.B) {
	tmpDir := b.TempDir()
	path := filepath.Join(tmpDir, "bench_save.json")
	state := &RunState{
		Version:    1,
		Seed:       12345,
		Genre:      "scifi",
		Wave:       100,
		Score:      999999,
		PlayerData: make([]byte, 1024),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Save(path, state)
	}
}

func BenchmarkLoad(b *testing.B) {
	tmpDir := b.TempDir()
	path := filepath.Join(tmpDir, "bench_load.json")
	state := &RunState{
		Version:    1,
		Seed:       12345,
		Genre:      "scifi",
		Wave:       100,
		Score:      999999,
		PlayerData: make([]byte, 1024),
	}
	Save(path, state)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Load(path)
	}
}
