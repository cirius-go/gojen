package gojen

import (
	"fmt"

	"github.com/cirius-go/gojen/color"
)

// D represents a template definition.
type D struct {
	Path         string
	Require      []string
	Name         string
	Context      map[string]any
	Select       []*DItem
	Strategy     Strategy
	Dependencies []string
	Description  string
	Confirm      bool
}

// With sets the required context keys and the template string.
func (d *D) With(requires []string, template string) *D {
	d.Select = append(d.Select, &DItem{
		Template: template,
		Require:  requires,
	})

	return d
}

// DItem represents a template item.
type DItem struct {
	Template string
	Require  []string
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
