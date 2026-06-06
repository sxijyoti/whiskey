package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Title       string `toml:"title"`
	Description string `toml:"description"`
	BaseURL     string `toml:"base_url"`
	Theme       string `toml:"theme"`
}

func Load(siteRoot string) (*Config, error) {

	cfg := Default()

	path := filepath.Join(
		siteRoot,
		"whiskey.toml",
	)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(
		path,
		cfg,
	); err != nil {
		return nil, err
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func Default() *Config {

	return &Config{
		Title:       "Whiskey Site",
		Description: "",
		BaseURL:     "",
		Theme:       "default",
	}
}

func Validate(
	cfg *Config,
) error {

	if cfg.Theme == "" {
		return fmt.Errorf(
			"theme cannot be empty",
		)
	}

	return nil
}
