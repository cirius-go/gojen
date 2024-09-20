package gojen

import (
	"strings"
	"text/template"

	"github.com/danielgtaylor/casing"
	"github.com/gertd/go-pluralize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ENUM(plural,singular,irregular)
//
//go:generate go-enum -f=$GOFILE --marshal --names --values
type PluralizeType string

type (
	// PipelineConfig contains configurations for pipeline.
	PipelineConfig struct {
	}

	// type pipeline manage the pipeline functions to handle variables in the
	// template.
	pipeline struct {
		cfg        *PipelineConfig
		pl         *pluralize.Client
		titleCaser cases.Caser
		funcs      template.FuncMap
	}
)

// PipelineC returns a new PipelineConfig with default.
func PipelineC() *PipelineConfig {
	return &PipelineConfig{}
}

// NewPipeline returns a new pipeline with default configuration.
func NewPipeline() *pipeline {
	c := PipelineC()

	return NewPipelineWithConfig(c)
}

// newPipelineWithConfig returns a new pipeline with the given configuration.
func NewPipelineWithConfig(cfg *PipelineConfig) *pipeline {
	if cfg == nil {
		panic("pipeline config is required")
	}

	pl := pluralize.NewClient()
	tc := cases.Title(language.AmericanEnglish)
	p := &pipeline{
		cfg:        cfg,
		pl:         pl,
		titleCaser: tc,
		funcs: template.FuncMap{
			// employees -> employee
			"singular": pl.Singular,
			// employee -> employees
			"plural": pl.Plural,

			"s": pl.Singular,
			"p": pl.Plural,

			"title":  tc.String,
			"sTitle": func(s string) string { return tc.String(pl.Singular(s)) },
			"pTitle": func(s string) string { return tc.String(pl.Plural(s)) },

			"lower":  strings.ToLower,
			"sLower": func(s string) string { return strings.ToLower(pl.Singular(s)) },
			"pLower": func(s string) string { return strings.ToLower(pl.Plural(s)) },

			"upper":  strings.ToUpper,
			"sUpper": func(s string) string { return strings.ToUpper(pl.Singular(s)) },
			"pUpper": func(s string) string { return strings.ToUpper(pl.Plural(s)) },

			"snake":  casing.Snake,
			"sSnake": func(s string) string { return casing.Snake(pl.Singular(s)) },
			"pSnake": func(s string) string { return casing.Snake(pl.Plural(s)) },

			"titleSnake":  func(s string) string { return casing.Snake(s, strings.ToUpper) },
			"sTitleSnake": func(s string) string { return casing.Snake(pl.Singular(s), strings.ToUpper) },
			"pTitleSnake": func(s string) string { return casing.Snake(pl.Plural(s), strings.ToUpper) },

			"initialism": casing.Initialism,
			"identity":   casing.Identity,

			"camel":     casing.Camel,
			"sCamel":    func(s string) string { return casing.Camel(pl.Singular(s)) },
			"pCamel":    func(s string) string { return casing.Camel(pl.Plural(s)) },
			"iniCamel":  func(s string) string { return casing.Camel(s, casing.Initialism) },
			"sIniCamel": func(s string) string { return casing.Camel(pl.Singular(s), casing.Initialism) },
			"pIniCamel": func(s string) string { return casing.Camel(pl.Plural(s), casing.Initialism) },

			"lowerCamel":     casing.LowerCamel,
			"sLowerCamel":    func(s string) string { return casing.LowerCamel(pl.Singular(s)) },
			"pLowerCamel":    func(s string) string { return casing.LowerCamel(pl.Plural(s)) },
			"iniLowerCamel":  func(s string) string { return casing.LowerCamel(s, casing.Initialism) },
			"sIniLowerCamel": func(s string) string { return casing.LowerCamel(pl.Singular(s), casing.Initialism) },
			"pIniLowerCamel": func(s string) string { return casing.LowerCamel(pl.Plural(s), casing.Initialism) },

			"kebab":       casing.Kebab,
			"sKebab":      func(s string) string { return casing.Kebab(pl.Singular(s)) },
			"pKebab":      func(s string) string { return casing.Kebab(pl.Plural(s)) },
			"titleKebab":  func(s string) string { return casing.Kebab(s, strings.ToUpper) },
			"sTitleKebab": func(s string) string { return casing.Kebab(pl.Singular(s), strings.ToUpper) },
			"pTitleKebab": func(s string) string { return casing.Kebab(pl.Plural(s), strings.ToUpper) },
		},
	}

	return p
}

// GetFuncs returns the pipeline functions.
func (p *pipeline) GetFuncs() template.FuncMap {
	return p.funcs
}

// UpdateFuncs updates the pipeline functions.
func (p *pipeline) UpdateFuncs(fn func(current template.FuncMap) template.FuncMap) {
	p.funcs = fn(p.funcs)
}
