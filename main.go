package main

import (
	"fmt"
	"os"
	"strings"

	i "uwelcome/internal"
	"uwelcome/internal/config"
	"uwelcome/internal/locale"
	"uwelcome/internal/motd"
	"uwelcome/internal/render"
	"uwelcome/internal/state"
	"uwelcome/internal/symbols"
	"uwelcome/internal/system"

	"github.com/leonelquinteros/gotext"
)

const version = "0.4.1"

func main() {

	// Loads the locale based on the system's locale
	l := gotext.NewLocale("locales", locale.DetectLocale())
	l.AddDomain("default")

	i.InitCommonStrings(l)

	isDisabled, _ := i.DoesFileExist(state.DisabledFile)

	// Handles command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {

		// Prints the version
		case "version":
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
			i.Usage(l)
			return

		// Default output
		default:
			i.Warn(l.Get("Invalid command."), nil)
			i.Usage(l)
			return
		}
	}

	// Exits if the banner is disabled

	if isDisabled {
		return
	}

	// Loads the configuration from the system's config file

	cfg, _ := config.GetConfig()

	var out strings.Builder

	// Greets the user

	fmt.Fprintf(&out, "# %s", cfg.Greeting.Prefix)

	if len(cfg.Greeting.Message) > 0 {
		out.WriteString(cfg.Greeting.Message)
	} else {
		out.WriteString(l.Get("Welcome to %s", system.OSName))
		if system.OSName == "" {
			l.Get("your system")
		}
	}

	out.WriteString(cfg.Greeting.Suffix)
	out.WriteString("\n")

	// Gets the image info

	if imageInfo := system.GetImageInfo(); imageInfo.ImageRef != "" || imageInfo.ImageTag != "" {
		fmt.Fprintf(&out, " %s `%s/%s:%s`\n\n", symbols.GetSymbol("oci"), imageInfo.ImageVendor, imageInfo.ImageName, imageInfo.ImageTag)
	}

	// Gets the Greenboot status

	if greenboot := system.GetGreenbootInfo(); greenboot != "" {
		fmt.Fprintf(&out, "\n %s %s:", symbols.GetSymbol("boot"), l.Get("Boot Status"))
		if greenboot == "healthy" {
			fmt.Fprintf(&out, "%s", "`"+l.Get("Healthy")+" "+symbols.GetSymbol("checkmark")+"`")
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
			var commandDesc string
			switch cmd.Desc {
			case "cmd_list":
				commandDesc = l.Get("List all available commands")
			case "cmd_report":
				commandDesc = l.Get("Report an issue to the contributors")
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
			default:
				commandDesc = cmd.Desc
			}
			fmt.Fprintf(&out, "| `%s` | %s |\n", cmd.Cmd, commandDesc)
		}
		fmt.Fprintf(&out, "\n")
	}

	// Gets a random motd

	if len(cfg.Motd.Messages) > 0 || len(cfg.Motd.Commands) > 0 {
		fmt.Fprintf(&out, "%s", motd.GetRandomMotd(cfg))
		fmt.Fprintf(&out, "\n\n")
	}

	// Gets the links

	if len(cfg.Links) > 0 {
		for _, link := range cfg.Links {
			var linkLabel string
			switch link.Name {
			case "bluesky":
				linkLabel = symbols.GetSymbol("bluesky") + " [" + l.Get("Bluesky") + "]"
			case "docs":
				linkLabel = symbols.GetSymbol("docs") + " [" + l.Get("Documentation") + "]"
			case "donate":
				linkLabel = symbols.GetSymbol("donate") + " [" + l.Get("Donate") + "]"
			case "discord":
				linkLabel = symbols.GetSymbol("discord") + " [" + l.Get("Discord") + "]"
			case "discuss":
				linkLabel = symbols.GetSymbol("discuss") + " [" + l.Get("Discuss") + "]"
			case "issues":
				linkLabel = symbols.GetSymbol("issues") + " [" + l.Get("Report an issue") + "]"
			case "mastodon":
				linkLabel = symbols.GetSymbol("mastodon") + " [" + l.Get("Mastodon") + "]"
			case "matrix":
				linkLabel = symbols.GetSymbol("matrix") + " [" + l.Get("Matrix") + "]"
			case "source":
				linkLabel = symbols.GetSymbol("source") + " [" + l.Get("Source code") + "]"
			case "website":
				linkLabel = symbols.GetSymbol("website") + " [" + l.Get("Website") + "]"
			default:
				linkLabel = symbols.GetSymbol("link") + " [" + link.Name + "]"
			}
			fmt.Fprintf(&out, " - %s(%s)\n", linkLabel, link.URL)
		}
		fmt.Fprintf(&out, "\n")
	}

	// Renders the output

	fmt.Print(render.GetRender(cfg.Color, out.String(), l))
}
