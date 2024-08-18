package gojen

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileFlags(t *testing.T) {
	tests := []struct {
		strategy      Strategy
		expectedFlags int
	}{
		{StrategyTrunc, os.O_CREATE | os.O_RDWR | os.O_TRUNC},
		{StrategyAppendAtLast, os.O_CREATE | os.O_RDWR | os.O_APPEND},
		{StrategyIgnore, os.O_CREATE | os.O_RDWR},
	}

	for _, test := range tests {
		fflags := getFileFlags("test_fp", test.strategy)
		assert.Equal(t, test.expectedFlags, fflags)
	}

	testDir := t.TempDir()
	testFile := testDir + "/test.txt"
	if _, err := os.Create(testFile); err != nil {
		t.Fatal("failed to create test file")
	}

	fflags := getFileFlags(testFile, StrategyIgnore)
	assert.Equal(t, 0, fflags)
}

func TestIsConfirmed(t *testing.T) {
	assert.True(t, isConfirmed("y"))
	assert.True(t, isConfirmed("Y"))
	assert.False(t, isConfirmed("n"))
	assert.False(t, isConfirmed(""))
}

func TestMergeMaps(t *testing.T) {
	tests := []struct {
		name   string
		input  []map[string]int
		output map[string]int
	}{
		{
			name:   "No maps provided",
			input:  []map[string]int{},
			output: map[string]int{},
		},
		{
			name: "Single map provided",
			input: []map[string]int{
				{"a": 1, "b": 2},
			},
			output: map[string]int{
				"a": 1, "b": 2,
			},
		},
		{
			name: "Multiple maps with no overlapping keys",
			input: []map[string]int{
				{"a": 1, "b": 2},
				{"c": 3, "d": 4},
			},
			output: map[string]int{
				"a": 1, "b": 2, "c": 3, "d": 4,
			},
		},
		{
			name: "Multiple maps with overlapping keys",
			input: []map[string]int{
				{"a": 1, "b": 2},
				{"b": 20, "c": 3},
				{"d": 4, "a": 10},
			},
			output: map[string]int{
				"a": 10, "b": 20, "c": 3, "d": 4,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := mergeMaps(test.input...)
			assert.Equal(t, test.output, result)
		})
	}
}

func TestFilterMap(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]int
		keys   []string
		output map[string]int
	}{
		{
			name:   "Empty map and keys",
			input:  map[string]int{},
			keys:   []string{},
			output: map[string]int{},
		},
		{
			name:   "Non-empty map with no matching keys",
			input:  map[string]int{"a": 1, "b": 2, "c": 3},
			keys:   []string{"d", "e"},
			output: map[string]int{},
		},
		{
			name:   "Non-empty map with some matching keys",
			input:  map[string]int{"a": 1, "b": 2, "c": 3},
			keys:   []string{"b", "c"},
			output: map[string]int{"b": 2, "c": 3},
		},
		{
			name:   "Non-empty map with all matching keys",
			input:  map[string]int{"a": 1, "b": 2, "c": 3},
			keys:   []string{"a", "b", "c"},
			output: map[string]int{"a": 1, "b": 2, "c": 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := filterMap(test.input, test.keys)
			assert.Equal(t, test.output, result)
		})
	}
}
