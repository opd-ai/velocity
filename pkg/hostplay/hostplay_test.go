package hostplay

import "testing"

func TestNewHost(t *testing.T) {
	h := NewHost(7777, 4)
	if h == nil {
		t.Fatal("NewHost() returned nil")
	}
	if h.port != 7777 {
		t.Errorf("port = %d, want 7777", h.port)
	}
	if h.maxPlayers != 4 {
		t.Errorf("maxPlayers = %d, want 4", h.maxPlayers)
	}
	if h.running {
		t.Error("new host should not be running")
	}
}

func TestHostStart(t *testing.T) {
	h := NewHost(8080, 8)

	if h.IsRunning() {
		t.Error("host should not be running before Start()")
	}

	err := h.Start()
	if err != nil {
		t.Errorf("Start() returned error: %v", err)
	}

	if !h.IsRunning() {
		t.Error("host should be running after Start()")
	}
}

func TestHostStop(t *testing.T) {
	h := NewHost(8080, 8)
	h.Start()

	if !h.IsRunning() {
		t.Fatal("host should be running after Start()")
	}

	h.Stop()

	if h.IsRunning() {
		t.Error("host should not be running after Stop()")
	}
}

func TestHostIsRunning(t *testing.T) {
	h := NewHost(9999, 2)

	if h.IsRunning() {
		t.Error("new host should return false for IsRunning()")
	}

	h.Start()
	if !h.IsRunning() {
		t.Error("started host should return true for IsRunning()")
	}

	h.Stop()
	if h.IsRunning() {
		t.Error("stopped host should return false for IsRunning()")
	}
}

func TestHostStartStopCycle(t *testing.T) {
	h := NewHost(5000, 4)

	for i := 0; i < 3; i++ {
		err := h.Start()
		if err != nil {
			t.Errorf("iteration %d: Start() error: %v", i, err)
		}
		if !h.IsRunning() {
			t.Errorf("iteration %d: host should be running after Start()", i)
		}

		h.Stop()
		if h.IsRunning() {
			t.Errorf("iteration %d: host should not be running after Stop()", i)
		}
	}
}

func TestHostDifferentPorts(t *testing.T) {
	ports := []int{1234, 7777, 8080, 27015}

	for _, port := range ports {
		h := NewHost(port, 4)
		if h.port != port {
			t.Errorf("port = %d, want %d", h.port, port)
		}
	}
}

func TestHostDifferentMaxPlayers(t *testing.T) {
	maxPlayers := []int{1, 2, 4, 8, 16, 32}

	for _, max := range maxPlayers {
		h := NewHost(8080, max)
		if h.maxPlayers != max {
			t.Errorf("maxPlayers = %d, want %d", h.maxPlayers, max)
		}
	}
}

func TestHostStruct(t *testing.T) {
	h := &Host{
		port:       6666,
		running:    true,
		maxPlayers: 10,
	}

	if h.port != 6666 {
		t.Errorf("port = %d, want 6666", h.port)
	}
	if !h.running {
		t.Error("running should be true")
	}
	if h.maxPlayers != 10 {
		t.Errorf("maxPlayers = %d, want 10", h.maxPlayers)
	}
}
