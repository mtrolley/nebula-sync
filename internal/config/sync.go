package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Sync struct {
	FullSync        bool    `required:"true" envconfig:"FULL_SYNC"`
	Cron            *string `envconfig:"CRON"`
	RunGravity      bool    `default:"false" envconfig:"RUN_GRAVITY"`
	GravitySettings *GravitySettings
	ConfigSettings  *ConfigSettings `ignored:"true"`
}

func (c *Config) loadSync() error {
	sync := Sync{}
	if err := envconfig.Process("", &sync); err != nil {
		return fmt.Errorf("sync env vars: %w", err)
	}

	if err := sync.loadConfigSettings(); err != nil {
		return fmt.Errorf("load config settings: %w", err)
	}

	c.Sync = &sync
	return nil
}
