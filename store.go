package gojen

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/cirius-go/gojen/util"
)

type (
	Args map[string]any

	// StoreConfig contains configurations for store.
	StoreConfig struct {
		silent   bool
		stateDir string // store template with states.
	}

	// store manage template declarations with parameters.
	store struct {
		cfg         *StoreConfig
		decls       map[string]*D
		args        Args
		builtStates []*State

		fm FileManager
		c  ConsoleManager
	}
)

// NewArgs returns a new Args.
func NewArgs(args ...Args) Args {
	a := make(Args)
	return a.Merge(args...)
}

// Extract extracts the given keys from the args.
func (a Args) Extract(keys ...string) (Args, []string) {
	if len(keys) == 0 {
		return a, nil
	}
	var (
		notFound = util.MapExisting[string]{}
		args     = make(Args)
	)
	for _, k := range keys {
		if v, ok := a[k]; ok {
			args[k] = v
			continue
		}
		notFound.Add(k)
	}
	return args, notFound.Keys()
}

// Clone clones the args.
func (a Args) Clone() Args {
	clone := make(Args, len(a))
	for k, v := range a {
		clone[k] = v
	}
	return clone
}

// Merge merges the given args.
func (a Args) Merge(as ...Args) Args {
	if len(as) == 0 {
		return a
	}

	for _, n := range as {
		for k, v := range n {
			a[k] = v
		}
	}
	return a
}

// StoreC returns a new StoreConfig with default.
func StoreC() *StoreConfig {
	return &StoreConfig{
		silent:   false,
		stateDir: ".gojen-state",
	}
}

// NewStoreWithConfig returns a new store with the given configuration.
func NewStoreWithConfig(cfg *StoreConfig, c ConsoleManager, fm FileManager) *store {
	if cfg == nil {
		panic("store config is required")
	}

	return &store{
		cfg:         cfg,
		decls:       make(map[string]*D),
		args:        make(Args),
		builtStates: []*State{},
		fm:          fm,
		c:           c,
	}
}

// NewStore returns a new store with default configuration.
func NewStore(c ConsoleManager, fm FileManager) *store {
	cfg := StoreC()

	return NewStoreWithConfig(cfg, c, fm)
}

// LoadDir walks through the directory and loads the template definitions.
func (s *store) LoadDir(dir string) error {
	return s.fm.WalkDir(dir, true, func(e *FileInfo) error {
		var fileDecoder FileDecoder

		switch e.Ext {
		case ".yaml", ".yml":
			fileDecoder = yaml.NewDecoder(e.File)
		case ".json":
			fileDecoder = json.NewDecoder(e.File)
		default:
			return nil
		}

		d := &D{}
		if err := fileDecoder.Decode(&d); err != nil {
			return err
		}

		if err := d.Validate(); err != nil {
			return fmt.Errorf("error validating template '%s': %w", e.Path, err)
		}

		if s.SetDecl(d) {
			s.c.Infof(s.cfg.silent, "Loaded template definition from: '%s'\n", e.Path)
		}

		return nil
	})
}

// Get returns the template definition with the given name.
func (s *store) GetDecl(name string) *D {
	return s.decls[name]
}

// SetDecl stores the template definition in the map.
func (s *store) SetDecl(d *D) bool {
	if d.Name == "" {
		s.c.Warnf(s.cfg.silent, "Template name is required. Skipping...\n")
		return false
	}

	if _, ok := s.decls[d.Name]; ok {
		r := s.c.PerformYesNo("Declaration with name '%s' already exists. Do you want to override it?\n", d.Name)
		if !r {
			return false
		}
	}

	s.decls[d.Name] = d
	return true
}

// GetArgs returns the args.
func (s *store) GetArgs(keys ...string) (Args, []string) {
	return s.args.Extract(keys...)
}

// UpdateArgs adds args to Gojen instance.
func (g *store) UpdateArgs(args Args) {
	if args == nil {
		return
	}

	g.args = g.args.Merge(args)
}

// AddState adds the built state.
func (g *store) AddState(s *State) {
	g.builtStates = append(g.builtStates, s)
}

// GetStates returns the built states.
func (g *store) GetStates() []*State {
	return g.builtStates
}

// LastState returns the last built state.
func (g *store) LastState() *State {
	if len(g.builtStates) == 0 {
		return nil
	}
	return g.builtStates[len(g.builtStates)-1]
}
