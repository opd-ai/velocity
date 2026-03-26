package integration

import "testing"

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if r.hooks == nil {
		t.Error("hooks map should be initialized")
	}
	if len(r.hooks) != 0 {
		t.Errorf("new registry should have 0 hooks, got %d", len(r.hooks))
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()

	r.Register("oauth")
	r.Register("telemetry")
	r.Register("cdn")

	if len(r.hooks) != 3 {
		t.Errorf("expected 3 hooks, got %d", len(r.hooks))
	}

	// Hooks should be disabled by default
	for _, name := range []string{"oauth", "telemetry", "cdn"} {
		hook, ok := r.Get(name)
		if !ok {
			t.Errorf("hook %q not found", name)
			continue
		}
		if hook.Enabled {
			t.Errorf("hook %q should be disabled by default", name)
		}
		if hook.Name != name {
			t.Errorf("hook.Name = %q, want %q", hook.Name, name)
		}
	}
}

func TestRegistryEnable(t *testing.T) {
	r := NewRegistry()
	r.Register("oauth")

	hook, _ := r.Get("oauth")
	if hook.Enabled {
		t.Fatal("hook should start disabled")
	}

	r.Enable("oauth")

	hook, _ = r.Get("oauth")
	if !hook.Enabled {
		t.Error("hook should be enabled after Enable()")
	}
}

func TestRegistryDisable(t *testing.T) {
	r := NewRegistry()
	r.Register("oauth")
	r.Enable("oauth")

	hook, _ := r.Get("oauth")
	if !hook.Enabled {
		t.Fatal("hook should be enabled")
	}

	r.Disable("oauth")

	hook, _ = r.Get("oauth")
	if hook.Enabled {
		t.Error("hook should be disabled after Disable()")
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry()
	r.Register("telemetry")

	// Existing hook
	hook, ok := r.Get("telemetry")
	if !ok {
		t.Error("Get() should return true for existing hook")
	}
	if hook == nil {
		t.Fatal("Get() should return non-nil hook")
	}
	if hook.Name != "telemetry" {
		t.Errorf("hook.Name = %q, want %q", hook.Name, "telemetry")
	}

	// Non-existing hook
	hook, ok = r.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for non-existing hook")
	}
	if hook != nil {
		t.Error("Get() should return nil for non-existing hook")
	}
}

func TestRegistryEnableNonExistent(t *testing.T) {
	r := NewRegistry()

	// Should not panic
	defer func() {
		if rec := recover(); rec != nil {
			t.Errorf("Enable() panicked for non-existent hook: %v", rec)
		}
	}()

	r.Enable("nonexistent")
}

func TestRegistryDisableNonExistent(t *testing.T) {
	r := NewRegistry()

	// Should not panic
	defer func() {
		if rec := recover(); rec != nil {
			t.Errorf("Disable() panicked for non-existent hook: %v", rec)
		}
	}()

	r.Disable("nonexistent")
}

func TestServiceHookStruct(t *testing.T) {
	hook := ServiceHook{
		Name:    "custom",
		Enabled: true,
	}

	if hook.Name != "custom" {
		t.Errorf("Name = %q, want %q", hook.Name, "custom")
	}
	if !hook.Enabled {
		t.Error("Enabled should be true")
	}
}

func TestRegistryMultipleOperations(t *testing.T) {
	r := NewRegistry()

	r.Register("a")
	r.Register("b")
	r.Register("c")

	r.Enable("a")
	r.Enable("b")

	aHook, _ := r.Get("a")
	bHook, _ := r.Get("b")
	cHook, _ := r.Get("c")

	if !aHook.Enabled {
		t.Error("hook a should be enabled")
	}
	if !bHook.Enabled {
		t.Error("hook b should be enabled")
	}
	if cHook.Enabled {
		t.Error("hook c should be disabled")
	}

	r.Disable("a")
	aHook, _ = r.Get("a")
	if aHook.Enabled {
		t.Error("hook a should be disabled after Disable()")
	}
}
