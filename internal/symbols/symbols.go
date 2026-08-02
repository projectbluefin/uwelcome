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
	"bluesky":         "",
	"boot":            "󰟀",
	"command_palette": " ",
	"discord":         "󰙯",
	"discuss":         "󰊌",
	"docs":            "󰈙",
	"donate":          "󱢏",
	"healthy":         "󰄳",
	"info":            "󰋼",
	"issues":          "",
	"link":            "󰌹",
	"mastodon":        "󰫑",
	"matrix":          "󰊌",
	"oci":             "󱋩",
	"source":          "󰊢",
	"website":         "󰖟",
}

// asciiSymbols is a map of symbols that are available in ASCII.
var asciiSymbols = map[string]string{
	"command_palette": ">_",
	"oci":             "[Ci]",
	"healthy":         "✓",
	"info":            "(i)",
}

// hasNerdFontSymbols checks if the system has Nerd Fonts Symbols installed.
func hasNerdFontSymbols() bool {
	out, err := exec.Command("fc-list").Output()
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(out))
	return strings.Contains(lower, "symbolsnerdfont") ||
		strings.Contains(lower, "nerdfontssymbolsonly")
}

// GetSymbol returns the symbol for the given symbol name. If the system has Nerd Fonts Symbols installed, it will return the Nerd Fonts Symbols version of the symbol. Otherwise, it will return the ASCII version of the symbol.
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
