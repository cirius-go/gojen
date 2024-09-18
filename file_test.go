package gojen_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cirius-go/gojen"
	"github.com/cirius-go/gojen/util/testlib"
)

func TestWalkDir(t *testing.T) {
	t.Run("It should return an error if the directory does not exist", func(t *testing.T) {
		c := gojen.FileManagerC()
		f := gojen.NewFileManagerWithConfig(c)

		tmpDir := t.TempDir()
		nonExistentDir := filepath.Join(tmpDir, "nonexistent")

		err := f.WalkDir(nonExistentDir, false, func(e *gojen.FileInfo) error {
			return nil
		})

		assert.NotNil(t, err)
	})

	t.Run("It should open the file if openFile is true", func(t *testing.T) {
		c := gojen.FileManagerC()
		f := gojen.NewFileManagerWithConfig(c)

		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.txt")
		testlib.NewFileWithContent(t, testFile, "This is the text message")

		err := f.WalkDir(tmpDir, true, func(e *gojen.FileInfo) error {
			assert.Equal(t, "test.txt", e.Name)
			assert.Equal(t, ".txt", e.Ext)
			assert.NotNil(t, e.File)

			return nil
		})

		assert.Nilf(t, err, "Error walking directory: %v", err)
	})

	t.Run("It shouldn't open the file if openFile is false", func(t *testing.T) {
		c := gojen.FileManagerC()
		f := gojen.NewFileManagerWithConfig(c)

		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.txt")
		testlib.NewFileWithContent(t, testFile, "This is the text message")

		err := f.WalkDir(tmpDir, false, func(e *gojen.FileInfo) error {
			assert.Equal(t, "test.txt", e.Name)
			assert.Equal(t, ".txt", e.Ext)
			assert.Nil(t, e.File)

			return nil
		})

		assert.Nilf(t, err, "Error walking directory: %v", err)
	})
}
