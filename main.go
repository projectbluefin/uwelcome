package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"uwelcome/internal"
	"uwelcome/internal/config"
	"uwelcome/internal/locale"
	"uwelcome/internal/render"
	"uwelcome/internal/state"

	"github.com/leonelquinteros/gotext"
)

//go:embed all:locales
var localesFS embed.FS

const version = "0.5.0"

func main() {

	// Loads the locale based on the system's locale
	l := gotext.NewLocaleFSWithPath(locale.GetLocale(localesFS), localesFS, "locales")
	l.AddDomain("default")

	internal.InitCommonStrings(l)

	// Handles command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {

		// Prints the version
		case "version":
			fmt.Println(version)
			return

		// Toggles the banner
		case "toggle":
			state.Toggle(l)
			return

		// Enables the banner
		case "enable":
			state.Enable(l)
			return

		// Disables the banner
		case "disable":
			state.Disable(l)
			return

		// Shows the path and state of the config file
		case "status":
			config.CheckWrapper(l)
			return

		// Opens the default editor to the config file
		case "edit":
			config.Edit(l)
			return

		// Resets the config file
		case "reset":
			config.Reset(l)
			return

		case "help":
			internal.Usage(l)
			return

		// Default output
		default:
			internal.Warn(l.Get("Invalid command."), nil)
			internal.Usage(l)
			os.Exit(1)
		}
	}

	// Exits if the banner is disabled
	if state.IsDisabled {
		return
	}

	// Exits if the version is unsupported
	path := config.GetPath()
	if isSupported, _ := config.IsVersionSupported(path); path != "" && !isSupported {
		internal.Error(l.Get("Your config file's version isn't supported by this version of uwelcome. This version supports config version %d.", config.ConfigVersion), nil)
		os.Exit(1)
	}

	// Loads the configuration from the system's config file
	cfg, _ := config.GetConfig()

	// Grouping the sections to render
	sections := []string{
		render.GetGreeting(cfg, l),
		render.GetInfo(),
		render.GetGreenboot(l),
		render.GetCommands(cfg, l),
		render.GetMotd(cfg),
		render.GetLinks(cfg, l),
	}

	// Renders the output sections together
	fmt.Print(render.GetRender(cfg.Color, strings.Join(sections, "\n"), l))
}
