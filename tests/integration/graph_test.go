package integration

import (
	"testing"

	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestDependencyGraph_CycleDetection(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"`)

	// Create a cycle using Whiskey's @include syntax
	site.WritePage("a.md", "---\ntitle: A\n---\n@include local:content/b.md")
	site.WritePage("b.md", "---\ntitle: B\n---\n@include local:content/a.md")
	site.Build()

	res := helpers.BuildCmd(t, site.Dir)

	helpers.AssertBuildFails(t, res)
	helpers.AssertOutputContains(t, res, "Cycle detected")
}

func TestDependencyGraph_OrphanGarbageCollection(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"`)

	// Use markdown files and @include syntax
	site.WritePage("widget.md", "Widget Content")
	site.WritePage("a.md", "---\ntitle: Page A\n---\n@include local:content/widget.md")
	site.WritePage("b.md", "---\ntitle: Page B\n---\nIndependent page")
	site.Build()

	// 1. Initial full build
	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	helpers.AssertFileExists(t, site.Dir, "dist/a/index.html")
	helpers.AssertContains(t, site.Dir, "dist/a/index.html", "Widget Content")

	// 2. Delete the page
	site.DeletePage("a.md")

	// 3. Rebuild incrementally
	resIncremental := helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, resIncremental)

	// 4. Assert the orphan is removed
	helpers.AssertFileNotExists(t, site.Dir, "dist/a/index.html")
	helpers.AssertFileExists(t, site.Dir, "dist/b/index.html")
}
