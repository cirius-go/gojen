package gojen

// Config is a struct that holds the configuration for the Gojen instance.
type Config struct {
	dryRun bool
}

// SetDryRun sets the dryRun field of the Config struct.
func (c *Config) SetDryRun(dryRun bool) *Config {
	c.dryRun = dryRun
	return c
}

// C returns a new Config struct.
func C() *Config {
	return &Config{
		dryRun: false,
	}
}
