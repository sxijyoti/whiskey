package integration

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sxijyoti/whiskey/internal/build"
	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/devserver"
	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestBuild_UnlistedAndCustom404(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`
title = "My Site"
theme = "minimal"

[rss]
enabled = true
collections = ["posts"]
`)

	site.WritePage("index.md", "---\ntitle: Home\n---\n# Welcome")
	site.WritePage("posts/post1.md", "---\ntitle: Post One\ncollection: posts\ntags: [news]\n---\nPost 1 content")
	site.WritePage("easter-egg.md", "---\ntitle: Secret Easter Egg\nunlisted: true\ntags: [secret]\n---\n# You found the secret!")
	site.WritePage("404.md", "---\ntitle: Page Not Found\n---\n# Custom 404 Page")

	site.Build()

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	// Check output files exist
	helpers.AssertFileExists(t, site.Dir, "dist/index.html")
	helpers.AssertFileExists(t, site.Dir, "dist/posts/post1/index.html")
	helpers.AssertFileExists(t, site.Dir, "dist/easter-egg/index.html")
	helpers.AssertFileExists(t, site.Dir, "dist/404.html")
	helpers.AssertFileExists(t, site.Dir, "dist/404/index.html")

	helpers.AssertContains(t, site.Dir, "dist/easter-egg/index.html", "You found the secret!")
	helpers.AssertContains(t, site.Dir, "dist/404.html", "Custom 404 Page")

	// Check Navigation
	cfg, err := config.Load(site.Dir)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	pages := []string{
		filepath.Join(site.Dir, "content", "index.md"),
		filepath.Join(site.Dir, "content", "posts", "post1.md"),
		filepath.Join(site.Dir, "content", "easter-egg.md"),
		filepath.Join(site.Dir, "content", "404.md"),
	}

	idx, err := build.BuildIndex(site.Dir, pages)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	nav := build.BuildNav(cfg, idx)
	for _, item := range nav {
		if item.Title == "Secret Easter Egg" || item.URL == "/easter-egg/" {
			t.Errorf("Unlisted easter egg should not appear in automatic nav, but found: %+v", item)
		}
		if item.Title == "Page Not Found" || item.URL == "/404/" || item.URL == "/404.html" {
			t.Errorf("404 page should not appear in automatic nav, but found: %+v", item)
		}
	}

	// Check Sitemap
	helpers.AssertFileExists(t, site.Dir, "dist/sitemap.xml")
	helpers.AssertNotContains(t, site.Dir, "dist/sitemap.xml", "/easter-egg/")
	helpers.AssertNotContains(t, site.Dir, "dist/sitemap.xml", "/404")

	// Check DevServer 404 handling
	handler := devserver.NewFileHandler(site.Dir, "dist")
	req := httptest.NewRequest("GET", "/non-existent-page-url", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 from devserver for missing route, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Custom 404 Page") {
		t.Errorf("Expected devserver to serve custom 404 content, got: %s", rec.Body.String())
	}
}

func TestBuild_ToggleUnlistedIncremental(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteLayout("minimal", "base.html", `<!DOCTYPE html><html><body><nav>{{ range .Site.Nav }}<a href="{{ .URL }}">{{ .Title }}</a>{{ end }}</nav>{{ block "content" . }}{{ end }}</body></html>`)
	site.WriteConfig(`
title = "My Site"
theme = "minimal"
`)

	site.WritePage("index.md", "---\ntitle: Home\n---\n# Welcome")
	site.WritePage("secret.md", "---\ntitle: Secret Page\nunlisted: true\n---\n# Secret")

	site.Build()
	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	// Initially secret.md is unlisted -> index.html navigation should NOT contain Secret Page
	helpers.AssertNotContains(t, site.Dir, "dist/index.html", "Secret Page")

	// Now toggle unlisted to false
	site.EditPage("secret.md", "---\ntitle: Secret Page\nunlisted: false\n---\n# Secret")
	res2 := helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res2)

	// index.html navigation SHOULD now contain Secret Page
	helpers.AssertContains(t, site.Dir, "dist/index.html", "Secret Page")

	// Now toggle unlisted back to true
	site.EditPage("secret.md", "---\ntitle: Secret Page\nunlisted: true\n---\n# Secret")
	res3 := helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res3)

	// index.html navigation SHOULD NO LONGER contain Secret Page
	helpers.AssertNotContains(t, site.Dir, "dist/index.html", "Secret Page")
}
