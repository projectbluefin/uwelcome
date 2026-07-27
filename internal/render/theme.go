package render

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
)

// Uses OSC 11 query to detect the terminal background color
func queryTerminalBackground() string {

	if !term.IsTerminal(os.Stdout.Fd()) {
		return "dark"
	}

	oldState, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		return "dark"
	}
	defer func(fd uintptr, oldState *term.State) {
		err := term.Restore(fd, oldState)
		if err != nil {

		}
	}(os.Stdin.Fd(), oldState)

	fmt.Print("\033]11;?\033\\")

	// read response like "rgb:0000/0000/0000"
	buf := make([]byte, 64)
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		return "dark"
	}

	response := string(buf[:n])

	// parse rgb:RRRR/GGGG/BBBB
	_, after, ok := strings.Cut(response, "rgb:")
	if !ok {
		return "dark"
	}

	parts := strings.Split(after, "/")
	if len(parts) < 3 {
		return "dark"
	}

	// take just the first two chars (most significant byte)
	rVal := parts[0][:2]
	var r int
	_, err = fmt.Sscanf(rVal, "%x", &r)
	if err != nil {
		return ""
	}

	// luminance threshold — below 128 is dark
	if r < 128 {
		return "dark"
	}
	return "light"
}

func detectTheme() string {

	if s := os.Getenv("GLAMOUR_STYLE"); s == "light" || s == "dark" {
		return s
	}

	// TERM_BACKGROUND (Ghostty, iTerm2, etc.)
	if bg := os.Getenv("TERM_BACKGROUND"); bg == "light" || bg == "dark" {
		return bg
	}

	// COLORFGBG (most common, set by xterm, urxvt, kitty...)
	if cfg := os.Getenv("COLORFGBG"); cfg != "" {
		parts := strings.Split(cfg, ";")
		var bg int
		_, err := fmt.Sscanf(parts[len(parts)-1], "%d", &bg)
		if err != nil {
			return ""
		}
		if bg > 8 {
			return "light"
		}
		return "dark"
	}

	// Try OSC 11 query only if we have a TTY
	if term.IsTerminal(uintptr(os.Stdout.Fd())) {
		return queryTerminalBackground()
	}

	return "dark" // fallback
}
