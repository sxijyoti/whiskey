package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/graph"
)

func CopyAsset(
	root string,
	theme string,
	asset string,
) error {

	themeRoot := filepath.Join(
		"themes",
		theme,
		"static",
	)

	siteRoot := filepath.Join(
		root,
		"static",
	)

	var rel string

	if strings.HasPrefix(
		filepath.Base(asset),
		".",
	) {
		return nil
	}

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

		siteOverride := filepath.Join(
			siteRoot,
			rel,
		)

		if _, err := os.Stat(siteOverride); err == nil {
			asset = siteOverride
		}

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
			root,
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

		if strings.HasPrefix(
			filepath.Base(node),
			".",
		) {
			continue
		}

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
