package graph

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sxijyoti/whiskey/internal/dependency"
	"github.com/sxijyoti/whiskey/internal/parser"
)

func resolveThemeFile(siteRoot, theme, rel string) string {
	siteFile := filepath.Join(
		siteRoot,
		rel,
	)

	if _, err := os.Stat(siteFile); err == nil {
		return siteFile
	}

	return filepath.Join(
		"themes",
		theme,
		rel,
	)
}

func resolvePartials(siteRoot, theme string) ([]string, error) {
	partials := map[string]string{}

	themePartials, err := filepath.Glob(
		filepath.Join(
			"themes",
			theme,
			"layouts",
			"partials",
			"*.html",
		),
	)
	if err != nil {
		return nil, err
	}

	for _, partial := range themePartials {
		partials[filepath.Base(partial)] = partial
	}

	sitePartials, err := filepath.Glob(
		filepath.Join(
			siteRoot,
			"layouts",
			"partials",
			"*.html",
		),
	)
	if err != nil {
		return nil, err
	}

	for _, partial := range sitePartials {
		partials[filepath.Base(partial)] = partial
	}

	names := make([]string, 0, len(partials))
	for name := range partials {
		names = append(names, name)
	}

	sort.Strings(names)

	files := make([]string, 0, len(names))
	for _, name := range names {
		files = append(files, partials[name])
	}

	return files, nil
}

func addTemplateGraph(
	g *Graph,
	siteRoot string,
	theme string,
	pagePath string,
	layout string,
) error {

	layoutPath := resolveThemeFile(
		siteRoot,
		theme,
		filepath.Join(
			"layouts",
			layout+".html",
		),
	)

	basePath := resolveThemeFile(
		siteRoot,
		theme,
		filepath.Join(
			"layouts",
			"base.html",
		),
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

	partials, err := resolvePartials(
		siteRoot,
		theme,
	)
	if err != nil {
		return err
	}

	for _, partial := range partials {

		g.AddNode(
			partial,
			PartialNode,
		)

		g.AddEdge(
			basePath,
			partial,
		)
	}

	return nil
}

func BuildSiteGraph(
	siteRoot string,
	theme string,
) (*Graph, error) {

	g := New()

	contentRoot := filepath.Join(
		siteRoot,
		"content",
	)

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

			if strings.HasPrefix(
				filepath.Base(path),
				".",
			) {
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

			if err := addTemplateGraph(
				g,
				siteRoot,
				theme,
				path,
				layout,
			); err != nil {
				return err
			}

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
			siteRoot,
			"static",
		),
		filepath.Join(
			"themes",
			theme,
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

				if strings.HasPrefix(
					filepath.Base(path),
					".",
				) {
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
