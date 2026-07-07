package build

import (
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"fmt"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/dependency"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/template"
	"github.com/sxijyoti/whiskey/internal/source"
	"github.com/sxijyoti/whiskey/internal/graph"
)

func BuildPage(
	siteRoot string,
	cfg *config.Config,
	doc *parser.Document,
	output string,
) error {

	resolvedBody, err := dependency.ResolveIncludes(
		doc.Body,
		func(ref string) ([]byte, error) {

			if !source.WorkspaceExists(
				siteRoot,
				ref,
			) {

				return nil, fmt.Errorf(
					"workspace missing for %s",
					ref,
				)
			}

			return source.ReadWorkspace(
				siteRoot,
				ref,
			)
		},
	)
	if err != nil {
		return err
	}

	resolvedBody, err = parser.ExpandShortcodes(
		resolvedBody,
	)
	if err != nil {
		return err
	}

	html, err := parser.MdToHTML(
		resolvedBody,
	)
	if err != nil {
		return err
	}

	layout := doc.Meta.Layout
	if layout == "" {
		layout = "page"
	}

	var date string

	if !doc.Meta.Date.IsZero() {
		date = doc.Meta.Date.Format(
			"2006-01-02",
		)
	}

	page, err := template.RenderPage(
		siteRoot,
		cfg.Theme,
		layout,
		template.PageData{
			Site:        cfg,
			Title:       doc.Meta.Title,
			Description: doc.Meta.Description,
			Date:        date,
			Tags:        doc.Meta.Tags,
			Content:     htmltemplate.HTML(html),
		},
	)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(
		filepath.Dir(output),
		0755,
	); err != nil {
		return err
	}

	return os.WriteFile(
		output,
		page,
		0644,
	)
}

func BuildSite(
	root string,
) error {

	if err := EnsureWorkspace(
		root,
	); err != nil {
		return err
	}

	cfg, err := config.Load(
		root,
	)
	if err != nil {
		return err
	}

	contentRoot := filepath.Join(
		root,
		"content",
	)

	pages, err := DiscoverPages(
		contentRoot,
	)
	if err != nil {
		return err
	}

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

	index, err := BuildIndex(
		root,
		pages,
	)
	if err != nil {
		return err
	}

	cfg.Nav = BuildNav(
		index,
	)

	for _, page := range pages {

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
			continue
		}

		output, err := pageOutputPath(
			root,
			contentRoot,
			page,
		)
		if err != nil {
			return err
		}

		skip := false

		for _, dep := range g.Dependencies(
			page,
		) {

			if materialized.OfflineCached[dep] {
				continue
			}

			if _, failed := materialized.Failed[dep]; failed {

				_ = os.Remove(
					output,
				)

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

		if err := BuildPage(
			root,
			cfg,
			doc,
			output,
		); err != nil {

			_ = os.Remove(
				output,
			)

			fmt.Printf(
				"[build] failed %s: %v\n",
				page,
				err,
			)

			continue
		}
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

	if err := BuildCollections(
		root,
		cfg,
		index,
	); err != nil {
		return err
	}

	if err := CopyStatic(
		root,
		cfg.Theme,
	); err != nil {
		return err
	}

	return nil
}