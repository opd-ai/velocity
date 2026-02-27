// Package benchmark provides a per-system micro-benchmark harness.
package benchmark

import "time"

// Result holds the result of a benchmark run.
type Result struct {
	SystemName string
	Iterations int
	TotalMs    float64
	AvgMs      float64
}

// Run benchmarks a function for the given number of iterations.
func Run(systemName string, iterations int, fn func()) Result {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	total := time.Since(start).Seconds() * 1000
	return Result{
		SystemName: systemName,
		Iterations: iterations,
		TotalMs:    total,
		AvgMs:      total / float64(iterations),
	}
}
