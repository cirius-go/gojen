package gojen

import (
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	pl         = pluralize.NewClient()
	titleCaser = cases.Title(language.AmericanEnglish)
)

func singular(s string) string {
	return pl.Singular(s)
}

func plural(s string) string {
	return pl.Plural(s)
}

func title(s string) string {
	return titleCaser.String(s)
}

func lower(s string) string {
	return strings.ToLower(s)
}

func upper(s string) string {
	return strings.ToUpper(s)
}

var templateFuncs = template.FuncMap{
	"singular":       singular,
	"plural":         plural,
	"title":          title,
	"lower":          lower,
	"upper":          upper,
	"snake":          strcase.ToSnake,
	"camel":          strcase.ToCamel,
	"screamingSnake": strcase.ToScreamingSnake,
	"kebab":          strcase.ToKebab,
	"screamingKebab": strcase.ToScreamingKebab,
	"lowerCamel":     strcase.ToLowerCamel,
}
