package main

import (
	"fmt"
	"net/http"

	"github.com/sensu/sensu-go/types"
	sensuhttp "github.com/sensu/sensu-plugins-go-library/http"
	"github.com/sensu/sensu-plugins-go-library/sensu"
)

type HandlerConfig struct {
	sensu.PluginConfig

	url     string
	timeout int
}

type Payload struct {
	Source    string
	Component string
	Severity  string
	Summary   string
	Details   interface{}
}

var (
	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:  "sensu-handler-template",
			Short: "A Sensu Go handler template to use as a starting point for a Sensu Go handler",
		},
	}

	handlerConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      "url",
			Env:       "URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "https://staging.fakedomain.local",
			Usage:     "The URL of the API to send the event to",
			Value:     &config.url,
		},
		{
			Path:      "timeout",
			Env:       "TIMEOUT",
			Argument:  "timeout",
			Shorthand: "t",
			Default:   10,
			Usage:     "The duration (in seconds) to wait before timing out",
			Value:     &config.timeout,
		},
	}
)

func main() {
	goHandler := sensu.NewGoHandler(&config.PluginConfig, handlerConfigOptions, checkArgs, executeHandler)
	goHandler.Execute()
}

func checkArgs(_ *types.Event) error {
	if len(config.url) == 0 {
		return fmt.Errorf("--url or URL environment variable is required")
	}

	return nil
}

func executeHandler(event *types.Event) error {
	httpWrapper, err := sensuhttp.NewHttpWrapper(uint64(config.Timeout), "", "", "")
	if err != nil {
		return fmt.Errorf("could not create http wrapper: %s", err.Error())
	}

	body := createPayload(event)

	statusCode, httpResponse, err := httpWrapper.ExecuteRequest(http.MethodPost, config.url, body, nil)

	if err != nil {
		return err
	}

	if statusCode >= 400 {
		return fmt.Errorf("failed to send event: %s - %s", statusCode, httpResponse)
	}

	return nil
}

func createPayload(event *types.Event) Payload {
	return Payload{
		Source:    event.Entity.Name,
		Component: event.Check.Name,
		Severity:  getSeverity(event.Check.Status),
		Summary:   event.Check.Output,
		Details:   event,
	}
}

func getSeverity(status uint32) string {
	severity := "unknown"
	if status < 3 {
		severities := []string{"info", "warning", "critical"}
		severity = severities[status]
	}
	return severity
}
