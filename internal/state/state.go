package state

import (
	"fmt"
	"os"
	"path/filepath"
	"uwelcome/internal"

	"github.com/leonelquinteros/gotext"
)

// getDisabledFile retrieves the 'disabled' file's path
func getDisabledFile() (string, error) {
	dir, err := internal.UserDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "disabled"), nil
}

// IsDisabled returns if the 'disabled' file exists or not
func IsDisabled() bool {
	disabledFile, err := getDisabledFile()
	if err != nil {
		return false
	}
	exists, _ := internal.DoesFileExist(disabledFile)
	return exists
}

// SetEnabled deletes the 'disabled' file
func SetEnabled(l *gotext.Locale) {
	disabledFile, pathErr := getDisabledFile()
	if pathErr != nil {
		internal.Error(l.Get("Failed to get the disabled file path."), pathErr)
		os.Exit(1)
	}
	removeErr := os.Remove(disabledFile)
	if removeErr != nil && !os.IsNotExist(removeErr) {
		internal.Error(l.Get("Failed to enable the banner."), removeErr)
		os.Exit(1)
	}
	fmt.Println(l.Get("The banner has been enabled."))
}

// SetDisabled creates the 'disabled' file
func SetDisabled(l *gotext.Locale) {
	disabledFile, pathErr := getDisabledFile()
	if pathErr != nil {
		internal.Error(l.Get("Failed to get the disabled file path."), pathErr)
		os.Exit(1)
	}
	if makeErr := os.MkdirAll(filepath.Dir(disabledFile), 0755); makeErr != nil {
		internal.Error(l.Get("Failed to disable the banner."), makeErr)
		os.Exit(1)
	}
	if writeErr := os.WriteFile(disabledFile, nil, 0644); writeErr != nil {
		internal.Error(l.Get("Failed to disable the banner."), writeErr)
		os.Exit(1)
	}
	fmt.Println(l.Get("The banner has been disabled."))
}
