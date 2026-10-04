package state

import (
	"os"
	"path/filepath"
	i "uwelcome/internal"

	"github.com/leonelquinteros/gotext"
)

var DisabledFile string = os.ExpandEnv("$HOME/.config/uwelcome/disabled")

// Enable is a wrapper for os.remove
func Enable(l *gotext.Locale) {
	err := os.Remove(DisabledFile)
	if err != nil && !os.IsNotExist(err) {
		i.Error(l.Get("Failed to enable the banner."), err)
		return
	}
	i.Success(l.Get("The banner has been enabled."))
}

// Disable is a wrapper for os.MkdirAll
func Disable(l *gotext.Locale) {
	err := os.MkdirAll(filepath.Dir(DisabledFile), 0755)
	if err != nil {
		i.Error(l.Get("Failed to disable the banner."), err)
		return
	}
	disabledFile, err := os.Create(DisabledFile)
	if err != nil {
		i.Error(l.Get("Failed to disable the banner."), err)
		return
	}
	i.Success(l.Get("The banner has been disabled."))
	err = disabledFile.Close()
	if err != nil {
		return
	}
}
