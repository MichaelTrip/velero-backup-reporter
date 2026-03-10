package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Kubeconfig         string        `mapstructure:"kubeconfig"`
	Namespace          string        `mapstructure:"namespace"`
	Port               int           `mapstructure:"port"`
	CollectionInterval time.Duration `mapstructure:"collection-interval"`

	SMTP  SMTPConfig `mapstructure:",squash"`
	Email EmailConfig
}

type SMTPConfig struct {
	Host     string `mapstructure:"smtp-host"`
	Port     int    `mapstructure:"smtp-port"`
	Username string `mapstructure:"smtp-username"`
	Password string `mapstructure:"smtp-password"`
	From     string `mapstructure:"smtp-from"`
	To       []string `mapstructure:"smtp-to"`
	TLS      bool   `mapstructure:"smtp-tls"`
}

type EmailConfig struct {
	Enabled     bool   `mapstructure:"email-enabled"`
	Schedule    string `mapstructure:"email-schedule"`
	TestEnabled bool   `mapstructure:"email-test-enabled"`
}

func Load() (*Config, error) {
	v := viper.GetViper()

	// Environment variable bindings
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Load config file if specified
	configFile := v.GetString("config")
	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}

	cfg := &Config{}

	cfg.Kubeconfig = v.GetString("kubeconfig")
	cfg.Namespace = v.GetString("namespace")
	if cfg.Namespace == "" {
		cfg.Namespace = "velero"
	}
	cfg.Port = v.GetInt("port")
	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	intervalStr := v.GetString("collection-interval")
	if intervalStr == "" {
		intervalStr = "5m"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("parsing collection-interval: %w", err)
	}
	cfg.CollectionInterval = interval

	cfg.SMTP.Host = v.GetString("smtp-host")
	cfg.SMTP.Port = v.GetInt("smtp-port")
	if cfg.SMTP.Port == 0 {
		cfg.SMTP.Port = 587
	}
	cfg.SMTP.Username = v.GetString("smtp-username")
	cfg.SMTP.Password = v.GetString("smtp-password")
	cfg.SMTP.From = v.GetString("smtp-from")
	cfg.SMTP.To = v.GetStringSlice("smtp-to")
	cfg.SMTP.TLS = v.GetBool("smtp-tls")

	cfg.Email.Enabled = v.GetBool("email-enabled")
	cfg.Email.Schedule = v.GetString("email-schedule")
	if cfg.Email.Schedule == "" {
		cfg.Email.Schedule = "0 8 * * *"
	}
	cfg.Email.TestEnabled = v.GetBool("email-test-enabled")

	return cfg, nil
}
