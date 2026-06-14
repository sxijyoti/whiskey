package build

import (
	"path/filepath"

	"github.com/sxijyoti/whiskey/internal/fingerprint"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/planner"
)

func IncrementalBuild(
	root string,
) error {

	g, err := graph.BuildSiteGraph(
		filepath.Join(
			root,
			"content",
		),
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

	if len(store) == 0 {

		if err := BuildSite(
			root,
		); err != nil {
			return err
		}

		return fingerprint.Save(
			".whiskey/fingerprints.json",
			store,
		)
	}

	changed, err := fingerprint.
		ChangedSources(
			g,
			store,
		)

	if err != nil {
		return err
	}

	if len(changed) == 0 {

		return fingerprint.Save(
			".whiskey/fingerprints.json",
			store,
		)
	}

	dirty := planner.DirtyPages(
		g,
		changed,
	)

	for _, page := range dirty {

		println(
			"[incremental]",
			page,
		)
	}

	return fingerprint.Save(
		".whiskey/fingerprints.json",
		store,
	)
}