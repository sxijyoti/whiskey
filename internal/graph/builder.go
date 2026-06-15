package graph

import (
	"io/fs"
	"path/filepath"
	"os"

	"github.com/sxijyoti/whiskey/internal/dependency"
	"github.com/sxijyoti/whiskey/internal/parser"
)

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

	return g, nil
}