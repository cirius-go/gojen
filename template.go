package gojen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cirius-go/gojen/color"
)

type (
	// E represents a template item.
	E struct {
		Template string   `json:"template" yaml:"template"`
		Require  []string `json:"require" yaml:"require"`
		Strategy Strategy `json:"strategy" yaml:"strategy"`
		Confirm  bool     `json:"confirm" yaml:"confirm"`

		index int
	}

	// D represents a template definition.
	D struct {
		Path         string         `json:"path" yaml:"path"`
		Required     []string       `json:"require" yaml:"required"`
		Name         string         `json:"name" yaml:"name"`
		Context      map[string]any `json:"context" yaml:"context"`
		Select       []*E           `json:"select" yaml:"select"`
		Dependencies []string       `json:"dependencies" yaml:"dependencies"`
		Description  string         `json:"description" yaml:"description"`

		selectedTmpl int
	}

	// M represents a map of template definitions.
	M map[string]*D
)

// Store stores the template definition in the map.
func (m M) Store(d *D) M {
	for _, s := range d.Select {
		if s.Strategy == "" {
			s.Strategy = StrategyIgnore
		}
		if !s.Strategy.IsValid() {
			panic(fmt.Sprintf("invalid strategy: %s", s.Strategy))
		}
	}
	m[d.Name] = d
	return m
}

// MarshalJSON marshals the E struct to JSON.
func (e *E) MarshalJSON() ([]byte, error) {
	type Alias E

	return json.Marshal(&struct {
		*Alias
		Template []string `json:"template"`
	}{
		Alias:    (*Alias)(e),
		Template: strings.Split(e.Template, "\n"),
	})
}

// UnmarshalJSON unmarshals the E struct from JSON.
func (e *E) UnmarshalJSON(data []byte) error {
	type Alias E
	aux := &struct {
		Template []string `json:"template"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.Template = strings.Join(aux.Template, "\n")
	return nil
}

func (e *E) clone() *E {
	if e == nil {
		return nil
	}
	return &E{Template: e.Template, Require: cloneSlice(e.Require), Strategy: e.Strategy, Confirm: e.Confirm}
}

func (d *D) clone() *D {
	if d == nil {
		return nil
	}
	es := make([]*E, len(d.Select))
	for i, e := range d.Select {
		es[i] = e.clone()
	}
	return &D{
		Path:         d.Path,
		Required:     cloneSlice(d.Required),
		Name:         d.Name,
		Context:      cloneMap(d.Context),
		Select:       es,
		Dependencies: cloneSlice(d.Dependencies),
		Description:  d.Description,
	}
}

func (d *D) provideRequireCtx(required []string) error {
	for _, key := range required {
		if _, exists := d.Context[key]; !exists {
			fmt.Print(color.Redf("Provide missing value for context key '%s': ", key))
			var val string
			_, err := fmt.Scanln(&val)
			if err != nil {
				return err
			}

			if val != "" {
				d.Context[key] = val
			}
		}
	}

	return nil
}

// performSelect performs the template selection based on the sequence.
// It will ask the user to select the template if there are multiple options.
func (d *D) performSelect(fp string, seq *sequence) (*E, error) {
	// indexing
	for i, s := range d.Select {
		s.index = i + 1
	}

	if len(d.Select) == 0 {
		return nil, fmt.Errorf("no template to select")
	}

	items := d.Select
	if len(items) == 0 {
		return nil, nil
	}

	if seq != nil && seq.n == d.Name {
		items = seq.filter(items)
	}

	// Determine the selected template index
	var input int
	if len(items) == 1 {
		input = 1
	} else {
		fmt.Printf("Please select one of the following templates for '%s':\n", d.Name)
		for i, s := range items {
			fmt.Printf("%s%s\n\n", color.Bluef("Option %d: %s\n", i+1, fp), color.Greenf(s.Template))
		}

		fmt.Printf("Enter the option number: ")
		if _, err := fmt.Scanln(&input); err != nil || input < 1 || input > len(items) {
			return nil, fmt.Errorf("invalid option selected")
		}
	}

	e := items[input-1]
	if err := d.provideRequireCtx(e.Require); err != nil {
		return nil, err
	}

	if !e.Strategy.IsValid() {
		return nil, fmt.Errorf("invalid handle file strategy")
	}

	d.selectedTmpl = e.index

	return e, nil
}

// mergeContext merges the global context with the template context.
// It returns the merged context.
// The template context takes precedence over the global context.
func (d *D) mergeContext(ctx map[string]any) *D {
	d.Context = mergeMaps(ctx, d.Context)

	return d
}
