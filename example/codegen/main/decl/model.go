package decl

import "github.com/cirius-go/gojen"

var ModelDecl = &gojen.D{
	Path:    "internal/repo/model/{{ sSnake .Domain }}.go",
	Name:    "model",
	Require: []string{"Domain"},
	Templates: []*gojen.T{
		{
			Name:     "initModelFile",
			Strategy: gojen.StrategyInit,
			Template: `package model
// {{ sIniCamel .Domain }} model.
type {{ sIniCamel .Domain }} struct {
	Model
}

// +gojen:input=dto.createDto->createView
`,
		},
		{
			Name:     "createView",
			Strategy: gojen.StrategyAppendAtPos,
			Template: `package model
// View{{ sIniCamel .Domain }} view.
type View{{ sIniCamel .Domain }} struct {
  View
}
`,
		},
	},
}
