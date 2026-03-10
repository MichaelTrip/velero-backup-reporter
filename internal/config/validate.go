package config

import (
	"fmt"
	"log"
)

func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}

	if c.CollectionInterval <= 0 {
		return fmt.Errorf("collection-interval must be positive")
	}

	if c.Email.Enabled {
		if c.SMTP.Host == "" {
			log.Println("WARN: email enabled but smtp-host not set, disabling email")
			c.Email.Enabled = false
		}
		if c.SMTP.From == "" {
			log.Println("WARN: email enabled but smtp-from not set, disabling email")
			c.Email.Enabled = false
		}
		if len(c.SMTP.To) == 0 {
			log.Println("WARN: email enabled but smtp-to not set, disabling email")
			c.Email.Enabled = false
		}
	}

	return nil
}
