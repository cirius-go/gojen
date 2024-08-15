package gojen

import (
	"fmt"

	"github.com/cirius-go/gojen/color"
)

// Bluef prints a blue string with a format.
func Bluef(format string, a ...interface{}) {
	fmt.Print(color.Bluef(format, a...))
}

// Greenf prints a green string with a format.
func Greenf(format string, a ...interface{}) {
	fmt.Print(color.Greenf(format, a...))
}

// Redf prints a red string with a format.
func Redf(format string, a ...interface{}) {
	fmt.Print(color.Redf(format, a...))
}
