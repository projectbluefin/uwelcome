package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"uwelcome/internal/symbols"
)

// Error uses printWithErr to print an error
func Error(msg string, err error) {
	printWithErr("error", msg, err)
}

// Warn uses printWithErr to print a warning
func Warn(msg string, err error) {
	printWithErr("warning", msg, err)
}

// printWithErr prints a message to stderr followed by the go error that triggered it
func printWithErr(symbol, msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s %s\n", symbols.GetSymbol(symbol), msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", symbols.GetSymbol("arrow"), err)
	}
}

func printNoErr(symbol string, msg string) {
	fmt.Printf("%s %s\n", symbols.GetSymbol(symbol), msg)

}

// Success prints a success message
func Success(msg string) {
	printNoErr("checkmark", msg)
}

func Info(msg string) {
	printNoErr("info", msg)
}

// DoesFileExist checks if a file at a given path exists
func DoesFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist), err
}
