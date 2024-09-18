package gojen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew test new gojen with default configuration.
func TestNew(t *testing.T) {
	g := New()
	c := C()
	assert.Equal(t, g.cfg, c, "config should be equal")
}
