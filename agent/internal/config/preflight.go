package config

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tale/headplane/agent/internal/i18n"
)

// Checks to make sure all required environment variables are set
func validateRequired(config *Config) error {
	if config.Hostname == "" {
		return fmt.Errorf("%s", i18n.Message(
			"agent.config.required_env",
			"{{.Name}} is required",
			map[string]any{"Name": HostnameEnv},
		))
	}

	if config.TSControlURL == "" {
		return fmt.Errorf("%s", i18n.Message(
			"agent.config.required_env",
			"{{.Name}} is required",
			map[string]any{"Name": TSControlURLEnv},
		))
	}

	if config.TSAuthKey == "" {
		return fmt.Errorf("%s", i18n.Message(
			"agent.config.required_env",
			"{{.Name}} is required",
			map[string]any{"Name": TSAuthKeyEnv},
		))
	}

	if config.WorkDir == "" {
		return fmt.Errorf("%s", i18n.Message(
			"agent.config.required_env",
			"{{.Name}} is required",
			map[string]any{"Name": WorkDirEnv},
		))
	}

	return nil
}

// Pings the Tailscale control server to make sure it's up and running
func validateTSReady(config *Config) error {
	testURL := config.TSControlURL
	if strings.HasSuffix(testURL, "/") {
		testURL = testURL[:len(testURL)-1]
	}

	testURL = fmt.Sprintf("%s/health", testURL)
	resp, err := http.Get(testURL)
	if err != nil {
		return fmt.Errorf("%s", i18n.Message(
			"agent.config.ts_unreachable",
			"Failed to connect to TS control server: {{.Reason}}",
			map[string]any{"Reason": err},
		))
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", i18n.Message(
			"agent.config.ts_unreachable",
			"Failed to connect to TS control server: {{.Reason}}",
			map[string]any{"Reason": resp.Status},
		))
	}

	return nil
}
