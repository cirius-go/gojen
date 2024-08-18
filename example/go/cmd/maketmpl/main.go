package main

import "github.com/cirius-go/gojen"

func main() {
	c := gojen.C().SetDryRun(false).ParseArgs(true)
	g := gojen.NewWithConfig(c)
	g.PrintTemplateUsage()

	g.SetTemplate(&gojen.D{
		Path:     "example/go/internal/dto/{{ sLower .Domain }}.go",
		Name:     "dto",
		Required: []string{"Domain"},
		Select: []*gojen.E{
			{
				Strategy: gojen.StrategyIgnore,
				Require:  []string{"Domain"},
				Confirm:  false,
				Template: `package dto

// DO NOT REMOVE THIS COMMENT
// +gojen:append-template=dto`,
			},
			{
				Strategy: gojen.StrategyAppend,
				Require:  []string{"Domain", "Method"},
				Confirm:  true,
				Template: `// {{ sIniCamel .Method }}{{ sIniCamel .Domain }}Req represents the request body of '{{ sIniCamel .Domain }}.{{ sIniCamel .Method }}'.
type {{ sIniCamel .Method }}{{ sIniCamel .Domain }}Req struct {
}

// {{ sIniCamel .Method }}{{ sIniCamel .Domain }}Res represents the response body of '{{ sIniCamel .Domain }}.{{ sIniCamel .Method }}'.
type {{ sIniCamel .Method }}{{ sIniCamel .Domain }}Res struct {
}`,
			},
			{
				Strategy: gojen.StrategyAppend,
				Require:  []string{"Domain"},
				Confirm:  true,
				Template: `// Create{{ sIniCamel .Domain }}Req represents the request body of '{{ sIniCamel .Domain }}.Create'.
type Create{{ sIniCamel .Domain }}Req struct {
}

// Create{{ sIniCamel .Domain }}Res represents the response body of '{{ sIniCamel .Domain }}.Create'.
type Create{{ sIniCamel .Domain }}Res struct {
}

// List{{ pIniCamel .Domain }}Req represents the request body of '{{ sIniCamel .Domain }}.List'.
type List{{ pIniCamel .Domain }}Req struct {
}

// List{{ pIniCamel .Domain }}Res represents the response body of '{{ sIniCamel .Domain }}.List'.
type List{{ pIniCamel .Domain }}Res struct {
}

// Get{{ sIniCamel .Domain }}Req represents the request body of '{{ sIniCamel .Domain }}.Get'.
type Get{{ sIniCamel .Domain }}Req struct {
}

// Get{{ sIniCamel .Domain }}Res represents the response body of '{{ sIniCamel .Domain }}.Get'.
type Get{{ sIniCamel .Domain }}Res struct {
}

// Update{{ sIniCamel .Domain }}Req represents the request body of '{{ sIniCamel .Domain }}.Update'.
type Update{{ sIniCamel .Domain }}Req struct {
}

// Update{{ sIniCamel .Domain }}Res represents the response body of '{{ sIniCamel .Domain }}.Update'.
type Update{{ sIniCamel .Domain }}Res struct {
}

// Delete{{ sIniCamel .Domain }}Req represents the request body of '{{ sIniCamel .Domain }}.Delete'.
type Delete{{ sIniCamel .Domain }}Req struct {
}

// Delete{{ sIniCamel .Domain }}Res represents the response body of '{{ sIniCamel .Domain }}.Delete'.
type Delete{{ sIniCamel .Domain }}Res struct {
}`,
			},
		},
	})

	g.SetTemplate(&gojen.D{
		Path:        "example/go/internal/service/{{ sLower .Domain }}.go",
		Name:        "service",
		Required:    []string{"Domain"},
		Description: "Generate a service file or append a service method.",
		Select: []*gojen.E{
			{
				Strategy: gojen.StrategyIgnore,
				Require:  []string{"Domain"},
				Confirm:  false,
				Template: `package service

// {{ sIniCamel .Domain }} service.
type {{ sIniCamel .Domain }} struct {
}

// DO NOT REMOVE THIS COMMENT
// +gojen:append-template=service`,
			},
			{
				Strategy: gojen.StrategyAppend,
				Confirm:  true,
				Require:  []string{"Domain", "Method"},
				Template: `
// {{ sIniCamel .Method }} method of '{{ sIniCamel .Domain }}' service.
func (s *{{ sIniCamel .Domain }}) {{ sIniCamel .Method }}(ctx context.Context, req *dto.{{ sIniCamel .Method }}{{ sIniCamel .Domain }}Req) (*dto.{{ sIniCamel .Method }}{{ sIniCamel .Domain }}Res, error) {
  panic("not implemented")
}`,
			},
			{
				Strategy: gojen.StrategyAppend,
				Confirm:  true,
				Require:  []string{"Domain"},
				Template: `// Create new '{{ sIniCamel .Domain }}'.
func (s *{{ sIniCamel .Domain }}) Create(ctx context.Context, req *dto.Create{{ sIniCamel .Domain }}Req) (*dto.Create{{ sIniCamel .Domain }}Res, error) {
  panic("not implemented")
}

// List all '{{ pIniCamel .Domain }}'.
func (s *{{ sIniCamel .Domain }}) List(ctx context.Context, req *dto.List{{ pIniCamel .Domain }}Req) (*dto.List{{ pIniCamel .Domain }}Res, error) {
  panic("not implemented")
}

// Get one of '{{ pIniCamel .Domain }}'.
func (s *{{ sIniCamel .Domain }}) Get(ctx context.Context, req *dto.Get{{ sIniCamel .Domain }}Req) (*dto.Get{{ sIniCamel .Domain }}Res, error) {
  panic("not implemented")
}

// Update one of '{{ pIniCamel .Domain }}'.
func (s *{{ sIniCamel .Domain }}) Update(ctx context.Context, req *dto.Update{{ sIniCamel .Domain }}Req) (*dto.Update{{ sIniCamel .Domain }}Res, error) {
  panic("not implemented")
}

// Delete one of '{{ pIniCamel .Domain }}'.
func (s *{{ sIniCamel .Domain }}) Delete(ctx context.Context, req *dto.Delete{{ sIniCamel .Domain }}Req) (*dto.Delete{{ sIniCamel .Domain }}Res, error) {
  panic("not implemented")
}`,
			},
		},
	})

	if err := g.WriteDir("example/go/assets/templates", ".yaml"); err != nil {
		panic(err)
	}

	if err := g.WriteDir("example/go/assets/templates", ".json"); err != nil {
		panic(err)
	}
}
