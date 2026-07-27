package state

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leonelquinteros/gotext"
)

func getDisabledFile() string {
	return os.ExpandEnv("$HOME/.config/uwelcome/disabled")
}

func IsDisabled() bool {
	_, err := os.Stat(getDisabledFile())
	return err == nil
}

func Enable(l *gotext.Locale) {
	err := os.Remove(getDisabledFile())
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(l.Get("Failed to enable the banner."))
		println(l.Get("Error ~> %s", err.Error()))
		return
	}
	fmt.Println(l.Get("The banner has been enabled."))
}

func Disable(l *gotext.Locale) {
	err := os.MkdirAll(filepath.Dir(getDisabledFile()), 0755)
	if err != nil {
		fmt.Println(l.Get("Failed to disable the banner."))
		println(l.Get("Error ~> %s", err.Error()))
		return
	}
	disabledFile, err := os.Create(getDisabledFile())
	if err != nil {
		fmt.Println(l.Get("Failed to disable the banner."))
		println(l.Get("Error ~> %s", err.Error()))
		return
	}
	fmt.Println(l.Get("The banner has been disabled."))
	err = disabledFile.Close()
	if err != nil {
		return
	}
}
