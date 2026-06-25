package build

import (
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/dependency"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/template"
)

func BuildPage(
	siteRoot string,
	cfg *config.Config,
	doc *parser.Document,
	output string,
) error {

	resolvedBody, err := dependency.ResolveIncludes(
		doc.Body,
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

func BuildSite(root string) error {

	cfg, err := config.Load(root)
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

	nav, err := BuildNav(
		root,
		pages,
	)
	if err != nil {
		return err
	}
	cfg.Nav = nav

	for _, page := range pages {

		rel, err := filepath.Rel(
			contentRoot,
			page,
		)
		if err != nil {
			return err
		}

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

		if err := BuildPage(
			root,
			cfg,
			doc,
			output,
		); err != nil {
			return err
		}
	}

	if err := BuildTags(
		root,
		cfg,
		pages,
	); err != nil {
		return err
	}

	if err := BuildRSS(
		root,
		cfg,
		pages,
	); err != nil {
		return err
	}

	if err := BuildCollections(
		root,
		cfg,
		pages,
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
