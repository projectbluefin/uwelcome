package cmd

import (
	"uwelcome/internal/state"

	"github.com/leonelquinteros/gotext"
)

func Toggle(l *gotext.Locale) {
	if state.IsDisabled() {
		state.SetEnabled(l)
	} else {
		state.SetDisabled(l)
	}
}
