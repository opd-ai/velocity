// Package mods provides a mod loader and scripted mod API.
package mods

// Mod represents a loaded mod.
type Mod struct {
	Name    string
	Version string
	Enabled bool
}

// Loader discovers and loads mods.
type Loader struct {
	mods []Mod
}

// NewLoader creates a new mod loader.
func NewLoader() *Loader {
	return &Loader{}
}

// Register adds a mod to the loader.
func (l *Loader) Register(name, version string) {
	l.mods = append(l.mods, Mod{Name: name, Version: version, Enabled: true})
}

// List returns all registered mods.
func (l *Loader) List() []Mod {
	return l.mods
}

// Enable enables a mod by name.
func (l *Loader) Enable(name string) {
	for i := range l.mods {
		if l.mods[i].Name == name {
			l.mods[i].Enabled = true
			return
		}
	}
}

// Disable disables a mod by name.
func (l *Loader) Disable(name string) {
	for i := range l.mods {
		if l.mods[i].Name == name {
			l.mods[i].Enabled = false
			return
		}
	}
}
