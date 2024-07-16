package gojen

// D represents a template definition.
type D struct {
	Path            string
	Name            string
	TemplateString  string
	Strategy        Strategy
	RequiredContext []string
	Dependencies    []string
}
