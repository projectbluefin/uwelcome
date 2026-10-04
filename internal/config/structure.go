package config

type Config struct {
	Color    string    `json:"color"`
	Commands []Command `json:"commands"`
	Greeting Greeting  `json:"greeting"`
	Links    []Link    `json:"links"`
	Motd     Motd      `json:"motd"`
}

type Command struct {
	Cmd  string `json:"cmd"`
	Desc string `json:"desc"`
}

type Greeting struct {
	Prefix  string `json:"prefix"`
	Suffix  string `json:"suffix"`
	Message string `json:"message"`
}

type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Motd struct {
	Messages []string `json:"messages"`
	Commands []string `json:"commands"`
}

// defaultConfig returns a sensible default config
func defaultConfig() Config {
	return Config{
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
				"uWelcome is made with ❤️  by Project Bluefin.",
				"uWelcome is free and open source! [Check out the source code on GitHub](https://github.com/projectbluefin/uwelcome/)",
				"uWelcome is a terminal welcome banner written in Go.",
				"uWelcome is translatable! You can help translate it in your language. Read the [translation guide](https://github.com/projectbluefin/uwelcome/blob/main/TRANSLATION.md)",
				"You're currently using the default config. You can customize it by editing the config file with `uwelcome edit`",
			},
		},
		Greeting: Greeting{
			Suffix: "!",
		},
	}
}
