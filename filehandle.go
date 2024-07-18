package gojen

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

	flags := os.O_CREATE | os.O_WRONLY

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
