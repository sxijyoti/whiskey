package integration

import (
	"testing"

	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestBuild_SinglePage(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"`)
	site.WritePage("hello.md", "---\ntitle: Hello\n---\n# Hello World")
	site.Build()

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	helpers.AssertFileExists(t, site.Dir, "dist/hello/index.html")
	helpers.AssertContains(t, site.Dir, "dist/hello/index.html", "<h1")
	helpers.AssertContains(t, site.Dir, "dist/hello/index.html", "Hello World")
}

func TestBuild_DraftsSkipped(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"`)
	site.WritePage("published.md", "---\ntitle: Pub\ndraft: false\n---\nPub")
	site.WritePage("draft.md", "---\ntitle: Draft\ndraft: true\n---\nDraft")
	site.Build()

	helpers.BuildCmd(t, site.Dir)

	helpers.AssertFileExists(t, site.Dir, "dist/published/index.html")
	helpers.AssertFileNotExists(t, site.Dir, "dist/draft/index.html")
}