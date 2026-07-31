package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

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
			{Name: "discuss", URL: "https://universal-blue.discourse.group/"},
			{Name: "discord", URL: "https://discord.com/invite/8RZGC3uFzA"},
			{Name: "mastodon", URL: "https://fosstodon.org/@UniversalBlue"},
		},
		Greeting: Greeting{
			Suffix: "!",
		},
	}
}

func WriteDefaultConfig(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(defaultConfig(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddMotdMessage(msg string) error {
	cfg := GetConfig()
	cfg.Motd.Messages = append(cfg.Motd.Messages, msg)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(GetPath(), data, 0644)
}

func RemoveMotdMessage(msg string) error {
	// if strings.Contains(GetPath(), "/etc/") {
	// 	return nil
	// }
	cfg := GetConfig()
	for i, preset := range cfg.Motd.Messages {
		if preset == msg {
			cfg.Motd.Messages = append(cfg.Motd.Messages[:i], cfg.Motd.Messages[i+1:]...)
			break
		}
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(GetPath(), data, 0644)
}

func ListMotdMessages() []string {
	cfg := GetConfig()
	return cfg.Motd.Messages
}

// isConfigOkay checks if the config file at the given path is valid
func isConfigOkay(path string) bool {

	var noError bool = true

	_, err := os.Stat(path)
	if err != nil {
		noError = false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		noError = false
	}

	var tempCfg Config
	if err := json.Unmarshal(data, &tempCfg); err != nil {
		noError = false
	}

	return noError
}

// GetPath returns the path to a valid config file, returns "" if no valid config file is found
func GetPath() string {
	path := []string{
		os.ExpandEnv("$HOME/.config/uwelcome/config.json"),
		"/etc/uwelcome/config.json",
	}
	for _, path := range path {
		if isConfigOkay(path) {
			return path
		}
	}
	return ""
}

// GetConfig returns the config file at the given path, or a default config if no valid config file is found
func GetConfig() Config {
	cfg := Config{}
	path := GetPath()

	if path == "" {
		return defaultConfig()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig()
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig()
	}

	return cfg
}
