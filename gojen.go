package gojen

import (
	"fmt"
	"os"
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

// Build builds the templates.
func (g *Gojen) Build() error {
	built := []string{}
	for name, def := range g.defs {
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
	t := template.New(name)
	t = t.Funcs(template.FuncMap{
		"singular": singular,
		"plural":   plural,
		"title":    title,
		"lower":    lower,
	})
	if _, err := t.Parse(templateString); err != nil {
		return nil, err
	}
	return t, nil
}

// execTemplate executes a template.
func (g *Gojen) execTemplate(t *template.Template, writer *os.File) error {
	return t.Execute(writer, &g.context)
}

func (g *Gojen) execTemplateToString(t *template.Template) (string, error) {
	writer := &strings.Builder{}
	if err := t.Execute(writer, &g.context); err != nil {
		return "", err
	}

	return writer.String(), nil
}

// buildTemplate builds a template.
func (g *Gojen) buildTemplate(name string, d *D) (string, error) {
	for _, key := range d.RequiredContext {
		if _, exists := g.context[key]; !exists {
			return "", fmt.Errorf("missing required context key: %s", key)
		}
	}

	path, err := g.makeTemplate(d.Path, d.Path)
	if err != nil {
		return "", err
	}

	pathStr, err := g.execTemplateToString(path)
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
		contentStr, err := g.execTemplateToString(content)
		if err != nil {
			return "", err
		}

		fmt.Println(contentStr)
		return "", nil
	}

	if err := makeDirAll(pathStr); err != nil {
		return "", err
	}

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

	if err := g.execTemplate(content, file); err != nil {
		return "", err
	}

	fmt.Printf("generated %s\n", pathStr)
	return pathStr, nil
}
