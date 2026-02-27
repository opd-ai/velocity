// Package security provides E2E encrypted transport and authentication tokens.
package security

import "errors"

// Token represents an authentication token.
type Token struct {
	Value   string
	Expires int64
}

// NewToken creates a new placeholder token.
func NewToken(value string) *Token {
	return &Token{Value: value}
}

// ErrNotImplemented is returned by stub cryptographic functions.
var ErrNotImplemented = errors.New("security: encryption not yet implemented")

// Encrypt encrypts a byte slice using E2E encryption.
func Encrypt(data []byte, key []byte) ([]byte, error) {
	return nil, ErrNotImplemented
}

// Decrypt decrypts a byte slice using E2E encryption.
func Decrypt(data []byte, key []byte) ([]byte, error) {
	return nil, ErrNotImplemented
}
