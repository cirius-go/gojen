package gojen

import "github.com/cirius-go/gojen/util"

// config is a struct that holds the configuration for the Gojen instance.
type config struct {
	console              *ConsoleConfig
	pipeline             *PipelineConfig
	fileManager          *FileManagerConfig
	store                *StoreConfig
	silent               bool
	commentQuote         string
	storePath            string
	ignoreComparingLines util.MapExisting[string]
}

// SetCommentQuote sets the commentQuote field of the Config struct.
func (c *config) SetCommentQuote(commentQuote string) *config {
	c.commentQuote = commentQuote
	return c
}

// SetSilent sets the silent field of the Config struct.
func (c *config) SetSilent(silent bool) *config {
	c.silent = silent
	return c
}

// SetStorePath sets the storePath field of the Config struct.
func (c *config) SetStorePath(storePath string) *config {
	c.storePath = storePath
	return c
}

// IgnoreComparingLine adds a new ignoreCompareLineWith to the Config struct.
func (c *config) IgnoreComparingLine(lines ...string) *config {
	c.ignoreComparingLines.Add(lines...)
	return c
}

func (c *config) SetConsoleConfig(cfg *ConsoleConfig) *config {
	c.console = cfg
	return c
}

func (c *config) SetPipelineConfig(cfg *PipelineConfig) *config {
	c.pipeline = cfg
	return c
}

func (c *config) SetFileManagerConfig(cfg *FileManagerConfig) *config {
	c.fileManager = cfg
	return c
}

func (c *config) SetStoreConfig(cfg *StoreConfig) *config {
	c.store = cfg
	return c
}

// C returns a new Config struct.
func C() *config {
	return &config{
		console:              ConsoleC(),
		pipeline:             PipelineC(),
		fileManager:          FileManagerC(),
		store:                StoreC(),
		silent:               false,
		commentQuote:         "//",
		storePath:            ".gojen",
		ignoreComparingLines: make(util.MapExisting[string]),
	}
}
