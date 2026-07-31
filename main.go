package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"uwelcome/internal/config"
	"uwelcome/internal/locale"
	"uwelcome/internal/motd"
	"uwelcome/internal/render"
	"uwelcome/internal/state"
	"uwelcome/internal/symbols"
	"uwelcome/internal/system"

	"github.com/leonelquinteros/gotext"
)

const VERSION = "0.3.3"

//go:embed all:locales
var localesFS embed.FS

func main() {

	// Loads the locale based on the system's locale
	currentLocale := locale.DetectLocale(localesFS)
	l := gotext.NewLocaleFSWithPath(currentLocale, localesFS, "locales")
	l.AddDomain("default")

	isDisabled := state.IsDisabled()

	// Handles command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {

		// Prints the version
		case "--version", "-v", "version":
			fmt.Println(VERSION)
			return

		// Toggles the banner
		case "toggle":
			if isDisabled {
				state.Enable(l)
				return
			} else {
				state.Disable(l)
				return
			}

		// Enables the banner
		case "enable":
			if isDisabled {
				state.Enable(l)
				return
			} else {
				fmt.Println(l.Get("The banner is already enabled."))
				return
			}

		// Disables the banner
		case "disable":
			if isDisabled {
				fmt.Println(l.Get("The banner is already disabled."))
				return
			} else {
				state.Disable(l)
				return
			}

		// Returns the path to the current file
		case "config-path":
			fmt.Println(config.GetPath())
			return
		default:
			fmt.Println(l.Get("Invalid command"))
			return
		}
	}

	// Exits if the banner is disabled

	if isDisabled {
		os.Exit(0)
	}

	// Loads the configuration from the system's config file

	cfg := config.GetConfig()

	// Greets the user

	in := "# " + cfg.Greeting.Prefix

	if len(cfg.Greeting.Message) > 0 {
		in += cfg.Greeting.Message
	} else {
		in += l.Get("Welcome to %s", system.GetOSName())
	}

	in += cfg.Greeting.Suffix + "\n"

	// Gets the image info

	if imageInfo := system.GetImageInfo(); imageInfo.ImageRef != "" || imageInfo.ImageTag != "" {
		in += " " + symbols.GetSymbol("oci") + " `" + imageInfo.ImageRef + ":" + imageInfo.ImageTag + "` \n"
	} else if system.IsBootcSystem() {
		in += " " + symbols.GetSymbol("oci") + " `" + l.Get("Unknown system") + "` \n"
	}

	// Gets the Greenboot status

	if greenboot := system.GetGreenbootInfo(); greenboot != "" {
		in += "\n " + symbols.GetSymbol("boot") + " " + l.Get("Boot Status") + ":"
		if greenboot == "healthy" {
			in += "`" + l.Get("Healthy") + " " + symbols.GetSymbol("healthy") + "`"
		} else {
			in += "`" + greenboot + "`"
		}
		in += " \n"
	}

	// Command list

	if len(cfg.Commands) > 0 {
		in += " | " + symbols.GetSymbol("command_palette") + " " + l.Get("Command") + " | " + l.Get("Description") + " | \n"
		in += "| ------------ | ----------- |\n"
		var cmdSb strings.Builder
		for _, cmd := range cfg.Commands {
			switch cmd.Desc {
			case "cmd_list":
				cmd.Desc = l.Get("List all available commands")
			case "cli_pkg":
				cmd.Desc = l.Get("Manage command line packages")
			case "term_bling":
				cmd.Desc = l.Get("Enable terminal bling")
			case "banner_toggle":
				cmd.Desc = l.Get("Toggle this banner on/off")
			case "sys_info":
				cmd.Desc = l.Get("View system information")
			case "man_upd":
				cmd.Desc = l.Get("Manually update the system")
			}
			_, err := fmt.Fprintf(&cmdSb, "| `%s` | %s |\n", cmd.Cmd, cmd.Desc)
			if err != nil {
				return
			}
		}
		in += cmdSb.String()
		in += "\n"
	}

	// Gets a random tip

	if len(cfg.Motd.Messages) > 0 || len(cfg.Motd.Commands) > 0 {
		in += motd.GetRandomMessage(cfg) + "\n\n"
	}

	// Gets the links

	if len(cfg.Links) > 0 {
		var linkSb strings.Builder
		for _, link := range cfg.Links {
			switch link.Name {
			case "website":
				link.Name = symbols.GetSymbol("website") + " [" + l.Get("Website") + "]"
			case "issues":
				link.Name = symbols.GetSymbol("issues") + " [" + l.Get("Report an issue") + "]"
			case "docs":
				link.Name = symbols.GetSymbol("docs") + " [" + l.Get("Documentation") + "]"
			case "discuss":
				link.Name = symbols.GetSymbol("discuss") + " [" + l.Get("Discuss") + "]"
			case "discord":
				link.Name = symbols.GetSymbol("discord") + " [" + l.Get("Discord") + "]"
			case "matrix":
				link.Name = symbols.GetSymbol("matrix") + " [" + l.Get("Matrix") + "]"
			case "bluesky":
				link.Name = symbols.GetSymbol("bluesky") + " [" + l.Get("Bluesky") + "]"
			case "mastodon":
				link.Name = symbols.GetSymbol("mastodon") + " [" + l.Get("Mastodon") + "]"
			case "donate":
				link.Name = symbols.GetSymbol("donate") + " [" + l.Get("Donate") + "]"
			default:
				link.Name = symbols.GetSymbol("link") + " [" + link.Name + "]"
			}
			_, err := fmt.Fprintf(&linkSb, " - %s(%s)\n", link.Name, link.URL)
			if err != nil {
				return
			}
		}
		in += linkSb.String()
		in += "\n"
	}

	// Renders the output

	fmt.Print(render.GetRender(cfg.Color, in))
}
