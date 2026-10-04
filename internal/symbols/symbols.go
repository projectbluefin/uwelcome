package symbols

import (
	"os/exec"
	"strings"
)

/*
nerdFontSymbols is a map of symbols that are available in Nerd Fonts Symbols.

https://nerdfonts.ytyng.com/
*/
var nerdFontSymbols = map[string]string{
	"arrow":           "",
	"bluesky":         "",
	"boot":            "󰟀",
	"checkmark":       "󰄳",
	"command_palette": " ",
	"discord":         "󰙯",
	"discuss":         "󰊌",
	"docs":            "󰈙",
	"donate":          "󱢏",
	"error":           "󰅙",
	"info":            "󰋼",
	"issues":          "",
	"link":            "󰌹",
	"mastodon":        "󰫑",
	"matrix":          "󰊌",
	"oci":             "󱋩",
	"source":          "󰊢",
	"warning":         "",
	"website":         "󰖟",
}

// asciiSymbols is a map of symbols that are available in ASCII.
var asciiSymbols = map[string]string{
	"arrow":           "~>",
	"checkmark":       "✓",
	"command_palette": ">_",
	"error":           "×",
	"oci":             "Oci:",
	"info":            "(i)",
	"warning":         "/!\\",
}

// hasNerdFontSymbols checks if the system has Nerd Fonts Symbols installed.
func hasNerdFontSymbols() bool {
	fontList, err := exec.Command("fc-list").Output()
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(fontList))
	return strings.Contains(lower, "nerdfont") // Trying something out with the nerd fonts
}

/*
GetSymbol returns the symbol for the given symbol name.

If the system has a Nerd Font installed, it will return the symbol embedded in the font. Otherwise, it will return the ASCII version of the symbol.
*/
func GetSymbol(symbolName string) string {
	if hasNerdFontSymbols() {
		if symbol, ok := nerdFontSymbols[symbolName]; ok {
			return symbol
		}
	}
	if symbol, ok := asciiSymbols[symbolName]; ok {
		return symbol
	}
	return ""
}
