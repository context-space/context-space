package registry

import (
	"sync"

	"github.com/context-space/context-space/backend/internal/provideradapter/application"
)

var (
	// Global registry for adapter templates
	templateRegistry = NewAdapterTemplateRegistry()
)

// AdapterTemplateRegistry manages the registration of adapter templates
type AdapterTemplateRegistry struct {
	templates map[string]application.AdapterTemplate
	mu        sync.RWMutex
}

// NewAdapterTemplateRegistry creates a new template registry
func NewAdapterTemplateRegistry() *AdapterTemplateRegistry {
	return &AdapterTemplateRegistry{
		templates: make(map[string]application.AdapterTemplate),
		mu:        sync.RWMutex{},
	}
}

// RegisterTemplate registers an adapter template
func (r *AdapterTemplateRegistry) RegisterTemplate(name string, template application.AdapterTemplate) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[name] = template
}

// GetTemplate returns a registered template by name
func (r *AdapterTemplateRegistry) GetTemplate(name string) (application.AdapterTemplate, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	template, ok := r.templates[name]
	return template, ok
}

// ListTemplates returns the names of all registered templates
func (r *AdapterTemplateRegistry) ListTemplates() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.templates))
	for name := range r.templates {
		names = append(names, name)
	}

	return names
}

// RegisterAdapterTemplate registers a template with the global registry
func RegisterAdapterTemplate(name string, template application.AdapterTemplate) {
	templateRegistry.RegisterTemplate(name, template)
}

// GetAdapterTemplate returns a template from the global registry
func GetAdapterTemplate(name string) (application.AdapterTemplate, bool) {
	return templateRegistry.GetTemplate(name)
}

// ListAdapterTemplates returns the names of all templates in the global registry
func ListAdapterTemplates() []string {
	return templateRegistry.ListTemplates()
}
