package filemanager

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cirius-go/gojen/util"
)

type (
	// Config contains the configuration for the file manager.
	Config struct {
	}

	FileManager struct {
		cfg        *Config
		builtFiles map[string]string
	}
)

// C returns a new config with default params.
func C() *Config {
	return &Config{}
}

// New returns a new file manager instance.
func New() *FileManager {
	c := C()
	return NewWithConfig(c)
}

func NewWithConfig(c *Config) *FileManager {
	return &FileManager{
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
func (f *FileManager) WalkDir(dirPath string, openFile bool, handler func(e *FileInfo) error) error {
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
func (f *FileManager) CreateFileIfNotExist(path string, content string) (created bool, err error) {
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
func (f *FileManager) TruncWithContent(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// FileExists checks if the file exists.
func (f *FileManager) FileExists(path string) bool {
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
func (f *FileManager) AppendContent(path string, content string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

// AppendContentAfter appends the content after the line identified by lineIdent.
func (f *FileManager) AppendContentAfter(path string, lineIdent, content string) error {
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

func (f *FileManager) CopyFile(src, dst string) error {
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

func (f *FileManager) getLinesFromFile(path string) (map[string]bool, []string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	return f.getLines(string(content))
}

func (f *FileManager) getLines(content string) (map[string]bool, []string, error) {
	lines := make(map[string]bool)
	var orderedLines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines[strings.ToLower(line)] = true
			orderedLines = append(orderedLines, line)
		}
	}
	return lines, orderedLines, scanner.Err()
}

var (
	resetColor = "\033[0m"
	redColor   = "\033[31m"
)

func (f *FileManager) findMostDuplicatedSegment(linesA, linesB []string) (string, int) {
	maxCount := 0
	bestMatch := ""

	for i := 0; i < len(linesA); i++ {
		for j := 0; j < len(linesB); j++ {
			if linesA[i] == linesB[j] {
				matchCount := 0
				var matchedSegment strings.Builder

				for x, y := i, j; x < len(linesA) && y < len(linesB) && linesA[x] == linesB[y]; x, y = x+1, y+1 {
					matchedSegment.WriteString(redColor + linesA[x] + resetColor + "\n")
					matchCount++
				}

				if matchCount > maxCount {
					maxCount = matchCount
					bestMatch = matchedSegment.String()
				}
			}
		}
	}

	return bestMatch, maxCount
}

func (f *FileManager) CompareContentWithFile(content, dst string, ignoreLines util.MapExisting[string]) (percent float64, dstHighlighted string, err error) {
	_, orderedWordsA, err := f.getLines(content)
	if err != nil {
		return 0, "", err
	}
	_, orderedWordsB, err := f.getLinesFromFile(dst)
	if err != nil {
		return 0, "", err
	}
	dstHighlighted, matchCount := f.findMostDuplicatedSegment(orderedWordsA, orderedWordsB)
	percent = float64(matchCount) / float64(len(orderedWordsA)) * 100
	return
}

func (f *FileManager) CompareFile(src, dst string, ignoreLines util.MapExisting[string]) (percent float64, dstHighlighted string, err error) {
	_, orderedWordsA, err := f.getLinesFromFile(src)
	if err != nil {
		return 0, "", err
	}
	_, orderedWordsB, err := f.getLinesFromFile(dst)
	if err != nil {
		return 0, "", err
	}
	dstHighlighted, matchCount := f.findMostDuplicatedSegment(orderedWordsA, orderedWordsB)
	percent = float64(matchCount) / float64(len(orderedWordsA)) * 100
	return
}
