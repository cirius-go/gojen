package gojen

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gertd/go-pluralize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	pl         = pluralize.NewClient()
	titleCaser = cases.Title(language.AmericanEnglish)
)

func singular(s string) string {
	return pl.Singular(s)
}

func plural(s string) string {
	return pl.Plural(s)
}

func title(s string) string {
	return titleCaser.String(s)
}

func lower(s string) string {
	return strings.ToLower(s)
}

func upper(s string) string {
	return strings.ToUpper(s)
}

func makeDirAll(path string) error {
	dir, _ := filepath.Split(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return nil
}
