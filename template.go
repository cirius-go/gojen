package gojen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cirius-go/generic/slice"

	"github.com/cirius-go/gojen/color"
)

// D represents a template definition.
type D struct {
	Path         string         `json:"path"`
	Require      []string       `json:"require"`
	Name         string         `json:"name"`
	Context      map[string]any `json:"context"`
	Select       []*DItem       `json:"select"`
	Dependencies []string       `json:"dependencies"`
	Description  string         `json:"description"`
}

// MarshalJSON marshals the D struct to JSON.
func (d *D) MarshalJSON() ([]byte, error) {
	type Alias D

	type ADItem struct {
		Template []string `json:"template"`
		Require  []string `json:"require"`
		Strategy Strategy `json:"strategy"`
		Confirm  bool     `json:"confirm"`
	}

	return json.Marshal(&struct {
		Select []*ADItem `json:"select"`
		*Alias
	}{
		Select: slice.Map(func(e *DItem) *ADItem {
			return &ADItem{
				Template: strings.Split(e.Template, "\n"),
				Require:  e.Require,
				Strategy: e.Strategy,
				Confirm:  e.Confirm,
			}
		}, d.Select...),
		Alias: (*Alias)(d),
	})
}

// UnmarshalJSON unmarshals the D struct from JSON.
func (d *D) UnmarshalJSON(data []byte) error {
	type Alias D
	type ADItem struct {
		Template []string `json:"template"`
		Require  []string `json:"require"`
		Strategy Strategy `json:"strategy"`
		Confirm  bool     `json:"confirm"`
	}
	aux := &struct {
		Select []*ADItem `json:"select"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	d.Select = slice.Map(func(e *ADItem) *DItem {
		return &DItem{
			Template: strings.Join(e.Template, "\n"),
			Require:  e.Require,
			Confirm:  e.Confirm,
			Strategy: e.Strategy,
		}
	}, aux.Select...)
	return nil
}

// DItem represents a template item.
type DItem struct {
	Template string   `json:"template"`
	Require  []string `json:"require"`
	Strategy Strategy `json:"strategy"`
	Confirm  bool     `json:"confirm"`
}

func (d *D) tryProvideCtx(required []string) error {
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

// mergeContext merges the global context with the template context.
// It returns the merged context.
// The template context takes precedence over the global context.
func (d *D) mergeContext(ctx map[string]any) *D {
	d.Context = mergeMaps(ctx, d.Context)

	return d
}
