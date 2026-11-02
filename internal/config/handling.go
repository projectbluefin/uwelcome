package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"uwelcome/internal"

	"github.com/BurntSushi/toml"
)

const etcPath = "/etc/uwelcome/config.toml"

var homePath = getHomePath()
var ErrUnsupportedConfig = errors.New("Config file version not supported.")
var ErrUndecoded = errors.New("Unknown keys present in config file.")

// getHomePath returns the path of the config if it's in the user's Home directory
func getHomePath() string {
	dir, _ := os.UserConfigDir()
	if dir == "" {
		return dir
	}
	return filepath.Join(dir, "uwelcome", "config.toml")
}

// createDefault create the default config to the given path
func createDefault(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := toml.Marshal(Default())
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func CreateEtcConfig() error {
	return createDefault(etcPath)
}

func CreateHomeConfig() error {
	_, err := Load(etcPath)

	if exists, _ := internal.DoesFileExist(etcPath); exists && err == nil {
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

// loadConfig loads the config at the given path or returns an error
func Load(path string) (Config, error) {
	var fileVersion ConfigCheck
	var c Config = Default()

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
		// strings.Join(md.Undecoded(), " ")
		return c, ErrUndecoded
	}

	return c, nil
}

// Target detects if the current user is root and then returns the config path and the function to create the config file accordingly
func Target() (path string, create func() error) {
	if os.Geteuid() == 0 {
		return etcPath, CreateEtcConfig
	}
	return homePath, CreateHomeConfig
}

// Get returns the config file at the first found path, or a default config if no valid config file is found
func Get() (cfg Config, path string, warn error) {
	for _, p := range []string{homePath, etcPath} {
		c, err := Load(p)
		if err == nil {
			return c, p, warn
		}
		if !errors.Is(err, fs.ErrNotExist) && warn == nil {
			warn = fmt.Errorf("%s: %w", p, err)
		}
	}
	return Default(), "", warn
}
