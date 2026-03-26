package benchmark

import (
	"testing"
	"time"
)

func TestRunBasic(t *testing.T) {
	counter := 0
	result := Run("test-system", 100, func() {
		counter++
	})

	if result.SystemName != "test-system" {
		t.Errorf("SystemName = %q, want %q", result.SystemName, "test-system")
	}
	if result.Iterations != 100 {
		t.Errorf("Iterations = %d, want 100", result.Iterations)
	}
	if counter != 100 {
		t.Errorf("counter = %d, want 100 (function should be called iterations times)", counter)
	}
}

func TestRunTiming(t *testing.T) {
	result := Run("sleep-test", 10, func() {
		time.Sleep(1 * time.Millisecond)
	})

	// Should take at least 10ms total
	if result.TotalMs < 10.0 {
		t.Errorf("TotalMs = %f, want >= 10.0", result.TotalMs)
	}
	// Average should be at least 1ms
	if result.AvgMs < 1.0 {
		t.Errorf("AvgMs = %f, want >= 1.0", result.AvgMs)
	}
}

func TestRunZeroIterations(t *testing.T) {
	called := false
	result := Run("zero-iter", 0, func() {
		called = true
	})

	if called {
		t.Error("function should not be called with 0 iterations")
	}
	if result.Iterations != 0 {
		t.Errorf("Iterations = %d, want 0", result.Iterations)
	}
	if result.AvgMs != 0 {
		t.Errorf("AvgMs = %f, want 0 (avoid division by zero)", result.AvgMs)
	}
}

func TestRunResultFields(t *testing.T) {
	result := Run("field-test", 50, func() {
		// Fast no-op
	})

	if result.SystemName != "field-test" {
		t.Errorf("SystemName = %q, want %q", result.SystemName, "field-test")
	}
	if result.Iterations != 50 {
		t.Errorf("Iterations = %d, want 50", result.Iterations)
	}
	if result.TotalMs < 0 {
		t.Errorf("TotalMs = %f, should be non-negative", result.TotalMs)
	}
	if result.AvgMs < 0 {
		t.Errorf("AvgMs = %f, should be non-negative", result.AvgMs)
	}
}

func TestRunAverageCalculation(t *testing.T) {
	// Test that average is correctly calculated
	result := Run("avg-test", 100, func() {
		// Very fast operation
	})

	// Average should be TotalMs / Iterations
	expectedAvg := result.TotalMs / float64(result.Iterations)
	diff := result.AvgMs - expectedAvg
	if diff < -0.001 || diff > 0.001 {
		t.Errorf("AvgMs = %f, expected %f (TotalMs/Iterations)", result.AvgMs, expectedAvg)
	}
}

func TestResultStruct(t *testing.T) {
	r := Result{
		SystemName: "manual",
		Iterations: 42,
		TotalMs:    84.0,
		AvgMs:      2.0,
	}

	if r.SystemName != "manual" {
		t.Errorf("SystemName = %q, want %q", r.SystemName, "manual")
	}
	if r.Iterations != 42 {
		t.Errorf("Iterations = %d, want 42", r.Iterations)
	}
	if r.TotalMs != 84.0 {
		t.Errorf("TotalMs = %f, want 84.0", r.TotalMs)
	}
	if r.AvgMs != 2.0 {
		t.Errorf("AvgMs = %f, want 2.0", r.AvgMs)
	}
}
