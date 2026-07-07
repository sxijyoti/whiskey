package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/fingerprint"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/planner"
	"github.com/sxijyoti/whiskey/internal/source"
)

func IncrementalBuild(
	root string,
) error {

	if err := EnsureWorkspace(
		root,
	); err != nil {
		return err
	}

	contentRoot := filepath.Join(
		root,
		"content",
	)

	cfg, err := config.Load(root)
	if err != nil {
		return err
	}

	store, err := fingerprint.Load(
		filepath.Join(
			root,
			".whiskey",
			"fingerprints.json",
		),
	)
	if err != nil {
		return err
	}

	allPages, err := DiscoverPages(
		contentRoot,
	)
	if err != nil {
		return err
	}

	index, err := BuildIndex(
		root,
		allPages,
	)
	if err != nil {
		return err
	}

	cfg.Nav = BuildNav(
		index,
	)

	g, err := graph.BuildSiteGraph(
		root,
		cfg.Theme,
	)
	if err != nil {
		return err
	}

	manifest, err := source.LoadManifest(
		root,
	)
	if err != nil {
		return err
	}

	materialized, err := MaterializeSources(
		root,
		g,
		manifest,
	)
	if err != nil {
		return err
	}

	changedSources, err := fingerprint.ChangedSources(
		root,
		g,
		store,
	)
	if err != nil {
		return err
	}

	fullBuild := false
	reason := ""

	configHash := fingerprint.ConfigHash(
		cfg,
	)

	if store["__config__"].Hash != configHash {

		store["__config__"] = fingerprint.Entry{
			Hash: configHash,
		}

		fullBuild = true
		reason = "config"
	}

	for _, source := range changedSources {

		node := g.Nodes[source]
		if node == nil {
			continue
		}

		switch node.Type {

		case graph.LayoutNode,
			graph.PartialNode:

			fullBuild = true

			if reason == "" {
				reason = "layout"
			}
		}
	}

	if fullBuild {

		switch reason {

		case "layout":

			fmt.Println("[build] full rebuild (layout changed)")

		case "config":

			fmt.Println("[build] full rebuild (config changed)")

		default:
			fmt.Println("[build] full rebuild")
		}

		if err := BuildSite(root); err != nil {
			return err
		}

		if err := fingerprint.UpdateLocalPages(
			contentRoot,
			store,
		); err != nil {
			return err
		}


		return fingerprint.Save(
			filepath.Join(
				root,
				".whiskey",
				"fingerprints.json",
			),
			store,
		)
	}

	dirtyAssets := DirtyAssets(
		g,
		changedSources,
	)

	localDirty, err := planner.LocalDirtyPages(
		contentRoot,
		store,
	)
	if err != nil {
		return err
	}

	var filteredDirty []string

	for _, page := range localDirty {

		raw, err := os.ReadFile(page)
		if err != nil {
			return err
		}

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return err
		}

		if doc.Meta.Draft {
			continue
		}

		filteredDirty = append(
			filteredDirty,
			page,
		)
	}

	dirty := planner.IncrementalDirtySet(
		g,
		filteredDirty,
		changedSources,
	)

	if len(dirty) == 0 &&
		len(dirtyAssets) == 0 {

		if _, err := os.Stat(filepath.Join(root, "dist")); os.IsNotExist(err) {

			if err := BuildSite(root); err != nil {
				return err
			}

			if err := fingerprint.UpdateLocalPages(
				contentRoot,
				store,
			); err != nil {
				return err
			}


			return fingerprint.Save(
				filepath.Join(
					root,
					".whiskey",
					"fingerprints.json",
				),
				store,
			)
		}

		fmt.Println("[build] nothing changed")

		return nil
	}

	fmt.Println("[build] incremental")

	for _, asset := range dirtyAssets {

		fmt.Printf(
			"[asset] %s\n",
			asset,
		)

		if err := CopyAsset(
			root,
			cfg.Theme,
			asset,
		); err != nil {
			return err
		}
	}

	var pagesToBuild []string

	for _, page := range dirty {

		raw, err := os.ReadFile(page)
		if err != nil {
			return err
		}

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return err
		}

		if doc.Meta.Draft {
			continue
		}

		skip := false
		for _, dep := range g.Dependencies(page) {
			if materialized.OfflineCached[dep] {
				continue
			}

			if _, failed := materialized.Failed[dep]; failed {
				output, err := pageOutputPath(
					root,
					contentRoot,
					page,
				)
				if err == nil {
					_ = os.Remove(output)
				}
				fmt.Printf(
					"[build] skipped %s (missing %s)\n",
					page,
					dep,
				)
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		pagesToBuild = append(
			pagesToBuild,
			page,
		)
	}

	if len(pagesToBuild) > 0 {

		fmt.Printf(
			"[build] %d dirty page(s)\n",
			len(pagesToBuild),
		)

		for _, page := range pagesToBuild {

			fmt.Printf(
				"[build] %s\n",
				page,
			)

			if err := rebuildPage(
				root,
				cfg,
				contentRoot,
				page,
			); err != nil {
				return err
			}
		}
	}

	if len(filteredDirty) > 0 {

		if err := BuildCollections(
			root,
			cfg,
			index,
		); err != nil {
			return err
		}

		if err := BuildTags(
			root,
			cfg,
			index,
		); err != nil {
			return err
		}

		if err := BuildRSS(
			root,
			cfg,
			index,
		); err != nil {
			return err
		}

		if err := BuildSitemap(
			root,
			cfg,
			index,
		); err != nil {
			return err
		}
	}

	return fingerprint.Save(
		filepath.Join(
			root,
			".whiskey",
			"fingerprints.json",
		),
		store,
	)
}

func rebuildPage(
	root string,
	cfg *config.Config,
	contentRoot string,
	page string,
) error {

	raw, err := os.ReadFile(
		page,
	)
	if err != nil {
		return err
	}

	doc, err := parser.ParseFrontmatter(
		raw,
	)
	if err != nil {
		return err
	}

	if doc.Meta.Draft {
		return nil
	}

	output, err := pageOutputPath(
		root,
		contentRoot,
		page,
	)
	if err != nil {
		return err
	}

	if err := BuildPage(
		root,
		cfg,
		doc,
		output,
	); err != nil {

		_ = os.Remove(output)

		return err
	}

	return nil
}

func pageOutputPath(
	root string,
	contentRoot string,
	page string,
) (string, error) {

	rel, err := filepath.Rel(
		contentRoot,
		page,
	)
	if err != nil {
		return "", err
	}

	slug := strings.TrimSuffix(
		rel,
		filepath.Ext(rel),
	)

	if slug == "index" {
		return filepath.Join(
			root,
			"dist",
			"index.html",
		), nil
	}

	return filepath.Join(
		root,
		"dist",
		slug,
		"index.html",
	), nil
}