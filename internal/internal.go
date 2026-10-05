package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"uwelcome/internal/symbols"

	"github.com/leonelquinteros/gotext"
)

var errorStr, warnStr string

func InitCommonStrings(l *gotext.Locale) {
	errorStr = l.GetC("Error:", "Before an error message")
	warnStr = l.GetC("Warning:", "Before a warning message")
}

// Error uses printWithErr to print an error
func Error(msg string, err error) {
	printWithErr(errorStr, msg, err)
}

// Warn uses printWithErr to print a warning
func Warn(msg string, err error) {
	printWithErr(warnStr, msg, err)
}

// printWithErr prints a message to stderr followed by the go error that triggered it
func printWithErr(tag, msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s %s\n", tag, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", symbols.GetSymbol("arrow"), err)
	}
}

// DoesFileExist checks if a file at a given path exists
func DoesFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist), err
}
