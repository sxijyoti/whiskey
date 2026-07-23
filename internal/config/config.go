package config

import (
	"fmt"
	"net/url"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type RSSConfig struct {
	Enabled     bool     `toml:"enabled"`
	Collections []string `toml:"collections"`
}

type NavItem struct {
	Title string `toml:"title"`
	URL   string `toml:"url"`
}

type Config struct {
	Title       string    `toml:"title"`
	Description string    `toml:"description"`
	BaseURL     string    `toml:"base_url"`
	Theme       string    `toml:"theme"`
	Favicon     string    `toml:"favicon"`
	ExplicitNav []NavItem `toml:"nav"`
	Nav         []NavItem
	RSS         RSSConfig `toml:"rss"`
	RSSURL      string
}

func Load(siteRoot string) (*Config, error) {

	cfg := Default()

	path := filepath.Join(
		siteRoot,
		"whiskey.toml",
	)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		normalizePaths(cfg)
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

	normalizePaths(cfg)

	if cfg.RSS.Enabled &&
		len(cfg.RSS.Collections) > 0 {
		cfg.RSSURL = "/feed.xml"
	}

	return cfg, nil
}

func Default() *Config {

	return &Config{
		Title:       "Whiskey Site",
		Description: "",
		BaseURL:     "",
		Theme:       "minimal",
		Favicon:     "images/favicon.ico",

		RSS: RSSConfig{
			Enabled: true,
			Collections: []string{
				"blog",
			},
		},
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

func normalizePaths(cfg *Config) {
	if cfg.Favicon == "" {
		return
	}

	favicon := strings.ReplaceAll(cfg.Favicon, "\\", "/")
	if strings.HasPrefix(favicon, "http://") ||
		strings.HasPrefix(favicon, "https://") {
		cfg.Favicon = favicon
		return
	}

	publicPath := strings.TrimPrefix(pathpkg.Clean(favicon), "/")
	basePath := ""

	if cfg.BaseURL != "" {
		if parsed, err := url.Parse(cfg.BaseURL); err == nil {
			basePath = parsed.Path
		}
	}

	cfg.Favicon = pathpkg.Join("/", basePath, publicPath)
}
