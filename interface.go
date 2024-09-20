package gojen

import (
	"text/template"

	"github.com/cirius-go/gojen/util"
)

// Service Dependencies.
type (
	// FileManager is an interface is use for managing files.
	//go:generate mockery --name FileManager
	FileManager interface {
		WalkDir(dirPath string, openFile bool, handler func(e *FileInfo) error) error
		CreateFileIfNotExist(path string, content string) (created bool, err error)
		TruncWithContent(path string, content string) error
		FileExists(path string) bool
		AppendContent(path string, content string) error
		AppendContentAfter(path string, lineIdent, content string) error
		CompareFile(src, dst string) (percent float64, dstHighlighted string, err error)
		CompareContentWithFile(content, dst string) (percent float64, dstHighlighted string, err error)
	}

	// ConsoleManager is an interface that defines the methods for interacting with the
	// console.
	//go:generate mockery --name ConsoleManager
	ConsoleManager interface {
		TermWidth() int
		Printf(l bool, msg string, args ...any)
		Infof(l bool, msg string, args ...any)
		InfoStringf(msg string, args ...any) string
		Successf(l bool, msg string, args ...any)
		SuccessStringf(msg string, args ...any) string
		Dangerf(l bool, msg string, args ...any)
		DangerStringf(msg string, args ...any) string
		Warnf(l bool, msg string, args ...any)
		WarnStringf(msg string, args ...any) string
		PerformYesNo(msg string, args ...any) bool
		Scanln() ([]byte, error)
	}

	// PipelineManager is an interface that defines the methods for managing the
	// template pipeline functions.
	PipelineManager interface {
		GetFuncs() template.FuncMap
		UpdateFuncs(fn func(template.FuncMap) template.FuncMap)
	}

	// StoreManager is an interface that defines the methods for managing the
	// template declarations and it's parameters.
	StoreManager interface {
		LoadDir(dirPath string) error
		GetDecl(name string) *D
		SetDecl(d *D) bool
		GetArgs(keys ...string) (Args, []string)
		UpdateArgs(args Args)
		AddState(s *State)
		GetStates() []*State
		LastState() *State
	}

	// SeqSwitcher represents a branch switcher of the seq.
	SeqSwitcher interface {
		When(eName string, handler util.PRFunc[*Seq, *Seq]) SeqSwitcher
	}
)

type (
	// FileDecoder is an interface that defines the Decode method.
	FileDecoder interface {
		Decode(v any) error
	}
)
