package pipeline

import (
	"html/template"
	"strings"

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
	// Config contains configurations for pipeline.
	Config struct {
		irregularMap  map[string]string
		uncountable   []string
		updateFuncFns []func(def template.FuncMap) template.FuncMap
	}

	// type Pipeline manage the Pipeline functions to handle variables in the
	// template.
	Pipeline struct {
		cfg        *Config
		pl         *pluralize.Client
		titleCaser cases.Caser
		funcs      template.FuncMap
	}
)

// AddIrregularMap set the irregular map.
func (c *Config) AddIrregularMap(irregularMap map[string]string) *Config {
	for k, v := range irregularMap {
		c.irregularMap[k] = v
	}
	return c
}

// AddUncountable set the uncountable words.
func (c *Config) AddUncountable(uncountable []string) *Config {
	c.uncountable = append(c.uncountable, uncountable...)
	return c
}

// C returns a new PipelineConfig with default.
func C() *Config {
	return &Config{
		irregularMap: make(map[string]string),
		uncountable:  []string{},
	}
}

// New returns a new pipeline with default configuration.
func New() *Pipeline {
	c := C()

	return NewWithConfig(c)
}

// newPipelineWithConfig returns a new pipeline with the given configuration.
func NewWithConfig(cfg *Config) *Pipeline {
	if cfg == nil {
		panic("pipeline config is required")
	}

	p := &Pipeline{
		cfg:   cfg,
		funcs: makeDefaultFuncs(cfg),
	}

	if len(cfg.updateFuncFns) > 0 {
		for _, fn := range cfg.updateFuncFns {
			def := makeDefaultFuncs(cfg)
			newFuncs := fn(def)
			for k, v := range newFuncs {
				p.funcs[k] = v
			}
		}
	}

	return p
}

// GetFuncs returns the pipeline functions.
func (p *Pipeline) GetFuncs() template.FuncMap {
	return p.funcs
}

// UpdateFuncs updates the pipeline functions.
func (p *Config) UpdateFuncs(fn func(def template.FuncMap) template.FuncMap) {
	p.updateFuncFns = append(p.updateFuncFns, fn)
}

func makeDefaultFuncs(cfg *Config) template.FuncMap {
	pl := pluralize.NewClient()
	for k, v := range cfg.irregularMap {
		pl.AddIrregularRule(k, v)
	}
	for _, w := range cfg.uncountable {
		pl.AddUncountableRule(w)
	}
	tc := cases.Title(language.AmericanEnglish)

	return template.FuncMap{
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
	}
}
