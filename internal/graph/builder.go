package graph

import (
	"io/fs"
	"path/filepath"
	"os"

	"github.com/sxijyoti/whiskey/internal/dependency"
	"github.com/sxijyoti/whiskey/internal/parser"
)

func addTemplateGraph(
	g *Graph,
	pagePath string,
	layout string,
) {

	layoutPath := filepath.Join(
		"themes",
		"default",
		"layouts",
		layout+".html",
	)

	basePath := filepath.Join(
		"themes",
		"default",
		"layouts",
		"base.html",
	)

	g.AddNode(
		layoutPath,
		LayoutNode,
	)

	g.AddNode(
		basePath,
		LayoutNode,
	)

	g.AddEdge(
		pagePath,
		layoutPath,
	)

	g.AddEdge(
		layoutPath,
		basePath,
	)

	partials := []string{
		"head.html",
		"header.html",
		"footer.html",
	}

	for _, partial := range partials {

		p := filepath.Join(
			"themes",
			"default",
			"layouts",
			"partials",
			partial,
		)

		g.AddNode(
			p,
			PartialNode,
		)

		g.AddEdge(
			basePath,
			p,
		)
	}
}

func BuildSiteGraph(
	contentRoot string,
) (*Graph, error) {

	g := New()

	err := filepath.WalkDir(
		contentRoot,
		func(
			path string,
			d fs.DirEntry,
			err error,
		) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".md" {
				return nil
			}

			raw, err := os.ReadFile(path)
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

			g.AddNode(
				path,
				PageNode,
			)

			layout := doc.Meta.Layout

			if layout == "" {
				layout = "page"
			}

			addTemplateGraph(
				g,
				path,
				layout,
			)

			directives := dependency.Extract(
				string(raw),
			)

			for _, dir := range directives {

				g.AddNode(
					dir.Ref,
					SourceNode,
				)

				g.AddEdge(
					path,
					dir.Ref,
				)
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	assetRoots := []string{
		filepath.Join(
			contentRoot,
			"..",
			"static",
		),
		filepath.Join(
			"themes",
			"default",
			"static",
		),
	}

	for _, root := range assetRoots {

		if _, err := os.Stat(root); os.IsNotExist(err) {
			continue
		}

		err := filepath.WalkDir(
			root,
			func(
				path string,
				d fs.DirEntry,
				err error,
			) error {

				if err != nil {
					return err
				}

				if d.IsDir() {
					return nil
				}

				g.AddNode(
					path,
					AssetNode,
				)

				return nil
			},
		)

		if err != nil {
			return nil, err
		}
	}

	return g, nil
}