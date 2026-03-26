//go:build noebiten

package audio

// stubBackend implements AudioBackend as a no-op for headless/test environments.
type stubBackend struct{}

// newAudioBackend creates a stub audio backend for headless environments.
func newAudioBackend() AudioBackend {
	return &stubBackend{}
}

// Initialize is a no-op for the stub backend.
func (b *stubBackend) Initialize() {}

// PlayBytes is a no-op for the stub backend.
func (b *stubBackend) PlayBytes(data []byte) {}
