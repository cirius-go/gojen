package gojen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfig_C tests the C function of the Config struct.
func TestConfig_C(t *testing.T) {
	// Create a new Config struct.
	config := C()
	assert.False(t, config.dryRun)
}

// TestConfig_SetDryRun tests the SetDryRun method of the Config struct.
func TestConfig_SetDryRun(t *testing.T) {
	// Create a new Config struct.
	config := C()
	// Set the dryRun field to true.
	config.SetDryRun(true)

	assert.True(t, config.dryRun)
}
