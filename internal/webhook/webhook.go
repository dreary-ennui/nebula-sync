package webhook

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lovelaze/nebula-sync/internal/config"
	"github.com/lovelaze/nebula-sync/version"
	"github.com/rs/zerolog/log"
)

type Client struct {
	success    config.WebhookRequest
	failure    config.WebhookRequest
	httpClient *http.Client
}

func NewClient(c *config.WebhookSettings) *Client {
	return &Client{
		success: c.Success,
		failure: c.Failure,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: c.Client.SkipTLSVerification},
			}},
	}
}

func (c *Client) OnSuccess() {
	if err := c.triggerSuccess(); err != nil {
		log.Warn().Err(err).Msg("Webhook trigger failed")
	}
}

func (c *Client) OnFailure(err error) {
	if err := c.triggerFailure(); err != nil {
		log.Warn().Err(err).Msg("Webhook trigger failed")
	}
}

func (c *Client) triggerSuccess() error {
	return invoke(c.httpClient, c.success)
}

func (c *Client) triggerFailure() error {
	return invoke(c.httpClient, c.failure)
}

func invoke(client *http.Client, settings config.WebhookRequest) error {
	if settings.Url == "" {
		return nil
	}

	log.Debug().
		Str("url", settings.Url).
		Str("method", settings.Method).
		Str("body", settings.Body).
		Interface("headers", settings.Headers).
		Msg("Invoking webhook")

	req, err := http.NewRequest(settings.Method, settings.Url, strings.NewReader(settings.Body))
	if err != nil {
		return fmt.Errorf("create webhook request: %w", err)
	}

	req.Header.Set("User-Agent", fmt.Sprintf("nebula-sync/%s", version.Version))

	for key, value := range settings.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
