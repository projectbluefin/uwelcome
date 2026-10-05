package render

import (
	"fmt"
	"math"
	"os"
	"time"

	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
	"github.com/leonelquinteros/gotext"
	"github.com/rymdport/portal/settings"

	i "uwelcome/internal"
)

// colorMap maps color names to their corresponding hex values. These colors are used for accent and link styling in the rendered output.
var colorMap = map[string]string{
	"blue":   "#3584E4",
	"teal":   "#2190A4",
	"green":  "#3A944A",
	"yellow": "#C88800",
	"orange": "#ED5B00",
	"red":    "#E62D42",
	"pink":   "#D56199",
	"purple": "#9141AC",
	"slate":  "#6F8396",
}

// RGB is a struct that matches the (ddd) DBus type
type RGB struct {
	R, G, B float64
}

/*
accentColorTimeout bounds the portal read.

When xdg-desktop-portal cannot start (root shells, TTY logins, a session that has not reached graphical-session.target), the D-Bus activation blocks for its full 120-second timeout instead of failing fast; a greeting must never hold a login hostage that long.
*/
const accentColorTimeout = 2 * time.Second

/*
GetRender takes a color and input string, and returns the rendered output based on the specified color scheme.

- If the color is set to "auto", it detects the accent color from the system settings.

- If the color is one of the predefined colors, it uses that color for rendering.

- If the color is not recognized, it defaults to a non-colored rendering.
*/
func GetRender(color string, input string, l *gotext.Locale) string {
	var withColor bool = true
	var processColor, out string
	var err error
	var r *glamour.TermRenderer

	// Determine the color scheme based on the provided color argument
	switch color {
	case "auto":
		processColor = getAccentColor()
	case "blue", "green", "orange", "pink", "purple", "red", "slate", "teal", "yellow":
		processColor = colorMap[color]
	default:
		withColor = false
	}

	// Render the input string using the appropriate color scheme
	if !withColor {
		colorScheme := "light"
		if lipgloss.HasDarkBackground(os.Stdin, os.Stdout) {
			colorScheme = "dark"
		}
		out, err = glamour.Render(input, colorScheme)
		if err != nil {
			i.Error(l.Get("An error occured while rendering the banner."), err)
			return ""
		}
	} else {
		r, err = glamour.NewTermRenderer(
			glamour.WithStyles(getColorizedStyle(processColor)),
		)
		if err != nil {
			i.Error(l.Get("An error occured while loading the renderer."), err)
			r.Close()
			return ""
		}
		out, err = r.Render(input)
		if err != nil {
			i.Error(l.Get("An error occured while rendering the banner."), err)
			r.Close()
			return ""
		}
		r.Close()
	}

	return out
}

// xyzToRgb converts a DBus (ddd) type to an RGB struct. The input is expected to be a slice of three float64 values representing the RGB components in the range [0, 1]. The function scales these values to the range [0, 255] and returns an RGB struct.
func xyzToRgb(xyz any) (RGB, error) {
	// Type assert the input as []interface{}
	rgbSlice, ok := xyz.([]interface{})
	if !ok {
		return RGB{R: 0, G: 0, B: 0}, fmt.Errorf("input must be a []interface{} slice of RGB values (0-1)")
	}
	if len(rgbSlice) != 3 {
		return RGB{R: 0, G: 0, B: 0}, fmt.Errorf("input must contain exactly 3 RGB values")
	}

	// Extract R, G, B as float64
	r, okR := rgbSlice[0].(float64)
	g, okG := rgbSlice[1].(float64)
	b, okB := rgbSlice[2].(float64)
	if !okR || !okG || !okB {
		return RGB{R: 0, G: 0, B: 0}, fmt.Errorf("R, G, B values must be float64")
	}

	// Scale to 0-255 (convert from 0-1 range to 0-255)
	scale := func(v float64) int {
		return int(math.Round(v * 255))
	}
	ri, gi, bi := scale(r), scale(g), scale(b)
	return RGB{R: float64(ri), G: float64(gi), B: float64(bi)}, nil
}

/*
rgbToHex converts an RGB struct to a hexadecimal color string. The input is expected to be an RGB struct with R, G, B values in the range [0, 255].

The function returns a hex string in the format "#RRGGBB". If the RGB values are out of range, it returns an error.
*/
func rgbToHex(rgb RGB) (string, error) {
	// Ensure RGB values are within the valid range
	if rgb.R < 0 || rgb.R > 255 || rgb.G < 0 || rgb.G > 255 || rgb.B < 0 || rgb.B > 255 {
		return "", fmt.Errorf("RGB values must be in the range [0, 255]")
	}

	// Convert RGB to hex string
	hex := fmt.Sprintf("#%02X%02X%02X", int(rgb.R), int(rgb.G), int(rgb.B))
	return hex, nil
}

/*
getAccentColor retrieves the accent color from the system settings using D-Bus. It reads the "accent-color" property from the "org.freedesktop.appearance" interface.

If successful, it converts the retrieved XYZ color to RGB and then to a hexadecimal string.

If any step fails or the portal does not answer within the accentColorTimeout, it defaults to returning "blue".
*/
func getAccentColor() string {

	// Read the accent color from D-Bus, bounded by accentColorTimeout.
	// The goroutine is leaked if the portal never answers; the process
	// is short-lived so that is acceptable.
	type readResult struct {
		value any
		err   error
	}
	resultCh := make(chan readResult, 1)
	go func() {
		value, err := settings.ReadOne("org.freedesktop.appearance", "accent-color")
		resultCh <- readResult{value, err}
	}()

	var xyzValue any
	select {
	case result := <-resultCh:
		if result.err != nil {
			return "blue"
		}
		xyzValue = result.value
	case <-time.After(accentColorTimeout):
		return "blue"
	}

	// Convert the XYZ value to RGB
	rgbValue, err := xyzToRgb(xyzValue)
	if err != nil {
		return "blue"
	}

	// Convert the RGB value to Hex
	hexValue, err := rgbToHex(rgbValue)
	if err != nil {
		return "blue"
	}

	return hexValue
}
