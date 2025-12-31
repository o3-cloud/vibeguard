// Package templates provides predefined configuration templates for different project types.
package templates

import (
	"fmt"
	"sort"
)

// Template represents a predefined vibeguard configuration template.
type Template struct {
	Name        string // Unique identifier (e.g., "go-standard")
	Description string // Human-readable description
	Content     string // YAML configuration content
}

// registry holds all available templates
var registry = map[string]Template{}

// Register adds a template to the registry. Panics if duplicate name.
func Register(t Template) {
	if _, exists := registry[t.Name]; exists {
		panic(fmt.Sprintf("template %q already registered", t.Name))
	}
	registry[t.Name] = t
}

// Get returns a template by name, or an error if not found.
func Get(name string) (Template, error) {
	t, ok := registry[name]
	if !ok {
		return Template{}, fmt.Errorf("template %q not found", name)
	}
	return t, nil
}

// List returns all registered templates sorted by name.
func List() []Template {
	templates := make([]Template, 0, len(registry))
	for _, t := range registry {
		templates = append(templates, t)
	}
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})
	return templates
}

// Names returns all registered template names sorted alphabetically.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Exists returns true if a template with the given name exists.
func Exists(name string) bool {
	_, ok := registry[name]
	return ok
}
