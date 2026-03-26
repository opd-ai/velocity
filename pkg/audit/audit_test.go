package audit

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}
	if len(logger.Samples()) != 0 {
		t.Errorf("new logger should have 0 samples, got %d", len(logger.Samples()))
	}
}

func TestLoggerRecord(t *testing.T) {
	logger := NewLogger()

	before := time.Now()
	logger.Record(16.5, 100)
	after := time.Now()

	samples := logger.Samples()
	if len(samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(samples))
	}

	sample := samples[0]
	if sample.DurationMs != 16.5 {
		t.Errorf("DurationMs = %f, want 16.5", sample.DurationMs)
	}
	if sample.EntityCount != 100 {
		t.Errorf("EntityCount = %d, want 100", sample.EntityCount)
	}
	if sample.Timestamp.Before(before) || sample.Timestamp.After(after) {
		t.Errorf("Timestamp %v not in expected range [%v, %v]", sample.Timestamp, before, after)
	}
}

func TestLoggerRecordMultiple(t *testing.T) {
	logger := NewLogger()

	logger.Record(16.0, 50)
	logger.Record(17.0, 75)
	logger.Record(18.0, 100)

	samples := logger.Samples()
	if len(samples) != 3 {
		t.Fatalf("expected 3 samples, got %d", len(samples))
	}

	expectedDurations := []float64{16.0, 17.0, 18.0}
	expectedCounts := []int{50, 75, 100}

	for i, sample := range samples {
		if sample.DurationMs != expectedDurations[i] {
			t.Errorf("sample[%d].DurationMs = %f, want %f", i, sample.DurationMs, expectedDurations[i])
		}
		if sample.EntityCount != expectedCounts[i] {
			t.Errorf("sample[%d].EntityCount = %d, want %d", i, sample.EntityCount, expectedCounts[i])
		}
	}
}

func TestLoggerReset(t *testing.T) {
	logger := NewLogger()

	logger.Record(16.0, 50)
	logger.Record(17.0, 75)

	if len(logger.Samples()) != 2 {
		t.Fatalf("expected 2 samples before reset, got %d", len(logger.Samples()))
	}

	logger.Reset()

	if len(logger.Samples()) != 0 {
		t.Errorf("expected 0 samples after reset, got %d", len(logger.Samples()))
	}
}

func TestLoggerResetPreservesCapacity(t *testing.T) {
	logger := NewLogger()

	// Record several samples
	for i := 0; i < 100; i++ {
		logger.Record(float64(i), i)
	}

	logger.Reset()

	// Should be able to record again without issues
	logger.Record(1.0, 1)

	if len(logger.Samples()) != 1 {
		t.Errorf("expected 1 sample after reset and new record, got %d", len(logger.Samples()))
	}
}

func TestFrameSampleFields(t *testing.T) {
	now := time.Now()
	sample := FrameSample{
		Timestamp:   now,
		DurationMs:  16.67,
		EntityCount: 200,
	}

	if sample.Timestamp != now {
		t.Errorf("Timestamp = %v, want %v", sample.Timestamp, now)
	}
	if sample.DurationMs != 16.67 {
		t.Errorf("DurationMs = %f, want 16.67", sample.DurationMs)
	}
	if sample.EntityCount != 200 {
		t.Errorf("EntityCount = %d, want 200", sample.EntityCount)
	}
}
