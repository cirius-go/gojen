package main

import (
	"regexp"
	"strings"
	"text/template"

	"github.com/cirius-go/gojen"
	"github.com/cirius-go/gojen/example/codegen/main/decl"
	"github.com/cirius-go/gojen/lib/pipeline"
)

type act string

var (
	Action act = ""
)

func main() {
	pc := pipeline.C()
	pc.UpdateFuncs(func(def template.FuncMap) template.FuncMap {
		def["gopkg"] = func(domain string) string {
			// special case
			if domain == "cms" {
				return "cms"
			}
			// special case
			switch snakePipeline := def["sSnake"]; snakePipeline.(type) {
			case func(string) string:
				snaked := snakePipeline.(func(string) string)(domain)
				return strings.ReplaceAll(snaked, "_", "")
			default:
				panic("sSnake is not a function")
			}
		}

		def["parseEchoSlug"] = func(path string) string {
			re := regexp.MustCompile(`:([a-zA-Z0-9-_]+)`)
			return re.ReplaceAllString(path, `{$1}`)
		}
		return def
	})
	c := gojen.C().SetPipelineConfig(pc)
	g := gojen.NewWithConfig(c)
	g.UpdateArgs(gojen.Args{"Domain": "customer", "BaseAPI": "cms"})
	g.SetDecls(decl.APIDecl)

	// crud
	s := gojen.
		NewSeq("model", "initModelFile").
		// AppendWiths([]gojen.Args{
		// 	{"Method": "get", "View": ""},
		// 	{"Method": "update"},
		// 	{"Method": "list"},
		// }, "model", "createView").
		AppendWiths([]gojen.Args{
			{"Method": "get"},
			{"Method": "list"},
			{"Method": "create"},
			{"Method": "update"},
			{"Method": "delete"},
		}, "dto", "createDto").
		Append("svc", "initIntfFile", "createIntf", "initSvcFile").
		AppendWiths([]gojen.Args{
			{"Method": "get"},
			{"Method": "list"},
			{"Method": "create"},
			{"Method": "update"},
			{"Method": "delete"},
		}, "svc", "createHandler").
		Append("api", "initIntfFile", "createIntf", "initApiFile").
		AppendWiths([]gojen.Args{
			{"Method": "get", "HTTPMethod": "get", "Slug": "/:id"},
			{"Method": "list", "HTTPMethod": "get", "Slug": ""},
			{"Method": "create", "HTTPMethod": "post", "Slug": ""},
			{"Method": "update", "HTTPMethod": "patch", "Slug": "/:id"},
			{"Method": "delete", "HTTPMethod": "delete", "Slug": "/:id"},
		}, "api", "createHandler")

	if err := g.Build(s); err != nil {
		panic(err)
	}

	if err := g.Apply(); err != nil {
		panic(err)
	}
}
