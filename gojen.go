package gojen

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/cirius-go/gojen/util"
)

// Gojen is a struct that holds the Gojen instance.
type Gojen struct {
	cfg *config

	p PipelineManager
	f FileManager
	s StoreManager
	c ConsoleManager
}

// New returns a new Gojen instance with default configuration.
func New() *Gojen {
	c := C()

	return NewWithConfig(c)
}

// NewWithConfig returns a new Gojen instance.
func NewWithConfig(cfg *config) *Gojen {
	if cfg == nil {
		panic("config is required")
	}

	// template pipeline manager.
	pipelineCfg := PipelineC()
	p := NewPipelineWithConfig(pipelineCfg)

	c := NewCLI()

	f := NewFileManager()

	// store
	storeCfg := StoreC()
	s := NewStoreWithConfig(storeCfg, c, f)

	g := &Gojen{
		cfg: cfg,
		p:   p,
		f:   f,
		s:   s,
		c:   c,
	}

	return g
}

// LoadDecls loads the declarations from the given directory.
func (g *Gojen) LoadDecls(dirPaths ...string) error {
	for _, dirPath := range dirPaths {
		if err := g.s.LoadDir(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (g *Gojen) SetDecls(decls ...*D) {
	for _, d := range decls {
		g.s.SetDecl(d)
	}
}

// parseTemplate creates, executes a template and returns the result as a string.
func (g *Gojen) parseTemplate(args map[string]any, name string, templateString string) (string, error) {
	pipelineFns := g.p.GetFuncs()

	t, err := template.
		New(name).
		Funcs(pipelineFns).
		Option("missingkey=zero").
		Parse(templateString)

	if err != nil {
		return "", err
	}

	w := strings.Builder{}
	if err := t.Execute(&w, &args); err != nil {
		return "", err
	}

	return w.String(), nil
}

func (g *Gojen) build(seq *Seq) error {
	decl := g.s.GetDecl(seq.DName)
	if decl == nil {
		return fmt.Errorf("Declaration '%s' not found", seq.DName)
	}
	declElem := decl.GetElements(seq.EName)
	if declElem == nil {
		return fmt.Errorf("Element '%s' not found in declaration '%s'", seq.EName, seq.DName)
	}

	var (
		rawPath          = util.IfValue("", declElem.Path, decl.Path)
		storeArgs, _     = g.s.GetArgs()
		args             = NewArgs(storeArgs, decl.Args, declElem.Args)
		requiredArgNames = util.NewSlice(decl.Require, declElem.Require, seq.ForwardArgs.Keys())
	)

	if _, notFoundArgNames := args.Extract(requiredArgNames...); len(notFoundArgNames) > 0 {
		g.c.Dangerf(true, "Please provide the missing arguments [%s] in JSON format: ", strings.Join(notFoundArgNames, ", "))
		jsonArgs, err := g.c.Scanln()
		if err != nil {
			return err
		}
		parsedJSONArgs := NewArgs()
		if err = json.Unmarshal(jsonArgs, &parsedJSONArgs); err != nil {
			return err
		}
		for _, k := range notFoundArgNames {
			v, ok := parsedJSONArgs[k]
			if !ok {
				return fmt.Errorf("Required argument '%s' not found in the JSON", k)
			}
			args[k] = v
		}
	}

	forwardArgs := make(Args)
	for k := range seq.ForwardArgs {
		v, ok := args[k]
		if !ok {
			continue
		}

		forwardArgs[k] = v
	}
	g.s.UpdateArgs(forwardArgs)

	parsedPath, err := g.parseTemplate(args, "path", rawPath)
	if err != nil {
		return err
	}

	parsedTmpl, err := g.parseTemplate(args, "content", declElem.Template)
	if err != nil {
		return err
	}

	st := &State{
		seq:           seq,
		d:             decl,
		e:             declElem,
		Strategy:      declElem.Strategy,
		Confirm:       declElem.Confirm,
		DName:         decl.Name,
		EName:         declElem.Name,
		EAlias:        declElem.Alias,
		Args:          args,
		RawTmpl:       declElem.Template,
		ParsedTmpl:    parsedTmpl,
		RawPath:       rawPath,
		ParsedPath:    parsedPath,
		ForwardedArgs: forwardArgs,
	}
	g.s.AddState(st)

	g.c.Infof(!g.cfg.silent, "Built state: %s.%s\n", decl.Name, declElem.Name)
	g.c.Successf(!g.cfg.silent, fmt.Sprintf("%s\n", st))

	return nil
}

// Build builds the templates.
func (g *Gojen) Build(seq *Seq) error {
	var travelSeq func(n *Seq) error
	var flow = []string{}

	travelSeq = func(n *Seq) error {
		if n.DName == "" && n.EName == "" {
			return errors.New("invalid case selected")
		}

		if err := g.build(n); err != nil {
			return err
		}

		flow = append(flow, fmt.Sprintf("%s.%s", n.DName, n.EName))

		if len(n.Cases) == 0 {
			if n.Next != nil {
				return travelSeq(n.Next)
			}

			return nil
		}

		g.c.Dangerf(true, "Which case do you want to choose?\n")
		i := 1
		mapIndexCases := map[int]*Seq{}
		util.LoopStrMap(n.Cases, func(k string, c *Seq) {
			g.c.Infof(true, "%d) %s.%s\n", i, c.DName, c.EName)
			mapIndexCases[i] = c
			i++
		})

		g.c.Dangerf(true, "Select case: ")
		selectedBytes, err := g.c.Scanln()
		if err != nil {
			return err
		}

		selected, err := strconv.Atoi(string(selectedBytes))
		if err != nil {
			return err
		}

		c, ok := mapIndexCases[selected]
		if !ok {
			return fmt.Errorf("invalid case selected")
		}

		if err := g.build(c); err != nil {
			return err
		}

		flow = append(flow, fmt.Sprintf("%s.%s", c.DName, c.EName))

		if c.Next != nil {
			return travelSeq(c.Next)
		}

		return nil
	}
	if err := travelSeq(seq.root); err != nil {
		return err
	}

	g.c.Successf(!g.cfg.silent, "built sequences: %s\n", strings.Join(flow, " -> "))

	return nil
}

func (g *Gojen) applyState(s *State) error {
	if s.Confirm {
		ok := g.c.PerformYesNo("Do you want to apply template '%s.%s'? ", s.DName, s.EName)
		if !ok {
			return nil
		}
	}

	switch s.Strategy {
	case StrategyInit:
		created, err := g.f.CreateFileIfNotExist(s.ParsedPath, s.ParsedTmpl)
		if err != nil {
			return err
		}

		if created {
			g.c.Successf(!g.cfg.silent, "Created file '%s' with:\n%s\n", s.ParsedPath, s.ParsedTmpl)

			return nil
		}

		g.c.Infof(!g.cfg.silent, "File already exists: '%s'. Skipped to init the file\n", s.ParsedPath)
		return nil
	case StrategyTrunc:
		if ok := g.c.PerformYesNo("Do you want to truncate %s?", s.ParsedPath); !ok {
			g.c.Infof(!g.cfg.silent, "User skipped to truncate file: %s\n", s.ParsedPath)
			return nil
		}

		return g.f.TruncWithContent(s.ParsedPath, s.ParsedTmpl)
	case StrategyAppendAtLast:
		exist := g.f.FileExists(s.ParsedPath)
		if !exist {
			if ok := g.c.PerformYesNo("File %s does not exist. Do you want to create it to append parsed content?", s.ParsedPath); !ok {
				g.c.Infof(!g.cfg.silent, "User skipped to create file: %s\n", s.ParsedPath)
				return nil
			}

			return g.f.TruncWithContent(s.ParsedPath, s.ParsedTmpl)
		}

		return g.f.AppendContent(s.ParsedPath, s.ParsedTmpl)
	case StrategyAppend:
		exist := g.f.FileExists(s.ParsedPath)
		if !exist {
			g.c.Infof(!g.cfg.silent, "File %s does not exist. Skipped to append parsed content\n", s.ParsedPath)
			return nil
		}

		alias := util.IfValue("", s.EAlias, s.EName)
		rawLineIdent := fmt.Sprintf("%s +gojen:append=%s", g.cfg.commentQuote, alias)
		parsedLineIdent, err := g.parseTemplate(s.Args, alias, rawLineIdent)
		if err != nil {
			return err
		}
		return g.f.AppendContentAfter(s.ParsedPath, parsedLineIdent, s.ParsedTmpl)
	default:
		fmt.Println("TODO: implement other strategies")
	}

	return nil
}

// Gojen applies the built templates.
func (g *Gojen) Apply() error {
	for _, bs := range g.s.GetStates() {
		if err := g.applyState(bs); err != nil {
			return err
		}
	}

	return nil
}
