// Package integration provides external service integration hooks
// for OAuth identity, CDN asset delivery, and telemetry pipelines.
package integration

// ServiceHook represents an external service integration point.
type ServiceHook struct {
	Name    string
	Enabled bool
}

// Registry holds all registered service hooks.
type Registry struct {
	hooks map[string]*ServiceHook
}

// NewRegistry creates a new integration registry.
func NewRegistry() *Registry {
	return &Registry{hooks: make(map[string]*ServiceHook)}
}

// Register adds a new service hook.
func (r *Registry) Register(name string) {
	r.hooks[name] = &ServiceHook{Name: name, Enabled: false}
}

// Enable activates a registered hook by name.
func (r *Registry) Enable(name string) {
	if h, ok := r.hooks[name]; ok {
		h.Enabled = true
	}
}

// Disable deactivates a registered hook by name.
func (r *Registry) Disable(name string) {
	if h, ok := r.hooks[name]; ok {
		h.Enabled = false
	}
}

// Get returns a hook by name.
func (r *Registry) Get(name string) (*ServiceHook, bool) {
	h, ok := r.hooks[name]
	return h, ok
}
