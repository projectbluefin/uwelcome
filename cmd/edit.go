package cmd

import (
	"os"
	"os/exec"
	"strings"
	"uwelcome/internal"
	"uwelcome/internal/config"

	"github.com/leonelquinteros/gotext"
)

// Edit detects if the user has admin privileges and lets them edit the corresponding configuration file
func Edit(l *gotext.Locale) {
	configPath, create := config.Target()

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

	Status(l)
}
