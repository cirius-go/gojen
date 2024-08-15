package gojen

// Config is a struct that holds the configuration for the Gojen instance.
type Config struct {
	debug          bool
	dryRun         bool
	parseArgs      bool
	customPipeline map[string]any
}

// SetDryRun sets the dryRun field of the Config struct.
func (c *Config) SetDryRun(dryRun bool) *Config {
	c.dryRun = dryRun
	return c
}

// RegisterPipeline registers a custom pipeline for templates.
func (c *Config) RegisterPipeline(name string, pipeline any) *Config {
	c.customPipeline[name] = pipeline
	return c
}

// ParseArgs sets the parseArgs field of the Config struct.
func (c *Config) ParseArgs(parseArgs bool) *Config {
	c.parseArgs = parseArgs
	return c
}

// SetDebug sets the debug field of the Config struct.
func (c *Config) SetDebug(d bool) *Config { c.debug = d; return c }

// C returns a new Config struct.
func C() *Config {
	return &Config{
		dryRun:         false,
		customPipeline: make(map[string]any),
	}
}
