package gojen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type (
	// FileManagerConfig contains the configuration for the file manager.
	FileManagerConfig struct {
	}

	fileManager struct {
		cfg        *FileManagerConfig
		builtFiles map[string]string
	}
)

// FileManagerC returns a new config with default params.
func FileManagerC() *FileManagerConfig {
	return &FileManagerConfig{}
}

// NewFileManager returns a new file manager instance.
func NewFileManager() *fileManager {
	c := FileManagerC()
	return NewFileManagerWithConfig(c)
}

func NewFileManagerWithConfig(c *FileManagerConfig) *fileManager {
	return &fileManager{
		cfg:        c,
		builtFiles: make(map[string]string),
	}
}

// FileInfo contains simple required information only.
type FileInfo struct {
	Name string
	Ext  string
	Path string
	File *os.File
}

// WalkDir walks through the directory and calls the handler for each file.
func (f *fileManager) WalkDir(dirPath string, openFile bool, handler func(e *FileInfo) error) error {
	stat, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("'%s' is not a directory", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		i := &FileInfo{
			Name: e.Name(),
		}
		i.Ext = filepath.Ext(i.Name)
		i.Path = filepath.Join(dirPath, e.Name())

		var next func(e *FileInfo) error = handler

		if openFile {
			next = func(e *FileInfo) error {
				file, err := os.Open(e.Path)
				if err != nil {
					return err
				}
				defer file.Close()

				e.File = file
				return handler(e)
			}
		}

		if err = next(i); err != nil {
			return err
		}
	}

	return nil
}

// CreateIfNotExist creates a file with the given content if it does not exist.
func (f *fileManager) CreateFileIfNotExist(path string, content string) (created bool, err error) {
	dir, _ := filepath.Split(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return false, err
		}
	}

	// Check if file exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// File already exists
	return false, nil
}

// TruncWithContent truncates the file with the given content.
func (f *fileManager) TruncWithContent(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// FileExists checks if the file exists.
func (f *fileManager) FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return false
	}

	return true
}

// AppendContent appends the content to the file.
func (f *fileManager) AppendContent(path string, content string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

// AppendContentAfter appends the content after the line identified by lineIdent.
func (f *fileManager) AppendContentAfter(path string, lineIdent, content string) error {
	lineIdent = strings.TrimSpace(lineIdent)
	// Read the entire file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Split the content into lines
	lines := strings.Split(string(fileContent), "\n")

	// Find the line with lineIdent and insert the new content
	newLines := []string{}
	identFound := false
	for _, line := range lines {
		newLines = append(newLines, line)
		if strings.TrimSpace(line) == lineIdent {
			newLines = append(newLines, content)
			identFound = true
		}
	}

	// If lineIdent is not found, return without modifying the file
	if !identFound {
		return nil
	}

	// Join the lines back into a single string
	newContent := strings.Join(newLines, "\n")

	// Write the modified contents back to the file
	err = os.WriteFile(path, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func (f *fileManager) CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	if err != nil {
		return err
	}

	return nil
}

func (f *fileManager) getWordsFromFile(path string) (map[string]bool, []string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	return f.getWords(string(content))
}

func (f *fileManager) getWords(content string) (map[string]bool, []string, error) {
	words := make(map[string]bool)
	var orderedWords []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		words[strings.ToLower(word)] = true
		orderedWords = append(orderedWords, word)
	}
	return words, orderedWords, scanner.Err()
}

func (f *fileManager) compareWords(wordsA, wordsB map[string]bool, orderedWordsB []string) (float64, string) {
	commonWords := 0
	for word := range wordsA {
		if wordsB[word] {
			commonWords++
		}
	}
	percentage := float64(commonWords) / float64(len(wordsA)) * 100

	var highlighted strings.Builder
	for _, word := range orderedWordsB {
		if wordsA[strings.ToLower(word)] {
			highlighted.WriteString(word)
		}
		highlighted.WriteString(" ")
	}

	return percentage, highlighted.String()
}

func (f *fileManager) CompareContentWithFile(content, dst string) (percent float64, dstHighlighted string, err error) {
	wordsA, _, err := f.getWords(content)
	if err != nil {
		return 0, "", err
	}
	wordsB, orderedWordsB, err := f.getWordsFromFile(dst)
	if err != nil {
		return 0, "", err
	}
	percent, dstHighlighted = f.compareWords(wordsA, wordsB, orderedWordsB)
	return
}

func (f *fileManager) CompareFile(src, dst string) (percent float64, dstHighlighted string, err error) {
	wordsA, _, err := f.getWordsFromFile(src)
	if err != nil {
		return 0, "", err
	}
	wordsB, orderedWordsB, err := f.getWordsFromFile(dst)
	if err != nil {
		return 0, "", err
	}
	percent, dstHighlighted = f.compareWords(wordsA, wordsB, orderedWordsB)
	return
}
