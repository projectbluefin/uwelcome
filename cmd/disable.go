package cmd

import (
	"fmt"
	"uwelcome/internal/state"

	"github.com/leonelquinteros/gotext"
)

func Disable(l *gotext.Locale) {
	if state.IsDisabled {
		fmt.Println(l.Get("The banner is already disabled."))
	} else {
		state.SetDisabled(l)
	}
}
