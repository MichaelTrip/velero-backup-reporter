package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func resetViper() {
	viper.Reset()
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := &Config{
		Namespace:          "velero",
		Port:               8080,
		CollectionInterval: 5 * time.Minute,
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"zero", 0},
		{"negative", -1},
		{"too high", 70000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:               tt.port,
				CollectionInterval: 5 * time.Minute,
			}
			if err := cfg.Validate(); err == nil {
				t.Fatal("expected error for invalid port")
			}
		})
	}
}

func TestValidate_EmailDisabledWhenSMTPIncomplete(t *testing.T) {
	cfg := &Config{
		Port:               8080,
		CollectionInterval: 5 * time.Minute,
		Email: EmailConfig{
			Enabled:  true,
			Schedule: "0 8 * * *",
		},
		SMTP: SMTPConfig{
			Host: "",
			From: "test@example.com",
			To:   []string{"admin@example.com"},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Email.Enabled {
		t.Fatal("expected email to be disabled when smtp-host is empty")
	}
}

func TestValidate_EmailEnabledWithCompleteSMTP(t *testing.T) {
	cfg := &Config{
		Port:               8080,
		CollectionInterval: 5 * time.Minute,
		Email: EmailConfig{
			Enabled:  true,
			Schedule: "0 8 * * *",
		},
		SMTP: SMTPConfig{
			Host: "smtp.example.com",
			Port: 587,
			From: "test@example.com",
			To:   []string{"admin@example.com"},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !cfg.Email.Enabled {
		t.Fatal("expected email to remain enabled with complete SMTP config")
	}
}

func TestLoad_Defaults(t *testing.T) {
	resetViper()
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Namespace != "velero" {
		t.Errorf("expected namespace 'velero', got '%s'", cfg.Namespace)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.CollectionInterval != 5*time.Minute {
		t.Errorf("expected 5m interval, got %v", cfg.CollectionInterval)
	}
	if cfg.Email.Schedule != "0 8 * * *" {
		t.Errorf("expected default email schedule, got '%s'", cfg.Email.Schedule)
	}
}

func TestLoad_FromConfigFile(t *testing.T) {
	resetViper()

	content := []byte("namespace: monitoring\nport: 9090\n")
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	viper.Set("config", cfgFile)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Namespace != "monitoring" {
		t.Errorf("expected namespace 'monitoring', got '%s'", cfg.Namespace)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
}
