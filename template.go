package gojen

import (
	"strings"
)

// D represents a template definition.
type D struct {
	Path            string
	Name            string
	Context         map[string]any
	TemplateString  string
	Select          []string
	Strategy        Strategy
	RequiredContext []string
	Dependencies    []string
	Description     string
	Confirm         bool
}

func (d *D) isCtxSatisfied() (bool, []string) {
	notSatisfied := []string{}
	for _, key := range d.RequiredContext {
		if _, exists := d.Context[key]; !exists {
			notSatisfied = append(notSatisfied, key)
		}
	}
	return len(notSatisfied) == 0, notSatisfied
}

// mergeGlobalCtx merges the global context with the template context.
// It returns the merged context.
// The template context takes precedence over the global context.
func (d *D) mergeGlobalCtx(ctx map[string]any) *D {
	d.Context = mergeMaps(ctx, d.Context)

	return d
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
	Confirm         bool           `json:"confirm"`
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
		Confirm:         j.Confirm,
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
		Confirm:         j.Confirm,
	}
}
