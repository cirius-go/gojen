package testlib

import (
	"os"
	"path/filepath"
	"testing"
)

// CreateDir creates a directory for testing.
func CreateDir(t *testing.T, customPaths ...string) string {
	t.Helper()

	declDir := filepath.Join(append([]string{t.TempDir()}, customPaths...)...)
	if err := os.MkdirAll(declDir, 0755); err != nil {
		t.Logf("Error creating testing directory: %v", err)
		t.FailNow()
	}

	return declDir
}

func NewFileWithContent(t *testing.T, path string, content string) {
	t.Helper()

	file, err := os.Create(path)
	if err != nil {
		t.Logf("Error creating file: %v", err)
		t.FailNow()
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		t.Logf("Error writing to file: %v", err)
		t.FailNow()
	}
}
