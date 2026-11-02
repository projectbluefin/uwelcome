package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"uwelcome/cmd"
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
			cmd.Toggle(l)
			return

		// Enables the banner
		case "enable":
			cmd.Enable(l)
			return

		// Disables the banner
		case "disable":
			cmd.Disable(l)
			return

		// Shows the path and state of the config file
		case "status":
			cmd.Status(l)
			return

		// Opens the default editor to the config file
		case "edit":
			cmd.Edit(l)
			return

		// Resets the config file
		case "reset":
			cmd.Reset(l)
			return

		case "help":
			cmd.Usage(l)
			return

		// Default output
		default:
			internal.Warn(l.Get("Invalid command."), nil)
			cmd.Usage(l)
			os.Exit(1)
		}
	}

	// Exits if the banner is disabled
	if state.IsDisabled {
		return
	}

	// Loads the configuration from the system's config file
	cfg, _, _ := config.Get()

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
