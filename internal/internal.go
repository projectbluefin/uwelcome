package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"uwelcome/internal/symbols"

	"github.com/leonelquinteros/gotext"
)

var errorStr, warnStr string

// InitCommonStrings sets translatable strings that can be used globally without squeezing locale in
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
		fmt.Fprintf(os.Stderr, "%s %s\n", symbols.GetSymbol("arrow"), err)
	}
}

// DoesFileExist checks if a file at a given path exists
func DoesFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, fs.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}

func UserDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "uwelcome"), nil
}
