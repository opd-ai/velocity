package networking

import "testing"

func TestNewServer(t *testing.T) {
	s := NewServer(7777)
	if s == nil {
		t.Fatal("NewServer() returned nil")
	}
	if s.port != 7777 {
		t.Errorf("port = %d, want 7777", s.port)
	}
	if s.running {
		t.Error("new server should not be running")
	}
}

func TestServerStart(t *testing.T) {
	s := NewServer(8080)

	err := s.Start()
	if err != nil {
		t.Errorf("Start() returned error: %v", err)
	}

	if !s.running {
		t.Error("server should be running after Start()")
	}
}

func TestServerStop(t *testing.T) {
	s := NewServer(8080)
	s.Start()

	if !s.running {
		t.Fatal("server should be running after Start()")
	}

	s.Stop()

	if s.running {
		t.Error("server should not be running after Stop()")
	}
}

func TestServerStartStopCycle(t *testing.T) {
	s := NewServer(5000)

	for i := 0; i < 3; i++ {
		err := s.Start()
		if err != nil {
			t.Errorf("iteration %d: Start() error: %v", i, err)
		}
		if !s.running {
			t.Errorf("iteration %d: server should be running", i)
		}

		s.Stop()
		if s.running {
			t.Errorf("iteration %d: server should not be running", i)
		}
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("localhost:7777")
	if c == nil {
		t.Fatal("NewClient() returned nil")
	}
	if c.address != "localhost:7777" {
		t.Errorf("address = %q, want %q", c.address, "localhost:7777")
	}
	if c.connected {
		t.Error("new client should not be connected")
	}
}

func TestClientConnect(t *testing.T) {
	c := NewClient("127.0.0.1:8080")

	err := c.Connect()
	if err != nil {
		t.Errorf("Connect() returned error: %v", err)
	}

	if !c.connected {
		t.Error("client should be connected after Connect()")
	}
}

func TestClientDisconnect(t *testing.T) {
	c := NewClient("127.0.0.1:8080")
	c.Connect()

	if !c.connected {
		t.Fatal("client should be connected after Connect()")
	}

	c.Disconnect()

	if c.connected {
		t.Error("client should not be connected after Disconnect()")
	}
}

func TestClientConnectDisconnectCycle(t *testing.T) {
	c := NewClient("game.example.com:27015")

	for i := 0; i < 3; i++ {
		err := c.Connect()
		if err != nil {
			t.Errorf("iteration %d: Connect() error: %v", i, err)
		}
		if !c.connected {
			t.Errorf("iteration %d: client should be connected", i)
		}

		c.Disconnect()
		if c.connected {
			t.Errorf("iteration %d: client should not be connected", i)
		}
	}
}

func TestServerStruct(t *testing.T) {
	s := &Server{
		port:    9999,
		running: true,
	}

	if s.port != 9999 {
		t.Errorf("port = %d, want 9999", s.port)
	}
	if !s.running {
		t.Error("running should be true")
	}
}

func TestClientStruct(t *testing.T) {
	c := &Client{
		address:   "example.com:1234",
		connected: true,
	}

	if c.address != "example.com:1234" {
		t.Errorf("address = %q, want %q", c.address, "example.com:1234")
	}
	if !c.connected {
		t.Error("connected should be true")
	}
}

func TestNewServerDifferentPorts(t *testing.T) {
	ports := []int{80, 443, 7777, 27015, 65535}

	for _, port := range ports {
		s := NewServer(port)
		if s.port != port {
			t.Errorf("port = %d, want %d", s.port, port)
		}
	}
}

func TestNewClientDifferentAddresses(t *testing.T) {
	addresses := []string{
		"localhost:8080",
		"127.0.0.1:7777",
		"192.168.1.1:27015",
		"game.example.com:443",
		"[::1]:8080",
	}

	for _, addr := range addresses {
		c := NewClient(addr)
		if c.address != addr {
			t.Errorf("address = %q, want %q", c.address, addr)
		}
	}
}
