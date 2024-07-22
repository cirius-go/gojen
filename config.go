package gojen

import (
	"text/template"
)

// Config is a struct that holds the configuration for the Gojen instance.
type Config struct {
	dryRun         bool
	customPipeline map[string]any
}

// SetDryRun sets the dryRun field of the Config struct.
func (c *Config) SetDryRun(dryRun bool) *Config {
	c.dryRun = dryRun
	return c
}

// RegisterPipeline registers a custom pipeline.
func (c *Config) RegisterPipeline(name string, pipeline any) *Config {
	switch v := pipeline.(type) {
	case func(template.FuncMap, string) string:
		// The custom pipeline is a function that takes a FuncMap and a string as
		// arguments and returns a string.
		c.customPipeline[name] = func(s string) string {
			return v(templateFuncs, s)
		}
	default:
		c.customPipeline[name] = pipeline
	}
	c.customPipeline[name] = pipeline
	return c
}

// C returns a new Config struct.
func C() *Config {
	return &Config{
		dryRun:         false,
		customPipeline: make(map[string]any),
	}
}
