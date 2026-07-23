package integration

import (
	"testing"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestConfig_ExplicitNav(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("portfolio")
	site.WriteConfig(`
title = "Portfolio"
theme = "portfolio"

[[nav]]
title = "Home"
url = "/"

[[nav]]
title = "Projects"
url = "/projects/"
`)
	site.Build()

	cfg, err := config.Load(site.Dir)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if len(cfg.ExplicitNav) != 2 {
		t.Fatalf("Expected 2 explicit nav items, got %d", len(cfg.ExplicitNav))
	}
	if cfg.ExplicitNav[0].Title != "Home" || cfg.ExplicitNav[1].Title != "Projects" {
		t.Errorf("Unexpected explicit nav items: %+v", cfg.ExplicitNav)
	}
}
