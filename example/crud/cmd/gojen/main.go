package main

import "github.com/cirius-go/gojen"

func main() {
	c := gojen.C()

	g := gojen.New(c)

	g.SetTemplate("repo", &gojen.Template{
		Path:            "internal/repo/{{ lower .Domain }}.go",
		RequiredContext: []string{"Domain"},
		TemplateString: `
  package repo

  // {{ plural .Domain | title }} represents the {{ lower .Domain }} repository.
  type {{ plural .Domain | title }} struct {
    db *gorm.DB
    *Repo[model.{{ singular .Domain | title }}]
  }

  // New{{ plural .Domain | title }} returns a new {{ plural .Domain | title }} instance.
  func New{{ plural .Domain | title }}(db *gorm.DB) *{{ plural .Domain | title }} {
    return &{{ plural .Domain | title }}{
      db: db,
      Repo: NewRepo[{{ singular .Domain | title }}](db),
    }
  }
  `,
	})

	g.SetTemplate("model", &gojen.Template{
		Path: "internal/repo/model/{{ lower .Domain }}.go",
	})
}
