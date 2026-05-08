// Package tcolor wraps stdout text in ANSI color codes when stdout is a
// terminal and NO_COLOR is not set. Off when piped to files or in tests.
package tcolor

import (
	"io"
	"os"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

var nameToANSI = map[string]string{
	"red":     Red,
	"green":   Green,
	"yellow":  Yellow,
	"blue":    Blue,
	"magenta": Magenta,
	"cyan":    Cyan,
	"white":   White,
}

func Enabled(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func Wrap(w io.Writer, ansi, text string) string {
	if !Enabled(w) {
		return text
	}
	return ansi + text + Reset
}

func WrapName(w io.Writer, colorName, text string) string {
	ansi, ok := nameToANSI[colorName]
	if !ok {
		ansi = Green
	}
	return Wrap(w, ansi, text)
}
