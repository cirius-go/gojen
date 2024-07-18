package gojen

import (
	"strings"
)

// D represents a template definition.
type D struct {
	Path            string
	Context         map[string]any
	TemplateString  string
	Strategy        Strategy
	RequiredContext []string
	Dependencies    []string
	Description     string
}

// JSOND represents a JSON template definition.
type JSOND struct {
	Path            string         `json:"path"`
	TemplateBody    []string       `json:"templateString"`
	Strategy        Strategy       `json:"strategy"`
	RequiredContext []string       `json:"requiredContext"`
	Dependencies    []string       `json:"dependencies"`
	Context         map[string]any `json:"context"`
	Description     string         `json:"description"`
}

func (j *JSOND) toD() *D {
	return &D{
		Path:            j.Path,
		TemplateString:  strings.Join(j.TemplateBody, "\n"),
		Context:         j.Context,
		Strategy:        j.Strategy,
		RequiredContext: j.RequiredContext,
		Dependencies:    j.Dependencies,
		Description:     j.Description,
	}
}

func (j *D) toJSOND() *JSOND {
	return &JSOND{
		Path:            j.Path,
		TemplateBody:    strings.Split(j.TemplateString, "\n"),
		Context:         j.Context,
		Strategy:        j.Strategy,
		RequiredContext: j.RequiredContext,
		Dependencies:    j.Dependencies,
		Description:     j.Description,
	}
}
