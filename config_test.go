package gojen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_SetDryRun(t *testing.T) {
	config := C()

	config.SetDryRun(true)
	assert.True(t, config.dryRun, "dryRun should be true")

	config.SetDryRun(false)
	assert.False(t, config.dryRun, "dryRun should be false")
}

func TestConfig_RegisterPipeline(t *testing.T) {
	config := C()

	// Register a custom pipeline
	pipelineName := "testPipeline"
	plFn := func(s string) string { return s }
	config.RegisterPipeline(pipelineName, plFn)

	// Assert the custom pipeline is registered correctly
	assert.Contains(t, config.customPipeline, pipelineName, "Pipeline should be registered")

	r1 := config.customPipeline[pipelineName].(func(s string) string)("test")
	r2 := plFn("test")
	assert.Equal(t, r1, r2, "Pipeline function should be stored correctly")
}

func TestConfig_ParseArgs(t *testing.T) {
	config := C()

	config.ParseArgs(true)
	assert.True(t, config.parseArgs, "parseArgs should be true")

	config.ParseArgs(false)
	assert.False(t, config.parseArgs, "parseArgs should be false")
}

func TestConfig_C(t *testing.T) {
	config := C()

	// Assert that the initial values of the struct fields are correct
	assert.False(t, config.dryRun, "dryRun should be false by default")
	assert.False(t, config.parseArgs, "parseArgs should be false by default")
	assert.NotNil(t, config.customPipeline, "customPipeline should be initialized")
	assert.Empty(t, config.customPipeline, "customPipeline should be empty initially")
}
