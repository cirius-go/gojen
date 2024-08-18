package gojen

import (
	"os"
)

// config is a struct that holds the configuration for the Gojen instance.
type config struct {
	wd             string
	dryRun         bool
	parseArgs      bool
	customPipeline map[string]any
}

// SetDryRun sets the dryRun field of the Config struct.
func (c *config) SetDryRun(dryRun bool) *config {
	c.dryRun = dryRun
	return c
}

// RegisterPipeline registers a custom pipeline for templates.
func (c *config) RegisterPipeline(name string, pipeline any) *config {
	c.customPipeline[name] = pipeline
	return c
}

// ParseArgs sets the parseArgs field of the Config struct.
func (c *config) ParseArgs(parseArgs bool) *config {
	c.parseArgs = parseArgs
	return c
}

// C returns a new Config struct.
func C() *config {
	wd, err := os.Getwd()
	if err != nil {
		panic("could not get current working directory")
	}
	return &config{
		wd:             wd,
		dryRun:         false,
		parseArgs:      false,
		customPipeline: make(map[string]any),
	}
}
