package build

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/fingerprint"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/planner"
	"github.com/sxijyoti/whiskey/internal/source"
)

// IncrementalBuild runs the dependency-aware incremental build pipeline:
//
//  1. Load config and fingerprint store.
//  2. Build dependency graph (all pages, sources, layouts, partials, assets).
//  3. Materialize remote sources into the workspace; garbage-collect stale ones.
//  4. Fingerprint every graph node → detect changed nodes.
//  5. If config, layout, or partial changed → full rebuild.
//  6. Otherwise → graph propagation → rebuild only dirty publishable pages.
//  7. Rebuild secondary indexes (collections, tags, RSS, sitemap).
//  8. Save fingerprints.
func IncrementalBuild(root string) error {

	if err := EnsureWorkspace(root); err != nil {
		return err
	}

	contentRoot := filepath.Join(root, "content")

	// --- Config ----------------------------------------------------------

	cfg, err := config.Load(root)
	if err != nil {
		return err
	}

	// --- Fingerprint store -----------------------------------------------

	store, err := fingerprint.Load(fingerprintPath(root))
	if err != nil {
		return err
	}

	// --- Dependency graph ------------------------------------------------

	g, err := graph.BuildSiteGraph(root, cfg.Theme)
	if err != nil {
		return err
	}

	// --- Navigation index (needed by rendering) --------------------------

	allPages, err := DiscoverPages(contentRoot)
	if err != nil {
		return err
	}

	index, err := BuildIndex(root, allPages)
	if err != nil {
		return err
	}

	cfg.Nav = BuildNav(cfg, index)

	// --- Materialize remote sources --------------------------------------

	manifest, err := source.LoadManifest(root)
	if err != nil {
		return err
	}

	materialized, err := MaterializeSources(root, g, manifest)
	if err != nil {
		return err
	}

	// --- Garbage-collect stale workspace entries -------------------------

	if err := GarbageCollectWorkspace(root, g, manifest); err != nil {
		return err
	}

	// --- Detect changed nodes via unified fingerprinting -----------------

	changedNodes, err := fingerprint.ChangedNodes(root, g, store)
	if err != nil {
		return err
	}

	// --- Clean up deleted pages/nodes ------------------------------------
	for id := range store {
		if strings.HasPrefix(id, "__") {
			continue
		}
		if g.Nodes[id] == nil {
			if strings.HasPrefix(id, contentRoot) && filepath.Ext(id) == ".md" {
				if output, err := pageOutputPath(root, contentRoot, id); err == nil {
					_ = os.Remove(output)
					parent := filepath.Dir(output)
					if parent != filepath.Join(root, "dist") {
						_ = os.Remove(parent)
					}
				}
			}
			delete(store, id)
		}
	}

	// --- Config change? --------------------------------------------------

	configHash := fingerprint.ConfigHash(cfg)
	if store["__config__"].Hash != configHash {
		store["__config__"] = fingerprint.Entry{Hash: configHash}
		return fullBuildAndSave(root, "config changed")
	}

	// --- Layout or partial changed? → full rebuild -----------------------

	for _, id := range changedNodes {
		n := g.Nodes[id]
		if n == nil {
			continue
		}
		if n.Type == graph.LayoutNode || n.Type == graph.PartialNode {
			return fullBuildAndSave(root, rebuildReason(n.Type))
		}
	}

	// --- Compute dirty pages via graph propagation -----------------------

	dirtyPages := planner.DirtyPages(g, changedNodes)
	dirtyDraftPages := changedDraftPages(g, changedNodes)
	dirtyAssets := DirtyAssets(g, changedNodes)

	if len(dirtyPages) == 0 && len(dirtyDraftPages) == 0 && len(dirtyAssets) == 0 {
		if _, err := os.Stat(filepath.Join(root, "dist")); os.IsNotExist(err) {
			return fullBuildAndSave(root, "missing output")
		}
		if !LogNoopBuilds {
			return fingerprint.Save(fingerprintPath(root), store)
		}
		start := logBuildStart("Incremental build", "no changes")
		logSources(nil, materialized)
		logDirtyPages(0)
		logDirtyAssets(0)
		logBuildDone("Incremental", start, 0, materialized)
		return fingerprint.Save(fingerprintPath(root), store)
	}

	localChanged := changedLocalSources(root, g, changedNodes)
	start := logBuildStart("Incremental build", incrementalReason(g, changedNodes, materialized))
	logDirtyPages(len(dirtyPages))
	logDirtyAssets(len(dirtyAssets))
	logSources(localChanged, materialized)

	// --- Copy dirty assets -----------------------------------------------

	for _, asset := range dirtyAssets {
		if err := CopyAsset(root, cfg.Theme, asset); err != nil {
			return err
		}
	}

	// --- Rebuild dirty pages ---------------------------------------------

	failedPages := make(map[string]bool)

	if len(dirtyPages) > 0 {
		for _, page := range dirtyPages {
			logRenderPage(contentRoot, page)
			rebuildOnePage(root, cfg, contentRoot, page, g, materialized, failedPages)
		}
	}

	for _, page := range dirtyDraftPages {
		rebuildOnePage(root, cfg, contentRoot, page, g, materialized, failedPages)
	}

	// --- Rebuild secondary indexes ---------------------------------------
	//
	// Always regenerate after any incremental rebuild. Correctness first.

	if len(dirtyPages) > 0 || len(dirtyDraftPages) > 0 || len(dirtyAssets) > 0 {
		if err := rebuildIndexes(root, cfg, contentRoot, allPages, failedPages); err != nil {
			return err
		}
	}

	// --- Save fingerprints -----------------------------------------------

	if err := fingerprint.Save(fingerprintPath(root), store); err != nil {
		return err
	}

	if len(failedPages) > 0 {
		return fmt.Errorf("incremental build completed with %d failure(s)", len(failedPages))
	}

	logBuildDone("Incremental", start, len(dirtyPages), materialized)

	return nil
}

// fullBuildAndSave runs a complete site build then persists the fingerprint store.
func fullBuildAndSave(root, reason string) error {
	return BuildSiteWithReason(root, reason)
}

func rebuildReason(nodeType graph.NodeType) string {
	switch nodeType {
	case graph.LayoutNode:
		return "layout changed"
	case graph.PartialNode:
		return "partial changed"
	default:
		return "site changed"
	}
}

func incrementalReason(g *graph.Graph, changed []string, materialized *MaterializationResult) string {
	if materialized != nil && len(materialized.Updated) > 0 {
		return "remote source updated"
	}

	for _, id := range changed {
		n := g.Nodes[id]
		if n == nil {
			continue
		}
		switch n.Type {
		case graph.PageNode:
			return "page changed"
		case graph.AssetNode:
			return "asset changed"
		case graph.SourceNode:
			if strings.HasPrefix(id, "local:") {
				return "page changed"
			}
			return "remote source updated"
		}
	}

	return "site changed"
}

// changedLocalSources returns a sorted slice of site-root-relative paths for
// every local SourceNode in changedNodes. Local source IDs have the form
// "local:/absolute/path"; the prefix is stripped and the path is made relative
// to root so the output is legible (e.g. "content/getting-started.md").
func changedLocalSources(root string, g *graph.Graph, changed []string) []string {
	var paths []string
	for _, id := range changed {
		n := g.Nodes[id]
		if n == nil || n.Type != graph.SourceNode {
			continue
		}
		if !strings.HasPrefix(id, "local:") {
			continue
		}
		abs := strings.TrimPrefix(id, "local:")
		rel, err := filepath.Rel(root, abs)
		if err != nil {
			// Fall back to the absolute path if Rel fails.
			rel = abs
		}
		paths = append(paths, filepath.ToSlash(rel))
	}
	sort.Strings(paths)
	return paths
}

func changedDraftPages(g *graph.Graph, changed []string) []string {
	var pages []string
	seen := map[string]bool{}

	for _, id := range changed {
		node := g.Nodes[id]
		if node == nil || node.Type != graph.PageNode || !node.Draft || seen[id] {
			continue
		}
		seen[id] = true
		pages = append(pages, id)
	}

	return pages
}

// rebuildOnePage reads, parses, and renders a single page.
// On failure it removes any stale output and records the slug in failedPages.
// Errors are logged but not returned so the caller can continue with
// the remaining pages.
func rebuildOnePage(
	root string,
	cfg *config.Config,
	contentRoot string,
	page string,
	g *graph.Graph,
	materialized *MaterializationResult,
	failedPages map[string]bool,
) {

	raw, err := os.ReadFile(page)
	if err != nil {
		logPageError(contentRoot, page, err)
		failedPages[pageSlug(contentRoot, page)] = true
		return
	}

	doc, err := parser.ParseFrontmatter(raw)
	if err != nil {
		logPageError(contentRoot, page, err)
		failedPages[pageSlug(contentRoot, page)] = true
		return
	}

	// Draft pages are PageNodes in the graph (for dependency tracking) but
	// are never rendered.
	if doc.Meta.Draft {
		if output, err := pageOutputPath(root, contentRoot, page); err == nil {
			_ = os.Remove(output)
			parent := filepath.Dir(output)
			if parent != filepath.Join(root, "dist") {
				_ = os.Remove(parent)
			}
		}
		return
	}

	output, err := pageOutputPath(root, contentRoot, page)
	if err != nil {
		logPageError(contentRoot, page, err)
		failedPages[pageSlug(contentRoot, page)] = true
		return
	}

	slug := pageSlug(contentRoot, page)

	// If any required remote source failed to materialize, skip this page.
	for _, dep := range g.Dependencies(page) {
		if materialized.OfflineCached[dep] {
			continue
		}
		if _, failed := materialized.Failed[dep]; failed {
			_ = os.Remove(output)
			logPageError(contentRoot, page, fmt.Errorf("workspace missing for %s", dep))
			failedPages[slug] = true
			return
		}
	}

	if err := BuildPage(root, cfg, page, doc, output); err != nil {
		_ = os.Remove(output)
		logPageError(contentRoot, page, err)
		failedPages[slug] = true
	}
}

// rebuildIndexes regenerates collections, tags, RSS, and sitemap from the set
// of all known pages, excluding any that failed during this build.
func rebuildIndexes(
	root string,
	cfg *config.Config,
	contentRoot string,
	allPages []string,
	failedPages map[string]bool,
) error {

	// Restrict to pages that have a valid output on disk and did not fail.
	var verified []string

	for _, p := range allPages {
		slug := pageSlug(contentRoot, p)
		if failedPages[slug] {
			continue
		}
		output, err := pageOutputPath(root, contentRoot, p)
		if err != nil {
			continue
		}
		if _, err := os.Stat(output); err == nil {
			verified = append(verified, p)
		}
	}

	cleanIndex, err := BuildIndex(root, verified)
	if err != nil {
		return err
	}

	cfg.Nav = BuildNav(cfg, cleanIndex)

	for _, page := range verified {
		raw, err := os.ReadFile(page)
		if err != nil {
			return err
		}

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return err
		}

		output, err := pageOutputPath(root, contentRoot, page)
		if err != nil {
			return err
		}

		if err := BuildPage(root, cfg, page, doc, output); err != nil {
			return err
		}
	}

	if err := cleanupSecondaryIndexes(root, allPages); err != nil {
		return err
	}

	if err := BuildCollections(root, cfg, cleanIndex); err != nil {
		return err
	}

	if err := BuildTags(root, cfg, cleanIndex); err != nil {
		return err
	}

	if err := BuildRSS(root, cfg, cleanIndex); err != nil {
		return err
	}

	return BuildSitemap(root, cfg, cleanIndex)
}

func cleanupSecondaryIndexes(root string, allPages []string) error {
	_ = os.RemoveAll(filepath.Join(root, "dist", "tags"))

	collections := map[string]bool{}
	for _, page := range allPages {
		raw, err := os.ReadFile(page)
		if err != nil {
			continue
		}
		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			continue
		}
		if doc.Meta.Collection != "" {
			collections[doc.Meta.Collection] = true
		}
	}

	for collection := range collections {
		path := filepath.Join(root, "dist", collection, "index.html")
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// pageSlug derives the URL slug from a page's absolute path.
func pageSlug(contentRoot, page string) string {
	rel, _ := filepath.Rel(contentRoot, page)
	return strings.TrimSuffix(rel, filepath.Ext(rel))
}

// pageOutputPath derives the dist output path for a content page.
func pageOutputPath(root, contentRoot, page string) (string, error) {
	rel, err := filepath.Rel(contentRoot, page)
	if err != nil {
		return "", err
	}
	slug := strings.TrimSuffix(rel, filepath.Ext(rel))
	if slug == "index" {
		return filepath.Join(root, "dist", "index.html"), nil
	}
	return filepath.Join(root, "dist", slug, "index.html"), nil
}

// fingerprintPath returns the canonical location of the fingerprint store.
func fingerprintPath(root string) string {
	return filepath.Join(root, ".whiskey", "fingerprints.json")
}
