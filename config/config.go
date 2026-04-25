package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                  string
	DatabaseURL           string
	AnthropicAPIKey       string
	SupabaseURL           string
	SupabaseServiceRoleKey string
	StripeSecretKey       string
	StripePriceID         string
}

func Load() (*Config, error) {
	required := []string{
		"PORT",
		"DATABASE_URL",
		"SUPABASE_URL",
		"SUPABASE_SERVICE_ROLE_KEY",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return &Config{
		Port:                  os.Getenv("PORT"),
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		AnthropicAPIKey:       os.Getenv("ANTHROPIC_API_KEY"),
		SupabaseURL:           os.Getenv("SUPABASE_URL"),
		SupabaseServiceRoleKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		StripeSecretKey:       os.Getenv("STRIPE_SECRET_KEY"),
		StripePriceID:         os.Getenv("STRIPE_PRICE_ID"),
	}, nil
}
