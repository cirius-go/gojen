package main

import (
	"github.com/cirius-go/gojen"
)

func main() {
	c := gojen.C().SetDryRun(false).ParseArgs(true)
	g := gojen.NewWithConfig(c)
	g.PrintTemplateUsage()
	g.WalkDir("example/go/assets/templates")

	seq := gojen.
		S("service", 1). // init
		ForwardCtx("Domain").
		S("service", 2, 3). // custom method or crud
		ForwardCtx("Method").
		When(2, func(sub gojen.Sequence) {
			sub.S("dto", 1).S("dto", 2)
		}).
		When(3, func(sub gojen.Sequence) {
			sub.S("dto", 1).S("dto", 3)
		})

	if err := g.BuildSeqs(seq); err != nil {
		panic(err)
	}
}
