package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sxijyoti/whiskey/internal/dependency"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestRegression_NoFullRebuildOnLocalInclude(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("inc.md", "---\ndraft: true\n---\nContent")
	site.WritePage("a.md", "---\ntitle: A\n---\n@include local:content/inc.md")
	site.WritePage("b.md", "---\ntitle: B\n---\nBody B")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	// Edit local include
	site.EditPage("inc.md", "---\ndraft: true\n---\nNew Content")

	res := helpers.BuildCmd(t, site.Dir)
	// Must be incremental, NOT full
	helpers.AssertIncrementalBuild(t, res)
	helpers.AssertOutputContains(t, res, "1 dirty page")
}

func TestRegression_NoFullRebuildOnRemoteInclude(t *testing.T) {
	site := setupMinimalSite(t)
	site.WriteRemoteWorkspace("https://example.com/source.md", "Content")
	site.WritePage("a.md", "---\ntitle: A\n---\n@include https://example.com/source.md")
	site.WritePage("b.md", "---\ntitle: B\n---\nBody B")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	// Edit remote cache
	site.WriteRemoteWorkspace("https://example.com/source.md", "New Content")
	site.Build()

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res)
	helpers.AssertOutputContains(t, res, "1 dirty page")
}

func TestRegression_MarkdownEditOnlyRebuildsAffected(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("a.md", "---\ntitle: A\n---\nBody A")
	site.WritePage("b.md", "---\ntitle: B\n---\nBody B")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	site.EditPage("a.md", "---\ntitle: A\n---\nNew Body A")

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res)
	helpers.AssertOutputContains(t, res, "1 dirty page")
}

func TestRegression_UnchangedBuildNothingChanged(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("a.md", "---\ntitle: A\n---\nBody A")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertNothingChanged(t, res)
}

func TestRegression_ConfigChangeFullRebuild(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("a.md", "---\ntitle: A\n---\nBody A")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	// Edit config
	site.WriteConfig(`theme="minimal"
title="Updated Test Site"`)
	site.Build()

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertFullBuild(t, res, "config changed")
}

func TestRegression_LayoutChangeFullRebuild(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("a.md", "---\ntitle: A\n---\nBody A")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	// Edit layout page.html
	site.WriteLayout("minimal", "page.html", `{{ define "content" }}Layout override: {{ .Content }}{{ end }}`)
	site.Build()

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertFullBuild(t, res, "layout changed")
}

func TestRegression_PartialChangeFullRebuild(t *testing.T) {
	site := setupMinimalSite(t)
	site.WriteLayout("minimal", "base.html", `<!DOCTYPE html><html><body>{{ block "content" . }}{{ end }}{{ template "footer" . }}</body></html>`)
	site.WritePartial("minimal", "footer.html", `{{ define "footer" }}<footer>Footer</footer>{{ end }}`)
	site.WritePage("a.md", "---\ntitle: A\n---\nBody A")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))

	// Edit partial
	site.WritePartial("minimal", "footer.html", `{{ define "footer" }}<footer>New Footer</footer>{{ end }}`)
	site.Build()

	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertFullBuild(t, res, "partial changed")
}

func TestRegression_RelativePathsAndCanonicalIncludes(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"
favicon="icons/favicon.ico"`)
	site.WriteLayout("minimal", "base.html", `<!DOCTYPE html><html><head>{{ template "head" . }}</head><body><nav>{{ range .Site.Nav }}<a href="{{ .URL }}">{{ .Title }}</a>{{ end }}</nav>{{ block "content" . }}{{ end }}</body></html>`)
	site.WritePartial("minimal", "head.html", `{{ define "head" }}{{ if .Site.Favicon }}<link rel="icon" href="{{ .Site.Favicon }}">{{ end }}{{ end }}`)
	site.WriteAsset("", "icons/favicon.ico", "ico")
	site.WritePage("index.md", "---\ntitle: Home\n---\n@include local:getting-started.md")
	site.WritePage("getting-started.md", "---\ndraft: true\n---\nGetting started")
	site.WritePage("blog/post.md", "---\ntitle: Post\ncollection: blog\ndate: 2026-07-01\ntags: [go]\n---\n@include local:snippets/intro.md")
	site.WritePage("blog/snippets/intro.md", "---\ndraft: true\n---\nNested intro")
	site.Build()

	res := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	helpers.AssertContains(t, site.Dir, "dist/index.html", "Getting started")
	helpers.AssertContains(t, site.Dir, "dist/blog/post/index.html", "Nested intro")
	helpers.AssertContains(t, site.Dir, "dist/index.html", `href="/icons/favicon.ico"`)
	helpers.AssertContains(t, site.Dir, "dist/blog/post/index.html", `href="/icons/favicon.ico"`)
	helpers.AssertFileExists(t, site.Dir, "dist/icons/favicon.ico")

	g, err := graph.BuildSiteGraph(site.Dir, "minimal")
	if err != nil {
		t.Fatalf("failed to build graph: %v", err)
	}

	expected, err := dependency.ResolveLocalRef(
		filepath.Join(site.Dir, "content", "blog", "post.md"),
		"local:snippets/../snippets/intro.md",
	)
	if err != nil {
		t.Fatalf("failed to resolve include: %v", err)
	}

	actual, err := dependency.ResolveLocalRef(
		filepath.Join(site.Dir, "content", "blog", "post.md"),
		"local:snippets/intro.md",
	)
	if err != nil {
		t.Fatalf("failed to resolve include: %v", err)
	}

	if actual != expected {
		t.Fatalf("expected canonical paths to match:\nactual:   %s\nexpected: %s", actual, expected)
	}

	if g.Nodes[actual] == nil {
		t.Fatalf("expected graph to contain canonical include node %s", actual)
	}
}

func TestRegression_FaviconUsesBaseURLPath(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"
base_url="https://example.com/docs"
favicon="images/favicon.ico"`)
	site.WriteLayout("minimal", "base.html", `<!DOCTYPE html><html><head>{{ template "head" . }}</head><body>{{ block "content" . }}{{ end }}</body></html>`)
	site.WritePartial("minimal", "head.html", `{{ define "head" }}<link rel="icon" href="{{ .Site.Favicon }}">{{ end }}`)
	site.WriteAsset("", "images/favicon.ico", "ico")
	site.WritePage("index.md", "---\ntitle: Home\n---\nHome")
	site.Build()

	res := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	helpers.AssertContains(t, site.Dir, "dist/index.html", `href="/docs/images/favicon.ico"`)
	helpers.AssertFileExists(t, site.Dir, "dist/images/favicon.ico")
}

func TestRegression_FaviconRenderedOnEveryPageType(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"
favicon="images/favicon.ico"
[rss]
enabled=true
collections=["blog"]`)
	site.WriteLayout("minimal", "base.html", `<!DOCTYPE html><html><head>{{ template "head" . }}</head><body>{{ block "content" . }}{{ end }}</body></html>`)
	site.WritePartial("minimal", "head.html", `{{ define "head" }}{{ if .Site.Favicon }}<link rel="icon" href="{{ .Site.Favicon }}">{{ end }}{{ end }}`)
	site.WriteLayout("minimal", "tags.html", `{{ define "content" }}<h1>Tags</h1>{{ range .Tags }}<a href="/tags/{{ . }}/">{{ . }}</a>{{ end }}{{ end }}`)
	site.WriteAsset("", "images/favicon.ico", "ico")
	site.WritePage("index.md", "---\ntitle: Home\n---\nHome")
	site.WritePage("about.md", "---\ntitle: About\n---\nAbout")
	site.WritePage("blog/post.md", "---\ntitle: Post\ncollection: blog\ndate: 2026-07-01\ntags: [go]\n---\nPost")
	site.Build()

	res := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	for _, path := range []string{
		"dist/index.html",
		"dist/about/index.html",
		"dist/blog/post/index.html",
		"dist/blog/index.html",
		"dist/tags/index.html",
		"dist/tags/go/index.html",
	} {
		helpers.AssertContains(t, site.Dir, path, `href="/images/favicon.ico"`)
	}
}

func TestRegression_RepositoryReadmeIncludeIsNotCycle(t *testing.T) {
	site := setupMinimalSite(t)
	site.WriteFile("README.md", "# Project README\n\n@include local:README.md\n")
	site.WritePage("docs/readme.md", "---\ntitle: Docs\n---\n@include local:../../README.md")
	site.Build()

	res := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	helpers.AssertContains(t, site.Dir, "dist/docs/readme/index.html", "Project README")
	helpers.AssertContains(t, site.Dir, "dist/docs/readme/index.html", "@include local:README.md")
}

func TestRegression_RealCyclicIncludesStillFail(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("a.md", "---\ntitle: A\n---\n@include local:b.md")
	site.WritePage("b.md", "---\ntitle: B\n---\n@include local:a.md")
	site.Build()

	res := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildFails(t, res)
	helpers.AssertOutputContains(t, res, "Cycle detected")
}

func TestRegression_GitHubBlobURLNormalizesToRaw(t *testing.T) {
	site := setupMinimalSite(t)
	blob := "https://github.com/sxijyoti/whiskey/blob/main/README.md"
	raw := "https://raw.githubusercontent.com/sxijyoti/whiskey/main/README.md"
	site.WriteRemoteWorkspace(raw, "Normalized README")
	site.WritePage("index.md", "---\ntitle: Home\n---\n@include "+blob)
	site.Build()

	res := helpers.BuildOfflineCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)
	helpers.AssertContains(t, site.Dir, "dist/index.html", "Normalized README")
}

func TestRegression_RawGitHubURLStillWorks(t *testing.T) {
	site := setupMinimalSite(t)
	raw := "https://raw.githubusercontent.com/sxijyoti/whiskey/main/README.md"
	site.WriteRemoteWorkspace(raw, "Raw README")
	site.WritePage("index.md", "---\ntitle: Home\n---\n@include "+raw)
	site.Build()

	res := helpers.BuildOfflineCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)
	helpers.AssertContains(t, site.Dir, "dist/index.html", "Raw README")
}

func TestRegression_NavigationIdenticalAcrossPageTypes(t *testing.T) {
	site := helpers.NewTestSite(t)
	site.Theme("minimal")
	site.WriteConfig(`theme="minimal"
[rss]
enabled=true
collections=["blog"]`)
	site.WriteLayout("minimal", "base.html", `<!DOCTYPE html><html><body><nav data-nav>{{ range .Site.Nav }}<a href="{{ .URL }}">{{ .Title }}</a>{{ end }}</nav>{{ block "content" . }}{{ end }}</body></html>`)
	site.WriteLayout("minimal", "tags.html", `{{ define "content" }}<h1>Tags</h1>{{ range .Tags }}<a href="/tags/{{ . }}/">{{ . }}</a>{{ end }}{{ end }}`)
	site.WritePage("index.md", "---\ntitle: Home\n---\nHome")
	site.WritePage("about.md", "---\ntitle: About\n---\nAbout")
	site.WritePage("blog/post.md", "---\ntitle: Post\ncollection: blog\ndate: 2026-07-01\ntags: [go]\n---\nPost")
	site.Build()

	res := helpers.BuildFullCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)

	paths := []string{
		"dist/index.html",
		"dist/about/index.html",
		"dist/blog/post/index.html",
		"dist/blog/index.html",
		"dist/tags/index.html",
		"dist/tags/go/index.html",
	}

	var nav string
	for _, path := range paths {
		html := readOutput(t, site.Dir, path)
		current := between(t, html, `<nav data-nav>`, `</nav>`)
		if nav == "" {
			nav = current
			continue
		}
		if current != nav {
			t.Fatalf("navigation mismatch in %s:\nexpected: %s\nactual:   %s", path, nav, current)
		}
	}
}

func TestRegression_DraftToggleRemovesAndRepublishesOutputs(t *testing.T) {
	site := setupMinimalSite(t)
	site.WriteConfig(`theme="minimal"
[rss]
enabled=true
collections=["blog"]`)
	site.WritePage("blog/post.md", "---\ntitle: Post\ncollection: blog\ndate: 2026-07-01\ntags: [go]\ndraft: false\n---\nPublished")
	site.Build()

	helpers.AssertBuildSucceeds(t, helpers.BuildFullCmd(t, site.Dir))
	helpers.AssertFileExists(t, site.Dir, "dist/blog/post/index.html")
	helpers.AssertContains(t, site.Dir, "dist/blog/index.html", "Post")
	helpers.AssertContains(t, site.Dir, "dist/tags/go/index.html", "Post")
	helpers.AssertContains(t, site.Dir, "dist/feed.xml", "Post")
	helpers.AssertContains(t, site.Dir, "dist/sitemap.xml", "/blog/post/")

	site.EditPage("blog/post.md", "---\ntitle: Post\ncollection: blog\ndate: 2026-07-01\ntags: [go]\ndraft: true\n---\nDraft")
	res := helpers.BuildCmd(t, site.Dir)
	helpers.AssertIncrementalBuild(t, res)
	helpers.AssertFileNotExists(t, site.Dir, "dist/blog/post/index.html")
	helpers.AssertFileNotExists(t, site.Dir, "dist/blog/index.html")
	helpers.AssertNotContains(t, site.Dir, "dist/tags/index.html", "go")
	helpers.AssertFileNotExists(t, site.Dir, "dist/tags/go/index.html")
	helpers.AssertNotContains(t, site.Dir, "dist/feed.xml", "Post")
	helpers.AssertNotContains(t, site.Dir, "dist/sitemap.xml", "/blog/post/")

	site.EditPage("blog/post.md", "---\ntitle: Post\ncollection: blog\ndate: 2026-07-01\ntags: [go]\ndraft: false\n---\nRepublished")
	res = helpers.BuildCmd(t, site.Dir)
	helpers.AssertBuildSucceeds(t, res)
	helpers.AssertFileExists(t, site.Dir, "dist/blog/post/index.html")
	helpers.AssertContains(t, site.Dir, "dist/blog/post/index.html", "Republished")
	helpers.AssertContains(t, site.Dir, "dist/blog/index.html", "Post")
	helpers.AssertContains(t, site.Dir, "dist/tags/go/index.html", "Post")
	helpers.AssertContains(t, site.Dir, "dist/feed.xml", "Post")
	helpers.AssertContains(t, site.Dir, "dist/sitemap.xml", "/blog/post/")
}

func readOutput(t *testing.T, dir, path string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, path))
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return string(raw)
}

func between(t *testing.T, body, start, end string) string {
	t.Helper()
	startAt := strings.Index(body, start)
	if startAt < 0 {
		t.Fatalf("missing start marker %q in:\n%s", start, body)
	}
	startAt += len(start)
	endAt := strings.Index(body[startAt:], end)
	if endAt < 0 {
		t.Fatalf("missing end marker %q in:\n%s", end, body)
	}
	return body[startAt : startAt+endAt]
}

func TestRegression_GraphJSONUsesCanonicalLocalPaths(t *testing.T) {
	site := setupMinimalSite(t)
	site.WritePage("blog/snippets/intro.md", "---\ndraft: true\n---\nIntro")
	site.WritePage("blog/post.md", "---\ntitle: Post\n---\n@include local:snippets/../snippets/intro.md")
	site.Build()

	g, err := graph.BuildSiteGraph(site.Dir, "minimal")
	if err != nil {
		t.Fatalf("failed to build graph: %v", err)
	}

	data, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("failed to marshal graph: %v", err)
	}

	canonical, err := dependency.ResolveLocalRef(
		filepath.Join(site.Dir, "content", "blog", "post.md"),
		"local:snippets/intro.md",
	)
	if err != nil {
		t.Fatalf("failed to resolve canonical path: %v", err)
	}

	if !strings.Contains(string(data), canonical) {
		t.Fatalf("expected graph JSON to contain canonical path %s\n%s", canonical, string(data))
	}
	if strings.Contains(string(data), "snippets/../snippets") {
		t.Fatalf("expected graph JSON to avoid uncanonical include path\n%s", string(data))
	}
}
