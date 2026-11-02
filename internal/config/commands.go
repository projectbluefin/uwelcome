package config

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"slices"
	"strings"
	"uwelcome/internal"
	"uwelcome/internal/symbols"

	"github.com/leonelquinteros/gotext"
)

const etcPath = "/etc/uwelcome/config.toml"

var homePath = getHomePath()
var isRoot = (os.Geteuid() == 0)

// GetPath returns the path to a valid config file or a blank string if none is found
func GetPath() string {
	_, cfgPath := GetConfig()
	return cfgPath
}

func IsVersionSupported(path string) (bool, error) {
	version, err := getConfigVersion(path)
	return ConfigVersion == version, err
}

// CheckWrapper uses check() on a path and prints the status of the config of that path
func CheckWrapper(l *gotext.Locale) {
	path, _ := Target()
	fmt.Println(path)

	if err := check(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			internal.Warn(l.Get("The config file doesn't exist (yet). uWelcome will use either your system's config or the default."), nil)
			return
		}

		internal.Warn(l.Get("There is something wrong with the config file. uWelcome will use either your system's config or the default."), err)
		return
	}

	isSupported, _ := IsVersionSupported(path)

	if !isSupported {
		internal.Warn(l.Get("Your config file's version isn't supported by this version of uwelcome. This version supports config version %d.", ConfigVersion), nil)
	}

	fmt.Println(l.Get("The config file is OK.") + " " + symbols.GetSymbol("checkmark"))
}

// GetConfig returns the config file at the first found path, or a default config if no valid config file is found
func GetConfig() (Config, string) {
	possiblePaths := []string{
		homePath,
		etcPath,
	}

	for _, p := range possiblePaths {
		if cfg, err := loadConfig(p); err == nil {
			return cfg, p
		}
	}

	return defaultConfig(), ""
}

// Target detects if the current user is root and then returns the config path and the function to create the config file accordingly
func Target() (path string, create func() error) {
	if isRoot {
		return etcPath, createEtcConfig
	}
	return homePath, createHomeConfig
}

// Edit detects if the user has admin privileges and lets them edit the corresponding configuration file
func Edit(l *gotext.Locale) {
	configPath, create := Target()

	if exists, err := internal.DoesFileExist(configPath); exists && err != nil {
		internal.Error(l.Get("Failed to access the config path."), err)
		return
	} else if !exists {
		if createErr := create(); createErr != nil {
			internal.Error(l.Get("Failed to create the default config."), createErr)
			return
		}
	}

	parts := strings.Fields(os.Getenv("EDITOR"))
	if len(parts) == 0 {
		parts = []string{"nano"}
	}
	cmdEdit := exec.Command(parts[0], append(parts[1:], configPath)...)

	// Giving access to terminal input, output and error to the editor
	cmdEdit.Stdin = os.Stdin
	cmdEdit.Stdout = os.Stdout
	cmdEdit.Stderr = os.Stderr

	if err := cmdEdit.Run(); err != nil {
		internal.Error(l.Get("Failed to edit the config file."), err)
		return
	}

	CheckWrapper(l)
}

// Reset resets the config to the system/default config
func Reset(l *gotext.Locale) {
	configPath, create := Target()
	parsedResponse := ""

	fmt.Println(l.Get("Are you sure? This will reset your config file (%s) (yes/NO)", configPath))
	r := bufio.NewReader(os.Stdin)

	response, err := r.ReadString('\n')
	if err != nil {
		internal.Error(l.Get("There was an error parsing your answer."), err)
		return
	}

	tmpTable := strings.Fields(response)
	if len(tmpTable) > 0 {
		parsedResponse = tmpTable[0]
	}

	parsedResponse = strings.ToLower(parsedResponse)

	if slices.Contains([]string{
		l.GetC("yes", "Word used to give an affirmative response (lower case only)"),
		l.GetC("y", "First letter of the word used to give an affirmative response (lower case only)"),
		"yes", "y",
	}, parsedResponse) {
		if err := create(); err != nil {
			internal.Error(l.Get("Failed to reset the config file."), err)
		} else {
			fmt.Println(l.Get("Config file reset successfully."))
		}
	} else {
		fmt.Println(l.Get("Config file untouched."))
	}
}
