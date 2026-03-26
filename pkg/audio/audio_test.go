package audio

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("expected non-nil Manager")
	}
	if m.genreID != "scifi" {
		t.Errorf("expected default genre 'scifi', got %s", m.genreID)
	}
}

func TestManager_SetGenre(t *testing.T) {
	m := NewManager()
	m.SetGenre("fantasy")
	if m.genreID != "fantasy" {
		t.Errorf("expected genre 'fantasy', got %s", m.genreID)
	}
}

func TestManager_SetVolumes(t *testing.T) {
	m := NewManager()
	m.SetVolumes(0.5, 0.4, 0.6)

	if m.masterVolume != 0.5 {
		t.Errorf("expected master 0.5, got %f", m.masterVolume)
	}
	if m.musicVolume != 0.4 {
		t.Errorf("expected music 0.4, got %f", m.musicVolume)
	}
	if m.sfxVolume != 0.6 {
		t.Errorf("expected sfx 0.6, got %f", m.sfxVolume)
	}
}

func TestManager_SetVolumes_Clamping(t *testing.T) {
	m := NewManager()
	m.SetVolumes(-0.5, 1.5, 0.5)

	if m.masterVolume != 0 {
		t.Errorf("expected master clamped to 0, got %f", m.masterVolume)
	}
	if m.musicVolume != 1.0 {
		t.Errorf("expected music clamped to 1, got %f", m.musicVolume)
	}
}

func TestManager_SetIntensity(t *testing.T) {
	m := NewManager()
	m.SetIntensity(0.7)
	if m.intensity != 0.7 {
		t.Errorf("expected intensity 0.7, got %f", m.intensity)
	}
}

func TestManager_PlayMusic(t *testing.T) {
	m := NewManager()
	m.PlayMusic()
	if !m.musicPlaying {
		t.Error("expected musicPlaying to be true")
	}
}

func TestManager_StopMusic(t *testing.T) {
	m := NewManager()
	m.PlayMusic()
	m.StopMusic()
	if m.musicPlaying {
		t.Error("expected musicPlaying to be false")
	}
}

func TestManager_PlaySFX(t *testing.T) {
	m := NewManager()
	// Should not panic
	m.PlaySFX("laser")
	m.PlaySFX("explosion")
}

func TestManager_PlaySFXAt(t *testing.T) {
	m := NewManager()
	m.PlaySFXAt("explosion", 100, 200, true)

	m.sfxMu.Lock()
	count := len(m.sfxQueue)
	m.sfxMu.Unlock()

	if count != 1 {
		t.Errorf("expected 1 queued SFX, got %d", count)
	}
}

func TestManager_Update(t *testing.T) {
	m := NewManager()
	m.PlaySFXAt("laser", 0, 0, false)
	m.Update()

	m.sfxMu.Lock()
	count := len(m.sfxQueue)
	m.sfxMu.Unlock()

	if count != 0 {
		t.Errorf("expected queue cleared, got %d items", count)
	}
}

func TestGenerateTone(t *testing.T) {
	data := GenerateTone(440.0, 0.1)
	if len(data) == 0 {
		t.Fatal("expected non-empty audio data")
	}

	// 44100 samples/sec * 0.1 sec * 4 bytes/sample
	expectedLen := int(0.1 * 44100 * 4)
	if len(data) != expectedLen {
		t.Errorf("expected %d bytes, got %d", expectedLen, len(data))
	}
}

func TestGenerateLaserSFX(t *testing.T) {
	data := GenerateLaserSFX()
	if len(data) == 0 {
		t.Fatal("expected non-empty laser SFX data")
	}
}

func TestGenerateExplosionSFX(t *testing.T) {
	data := GenerateExplosionSFX()
	if len(data) == 0 {
		t.Fatal("expected non-empty explosion SFX data")
	}
}

func TestGeneratePowerupSFX(t *testing.T) {
	data := GeneratePowerupSFX()
	if len(data) == 0 {
		t.Fatal("expected non-empty powerup SFX data")
	}
}

func TestGenerateMenuSelectSFX(t *testing.T) {
	data := GenerateMenuSelectSFX()
	if len(data) == 0 {
		t.Fatal("expected non-empty menu select SFX data")
	}
}

func TestGetSFXData(t *testing.T) {
	tests := []struct {
		name     string
		sfxName  string
		wantData bool
	}{
		{"laser", "laser", true},
		{"explosion", "explosion", true},
		{"powerup", "powerup", true},
		{"menu_select", "menu_select", true},
		{"wave_start", "wave_start", true},
		{"wave_complete", "wave_complete", true},
		{"unknown", "unknown_sfx", true}, // Returns default tone
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := GetSFXData(tt.sfxName)
			if tt.wantData && len(data) == 0 {
				t.Errorf("expected SFX data for %s", tt.sfxName)
			}
		})
	}
}

func TestCalculateSpatialVolume(t *testing.T) {
	tests := []struct {
		name         string
		listenerX    float64
		listenerY    float64
		sourceX      float64
		sourceY      float64
		maxDist      float64
		wantVolume   float64
		wantPan      float64
		volumeTol    float64
		panTol       float64
	}{
		{"same position", 0, 0, 0, 0, 100, 1.0, 0, 0.01, 0.01},
		{"max distance", 0, 0, 100, 0, 100, 0, 0, 0.01, 0.01},
		{"beyond max", 0, 0, 150, 0, 100, 0, 0, 0.01, 0.01},
		{"half distance", 0, 0, 50, 0, 100, 0.5, 0.5, 0.01, 0.01},
		{"left side", 0, 0, -50, 0, 100, 0.5, -0.5, 0.01, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vol, pan := CalculateSpatialVolume(tt.listenerX, tt.listenerY, tt.sourceX, tt.sourceY, tt.maxDist)

			if vol < tt.wantVolume-tt.volumeTol || vol > tt.wantVolume+tt.volumeTol {
				t.Errorf("volume = %f, want %f", vol, tt.wantVolume)
			}
			if pan < tt.wantPan-tt.panTol || pan > tt.wantPan+tt.panTol {
				t.Errorf("pan = %f, want %f", pan, tt.wantPan)
			}
		})
	}
}

func TestApplySpatialAudio(t *testing.T) {
	data := GenerateTone(440.0, 0.05)
	original := make([]byte, len(data))
	copy(original, data)

	// Apply spatial audio
	result := ApplySpatialAudio(data, 0.5, 0.0)

	if len(result) != len(data) {
		t.Errorf("expected same length, got %d vs %d", len(result), len(data))
	}

	// Verify volume was reduced
	// (simplified check - just verify data was modified)
	modified := false
	for i := 0; i < len(result) && i < len(original); i++ {
		if result[i] != original[i] {
			modified = true
			break
		}
	}
	if !modified {
		t.Error("expected spatial audio to modify data")
	}
}

func TestGetGenreParams(t *testing.T) {
	genres := []string{"scifi", "fantasy", "horror", "cyberpunk", "postapoc"}

	for _, genre := range genres {
		t.Run(genre, func(t *testing.T) {
			params := GetGenreParams(genre)
			if params.BaseFrequency <= 0 {
				t.Error("expected positive base frequency")
			}
			if len(params.Scale) == 0 {
				t.Error("expected non-empty scale")
			}
			if params.Tempo <= 0 {
				t.Error("expected positive tempo")
			}
		})
	}

	// Test unknown genre
	params := GetGenreParams("unknown")
	if params.BaseFrequency <= 0 {
		t.Error("expected default params for unknown genre")
	}
}

func TestClampVolume(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.5, 0.5},
		{-0.5, 0},
		{1.5, 1},
		{0, 0},
		{1, 1},
	}

	for _, tt := range tests {
		result := clampVolume(tt.input)
		if result != tt.expected {
			t.Errorf("clampVolume(%f) = %f, want %f", tt.input, result, tt.expected)
		}
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		v, min, max, expected float64
	}{
		{0.5, 0, 1, 0.5},
		{-0.5, 0, 1, 0},
		{1.5, 0, 1, 1},
		{5, 0, 10, 5},
	}

	for _, tt := range tests {
		result := clamp(tt.v, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("clamp(%f, %f, %f) = %f, want %f", tt.v, tt.min, tt.max, result, tt.expected)
		}
	}
}

func TestApplyEnvelope(t *testing.T) {
	total := 10000

	// Start should be low (attack)
	startEnv := applyEnvelope(0, total)
	if startEnv > 0.1 {
		t.Errorf("expected low envelope at start, got %f", startEnv)
	}

	// Middle should be 1.0
	midEnv := applyEnvelope(total/2, total)
	if midEnv != 1.0 {
		t.Errorf("expected envelope 1.0 in middle, got %f", midEnv)
	}

	// End should be low (release)
	endEnv := applyEnvelope(total-1, total)
	if endEnv > 0.1 {
		t.Errorf("expected low envelope at end, got %f", endEnv)
	}
}

func TestSquareWave(t *testing.T) {
	// Test at various phases
	if squareWave(0) != 1.0 {
		t.Error("expected 1.0 at phase 0")
	}
	// sin(3.5) is negative
	if squareWave(3.5) != -1.0 {
		t.Error("expected -1.0 at phase 3.5")
	}
}

func TestPseudoRandom(t *testing.T) {
	// Test determinism
	val1 := pseudoRandom(12345)
	val2 := pseudoRandom(12345)
	if val1 != val2 {
		t.Error("expected deterministic results")
	}

	// Test range (0-1)
	for i := uint32(0); i < 100; i++ {
		val := pseudoRandom(i)
		if val < 0 || val > 1 {
			t.Errorf("pseudoRandom(%d) = %f, expected [0,1]", i, val)
		}
	}
}

func BenchmarkGenerateTone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateTone(440.0, 0.1)
	}
}

func BenchmarkGenerateLaserSFX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateLaserSFX()
	}
}

func BenchmarkGenerateExplosionSFX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateExplosionSFX()
	}
}
