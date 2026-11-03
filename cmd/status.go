package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"uwelcome/internal"
	"uwelcome/internal/config"
	"uwelcome/internal/symbols"

	"github.com/leonelquinteros/gotext"
)

// Status uses check() on a path and prints the status of the config of that path
func Status(l *gotext.Locale) {
	path, _ := config.Target(l)
	fmt.Println(path)

	if _, err := config.Load(path, l); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			internal.Warn(l.Get("The config file doesn't exist (yet). uWelcome will use either your system's config or the default."), nil)
			return
		}

		if errors.Is(err, config.ErrUnsupportedConfig) {
			internal.Warn(l.Get("Your config file's version isn't supported by this version of uWelcome. This version supports config version %d. uWelcome will use either your system's config or the default.", config.Version), nil)
			return
		}

		if !errors.Is(err, config.ErrUndecoded) {
			internal.Warn(l.Get("There is something wrong with the config file. uWelcome will use either your system's config or the default."), err)
			return
		}

		internal.Warn(l.Get("Unknown keys have been detected in your config file."), nil)

	}

	fmt.Println(l.Get("The config file is OK.") + " " + symbols.GetSymbol("checkmark"))
}
