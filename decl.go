package gojen

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Strategy is a type that represents the strategy for setting a template.
// ENUM(trunc,append,append_at_last,init)
// trunc: Truncate the destination file. (please commit this file to git before running gojen)
// append: Append at the end to the destination file.
// ignore: Ignore the destination file.
//
//go:generate go-enum -f=$GOFILE --marshal --names --values
type Strategy string

type (
	// E represents a element.
	E struct {
		Path     string         `json:"path" yaml:"path"`
		Name     string         `json:"name" yaml:"name"`
		Alias    string         `json:"alias" yaml:"alias"`
		Require  []string       `json:"require" yaml:"require"`
		Args     map[string]any `json:"args" yaml:"args"`
		Template string         `json:"template" yaml:"template"`
		Strategy Strategy       `json:"strategy" yaml:"strategy"`
		Confirm  bool           `json:"confirm" yaml:"confirm"`
	}

	// D represents a declaration for templates.
	D struct {
		Path        string         `json:"path" yaml:"path"`
		Name        string         `json:"name" yaml:"name"`
		Require     []string       `json:"require" yaml:"require"`
		Args        map[string]any `json:"args" yaml:"args"`
		Elements    []*E           `json:"elements" yaml:"elements"`
		Description string         `json:"description" yaml:"description"`
		selected    string
	}

	// M represents a map of template definitions.
	M map[string]*D
)

func (e *E) Validate() error {
	if e.Path == "" {
		return fmt.Errorf("path is required")
	}
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func (d *D) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(d.Elements) == 0 {
		return fmt.Errorf("elements are required")
	}

	for _, e := range d.Elements {
		if err := e.Validate(); err != nil {
			return err
		}
	}

	return nil
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

// GetElements returns the element with the given name.
func (d *D) GetElements(name string) *E {
	for _, el := range d.Elements {
		if el.Name == name {
			return el
		}
	}

	return nil
}

// performSelect performs the template selection based on the sequence.
// It will ask the user to select the template if there are multiple options.
// TODO: handle this
// func (d *D) performSelect(fp string, seq *Sequence) (*E, error) {
// 	items := d.Elements
// 	if len(items) == 0 {
// 		return nil, nil
// 	}
//
// 	// filter the selected template from sequence.
// 	if seq != nil && seq.DeclName == d.Name && len(seq.elements) > 0 {
// 		items = make([]*E, 0)
// 		for _, sel := range d.Elements {
// 			if _, ok := seq.elements[sel.Name]; ok {
// 				items = append(items, sel)
// 			}
// 		}
// 	}
//
// 	// Ask user to select template if there are multiple options.
// 	// if only one option, select it automatically.
// 	var input int
// 	if len(items) == 1 {
// 		input = 1
// 	} else {
// 		// console.Dangerf("Please select one of the following templates for '%s':\n", d.Name)
// 		// for i, s := range items {
// 		// 	console.Infof("Option %d: %s\n", i+1, fp)
// 		// 	console.Successf("%s\n", s.Template)
// 		// }
// 		//
// 		// console.Dangerf("Enter the option number: ")
// 		// if _, err := fmt.Scanln(&input); err != nil || input < 1 || input > len(items) {
// 		// 	return nil, fmt.Errorf("invalid option selected")
// 		// }
// 	}
//
// 	// provide required context for the selected template.
// 	e := items[input-1]
// 	// if err := d.provideRequireCtx(e.Require); err != nil {
// 	// 	return nil, err
// 	// }
//
// 	if !e.Strategy.IsValid() {
// 		return nil, fmt.Errorf("invalid handle file strategy")
// 	}
//
// 	d.selected = e.Name
// 	return e, nil
// }
//
// // mergeContext merges the global context with the template context.
// // It returns the merged context.
// // The template context takes precedence over the global context.
// func (d *D) mergeContext(ctx map[string]any) *D {
// 	d.Args = mergeMaps(ctx, d.Args)
//
// 	return d
// }
