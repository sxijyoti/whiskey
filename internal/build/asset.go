package build

import (
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/graph"
)

func CopyAsset(
	root string,
	asset string,
) error {

	themeRoot := filepath.Join(
		"themes",
		"default",
		"static",
	)

	siteRoot := filepath.Join(
		root,
		"static",
	)

	var rel string

	if strings.HasPrefix(
		asset,
		themeRoot,
	) {

		r, err := filepath.Rel(
			themeRoot,
			asset,
		)

		if err != nil {
			return err
		}

		rel = r

	} else {

		r, err := filepath.Rel(
			siteRoot,
			asset,
		)

		if err != nil {
			return err
		}

		rel = r
	}

	return copyFile(
		asset,
		filepath.Join(
			"dist",
			rel,
		),
	)
}

func DirtyAssets(
	g *graph.Graph,
	changed []string,
) []string {

	var assets []string

	for _, node := range changed {

		n := g.Nodes[node]

		if n == nil {
			continue
		}

		if n.Type != graph.AssetNode {
			continue
		}

		assets = append(
			assets,
			node,
		)
	}

	return assets
}