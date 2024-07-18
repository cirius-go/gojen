package gojen

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

// Strategy is a type that represents the strategy for setting a template.
// ENUM(trunc,append,ignore)
// trunc: Truncate the destination file. (please commit this file to git before running gojen)
// append: Append at the end to the destination file.
// update: update methods of the destination struct/interface.
// override: override the destination struct/interface.
//
//go:generate go-enum -f=$GOFILE --marshal --names --values
type Strategy string

// Gojen
type Gojen struct {
	cfg          *Config
	context      map[string]any
	defs         map[string]*D
	dependencies map[string][]string
	WrittenFiles []string
}

// New returns a new Gojen instance.
func New(cfg *Config) *Gojen {
	return &Gojen{
		cfg:          cfg,
		context:      make(map[string]any),
		defs:         make(map[string]*D),
		dependencies: make(map[string][]string),
	}
}

// LoadJSON loads template config from file.
func (g *Gojen) LoadJSON(jsonPath string) error {
	f, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer f.Close()

	defs := map[string]*JSOND{}
	if err := json.NewDecoder(f).Decode(&defs); err != nil {
		return err
	}

	for k, def := range defs {
		g.SetTemplate(k, def.toD())
	}

	return nil
}

// SetTemplate sets a template.
func (g *Gojen) SetTemplate(name string, template *D) {
	if template == nil {
		return
	}

	if !template.Strategy.IsValid() {
		template.Strategy = StrategyAppend
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

	g.defs[name] = template
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

// PrintJSONDefinitions prints the JSON definitions.
func (g *Gojen) PrintJSONDefinitions() error {
	mapJSOND := map[string]*JSOND{}
	for k, v := range g.defs {
		mapJSOND[k] = v.toJSOND()
	}
	b, err := json.MarshalIndent(mapJSOND, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

// Build builds the templates.
func (g *Gojen) Build(tmplNames ...string) error {
	built := []string{}

	defs := g.defs
	if len(tmplNames) > 0 {
		defs = map[string]*D{}
		for _, name := range tmplNames {
			if tmpl, exists := g.defs[name]; exists {
				defs[name] = tmpl
			}
		}
	}

	for name, def := range defs {
		f, err := g.buildTemplate(name, def)
		if err != nil {
			return err
		}

		if f != "" {
			built = append(built, f)
		}
	}

	g.WrittenFiles = built
	return nil
}

// makeTemplate creates a new template.
func (g *Gojen) makeTemplate(name string, templateString string) (*template.Template, error) {
	t := template.
		New(name).
		Funcs(template.FuncMap{
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
		}).
		Option("missingkey=zero")
	if _, err := t.Parse(templateString); err != nil {
		return nil, err
	}

	return t, nil
}

// execTemplate executes a template.
func (g *Gojen) execTemplate(t *template.Template, writer io.Writer, ctx map[string]any) error {
	return t.Execute(writer, &ctx)
}

func (g *Gojen) execTemplateToString(t *template.Template, ctx map[string]any) (string, error) {
	writer := &strings.Builder{}
	if err := t.Execute(writer, &ctx); err != nil {
		return "", err
	}

	return writer.String(), nil
}

// buildTemplate builds a template.
func (g *Gojen) buildTemplate(name string, d *D) (string, error) {
	myCtx := make(map[string]any)

	// Merge g.context and d.Context into myCtx
	for _, contextMap := range []map[string]any{g.context, d.Context} {
		for k, v := range contextMap {
			myCtx[k] = v
		}
	}

	// Check for required context keys
	for _, key := range d.RequiredContext {
		if _, exists := myCtx[key]; !exists {
			return "", fmt.Errorf("missing required context key: %s", key)
		}
	}

	path, err := g.makeTemplate(d.Path, d.Path)
	if err != nil {
		return "", err
	}

	pathStr, err := g.execTemplateToString(path, myCtx)
	if err != nil {
		return "", err
	}

	content, err := g.makeTemplate(name, d.TemplateString)
	if err != nil {
		return "", err
	}

	if g.cfg.dryRun {
		fmt.Println("== DRY RUN ==")
		fmt.Println(pathStr)
		contentStr, err := g.execTemplateToString(content, myCtx)
		if err != nil {
			return "", err
		}

		fmt.Println(contentStr)
		return "", nil
	}

	if err := makeDirAll(pathStr); err != nil {
		return "", err
	}

	// 0777 is the permission
	flag := os.O_CREATE | os.O_WRONLY
	switch d.Strategy {
	case StrategyTrunc:
		flag |= os.O_TRUNC
	case StrategyAppend:
		flag |= os.O_APPEND
	case StrategyIgnore:
		_, err := os.Stat(pathStr)
		if err != nil {
			if os.IsNotExist(err) {
				flag |= os.O_APPEND
			}
		} else {
			fmt.Printf("skipped %s\n", pathStr)
		}
	}

	file, err := os.OpenFile(pathStr, flag, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := g.execTemplate(content, file, myCtx); err != nil {
		return "", err
	}

	fmt.Printf("generated %s\n", pathStr)
	return pathStr, nil
}
