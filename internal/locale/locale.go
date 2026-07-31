package locale

import (
	"embed"
	"os"
	"strings"

	"golang.org/x/text/language"
)

// DetectLocale detects the current language
func DetectLocale(localesFS embed.FS) string {
	langDir, _ := localesFS.ReadDir("locales")

	tags := []language.Tag{language.English}
	for _, languageDir := range langDir {
		if languageDir.Name() == "en" {
			continue
		}
		tags = append(tags, language.Make(languageDir.Name()))
	}

	var supported = language.NewMatcher(tags)

	raw := os.Getenv("LC_ALL")
	if raw == "" {
		raw = os.Getenv("LANGUAGE")
	}
	if raw == "" {
		raw = os.Getenv("LANG")
	}
	raw = strings.Split(raw, ".")[0]
	raw = strings.ReplaceAll(raw, "_", "-")

	tag := language.Make(raw)
	match, _, _ := supported.Match(tag)

	return match.String()
}
