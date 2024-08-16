package main

import "github.com/cirius-go/gojen"

func main() {
	c := gojen.C().SetDryRun(false).ParseArgs(true).SetDebug(true)
	g := gojen.NewWithConfig(c)
	g.PrintUsages()

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

// {{ camel .Domain | singular }} service.
type {{ camel .Domain | singular }} struct {
}

// Please select option 2,3 to append service method for service '{{ camel .Domain | singular }}'.
// +gojen:append-template=service`,
			},
			{
				Strategy: gojen.StrategyAppend,
				Confirm:  true,
				Require:  []string{"Domain", "Method"},
				Template: `
// {{ camel .Method | singular }} method of '{{ camel .Domain | singular }}' service.
func (s *{{ camel .Domain | singular }}) {{ camel .Method | singular }}(ctx context.Context, req *dto.{{ camel .Method | singular }}{{ camel .Domain | singular }}Req) (*dto.{{ camel .Method | singular }}{{ camel .Domain | singular }}Res, error) {
  panic("not implemented")
}`,
			},
			{
				Strategy: gojen.StrategyAppend,
				Confirm:  true,
				Require:  []string{"Domain"},
				Template: `// Create new '{{ camel .Domain | singular }}'.
func (s *{{ camel .Domain | singular }}) Create(ctx context.Context, req *dto.Create{{ camel .Domain | singular }}Req) (*dto.Create{{ camel .Domain | singular }}Res, error) {
  panic("not implemented")
}

// List all '{{ camel .Domain | plural }}'.
func (s *{{ camel .Domain | singular }}) List(ctx context.Context, req *dto.List{{ camel .Domain | plural }}Req) (*dto.List{{ camel .Domain | plural }}Res, error) {
  panic("not implemented")
}

// Get one of '{{ camel .Domain | plural }}'.
func (s *{{ camel .Domain | singular }}) Get(ctx context.Context, req *dto.Get{{ camel .Domain | singular }}Req) (*dto.Get{{ camel .Domain | singular }}Res, error) {
  panic("not implemented")
}

// Update one of '{{ camel .Domain | plural }}'.
func (s *{{ camel .Domain | singular }}) Update(ctx context.Context, req *dto.Update{{ camel .Domain | singular }}Req) (*dto.Update{{ camel .Domain | singular }}Res, error) {
  panic("not implemented")
}

// Delete one of '{{ camel .Domain | plural }}'.
func (s *{{ camel .Domain | singular }}) Delete(ctx context.Context, req *dto.Delete{{ camel .Domain | singular }}Req) (*dto.Delete{{ camel .Domain | singular }}Res, error) {
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
