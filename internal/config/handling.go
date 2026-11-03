package config

import (
	"errors"
	"os"
	"path/filepath"

	"uwelcome/internal"

	"github.com/BurntSushi/toml"
	"github.com/leonelquinteros/gotext"
)

const etcPath = "/etc/uwelcome/config.toml"

var ErrUnsupportedConfig = errors.New("Config file version not supported.")
var ErrUndecoded = errors.New("Unknown keys present in config file.")

// getHomePath returns the path of the config if it's in the user's Home directory
func getHomePath() (string, error) {
	dir, err := internal.UserDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

// createDefault create the default config to the given path
func createDefault(path string, l *gotext.Locale) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := toml.Marshal(Default(l))
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func CreateEtcConfig(l *gotext.Locale) error {
	return createDefault(etcPath, l)
}

func CreateHomeConfig(l *gotext.Locale) error {
	_, err := Load(etcPath, l)

	homePath, homeErr := getHomePath()
	if homeErr != nil {
		return homeErr
	}

	if exists, _ := internal.DoesFileExist(etcPath); exists && err == nil {
		return copyFile(etcPath, homePath)
	}

	return createDefault(homePath, l)
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

// Load loads the config at the given path or returns an error
func Load(path string, l *gotext.Locale) (Config, error) {
	var fileVersion ConfigCheck
	c := Default(l)

	_, versionErr := toml.DecodeFile(path, &fileVersion)
	if versionErr != nil {
		return c, versionErr
	}
	if fileVersion.Version != Version {
		return c, ErrUnsupportedConfig
	}

	md, decodeErr := toml.DecodeFile(path, &c)
	if decodeErr != nil {
		return c, decodeErr
	}

	if len(md.Undecoded()) > 0 {
		return c, ErrUndecoded
	}

	return c, nil
}

// Target detects if the current user is root and then returns the config path and the function to create the config file accordingly
func Target(l *gotext.Locale) (path string, create func(l *gotext.Locale) error) {
	if os.Geteuid() == 0 {
		return etcPath, CreateEtcConfig
	}
	homePath, _ := getHomePath()
	return homePath, CreateHomeConfig
}

// Get returns the config file at the first found path, or a default config if no valid config file is found
func Get(l *gotext.Locale) Config {
	homePath, homeErr := getHomePath()
	if homeErr != nil {
		return Default(l)
	}
	for _, p := range []string{homePath, etcPath} {
		c, err := Load(p, l)
		if err == nil || errors.Is(err, ErrUndecoded) {
			return c
		}
	}
	return Default(l)
}
