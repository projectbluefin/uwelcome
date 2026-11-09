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
var version = "0.5.1"

func main() {

	hasArgs := len(os.Args) > 1
	if hasArgs && os.Args[1] == "version" {
		fmt.Println(version)
		return
	}

	// Loads the locale based on the system's locale
	l := gotext.NewLocaleFSWithPath(locale.GetLocale(localesFS), localesFS, "locales")
	l.AddDomain("default")

	internal.InitCommonStrings(l)

	// Handles command line arguments
	if hasArgs {
		commands := map[string]func(*gotext.Locale){
			"toggle": cmd.Toggle, "enable": cmd.Enable, "disable": cmd.Disable,
			"status": cmd.Status, "edit": cmd.Edit, "reset": cmd.Reset, "help": cmd.Usage,
		}

		toRun := commands[os.Args[1]]
		if toRun == nil {
			internal.Warn(l.Get("Invalid command."), nil)
			cmd.Usage(l)
			os.Exit(1)
		}
		toRun(l)
		return
	}

	// Exits if the banner is disabled
	if state.IsDisabled() {
		return
	}

	// Loads the configuration from the system's config file
	cfg := config.Get(l)

	// Groups the sections to render
	sections := []string{
		render.GetGreeting(cfg, l),
		render.GetInfo(),
		render.GetGreenboot(l),
		render.GetCommands(cfg, l),
		render.GetMotd(cfg),
		render.GetLinks(cfg, l),
	}

	// Renders the output sections together
	fmt.Print(render.GetRender(cfg.Color, strings.Join(sections, ""), l))
}
