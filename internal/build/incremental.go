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
)

func IncrementalBuild(
	root string,
) error {

	contentRoot := filepath.Join(
		root,
		"content",
	)

	cfg, err := config.Load(root)
	if err != nil {
		return err
	}

	store, err := fingerprint.Load(
		".whiskey/fingerprints.json",
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

	changedSources, err := fingerprint.ChangedSources(
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
			".whiskey/fingerprints.json",
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

		if _, err := os.Stat("dist"); os.IsNotExist(err) {

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
				".whiskey/fingerprints.json",
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
	}

	return fingerprint.Save(
		".whiskey/fingerprints.json",
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

	rel, err := filepath.Rel(
		contentRoot,
		page,
	)
	if err != nil {
		return err
	}

	slug := strings.TrimSuffix(
		rel,
		filepath.Ext(rel),
	)

	var output string

	if slug == "index" {

		output = filepath.Join(
			"dist",
			"index.html",
		)

	} else {

		output = filepath.Join(
			"dist",
			slug,
			"index.html",
		)
	}

	return BuildPage(
		root,
		cfg,
		doc,
		output,
	)
}