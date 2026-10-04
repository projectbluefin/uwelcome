package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	i "uwelcome/internal"

	"github.com/leonelquinteros/gotext"
)

const etcPath = "/etc/uwelcome/config.json"

var homePath = getHomePath()
var isRoot = (os.Geteuid() == 0)

// getHomePath returns the path of the config if it's in the user's Home directory
func getHomePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "uwelcome", "config.json")
}

// createDefault create the default config to the given path
func createDefault(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(defaultConfig(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func createEtcConfig() error {
	return createDefault(etcPath)
}

func createHomeConfig() error {
	if exists, _ := i.DoesFileExist(etcPath); exists && check(etcPath) == nil {
		return copyFile(etcPath, homePath)
	} else {
		return createDefault(homePath)
	}
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

// CheckWrapper uses check() on a path and prints the status of the config of that path
func CheckWrapper(l *gotext.Locale) {
	path, _ := target()
	fmt.Println(path)
	if err := check(path); err == nil {
		i.Success(l.Get("The config file is good to go."))
	} else if errors.Is(err, fs.ErrNotExist) {
		i.Warn(l.Get("The config file doesn't exist (yet). uWelcome will use either your system's config or the default."), nil)
	} else {
		i.Warn(l.Get("There is something wrong with the config file. uWelcome will use either your system's config or the default."), err)
	}
}

// loadConfig loads the config at the given path or returns an error
func loadConfig(path string) (Config, error) {
	var c Config
	data, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	return c, json.Unmarshal(data, &c)
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

// target detects if the current user is root and then returns the config path and the function to create the config file accordingly
func target() (path string, create func() error) {
	if isRoot {
		return etcPath, createEtcConfig
	}
	return homePath, createHomeConfig
}

// Edit detects if the user has admin privileges and lets them edit the corresponding configuration file
func Edit(l *gotext.Locale) {
	configPath, create := target()

	if exists, err := i.DoesFileExist(configPath); exists && err != nil {
		i.Error(l.Get("Failed to access the config path."), err)
		return
	} else if !exists && err != nil {
		if err := create(); err != nil {
			i.Error(l.Get("Failed to create the default config."), err)
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
		i.Error(l.Get("Failed to edit the config file."), err)
		return
	}

	CheckWrapper(l)
}

// Reset resets the config to the system/default configuration
func Reset(l *gotext.Locale) {
	configPath, create := target()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(l.Get("Are you sure? This will reset your configuration file (%s) (yes/NO)", configPath))
	output, err := reader.ReadString('\n')
	if err != nil {
		i.Error(l.Get("There was an error in processing your answer"), err)
		return
	}
	output = strings.ToLower(output)
	output = strings.TrimSpace(output)
	tmpTable := strings.Fields(output)
	if len(tmpTable) > 0 {
		output = tmpTable[0]
	} else {
		output = ""
	}

	responses := []string{}

	responses = append(responses, []string{
		l.GetC("yes", "Single word affirmative response in response to a question (lower case only)"),
		l.GetC("y", "First letter of the single word affirmative response in response to a question (lower case only)"),
		"yes", "y"}...)

	if slices.Contains(responses, output) {
		if err := create(); err != nil {
			i.Error(l.Get("Failed to reset the configuration file."), err)
		} else {
			i.Success(l.Get("Config file reset successfully."))
		}
	}
}
