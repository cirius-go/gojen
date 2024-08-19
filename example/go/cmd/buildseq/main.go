package main

import (
	"github.com/cirius-go/gojen"
)

func main() {
	c := gojen.C().SetDryRun(false).ParseArgs(true)
	g := gojen.NewWithConfig(c)
	g.PrintTemplateUsage()
	g.WalkDir(".gojen/templates")

	s := gojen.
		S("service", 1). // init
		ForwardCtx("Domain").
		S("service", 2, 3). // custom method or crud
		ForwardCtx("Method").
		When(2, func(sub gojen.Sequence) {
			sub.
				M("dto", 1, 2).
				S("api", 1). // file api
				ForwardCtx("APIGroup").
				S("api", 2, 3, 4, 5, 6). // custom service method
				M("api-intf", 1, 2).
				S("api-intf-{{ sIniLowerCamel .Domain }}", 1)
		}).
		When(3, func(sub gojen.Sequence) {
			sub.
				M("dto", 1, 3).
				M("api", 1, 7).
				M("api-intf", 1, 2).
				S("api-intf-{{ sIniLowerCamel .Domain }}", 2)
		}).
		S("model").
		S("repo")

	if err := g.BuildSeqs(s); err != nil {
		panic(err)
	}
}
