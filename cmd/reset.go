package cmd

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"uwelcome/internal"
	"uwelcome/internal/config"

	"github.com/leonelquinteros/gotext"
)

// Reset resets the config to the system/default config
func Reset(l *gotext.Locale) {
	configPath, create := config.Target(l)
	parsedResponse := ""

	fmt.Println(l.Get("Are you sure? This will reset your config file (%s) (yes/NO)", configPath))
	r := bufio.NewReader(os.Stdin)

	response, err := r.ReadString('\n')
	if err != nil {
		internal.Error(l.Get("There was an error parsing your answer."), err)
		os.Exit(1)
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
		if err := create(l); err != nil {
			internal.Error(l.Get("Failed to reset the config file."), err)
			os.Exit(1)
		} else {
			fmt.Println(l.Get("Config file reset successfully."))
		}
	} else {
		fmt.Println(l.Get("Config file untouched."))
	}
}
