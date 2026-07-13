package integration

import (
	"testing"

	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestOffline_BuildWithCache(t *testing.T) {
	site := setupMinimalSite(t)
	// Write cached workspace entry for remote include
	site.WriteRemoteWorkspace("https://example.com/cache.md", "Cached Content")
	site.WritePage("index.md", "---\ntitle: Home\n---\n@include https://example.com/cache.md")
	site.Build()

	// Run offline build
	res := helpers.BuildOfflineCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)
	helpers.AssertContains(t, site.Dir, "dist/index.html", "Cached Content")
}

func TestOffline_BuildWithoutCacheFails(t *testing.T) {
	site := setupMinimalSite(t)
	// No cache entry created
	site.WritePage("index.md", "---\ntitle: Home\n---\n@include https://example.com/nocache.md")
	site.Build()

	// Run offline build - should fail
	res := helpers.BuildOfflineCmd(t, site.Dir)
	helpers.AssertBuildFails(t, res)
	helpers.AssertOutputContains(t, res, "missing https://example.com/nocache.md")
}
