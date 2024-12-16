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
		Path                 string   `json:"path" yaml:"path" validate:"required"`
		Name                 string   `json:"name" yaml:"name" validate:"required"`
		Alias                string   `json:"alias" yaml:"alias"`
		Require              []string `json:"require" yaml:"require"`
		Args                 Args     `json:"args" yaml:"args" validate:"required"`
		Template             string   `json:"template" yaml:"template" validate:"required"`
		Strategy             Strategy `json:"strategy" yaml:"strategy" validate:"required"`
		IgnoreComparingLines []string `json:"ignore_comparing_lines" yaml:"ignore_comparing_lines"`
	}

	// D represents a declaration for templates.
	D struct {
		Path        string   `json:"path" yaml:"path"`
		Name        string   `json:"name" yaml:"name"`
		Require     []string `json:"require" yaml:"require"`
		Args        Args     `json:"args" yaml:"args"`
		Elements    []*E     `json:"elements" yaml:"elements"`
		Description string   `json:"description" yaml:"description"`
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
