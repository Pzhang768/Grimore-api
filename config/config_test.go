package config

import (
	"os"
	"testing"
)

var allVars = map[string]string{
	"PORT":                      "8080",
	"DATABASE_URL":              "postgres://localhost/test",
	"SUPABASE_URL":              "https://test.supabase.co",
	"SUPABASE_SERVICE_ROLE_KEY": "test-role-key",
}

func setEnv(vars map[string]string) {
	for k, v := range vars {
		os.Setenv(k, v)
	}
}

func unsetEnv(vars map[string]string) {
	for k := range vars {
		os.Unsetenv(k)
	}
}

func TestLoad_Success(t *testing.T) {
	setEnv(allVars)
	defer unsetEnv(allVars)

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

func TestLoad_MissingVar(t *testing.T) {
	required := []string{
		"PORT",
		"DATABASE_URL",
		"SUPABASE_URL",
		"SUPABASE_SERVICE_ROLE_KEY",
	}

	for _, missing := range required {
		t.Run(missing, func(t *testing.T) {
			setEnv(allVars)
			defer unsetEnv(allVars)
			os.Unsetenv(missing)

			_, err := Load()
			if err == nil {
				t.Fatalf("expected error for missing %s, got nil", missing)
			}
		})
	}
}
