package main

import "github.com/cirius-go/gojen"

func main() {
	c := gojen.C().SetDryRun(false).ParseArgs(true).SetDebug(true)
	g := gojen.NewWithConfig(c)
	g.PrintTemplateUsage()

	g.SetTemplate("service", &gojen.D{
		Path:        "example/go/internal/service/{{ lower .Domain | singular }}.go",
		Require:     []string{"Domain"},
		Description: "Generate a service file or append a service method.",
		Select: []*gojen.DItem{
			{
				Strategy: gojen.StrategyIgnore,
				Require:  []string{"Domain"},
				Confirm:  false,
				Template: `package service

// {{ siniCamel .Domain }} service.
type {{ siniCamel .Domain }} struct {
}

// Please select option 2,3 to append service method for service '{{ siniCamel .Domain }}'.
// +gojen:append-template=service`,
			},
			{
				Strategy: gojen.StrategyAppend,
				Confirm:  true,
				Require:  []string{"Domain", "Method"},
				Template: `
// {{ siniCamel .Method }} method of '{{ siniCamel .Domain }}' service.
func (s *{{ siniCamel .Domain }}) {{ siniCamel .Method }}(ctx context.Context, req *dto.{{ siniCamel .Method }}{{ siniCamel .Domain }}Req) (*dto.{{ siniCamel .Method }}{{ siniCamel .Domain }}Res, error) {
  panic("not implemented")
}`,
			},
			{
				Strategy: gojen.StrategyAppend,
				Confirm:  true,
				Require:  []string{"Domain"},
				Template: `// Create new '{{ siniCamel .Domain }}'.
func (s *{{ siniCamel .Domain }}) Create(ctx context.Context, req *dto.Create{{ siniCamel .Domain }}Req) (*dto.Create{{ siniCamel .Domain }}Res, error) {
  panic("not implemented")
}

// List all '{{ piniCamel .Domain }}'.
func (s *{{ siniCamel .Domain }}) List(ctx context.Context, req *dto.List{{ piniCamel .Domain }}Req) (*dto.List{{ piniCamel .Domain }}Res, error) {
  panic("not implemented")
}

// Get one of '{{ piniCamel .Domain }}'.
func (s *{{ siniCamel .Domain }}) Get(ctx context.Context, req *dto.Get{{ siniCamel .Domain }}Req) (*dto.Get{{ siniCamel .Domain }}Res, error) {
  panic("not implemented")
}

// Update one of '{{ piniCamel .Domain }}'.
func (s *{{ siniCamel .Domain }}) Update(ctx context.Context, req *dto.Update{{ siniCamel .Domain }}Req) (*dto.Update{{ siniCamel .Domain }}Res, error) {
  panic("not implemented")
}

// Delete one of '{{ piniCamel .Domain }}'.
func (s *{{ siniCamel .Domain }}) Delete(ctx context.Context, req *dto.Delete{{ siniCamel .Domain }}Req) (*dto.Delete{{ siniCamel .Domain }}Res, error) {
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
