package stability

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func TestNewWatchdog(t *testing.T) {
	w := NewWatchdog(500 * time.Millisecond)
	if w == nil {
		t.Fatal("NewWatchdog() returned nil")
	}
	if w.timeout != 500*time.Millisecond {
		t.Errorf("timeout = %v, want 500ms", w.timeout)
	}
	if w.active {
		t.Error("new watchdog should not be active")
	}
}

func TestWatchdogStart(t *testing.T) {
	w := NewWatchdog(100 * time.Millisecond)

	before := time.Now()
	w.Start()
	after := time.Now()

	if !w.active {
		t.Error("watchdog should be active after Start()")
	}
	if w.lastPing.Before(before) || w.lastPing.After(after) {
		t.Error("lastPing should be set to current time on Start()")
	}
}

func TestWatchdogPing(t *testing.T) {
	w := NewWatchdog(100 * time.Millisecond)
	w.Start()

	time.Sleep(10 * time.Millisecond)
	before := time.Now()
	w.Ping()
	after := time.Now()

	if w.lastPing.Before(before) || w.lastPing.After(after) {
		t.Error("lastPing should be updated on Ping()")
	}
}

func TestWatchdogCheckNotActive(t *testing.T) {
	w := NewWatchdog(1 * time.Millisecond)
	// Not started, so not active

	time.Sleep(5 * time.Millisecond)
	if w.Check() {
		t.Error("Check() should return false when watchdog is not active")
	}
}

func TestWatchdogCheckWithinTimeout(t *testing.T) {
	w := NewWatchdog(100 * time.Millisecond)
	w.Start()

	// Check immediately - should be within timeout
	if w.Check() {
		t.Error("Check() should return false when within timeout")
	}
}

func TestWatchdogCheckExceedsTimeout(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	w := NewWatchdog(10 * time.Millisecond)
	w.Start()

	time.Sleep(20 * time.Millisecond)
	if !w.Check() {
		t.Error("Check() should return true when timeout exceeded")
	}

	// Should have logged a warning
	if !bytes.Contains(buf.Bytes(), []byte("watchdog: frame exceeded timeout")) {
		t.Error("Check() should log a warning when timeout exceeded")
	}
}

func TestWatchdogStop(t *testing.T) {
	w := NewWatchdog(10 * time.Millisecond)
	w.Start()

	if !w.active {
		t.Error("watchdog should be active after Start()")
	}

	w.Stop()

	if w.active {
		t.Error("watchdog should not be active after Stop()")
	}

	// After Stop, Check should return false even if timeout exceeded
	time.Sleep(20 * time.Millisecond)
	if w.Check() {
		t.Error("Check() should return false after Stop()")
	}
}

func TestWatchdogPingPreventsTimeout(t *testing.T) {
	w := NewWatchdog(50 * time.Millisecond)
	w.Start()

	// Ping multiple times, each before timeout
	for i := 0; i < 5; i++ {
		time.Sleep(20 * time.Millisecond)
		if w.Check() {
			t.Errorf("Check() returned true on iteration %d, but we pinged in time", i)
		}
		w.Ping()
	}
}

func TestWatchdogRestart(t *testing.T) {
	w := NewWatchdog(50 * time.Millisecond)

	// Start, stop, start again
	w.Start()
	w.Stop()
	w.Start()

	if !w.active {
		t.Error("watchdog should be active after restart")
	}

	// Should not trigger timeout immediately after restart
	if w.Check() {
		t.Error("Check() should return false immediately after restart")
	}
}
