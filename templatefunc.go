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

// ENUM(plural,singular,irregular)
//
//go:generate go-enum -f=$GOFILE --marshal --names --values
type PluralizeType string

// SetPluralRules sets the plural rules for the strcase package.
func SetPluralRules(rt PluralizeType, rules map[string]string) {
	for k, v := range rules {
		switch rt {
		case PluralizeTypePlural:
			pl.AddPluralRule(k, v)
		case PluralizeTypeSingular:
			pl.AddSingularRule(k, v)
		case PluralizeTypeIrregular:
			pl.AddIrregularRule(k, v)
		default:
			panic("invalid PluralizeType")
		}
	}
}

// SetAcronyms sets the acronyms for the strcase package.
func SetAcronyms(acs map[string]string) {
	for k, v := range acs {
		strcase.ConfigureAcronym(k, v)
	}
}

func Singular(s string) string {
	return pl.Singular(s)
}

func Plural(s string) string {
	return pl.Plural(s)
}

func Title(s string) string {
	return titleCaser.String(s)
}

func Lower(s string) string {
	return strings.ToLower(s)
}

func Upper(s string) string {
	return strings.ToUpper(s)
}

// strcase functions
var (
	Snake          = strcase.ToSnake
	Camel          = strcase.ToCamel
	ScreamingSnake = strcase.ToScreamingSnake
	Kebab          = strcase.ToKebab
	ScreamingKebab = strcase.ToScreamingKebab
	LowerCamel     = strcase.ToLowerCamel
)

var templateFuncs = template.FuncMap{
	"singular":       Singular,
	"plural":         Plural,
	"title":          Title,
	"lower":          Lower,
	"upper":          Upper,
	"snake":          Snake,
	"camel":          Camel,
	"screamingSnake": ScreamingSnake,
	"kebab":          Kebab,
	"screamingKebab": ScreamingKebab,
	"lowerCamel":     LowerCamel,
}
