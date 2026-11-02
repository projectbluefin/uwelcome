package config

import (
	"os"
	"path/filepath"

	"uwelcome/internal"

	"github.com/BurntSushi/toml"
)

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

// loadConfig loads the config at the given path or returns an error
func loadConfig(path string) (Config, error) {
	var c Config
	_, decodeErr := toml.DecodeFile(path, &c)
	return c, decodeErr
}
