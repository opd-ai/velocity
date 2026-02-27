// Package audit provides frame-time telemetry and entity-count logging.
package audit

import "time"

// FrameSample records timing data for a single frame.
type FrameSample struct {
	Timestamp   time.Time
	DurationMs  float64
	EntityCount int
}

// Logger collects audit samples.
type Logger struct {
	samples []FrameSample
}

// NewLogger creates a new audit logger.
func NewLogger() *Logger {
	return &Logger{}
}

// Record adds a frame sample to the log.
func (l *Logger) Record(duration float64, entityCount int) {
	l.samples = append(l.samples, FrameSample{
		Timestamp:   time.Now(),
		DurationMs:  duration,
		EntityCount: entityCount,
	})
}

// Samples returns all recorded samples.
func (l *Logger) Samples() []FrameSample {
	return l.samples
}

// Reset clears all recorded samples.
func (l *Logger) Reset() {
	l.samples = l.samples[:0]
}
