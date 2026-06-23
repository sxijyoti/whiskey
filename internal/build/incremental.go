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

	g, err := graph.BuildSiteGraph(
		contentRoot,
	)
	if err != nil {
		return err
	}

	store, err := fingerprint.Load(
		".whiskey/fingerprints.json",
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

	localDirty, err := planner.LocalDirtyPages(
		contentRoot,
		store,
	)
	if err != nil {
		return err
	}

	dirty := planner.IncrementalDirtySet(
		g,
		localDirty,
		changedSources,
	)

	if len(dirty) == 0 {

		if _, err := os.Stat("dist"); os.IsNotExist(err) {
			return BuildSite(root)
		}

		fmt.Println("[build] nothing changed")
		return nil
	}

	fmt.Printf(
		"[build] %d dirty page(s)\n",
		len(dirty),
	)

	for _, page := range dirty {

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

	if err := fingerprint.Save(
		".whiskey/fingerprints.json",
		store,
	); err != nil {
		return err
	}


	if err := CopyStatic(root); err != nil {
		return err
	}

	return nil
}

func rebuildPage(
	root string,
	cfg *config.Config,
	contentRoot string,
	page string,
) error {

	raw, err := os.ReadFile(page)
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

	if doc.Meta.Draft {

		fmt.Printf(
			"[skip] draft %s\n",
			page,
		)

		return nil
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