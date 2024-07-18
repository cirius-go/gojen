package gojen

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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
	cfg           *Config
	context       map[string]any
	defs          map[string]*D
	dependencies  map[string][]string
	ModifiedFiles []string
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
	defs := map[string]*JSOND{}
	if err := loadJSON(jsonPath, &defs); err != nil {
		return err
	}

	for k, def := range defs {
		g.SetTemplate(k, def.toD())
	}

	return nil
}

// LoadJSONInDir loads template config from files in dir.
func (g *Gojen) LoadJSONInDir(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fn := file.Name()
		ext := filepath.Ext(fn)
		if ext != ".json" {
			continue
		}
		if err := g.LoadJSON(fn); err != nil {
			return err
		}
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

// makeTemplate creates a new template.
func (g *Gojen) makeTemplate(name string, templateString string) (*template.Template, error) {
	t := template.
		New(name).
		Funcs(templateFuncs).
		Option("missingkey=zero")
	if _, err := t.Parse(templateString); err != nil {
		return nil, err
	}

	return t, nil
}

// exec executes a template.
func (g *Gojen) exec(t *template.Template, writer io.Writer, ctx map[string]any) error {
	return t.Execute(writer, &ctx)
}

// execToStr executes a template and returns the result as a string.
func (g *Gojen) execToStr(t *template.Template, ctx map[string]any) (string, error) {
	writer := &strings.Builder{}
	if err := t.Execute(writer, &ctx); err != nil {
		return "", err
	}

	return writer.String(), nil
}

// makeAndExec creates and executes a template.
func (g *Gojen) makeAndExec(name string, templateString string, writer io.Writer, ctx map[string]any) error {
	t, err := g.makeTemplate(name, templateString)
	if err != nil {
		return err
	}
	return g.exec(t, writer, ctx)
}

// makeAndExecToStr creates, executes a template, and returns the result as a string.
func (g *Gojen) makeAndExecToStr(name string, templateString string, ctx map[string]any) (string, error) {
	t, err := g.makeTemplate(name, templateString)
	if err != nil {
		return "", err
	}
	return g.execToStr(t, ctx)
}

// buildTemplate builds a template.
// It returns the path of the modified file and an error if any.
func (g *Gojen) buildTemplate(name string, d *D) (string, error) {
	// Prepare the context for template 'd'
	d.mergeGlobalCtx(g.context)

	// Check all required context keys are satisfied
	if satisfied, notSatisfied := d.isCtxSatisfied(); !satisfied {
		return "", fmt.Errorf("required context keys not satisfied: %v", notSatisfied)
	}

	// make the file path with context.
	fp, err := g.makeAndExecToStr(d.Path, d.Path, d.Context)
	if err != nil {
		return "", err
	}

	if g.cfg.dryRun {
		contentStr, err := g.makeAndExecToStr(name, d.TemplateString, d.Context)
		if err != nil {
			return "", err
		}

		fmt.Printf("== DRY RUN: %s ==\n%s\n%s\n", fp, d.Description, contentStr)
		// INFO: no path generated or modified in dry run
		return "", nil
	}

	file, err := openFileWithStrategy(fp, d.Strategy, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := g.makeAndExec(name, d.TemplateString, file, d.Context); err != nil {
		return "", err
	}

	fmt.Printf("modified '%s'\n", fp)
	return fp, nil
}

// Build builds the templates.
func (g *Gojen) Build(tmplNames ...string) error {
	defs := g.defs
	if len(tmplNames) > 0 {
		defs = filterMap(defs, tmplNames)
	}

	for name, def := range defs {
		f, err := g.buildTemplate(name, def)
		if err != nil {
			return err
		}

		if f != "" {
			g.ModifiedFiles = append(g.ModifiedFiles, f)
		}
	}

	return nil
}
