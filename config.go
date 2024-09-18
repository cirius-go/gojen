package gojen

import (
	"os"
	"path/filepath"
)

// config is a struct that holds the configuration for the Gojen instance.
type config struct {
	silent       bool
	wd           string
	dryRun       bool
	commentQuote string
}

// SetDryRun sets the dryRun field of the Config struct.
func (c *config) SetDryRun(dryRun bool) *config {
	c.dryRun = dryRun
	return c
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

// TmpDir returns the temporary directory path.
func (c *config) TmpDir() string {
	return filepath.Join(c.wd, ".gojen", "tmp")
}

// SetWorkingDir sets the working directory of the Config struct.
func (c *config) SetWorkingDir(wd string) *config {
	c.wd = wd
	return c
}

// C returns a new Config struct.
func C() *config {
	wd, err := os.Getwd()
	if err != nil {
		panic("could not get current working directory")
	}
	return &config{
		silent:       false,
		wd:           wd,
		dryRun:       false,
		commentQuote: "//",
	}
}
