package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

var settings Config

type Config struct {
	PaymentCallbackUrl string `envconfig:"PAYMENT_CALLBACK_URL" required:"true"`
	DodoWebhookSecret  string `envconfig:"DODO_WEBHOOK_SECRET" required:"true"`
	DodoAPIKey         string `envconfig:"DODO_API_KEY" required:"true"`
	DodoCheckoutURL    string `envconfig:"DODO_CHECKOUT_URL" required:"true"`
	KVNamespace        string `envconfig:"KV_NAMESPACE" required:"true"`
}

func LoadConfig() error {
	if err := envconfig.Process("", &settings); err != nil {
		return fmt.Errorf("failed to load config from environment: %w", err)
	}
	if err := settings.validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

func GetConfig() Config {
	return settings
}

func (c Config) validate() error {
	if !isValidURL(c.PaymentCallbackUrl) {
		return fmt.Errorf("PAYMENT_CALLBACK_URL is invalid: %s", c.PaymentCallbackUrl)
	}
	if !isValidURL(c.DodoCheckoutURL) {
		return fmt.Errorf("DODO_CHECKOUT_URL is invalid: %s", c.DodoCheckoutURL)
	}
	if len(c.KVNamespace) < 4 {
		return fmt.Errorf("KV_NAMESPACE is too short")
	}
	if len(c.DodoWebhookSecret) < 8 {
		return fmt.Errorf("DODO_WEBHOOK_SECRET is too short")
	}
	if len(c.DodoAPIKey) < 8 {
		return fmt.Errorf("DODO_API_KEY is too short")
	}
	return nil
}

func isValidURL(url string) bool {
	return len(url) > 0 && strings.HasPrefix(url, "https://")
}
