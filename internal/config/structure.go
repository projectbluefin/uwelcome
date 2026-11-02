package config

import "uwelcome/internal"

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

// defaultConfig returns a sensible default config
func Default() Config {
	return Config{
		Version: Version,
		Commands: []Command{
			{Cmd: "uwelcome toggle", Desc: "banner_toggle"},
			{Cmd: "fastfetch", Desc: "sys_info"},
			{Cmd: "brew help", Desc: "cli_pkg"},
		},
		Links: []Link{
			{Name: "source", URL: "https://github.com/ublue-os/uwelcome/"},
			{Name: "issues", URL: "https://github.com/ublue-os/uwelcome/issues"},
			{Name: "docs", URL: "https://github.com/ublue-os/uwelcome/tree/main/docs"},
		},
		Motd: Motd{
			Messages: []string{
				internal.Default1,
				internal.Default2,
				internal.Default3,
				internal.Default4,
				internal.Default5,
			},
		},
		Greeting: Greeting{
			Suffix: "!",
		},
	}
}
