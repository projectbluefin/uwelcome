package state

import (
	"fmt"
	"os"
	"path/filepath"
	"uwelcome/internal"

	"github.com/leonelquinteros/gotext"
)

var disabledFile string = getDisabledFile()
var IsDisabled, _ = internal.DoesFileExist(disabledFile)

func getDisabledFile() string {
	dir, _ := os.UserConfigDir()
	if dir == "" {
		return dir
	}
	return filepath.Join(dir, "uwelcome", "disabled")
}

// setEnabled is a wrapper for os.remove
func SetEnabled(l *gotext.Locale) {
	err := os.Remove(disabledFile)
	if err != nil && !os.IsNotExist(err) {
		internal.Error(l.Get("Failed to enable the banner."), err)
		return
	}
	fmt.Println(l.Get("The banner has been enabled."))
}

// setDisabled is a wrapper for os.MkdirAll
func SetDisabled(l *gotext.Locale) {
	err := os.MkdirAll(filepath.Dir(disabledFile), 0755)
	if err != nil {
		internal.Error(l.Get("Failed to disable the banner."), err)
		return
	}
	newDisabledFile, err := os.Create(disabledFile)
	if err != nil {
		internal.Error(l.Get("Failed to disable the banner."), err)
		return
	}
	fmt.Println(l.Get("The banner has been disabled."))
	err = newDisabledFile.Close()
	if err != nil {
		return
	}
}
