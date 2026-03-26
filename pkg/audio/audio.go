// Package audio provides adaptive music, SFX playback, and positional audio.
package audio

import (
	"encoding/binary"
	"math"
	"sync"
)

// Audio system constants
const (
	// SampleRate is the audio sample rate in Hz.
	SampleRate = 44100
	// BytesPerSample is the number of bytes per sample (16-bit stereo).
	BytesPerSample = 4
)

// Pentatonic scale note frequencies in Hz.
// C4, D4, E4, G4, A4 - a pleasant, non-dissonant scale for procedural music.
const (
	NoteC4 = 261.63
	NoteD4 = 293.66
	NoteE4 = 329.63
	NoteG4 = 392.00
	NoteA4 = 440.00
)

// PentatonicScale returns the standard PentatonicScale() scale frequencies.
func PentatonicScale() []float64 {
	return []float64{NoteC4, NoteD4, NoteE4, NoteG4, NoteA4}
}

// Manager handles all audio playback.
type Manager struct {
	genreID      string
	masterVolume float64
	musicVolume  float64
	sfxVolume    float64

	intensity    float64 // Music intensity level (0.0-1.0)
	musicPlaying bool

	sfxQueue []SFXRequest
	sfxMu    sync.Mutex

	playerX float64
	playerY float64

	// audioBackend is the platform-specific audio backend (Ebiten or stub)
	audioBackend AudioBackend
}

// AudioBackend defines the interface for platform-specific audio playback.
type AudioBackend interface {
	// PlayBytes plays raw PCM audio data.
	PlayBytes(data []byte)
	// Initialize initializes the audio backend.
	Initialize()
}

// SFXRequest represents a queued sound effect.
type SFXRequest struct {
	Name    string
	X, Y    float64 // Position for spatial audio
	Spatial bool    // Whether to apply positional audio
}

// NewManager creates a new audio manager.
func NewManager() *Manager {
	m := &Manager{
		genreID:      "scifi",
		masterVolume: 0.8,
		musicVolume:  0.6,
		sfxVolume:    0.8,
		sfxQueue:     make([]SFXRequest, 0, 16),
		audioBackend: newAudioBackend(),
	}
	return m
}

// SetGenre switches audio assets to match the given genre.
func (m *Manager) SetGenre(genreID string) {
	m.genreID = genreID
}

// SetVolumes sets the master, music, and SFX volume levels.
func (m *Manager) SetVolumes(master, music, sfx float64) {
	m.masterVolume = clampVolume(master)
	m.musicVolume = clampVolume(music)
	m.sfxVolume = clampVolume(sfx)
}

// SetIntensity sets the music intensity level (0.0-1.0).
func (m *Manager) SetIntensity(intensity float64) {
	m.intensity = clampVolume(intensity)
}

// SetPlayerPosition updates player position for spatial audio.
func (m *Manager) SetPlayerPosition(x, y float64) {
	m.playerX = x
	m.playerY = y
}

// PlaySFX plays a named sound effect.
func (m *Manager) PlaySFX(name string) {
	m.PlaySFXAt(name, m.playerX, m.playerY, false)
}

// PlaySFXAt plays a sound effect at a specific position.
func (m *Manager) PlaySFXAt(name string, x, y float64, spatial bool) {
	m.sfxMu.Lock()
	defer m.sfxMu.Unlock()

	// Limit queue size to prevent memory issues
	if len(m.sfxQueue) < 32 {
		m.sfxQueue = append(m.sfxQueue, SFXRequest{
			Name:    name,
			X:       x,
			Y:       y,
			Spatial: spatial,
		})
	}
}

// PlayMusic starts adaptive background music.
func (m *Manager) PlayMusic() {
	m.musicPlaying = true
}

// StopMusic stops background music.
func (m *Manager) StopMusic() {
	m.musicPlaying = false
}

// Update advances the audio state each frame.
func (m *Manager) Update() {
	// Lazy-initialize audio backend on first update
	if m.audioBackend != nil {
		m.audioBackend.Initialize()
	}

	// Process any queued SFX
	m.sfxMu.Lock()
	queue := make([]SFXRequest, len(m.sfxQueue))
	copy(queue, m.sfxQueue)
	m.sfxQueue = m.sfxQueue[:0] // Clear processed queue
	m.sfxMu.Unlock()

	// Play each queued SFX
	for _, req := range queue {
		m.playSFXNow(req)
	}
}

// playSFXNow plays a sound effect immediately using the audio backend.
func (m *Manager) playSFXNow(req SFXRequest) {
	if m.audioBackend == nil {
		return
	}

	// Get the raw PCM data for this SFX
	data := GetSFXData(req.Name)
	if len(data) == 0 {
		return
	}

	// Apply spatial audio if requested
	if req.Spatial {
		vol, pan := CalculateSpatialVolume(m.playerX, m.playerY, req.X, req.Y, 500.0)
		data = ApplySpatialAudio(data, vol*m.sfxVolume*m.masterVolume, pan)
	} else {
		// Apply volume scaling for non-spatial audio
		data = ApplySpatialAudio(data, m.sfxVolume*m.masterVolume, 0)
	}

	// Play the audio data
	m.audioBackend.PlayBytes(data)
}

// GenerateTone creates PCM audio data for a simple tone.
func GenerateTone(frequency, duration float64) []byte {
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		sample := math.Sin(2 * math.Pi * frequency * t)

		// Apply envelope to avoid clicks
		envelope := applyEnvelope(i, numSamples)
		sample *= envelope

		// Convert to 16-bit signed int
		intSample := int16(sample * 32767 * 0.5)                      // 50% volume
		binary.LittleEndian.PutUint16(buf[i*4:], uint16(intSample))   // Left
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(intSample)) // Right
	}

	return buf
}

// GenerateLaserSFX creates PCM data for a laser sound.
func GenerateLaserSFX() []byte {
	duration := 0.1
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	baseFreq := 880.0 // A5

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		progress := float64(i) / float64(numSamples)

		// Frequency sweep down
		freq := baseFreq * (1.0 - progress*0.5)

		// Square wave for laser sound
		sample := squareWave(2 * math.Pi * freq * t)

		// Apply envelope
		envelope := applyEnvelope(i, numSamples)
		sample *= envelope * 0.3

		intSample := int16(sample * 32767)
		binary.LittleEndian.PutUint16(buf[i*4:], uint16(intSample))
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(intSample))
	}

	return buf
}

// GenerateExplosionSFX creates PCM data for an explosion sound.
func GenerateExplosionSFX() []byte {
	duration := 0.3
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	for i := 0; i < numSamples; i++ {
		progress := float64(i) / float64(numSamples)

		// White noise burst with exponential decay
		noise := (pseudoRandom(uint32(i)) - 0.5) * 2.0
		envelope := math.Exp(-progress * 5.0)

		// Add low frequency rumble
		t := float64(i) / SampleRate
		rumble := math.Sin(2*math.Pi*60*t) * 0.3

		sample := (noise*0.7 + rumble) * envelope * 0.4

		intSample := int16(clamp(sample, -1, 1) * 32767)
		binary.LittleEndian.PutUint16(buf[i*4:], uint16(intSample))
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(intSample))
	}

	return buf
}

// GeneratePowerupSFX creates PCM data for a powerup collect sound.
func GeneratePowerupSFX() []byte {
	duration := 0.2
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	baseFreq := 440.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		progress := float64(i) / float64(numSamples)

		// Rising frequency sweep
		freq := baseFreq * (1.0 + progress*1.5)

		// Sine wave for pleasing sound
		sample := math.Sin(2 * math.Pi * freq * t)

		envelope := applyEnvelope(i, numSamples)
		sample *= envelope * 0.4

		intSample := int16(sample * 32767)
		binary.LittleEndian.PutUint16(buf[i*4:], uint16(intSample))
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(intSample))
	}

	return buf
}

// GenerateMenuSelectSFX creates PCM data for a menu selection sound.
func GenerateMenuSelectSFX() []byte {
	duration := 0.05
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	freq := 660.0 // E5

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		sample := math.Sin(2 * math.Pi * freq * t)

		envelope := applyEnvelope(i, numSamples)
		sample *= envelope * 0.3

		intSample := int16(sample * 32767)
		binary.LittleEndian.PutUint16(buf[i*4:], uint16(intSample))
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(intSample))
	}

	return buf
}

// GetSFXData returns PCM data for a named sound effect.
func GetSFXData(name string) []byte {
	switch name {
	case "laser", "fire", "shoot":
		return GenerateLaserSFX()
	case "explosion", "death", "enemy_death":
		return GenerateExplosionSFX()
	case "powerup", "pickup", "collect":
		return GeneratePowerupSFX()
	case "menu_select", "select", "ui_click":
		return GenerateMenuSelectSFX()
	case "wave_start":
		return GenerateTone(440.0, 0.2)
	case "wave_complete":
		return GeneratePowerupSFX()
	default:
		return GenerateTone(440.0, 0.1)
	}
}

// CalculateSpatialVolume returns volume and pan based on distance and angle.
func CalculateSpatialVolume(listenerX, listenerY, sourceX, sourceY, maxDistance float64) (volume, pan float64) {
	dx := sourceX - listenerX
	dy := sourceY - listenerY
	distance := math.Sqrt(dx*dx + dy*dy)

	// Volume attenuation based on distance
	if distance >= maxDistance {
		return 0, 0
	}
	volume = 1.0 - (distance / maxDistance)
	volume = math.Max(0, math.Min(1, volume))

	// Pan based on horizontal position (-1 = left, +1 = right)
	if distance > 0 {
		pan = dx / maxDistance
		pan = math.Max(-1, math.Min(1, pan))
	}

	return volume, pan
}

// ApplySpatialAudio modifies PCM data with volume and pan.
func ApplySpatialAudio(data []byte, volume, pan float64) []byte {
	result := make([]byte, len(data))
	copy(result, data)

	leftVol := volume * (1.0 - math.Max(0, pan))
	rightVol := volume * (1.0 + math.Min(0, pan))

	for i := 0; i < len(result); i += 4 {
		// Left channel
		left := int16(binary.LittleEndian.Uint16(result[i:]))
		left = int16(float64(left) * leftVol)
		binary.LittleEndian.PutUint16(result[i:], uint16(left))

		// Right channel
		right := int16(binary.LittleEndian.Uint16(result[i+2:]))
		right = int16(float64(right) * rightVol)
		binary.LittleEndian.PutUint16(result[i+2:], uint16(right))
	}

	return result
}

// applyEnvelope applies attack/release envelope to avoid clicks.
func applyEnvelope(sampleIndex, totalSamples int) float64 {
	attackSamples := int(0.01 * SampleRate)
	releaseSamples := int(0.01 * SampleRate)

	if sampleIndex < attackSamples {
		return float64(sampleIndex) / float64(attackSamples)
	}
	if sampleIndex > totalSamples-releaseSamples {
		return float64(totalSamples-sampleIndex) / float64(releaseSamples)
	}
	return 1.0
}

// squareWave generates a square wave value.
func squareWave(phase float64) float64 {
	if math.Sin(phase) >= 0 {
		return 1.0
	}
	return -1.0
}

// pseudoRandom generates a pseudo-random value for noise (0-1).
func pseudoRandom(seed uint32) float64 {
	seed = seed*1103515245 + 12345
	return float64(seed&0x7FFFFFFF) / float64(0x7FFFFFFF)
}

// clampVolume ensures volume is between 0 and 1.
func clampVolume(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// clamp ensures value is within bounds.
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// GenreAudioParams holds genre-specific audio parameters.
type GenreAudioParams struct {
	BaseFrequency float64
	Scale         []float64
	Tempo         float64
	WaveformMix   float64 // 0 = sine, 1 = square
}

// GetGenreParams returns audio parameters for a genre.
func GetGenreParams(genreID string) GenreAudioParams {
	switch genreID {
	case "scifi":
		return GenreAudioParams{
			BaseFrequency: 220.0,
			Scale:         PentatonicScale(),
			Tempo:         120.0,
			WaveformMix:   0.3,
		}
	case "fantasy":
		return GenreAudioParams{
			BaseFrequency: 261.63,
			Scale:         []float64{261.63, 293.66, 329.63, 349.23, 392.00, 440.00},
			Tempo:         100.0,
			WaveformMix:   0.1,
		}
	case "horror":
		return GenreAudioParams{
			BaseFrequency: 110.0,
			Scale:         []float64{110.0, 116.54, 130.81, 146.83, 155.56},
			Tempo:         80.0,
			WaveformMix:   0.5,
		}
	case "cyberpunk":
		return GenreAudioParams{
			BaseFrequency: 330.0,
			Scale:         PentatonicScale(),
			Tempo:         140.0,
			WaveformMix:   0.7,
		}
	case "postapoc":
		return GenreAudioParams{
			BaseFrequency: 196.0,
			Scale:         []float64{196.0, 220.0, 246.94, 293.66, 329.63},
			Tempo:         90.0,
			WaveformMix:   0.4,
		}
	default:
		return GenreAudioParams{
			BaseFrequency: 220.0,
			Scale:         PentatonicScale(),
			Tempo:         120.0,
			WaveformMix:   0.3,
		}
	}
}
