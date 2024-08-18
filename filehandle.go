package gojen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// makeDirAll creates directory and its parents if not exist.
func makeDirAll(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil // No directory to create if it's just a file name
	}

	return os.MkdirAll(dir, os.ModePerm)
}

func readFileContent(path string) ([]byte, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file '%s' is not exist", path)
		}

		return nil, err
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("'%s' is directory", path)
	}

	return os.ReadFile(path)
}

func handleOnStrategyAppend(path string, predicate func(l string) bool, content string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	if stat.IsDir() {
		return false, fmt.Errorf("'%s' is directory", path)
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(fileContent), "\n")
	loc := -1
	for i, line := range lines {
		if predicate(line) {
			loc = i + 1
			break
		}
	}

	if loc == -1 {
		return false, nil // no line matched
	}

	if loc == len(lines) {
		lines = append(lines, "")
	}

	first := lines[:loc]
	last := lines[loc:]
	fcontent := strings.Join(first, "\n") + "\n\n" + content + "\n" + strings.Join(last, "\n")

	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return true, err
	}
	defer f.Close()

	if _, err := f.WriteString(fcontent); err != nil {
		return true, err
	}

	return true, nil
}
