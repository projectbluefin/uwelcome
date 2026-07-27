package motd

import (
	"math/rand"
	"os/exec"
	"uwelcome/internal/config"
)

// GetRandomMessage returns a random string or the output of a random command
func GetRandomMessage(cfg config.Config) string {

	var messages []string

	if len(cfg.Motd.Messages) > 0 {
		for _, msg := range cfg.Motd.Messages {
			messages = append(messages, msg)
		}
	}

	if len(cfg.Motd.Commands) > 0 {
		for _, msg := range cfg.Motd.Commands {
			out, _ := exec.Command(msg).Output()
			messages = append(messages, string(out))
		}
	}

	return messages[rand.Intn(len(messages))]
}
