package gojen

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
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

	cfg          *Config
	context      map[string]any
	defs         map[string]*D
	dependencies map[string][]string

	useTmpls   []string
	titleCaser cases.Caser
	pluralize  *pluralize.Client
	tmplFuncs  template.FuncMap
	argCMD     *cobra.Command
}

// New returns a new Gojen instance.
func New() *Gojen {
	c := C()

	return NewWithConfig(c)
}

// NewWithConfig returns a new Gojen instance.
func NewWithConfig(cfg *Config) *Gojen {
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
			"singular":   pl.Singular,
			"plural":     pl.Plural,
			"title":      tc.String,
			"lower":      strings.ToLower,
			"upper":      strings.ToUpper,
			"snake":      strcase.ToSnake,
			"titleSnake": strcase.ToScreamingSnake,
			"camel":      strcase.ToCamel,
			"lowerCamel": strcase.ToLowerCamel,
			"kebab":      strcase.ToKebab,
			"titleKebab": strcase.ToScreamingKebab,
		},
		argCMD: &cobra.Command{
			Use:   "gojen",
			Short: "gojen is a code generator that uses Go templates.",
		},
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

	if len(g.cfg.customPipeline) > 0 {
		g.tmplFuncs = mergeMaps(g.tmplFuncs, g.cfg.customPipeline)
	}

	if cfg.parseArgs {
		g.argCMD.Execute()
	}

	return g
}

func (g *Gojen) AddMoreCommand(cmds ...*cobra.Command) {
	if !g.cfg.parseArgs {
		fmt.Println("'parseArgs' is disabled. Skipping to add more commands.")
		return
	}

	for _, c := range cmds {
		g.argCMD.AddCommand(c)
	}
}

// LoadDir loads template definitions from a directory.
func (g *Gojen) LoadDir(dir string) error {
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
		if err := g.LoadDef(fp); err != nil {
			return err
		}

		fmt.Println("Loaded template definition from:", fp)
	}

	return nil
}

// LoadDef loads a template definition from a file.
func (g *Gojen) LoadDef(fp string) error {
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
		fmt.Printf("template '%s' already exists. Skipped to import definition from: %s\n", d.Name, fp)
		return nil
	}
	g.defs[d.Name] = d
	return nil
}

// SetPluralRules sets the plural rules for the strcase package.
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
func (g *Gojen) SetTemplate(name string, template *D) *D {
	if template == nil {
		return nil
	}

	for _, s := range template.Select {
		// only set the default strategy if it is empty.
		if s.Strategy == "" {
			s.Strategy = StrategyIgnore
		}

		if !s.Strategy.IsValid() {
			panic(fmt.Sprintf("invalid strategy: %s", s.Strategy))
		}
	}

	if len(template.Dependencies) > 0 {
		if _, exists := g.dependencies[name]; !exists {
			g.dependencies[name] = template.Dependencies
		}

		// Check for dependency cycles using DFS
		visited := map[string]bool{}
		stack := map[string]bool{}
		if g.hasCycle(name, visited, stack) {
			panic(fmt.Sprintf("cycle detected with template: %s", name))
		}
	}

	template.Name = name
	g.defs[name] = template

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

// AddContext adds a key-value pair to the context map.
func (g *Gojen) AddContext(key string, value any) {
	g.context[key] = value
}

// SetContext sets the context map.
func (g *Gojen) SetContext(ctx map[string]any) {
	g.context = ctx
}

// parseTemplate creates, executes a template, and returns the result as a string.
func (g *Gojen) parseTemplate(name string, templateString string, ctx map[string]any) (string, error) {
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
func (g *Gojen) buildTemplate(name string, d *D, seq *sequence) (string, error) {
	Redf("== Building template: %s ==\n", name)
	// Prepare the context for template 'd'
	d = d.mergeContext(g.context)

	if err := d.tryProvideCtx(d.Require); err != nil {
		return "", err
	}

	// make the file path with context.
	fp, err := g.parseTemplate(d.Path, d.Path, d.Context)
	if err != nil {
		return "", err
	}

	if len(d.Select) == 0 {
		return "", fmt.Errorf("no template to select")
	}

	var (
		parsedContent string
		strategy      Strategy
		confirm       bool
	)
	if len(d.Select) > 0 {
		selected := 0
		if seq != nil && seq.n == name {
			selected = seq.i
		} else {
			Redf("Please select one of the following template for '%s':\n", name)

			for i, s := range d.Select {
				msg := color.Bluef("Option %d: %s\n", i+1, fp)
				msg += color.Greenf("%s\n\n", s)
				fmt.Printf(msg)
			}

			fmt.Printf("Enter the option number: ")

			_, err = fmt.Scanln(&selected)
			if err != nil {
				return "", err
			}

			if selected < 1 || selected > len(d.Select) {
				return "", fmt.Errorf("invalid option selected")
			}
		}

		decl := d.Select[selected-1]
		if err = d.tryProvideCtx(decl.Require); err != nil {
			return "", err
		}
		parsedContent, err = g.parseTemplate(name, decl.Template, d.Context)
		if err != nil {
			return "", err
		}

		strategy = decl.Strategy
		confirm = decl.Confirm
	}

	if !strategy.IsValid() {
		return "", fmt.Errorf("invalid handle file strategy")
	}

	// no confirm, no modify file, just print the content.
	if g.cfg.dryRun {
		Bluef("== DRY RUN: %s ==\n%s\n", fp, d.Description)
		Greenf("%s\n", parsedContent)
		return "", nil
	}

	if confirm {
		Bluef("== Will modify content: %s ==\n", fp)
		Greenf("%s\n", parsedContent)
		Redf("Do you want to apply template '%s'? (y/N)\n", name)

		var confirm = ""
		_, err = fmt.Scanln(&confirm)
		if err != nil {
			return "", err
		}

		// if confirm is empty, it will be treated as "N" and not run the template.
		switch confirm {
		case "y", "Y", "true", "1":
		default:
			return "", nil
		}
	}

	if err := makeDirAll(fp); err != nil {
		return "", err
	}

	// special case for append.
	if strategy == StrategyAppend {
		found, err := handleOnStrategyAppend(fp, func(l string) bool {
			return strings.TrimSpace(l) == fmt.Sprintf("// +gojen:append-template=%s", name)
		}, parsedContent)
		if err != nil {
			return "", err
		}

		if found {
			fmt.Printf("modified '%s'\n", fp)
			return fp, nil
		}

		return "", fmt.Errorf("'%s' not found in the file", fmt.Sprintf("// +gojen:append-template=%s", name))
	}

	fflags := os.O_CREATE | os.O_RDWR
	switch strategy {
	case StrategyTrunc:
		fflags |= os.O_TRUNC
	case StrategyAppendAtLast:
		fflags |= os.O_APPEND
	case StrategyIgnore:
		_, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				fflags |= os.O_APPEND
			}
		} else {
			fmt.Printf("skipped to modify '%s'. This file is exist.\n", fp)
		}
	}

	file, err := os.OpenFile(fp, fflags, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.Write([]byte(parsedContent)); err != nil {
		return "", err
	}

	fmt.Printf("modified '%s'\n", fp)
	return fp, nil
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
		}

		for _, dep := range v.Dependencies {
			res[k] = append(res[k], fmt.Sprintf("'%s'", dep))
		}
	}

	return res
}

// PrintUsages prints the usages to the cli.
func (g *Gojen) PrintUsages() {
	Bluef("Usage: go run main.go [flags]\n")
	if g.cfg.parseArgs {
		flag.PrintDefaults()
	}
	if len(g.cfg.customPipeline) > 0 {
		Bluef("Current pipelines:\n")
		for k := range g.cfg.customPipeline {
			fmt.Printf("  + %s\n", k)
		}
	}

	tmplUsages := g.makeTmplUsages()
	if len(tmplUsages) == 0 {
		Bluef("No template definitions found.\n")
		return
	}
	Bluef("Template Usages:\n")
	for k, v := range tmplUsages {
		fmt.Printf("- Template: '%s'\n", k)
		for _, u := range v {
			fmt.Printf("  + %s\n", u)
		}
	}
}

// Build builds the templates.
func (g *Gojen) Build(tmplNames ...string) error {
	defs := g.defs

	useTmpls := append(g.useTmpls, tmplNames...)
	if len(useTmpls) > 0 {
		defs = filterMap(defs, useTmpls)
	}

	for name, def := range defs {
		f, err := g.buildTemplate(name, def, nil)
		if err != nil {
			return err
		}

		if f != "" {
			g.ModifiedFiles = append(g.ModifiedFiles, f)
		}
	}

	return nil
}

type sequence struct {
	n    string
	i    int
	next *sequence
	root *sequence
}

// S returns a new sequence.
func S(n string, i int) *sequence {
	s := &sequence{
		n: n,
		i: i,
	}

	s.root = s

	return s
}

func (s *sequence) S(n string, i int) *sequence {
	s.next = &sequence{
		n:    n,
		i:    i,
		root: s.root,
	}

	return s.next
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
	if err := makeDirAll(dir); err != nil {
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

	s := seq.root
	for s != nil {
		def, exists := defs[s.n]
		if !exists {
			return fmt.Errorf("template '%s' not found", s.n)
		}
		f, err := g.buildTemplate(s.n, def, s)
		if err != nil {
			return err
		}
		if f != "" {
			g.ModifiedFiles = append(g.ModifiedFiles, f)
		}
		s = s.next
	}

	return nil
}
