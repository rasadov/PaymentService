package config

import (
	"fmt"
	"os"
	"strings"
)

var settings Config

type Config struct {
	Environment             string
	PaymentCallbackUrl      string
	DodoWebhookSecret       string
	DodoAPIKey              string
	DodoCheckoutURL         string
	DodoCheckoutRedirectUrl string
	KVNamespace             string
}

func LoadConfig() error {
	config := Config{
		Environment:             getEnvWithDefault("ENVIRONMENT", "development"),
		PaymentCallbackUrl:      getEnvRequired("PAYMENT_CALLBACK_URL"),
		DodoWebhookSecret:       getEnvRequired("DODO_WEBHOOK_SECRET"),
		DodoAPIKey:              getEnvRequired("DODO_API_KEY"),
		DodoCheckoutURL:         getEnvRequired("DODO_CHECKOUT_URL"),
		DodoCheckoutRedirectUrl: getEnvRequired("DODO_CHECKOUT_REDIRECT_URL"),
		KVNamespace:             getEnvRequired("KV_NAMESPACE"),
	}

	// Check for missing required environment variables
	var missingVars []string
	if config.PaymentCallbackUrl == "" {
		missingVars = append(missingVars, "PAYMENT_CALLBACK_URL")
	}
	if config.DodoWebhookSecret == "" {
		missingVars = append(missingVars, "DODO_WEBHOOK_SECRET")
	}
	if config.DodoAPIKey == "" {
		missingVars = append(missingVars, "DODO_API_KEY")
	}
	if config.DodoCheckoutURL == "" {
		missingVars = append(missingVars, "DODO_CHECKOUT_URL")
	}
	if config.DodoCheckoutRedirectUrl == "" {
		missingVars = append(missingVars, "DODO_CHECKOUT_REDIRECT_URL")
	}
	if config.KVNamespace == "" {
		missingVars = append(missingVars, "KV_NAMESPACE")
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	if err := config.validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	settings = config
	return nil
}

func GetConfig() Config {
	return settings
}

func getEnvRequired(key string) string {
	return os.Getenv(key)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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
