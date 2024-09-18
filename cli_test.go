package gojen

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetOutput(t *testing.T) {
	b := &bytes.Buffer{}

	cfg := ConsoleC()
	c := NewConsoleWithConfig(cfg)

	c.SetOutput(b)

	assert.Equal(t, c.writer, b)
}

func TestPrintWithColor(t *testing.T) {
	t.Run("With color", func(t *testing.T) {
		i := &bytes.Buffer{}
		o := &bytes.Buffer{}

		cfg := ConsoleC()
		c := NewConsoleWithConfig(cfg)

		c.SetInput(i)
		c.SetOutput(o)

		c.Infof(true, "info with blue color")
		c.Successf(true, "success with green color")
		c.Warnf(true, "warning with yellow color")
		c.Dangerf(true, "danger with red color")

		blueMsg := "\x1b[34minfo with blue color\x1b[0m"
		greenMSg := "\x1b[32msuccess with green color\x1b[0m"
		yellowMsg := "\x1b[33mwarning with yellow color\x1b[0m"
		redMsg := "\x1b[31mdanger with red color\x1b[0m"

		lines := []string{blueMsg, greenMSg, yellowMsg, redMsg}

		assert.Equal(t, o.String(), strings.Join(lines, ""))
	})

	t.Run("Without color", func(t *testing.T) {
		i := &bytes.Buffer{}
		o := &bytes.Buffer{}

		cfg := ConsoleC().WithColor(false)
		c := NewConsoleWithConfig(cfg)

		c.SetInput(i)
		c.SetOutput(o)

		c.Infof(true, "info with blue color")
		c.Successf(true, "success with green color")
		c.Warnf(true, "warning with yellow color")
		c.Dangerf(true, "danger with red color")

		blueMsg := "info with blue color"
		greenMsg := "success with green color"
		yellowMsg := "warning with yellow color"
		redMsg := "danger with red color"

		lines := []string{blueMsg, greenMsg, yellowMsg, redMsg}

		assert.Equal(t, o.String(), strings.Join(lines, ""))
	})
}

func TestPrintWithSilent(t *testing.T) {
	i := &bytes.Buffer{}
	o := &bytes.Buffer{}

	cfg := ConsoleC()
	c := NewConsoleWithConfig(cfg)

	c.SetInput(i)
	c.SetOutput(o)

	c.Infof(false, "info with blue color")
	c.Successf(false, "success with green color")
	c.Warnf(false, "warning with yellow color")
	c.Dangerf(false, "danger with red color")

	assert.Equal(t, o.String(), "")
}

func TestPerformYesNo(t *testing.T) {
	t.Run("It should ask the user to confirm the action with danger color", func(t *testing.T) {
		msg := "Are you sure? (y/n)"

		t.Run("It should return true if the input is not 'y'", func(t *testing.T) {
			t.Parallel()
			i := &bytes.Buffer{}
			o := &bytes.Buffer{}

			cfg := ConsoleC().WithDelayReadInput(1 * time.Second)
			c := NewConsoleWithConfig(cfg)

			c.SetInput(i)
			c.SetOutput(o)

			go func(startTime time.Time) {
				for time.Since(startTime) < 1*time.Second {
					time.Sleep(1 * time.Second / 2)
					if i.String() != "" {
						break
					}

					if o.String() != "" {
						i.Write([]byte("t\n"))
					}
				}
			}(time.Now())

			result := c.PerformYesNo(msg)

			assert.False(t, result)
		})

		t.Run("It should return true if the input is 'y'", func(t *testing.T) {
			t.Parallel()
			i := &bytes.Buffer{}
			o := &bytes.Buffer{}

			cfg := ConsoleC().WithDelayReadInput(1 * time.Second)
			c := NewConsoleWithConfig(cfg)

			c.SetInput(i)
			c.SetOutput(o)

			go func(startTime time.Time) {
				for time.Since(startTime) < 1*time.Second {
					time.Sleep(1 * time.Second / 2)
					if i.String() != "" {
						break
					}

					if o.String() != "" {
						i.Write([]byte("y\n"))
					}
				}
			}(time.Now())

			result := c.PerformYesNo(msg)

			assert.True(t, result)
		})
	})
}

func TestSetInput(t *testing.T) {
	b := &bytes.Buffer{}

	cfg := ConsoleC()
	c := NewConsoleWithConfig(cfg)

	c.SetInput(b)

	assert.Equal(t, c.reader, b)
}
