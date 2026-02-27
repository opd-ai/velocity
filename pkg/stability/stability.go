// Package stability provides crash detection, watchdog timer, and auto-save triggers.
package stability

import (
	"log"
	"time"
)

// Watchdog monitors the game loop and triggers recovery if a frame takes too long.
type Watchdog struct {
	timeout  time.Duration
	lastPing time.Time
	active   bool
}

// NewWatchdog creates a new watchdog with the given timeout.
func NewWatchdog(timeout time.Duration) *Watchdog {
	return &Watchdog{timeout: timeout}
}

// Start activates the watchdog.
func (w *Watchdog) Start() {
	w.active = true
	w.lastPing = time.Now()
}

// Ping resets the watchdog timer. Call once per frame.
func (w *Watchdog) Ping() {
	w.lastPing = time.Now()
}

// Check returns true if the watchdog has not been pinged within the timeout.
func (w *Watchdog) Check() bool {
	if !w.active {
		return false
	}
	if time.Since(w.lastPing) > w.timeout {
		log.Printf("watchdog: frame exceeded timeout (%v)", w.timeout)
		return true
	}
	return false
}

// Stop deactivates the watchdog.
func (w *Watchdog) Stop() {
	w.active = false
}
