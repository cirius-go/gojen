package gojen

import (
	"fmt"
)

// Strategy is a type that represents the strategy for setting a template.
// ENUM(init,prepend_at_head,prepend,append,append_at_pos,edit)
// init: Create file and set content by template. If this file exists, ignore.
// prepend_at_head: Prepend content at head of file.
// prepend: Prepend at anchor position.
// append: Append at anchor position.
// append_at_pos: append at the end of file.
// edit: Edit at anchor position.
// output: Output of a seq.
//
//go:generate go-enum -f=$GOFILE --marshal --names --values
type Strategy string

type (
	Output struct {
		Path     string `json:"path" yaml:"path"`
		Template string `json:"template" yaml:"template"`
	}

	// T represents a element.
	T struct {
		Path     string             `json:"path" yaml:"path" validate:"required"`
		Name     string             `json:"name" yaml:"name" validate:"required"`
		Alias    string             `json:"alias" yaml:"alias"`
		Require  []string           `json:"require" yaml:"require"`
		Args     Args               `json:"args" yaml:"args" validate:"required"`
		Template string             `json:"template" yaml:"template" validate:"required"`
		Strategy Strategy           `json:"strategy" yaml:"strategy" validate:"required"`
		Output   map[string]*Output `json:"output" yaml:"output"`
	}

	// D represents a group of declaration for templates.
	D struct {
		Path        string   `json:"path" yaml:"path"`
		Name        string   `json:"name" yaml:"name"`
		Require     []string `json:"require" yaml:"require"`
		Args        Args     `json:"args" yaml:"args"`
		Templates   []*T     `json:"elements" yaml:"elements"`
		Description string   `json:"description" yaml:"description"`
		selected    string
	}

	// M represents a map of template definitions.
	M map[string]*D
)

func (e *T) Validate() error {
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
	if len(d.Templates) == 0 {
		return fmt.Errorf("elements are required")
	}

	for _, e := range d.Templates {
		if err := e.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// GetElements returns the element with the given name.
func (d *D) GetElements(name string) *T {
	for _, el := range d.Templates {
		if el.Name == name {
			return el
		}
	}

	return nil
}
