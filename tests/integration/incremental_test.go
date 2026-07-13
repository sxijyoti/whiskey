package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sxijyoti/whiskey/tests/helpers"
)

func setupMinimalSite(t *testing.T) *helpers.SiteBuilder {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"
title="Test Site"
description="Description"
base_url="http://localhost"
[rss]
enabled=true
collections=["blog"]`)
	return site
}

func TestIncremental_MarkdownBodyEdit(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("a.md", "---\ntitle: Page A\n---\nContent A")
	site.WritePage("b.md", "---\ntitle: Page B\n---\nContent B")
	site.Build()

	// Initial full build
	res1 := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res1)

	// Mutate body of a.md
	site.EditPage("a.md", "---\ntitle: Page A\n---\nUpdated Body A")

	// Incremental build
	res2 := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res2)
	helpers.AssertOutputContains(t, res2, "1 dirty page")

	helpers.AssertContains(t, site.Dir, "dist/a/index.html", "Updated Body A")
	helpers.AssertContains(t, site.Dir, "dist/b/index.html", "Content B")
}

func TestIncremental_FrontmatterEdit(t *testing.T) {
	site := setupMinimalSite(t)
	site.WriteLayout("minimal", "page.html", `{{ define "content" }}<h1>{{ .Title }}</h1>{{ .Content }}{{ end }}`)
	site.WritePage("a.md", "---\ntitle: Title A\n---\nContent A")
	site.Build()

	res1 := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res1)

	// Edit frontmatter title
	site.EditPage("a.md", "---\ntitle: Title A Updated\n---\nContent A")

	res2 := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res2)
	helpers.AssertContains(t, site.Dir, "dist/a/index.html", "Title A Updated")
}

func TestIncremental_LocalIncludeEdit(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("inc.md", "---\ndraft: true\n---\nInclude Content")
	site.WritePage("index.md", "---\ntitle: Index\n---\n@include local:content/inc.md")
	site.Build()

	res1 := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res1)
	helpers.AssertContains(t, site.Dir, "dist/index.html", "Include Content")

	// Edit local include file
	site.EditPage("inc.md", "---\ndraft: true\n---\nUpdated Include Content")

	res2 := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res2)
	helpers.AssertContains(t, site.Dir, "dist/index.html", "Updated Include Content")
}

func TestIncremental_RemoteIncludeEdit(t *testing.T) {
	site := setupMinimalSite(t)
	site.WriteRemoteWorkspace("https://example.com/remote.md", "Remote Content")
	site.WritePage("index.md", "---\ntitle: Index\n---\n@include https://example.com/remote.md")
	site.Build()

	res1 := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res1)
	helpers.AssertContains(t, site.Dir, "dist/index.html", "Remote Content")

	// Update simulated remote workspace file content
	site.WriteRemoteWorkspace("https://example.com/remote.md", "Updated Remote Content")
	site.Build()

	res2 := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res2)
	helpers.AssertContains(t, site.Dir, "dist/index.html", "Updated Remote Content")
}

func TestIncremental_AssetEdit(t *testing.T) {
	site := setupMinimalSite(t)
	site.WriteAsset("", "style.css", "body { color: red; }")
	site.WritePage("index.md", "---\ntitle: Home\n---\nHello")
	site.Build()

	res1 := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res1)

	// Update asset
	site.WriteAsset("", "style.css", "body { color: blue; }")
	site.Build()

	res2 := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res2)
	helpers.AssertOutputContains(t, res2, "Dirty assets: 1")
	helpers.AssertOutputContains(t, res2, "Dirty pages: 0")
	helpers.AssertContains(t, site.Dir, "dist/style.css", "blue")
}

func TestIncremental_MultipleDependentPages(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("shared.md", "---\ndraft: true\n---\nShared content")
	site.WritePage("a.md", "---\ntitle: A\n---\n@include local:content/shared.md")
	site.WritePage("b.md", "---\ntitle: B\n---\n@include local:content/shared.md")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	site.EditPage("shared.md", "---\ndraft: true\n---\nNew Shared Content")

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res)
	helpers.AssertOutputContains(t, res, "2 dirty page(s)")
	helpers.AssertContains(t, site.Dir, "dist/a/index.html", "New Shared Content")
	helpers.AssertContains(t, site.Dir, "dist/b/index.html", "New Shared Content")
}

func TestIncremental_NestedIncludes(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("inner.md", "---\ndraft: true\n---\nInner Leaf")
	site.WritePage("mid.md", "---\ndraft: true\n---\n@include local:content/inner.md")
	site.WritePage("outer.md", "---\ntitle: Outer\n---\n@include local:content/mid.md")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	site.EditPage("inner.md", "---\ndraft: true\n---\nUpdated Leaf")

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res)
	helpers.AssertContains(t, site.Dir, "dist/outer/index.html", "Updated Leaf")
}

func TestIncremental_GarbageCollection(t *testing.T) {
	site := setupMinimalSite(t)
	site.WriteRemoteWorkspace("https://example.com/gc.md", "Gc Content")
	site.WritePage("index.md", "---\ntitle: Index\n---\n@include https://example.com/gc.md")
	site.Build()

	// Initial build materializes and caches gc.md
	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))
	helpers.AssertFileExists(t, site.Dir, ".whiskey/workspace")

	// Verify workspace file exists
	entries, _ := os.ReadDir(filepath.Join(site.Dir, ".whiskey", "workspace"))
	if len(entries) == 0 {
		t.Fatal("Expected cached workspace files")
	}

	// Remove dependency on remote source
	site.WritePage("index.md", "---\ntitle: Index\n---\nNo include")
	site.Build()

	// Run incremental
	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res)
	helpers.AssertOutputContains(t, res, "[gc] removed")

	// Workspace should now be empty (since no pages reference any remote includes)
	entries2, _ := os.ReadDir(filepath.Join(site.Dir, ".whiskey", "workspace"))
	if len(entries2) > 0 {
		t.Fatalf("Expected workspace to be cleaned up, still has %d files", len(entries2))
	}
}

func TestIncremental_SecondaryIndexesRegeneration(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("a.md", "---\ntitle: Blog 1\ncollection: blog\ntags: [tag1]\ndate: 2026-07-01\n---\nBody 1")
	site.WritePage("b.md", "---\ntitle: Blog 2\ncollection: blog\ntags: [tag2]\ndate: 2026-07-02\n---\nBody 2")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))
	helpers.AssertFileExists(t, site.Dir, "dist/feed.xml")
	helpers.AssertFileExists(t, site.Dir, "dist/sitemap.xml")
	helpers.AssertFileExists(t, site.Dir, "dist/tags/tag1/index.html")

	// Edit b.md title
	site.EditPage("b.md", "---\ntitle: Blog 2 Updated\ncollection: blog\ntags: [tag2, tag3]\ndate: 2026-07-02\n---\nBody 2")

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res)

	// Verify index updates
	helpers.AssertContains(t, site.Dir, "dist/feed.xml", "Blog 2 Updated")
	helpers.AssertContains(t, site.Dir, "dist/sitemap.xml", "/b/")
	helpers.AssertFileExists(t, site.Dir, "dist/tags/tag3/index.html")
}
