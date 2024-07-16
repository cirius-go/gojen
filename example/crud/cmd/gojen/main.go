package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cirius-go/gojen"
)

var (
	projectTemplate = flag.String("templ", "", "project template name")
	ctxReq          = flag.String("context", "", "context")
)

func main() {
	flag.Parse()

	c := gojen.C()

	g := gojen.New(c)

	if len(*ctxReq) > 0 {
		mapCtx := map[string]string{}

		if err := json.Unmarshal([]byte(*ctxReq), &mapCtx); err != nil {
			panic(err)
		}

		for k, v := range mapCtx {
			g.AddContext(k, v)
		}
	}

	g.SetTemplate("dto", &gojen.D{
		Path:            "example/crud/internal/dto/{{ lower .Domain }}.go",
		RequiredContext: []string{"Domain"},
		Strategy:        gojen.StrategyIgnore,
		TemplateString: `
    package dto

    type (
      // Create{{ singular .Domain | title }}Req represents the request to create a new {{ lower .Domain }}.
      // swagger:model
      Create{{ singular .Domain | title }}Req struct {
      }

      // Create{{ singular .Domain | title }}Resp represents the response to create a new {{ lower .Domain }}.
      // swagger:model
      Create{{ singular .Domain | title }}Resp struct {
      }
    )

    type (
      // Update{{ singular .Domain | title }}Req represents the request to update a {{ lower .Domain }}.
      // swagger:model
      Update{{ singular .Domain | title }}Req struct {
        ID string ` + "`param:\"id\"`" + `
      }

      // Update{{ singular .Domain | title }}Resp represents the response to update a {{ lower .Domain }}.
      // swagger:model
      Update{{ singular .Domain | title }}Resp struct {
      }
    )

    type (
      // Get{{ singular .Domain | title }}Req represents the request to get a {{ lower .Domain }}.
      // swagger:model
      Get{{ singular .Domain | title }}Req struct {
        ID string ` + "`param:\"id\"`" + `
      }

      // Get{{ singular .Domain | title }}Resp represents the response to get a {{ lower .Domain }}.
      // swagger:model
      Get{{ singular .Domain | title }}Resp struct {
      }
    )

    type (
      // List{{ plural .Domain | title }}Req represents the request to list {{ plural .Domain }}.
      // swagger:model
      List{{ plural .Domain | title }}Req struct {
      }

      // List{{ plural .Domain | title }}Resp represents the response to list {{ plural .Domain }}.
      // swagger:model
      List{{ plural .Domain | title }}Resp struct {
      }
    )

    type (
      // Delete{{ singular .Domain | title }}Req represents the request to delete a {{ lower .Domain }}.
      // swagger:model
      Delete{{ singular .Domain | title }}Req struct {
        ID string ` + "`param:\"id\"`" + `
      }

      // Delete{{ singular .Domain | title }}Resp represents the response to delete a {{ lower .Domain }}.
      // swagger:model
      Delete{{ singular .Domain | title }}Resp struct {
      }
    )
    `,
	})

	g.SetTemplate("api", &gojen.D{
		Path:            "example/crud/internal/api/{{ lower .Domain }}.go",
		RequiredContext: []string{"Domain"},
		Strategy:        gojen.StrategyIgnore,
		TemplateString: `
    package api

    // {{ singular .Domain | title }} represents the {{ lower .Domain }} API.
    type {{ singular .Domain | title }} struct {
      svc {{ singular .Domain | title }}Svc
    }

    // New{{ singular .Domain | title }} returns a new {{ singular .Domain | title }} instance.
    func New{{ singular .Domain | title }}(svc {{ singular .Domain | title }}Svc) *{{ singular .Domain | title }} {
      return &{{ singular .Domain | title }}{
        svc: svc,
      }
    }
    `,
		Dependencies: []string{"api_svc_intf"},
	})

	g.SetTemplate("api_svc_intf", &gojen.D{
		Path:            "example/crud/internal/api/interfaces.go",
		RequiredContext: []string{"Domain"},
		Strategy:        gojen.StrategyAppend,
		TemplateString: `
    // {{ singular .Domain | title }}Svc represents the {{ lower .Domain }} service.
    type {{ singular .Domain | title }}Svc interface {
      Create(ctx context.Context, req *dto.Create{{ singular .Domain | title }}Req) (*dto.Create{{ singular .Domain | title }}Resp, error)
      Update(ctx context.Context, req *dto.Update{{ singular .Domain | title }}Req) (*dto.Update{{ singular .Domain | title }}Resp, error)
      Get(ctx context.Context, req *dto.Get{{ singular .Domain | title }}Req) (*dto.Get{{ singular .Domain | title }}Resp, error)
      List(ctx context.Context, req *dto.List{{ plural .Domain | title }}Req) (*dto.List{{ plural .Domain | title }}Resp, error)
      Delete(ctx context.Context, req *dto.Delete{{ singular .Domain | title }}Req) (*dto.Delete{{ singular .Domain | title }}Resp, error)
    }
    `,
	})

	g.SetTemplate("svc", &gojen.D{
		Path:            "example/crud/internal/service/{{ .Domain }}.go",
		RequiredContext: []string{"Domain"},
		Strategy:        gojen.StrategyIgnore,
		TemplateString: `
    package service

    // {{ singular .Domain | title }} represents the {{ lower .Domain }} service.
    type {{ singular .Domain | title }} struct {
      uow uow.UnitOfWork
    }

    // New{{ singular .Domain | title }} returns a new {{ singular .Domain | title }} instance.
    func New{{ singular .Domain | title }}(uow uow.UnitOfWork) *{{ singular .Domain | title }} {
      return &{{ singular .Domain | title }}{
        uow: uow,
      }
    }

    // Create creates a new {{ lower .Domain }}.
    func (s *{{ singular .Domain | title }}) Create(ctx context.Context, req *dto.Create{{ singular .Domain | title }}Req) (*dto.Create{{ singular .Domain | title }}Resp, error) {
      panic("not implemented")
    }

    // Update updates a {{ lower .Domain }}.
    func (s *{{ singular .Domain | title }}) Update(ctx context.Context, req *dto.Update{{ singular .Domain | title }}Req) (*dto.Update{{ singular .Domain | title }}Resp, error) {
      panic("not implemented")
    }

    // Get gets a {{ lower .Domain }}.
    func (s *{{ singular .Domain | title }}) Get(ctx context.Context, req *dto.Get{{ singular .Domain | title }}Req) (*dto.Get{{ singular .Domain | title }}Resp, error) {
      panic("not implemented")
    }

    // List lists {{ plural .Domain }}.
    func (s *{{ singular .Domain | title }}) List(ctx context.Context, req *dto.List{{ plural .Domain | title }}Req) (*dto.List{{ plural .Domain | title }}Resp, error) {
      panic("not implemented")
    }

    // Delete deletes a {{ lower .Domain }}.
    func (s *{{ singular .Domain | title }}) Delete(ctx context.Context, req *dto.Delete{{ singular .Domain | title }}Req) (*dto.Delete{{ singular .Domain | title }}Resp, error) {
      panic("not implemented")
    }
    `,
	})

	g.SetTemplate("model", &gojen.D{
		Path:            "example/crud/internal/repo/model/{{ lower .Domain }}.go",
		RequiredContext: []string{"Domain"},
		Strategy:        gojen.StrategyIgnore,
		TemplateString: `
      package model

      // {{ singular .Domain | title }} represents the {{ lower .Domain }} model.
      type {{ singular .Domain | title }} struct {
        Model
      }
    `,
	})

	if err := g.Build(); err != nil {
		panic(err)
	}

	if len(g.WrittenFiles) > 0 {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		args := []string{"-w"}
		for _, file := range g.WrittenFiles {
			p := filepath.Join(wd, file)
			fmt.Println(p)
			args = append(args, p)
		}

		formatCMD := exec.Command("goimports", args...)
		cmdResult, err := formatCMD.CombinedOutput()
		if err != nil {
			fmt.Println(string(cmdResult))
			panic(err)
		}

		fmt.Println(string(cmdResult))
	}
}
