package gameerrors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("TEST_CODE", "test message")

	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Code != "TEST_CODE" {
		t.Errorf("expected code 'TEST_CODE', got %s", err.Code)
	}
	if err.Message != "test message" {
		t.Errorf("expected message 'test message', got %s", err.Message)
	}
	if err.Err != nil {
		t.Error("expected nil wrapped error")
	}
}

func TestWrap(t *testing.T) {
	inner := errors.New("inner error")
	err := Wrap("WRAP_CODE", "outer message", inner)

	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Code != "WRAP_CODE" {
		t.Errorf("expected code 'WRAP_CODE', got %s", err.Code)
	}
	if err.Message != "outer message" {
		t.Errorf("expected message 'outer message', got %s", err.Message)
	}
	if err.Err != inner {
		t.Error("expected wrapped error to be preserved")
	}
}

func TestGameError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *GameError
		expected string
	}{
		{
			name:     "without wrapped error",
			err:      New("CODE", "message"),
			expected: "[CODE] message",
		},
		{
			name:     "with wrapped error",
			err:      Wrap("CODE", "message", errors.New("inner")),
			expected: "[CODE] message: inner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("Error() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestGameError_Unwrap(t *testing.T) {
	inner := errors.New("inner error")
	err := Wrap("CODE", "message", inner)

	unwrapped := err.Unwrap()
	if unwrapped != inner {
		t.Error("Unwrap() should return inner error")
	}
}

func TestGameError_Unwrap_Nil(t *testing.T) {
	err := New("CODE", "message")

	unwrapped := err.Unwrap()
	if unwrapped != nil {
		t.Error("Unwrap() should return nil for non-wrapped error")
	}
}

func TestGameError_ErrorsIs(t *testing.T) {
	inner := errors.New("specific error")
	err := Wrap("CODE", "message", inner)

	if !errors.Is(err, inner) {
		t.Error("errors.Is should match inner error")
	}
}

func TestGameError_Fields(t *testing.T) {
	ge := &GameError{
		Code:    "FIELD_TEST",
		Message: "field test message",
		Err:     nil,
	}

	if ge.Code != "FIELD_TEST" {
		t.Error("code mismatch")
	}
	if ge.Message != "field test message" {
		t.Error("message mismatch")
	}
}

func TestNew_CommonCodes(t *testing.T) {
	codes := []string{
		"CONFIG",
		"WAVE_GEN",
		"RENDER",
		"AUDIO",
		"INPUT",
		"SAVE",
		"LOAD",
	}

	for _, code := range codes {
		t.Run(code, func(t *testing.T) {
			err := New(code, "test")
			if err.Code != code {
				t.Errorf("code mismatch: got %s, want %s", err.Code, code)
			}
		})
	}
}

func TestWrap_ChainedErrors(t *testing.T) {
	level1 := errors.New("root cause")
	level2 := Wrap("LEVEL2", "middle layer", level1)
	level3 := Wrap("LEVEL3", "top layer", level2)

	// Should be able to unwrap through the chain
	if !errors.Is(level3, level1) {
		t.Error("should be able to find root error through chain")
	}
	if !errors.Is(level3, level2) {
		t.Error("should be able to find middle error through chain")
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New("CODE", "message")
	}
}

func BenchmarkWrap(b *testing.B) {
	inner := errors.New("inner")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Wrap("CODE", "message", inner)
	}
}

func BenchmarkGameError_Error(b *testing.B) {
	err := Wrap("CODE", "message", errors.New("inner"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
