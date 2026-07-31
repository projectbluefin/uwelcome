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

const version = "0.3.4"

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
			fmt.Println(version)
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
		return
	}

	// Loads the configuration from the system's config file

	cfg := config.GetConfig()

	// Greets the user

	var out strings.Builder

	fmt.Fprintf(&out, "# %s", cfg.Greeting.Prefix)

	if len(cfg.Greeting.Message) > 0 {
		out.WriteString(cfg.Greeting.Message)
	} else {
		out.WriteString(l.Get("Welcome to %s", system.GetOSName()))
	}

	out.WriteString(cfg.Greeting.Suffix)
	out.WriteString("\n")

	// Gets the image info

	if imageInfo := system.GetImageInfo(); imageInfo.ImageRef != "" || imageInfo.ImageTag != "" {
		fmt.Fprintf(&out, " %s `%s:%s` \n", symbols.GetSymbol("oci"), imageInfo.ImageRef, imageInfo.ImageTag)
	} else if system.IsBootcSystem() {
		fmt.Fprintf(&out, " %s `%s` \n", symbols.GetSymbol("oci"), l.Get("Unknown system"))
	}

	// Gets the Greenboot status

	if greenboot := system.GetGreenbootInfo(); greenboot != "" {
		fmt.Fprintf(&out, "\n %s %s:", symbols.GetSymbol("boot"), l.Get("Boot Status"))
		if greenboot == "healthy" {
			fmt.Fprintf(&out, "%s", "`"+l.Get("Healthy")+" "+symbols.GetSymbol("healthy")+"`")
		} else {
			fmt.Fprintf(&out, "%s", "`"+greenboot+"`")
		}
		fmt.Fprintf(&out, " \n")
	}

	// Command list

	if len(cfg.Commands) > 0 {
		fmt.Fprintf(&out, " | %s %s | %s | \n", symbols.GetSymbol("command_palette"), l.Get("Command"), l.Get("Description"))
		fmt.Fprintf(&out, "| ------------ | ----------- |\n")
		for _, cmd := range cfg.Commands {
			commandDesc := cmd.Desc
			switch cmd.Desc {
			case "cmd_list":
				commandDesc = l.Get("List all available commands")
			case "cli_pkg":
				commandDesc = l.Get("Manage command line packages")
			case "term_bling":
				commandDesc = l.Get("Enable terminal bling")
			case "banner_toggle":
				commandDesc = l.Get("Toggle this banner on/off")
			case "sys_info":
				commandDesc = l.Get("View system information")
			case "man_upd":
				commandDesc = l.Get("Manually update the system")
			}
			fmt.Fprintf(&out, "| `%s` | %s |\n", cmd.Cmd, commandDesc)
		}
		fmt.Fprintf(&out, "\n")
	}

	// Gets a random tip

	if len(cfg.Motd.Messages) > 0 || len(cfg.Motd.Commands) > 0 {
		fmt.Fprintf(&out, "%s", motd.GetRandomMessage(cfg))
		fmt.Fprintf(&out, "\n\n")
	}

	// Gets the links

	if len(cfg.Links) > 0 {
		for _, link := range cfg.Links {
			linkLabel := link.Name
			switch link.Name {
			case "website":
				linkLabel = symbols.GetSymbol("website") + " [" + l.Get("Website") + "]"
			case "issues":
				linkLabel = symbols.GetSymbol("issues") + " [" + l.Get("Report an issue") + "]"
			case "docs":
				linkLabel = symbols.GetSymbol("docs") + " [" + l.Get("Documentation") + "]"
			case "discuss":
				linkLabel = symbols.GetSymbol("discuss") + " [" + l.Get("Discuss") + "]"
			case "discord":
				linkLabel = symbols.GetSymbol("discord") + " [" + l.Get("Discord") + "]"
			case "matrix":
				linkLabel = symbols.GetSymbol("matrix") + " [" + l.Get("Matrix") + "]"
			case "bluesky":
				linkLabel = symbols.GetSymbol("bluesky") + " [" + l.Get("Bluesky") + "]"
			case "mastodon":
				linkLabel = symbols.GetSymbol("mastodon") + " [" + l.Get("Mastodon") + "]"
			case "donate":
				linkLabel = symbols.GetSymbol("donate") + " [" + l.Get("Donate") + "]"
			default:
				linkLabel = symbols.GetSymbol("link") + " [" + link.Name + "]"
			}
			fmt.Fprintf(&out, " - %s(%s)\n", linkLabel, link.URL)
		}
		fmt.Fprintf(&out, "\n")
	}

	// Renders the output

	fmt.Print(render.GetRender(cfg.Color, out.String()))
}
