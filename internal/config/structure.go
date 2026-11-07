package config

import (
	"github.com/leonelquinteros/gotext"
)

const Version = 1

type ConfigCheck struct {
	Version int `toml:"version"`
}

type Config struct {
	Version  int       `toml:"version"`
	Color    string    `toml:"color"`
	Greeting Greeting  `toml:"greeting"`
	Commands []Command `toml:"commands"`
	Links    []Link    `toml:"links"`
	Motd     Motd      `toml:"motd"`
}

type Command struct {
	Cmd  string `toml:"cmd"`
	Desc string `toml:"desc"`
}

type Greeting struct {
	Prefix  string `toml:"prefix"`
	Suffix  string `toml:"suffix"`
	Message string `toml:"message"`
}

type Link struct {
	Name string `toml:"name"`
	URL  string `toml:"url"`
}

type Motd struct {
	Messages []string `toml:"messages"`
	Commands []string `toml:"commands"`
}

// Default returns a default config
func Default(l *gotext.Locale) Config {
	return Config{
		Version: Version,
		Commands: []Command{
			{Cmd: "uwelcome toggle", Desc: "banner_toggle"},
			{Cmd: "fastfetch", Desc: "sys_info"},
			{Cmd: "brew help", Desc: "cli_pkg"},
		},
		Links: []Link{
			{Name: "source", URL: "https://github.com/projectbluefin/uwelcome/"},
			{Name: "issues", URL: "https://github.com/projectbluefin/uwelcome/issues"},
			{Name: "docs", URL: "https://github.com/projectbluefin/uwelcome/tree/main/docs"},
		},
		Motd: Motd{
			Messages: []string{
				l.Get("uWelcome is made with ❤️ by Universal Blue."),
				l.Get("uWelcome is free and open source! [Check out the source code on GitHub](%s)", "https://github.com/projectbluefin/uwelcome/"),
				l.Get("uWelcome is a terminal welcome banner written in Go."),
				l.Get("uWelcome is translatable! You can help translate it in your language. Read the [translation guide](%s)", "https://github.com/projectbluefin/uwelcome/blob/main/docs/translating.md"),
				l.Get("You're currently using the default config. You can customize it by editing the config file with `uwelcome edit`"),
			},
		},
		Greeting: Greeting{
			Suffix: "!",
		},
	}
}
