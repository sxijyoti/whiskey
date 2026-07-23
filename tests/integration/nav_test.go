package integration

import (
	"testing"

	"github.com/sxijyoti/whiskey/internal/build"
	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestNav_ExplicitOverAuto(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("portfolio")
	site.WriteConfig(`
title = "Portfolio"
theme = "portfolio"

[[nav]]
title = "Home"
url = "/"

[[nav]]
title = "Ink & Light"
url = "/ink-and-light/"
`)
	site.WritePage("writeups/leaves.md", "---\ntitle: Leaves\ncollection: writeups\n---\nContent")
	site.Build()

	cfg, err := config.Load(site.Dir)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	nav := build.BuildNav(cfg, nil)
	if len(nav) != 2 {
		t.Fatalf("Expected 2 nav items from explicit config, got %d: %+v", len(nav), nav)
	}
	if nav[1].Title != "Ink & Light" {
		t.Errorf("Expected nav[1] title to be 'Ink & Light', got %q", nav[1].Title)
	}
}
