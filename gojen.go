package gojen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/danielgtaylor/casing"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	"github.com/cirius-go/gojen/color"
)

// Strategy is a type that represents the strategy for setting a template.
// ENUM(trunc,append,append_at_last,ignore)
// trunc: Truncate the destination file. (please commit this file to git before running gojen)
// append: Append at the end to the destination file.
// ignore: Ignore the destination file.
//
//go:generate go-enum -f=$GOFILE --marshal --names --values
type Strategy string

// ENUM(plural,singular,irregular)
//
//go:generate go-enum -f=$GOFILE --marshal --names --values
type PluralizeType string

// Gojen is a struct that holds the Gojen instance.
type Gojen struct {
	ModifiedFiles []string

	cfg          *config
	context      map[string]any
	defs         M
	dependencies map[string][]string

	useTmpls   []string
	titleCaser cases.Caser
	pluralize  *pluralize.Client
	tmplFuncs  template.FuncMap
	argCMD     *cobra.Command

	usage *strings.Builder
}

// New returns a new Gojen instance.
func New() *Gojen {
	c := C()

	return NewWithConfig(c)
}

// NewWithConfig returns a new Gojen instance.
func NewWithConfig(cfg *config) *Gojen {
	if cfg == nil {
		panic("config is required")
	}
	pl := pluralize.NewClient()
	tc := cases.Title(language.AmericanEnglish)

	// special case...
	pl.AddIrregularRule("staff", "Staffs")
	pl.AddIrregularRule("Staff", "Staffs")

	g := &Gojen{
		cfg:          cfg,
		context:      make(map[string]any),
		defs:         make(map[string]*D),
		dependencies: make(map[string][]string),
		titleCaser:   tc,
		pluralize:    pl,
		tmplFuncs: template.FuncMap{
			"singular": pl.Singular,
			"plural":   pl.Plural,

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
		usage: &strings.Builder{},
	}

	if len(g.cfg.customPipeline) > 0 {
		g.tmplFuncs = mergeMaps(g.tmplFuncs, g.cfg.customPipeline)
	}

	if cfg.parseArgs {
		g.argCMD = &cobra.Command{
			Use:   "gojen",
			Short: "gojen is a code generator that uses Go templates.",
		}

		g.argCMD.AddCommand(&cobra.Command{
			Use:   "ctx",
			Short: "context to be used in the template",
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, arg := range args {
					c := map[string]any{}
					if err := json.Unmarshal([]byte(arg), &c); err != nil {
						return err
					}

					g.context = mergeMaps(g.context, c)
				}

				return nil
			},
		})

		g.argCMD.Execute()
	}

	return g
}

// AddCommand wraps the cobra command.
//
// if parseArgs config is disabled, it will not add the command.
func (g *Gojen) AddCommand(cmds ...*cobra.Command) {
	if !g.cfg.parseArgs {
		fmt.Println("'parseArgs' is disabled. Skipping to add more commands.")
		return
	}

	for _, c := range cmds {
		g.argCMD.AddCommand(c)
	}
}

// WalkDir loads template definitions from a directory.
func (g *Gojen) WalkDir(dir string) error {
	if dir == "" {
		return nil
	}

	stat, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("'%s' is not a directory", dir)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		fp := filepath.Join(dir, file.Name())
		if err := g.LoadD(fp); err != nil {
			return err
		}
	}

	return nil
}

// LoadD loads a template definition from a file.
func (g *Gojen) LoadD(fp string) error {
	file, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer file.Close()
	ext := filepath.Ext(fp)

	var decoder interface{ Decode(v any) error }
	switch ext {
	case ".yaml", ".yml":
		decoder = yaml.NewDecoder(file)
	case ".json":
		decoder = json.NewDecoder(file)
	}

	d := &D{}
	if err := decoder.Decode(&d); err != nil {
		return err
	}
	if d.Name == "" {
		return fmt.Errorf("name of template is required")
	}
	if g.defs[d.Name] != nil {
		Redf("Template '%s' already exists. Skipped to import: '%s'\n", d.Name, fp)
		return nil
	}
	g.defs[d.Name] = d
	Bluef("Loaded template definition from: '%s'\n", fp)
	return nil
}

// SetPluralRules sets the plural rules for the casing package.
func (g *Gojen) SetPluralRules(rt PluralizeType, rules map[string]string) {
	for k, v := range rules {
		switch rt {
		case PluralizeTypePlural:
			g.pluralize.AddPluralRule(k, v)
		case PluralizeTypeSingular:
			g.pluralize.AddSingularRule(k, v)
		case PluralizeTypeIrregular:
			g.pluralize.AddIrregularRule(k, v)
		default:
			panic("invalid PluralizeType")
		}
	}
}

// SetTemplate sets a template.
func (g *Gojen) SetTemplate(template *D) *D {
	// perform DFS to detect cycles in the dependency graph.
	name := template.Name
	if name == "" {
		panic("name of template is required")
	}
	if len(template.Dependencies) > 0 {
		if _, exists := g.dependencies[name]; !exists {
			g.dependencies[name] = template.Dependencies
		}
		visited := map[string]bool{}
		stack := map[string]bool{}
		if g.hasCycle(name, visited, stack) {
			panic(fmt.Sprintf("cycle detected with template: %s", name))
		}
	}

	g.defs.Store(template)
	return template
}

// hasCycle performs DFS to detect cycles in the dependency graph.
func (g *Gojen) hasCycle(current string, visited map[string]bool, stack map[string]bool) bool {
	if stack[current] {
		return true
	}
	if visited[current] {
		return false
	}

	visited[current] = true
	stack[current] = true

	for _, dep := range g.dependencies[current] {
		if g.hasCycle(dep, visited, stack) {
			return true
		}
	}

	stack[current] = false
	return false
}

// AddContext adds a key-value pair to the existing context map.
func (g *Gojen) AddContext(key string, value any) {
	g.context[key] = value
}

// SetContext overrites the context map.
func (g *Gojen) SetContext(ctx map[string]any) {
	// better than nil
	if ctx == nil {
		ctx = make(map[string]any)
	}
	g.context = ctx
}

// parseTmpl creates, executes a template and returns the result as a string.
func (g *Gojen) parseTmpl(name string, templateString string, ctx map[string]any) (string, error) {
	t, err := template.New(name).Funcs(g.tmplFuncs).Option("missingkey=zero").Parse(templateString)
	if err != nil {
		return "", err
	}

	writer := strings.Builder{}
	if err := t.Execute(&writer, &ctx); err != nil {
		return "", err
	}

	return writer.String(), nil
}

// buildTemplate builds a template.
// It returns the path of the modified file and an error if any.
func (g *Gojen) buildTemplate(name string, d *D, seq *sequence) (string, *D, error) {
	Redf("== Building template: %s ==\n", name)

	// Prepare the context for template 'd'
	if seq != nil {
		d = d.clone() // clone the template to avoid modifying the original context.
	}
	// merge the global context with the template context.
	d = d.mergeContext(g.context)

	// ask for missing context values
	if err := d.provideRequireCtx(d.Required); err != nil {
		return "", nil, err
	}

	// Generate file path with the context
	fp, err := g.parseTmpl(d.Path, d.Path, d.Context)
	if err != nil {
		return "", nil, err
	}

	e, err := d.performSelect(fp, seq)
	if err != nil {
		return "", nil, err
	}

	if e == nil {
		return "", d, nil
	}

	// Parse the template from decl with the context
	parsedContent, err := g.parseTmpl(name, e.Template, d.Context)
	if err != nil {
		return "", nil, err
	}

	// If dry run is enabled, print the parsed content and do nothing to the file.
	if g.cfg.dryRun {
		Bluef("== DRY RUN: %s ==\n%s\n", fp, d.Description)
		Greenf("%s\n", parsedContent)
		return "", d, nil
	}

	// Confirm if needed
	if e.Confirm {
		Bluef("== Will modify content: %s ==\n", fp)
		Greenf("%s\n", parsedContent)
		Redf("Do you want to apply template '%s'? (y/N)\n", name)

		var confirm string
		if _, err = fmt.Scanln(&confirm); err != nil || !isConfirmed(confirm) {
			return "", d, nil
		}
	}

	// ensure the directory exists
	if err := makeDirAll(fp); err != nil {
		return "", nil, err
	}

	if e.Strategy == StrategyAppend {
		parsedName, err := g.parseTmpl(name, name, d.Context)
		if err != nil {
			return "", nil, err
		}
		found, err := handleOnStrategyAppend(fp, func(l string) bool {
			return strings.TrimSpace(l) == fmt.Sprintf("// +gojen:append-template=%s", parsedName)
		}, parsedContent)
		if err != nil {
			return "", nil, err
		}
		if found {
			fmt.Printf("modified '%s'\n", fp)
			return fp, d, nil
		}
		return "", nil, fmt.Errorf("'%s' not found in the file", fmt.Sprintf("// +gojen:append-template=%s", parsedName))
	}

	fflags := getFileFlags(fp, e.Strategy)
	if fflags == 0 {
		return "", d, nil
	}
	file, err := os.OpenFile(fp, fflags, 0644)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	if _, err := file.Write([]byte(parsedContent)); err != nil {
		return "", nil, err
	}

	fmt.Printf("modified '%s'\n", fp)
	return fp, d, nil
}

func (g *Gojen) makeTmplUsages() map[string][]string {
	res := map[string][]string{}

	for k, v := range g.defs {
		if _, ok := res[k]; !ok {
			res[k] = []string{}
		}

		res[k] = append(res[k], v.Description)
		if len(v.Dependencies) > 0 {
			res[k] = append(res[k], "Dependencies: ")

			for _, dep := range v.Dependencies {
				res[k] = append(res[k], fmt.Sprintf("'%s'", dep))
			}
		}

		if len(v.Select) > 0 {
			res[k] = append(res[k], "Select: ")
			for i, s := range v.Select {
				res[k] = append(res[k], fmt.Sprintf("  Option %d uses '%s' strategy and requires context: %s", i+1, s.Strategy, strings.Join(s.Require, ", ")))
			}
		}
	}

	return res
}

func (g *Gojen) PrintTemplateUsage() {
	if len(g.cfg.customPipeline) > 0 {
		g.usage.WriteString(color.Bluef("Current pipelines:\n"))
		for k := range g.cfg.customPipeline {
			g.usage.WriteString(fmt.Sprintf("  + %s\n", k))
		}
	}

	tmplUsages := g.makeTmplUsages()
	if len(tmplUsages) == 0 {
		g.usage.WriteString(color.Bluef("No template definitions found.\n"))
		return
	}
	g.usage.WriteString(color.Bluef("Template Usages:\n"))
	for k, v := range tmplUsages {
		g.usage.WriteString(fmt.Sprintf("- Template: '%s'\n", k))
		for _, u := range v {
			g.usage.WriteString(fmt.Sprintf("  + %s\n", u))
		}
	}

	fmt.Println(g.usage.String())
}

// Build builds the templates.
func (g *Gojen) Build(tmplNames ...string) error {
	defs := g.defs

	useTmpls := append(g.useTmpls, tmplNames...)
	if len(useTmpls) > 0 {
		defs = filterMap(defs, useTmpls)
	}

	for name, def := range defs {
		f, _, err := g.buildTemplate(name, def, nil)
		if err != nil {
			return err
		}

		if f != "" {
			g.ModifiedFiles = append(g.ModifiedFiles, f)
		}
	}

	return nil
}

// WriteDir writes definitions to a directory.
func (g *Gojen) WriteDir(dir string, ext string) error {
	if ext != ".json" && ext != ".yaml" && ext != ".yml" {
		return fmt.Errorf("invalid file extension: %s", ext)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir = filepath.Join(wd, dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for k, v := range g.defs {
		if err := func() error {
			fp := filepath.Join(dir, fmt.Sprintf("%s%s", k, ext))
			file, err := os.Create(fp)
			if err != nil {
				return err
			}
			defer file.Close()

			var encoder interface{ Encode(v any) error }
			switch ext {
			case ".yaml", ".yml":
				encoder = yaml.NewEncoder(file)
			default:
				encoder = json.NewEncoder(file)
			}

			if err := encoder.Encode(v); err != nil {
				return err
			}

			fmt.Printf("written template definition to: %s\n", fp)
			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}

// BuildSeqs builds the templates in sequence.
func (g *Gojen) BuildSeqs(seq *sequence, tmplNames ...string) error {
	defs := g.defs

	useTmpls := append(g.useTmpls, tmplNames...)
	if len(useTmpls) > 0 {
		defs = filterMap(defs, useTmpls)
	}

	var (
		s          = seq.root
		forwardCtx = map[string]any{}
	)
	for s != nil {
		// select template and build.
		d, exists := defs[s.n]
		if !exists {
			return fmt.Errorf("template '%s' not found", s.n)
		}

		if len(forwardCtx) > 0 {
			d.Context = mergeMaps(d.Context, forwardCtx)
		}

		f, builtD, err := g.buildTemplate(s.n, d, s)
		if err != nil {
			return err
		}
		if f != "" {
			g.ModifiedFiles = append(g.ModifiedFiles, f)
		}

		// forward ctx
		if s.forwardCtx != nil {
			if len(*s.forwardCtx) > 0 {
				for _, fc := range *s.forwardCtx {
					val, ok := builtD.Context[fc]
					if !ok {
						continue
					}

					forwardCtx[fc] = val
				}
			} else {
				forwardCtx = mergeMaps(forwardCtx, builtD.Context)
			}
		}

		if len(s.when) > 0 {
			branch, ok := s.when[builtD.selectedTmpl]
			if !ok {
				s = s.next
			} else {
				last := branch.last()
				last.next = s.next
				s = branch.next
			}
		} else {
			s = s.next
		}
	}

	return nil
}
