package build

import (
	"html"
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/template"
)

type CollectionPage struct {
	Title       string
	URL         string
	Description string
	Date        time.Time
}

func BuildCollections(
	root string,
	cfg *config.Config,
	pages []string,
) error {

	collections := map[string][]CollectionPage{}

	contentRoot := filepath.Join(
		root,
		"content",
	)

	for _, page := range pages {

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

		if doc.Meta.Collection == "" {
			continue
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

		collections[doc.Meta.Collection] = append(
			collections[doc.Meta.Collection],
			CollectionPage{
				Title:       doc.Meta.Title,
				URL:         "/" + slug + "/",
				Description: doc.Meta.Description,
				Date:        doc.Meta.Date,
			},
		)
	}

	for name, entries := range collections {

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Date.After(
				entries[j].Date,
			)
		})

		var content strings.Builder

		content.WriteString(
			`<div class="collection-list">`,
		)

		for _, p := range entries {

			content.WriteString(
				`<a class="collection-entry" href="`,
			)

			content.WriteString(
				html.EscapeString(
					p.URL,
				),
			)

			content.WriteString(`">`)

			content.WriteString(
				`<time class="collection-date">`,
			)

			content.WriteString(
				p.Date.Format(
					"2006-01-02",
				),
			)

			content.WriteString(`</time>`)

			content.WriteString(
				`<span class="collection-title">`,
			)

			content.WriteString(
				html.EscapeString(
					p.Title,
				),
			)

			content.WriteString(`</span>`)

			content.WriteString(`</a>`)
		}

		content.WriteString(
			`</div>`,
		)

		page, err := template.RenderPage(
			root,
			cfg.Theme,
			"list",
			template.PageData{
				Site:        cfg,
				Title:       DisplayName(name),
				Description: "",
				Content:     htmltemplate.HTML(content.String()),
			},
		)
		if err != nil {
			return err
		}

		out := filepath.Join(
			"dist",
			name,
			"index.html",
		)

		if err := os.MkdirAll(
			filepath.Dir(out),
			0755,
		); err != nil {
			return err
		}

		if err := os.WriteFile(
			out,
			page,
			0644,
		); err != nil {
			return err
		}
	}

	return nil
}