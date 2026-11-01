package config

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"uwelcome/internal"
	"uwelcome/internal/symbols"

	"github.com/BurntSushi/toml"
	"github.com/leonelquinteros/gotext"
)

const etcPath = "/etc/uwelcome/config.toml"

var homePath = getHomePath()
var isRoot = (os.Geteuid() == 0)

// getHomePath returns the path of the config if it's in the user's Home directory
func getHomePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "uwelcome", "config.toml")
}

// createDefault create the default config to the given path
func createDefault(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := toml.Marshal(defaultConfig())
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func createEtcConfig() error {
	return createDefault(etcPath)
}

func createHomeConfig() error {
	if exists, _ := internal.DoesFileExist(etcPath); exists && check(etcPath) == nil {
		return copyFile(etcPath, homePath)
	}

	return createDefault(homePath)
}

// copyFile copies a file from a source path to the destination path
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// check returns any problem from the given config file
func check(path string) error {
	_, err := loadConfig(path)
	return err
}

// getConfigVersion returns the config version & any problems related to its retrieval
func getConfigVersion(path string) (int, error) {
	var fileVersion ConfigCheck

	_, decodeErr := toml.DecodeFile(path, &fileVersion)
	if decodeErr != nil {
		return 0, decodeErr
	}

	return fileVersion.Version, nil
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

// loadConfig loads the config at the given path or returns an error
func loadConfig(path string) (Config, error) {
	var c Config
	_, decodeErr := toml.DecodeFile(path, &c)
	return c, decodeErr
}

// GetPath returns the path to a valid config file or a blank string if none is found
func GetPath() string {
	_, cfgPath := GetConfig()
	return cfgPath
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
