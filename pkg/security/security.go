// Package security provides E2E encrypted transport and authentication tokens.
//
// TODO(v5.0): This package is a stub. Full implementation planned for v5.0.
// See ROADMAP.md v5.0 milestone and GAPS.md "v5.0+ Features Are Stubs" section.
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
func Encrypt(data, key []byte) ([]byte, error) {
	return nil, ErrNotImplemented
}

// Decrypt decrypts a byte slice using E2E encryption.
func Decrypt(data, key []byte) ([]byte, error) {
	return nil, ErrNotImplemented
}
