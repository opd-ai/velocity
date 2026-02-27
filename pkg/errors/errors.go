// Package errors provides structured error types for the velocity engine.
package errors

import "fmt"

// GameError represents a structured error with a code and context.
type GameError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface.
func (e *GameError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *GameError) Unwrap() error {
	return e.Err
}

// New creates a new GameError.
func New(code, message string) *GameError {
	return &GameError{Code: code, Message: message}
}

// Wrap wraps an existing error with a GameError.
func Wrap(code, message string, err error) *GameError {
	return &GameError{Code: code, Message: message, Err: err}
}
