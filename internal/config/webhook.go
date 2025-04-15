package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type WebhookSettings struct {
	Failure WebhookRequest `ignored:"true"`
	Success WebhookRequest `ignored:"true"`
	Client  WebhookClient  `ignored:"true"`
}

type WebhookClient struct {
	SkipTLSVerification bool `default:"false" envconfig:"SKIP_TLS_VERIFICATION"`
}

type WebhookRequest struct {
	Body    string            `envconfig:"BODY"`
	Headers map[string]string `envconfig:"HEADERS"`
	Method  string            `default:"POST" envconfig:"METHOD"`
	Url     string            `envconfig:"URL"`
}

const webhookEnvPrefix = "SYNC_WEBHOOK_"

func (c *Config) loadWebhookSettings() error {
	webhookSettings := WebhookSettings{}

	if err := envconfig.Process(webhookEnvPrefix+"FAILURE", &webhookSettings.Failure); err != nil {
		return fmt.Errorf("process webhook env vars for failure: %w", err)
	}
	if err := envconfig.Process(webhookEnvPrefix+"SUCCESS", &webhookSettings.Success); err != nil {
		return fmt.Errorf("process webhook env vars for success: %w", err)
	}
	if err := envconfig.Process(webhookEnvPrefix+"CLIENT", &webhookSettings.Client); err != nil {
		return fmt.Errorf("process webhook env vars for client: %w", err)
	}

	c.Sync.WebhookSettings = &webhookSettings

	return nil
}
