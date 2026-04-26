package config

import (
	"strings"
	"testing"
)

var allVars = map[string]string{
	"PORT":                      "8080",
	"DATABASE_URL":              "postgres://localhost/test",
	"SUPABASE_URL":              "https://test.supabase.co",
	"SUPABASE_SERVICE_ROLE_KEY": "test-role-key",
	"SUPABASE_JWT_SECRET":       "test-jwt-secret",
}

func TestLoad_Success(t *testing.T) {
	for k, v := range allVars {
		t.Setenv(k, v)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected Port 8080, got %s", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://localhost/test" {
		t.Errorf("unexpected DatabaseURL: %s", cfg.DatabaseURL)
	}
	if cfg.SupabaseURL != "https://test.supabase.co" {
		t.Errorf("unexpected SupabaseURL: %s", cfg.SupabaseURL)
	}
	if cfg.SupabaseServiceRoleKey != "test-role-key" {
		t.Errorf("unexpected SupabaseServiceRoleKey: %s", cfg.SupabaseServiceRoleKey)
	}
}

func TestLoad_OptionalFieldsEmpty(t *testing.T) {
	for k, v := range allVars {
		t.Setenv(k, v)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.AnthropicAPIKey != "" {
		t.Errorf("expected AnthropicAPIKey empty, got %s", cfg.AnthropicAPIKey)
	}
	if cfg.StripeSecretKey != "" {
		t.Errorf("expected StripeSecretKey empty, got %s", cfg.StripeSecretKey)
	}
	if cfg.StripePriceID != "" {
		t.Errorf("expected StripePriceID empty, got %s", cfg.StripePriceID)
	}
}

func TestLoad_MissingVar(t *testing.T) {
	required := []string{
		"PORT",
		"DATABASE_URL",
		"SUPABASE_URL",
		"SUPABASE_SERVICE_ROLE_KEY",
		"SUPABASE_JWT_SECRET",
	}

	for _, missing := range required {
		t.Run(missing, func(t *testing.T) {
			for k, v := range allVars {
				t.Setenv(k, v)
			}
			t.Setenv(missing, "")

			_, err := Load()
			if err == nil {
				t.Fatalf("expected error for missing %s, got nil", missing)
			}
			if !strings.Contains(err.Error(), missing) {
				t.Errorf("expected error to mention %q, got: %v", missing, err)
			}
		})
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	for k, v := range allVars {
		t.Setenv(k, v)
	}
	t.Setenv("PORT", "notaport")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid PORT, got nil")
	}
	if !strings.Contains(err.Error(), "PORT") {
		t.Errorf("expected error to mention PORT, got: %v", err)
	}
}
