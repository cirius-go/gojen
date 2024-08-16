package main

import "github.com/cirius-go/gojen"

func main() {
	c := gojen.C().SetDryRun(false).ParseArgs(true).SetDebug(true)
	g := gojen.NewWithConfig(c)

	g.LoadDir("example/go/assets/templates")
	g.PrintTemplateUsage()

	// s := gojen.S("service", 1).S("service", 3)
	// if err := g.BuildSeqs(s); err != nil {
	// 	panic(err)
	// }
}
