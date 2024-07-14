package gojen

// Gojen
type Gojen struct {
	cfg       *Config
	context   map[string]any
	templates map[string]*Template
}

// New returns a new Gojen instance.
func New(cfg *Config) *Gojen {
	return &Gojen{
		cfg:       cfg,
		context:   make(map[string]any),
		templates: make(map[string]*Template),
	}
}

// SetTemplate sets a template.
func (g *Gojen) SetTemplate(name string, template *Template) {
	g.templates[name] = template
}

// AddContext adds a key-value pair to the context map.
func (g *Gojen) AddContext(key string, value any) {
	g.context[key] = value
}

// Build builds the templates.
func (g *Gojen) Build() error {

	return nil
}
