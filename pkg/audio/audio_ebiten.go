//go:build !noebiten

package audio

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// ebitenBackend implements AudioBackend using Ebitengine's audio system.
type ebitenBackend struct {
	context     *audio.Context
	initialized bool
}

// newAudioBackend creates the Ebiten-based audio backend.
func newAudioBackend() AudioBackend {
	return &ebitenBackend{}
}

// Initialize initializes the Ebiten audio context.
func (b *ebitenBackend) Initialize() {
	if b.initialized {
		return
	}
	b.context = audio.NewContext(SampleRate)
	b.initialized = true
}

// PlayBytes plays raw PCM audio data using Ebiten's audio system.
func (b *ebitenBackend) PlayBytes(data []byte) {
	if b.context == nil {
		return
	}
	player := b.context.NewPlayerFromBytes(data)
	player.Play()
}
