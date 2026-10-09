package internal

import (
	"fmt"

	"github.com/leonelquinteros/gotext"
)

const space = "          "

// Usage
func Usage(l *gotext.Locale) {
	fmt.Println(l.Get("uWelcome is a customizable terminal banner for Linux.") + "\n")
	fmt.Println(l.Get("Usage:") + "\n")
	println(space + "uwelcome")
	println(space + "uwelcome <command>\n")

	fmt.Println(l.Get("The commands are:") + "\n")
	println(space + "disable" + "    " + l.Get("disable the terminal banner"))
	println(space + "edit" + "       " + l.Get("open your uwelcome config file with your default `$EDITOR`"))
	println(space + "enable" + "     " + l.Get("enable the terminal banner"))
	println(space + "help" + "       " + l.Get("show this menu"))
	println(space + "reset" + "      " + l.Get("reset your uwelcome config file to the system default"))
	println(space + "status" + "     " + l.Get("show the location and status of your uwelcome config"))
	println(space + "toggle" + "     " + l.Get("toggle the terminal banner"))
	println(space + "version" + "    " + l.Get("print uwelcome version"))
}
