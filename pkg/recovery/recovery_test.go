package recovery

import (
	"bytes"
	"log"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	DefaultHandler("test panic value")

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("test panic value")) {
		t.Errorf("DefaultHandler output should contain panic value, got: %s", output)
	}
}

func TestWithRecovery_NoPanic(t *testing.T) {
	called := false
	handlerCalled := false

	WithRecovery(func() {
		called = true
	}, func(r interface{}) {
		handlerCalled = true
	})

	if !called {
		t.Error("wrapped function was not called")
	}
	if handlerCalled {
		t.Error("handler should not be called when no panic occurs")
	}
}

func TestWithRecovery_WithPanic(t *testing.T) {
	var recovered interface{}

	WithRecovery(func() {
		panic("test panic")
	}, func(r interface{}) {
		recovered = r
	})

	if recovered != "test panic" {
		t.Errorf("recovered = %v, want 'test panic'", recovered)
	}
}

func TestWithRecovery_NilHandler(t *testing.T) {
	// Should not panic when handler is nil
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("WithRecovery should not propagate panic: %v", r)
		}
	}()

	// Capture log output to avoid noise
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	WithRecovery(func() {
		panic("test panic with nil handler")
	}, nil)

	// Should use DefaultHandler
	if !bytes.Contains(buf.Bytes(), []byte("test panic with nil handler")) {
		t.Error("DefaultHandler should be called when handler is nil")
	}
}

func TestWithRecovery_PanicTypes(t *testing.T) {
	tests := []struct {
		name       string
		panicValue interface{}
	}{
		{"string panic", "string error"},
		{"int panic", 42},
		{"error panic", struct{ msg string }{"custom error"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var recovered interface{}

			WithRecovery(func() {
				panic(tt.panicValue)
			}, func(r interface{}) {
				recovered = r
			})

			if recovered != tt.panicValue {
				t.Errorf("recovered = %v, want %v", recovered, tt.panicValue)
			}
		})
	}
}
