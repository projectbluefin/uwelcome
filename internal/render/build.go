package render

import (
	"strings"
	"uwelcome/internal/config"
	"uwelcome/internal/motd"
	"uwelcome/internal/symbols"
	"uwelcome/internal/system"

	"github.com/leonelquinteros/gotext"
)

func GetGreeting(cfg config.Config, l *gotext.Locale) string {
	var out strings.Builder

	out.WriteString("# ")
	out.WriteString(cfg.Greeting.Prefix)

	if cfg.Greeting.Message != "" {
		out.WriteString(cfg.Greeting.Message)
	} else {
		name := system.GetOSName()
		if name == "" {
			name = l.Get("your system")
		}
		out.WriteString(l.Get("Welcome to %s", name))
	}

	out.WriteString(cfg.Greeting.Suffix)

	return out.String()
}

func GetInfo() string {
	var out strings.Builder
	if imageInfo := system.GetImageInfo(); imageInfo.ImageVendor != "" && imageInfo.ImageName != "" {
		out.WriteString(" ")
		out.WriteString(symbols.GetSymbol("oci"))
		out.WriteString(" `")
		out.WriteString(imageInfo.ImageVendor)
		out.WriteString("/")
		out.WriteString(imageInfo.ImageName)
		out.WriteString(":")
		out.WriteString(imageInfo.ImageTag)
		out.WriteString("`\n")
	}
	return out.String()
}

func GetGreenboot(l *gotext.Locale) string {
	var out strings.Builder
	if greenboot := system.GetGreenbootInfo(); greenboot != "" {
		out.WriteString("\n ")
		out.WriteString(symbols.GetSymbol("boot"))
		out.WriteString(" `")
		out.WriteString(l.Get("Boot Status:"))
		out.WriteString(" ")
		if greenboot == system.HealthyState {
			out.WriteString(l.Get("Healthy"))
			out.WriteString(" ")
			out.WriteString(symbols.GetSymbol("checkmark"))
		} else {
			out.WriteString(greenboot)
		}
		out.WriteString("`")
	}
	return out.String()
}

func commandDesc(l *gotext.Locale, key string) string {
	switch key {
	case "banner_toggle":
		return l.Get("Toggle this banner on/off")
	case "boot-restore":
		return l.Get("Launch Bootloader Restoring tool")
	case "cli_pkg":
		return l.Get("Manage command line packages")
	case "cmd_inst":
		return l.Get("Launch Installer")
	case "cmd_list":
		return l.Get("List all available commands")
	case "cmd_report":
		return l.Get("Report an issue to the contributors")
	case "man_upd":
		return l.Get("Manually update the system")
	case "sys_info":
		return l.Get("View system information")
	case "term_bling":
		return l.Get("Enable terminal bling")
	default:
		return key
	}
}

func GetCommands(cfg config.Config, l *gotext.Locale) string {
	var out strings.Builder
	if len(cfg.Commands) > 0 {
		out.WriteString("| ")
		out.WriteString(symbols.GetSymbol("command_palette"))
		out.WriteString(" ")
		out.WriteString(l.Get("Command"))
		out.WriteString(" | ")
		out.WriteString(l.Get("Description"))
		out.WriteString(" | \n")
		out.WriteString("| ------------ | ----------- |\n")
		for _, cmd := range cfg.Commands {
			out.WriteString("| `")
			out.WriteString(cmd.Cmd)
			out.WriteString("` | ")
			out.WriteString(commandDesc(l, cmd.Desc))
			out.WriteString(" |\n")
		}
	}
	return out.String()
}

func GetMotd(cfg config.Config) string {
	if len(cfg.Motd.Messages) > 0 || len(cfg.Motd.Commands) > 0 {
		return motd.GetRandomMotd(cfg)
	}
	return ""
}

func linkLabel(l *gotext.Locale, key string) string {
	switch key {
	case "bluesky":
		return symbols.GetSymbol("bluesky") + " [" + l.Get("Bluesky") + "]"
	case "discord":
		return symbols.GetSymbol("discord") + " [" + l.Get("Discord") + "]"
	case "discuss":
		return symbols.GetSymbol("discuss") + " [" + l.Get("Discuss") + "]"
	case "docs":
		return symbols.GetSymbol("docs") + " [" + l.Get("Documentation") + "]"
	case "donate":
		return symbols.GetSymbol("donate") + " [" + l.Get("Donate") + "]"
	case "issues":
		return symbols.GetSymbol("issues") + " [" + l.Get("Report an issue") + "]"
	case "mastodon":
		return symbols.GetSymbol("mastodon") + " [" + l.Get("Mastodon") + "]"
	case "matrix":
		return symbols.GetSymbol("matrix") + " [" + l.Get("Matrix") + "]"
	case "source":
		return symbols.GetSymbol("source") + " [" + l.Get("Source code") + "]"
	case "website":
		return symbols.GetSymbol("website") + " [" + l.Get("Website") + "]"
	default:
		return symbols.GetSymbol("link") + " [" + key + "]"
	}
}

func GetLinks(cfg config.Config, l *gotext.Locale) string {
	var out strings.Builder
	for _, link := range cfg.Links {
		out.WriteString(" - ")
		out.WriteString(linkLabel(l, link.Name))
		out.WriteString("(")
		out.WriteString(link.URL)
		out.WriteString(")\n")
	}
	return out.String()
}
