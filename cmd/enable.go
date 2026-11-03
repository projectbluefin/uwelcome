package cmd

import (
	"fmt"
	"uwelcome/internal/state"

	"github.com/leonelquinteros/gotext"
)

func Enable(l *gotext.Locale) {
	if state.IsDisabled() {
		state.SetEnabled(l)
	} else {
		fmt.Println(l.Get("The banner is already enabled."))
	}
}
