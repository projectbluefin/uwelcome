package motd

import (
	"math/rand"
	"os/exec"
	"uwelcome/internal/config"
)

// GetRandomMotd returns a random string or the output of a random command
func GetRandomMotd(cfg config.Config) string {

	var sizeMessages int = len(cfg.Motd.Messages)
	var sizeCommands int = len(cfg.Motd.Commands)

	switch {
	case sizeMessages > 0 && sizeCommands > 0:
		if pick := rand.Intn(2); pick == 0 {
			return getRandomMessage(cfg, sizeMessages)
		} else {
			return getRandomCommand(cfg, sizeCommands)
		}
	case sizeCommands > 0:
		return getRandomCommand(cfg, sizeCommands)
	case sizeMessages > 0:
		return getRandomMessage(cfg, sizeMessages)
	default:
		return ""
	}
}

func getRandomMessage(cfg config.Config, size int) string {
	return cfg.Motd.Messages[rand.Intn(size)]
}

func getRandomCommand(cfg config.Config, size int) string {
	out, err := exec.Command(cfg.Motd.Commands[rand.Intn(size)]).Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}
