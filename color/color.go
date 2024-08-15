package color

import (
	"fmt"
)

var (
	green = "\033[32m"
	reset = "\033[0m"
	red   = "\033[31m"
	blue  = "\033[34m"
)

func colorize(color string, msg string) string {
	return color + msg + reset
}

func Green(msg string) string {
	return colorize(green, msg)
}

func Red(msg string) string {
	return colorize(red, msg)
}

func Blue(msg string) string {
	return colorize(blue, msg)
}

func Bluef(format string, a ...interface{}) string {
	return colorize(blue, fmt.Sprintf(format, a...))
}

func Greenf(format string, a ...interface{}) string {
	return colorize(green, fmt.Sprintf(format, a...))
}

func Redf(format string, a ...interface{}) string {
	return colorize(red, fmt.Sprintf(format, a...))
}
