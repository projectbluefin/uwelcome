package render

import (
	"os/exec"
	"strings"

	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
)

const defaultMargin uint = 2

func GetRender(color string, input string) string {
	var withColor bool = true
	var processColor map[string]string
	var out string

	switch color {
	case "auto":
		processColor = getAccentColor()
	case "blue", "green", "orange", "pink", "purple", "red", "slate", "teal", "yellow":
		processColor = colorMap[color]
	default:
		withColor = false
	}

	if !withColor {
		colorScheme := detectTheme()
		out, _ = glamour.Render(input, colorScheme)
	} else {
		r, _ := glamour.NewTermRenderer(
			glamour.WithStyles(getColorizedStyle(processColor)),
		)
		out, _ = r.Render(input)
		r.Close()
	}

	return out
}

// The colors are ANSI color codes (ANSI 256 - https://en.wikipedia.org/wiki/ANSI_escape_code#Colors)

var colorMap = map[string]map[string]string{
	"blue":   {"accent": "33", "link": "69"},
	"green":  {"accent": "34", "link": "28"},
	"orange": {"accent": "208", "link": "130"},
	"pink":   {"accent": "212", "link": "163"},
	"purple": {"accent": "165", "link": "164"},
	"red":    {"accent": "203", "link": "124"},
	"slate":  {"accent": "104", "link": "104"},
	"teal":   {"accent": "44", "link": "38"},
	"yellow": {"accent": "220", "link": "178"},
}

func getAccentColor() map[string]string {
	defaultColor := colorMap["blue"]
	cmd := exec.Command("dconf", "read", "/org/gnome/desktop/interface/accent-color")
	output, err := cmd.Output()
	if err != nil {
		return defaultColor
	}
	accent := strings.Trim(strings.TrimSpace(string(output)), "'")
	if color, ok := colorMap[accent]; ok {
		return color
	}
	return defaultColor
}

func getColorizedStyle(color map[string]string) ansi.StyleConfig {
	accent := color["accent"]
	link := color["link"]

	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				// BlockPrefix: "\n",
				BlockSuffix: "\n",
			},
			Margin: new(defaultMargin),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Italic: new(true),
			},
			Indent:      new(uint(1)),
			IndentToken: new("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: defaultMargin,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       new(accent),
				Bold:        new(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{Prefix: "▌ "},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{Prefix: "┃ "},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{Prefix: "│ "},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{Prefix: "┆ "},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "┊ ",
				Bold:   new(false),
			},
		},
		Strikethrough: ansi.StylePrimitive{CrossedOut: new(true)},
		Emph:          ansi.StylePrimitive{Italic: new(true)},
		Strong: ansi.StylePrimitive{
			Color: new(accent),
			Bold:  new(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  new(accent),
			Format: "\n──────\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			Ticked:   "[✓] ",
			Unticked: "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     new(link),
			Underline: new(true),
		},
		LinkText: ansi.StylePrimitive{Bold: new(true)},
		Image:    ansi.StylePrimitive{Underline: new(true)},
		ImageText: ansi.StylePrimitive{
			Format: "Image: {{.text}}",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: " ",
				Suffix: " ",
				Color:  new(accent),
				Bold:   new(true),
			},
		},
		CodeBlock:             ansi.StyleCodeBlock{},
		Table:                 ansi.StyleTable{},
		DefinitionDescription: ansi.StylePrimitive{BlockPrefix: "\n🠶 "},
	}
}
