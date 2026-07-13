package integration

import (
	"testing"

	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestFailure_MissingLocalInclude(t *testing.T) {
	site := setupMinimalSite(t)
	// index.md includes a non-existent local file
	site.WritePage("index.md", "---\ntitle: Home\n---\n@include local:content/nonexistent.md")
	site.WritePage("other.md", "---\ntitle: Other\n---\nGood page")
	site.Build()

	res := helpers.BuildFullCmd(t, site.Dir)
	// Build should fail overall
	helpers.AssertBuildFails(t, res)
	helpers.AssertOutputContains(t, res, "no such file or directory")

	// The bad page output must not exist or must be removed
	helpers.AssertFileNotExists(t, site.Dir, "dist/index.html")

	// Remaining page should still be built successfully
	helpers.AssertFileExists(t, site.Dir, "dist/other/index.html")
}

func TestFailure_FailedPageRebuildRemovesStaleOutput(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("index.md", "---\ntitle: Home\n---\nHello")
	site.Build()

	// Initial successful build
	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))
	helpers.AssertFileExists(t, site.Dir, "dist/index.html")

	// Make the page fail by introducing an invalid include
	site.EditPage("index.md", "---\ntitle: Home\n---\n@include local:content/missing.md")

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildFails(t, res)

	// Output must be removed so we don't leave stale HTML behind
	helpers.AssertFileNotExists(t, site.Dir, "dist/index.html")
}

func TestFailure_ExcludedFromSecondaryIndexes(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("good.md", "---\ntitle: Good Blog\ncollection: blog\ndate: 2026-07-01\n---\nGood Body")
	site.WritePage("bad.md", "---\ntitle: Bad Blog\ncollection: blog\ndate: 2026-07-02\n---\n@include local:content/missing.md")
	site.Build()

	res := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildFails(t, res)

	// RSS, sitemap, collections should still build for successful pages
	helpers.AssertFileExists(t, site.Dir, "dist/feed.xml")
	helpers.AssertFileExists(t, site.Dir, "dist/sitemap.xml")

	// Only "Good Blog" should be in secondary indexes
	helpers.AssertContains(t, site.Dir, "dist/feed.xml", "Good Blog")
	helpers.AssertNotContains(t, site.Dir, "dist/feed.xml", "Bad Blog")

	helpers.AssertContains(t, site.Dir, "dist/sitemap.xml", "/good/")
	helpers.AssertNotContains(t, site.Dir, "dist/sitemap.xml", "/bad/")
}
