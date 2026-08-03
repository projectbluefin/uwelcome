package state

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leonelquinteros/gotext"
)

var disabledFile = os.ExpandEnv("$HOME/.config/uwelcome/disabled")

func IsDisabled() bool {
	_, err := os.Stat(disabledFile)
	return err == nil
}

func Enable(l *gotext.Locale) {
	err := os.Remove(disabledFile)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(l.Get("Failed to enable the banner."))
		println(l.Get("Error ~> %s", err.Error()))
		return
	}
	fmt.Println(l.Get("The banner has been enabled."))
}

func Disable(l *gotext.Locale) {
	err := os.MkdirAll(filepath.Dir(disabledFile), 0755)
	if err != nil {
		fmt.Println(l.Get("Failed to disable the banner."))
		println(l.Get("Error ~> %s", err.Error()))
		return
	}
	disabledFile, err := os.Create(disabledFile)
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
