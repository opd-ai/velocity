// Package recovery provides panic recovery middleware for the game loop.
package recovery

import (
	"log"
)

// Handler is a function called when a panic is recovered.
type Handler func(recovered interface{})

// DefaultHandler logs the recovered panic value.
func DefaultHandler(recovered interface{}) {
	log.Printf("recovered from panic: %v", recovered)
}

// WithRecovery wraps a function with panic recovery.
func WithRecovery(fn func(), handler Handler) {
	defer func() {
		if r := recover(); r != nil {
			if handler != nil {
				handler(r)
			} else {
				DefaultHandler(r)
			}
		}
	}()
	fn()
}
