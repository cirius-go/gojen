package main

import (
	"encoding/json"
	"flag"

	"github.com/cirius-go/gojen"
)

var (
	tmplFile = flag.String("tmpl", "", "JSON template file path")
	ctxReq   = flag.String("ctx", "", "context")
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

	if len(*tmplFile) != 0 {
		if err := g.LoadJSON(*tmplFile); err != nil {
			panic(err)
		}
	}

	if err := g.Build(); err != nil {
		panic(err)
	}
}
