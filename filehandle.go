package gojen

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func loadJSON[K any](jsonPath string, v *K) error {
	f, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(v); err != nil {
		return err
	}

	return nil
}

func makeDirAll(path string) error {
	dir, _ := filepath.Split(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func openFileWithStrategy(path string, strategy Strategy, perm fs.FileMode) (*os.File, error) {
	if err := makeDirAll(path); err != nil {
		return nil, err
	}

	flags := os.O_CREATE | os.O_RDWR

	switch strategy {
	case StrategyTrunc:
		flags |= os.O_TRUNC
	case StrategyAppend:
		flags |= os.O_APPEND
	case StrategyIgnore:
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				flags |= os.O_APPEND
			}
		} else {
			fmt.Printf("skipped to modify '%s'. This file is exist.\n", path)
		}
	}

	return os.OpenFile(path, flags, perm)
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

func readLines(path string) ([]string, error) {
	content, err := readFileContent(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}

func appendFileAtLine(path string, predicate func(l string) bool, content string) (bool, error) {
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
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	fcontent := strings.Join(first, "\n") + "\n" + content + strings.Join(last, "\n")

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
