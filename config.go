package gojen

// config is a struct that holds the configuration for the Gojen instance.
type config struct {
	silent       bool
	commentQuote string
	storePath    string
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

// C returns a new Config struct.
func C() *config {
	return &config{
		silent:       false,
		commentQuote: "//",
		storePath:    ".gojen",
	}
}
