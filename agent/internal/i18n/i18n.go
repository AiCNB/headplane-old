package i18n

import (
	"embed"
	"encoding/json"
	"os"
	"strings"
	"sync"

	gi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

var (
	bundleOnce sync.Once
	bundle     *gi18n.Bundle
	localizer  *gi18n.Localizer
)

func initBundle() {
	bundle = gi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files := []string{"en.json", "zh.json"}
	for _, file := range files {
		data, err := localeFS.ReadFile("locales/" + file)
		if err != nil {
			panic(err)
		}

		if _, err = bundle.ParseMessageFileBytes(data, file); err != nil {
			panic(err)
		}
	}
}

func initLocalizer() {
	bundleOnce.Do(initBundle)
	tag := detectLanguage()
	localizer = gi18n.NewLocalizer(bundle, tag, "en")
}

func getLocalizer() *gi18n.Localizer {
	if localizer == nil {
		initLocalizer()
	}

	return localizer
}

func detectLanguage() string {
	if value := os.Getenv("HEADPLANE_AGENT_LANG"); value != "" {
		return value
	}

	if value := os.Getenv("LANG"); value != "" {
		normalized := strings.Split(value, ".")[0]
		normalized = strings.Split(normalized, "_")[0]
		normalized = strings.Split(normalized, "-")[0]
		if normalized != "" {
			return normalized
		}
	}

	return "en"
}

// Message returns a localized string using the provided message id and default translation.
func Message(id string, defaultMessage string, data map[string]any) string {
	localizer := getLocalizer()
	message, err := localizer.Localize(&gi18n.LocalizeConfig{
		MessageID: id,
		DefaultMessage: &gi18n.Message{
			ID:    id,
			Other: defaultMessage,
		},
		TemplateData: data,
	})
	if err != nil {
		return defaultMessage
	}

	return message
}
