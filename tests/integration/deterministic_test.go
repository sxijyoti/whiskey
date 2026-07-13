package integration

import (
	"path/filepath"
	"testing"

	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestDeterministicOutput_DoubleBuild(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"`)

	// Generate a sufficiently complex site to ensure hash stability
	site.WritePage("a.md", "---\ntitle: A\n---\nContent A")
	site.WritePage("b.md", "---\ntitle: B\n---\nContent B")
	site.WritePage("c.md", "---\ntitle: C\ntags: [go, static]\n---\nContent C")
	site.Build()

	// First build
	res1 := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res1)
	distDir := filepath.Join(site.Dir, "dist")
	hash1 := helpers.HashDir(t, distDir)

	// Wipe output and rebuild without changing inputs
	helpers.RunWhiskey(t, site.Dir, "clean", ".")

	res2 := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res2)
	hash2 := helpers.HashDir(t, distDir)

	if hash1 != hash2 {
		t.Fatalf("Deterministic build failed!\nHash 1: %s\nHash 2: %s", hash1, hash2)
	}
}
