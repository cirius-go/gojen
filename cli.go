package gojen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/term"
)

var (
	resetColor  = "\033[0m"
	redColor    = "\033[31m"
	greenColor  = "\033[32m"
	yellowColor = "\033[33m"
	blueColor   = "\033[34m"
)

// ConsoleConfig is the configuration for the CLI.
type ConsoleConfig struct {
	Color          bool
	readInputDelay time.Duration
}

// WithColor sets the color for the CLI.
func (c *ConsoleConfig) WithColor(color bool) *ConsoleConfig {
	c.Color = color
	return c
}

// WithDelayReadInput sets the delay for reading input.
func (c *ConsoleConfig) WithDelayReadInput(delay time.Duration) *ConsoleConfig {
	c.readInputDelay = delay
	return c
}

type Console struct {
	cfg *ConsoleConfig

	writer io.Writer
	reader io.Reader
}

// ConsoleC returns a new CLIConfig instance.
func ConsoleC() *ConsoleConfig {
	return &ConsoleConfig{
		Color: true,
	}
}

func NewConsoleWithConfig(cfg *ConsoleConfig) *Console {
	return &Console{
		cfg,
		os.Stdout,
		os.Stdin,
	}
}

// NewConsole returns a new CLI instance.
func NewConsole() *Console {
	cfg := ConsoleC()

	return NewConsoleWithConfig(cfg)
}

// SetInput sets the input reader for the cli.
func (c *Console) SetInput(reader io.Reader) {
	c.reader = reader
}

// SetOutput sets the output writer for the cli.
func (c *Console) SetOutput(writer io.Writer) {
	c.writer = writer
}

func (c *Console) makeColor(color string, format string, args ...any) string {
	return color + fmt.Sprintf(format, args...) + resetColor
}

func (c *Console) sprintf(color string, format string, args ...any) string {
	if c.cfg.Color {
		return c.makeColor(color, format, args...)
	}
	return fmt.Sprintf(format, args...)
}

func (c *Console) printf(silent bool, color string, format string, args ...any) {
	if !silent {
		return
	}

	msg := c.sprintf(color, format, args...)
	fmt.Fprint(c.writer, msg)
}

// Printf logs a message to the console.
func (c *Console) Printf(l bool, msg string, args ...any) {
	c.printf(l, "", msg, args...)
}

// Warnf logs a warning message to the console.
func (c *Console) Warnf(l bool, msg string, args ...any) {
	c.printf(l, yellowColor, msg, args...)
}

func (c *Console) WarnStringf(msg string, args ...any) string {
	return c.sprintf(yellowColor, msg, args...)
}

// Dangerf logs a danger message to the console.
func (c *Console) Dangerf(l bool, msg string, args ...any) {
	c.printf(l, redColor, msg, args...)
}

func (c *Console) DangerStringf(msg string, args ...any) string {
	return c.sprintf(redColor, msg, args...)
}

// Infof logs an info message to the console.
func (c *Console) Infof(l bool, msg string, args ...any) {
	c.printf(l, blueColor, msg, args...)
}

func (c *Console) InfoStringf(msg string, args ...any) string {
	return c.sprintf(blueColor, msg, args...)
}

// Successf logs a success message to the console.
func (c *Console) Successf(l bool, msg string, args ...any) {
	c.printf(l, greenColor, msg, args...)
}

func (c *Console) SuccessStringf(msg string, args ...any) string {
	return c.sprintf(greenColor, msg, args...)
}

// Scanln scans the input from the user.
func (c *Console) Scanln() ([]byte, error) {
	scanner := bufio.NewScanner(c.reader)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return scanner.Bytes(), nil
}

// TermWidth returns the terminal width.
func (c *Console) TermWidth() int {
	fd := int(os.Stdout.Fd())

	width, _, err := term.GetSize(fd)
	if err != nil {
		panic(err)
	}
	return width
}

// PerformYesNo asks the user to confirm the action with danger color.
// It uses the default io.Reader to read the input.
func (c *Console) PerformYesNo(msg string, args ...any) bool {
	c.Dangerf(true, msg, args...)

	if c.cfg.readInputDelay > 0 {
		time.Sleep(c.cfg.readInputDelay)
	}

	var input string
	if _, err := fmt.Fscanln(c.reader, &input); err != nil {
		return false
	}

	switch input {
	case "y", "Y", "yes", "YES", "true", "TRUE", "1":
		return true
	default:
		return false
	}
}
