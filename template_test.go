package gojen

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestM_Store tests the Store method of the M type.
func TestM_Store(t *testing.T) {
	m := make(M)

	// Test storing a template definition with a valid strategy
	d := &D{
		Name: "testTemplate",
		Select: []*E{
			{Template: "testTemplate", Strategy: StrategyTrunc},
		},
	}

	result := m.Store(d)
	assert.Equal(t, d, result[d.Name], "The template definition should be stored correctly")

	// Test storing a template definition with an invalid strategy
	dInvalid := &D{
		Name: "invalidTemplate",
		Select: []*E{
			{Template: "invalidTemplate", Strategy: "invalidStrategy"},
		},
	}

	assert.Panics(t, func() { m.Store(dInvalid) }, "Storing a template with an invalid strategy should panic")
}

// TestE_MarshalUnmarshalJSON tests the MarshalJSON and UnmarshalJSON methods of the E type.
func TestE_MarshalUnmarshalJSON(t *testing.T) {
	e := &E{
		Template: "line1\nline2\nline3",
		Require:  []string{"key1", "key2"},
		Strategy: StrategyTrunc,
		Confirm:  true,
	}

	data, err := json.Marshal(e)
	assert.NoError(t, err, "Marshalling should not produce an error")

	expectedJSON := `{"template":["line1","line2","line3"],"require":["key1","key2"],"strategy":"trunc","confirm":true}`
	assert.JSONEq(t, expectedJSON, string(data), "The marshalled JSON should match the expected format")

	var eUnmarshalled E
	err = json.Unmarshal(data, &eUnmarshalled)
	assert.NoError(t, err, "Unmarshalling should not produce an error")
	assert.Equal(t, e, &eUnmarshalled, "The unmarshalled struct should match the original")
}

// TestD_Clone tests the clone method of the D type.
func TestD_Clone(t *testing.T) {
	d := &D{
		Path:         "/test/path",
		Required:     []string{"req1"},
		Name:         "testName",
		Context:      map[string]any{"key1": "value1"},
		Select:       []*E{{Template: "testTemplate"}},
		Dependencies: []string{"dep1"},
		Description:  "testDescription",
	}

	clone := d.clone()

	assert.Equal(t, d, clone, "The cloned D struct should be equal to the original")
	assert.NotSame(t, d, clone, "The cloned D struct should not be the same instance as the original")
}

// TestD_MergeContext tests the mergeContext method of the D type.
func TestD_MergeContext(t *testing.T) {
	d := &D{
		Context: map[string]any{"key1": "value1"},
	}

	merged := d.mergeContext(map[string]any{"key2": "value2", "key1": "newValue1"})

	expectedContext := map[string]any{
		"key1": "value1", // d.Context takes precedence
		"key2": "value2",
	}

	assert.Equal(t, expectedContext, merged.Context, "The context should be merged correctly with d.Context taking precedence")
}
