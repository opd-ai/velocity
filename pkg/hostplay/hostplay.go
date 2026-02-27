// Package hostplay provides the local-host authoritative server scaffold.
package hostplay

// Host represents a local authoritative game server.
type Host struct {
	port    int
	running bool
	maxPlayers int
}

// NewHost creates a new host-play server.
func NewHost(port, maxPlayers int) *Host {
	return &Host{port: port, maxPlayers: maxPlayers}
}

// Start begins the authoritative game server.
func (h *Host) Start() error {
	// Stub: will start local authoritative server loop.
	h.running = true
	return nil
}

// Stop shuts down the host server.
func (h *Host) Stop() {
	h.running = false
}

// IsRunning returns whether the host is currently active.
func (h *Host) IsRunning() bool {
	return h.running
}
